/**
 * RefreshManager - Centralized refresh coordination singleton
 * Prevents UI thread blocking through intelligent refresh batching and scheduling
 */

import * as vscode from 'vscode';
import { debounceRefresh, DebouncedFunction } from '../utils/debounce';
import { performanceMonitor, measureAsync } from '../utils/performance';
import { Logger } from '../utils/logger';

export enum RefreshReason {
    MANUAL = 'manual',
    INTERVAL = 'interval',
    CONFIG_CHANGE = 'config-change',
    DATA_CHANGE = 'data-change',
    COMMAND = 'command',
    FILTER_CHANGE = 'filter-change',
    PROJECT_CHANGE = 'project-change',
    TEMPLATE_CHANGE = 'template-change',
    STARTUP = 'startup',
    CONNECTION = 'connection',
    EXTERNAL_CHANGE = 'external-change' // WebSocket real-time updates
}

export enum RefreshPriority {
    LOW = 0,
    NORMAL = 1,
    HIGH = 2,
    CRITICAL = 3,
    REALTIME = 4 // WebSocket real-time updates (highest priority)
}

export interface RefreshRequest {
    id: string;
    reason: RefreshReason;
    priority: RefreshPriority;
    targets: RefreshTarget[];
    timestamp: number;
    metadata?: Record<string, any>;
}

export enum RefreshTarget {
    TASKS = 'tasks',
    PROJECTS = 'projects',
    TEMPLATES = 'templates',
    STATUS_BAR = 'status-bar',
    FILTER_TOOLBAR = 'filter-toolbar',
    ALL = 'all'
}

export interface RefreshProvider {
    refresh(): Promise<void>;
    getName(): string;
    supportsTarget(target: RefreshTarget): boolean;
}

export interface RefreshStats {
    totalRequests: number;
    processedRequests: number;
    deduplicatedRequests: number;
    averageProcessingTime: number;
    lastRefreshTime?: number;
    pendingRequests: number;
}

/**
 * Centralized refresh management system
 * Coordinates all refresh operations to prevent UI blocking
 */
export class RefreshManager {
    private static instance: RefreshManager;
    private providers: Map<RefreshTarget, RefreshProvider[]> = new Map();
    private pendingRequests: Map<string, RefreshRequest> = new Map();
    private requestQueue: RefreshRequest[] = [];
    private isProcessing = false;
    private stats: RefreshStats = {
        totalRequests: 0,
        processedRequests: 0,
        deduplicatedRequests: 0,
        averageProcessingTime: 0,
        pendingRequests: 0
    };

    // Debounced refresh functions for different priorities
    private debouncedRefreshLow!: DebouncedFunction<() => Promise<void>>;
    private debouncedRefreshNormal!: DebouncedFunction<() => Promise<void>>;
    private debouncedRefreshHigh!: DebouncedFunction<() => Promise<void>>;

    // Configuration
    private config = {
        debounceDelayLow: 2000,      // 2 seconds for low priority
        debounceDelayNormal: 500,    // 500ms for normal priority
        debounceDelayHigh: 100,      // 100ms for high priority
        maxBatchSize: 10,            // Maximum requests to process in one batch
        enableBatching: true,
        enableDeduplication: true
    };

    static getInstance(): RefreshManager {
        if (!RefreshManager.instance) {
            RefreshManager.instance = new RefreshManager();
        }
        return RefreshManager.instance;
    }

    private constructor() {
        this.initializeDebouncing();
        this.loadConfiguration();
        // Logger.debug('[RefreshManager] Initialized'); // Reduced logging
    }

    private initializeDebouncing(): void {
        this.debouncedRefreshLow = debounceRefresh(
            () => this.processRequestsByPriority(RefreshPriority.LOW),
            this.config.debounceDelayLow
        );

        this.debouncedRefreshNormal = debounceRefresh(
            () => this.processRequestsByPriority(RefreshPriority.NORMAL),
            this.config.debounceDelayNormal
        );

        this.debouncedRefreshHigh = debounceRefresh(
            () => this.processRequestsByPriority(RefreshPriority.HIGH),
            this.config.debounceDelayHigh
        );
    }

    private loadConfiguration(): void {
        try {
            const config = vscode.workspace.getConfiguration('gorev.refreshManager');
            this.config.debounceDelayLow = config.get('debounceDelayLow', 2000);
            this.config.debounceDelayNormal = config.get('debounceDelayNormal', 500);
            this.config.debounceDelayHigh = config.get('debounceDelayHigh', 100);
            this.config.maxBatchSize = config.get('maxBatchSize', 10);
            this.config.enableBatching = config.get('enableBatching', true);
            this.config.enableDeduplication = config.get('enableDeduplication', true);
        } catch (error) {
            Logger.warn('[RefreshManager] Failed to load configuration, using defaults');
        }
    }

