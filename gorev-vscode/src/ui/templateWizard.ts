import * as vscode from 'vscode';
import * as path from 'path';
import { MCPClient } from '../mcp/client';
import { GorevTemplate, TemplateAlan } from '../models/template';
import { MarkdownParser } from '../utils/markdownParser';
import { Logger } from '../utils/logger';

/**
 * Template Wizard - Görev şablonlarından görev oluşturma için zengin UI
 */
export class TemplateWizard {
    private static currentPanel: TemplateWizard | undefined;
    private readonly panel: vscode.WebviewPanel;
    private template: GorevTemplate | undefined;
    private readonly disposables: vscode.Disposable[] = [];

    constructor(
        private readonly mcpClient: MCPClient,
        private readonly extensionUri: vscode.Uri
    ) {
        this.panel = vscode.window.createWebviewPanel(
            'gorevTemplateWizard',
            vscode.l10n.t('templateWizard.title'),
            vscode.ViewColumn.One,
            {
                enableScripts: true,
                retainContextWhenHidden: true,
                localResourceRoots: [
                    vscode.Uri.joinPath(extensionUri, 'media'),
                    vscode.Uri.joinPath(extensionUri, 'dist')
                ]
            }
        );

        // Set icon
        this.panel.iconPath = {
            light: vscode.Uri.joinPath(extensionUri, 'media', 'template-light.svg'),
            dark: vscode.Uri.joinPath(extensionUri, 'media', 'template-dark.svg')
        };

        // Set HTML content
        this.panel.webview.html = this.getHtmlContent();

        // Handle messages from webview
        this.panel.webview.onDidReceiveMessage(
            async message => {
                switch (message.command) {
                    case 'selectTemplate':
                        await this.selectTemplate(message.templateId);
                        break;
                    case 'createTask':
                        await this.createTaskFromTemplate(message.values);
                        break;
                    case 'loadTemplates':
                        await this.loadTemplates(message.category);
                        break;
                    case 'searchTemplates':
                        await this.searchTemplates(message.query);
                        break;
                    case 'previewTask':
                        this.previewTask(message.values);
                        break;
                    case 'saveAsFavorite':
                        await this.saveAsFavorite(message.templateId);
                        break;
                    case 'loadFavorites':
                        await this.loadFavorites();
                        break;
                }
            },
            null,
            this.disposables
        );

        // Handle panel disposal
        this.panel.onDidDispose(() => this.dispose(), null, this.disposables);

        // Load initial templates
        this.loadTemplates();
    }

    public static async show(mcpClient: MCPClient, extensionUri: vscode.Uri, templateId?: string): Promise<void> {
        // If panel already exists, reveal it
        if (TemplateWizard.currentPanel) {
            TemplateWizard.currentPanel.panel.reveal(vscode.ViewColumn.One);
            if (templateId) {
                await TemplateWizard.currentPanel.selectTemplate(templateId);
            }
            return;
        }

        // Create new panel
        TemplateWizard.currentPanel = new TemplateWizard(mcpClient, extensionUri);
        if (templateId) {
            await TemplateWizard.currentPanel.selectTemplate(templateId);
        }
    }

    private async loadTemplates(category?: string): Promise<void> {
        try {
            const result = await this.mcpClient.callTool('template_listele', { 
                kategori: category 
            });
            
            const templates = MarkdownParser.parseTemplateListesi(result.content[0].text);
            
            // Send templates to webview
            await this.panel.webview.postMessage({
                command: 'templatesLoaded',
                templates: templates
            });
        } catch (error) {
            Logger.error('Failed to load templates:', error);
            vscode.window.showErrorMessage(vscode.l10n.t('templateWizard.loadFailed'));
        }
    }

    private async searchTemplates(query: string): Promise<void> {
        try {
            // Load all templates and filter client-side
            const result = await this.mcpClient.callTool('template_listele');
            const allTemplates = MarkdownParser.parseTemplateListesi(result.content[0].text);
            
            const filtered = allTemplates.filter(t => 
                t.isim.toLowerCase().includes(query.toLowerCase()) ||
                t.tanim?.toLowerCase().includes(query.toLowerCase()) ||
                t.kategori?.toLowerCase().includes(query.toLowerCase())
            );
            
            await this.panel.webview.postMessage({
                command: 'searchResults',
                templates: filtered
            });
        } catch (error) {
            Logger.error('Failed to search templates:', error);
        }
    }

