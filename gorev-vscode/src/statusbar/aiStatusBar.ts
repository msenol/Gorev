// AI Status Bar Indicator for Gorev VS Code Extension
// Shows AI configuration status and provides quick access to AI features

import * as vscode from 'vscode';
import { ApiClient } from '../api/client';

interface AIConfig {
  provider?: string;
  model?: string;
  success?: boolean;
}

/**
 * AI Status Bar Indicator
 * - Shows AI status: "AI not configured" or provider info (e.g., "openrouter: gpt-4o-mini")
 * - Click to open AI panel or configuration
 * - Rule 15 compliance: Gracefully shows "not configured" when AI disabled
 */
export class AIStatusBar {
    private static instance: AIStatusBar;
    private statusBarItem: vscode.StatusBarItem;
    private apiClient: ApiClient;
    private currentProjectId: string | null = null;
    private aiConfigured = false;
    private providerInfo = '';

    private constructor(apiClient: ApiClient) {
        this.apiClient = apiClient;
        this.statusBarItem = vscode.window.createStatusBarItem(
            'gorev.aiStatus',
            vscode.StatusBarAlignment.Right,
            100 // Priority: lower than more important items
        );
        this.statusBarItem.name = 'Gorev AI Status';
        this.statusBarItem.command = 'gorev.ai.togglePanel';

        // Register command to toggle AI panel
        vscode.commands.registerCommand('gorev.ai.togglePanel', () => {
            vscode.commands.executeCommand('gorev.ai.showPanel');
        });

        // Register command to configure AI
        vscode.commands.registerCommand('gorev.ai.configure', () => {
            vscode.commands.executeCommand('gorev.ai.showPanel');
        });

        // Initial update
        this.update();
    }

    /**
     * Get singleton instance
     */
    public static getInstance(apiClient: ApiClient): AIStatusBar {
        if (!AIStatusBar.instance) {
            AIStatusBar.instance = new AIStatusBar(apiClient);
        }
        return AIStatusBar.instance;
    }

    /**
     * Show the status bar item
     */
    public show(): void {
        this.statusBarItem.show();
    }

    /**
     * Hide the status bar item
     */
    public hide(): void {
        this.statusBarItem.hide();
    }

    /**
     * Update AI status based on current project configuration
     */
    public async update(projectId?: string): Promise<void> {
        if (projectId) {
            this.currentProjectId = projectId;
        }

        if (!this.currentProjectId) {
            this.setNotConfigured();
            return;
        }

        try {
            // Check AI configuration via MCP tool
            const result = await vscode.commands.executeCommand<unknown>(
                'gorev.mcp.call',
                'gorev_ai',
                'action=configure'
            );

            if (result && (result as { success: boolean }).success) {
                this.aiConfigured = true;
                this.providerInfo = this.formatProviderInfo(result);
                this.setConfigured();
            } else {
                this.aiConfigured = false;
                this.setNotConfigured();
            }
        } catch {
            // If AI check fails, assume not configured (Rule 15 compliance)
            this.aiConfigured = false;
            this.setNotConfigured();
        }
    }

    /**
     * Set status bar to "not configured" state
     */
    private setNotConfigured(): void {
        this.statusBarItem.text = '$(robot) Gorev AI';
        this.statusBarItem.tooltip = 'AI is not configured for this project. Click to configure.';
        this.statusBarItem.color = undefined;
        this.aiConfigured = false;
    }

    /**
     * Set status bar to "configured" state with provider info
     */
    private setConfigured(): void {
        this.statusBarItem.text = `$(robot) ${this.providerInfo}`;
        this.statusBarItem.tooltip = `AI configured: ${this.providerInfo}\nClick to open AI panel`;
        this.statusBarItem.color = new vscode.ThemeColor('terminal.ansiGreen');
        this.aiConfigured = true;
    }

    /**
     * Format provider information for display
     */
    private formatProviderInfo(config: unknown): string {
        const aiConfig = config as AIConfig;
        const provider = aiConfig.provider || 'unknown';
        const model = aiConfig.model || 'unknown';

        // Shorten provider names
        const providerShort = provider === 'openrouter' ? 'OR' :
                             provider === 'anannas' ? 'AN' :
                             provider.substring(0, 2).toUpperCase();

        // Shorten model names
        const modelShort = model.includes('/') ? model.split('/').pop() : model;
        const modelShortened = modelShort ? modelShort.substring(0, 15) : 'unknown';

        return `${providerShort}: ${modelShortened}`;
    }

    /**
     * Check if AI is configured for current project
     */
    public isAIConfigured(): boolean {
        return this.aiConfigured;
    }

    /**
     * Get current provider info
     */
    public getProviderInfo(): string {
        return this.providerInfo;
    }

    /**
     * Dispose of resources
     */
    public dispose(): void {
        this.statusBarItem.dispose();
    }
}
