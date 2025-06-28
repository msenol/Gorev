import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { Proje } from '../models/proje';
import { ICONS, CONTEXT_VALUES } from '../utils/constants';
import { Logger } from '../utils/logger';

export class ProjeTreeProvider implements vscode.TreeDataProvider<ProjeTreeItem> {
  private _onDidChangeTreeData = new vscode.EventEmitter<ProjeTreeItem | undefined | null | void>();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;
  
  private projects: Proje[] = [];
  private activeProjectId: string | null = null;

  constructor(private mcpClient: MCPClient) {}

  getTreeItem(element: ProjeTreeItem): vscode.TreeItem {
    return element;
  }

  async getChildren(element?: ProjeTreeItem): Promise<ProjeTreeItem[]> {
    if (!this.mcpClient.isConnected()) {
      return [];
    }

    if (!element) {
      // Root level - return all projects
      try {
        await this.loadProjects();
        await this.loadActiveProject();
        return this.projects.map((project) => 
          new ProjeTreeItem(project, project.id === this.activeProjectId)
        );
      } catch (error) {
        Logger.error('Failed to load projects:', error);
        return [];
      }
    }

    return [];
  }

  async refresh(): Promise<void> {
    await this.loadProjects();
    await this.loadActiveProject();
    this._onDidChangeTreeData.fire();
  }

  private async loadProjects(): Promise<void> {
    try {
      const result = await this.mcpClient.callTool('proje_listele');
      
      // Parse the markdown content to extract projects
      this.projects = this.parseProjectsFromContent(result.content[0].text);
    } catch (error) {
      Logger.error('Failed to load projects:', error);
      throw error;
    }
  }

  private async loadActiveProject(): Promise<void> {
    try {
      const result = await this.mcpClient.callTool('aktif_proje_goster');
      const content = result.content[0].text;
      
      // Check if no active project
      if (content.includes('Henüz aktif proje ayarlanmamış')) {
        this.activeProjectId = null;
        return;
      }
      
      // Parse the active project ID from the markdown
      const lines = content.split('\n');
      for (const line of lines) {
        const match = line.match(/\*\*ID:\*\*\s*(.+)/);
        if (match) {
          this.activeProjectId = match[1].trim();
          break;
        }
      }
    } catch (error) {
      // No active project is ok
      this.activeProjectId = null;
    }
  }

  private parseProjectsFromContent(content: string): Proje[] {
    const projects: Proje[] = [];
    
    // Check for empty project list
    if (content.includes('Henüz proje bulunmuyor')) {
      return projects;
    }
    
    // Split content by project headers (###)
    const projectSections = content.split(/^###\s+/m).filter(section => section.trim());
    
    for (const section of projectSections) {
      const lines = section.trim().split('\n');
      if (lines.length === 0) continue;
      
      // First line is the project name
      const projectName = lines[0].trim();
      const project: Partial<Proje> = {
        isim: projectName,
      };
      
      // Parse the rest of the lines for properties
      for (let i = 1; i < lines.length; i++) {
        const line = lines[i].trim();
        if (!line || !line.startsWith('- **')) continue;
        
        // Extract property name and value
        const match = line.match(/- \*\*([^:]+):\*\*\s*(.+)/);
        if (!match) continue;
        
        const [, property, value] = match;
        
        switch (property) {
          case 'ID':
            project.id = value.trim();
            break;
          case 'Tanım':
            project.tanim = value.trim();
            break;
          case 'Oluşturma':
            project.olusturma_tarih = value.trim();
            break;
          case 'Görev Sayısı':
            project.gorev_sayisi = parseInt(value.trim()) || 0;
            break;
        }
      }
      
      // Only add project if it has an ID
      if (project.id) {
        projects.push(project as Proje);
      }
    }
    
    return projects;
  }
}

export class ProjeTreeItem extends vscode.TreeItem {
  constructor(
    public readonly project: Proje,
    public readonly isActive: boolean = false,
    public readonly collapsibleState: vscode.TreeItemCollapsibleState = vscode.TreeItemCollapsibleState.None
  ) {
    super(project.isim, collapsibleState);
    
    this.tooltip = this.getTooltip();
    this.description = this.getDescription();
    this.iconPath = this.getIcon();
    this.contextValue = this.isActive ? CONTEXT_VALUES.PROJECT_ACTIVE : CONTEXT_VALUES.PROJECT;
  }

  private getTooltip(): string {
    const lines = [
      this.project.isim,
      this.project.tanim || 'No description',
      '',
      `Total tasks: ${this.project.gorev_sayisi || 0}`,
    ];
    
    if (this.project.tamamlanan_sayisi !== undefined) {
      lines.push(`Completed: ${this.project.tamamlanan_sayisi}`);
    }
    if (this.project.devam_eden_sayisi !== undefined) {
      lines.push(`In progress: ${this.project.devam_eden_sayisi}`);
    }
    if (this.project.bekleyen_sayisi !== undefined) {
      lines.push(`Pending: ${this.project.bekleyen_sayisi}`);
    }
    
    if (this.isActive) {
      lines.push('', '✓ Active project');
    }
    
    return lines.join('\n');
  }

  private getDescription(): string {
    const parts = [];
    
    if (this.isActive) {
      parts.push('(active)');
    }
    
    const taskCount = this.project.gorev_sayisi || 0;
    parts.push(`${taskCount} tasks`);
    
    if (this.project.tamamlanan_sayisi !== undefined && taskCount > 0) {
      const percentage = Math.round((this.project.tamamlanan_sayisi / taskCount) * 100);
      parts.push(`${percentage}% done`);
    }
    
    return parts.join(' • ');
  }

  private getIcon(): vscode.ThemeIcon {
    const iconId = this.isActive ? ICONS.PROJECT_ACTIVE : ICONS.PROJECT;
    return new vscode.ThemeIcon(iconId);
  }
}