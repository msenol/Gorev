/**
 * Performance monitoring utility for tracking and optimizing refresh operations
 * Provides metrics, timing, and diagnostic information
 */

import { Logger } from './logger';

export interface PerformanceMetrics {
    operation: string;
    startTime: number;
    endTime?: number;
    duration?: number;
    success: boolean;
    reason?: string;
    memoryUsage?: NodeJS.MemoryUsage;
    metadata?: Record<string, any>;
}

export interface PerformanceAggregates {
    totalOperations: number;
    successfulOperations: number;
    failedOperations: number;
    averageDuration: number;
    minDuration: number;
    maxDuration: number;
    totalDuration: number;
    operationsPerMinute: number;
    lastOperation?: PerformanceMetrics;
}

export class PerformanceMonitor {
    private static instance: PerformanceMonitor;
    private metrics: Map<string, PerformanceMetrics[]> = new Map();
    private activeOperations: Map<string, PerformanceMetrics> = new Map();
    private maxMetricsPerOperation = 100; // Keep last 100 metrics per operation
    private enabled = true;

    static getInstance(): PerformanceMonitor {
        if (!PerformanceMonitor.instance) {
            PerformanceMonitor.instance = new PerformanceMonitor();
        }
        return PerformanceMonitor.instance;
    }

    private constructor() {
        // Check if performance monitoring is enabled in configuration
        this.checkConfiguration();
    }

    private checkConfiguration(): void {
        try {
            const vscode = require('vscode');
            const config = vscode.workspace.getConfiguration('gorev');
            this.enabled = config.get('performance.enableMonitoring', true);
            this.maxMetricsPerOperation = config.get('performance.maxMetricsPerOperation', 100);
        } catch (error) {
            // Fallback if VS Code API not available
            this.enabled = true;
        }
    }

    /**
     * Start timing an operation
     */
    startOperation(operationId: string, operation: string, reason?: string, metadata?: Record<string, any>): void {
        if (!this.enabled) return;

        const metric: PerformanceMetrics = {
            operation,
            startTime: performance.now(),
            success: false,
            reason,
            memoryUsage: process.memoryUsage(),
            metadata
        };

        this.activeOperations.set(operationId, metric);

        Logger.debug(`[Performance] Started: ${operation} (ID: ${operationId})${reason ? ` - Reason: ${reason}` : ''}`);
    }

    /**
     * End timing an operation
     */
    endOperation(operationId: string, success = true, error?: Error): PerformanceMetrics | undefined {
        if (!this.enabled) return;

        const metric = this.activeOperations.get(operationId);
        if (!metric) {
            Logger.warn(`[Performance] Operation not found: ${operationId}`);
            return;
        }

        const endTime = performance.now();
        metric.endTime = endTime;
        metric.duration = endTime - metric.startTime;
        metric.success = success;

        if (error) {
            metric.metadata = { ...metric.metadata, error: error.message };
        }

        // Store the completed metric
        const operationMetrics = this.metrics.get(metric.operation) || [];
        operationMetrics.push(metric);

        // Keep only the most recent metrics
        if (operationMetrics.length > this.maxMetricsPerOperation) {
            operationMetrics.splice(0, operationMetrics.length - this.maxMetricsPerOperation);
        }

        this.metrics.set(metric.operation, operationMetrics);
        this.activeOperations.delete(operationId);

        // Log performance information
        const status = success ? 'SUCCESS' : 'FAILED';
        const duration = metric.duration!.toFixed(2);
        Logger.debug(`[Performance] ${status}: ${metric.operation} (${duration}ms)${metric.reason ? ` - ${metric.reason}` : ''}`);

        // Warn about slow operations
        if (metric.duration! > 1000) {
            Logger.warn(`[Performance] Slow operation detected: ${metric.operation} took ${duration}ms`);
        }

        return metric;
    }

    /**
     * Get performance aggregates for an operation
     */
    getAggregates(operation: string): PerformanceAggregates | undefined {
        const operationMetrics = this.metrics.get(operation);
        if (!operationMetrics || operationMetrics.length === 0) {
            return undefined;
        }

        const completedMetrics = operationMetrics.filter(m => m.duration !== undefined);
        const successfulMetrics = completedMetrics.filter(m => m.success);
        const durations = completedMetrics.map(m => m.duration!);

        const totalDuration = durations.reduce((sum, duration) => sum + duration, 0);
        const avgDuration = totalDuration / durations.length;

        // Calculate operations per minute based on time range
        const timeRange = Math.max(1, (Date.now() - (completedMetrics[0]?.startTime || Date.now())) / 1000 / 60);
        const operationsPerMinute = completedMetrics.length / timeRange;

        return {
            totalOperations: completedMetrics.length,
            successfulOperations: successfulMetrics.length,
            failedOperations: completedMetrics.length - successfulMetrics.length,
            averageDuration: avgDuration,
            minDuration: Math.min(...durations),
            maxDuration: Math.max(...durations),
            totalDuration,
            operationsPerMinute,
            lastOperation: operationMetrics[operationMetrics.length - 1]
        };
    }