    private async selectTemplate(templateId: string): Promise<void> {
        try {
            // Load template details
            const result = await this.mcpClient.callTool('template_listele');
            const templates = MarkdownParser.parseTemplateListesi(result.content[0].text);
            
            this.template = templates.find(t => t.id === templateId);
            if (!this.template) {
                throw new Error(vscode.l10n.t('templateWizard.notFound'));
            }

            // Send template details to webview
            await this.panel.webview.postMessage({
                command: 'templateSelected',
                template: this.template
            });
        } catch (error) {
            Logger.error('Failed to load template:', error);
            vscode.window.showErrorMessage(vscode.l10n.t('templateWizard.loadTemplateFailed'));
        }
    }

    private async createTaskFromTemplate(values: Record<string, any>): Promise<void> {
        if (!this.template) {
            vscode.window.showErrorMessage(vscode.l10n.t('templateWizard.notSelected'));
            return;
        }

        try {
            // Validate required fields
            const missingFields = this.template.alanlar
                .filter(field => field.zorunlu && !values[field.isim])
                .map(field => field.isim);

            if (missingFields.length > 0) {
                await this.panel.webview.postMessage({
                    command: 'validationError',
                    fields: missingFields
                });
                return;
            }

            // Create task from template
            const result = await this.mcpClient.callTool('templateden_gorev_olustur', {
                template_id: this.template.id,
                degerler: values
            });

            // Show success message
            vscode.window.showInformationMessage(vscode.l10n.t('templateWizard.taskCreated'));

            // Close wizard
            this.panel.dispose();

            // Refresh task list
            await vscode.commands.executeCommand('gorev.refreshTasks');
        } catch (error) {
            Logger.error('Failed to create task from template:', error);
            vscode.window.showErrorMessage(vscode.l10n.t('templateWizard.createFailed'));
        }
    }

    private previewTask(values: Record<string, any>): void {
        if (!this.template) return;

        // Generate preview based on template and values
        const preview = this.generateTaskPreview(this.template, values);
        
        this.panel.webview.postMessage({
            command: 'previewGenerated',
            preview: preview
        });
    }

    private generateTaskPreview(template: GorevTemplate, values: Record<string, any>): string {
        let preview = `# ${values.baslik || template.varsayilan_baslik || vscode.l10n.t('templateWizard.newTask')}\n\n`;
        
        if (values.aciklama || template.aciklama_template) {
            preview += `## ${vscode.l10n.t('templateWizard.description')}\n${values.aciklama || template.aciklama_template}\n\n`;
        }

        preview += `## ${vscode.l10n.t('templateWizard.details')}\n`;
        preview += `- **${vscode.l10n.t('templateWizard.priorityLabel')}** ${values.oncelik || vscode.l10n.t('templateWizard.mediumPriority')}\n`;
        if (values.son_tarih) {
            preview += `- **${vscode.l10n.t('templateWizard.dueDateLabel')}** ${values.son_tarih}\n`;
        }
        if (values.etiketler) {
            preview += `- **${vscode.l10n.t('templateWizard.tagsLabel')}** ${values.etiketler}\n`;
        }

        // Add custom fields
        template.alanlar.forEach(field => {
            if (values[field.isim] && !['baslik', 'aciklama', 'oncelik', 'son_tarih', 'etiketler'].includes(field.isim)) {
                preview += `- **${field.isim}:** ${values[field.isim]}\n`;
            }
        });

        return preview;
    }

    private async saveAsFavorite(templateId: string): Promise<void> {
        const favorites = this.getFavorites();
        if (!favorites.includes(templateId)) {
            favorites.push(templateId);
            await this.saveFavorites(favorites);
            
            vscode.window.showInformationMessage(vscode.l10n.t('templateWizard.addedToFavorites'));
            
            // Update UI
            this.panel.webview.postMessage({
                command: 'favoriteAdded',
                templateId: templateId
            });
        }
    }

    private async loadFavorites(): Promise<void> {
        const favoriteIds = this.getFavorites();
        
        try {
            const result = await this.mcpClient.callTool('template_listele');
            const allTemplates = MarkdownParser.parseTemplateListesi(result.content[0].text);
            
            const favorites = allTemplates.filter(t => favoriteIds.includes(t.id));
            
            await this.panel.webview.postMessage({
                command: 'favoritesLoaded',
                templates: favorites
            });
        } catch (error) {
            Logger.error('Failed to load favorite templates:', error);
        }
    }