    /**
     * Register a refresh provider for specific targets
     */
    registerProvider(targets: RefreshTarget[], provider: RefreshProvider): void {
        for (const target of targets) {
            if (!this.providers.has(target)) {
                this.providers.set(target, []);
            }
            this.providers.get(target)!.push(provider);
        }
        // Logger.debug(`[RefreshManager] Registered provider '${provider.getName()}' for targets: ${targets.join(', ')}`); // Reduced logging
    }

    /**
     * Unregister a refresh provider
     */
    unregisterProvider(provider: RefreshProvider): void {
        for (const [target, providers] of this.providers.entries()) {
            const index = providers.indexOf(provider);
            if (index !== -1) {
                providers.splice(index, 1);
                if (providers.length === 0) {
                    this.providers.delete(target);
                }
            }
        }
        Logger.debug(`[RefreshManager] Unregistered provider '${provider.getName()}'`);
    }

    /**
     * Request a refresh with specified parameters
     */
    async requestRefresh(
        reason: RefreshReason,
        targets: RefreshTarget[] = [RefreshTarget.ALL],
        priority: RefreshPriority = RefreshPriority.NORMAL,
        metadata?: Record<string, any>
    ): Promise<void> {
        const request: RefreshRequest = {
            id: this.generateRequestId(),
            reason,
            priority,
            targets,
            timestamp: Date.now(),
            metadata
        };

        this.stats.totalRequests++;

        // Check for deduplication
        if (this.config.enableDeduplication && this.isDuplicateRequest(request)) {
            this.stats.deduplicatedRequests++;
            // Logger.debug(`[RefreshManager] Deduplicated request: ${request.id} (${reason})`); // Reduced logging
            return;
        }

        this.pendingRequests.set(request.id, request);
        this.stats.pendingRequests = this.pendingRequests.size;

        // Logger.debug(`[RefreshManager] Queued refresh request: ${request.id} (${reason}, ${RefreshPriority[priority]})`); // Reduced logging

        // Schedule processing based on priority
        await this.scheduleProcessing(request);
    }

    /**
     * Request immediate refresh (bypasses debouncing)
     */
    async requestImmediateRefresh(
        reason: RefreshReason,
        targets: RefreshTarget[] = [RefreshTarget.ALL],
        metadata?: Record<string, any>
    ): Promise<void> {
        const request: RefreshRequest = {
            id: this.generateRequestId(),
            reason,
            priority: RefreshPriority.CRITICAL,
            targets,
            timestamp: Date.now(),
            metadata
        };

        // Logger.debug(`[RefreshManager] Processing immediate refresh: ${request.id} (${reason})`); // Reduced logging

        await measureAsync(
            'immediate-refresh',
            async () => {
                await this.processBatch([request]);
            },
            reason,
            { targets, metadata }
        );
    }

    private async scheduleProcessing(request: RefreshRequest): Promise<void> {
        switch (request.priority) {
            case RefreshPriority.REALTIME:
                // WebSocket real-time events - immediate refresh, no debouncing
                await this.requestImmediateRefresh(request.reason, request.targets, request.metadata);
                break;
            case RefreshPriority.CRITICAL:
                await this.requestImmediateRefresh(request.reason, request.targets, request.metadata);
                break;
            case RefreshPriority.HIGH:
                await this.debouncedRefreshHigh();
                break;
            case RefreshPriority.NORMAL:
                await this.debouncedRefreshNormal();
                break;
            case RefreshPriority.LOW:
                await this.debouncedRefreshLow();
                break;
        }
    }

    private async processRequestsByPriority(priority: RefreshPriority): Promise<void> {
        if (this.isProcessing) {
            // Logger.debug('[RefreshManager] Already processing, skipping'); // Reduced logging
            return;
        }

        this.isProcessing = true;

        try {
            const requests = Array.from(this.pendingRequests.values())
                .filter(req => req.priority === priority)
                .sort((a, b) => b.timestamp - a.timestamp); // Most recent first

            if (requests.length === 0) {
                return;
            }

            const batchSize = this.config.enableBatching ?
                Math.min(requests.length, this.config.maxBatchSize) : 1;

            const batch = requests.slice(0, batchSize);

            // Logger.debug(`[RefreshManager] Processing batch of ${batch.length} requests (priority: ${RefreshPriority[priority]})`); // Reduced logging

            await measureAsync(
                'batch-refresh',
                async () => {
                    await this.processBatch(batch);
                },
                `batch-${RefreshPriority[priority]}-${batch.length}`,
                { batchSize: batch.length, priority }
            );

            // Remove processed requests
            for (const request of batch) {
                this.pendingRequests.delete(request.id);
            }

            this.stats.processedRequests += batch.length;
            this.stats.pendingRequests = this.pendingRequests.size;
            this.stats.lastRefreshTime = Date.now();

        } finally {
            this.isProcessing = false;
        }
    }

