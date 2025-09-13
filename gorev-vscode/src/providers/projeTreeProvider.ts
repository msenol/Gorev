import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import { MCPClient } from '../mcp/client';
import { Proje } from '../models/proje';
import { ICONS, CONTEXT_VALUES } from '../utils/constants';
import { Logger } from '../utils/logger';
import { MarkdownParser } from '../utils/markdownParser';

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
    console.log('[ProjeTreeProvider] getChildren called with element:', element);
    
    if (!this.mcpClient.isConnected()) {
      console.log('[ProjeTreeProvider] MCP client not connected');
      return [];
    }

    if (!element) {
      // Root level - return all projects
      console.log('[ProjeTreeProvider] Loading root projects...');
      try {
        await this.loadProjects();
        await this.loadActiveProject();
        const items = this.projects.map((project) => 
          new ProjeTreeItem(project, project.id === this.activeProjectId)
        );
        console.log('[ProjeTreeProvider] Returning', items.length, 'project items');
        return items;
      } catch (error) {
        Logger.error('Failed to load projects:', error);
        console.error('[ProjeTreeProvider] Error loading projects:', error);
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
      
      // Debug: Log raw response
      console.log('[ProjeTreeProvider] Raw MCP response:', result);
      console.log('[ProjeTreeProvider] Content text:', result.content[0].text);
      
      // Parse the markdown content to extract projects
      this.projects = MarkdownParser.parseProjeListesi(result.content[0].text);
      
      // Debug: Log parsed projects
      console.log('[ProjeTreeProvider] Parsed projects count:', this.projects.length);
      console.log('[ProjeTreeProvider] Parsed projects:', this.projects);
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
      if (content.includes(t('projectTree.noActiveProject'))) {
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
      this.project.tanim || t('projectTree.noDescription'),
      '',
      t('projectTree.totalTasks', (this.project.gorev_sayisi || 0).toString()),
    ];
    
    if (this.project.tamamlanan_sayisi !== undefined) {
      lines.push(t('projectTree.completed', this.project.tamamlanan_sayisi.toString()));
    }
    if (this.project.devam_eden_sayisi !== undefined) {
      lines.push(t('projectTree.inProgress', this.project.devam_eden_sayisi.toString()));
    }
    if (this.project.bekleyen_sayisi !== undefined) {
      lines.push(t('projectTree.pending', this.project.bekleyen_sayisi.toString()));
    }
    
    if (this.isActive) {
      lines.push('', t('projectTree.activeProject'));
    }
    
    return lines.join('\n');
  }

  private getDescription(): string {
    const parts = [];
    
    if (this.isActive) {
      parts.push(t('projectTree.active'));
    }
    
    const taskCount = this.project.gorev_sayisi || 0;
    parts.push(t('projectTree.tasks', taskCount.toString()));
    
    if (this.project.tamamlanan_sayisi !== undefined && taskCount > 0) {
      const percentage = Math.round((this.project.tamamlanan_sayisi / taskCount) * 100);
      parts.push(t('projectTree.percentDone', percentage.toString()));
    }
    
    return parts.join(' • ');
  }

  private getIcon(): vscode.ThemeIcon {
    const iconId = this.isActive ? ICONS.PROJECT_ACTIVE : ICONS.PROJECT;
    return new vscode.ThemeIcon(iconId);
  }
}
