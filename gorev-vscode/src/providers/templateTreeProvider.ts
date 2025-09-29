import * as vscode from 'vscode';
import { ClientInterface } from '../interfaces/client';
import { ApiClient, ApiError, Template } from '../api/client';
import { GorevTemplate, TemplateKategori } from '../models/template';
import { ICONS, CONTEXT_VALUES } from '../utils/constants';
import { Logger } from '../utils/logger';
import { t } from '../utils/l10n';

export class TemplateTreeProvider implements vscode.TreeDataProvider<TemplateTreeItem | TemplateCategoryItem> {
  private _onDidChangeTreeData = new vscode.EventEmitter<TemplateTreeItem | TemplateCategoryItem | undefined | null | void>();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

  private templates: GorevTemplate[] = [];

  // API Client for REST API calls
  private apiClient: ApiClient;

  constructor(private mcpClient: ClientInterface) {
    // Initialize API client
    this.apiClient = mcpClient instanceof ApiClient ? mcpClient : new ApiClient();
  }

  getTreeItem(element: TemplateTreeItem | TemplateCategoryItem): vscode.TreeItem {
    return element;
  }

  async getChildren(element?: TemplateTreeItem | TemplateCategoryItem): Promise<(TemplateTreeItem | TemplateCategoryItem)[]> {
    console.log('[TemplateTreeProvider] getChildren called with element:', element);
    
    if (!this.mcpClient.isConnected()) {
      console.log('[TemplateTreeProvider] MCP client not connected');
      return [];
    }

    if (!element) {
      // Root level - return categories
      console.log('[TemplateTreeProvider] Loading root categories...');
      try {
        await this.loadTemplates();
        const categories = this.getCategories();
        console.log('[TemplateTreeProvider] Returning', categories.length, 'categories');
        return categories;
      } catch (error) {
        Logger.error('Failed to load templates:', error);
        console.error('[TemplateTreeProvider] Error loading templates:', error);
        return [];
      }
    } else if (element instanceof TemplateCategoryItem) {
      // Return templates for this category
      console.log('[TemplateTreeProvider] Loading templates for category:', element.category);
      const templates = this.templates
        .filter((template) => template.kategori === element.category)
        .map((template) => new TemplateTreeItem(template));
      console.log('[TemplateTreeProvider] Returning', templates.length, 'templates');
      return templates;
    }

    return [];
  }

  async refresh(): Promise<void> {
    await this.loadTemplates();
    this._onDidChangeTreeData.fire();
  }

  private async loadTemplates(): Promise<void> {
    try {
      // Use REST API to get templates
      const response = await this.apiClient.getTemplates();

      if (!response.success || !response.data) {
        Logger.warn('[TemplateTreeProvider] No templates returned from API');
        this.templates = [];
        return;
      }

      // Convert API Template[] to internal GorevTemplate[] model
      this.templates = response.data.map(template => this.convertTemplateToGorevTemplate(template));

      Logger.info(`[TemplateTreeProvider] Loaded ${this.templates.length} templates`);
    } catch (error) {
      if (error instanceof ApiError) {
        Logger.error(`[TemplateTreeProvider] API Error ${error.statusCode}:`, error.apiError);
      } else {
        Logger.error('[TemplateTreeProvider] Failed to load templates:', error);
      }
      throw error;
    }
  }

  private getCategories(): TemplateCategoryItem[] {
    const categories = new Set<TemplateKategori>();
    this.templates.forEach((template) => {
      categories.add(template.kategori);
    });

    return Array.from(categories).map((category) => {
      const count = this.templates.filter((t) => t.kategori === category).length;
      return new TemplateCategoryItem(category, count);
    });
  }

  /**
   * Convert API Template to internal GorevTemplate model
   */
  private convertTemplateToGorevTemplate(template: Template): GorevTemplate {
    return {
      id: template.id,
      isim: template.isim,
      tanim: template.tanim,
      varsayilan_baslik: '', // Not provided by API
      aciklama_template: template.tanim,
      alanlar: template.alanlar.map(field => ({
        isim: field.isim,
        tur: this.mapFieldType(field.tip),
        zorunlu: field.zorunlu,
        varsayilan: field.varsayilan,
        secenekler: field.secenekler,
      })),
      ornek_degerler: {}, // Not provided by current API
      kategori: template.kategori as TemplateKategori,
      aktif: template.aktif,
    };
  }

  /**
   * Map API field type to internal field type
   */
  private mapFieldType(apiType: string): 'metin' | 'sayi' | 'tarih' | 'secim' {
    switch (apiType) {
      case 'text':
        return 'metin';
      case 'number':
        return 'sayi';
      case 'date':
        return 'tarih';
      case 'select':
        return 'secim';
      default:
        return 'metin';
    }
  }

}

export class TemplateCategoryItem extends vscode.TreeItem {
  constructor(
    public readonly category: TemplateKategori,
    public readonly count: number
  ) {
    super(category, vscode.TreeItemCollapsibleState.Expanded);
    
    this.description = t('templateTree.templates', count.toString());
    this.iconPath = new vscode.ThemeIcon(ICONS.PROJECT);
  }
}

export class TemplateTreeItem extends vscode.TreeItem {
  constructor(
    public readonly template: GorevTemplate,
    public readonly collapsibleState: vscode.TreeItemCollapsibleState = vscode.TreeItemCollapsibleState.None
  ) {
    super(template.isim, collapsibleState);
    
    this.tooltip = this.getTooltip();
    this.description = template.tanim;
    this.iconPath = new vscode.ThemeIcon(ICONS.TEMPLATE);
    this.contextValue = CONTEXT_VALUES.TEMPLATE;
  }

  private getTooltip(): string {
    const lines = [
      this.template.isim,
      this.template.tanim,
      '',
      t('templateTree.fields'),
    ];
    
    this.template.alanlar.forEach((field) => {
      const required = field.zorunlu ? ` ${t('templateTree.required')}` : '';
      lines.push(`  • ${field.isim}: ${field.tur}${required}`);
    });
    
    if (this.template.ornek_degerler && Object.keys(this.template.ornek_degerler).length > 0) {
      lines.push('', t('templateTree.exampleValues'));
      Object.entries(this.template.ornek_degerler).forEach(([key, value]) => {
        lines.push(`  • ${key}: ${value}`);
      });
    }
    
    return lines.join('\n');
  }
}
