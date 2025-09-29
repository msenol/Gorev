import * as vscode from 'vscode';
import { t } from '../utils/l10n';
import { ClientInterface } from '../interfaces/client';
import { ApiClient, ApiError, Project } from '../api/client';
import { Proje } from '../models/proje';
import { ICONS, CONTEXT_VALUES } from '../utils/constants';
import { Logger } from '../utils/logger';

export class ProjeTreeProvider implements vscode.TreeDataProvider<ProjeTreeItem> {
  private _onDidChangeTreeData = new vscode.EventEmitter<ProjeTreeItem | undefined | null | void>();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;

  private projects: Proje[] = [];
  private activeProjectId: string | null = null;

  // API Client for REST API calls
  private apiClient: ApiClient;

  constructor(private mcpClient: ClientInterface) {
    // Initialize API client
    this.apiClient = mcpClient instanceof ApiClient ? mcpClient : new ApiClient();
  }

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
      // Use REST API to get projects
      const response = await this.apiClient.getProjects();

      if (!response.success || !response.data) {
        Logger.warn('[ProjeTreeProvider] No projects returned from API');
        this.projects = [];
        return;
      }

      // Convert API Project[] to internal Proje[] model
      this.projects = response.data.map(project => this.convertProjectToProje(project));

      Logger.info(`[ProjeTreeProvider] Loaded ${this.projects.length} projects`);
    } catch (error) {
      if (error instanceof ApiError) {
        Logger.error(`[ProjeTreeProvider] API Error ${error.statusCode}:`, error.apiError);
      } else {
        Logger.error('[ProjeTreeProvider] Failed to load projects:', error);
      }
      throw error;
    }
  }

  private async loadActiveProject(): Promise<void> {
    try {
      // Use REST API to get active project
      const response = await this.apiClient.getActiveProject();

      if (!response.success || !response.data) {
        // No active project
        this.activeProjectId = null;
        return;
      }

      // Extract project ID from response
      this.activeProjectId = response.data.id;
      Logger.info(`[ProjeTreeProvider] Active project: ${this.activeProjectId}`);
    } catch (error) {
      if (error instanceof ApiError && error.isNotFound()) {
        // 404 means no active project - this is ok
        this.activeProjectId = null;
      } else {
        Logger.warn('[ProjeTreeProvider] Failed to load active project:', error);
        this.activeProjectId = null;
      }
    }
  }

  /**
   * Convert API Project to internal Proje model
   */
  private convertProjectToProje(project: Project): Proje {
    return {
      id: project.id,
      isim: project.isim,
      tanim: project.tanim || '',
      gorev_sayisi: project.gorev_sayisi,
      olusturma_tarih: project.olusturma_tarihi,
      guncelleme_tarih: project.olusturma_tarihi, // API doesn't return update date for projects
      // Task statistics not available from API Project type yet
      // These would need to be fetched separately or added to API response
    };
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
