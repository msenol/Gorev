/**
 * DatabaseFileWatcher - Monitors SQLite database changes and triggers refresh
 * Provides real-time synchronization between web/AI changes and VS Code extension
 */

import * as vscode from 'vscode';
import { RefreshManager, RefreshTarget, RefreshReason, RefreshPriority } from './refreshManager';
import { debounceRefresh } from '../utils/debounce';
import { Logger } from '../utils/logger';

export class DatabaseFileWatcher {
    private fileWatcher: vscode.FileSystemWatcher | undefined;
    private debouncedRefresh: ReturnType<typeof debounceRefresh>;
    private refreshManager: RefreshManager;

    constructor(refreshManager: RefreshManager) {
        this.refreshManager = refreshManager;

        // Initialize debounced refresh with 500ms delay to prevent excessive refreshes
        this.debouncedRefresh = debounceRefresh(
            () => this.handleDatabaseChange(),
            500
        );
    }

    /**
     * Start watching the database file
     */
    start(): void {
        const workspaceFolders = vscode.workspace.workspaceFolders;
        if (!workspaceFolders || workspaceFolders.length === 0) {
            Logger.warn('[DatabaseFileWatcher] No workspace folder found, cannot watch database');
            return;
        }

        // Watch .gorev/gorev.db in workspace
        const dbPattern = new vscode.RelativePattern(
            workspaceFolders[0],
            '.gorev/gorev.db'
        );

        this.fileWatcher = vscode.workspace.createFileSystemWatcher(
            dbPattern,
            false, // ignoreCreateEvents - watch create
            false, // ignoreChangeEvents - watch changes
            false  // ignoreDeleteEvents - watch delete
        );

        // Listen for database changes
        this.fileWatcher.onDidChange(() => {
            Logger.debug('[DatabaseFileWatcher] Database file changed');
            this.debouncedRefresh();
        });

        this.fileWatcher.onDidCreate(() => {
            Logger.debug('[DatabaseFileWatcher] Database file created');
            this.debouncedRefresh();
        });

        this.fileWatcher.onDidDelete(() => {
            Logger.warn('[DatabaseFileWatcher] Database file deleted');
            this.debouncedRefresh();
        });

        Logger.info('[DatabaseFileWatcher] Started watching database file: .gorev/gorev.db');
    }

    /**
     * Handle database change with smart refresh
     *
     * Future enhancement: Analyze which table changed and refresh only relevant providers
     * For now, refresh all providers to ensure consistency
     */
    private async handleDatabaseChange(): Promise<void> {
        Logger.debug('[DatabaseFileWatcher] Processing database change...');

        try {
            // Refresh all providers with HIGH priority for immediate update
            // This ensures web/AI changes are immediately reflected in VS Code
            await this.refreshManager.requestRefresh(
                RefreshReason.DATA_CHANGE,
                [RefreshTarget.ALL],
                RefreshPriority.HIGH
            );

            Logger.debug('[DatabaseFileWatcher] Database change processed successfully');
        } catch (error) {
            Logger.error('[DatabaseFileWatcher] Failed to process database change:', error);
        }
    }

    /**
     * Stop watching and cleanup
     */
    dispose(): void {
        if (this.fileWatcher) {
            this.fileWatcher.dispose();
            this.fileWatcher = undefined;
            Logger.info('[DatabaseFileWatcher] Stopped watching database file');
        }

        if (this.debouncedRefresh) {
            this.debouncedRefresh.cancel();
        }
    }
}
