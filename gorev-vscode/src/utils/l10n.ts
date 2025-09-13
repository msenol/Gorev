import * as vscode from 'vscode';
import * as l10n from '@vscode/l10n';
import { readFileSync, existsSync } from 'fs';
import { join } from 'path';

/**
 * Enhanced L10n manager for Gorev VS Code extension
 * Provides robust localization with proper fallback mechanisms
 */
export class L10nManager {
    private static instance: L10nManager;
    private initialized = false;
    private bundles: Map<string, Record<string, string>> = new Map();
    private currentLocale: string = 'en';

    private constructor() {}

    public static getInstance(): L10nManager {
        if (!L10nManager.instance) {
            L10nManager.instance = new L10nManager();
        }
        return L10nManager.instance;
    }

    /**
     * Initialize the L10n system
     */
    public async initialize(context: vscode.ExtensionContext): Promise<void> {
        console.log('[GOREV-L10N] 2. L10n manager initializing at:', new Date().toISOString());
        console.log('[GOREV-L10N] 3. Extension path:', context.extensionPath);

        if (this.initialized) {
            console.log('[GOREV-L10N] 4. Already initialized, skipping');
            return;
        }

        // Get current locale
        this.currentLocale = vscode.env.language || 'en';
        console.log('[GOREV-L10N] 5. Current locale:', this.currentLocale);

        // Load bundles
        await this.loadBundles(context);

        // Configure @vscode/l10n
        const bundleUri = this.getBundleUri(context, this.currentLocale);
        if (bundleUri && existsSync(bundleUri.fsPath)) {
            try {
                const bundleContent = readFileSync(bundleUri.fsPath, 'utf8');
                const bundle = JSON.parse(bundleContent);
                l10n.config({
                    contents: bundle
                });
                console.log('[GOREV-L10N] 6. VS Code l10n configured successfully');
            } catch (error) {
                console.log('[GOREV-L10N] 6. Failed to configure @vscode/l10n:', error);
            }
        } else {
            console.log('[GOREV-L10N] 6. No bundle URI found for VS Code l10n');
        }

        this.initialized = true;
        console.log('[GOREV-L10N] 7. L10n initialization completed');
    }

    /**
     * Load localization bundles
     */
    private async loadBundles(context: vscode.ExtensionContext): Promise<void> {
        const l10nPath = join(context.extensionPath, 'l10n');
        console.log('[GOREV-L10N] 8. L10n path:', l10nPath);

        // Load English bundle (fallback)
        const enBundlePath = join(l10nPath, 'bundle.l10n.json');
        console.log('[GOREV-L10N] 9. EN bundle exists:', existsSync(enBundlePath));
        if (existsSync(enBundlePath)) {
            try {
                const content = readFileSync(enBundlePath, 'utf8');
                const bundle = JSON.parse(content);
                this.bundles.set('en', bundle);
                console.log('[GOREV-L10N] 10. EN bundle loaded with', Object.keys(bundle).length, 'keys');
            } catch (error) {
                console.log('[GOREV-L10N] 10. Failed to load English bundle:', error);
            }
        }

        // Load Turkish bundle
        const trBundlePath = join(l10nPath, 'bundle.l10n.tr.json');
        console.log('[GOREV-L10N] 11. TR bundle exists:', existsSync(trBundlePath));
        if (existsSync(trBundlePath)) {
            try {
                const content = readFileSync(trBundlePath, 'utf8');
                const bundle = JSON.parse(content);
                this.bundles.set('tr', bundle);
                console.log('[GOREV-L10N] 12. TR bundle loaded with', Object.keys(bundle).length, 'keys');
            } catch (error) {
                console.log('[GOREV-L10N] 12. Failed to load Turkish bundle:', error);
            }
        }
    }