    private async processBatch(requests: RefreshRequest[]): Promise<void> {
        // Collect all unique targets from the batch
        const targets = new Set<RefreshTarget>();
        const reasons: string[] = [];

        for (const request of requests) {
            for (const target of request.targets) {
                if (target === RefreshTarget.ALL) {
                    // Add all available targets
                    this.providers.forEach((_, key) => targets.add(key));
                } else {
                    targets.add(target);
                }
            }
            reasons.push(request.reason);
        }

        // Process each target
        const targetArray = Array.from(targets);
        Logger.debug(`[RefreshManager] Refreshing targets: ${targetArray.join(', ')} (reasons: ${reasons.join(', ')})`);

        // Use setTimeout to yield control back to the event loop for better UI responsiveness
        // requestIdleCallback is not available in Node.js context, so always use setTimeout
        await new Promise<void>((resolve) => {
            setTimeout(async () => {
                await this.refreshTargets(targetArray);
                resolve();
            }, 0);
        });
    }

    private async refreshTargets(targets: RefreshTarget[]): Promise<void> {
        const refreshPromises: Promise<void>[] = [];

        for (const target of targets) {
            const providers = this.providers.get(target);
            if (providers) {
                for (const provider of providers) {
                    refreshPromises.push(
                        measureAsync(
                            `refresh-${target}`,
                            () => provider.refresh(),
                            `provider-${provider.getName()}`,
                            { target, provider: provider.getName() }
                        ).catch(error => {
                            Logger.error(`[RefreshManager] Failed to refresh ${target} with provider ${provider.getName()}:`, error);
                        })
                    );
                }
            }
        }

        // Execute all refreshes in parallel but don't block
        await Promise.allSettled(refreshPromises);
    }

    private isDuplicateRequest(request: RefreshRequest): boolean {
        for (const pendingRequest of this.pendingRequests.values()) {
            if (
                pendingRequest.reason === request.reason &&
                pendingRequest.priority === request.priority &&
                this.arraysEqual(pendingRequest.targets, request.targets) &&
                (Date.now() - pendingRequest.timestamp) < 1000 // Within 1 second
            ) {
                return true;
            }
        }
        return false;
    }

    private arraysEqual<T>(a: T[], b: T[]): boolean {
        return a.length === b.length && a.every((val, i) => val === b[i]);
    }

    private generateRequestId(): string {
        return `refresh-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    }

    /**
     * Get current statistics
     */
    getStats(): RefreshStats {
        return { ...this.stats };
    }

    /**
     * Clear all pending requests
     */
    clearPendingRequests(): void {
        this.pendingRequests.clear();
        this.stats.pendingRequests = 0;
        Logger.debug('[RefreshManager] Cleared all pending requests');
    }

    /**
     * Update configuration
     */
    updateConfiguration(newConfig: Partial<typeof this.config>): void {
        this.config = { ...this.config, ...newConfig };
        this.initializeDebouncing();
        Logger.debug('[RefreshManager] Configuration updated');
    }

    /**
     * Clear all provider caches to force fresh data load
     */
    clearAllCaches(): void {
        Logger.info('[RefreshManager] Clearing all provider caches');

        // Clear any internal caches
        this.clearPendingRequests();

        // Force all providers to clear their caches if they have any
        this.providers.forEach((providers, target) => {
            Logger.debug(`[RefreshManager] Clearing cache for target: ${target}`);
            providers.forEach(provider => {
                // If provider has a cache clearing method, call it
                if ('clearCache' in provider && typeof provider.clearCache === 'function') {
                    try {
                        (provider as any).clearCache();
                        Logger.debug('[RefreshManager] Provider cache cleared successfully');
                    } catch (error) {
                        Logger.warn('[RefreshManager] Failed to clear provider cache:', error);
                    }
                }
            });
        });

        Logger.info('[RefreshManager] All caches cleared');
    }

    /**
     * Dispose and cleanup
     */
    dispose(): void {
        this.clearPendingRequests();
        this.debouncedRefreshLow.cancel();
        this.debouncedRefreshNormal.cancel();
        this.debouncedRefreshHigh.cancel();
        this.providers.clear();
        Logger.debug('[RefreshManager] Disposed');
    }
}

// Export singleton instance
export const refreshManager = RefreshManager.getInstance();