    /**
     * Get all operation names being tracked
     */
    getTrackedOperations(): string[] {
        return Array.from(this.metrics.keys());
    }

    /**
     * Get performance summary for all operations
     */
    getPerformanceSummary(): Record<string, PerformanceAggregates> {
        const summary: Record<string, PerformanceAggregates> = {};

        for (const operation of this.getTrackedOperations()) {
            const aggregates = this.getAggregates(operation);
            if (aggregates) {
                summary[operation] = aggregates;
            }
        }

        return summary;
    }

    /**
     * Clear all metrics
     */
    clear(): void {
        this.metrics.clear();
        this.activeOperations.clear();
        Logger.debug('[Performance] Cleared all metrics');
    }

    /**
     * Enable or disable performance monitoring
     */
    setEnabled(enabled: boolean): void {
        this.enabled = enabled;
        Logger.debug(`[Performance] Monitoring ${enabled ? 'enabled' : 'disabled'}`);
    }

    /**
     * Check if monitoring is enabled
     */
    isEnabled(): boolean {
        return this.enabled;
    }

    /**
     * Generate a performance report
     */
    generateReport(): string {
        const summary = this.getPerformanceSummary();
        const lines: string[] = [];

        lines.push('=== Gorev Extension Performance Report ===');
        lines.push(`Generated: ${new Date().toISOString()}`);
        lines.push(`Monitoring Enabled: ${this.enabled}`);
        lines.push('');

        if (Object.keys(summary).length === 0) {
            lines.push('No performance data available.');
            return lines.join('\n');
        }

        for (const [operation, aggregates] of Object.entries(summary)) {
            lines.push(`Operation: ${operation}`);
            lines.push(`  Total Operations: ${aggregates.totalOperations}`);
            lines.push(`  Success Rate: ${((aggregates.successfulOperations / aggregates.totalOperations) * 100).toFixed(1)}%`);
            lines.push(`  Average Duration: ${aggregates.averageDuration.toFixed(2)}ms`);
            lines.push(`  Min/Max Duration: ${aggregates.minDuration.toFixed(2)}ms / ${aggregates.maxDuration.toFixed(2)}ms`);
            lines.push(`  Operations/Minute: ${aggregates.operationsPerMinute.toFixed(1)}`);

            if (aggregates.lastOperation) {
                const lastReason = aggregates.lastOperation.reason || 'No reason provided';
                lines.push(`  Last Reason: ${lastReason}`);
            }

            lines.push('');
        }

        return lines.join('\n');
    }
}

/**
 * Decorator for automatic performance monitoring
 */
export function performanceMonitor(operation: string, logResults = true) {
    return function (target: any, propertyName: string, descriptor: PropertyDescriptor) {
        const method = descriptor.value;

        descriptor.value = async function (...args: any[]) {
            const monitor = PerformanceMonitor.getInstance();
            const operationId = `${operation}-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

            monitor.startOperation(operationId, operation);

            try {
                const result = await method.apply(this, args);
                const metrics = monitor.endOperation(operationId, true);

                if (logResults && metrics) {
                    Logger.debug(`[Performance] ${operation} completed in ${metrics.duration!.toFixed(2)}ms`);
                }

                return result;
            } catch (error) {
                monitor.endOperation(operationId, false, error as Error);
                throw error;
            }
        };

        return descriptor;
    };
}

/**
 * Utility function to measure async operations
 */
export async function measureAsync<T>(
    operation: string,
    func: () => Promise<T>,
    reason?: string,
    metadata?: Record<string, any>
): Promise<T> {
    const monitor = PerformanceMonitor.getInstance();
    const operationId = `${operation}-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

    monitor.startOperation(operationId, operation, reason, metadata);

    try {
        const result = await func();
        monitor.endOperation(operationId, true);
        return result;
    } catch (error) {
        monitor.endOperation(operationId, false, error as Error);
        throw error;
    }
}

/**
 * Utility function to measure synchronous operations
 */
export function measureSync<T>(
    operation: string,
    func: () => T,
    reason?: string,
    metadata?: Record<string, any>
): T {
    const monitor = PerformanceMonitor.getInstance();
    const operationId = `${operation}-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

    monitor.startOperation(operationId, operation, reason, metadata);

    try {
        const result = func();
        monitor.endOperation(operationId, true);
        return result;
    } catch (error) {
        monitor.endOperation(operationId, false, error as Error);
        throw error;
    }
}

// Export singleton instance for convenience
export const performanceMonitorInstance = PerformanceMonitor.getInstance();