import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { GorevTemplate, TemplateKategori } from '../models/template';
import { ICONS, CONTEXT_VALUES } from '../utils/constants';
import { Logger } from '../utils/logger';
import { MarkdownParser } from '../utils/markdownParser';

export class TemplateTreeProvider implements vscode.TreeDataProvider<TemplateTreeItem | TemplateCategoryItem> {
  private _onDidChangeTreeData = new vscode.EventEmitter<TemplateTreeItem | TemplateCategoryItem | undefined | null | void>();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;
  
  private templates: GorevTemplate[] = [];

  constructor(private mcpClient: MCPClient) {}

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
      const result = await this.mcpClient.callTool('template_listele');
      
      // Debug: Log raw response
      console.log('[TemplateTreeProvider] Raw MCP response:', result);
      console.log('[TemplateTreeProvider] Content text:', result.content[0].text);
      
      // Parse the markdown content to extract templates
      this.templates = MarkdownParser.parseTemplateListesi(result.content[0].text);
      
      // Debug: Log parsed templates
      console.log('[TemplateTreeProvider] Parsed templates count:', this.templates.length);
      console.log('[TemplateTreeProvider] Parsed templates:', this.templates);
    } catch (error) {
      Logger.error('Failed to load templates:', error);
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

}

export class TemplateCategoryItem extends vscode.TreeItem {
  constructor(
    public readonly category: TemplateKategori,
    public readonly count: number
  ) {
    super(category, vscode.TreeItemCollapsibleState.Expanded);
    
    this.description = `${count} templates`;
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
      'Fields:',
    ];
    
    this.template.alanlar.forEach((field) => {
      const required = field.zorunlu ? ' (required)' : '';
      lines.push(`  • ${field.isim}: ${field.tur}${required}`);
    });
    
    if (this.template.ornek_degerler && Object.keys(this.template.ornek_degerler).length > 0) {
      lines.push('', 'Example values:');
      Object.entries(this.template.ornek_degerler).forEach(([key, value]) => {
        lines.push(`  • ${key}: ${value}`);
      });
    }
    
    return lines.join('\n');
  }
}