    private getFavorites(): string[] {
        // For now, return empty array. In real implementation,
        // we would need to pass the extension context to this class
        return [];
    }

    private async saveFavorites(favorites: string[]): Promise<void> {
        // For now, do nothing. In real implementation,
        // we would need to pass the extension context to this class
    }

    private getHtmlContent(): string {
        const scriptUri = this.panel.webview.asWebviewUri(
            vscode.Uri.joinPath(this.extensionUri, 'media', 'templateWizard.js')
        );
        const styleUri = this.panel.webview.asWebviewUri(
            vscode.Uri.joinPath(this.extensionUri, 'media', 'templateWizard.css')
        );

        return `<!DOCTYPE html>
        <html lang="${vscode.env.language}">
        <head>
            <meta charset="UTF-8">
            <meta name="viewport" content="width=device-width, initial-scale=1.0">
            <title>${vscode.l10n.t('templateWizard.title')}</title>
            <link href="${styleUri}" rel="stylesheet">
        </head>
        <body>
            <div class="wizard-container">
                <!-- Step 1: Template Selection -->
                <div id="step-template-selection" class="wizard-step active">
                    <h2>${vscode.l10n.t('templateWizard.selectTemplate')}</h2>
                    
                    <div class="search-container">
                        <input type="text" id="template-search" placeholder="${vscode.l10n.t('templateWizard.searchPlaceholder')}" />
                    </div>

                    <div class="category-tabs">
                        <button class="category-tab active" data-category="">${vscode.l10n.t('templateWizard.categoryAll')}</button>
                        <button class="category-tab" data-category="Genel">${vscode.l10n.t('templateWizard.categoryGeneral')}</button>
                        <button class="category-tab" data-category="Teknik">${vscode.l10n.t('templateWizard.categoryTechnical')}</button>
                        <button class="category-tab" data-category="Özellik">${vscode.l10n.t('templateWizard.categoryFeature')}</button>
                        <button class="category-tab" data-category="Bug">${vscode.l10n.t('templateWizard.categoryBug')}</button>
                        <button class="category-tab" data-category="favorites">${vscode.l10n.t('templateWizard.categoryFavorites')}</button>
                    </div>

                    <div id="template-grid" class="template-grid">
                        <!-- Templates will be loaded here -->
                    </div>
                </div>

                <!-- Step 2: Form Fields -->
                <div id="step-form-fields" class="wizard-step">
                    <h2 id="template-name">${vscode.l10n.t('templateWizard.title')}</h2>
                    <p id="template-description" class="template-description"></p>

                    <form id="template-form">
                        <div id="form-fields">
                            <!-- Dynamic fields will be loaded here -->
                        </div>
                    </form>

                    <div class="form-actions">
                        <button id="btn-back" class="btn-secondary">${vscode.l10n.t('templateWizard.back')}</button>
                        <button id="btn-preview" class="btn-secondary">${vscode.l10n.t('templateWizard.preview')}</button>
                        <button id="btn-create" class="btn-primary">${vscode.l10n.t('templateWizard.create')}</button>
                    </div>
                </div>

                <!-- Step 3: Preview -->
                <div id="step-preview" class="wizard-step">
                    <h2>${vscode.l10n.t('templateWizard.taskPreview')}</h2>
                    
                    <div id="task-preview" class="task-preview">
                        <!-- Preview content will be shown here -->
                    </div>

                    <div class="form-actions">
                        <button id="btn-back-to-form" class="btn-secondary">${vscode.l10n.t('templateWizard.edit')}</button>
                        <button id="btn-confirm-create" class="btn-primary">${vscode.l10n.t('templateWizard.confirmCreate')}</button>
                    </div>
                </div>
            </div>

            <script src="${scriptUri}"></script>
        </body>
        </html>`;
    }

    private dispose(): void {
        TemplateWizard.currentPanel = undefined;

        // Clean up resources
        this.panel.dispose();

        while (this.disposables.length) {
            const disposable = this.disposables.pop();
            if (disposable) {
                disposable.dispose();
            }
        }
    }
}