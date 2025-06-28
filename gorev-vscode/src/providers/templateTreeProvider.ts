import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { GorevTemplate, TemplateKategori } from '../models/template';
import { ICONS, CONTEXT_VALUES } from '../utils/constants';
import { Logger } from '../utils/logger';

export class TemplateTreeProvider implements vscode.TreeDataProvider<TemplateTreeItem | TemplateCategoryItem> {
  private _onDidChangeTreeData = new vscode.EventEmitter<TemplateTreeItem | TemplateCategoryItem | undefined | null | void>();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;
  
  private templates: GorevTemplate[] = [];

  constructor(private mcpClient: MCPClient) {}

  getTreeItem(element: TemplateTreeItem | TemplateCategoryItem): vscode.TreeItem {
    return element;
  }

  async getChildren(element?: TemplateTreeItem | TemplateCategoryItem): Promise<(TemplateTreeItem | TemplateCategoryItem)[]> {
    if (!this.mcpClient.isConnected()) {
      return [];
    }

    if (!element) {
      // Root level - return categories
      try {
        await this.loadTemplates();
        return this.getCategories();
      } catch (error) {
        Logger.error('Failed to load templates:', error);
        return [];
      }
    } else if (element instanceof TemplateCategoryItem) {
      // Return templates for this category
      return this.templates
        .filter((template) => template.kategori === element.category)
        .map((template) => new TemplateTreeItem(template));
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
      
      // Parse the markdown content to extract templates
      this.templates = this.parseTemplatesFromContent(result.content[0].text);
    } catch (error) {
      Logger.error('Failed to load templates:', error);
      throw error;
    }
  }

  private getCategories(): TemplateCategoryItem[] {
    const categories = new Set<TemplateKategori>();
    this.templates.forEach((template) => {
      if (template.aktif) {
        categories.add(template.kategori);
      }
    });

    return Array.from(categories).map((category) => {
      const count = this.templates.filter((t) => t.kategori === category && t.aktif).length;
      return new TemplateCategoryItem(category, count);
    });
  }

  private parseTemplatesFromContent(content: string): GorevTemplate[] {
    const templates: GorevTemplate[] = [];
    
    // Check for empty template list
    if (content.includes('Henüz template bulunmuyor')) {
      return templates;
    }
    
    const lines = content.split('\n');
    let currentCategory: TemplateKategori | null = null;
    let currentTemplate: Partial<GorevTemplate> | null = null;
    let inFieldsSection = false;
    
    for (let i = 0; i < lines.length; i++) {
      const line = lines[i];
      
      // Check for category (### header)
      if (line.startsWith('### ')) {
        // Save previous template if exists
        if (currentTemplate && currentTemplate.id) {
          templates.push(currentTemplate as GorevTemplate);
        }
        
        currentCategory = line.substring(4).trim() as TemplateKategori;
        currentTemplate = null;
        inFieldsSection = false;
        continue;
      }
      
      // Check for template name (#### header)
      if (line.startsWith('#### ')) {
        // Save previous template if exists
        if (currentTemplate && currentTemplate.id) {
          templates.push(currentTemplate as GorevTemplate);
        }
        
        currentTemplate = {
          isim: line.substring(5).trim(),
          kategori: currentCategory || TemplateKategori.Genel,
          alanlar: [],
          aktif: true,
        };
        inFieldsSection = false;
        continue;
      }
      
      // Parse template details
      if (currentTemplate && line.trim().startsWith('- **')) {
        const idMatch = line.match(/- \*\*ID:\*\* `([^`]+)`/);
        if (idMatch) {
          currentTemplate.id = idMatch[1];
          continue;
        }
        
        const descMatch = line.match(/- \*\*Açıklama:\*\* (.+)/);
        if (descMatch) {
          currentTemplate.tanim = descMatch[1].trim();
          continue;
        }
        
        const titleMatch = line.match(/- \*\*Başlık Şablonu:\*\* `([^`]+)`/);
        if (titleMatch) {
          currentTemplate.varsayilan_baslik = titleMatch[1];
          continue;
        }
        
        if (line.includes('- **Alanlar:**')) {
          inFieldsSection = true;
          continue;
        }
      }
      
      // Parse fields
      if (inFieldsSection && currentTemplate && line.trim().startsWith('- `')) {
        const fieldMatch = line.match(/- `([^`]+)` \(([^)]+)\)(.+)?/);
        if (fieldMatch) {
          const [, fieldName, fieldType, extra] = fieldMatch;
          const field: any = {
            isim: fieldName,
            tur: fieldType,
            zorunlu: false,
            varsayilan: '',
          };
          
          if (extra) {
            // Check if required
            if (extra.includes('*(zorunlu)*')) {
              field.zorunlu = true;
            }
            
            // Extract default value
            const defaultMatch = extra.match(/varsayılan: ([^-]+)/);
            if (defaultMatch) {
              field.varsayilan = defaultMatch[1].trim();
            }
            
            // Extract options
            const optionsMatch = extra.match(/seçenekler: (.+)/);
            if (optionsMatch) {
              field.secenekler = optionsMatch[1].split(',').map(opt => opt.trim());
            }
          }
          
          if (!currentTemplate.alanlar) {
            currentTemplate.alanlar = [];
          }
          currentTemplate.alanlar.push(field);
        }
      }
    }
    
    // Don't forget the last template
    if (currentTemplate && currentTemplate.id) {
      templates.push(currentTemplate as GorevTemplate);
    }
    
    return templates;
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