    /**
     * Get bundle URI for a specific locale
     */
    private getBundleUri(context: vscode.ExtensionContext, locale: string): vscode.Uri | null {
        const l10nPath = join(context.extensionPath, 'l10n');
        let bundlePath: string;

        if (locale.startsWith('tr')) {
            bundlePath = join(l10nPath, 'bundle.l10n.tr.json');
        } else {
            bundlePath = join(l10nPath, 'bundle.l10n.json');
        }

        return existsSync(bundlePath) ? vscode.Uri.file(bundlePath) : null;
    }

    /**
     * Translate a key with arguments
     */
    public t(key: string, ...args: (string | number | boolean | Record<string, any>)[]): string {
        if (!this.initialized) {
            console.log('[GOREV-L10N] 13. Manager not initialized, returning key:', key);
            return key;
        }

        console.log('[GOREV-L10N] 14. Translating key:', key, 'with', args.length, 'args');
        return this.manualLookup(key, args);
    }

    /**
     * Manual bundle lookup with fallback
     */
    private manualLookup(key: string, args: (string | number | boolean | Record<string, any>)[]): string {
        const simpleLocale = this.getSimpleLocale(this.currentLocale);
        console.log('[GOREV-L10N] 15. Looking up key:', key, 'for locale:', simpleLocale);
        console.log('[GOREV-L10N] 16. Available bundles:', Array.from(this.bundles.keys()));

        // Try current locale first
        let bundle = this.bundles.get(simpleLocale);
        if (!bundle) {
            console.log('[GOREV-L10N] 17. No bundle for locale', simpleLocale, ', trying English');
            // Fallback to English
            bundle = this.bundles.get('en');
        }

        if (!bundle) {
            console.log('[GOREV-L10N] 18. No bundle found at all, returning key:', key);
            return key;
        }

        let translation = bundle[key];
        if (!translation) {
            // Try English as fallback
            const enBundle = this.bundles.get('en');
            translation = enBundle?.[key] || key;
            console.log('[GOREV-L10N] 19. Translation not found in primary bundle, EN fallback result:', translation);
        } else {
            console.log('[GOREV-L10N] 20. Translation found:', translation);
        }

        // Replace placeholders
        if (args.length > 0) {
            const result = this.replacePlaceholders(translation, args);
            console.log('[GOREV-L10N] 21. Final result after placeholder replacement:', result);
            return result;
        }

        return translation;
    }

    /**
     * Replace placeholders in translation
     */
    private replacePlaceholders(text: string, args: (string | number | boolean | Record<string, any>)[]): string {
        return text.replace(/\{(\d+)\}/g, (match, index) => {
            const argIndex = parseInt(index, 10);
            const arg = args[argIndex];
            if (arg !== undefined) {
                if (typeof arg === 'object') {
                    // For object parameters, try to extract the first value
                    const values = Object.values(arg);
                    return values.length > 0 ? String(values[0]) : match;
                }
                return String(arg);
            }
            return match;
        });
    }

    /**
     * Get simple locale (e.g., 'tr' from 'tr-TR')
     */
    private getSimpleLocale(locale: string): string {
        return locale.split('-')[0].toLowerCase();
    }

    /**
     * Get current locale
     */
    public getCurrentLocale(): string {
        return this.currentLocale;
    }

    /**
     * Check if manager is initialized
     */
    public isInitialized(): boolean {
        return this.initialized;
    }
}

// Convenience function for global use
let globalManager: L10nManager | null = null;

/**
 * Initialize global L10n manager
 */
export async function initializeL10n(context: vscode.ExtensionContext): Promise<void> {
    globalManager = L10nManager.getInstance();
    await globalManager.initialize(context);
}

/**
 * Translate a key using global manager
 */
export function t(key: string, ...args: (string | number | boolean | Record<string, any>)[]): string {
    if (!globalManager) {
        console.warn('L10n not initialized, returning key');
        return key;
    }
    return globalManager.t(key, ...args);
}

/**
 * Get current locale
 */
export function getCurrentLocale(): string {
    return globalManager?.getCurrentLocale() || 'en';
}