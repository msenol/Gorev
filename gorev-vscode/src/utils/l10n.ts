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
        if (this.initialized) {
            return;
        }

        // Get current locale
        this.currentLocale = vscode.env.language || 'en';

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
            } catch (error) {
                console.warn(`Failed to configure @vscode/l10n: ${error}`);
            }
        }

        this.initialized = true;
    }

    /**
     * Load localization bundles
     */
    private async loadBundles(context: vscode.ExtensionContext): Promise<void> {
        const l10nPath = join(context.extensionPath, 'l10n');

        // Load English bundle (fallback)
        const enBundlePath = join(l10nPath, 'bundle.l10n.json');
        if (existsSync(enBundlePath)) {
            try {
                const content = readFileSync(enBundlePath, 'utf8');
                this.bundles.set('en', JSON.parse(content));
            } catch (error) {
                console.warn(`Failed to load English bundle: ${error}`);
            }
        }

        // Load Turkish bundle
        const trBundlePath = join(l10nPath, 'bundle.l10n.tr.json');
        if (existsSync(trBundlePath)) {
            try {
                const content = readFileSync(trBundlePath, 'utf8');
                this.bundles.set('tr', JSON.parse(content));
            } catch (error) {
                console.warn(`Failed to load Turkish bundle: ${error}`);
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
            console.warn('L10nManager not initialized, returning key');
            return key;
        }

        // Try @vscode/l10n first (works when VS Code has proper locale)
        try {
            if (args.length === 0) {
                const result = l10n.t(key);
                if (result !== key) {
                    return result;
                }
            } else {
                // Convert our args to l10n compatible format
                const convertedArgs = args.map(arg => {
                    if (typeof arg === 'object') {
                        return String(Object.values(arg)[0] || '');
                    }
                    return arg;
                }) as (string | number | boolean)[];

                const result = l10n.t(key, ...convertedArgs);
                if (result !== key) {
                    return result;
                }
            }
        } catch (error) {
            // Fall through to manual lookup
        }

        // Fallback to manual bundle lookup
        return this.manualLookup(key, args);
    }

    /**
     * Manual bundle lookup with fallback
     */
    private manualLookup(key: string, args: (string | number | boolean | Record<string, any>)[]): string {
        // Try current locale first
        let bundle = this.bundles.get(this.getSimpleLocale(this.currentLocale));
        if (!bundle) {
            // Fallback to English
            bundle = this.bundles.get('en');
        }

        if (!bundle) {
            console.warn(`No bundle found for key: ${key}`);
            return key;
        }

        let translation = bundle[key];
        if (!translation) {
            // Try English as fallback
            const enBundle = this.bundles.get('en');
            translation = enBundle?.[key] || key;
        }

        if (!translation) {
            console.warn(`Translation not found for key: ${key}`);
            return key;
        }

        // Replace placeholders
        if (args.length > 0) {
            return this.replacePlaceholders(translation, args);
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