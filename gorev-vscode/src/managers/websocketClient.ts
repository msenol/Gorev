import * as vscode from 'vscode';
import WebSocket from 'ws';
import { RefreshManager, RefreshReason, RefreshTarget, RefreshPriority } from './refreshManager';

/**
 * WebSocket change event types
 */
export enum ChangeEventType {
    TASK_CREATED = 'task_created',
    TASK_UPDATED = 'task_updated',
    TASK_DELETED = 'task_deleted',
    PROJECT_CREATED = 'project_created',
    PROJECT_UPDATED = 'project_updated',
    PROJECT_DELETED = 'project_deleted',
    TEMPLATE_CHANGED = 'template_changed',
    WORKSPACE_SYNC = 'workspace_sync'
}

/**
 * WebSocket change event structure
 */
export interface ChangeEvent {
    type: ChangeEventType;
    workspace_id: string;
    entity_id?: string;
    entity_type?: string;
    action?: string;
    data?: Record<string, unknown>;
    timestamp: number;
}

/**
 * WebSocket client for real-time database change notifications
 */
export class WebSocketClient {
    private ws: WebSocket | null = null;
    private reconnectTimer: NodeJS.Timeout | null = null;
    private readonly reconnectDelay = 5000; // 5 seconds
    private readonly pingInterval = 30000; // 30 seconds
    private pingTimer: NodeJS.Timeout | null = null;
    private isConnecting = false;
    private isDisposed = false;

    constructor(
        private readonly daemonUrl: string,
        private readonly workspaceId: string,
        private readonly outputChannel: vscode.OutputChannel
    ) {}

    /**
     * Connect to WebSocket server
     */
    public connect(): void {
        if (this.isDisposed) {
            this.outputChannel.appendLine('[WebSocket] Cannot connect - client is disposed');
            return;
        }

        if (this.isConnecting || this.ws?.readyState === WebSocket.OPEN) {
            this.outputChannel.appendLine('[WebSocket] Already connected or connecting');
            return;
        }

        this.isConnecting = true;

        try {
            // Convert HTTP URL to WebSocket URL
            const wsUrl = this.daemonUrl.replace(/^http/, 'ws');
            const fullUrl = `${wsUrl}/ws?workspace_id=${this.workspaceId}`;

            this.outputChannel.appendLine(`[WebSocket] Connecting to ${fullUrl}`);

            this.ws = new WebSocket(fullUrl);

            this.ws.on('open', () => {
                this.isConnecting = false;
                this.outputChannel.appendLine('[WebSocket] ✓ Connected');

                // Start ping timer
                this.startPingTimer();
            });

            this.ws.on('message', (data: WebSocket.Data) => {
                try {
                    const event: ChangeEvent = JSON.parse(data.toString());
                    this.handleChangeEvent(event);
                } catch (error) {
                    this.outputChannel.appendLine(`[WebSocket] Parse error: ${error}`);
                }
            });

            this.ws.on('error', (error: Error) => {
                this.isConnecting = false;
                this.outputChannel.appendLine(`[WebSocket] ✗ Error: ${error.message}`);
            });

            this.ws.on('close', () => {
                this.isConnecting = false;
                this.stopPingTimer();
                this.outputChannel.appendLine('[WebSocket] ✗ Connection closed');

                // Reconnect if not disposed
                if (!this.isDisposed) {
                    this.scheduleReconnect();
                }
            });

        } catch (error) {
            this.isConnecting = false;
            this.outputChannel.appendLine(`[WebSocket] Connection error: ${error}`);
            this.scheduleReconnect();
        }
    }

    /**
     * Handle incoming change events
     */
    private handleChangeEvent(event: ChangeEvent): void {
        this.outputChannel.appendLine(
            `[WebSocket] Event: ${event.type} | Entity: ${event.entity_type || 'N/A'} | Action: ${event.action || 'N/A'}`
        );

        const refreshManager = RefreshManager.getInstance();
        const targets: RefreshTarget[] = [];

        // Map event types to refresh targets
        switch (event.type) {
            case ChangeEventType.TASK_CREATED:
            case ChangeEventType.TASK_UPDATED:
            case ChangeEventType.TASK_DELETED:
                targets.push(RefreshTarget.TASKS);
                break;

            case ChangeEventType.PROJECT_CREATED:
            case ChangeEventType.PROJECT_UPDATED:
            case ChangeEventType.PROJECT_DELETED:
                targets.push(RefreshTarget.PROJECTS);
                break;

            case ChangeEventType.TEMPLATE_CHANGED:
                targets.push(RefreshTarget.TEMPLATES);
                break;

            case ChangeEventType.WORKSPACE_SYNC:
                // Sync all targets
                targets.push(RefreshTarget.ALL);
                break;
        }

        // Request refresh with real-time priority
        if (targets.length > 0) {
            refreshManager.requestRefresh(
                RefreshReason.EXTERNAL_CHANGE,
                targets,
                RefreshPriority.REALTIME
            );
        }
    }

    /**
     * Start ping timer to keep connection alive
     */
    private startPingTimer(): void {
        this.stopPingTimer();

        this.pingTimer = setInterval(() => {
            if (this.ws?.readyState === WebSocket.OPEN) {
                const ping = {
                    type: 'ping',
                    timestamp: Date.now()
                };
                this.ws.send(JSON.stringify(ping));
            }
        }, this.pingInterval);
    }

    /**
     * Stop ping timer
     */
    private stopPingTimer(): void {
        if (this.pingTimer) {
            clearInterval(this.pingTimer);
            this.pingTimer = null;
        }
    }

    /**
     * Schedule reconnection attempt
     */
    private scheduleReconnect(): void {
        if (this.reconnectTimer || this.isDisposed) {
            return;
        }

        this.outputChannel.appendLine(`[WebSocket] Reconnecting in ${this.reconnectDelay / 1000}s...`);

        this.reconnectTimer = setTimeout(() => {
            this.reconnectTimer = null;
            if (!this.isDisposed) {
                this.connect();
            }
        }, this.reconnectDelay);
    }

    /**
     * Disconnect from WebSocket server
     */
    public disconnect(): void {
        this.isDisposed = true;

        // Clear timers
        this.stopPingTimer();
        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
        }

        // Close WebSocket
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }

        this.outputChannel.appendLine('[WebSocket] Disconnected');
    }

    /**
     * Check if connected
     */
    public isConnected(): boolean {
        return this.ws?.readyState === WebSocket.OPEN;
    }
}
