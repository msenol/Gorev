import * as vscode from 'vscode';
import { ApiClient, ApiError } from '../api/client';
import { EnhancedGorevTreeProvider } from '../providers/enhancedGorevTreeProvider';
import { ProjeTreeProvider } from '../providers/projeTreeProvider';
import { TemplateTreeProvider } from '../providers/templateTreeProvider';
import { StatusBarManager } from '../ui/statusBar';
import { FilterToolbar } from '../ui/filterToolbar';
import { COMMANDS } from '../utils/constants';
import { Logger } from '../utils/logger';
import { t } from '../utils/l10n';
import { registerGorevCommands } from './gorevCommands';
import { registerProjeCommands } from './projeCommands';
import { registerTemplateCommands } from './templateCommands';
import { registerEnhancedGorevCommands } from './enhancedGorevCommands';
import { registerInlineEditCommands } from './inlineEditCommands';
import { registerFilterCommands } from './filterCommands';
import { registerDataCommands } from './dataCommands';
import { registerDatabaseCommands } from './databaseCommands';

export interface CommandContext {
  gorevTreeProvider: EnhancedGorevTreeProvider;
  projeTreeProvider: ProjeTreeProvider;
  templateTreeProvider: TemplateTreeProvider;
  statusBarManager: StatusBarManager;
  filterToolbar?: FilterToolbar;
}

export function registerCommands(
  context: vscode.ExtensionContext,
  apiClient: ApiClient,
  providers: CommandContext
): void {
  // Register all command groups
  registerGorevCommands(context, apiClient, providers);
  registerProjeCommands(context, apiClient, providers);
  registerTemplateCommands(context, apiClient, providers);
  registerEnhancedGorevCommands(context, apiClient, providers);
  registerInlineEditCommands(context, apiClient, providers);
  registerDataCommands(context, apiClient, providers);
  registerDatabaseCommands(context, apiClient, providers);
  
  if (providers.filterToolbar) {
    registerFilterCommands(context, apiClient, providers);
  }

  // Initialize API client for general commands

  // Register general commands
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.SHOW_SUMMARY, async () => {
      try {
        if (!apiClient.isConnected()) {
          vscode.window.showWarningMessage('Not connected to Gorev server');
          return;
        }

        // Use REST API to get summary
        const response = await apiClient.getSummary();

        if (!response.success || !response.data) {
          vscode.window.showErrorMessage('Failed to get summary from server');
          return;
        }

        const summaryPanel = vscode.window.createWebviewPanel(
          'gorevSummary',
          'Gorev Summary',
          vscode.ViewColumn.One,
          {}
        );

        // Generate rich HTML dashboard
        summaryPanel.webview.html = getSummaryHtml(response.data as SummaryData);
      } catch (error) {
        if (error instanceof ApiError) {
          Logger.error(`[ShowSummary] API Error ${error.statusCode}:`, error.apiError);
          vscode.window.showErrorMessage(`Failed to show summary: ${error.apiError}`);
        } else {
          vscode.window.showErrorMessage(`Failed to show summary: ${error}`);
        }
      }
    })
  );

  // Connect command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.CONNECT, async () => {
      try {
        await connectToServer(apiClient, providers);
      } catch (error) {
        vscode.window.showErrorMessage(`Failed to connect: ${error}`);
      }
    })
  );

  // Disconnect command
  context.subscriptions.push(
    vscode.commands.registerCommand(COMMANDS.DISCONNECT, () => {
      apiClient.disconnect();
      providers.statusBarManager.setDisconnected();
      vscode.window.showInformationMessage('Disconnected from Gorev server');
    })
  );
}

async function connectToServer(client: ApiClient, providers: CommandContext): Promise<void> {
  providers.statusBarManager.setConnecting();

  try {
    // Use the unified client's connect method which handles auto-detection
    await client.connect();
    providers.statusBarManager.setConnected();

    // Refresh all views after connection - sequentially to avoid overwhelming the server
    await providers.gorevTreeProvider.refresh();
    await providers.projeTreeProvider.refresh();
    await providers.templateTreeProvider.refresh();

    vscode.window.showInformationMessage('Connected to Gorev server');
  } catch (error) {
    providers.statusBarManager.setDisconnected();
    throw error;
  }
}

// Summary data types
interface SummaryTask {
  id: string;
  title?: string;
  baslik?: string;
  status: string;
  priority: string;
}

interface SummaryData {
  total_tasks?: number;
  total_projects?: number;
  total_templates?: number;
  completion_rate?: number;
  status_counts?: { pending: number; in_progress: number; completed: number };
  priority_counts?: { high: number; medium: number; low: number };
  due_date_summary?: { overdue: number; due_today: number; due_this_week: number };
  active_project?: { name: string; definition?: string };
  high_priority_tasks?: SummaryTask[];
  overdue_tasks?: SummaryTask[];
  recent_tasks?: SummaryTask[];
}

/**
 * Generate rich HTML summary dashboard from API response
 */
function getSummaryHtml(data: SummaryData): string {
  const totalTasks = data.total_tasks || 0;
  const totalProjects = data.total_projects || 0;
  const totalTemplates = data.total_templates || 0;
  const completionRate = Math.round(data.completion_rate || 0);
  
  const statusCounts = data.status_counts || { pending: 0, in_progress: 0, completed: 0 };
  const priorityCounts = data.priority_counts || { high: 0, medium: 0, low: 0 };
  const dueDateSummary = data.due_date_summary || { overdue: 0, due_today: 0, due_this_week: 0 };
  
  const activeProject = data.active_project;
  const highPriorityTasks = data.high_priority_tasks || [];
  const recentTasks = data.recent_tasks || [];

  // Calculate progress bar widths
  const statusTotal = statusCounts.pending + statusCounts.in_progress + statusCounts.completed;
  const pendingWidth = statusTotal > 0 ? (statusCounts.pending / statusTotal) * 100 : 0;
  const inProgressWidth = statusTotal > 0 ? (statusCounts.in_progress / statusTotal) * 100 : 0;
  const completedWidth = statusTotal > 0 ? (statusCounts.completed / statusTotal) * 100 : 0;

  // Generate task list HTML
  const renderTaskList = (tasks: SummaryTask[], emptyMessage: string): string => {
    if (!tasks || tasks.length === 0) {
      return `<div class="empty-state">${emptyMessage}</div>`;
    }
    return tasks.map(task => `
      <div class="task-item">
        <span class="task-status ${task.status}">${getStatusIcon(task.status)}</span>
        <span class="task-title">${escapeHtml(task.title || task.baslik || 'Untitled')}</span>
        <span class="task-priority ${task.priority}">${getPriorityBadge(task.priority)}</span>
      </div>
    `).join('');
  };

  const now = new Date();
  const dateStr = now.toLocaleDateString('tr-TR', { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' });

  return `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        :root {
            --accent-primary: #6366f1;
            --accent-secondary: #8b5cf6;
            --success: #10b981;
            --warning: #f59e0b;
            --danger: #ef4444;
            --info: #3b82f6;
        }
        
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        
        body {
            font-family: var(--vscode-font-family);
            color: var(--vscode-foreground);
            background: var(--vscode-editor-background);
            padding: 24px;
            line-height: 1.6;
        }
        
        .dashboard {
            max-width: 1200px;
            margin: 0 auto;
        }
        
        /* Header */
        .header {
            display: flex;
            align-items: center;
            gap: 16px;
            margin-bottom: 32px;
            padding-bottom: 20px;
            border-bottom: 1px solid var(--vscode-panel-border);
        }
        
        .header-icon {
            font-size: 40px;
            filter: drop-shadow(0 4px 6px rgba(0,0,0,0.3));
        }
        
        .header-content h1 {
            font-size: 28px;
            font-weight: 700;
            background: linear-gradient(135deg, var(--accent-primary), var(--accent-secondary));
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }
        
        .header-content .subtitle {
            color: var(--vscode-descriptionForeground);
            font-size: 14px;
        }
        
        /* Stats Grid */
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
            gap: 16px;
            margin-bottom: 32px;
        }
        
        .stat-card {
            background: var(--vscode-input-background);
            border: 1px solid var(--vscode-panel-border);
            border-radius: 12px;
            padding: 20px;
            text-align: center;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        
        .stat-card:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 25px rgba(0,0,0,0.2);
        }
        
        .stat-icon {
            font-size: 32px;
            margin-bottom: 8px;
        }
        
        .stat-value {
            font-size: 36px;
            font-weight: 800;
            color: var(--accent-primary);
        }
        
        .stat-label {
            font-size: 12px;
            color: var(--vscode-descriptionForeground);
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        
        /* Completion Ring */
        .completion-card {
            background: linear-gradient(135deg, var(--accent-primary), var(--accent-secondary));
            color: white;
        }
        
        .completion-card .stat-value {
            color: white;
        }
        
        .completion-card .stat-label {
            color: rgba(255,255,255,0.8);
        }
        
        /* Main Content Grid */
        .content-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
            gap: 24px;
            margin-bottom: 32px;
        }
        
        .card {
            background: var(--vscode-input-background);
            border: 1px solid var(--vscode-panel-border);
            border-radius: 12px;
            padding: 20px;
        }
        
        .card-header {
            display: flex;
            align-items: center;
            gap: 10px;
            margin-bottom: 16px;
            padding-bottom: 12px;
            border-bottom: 1px solid var(--vscode-panel-border);
        }
        
        .card-header h2 {
            font-size: 16px;
            font-weight: 600;
        }
        
        .card-header .icon {
            font-size: 20px;
        }
        
        /* Status Distribution */
        .status-bar {
            display: flex;
            height: 12px;
            border-radius: 6px;
            overflow: hidden;
            margin-bottom: 16px;
            background: var(--vscode-panel-border);
        }
        
        .status-segment {
            transition: width 0.5s ease;
        }
        
        .status-segment.pending { background: var(--warning); }
        .status-segment.in-progress { background: var(--info); }
        .status-segment.completed { background: var(--success); }
        
        .status-legend {
            display: flex;
            flex-wrap: wrap;
            gap: 16px;
        }
        
        .legend-item {
            display: flex;
            align-items: center;
            gap: 8px;
            font-size: 13px;
        }
        
        .legend-dot {
            width: 12px;
            height: 12px;
            border-radius: 50%;
        }
        
        .legend-dot.pending { background: var(--warning); }
        .legend-dot.in-progress { background: var(--info); }
        .legend-dot.completed { background: var(--success); }
        
        .legend-count {
            font-weight: 600;
            color: var(--vscode-foreground);
        }
        
        /* Priority Distribution */
        .priority-bars {
            display: flex;
            flex-direction: column;
            gap: 12px;
        }
        
        .priority-row {
            display: flex;
            align-items: center;
            gap: 12px;
        }
        
        .priority-label {
            width: 80px;
            font-size: 13px;
            display: flex;
            align-items: center;
            gap: 6px;
        }
        
        .priority-bar-container {
            flex: 1;
            height: 24px;
            background: var(--vscode-panel-border);
            border-radius: 4px;
            overflow: hidden;
        }
        
        .priority-bar {
            height: 100%;
            border-radius: 4px;
            transition: width 0.5s ease;
            display: flex;
            align-items: center;
            justify-content: flex-end;
            padding-right: 8px;
            font-size: 12px;
            font-weight: 600;
            color: white;
            min-width: 30px;
        }
        
        .priority-bar.high { background: var(--danger); }
        .priority-bar.medium { background: var(--warning); }
        .priority-bar.low { background: var(--info); }
        
        /* Due Date Cards */
        .due-date-grid {
            display: grid;
            grid-template-columns: repeat(3, 1fr);
            gap: 12px;
        }
        
        .due-card {
            text-align: center;
            padding: 16px 8px;
            border-radius: 8px;
            background: var(--vscode-editor-background);
        }
        
        .due-card.overdue { border-left: 3px solid var(--danger); }
        .due-card.today { border-left: 3px solid var(--warning); }
        .due-card.week { border-left: 3px solid var(--info); }
        
        .due-count {
            font-size: 28px;
            font-weight: 700;
        }
        
        .due-card.overdue .due-count { color: var(--danger); }
        .due-card.today .due-count { color: var(--warning); }
        .due-card.week .due-count { color: var(--info); }
        
        .due-label {
            font-size: 11px;
            color: var(--vscode-descriptionForeground);
            text-transform: uppercase;
        }
        
        /* Task Lists */
        .task-list {
            display: flex;
            flex-direction: column;
            gap: 8px;
            max-height: 300px;
            overflow-y: auto;
        }
        
        .task-item {
            display: flex;
            align-items: center;
            gap: 10px;
            padding: 10px 12px;
            background: var(--vscode-editor-background);
            border-radius: 8px;
            transition: background 0.2s;
        }
        
        .task-item:hover {
            background: var(--vscode-list-hoverBackground);
        }
        
        .task-status {
            font-size: 14px;
        }
        
        .task-title {
            flex: 1;
            font-size: 13px;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
        }
        
        .task-priority {
            font-size: 11px;
            padding: 2px 8px;
            border-radius: 10px;
            font-weight: 600;
        }
        
        .task-priority.high, .task-priority.yuksek { 
            background: rgba(239, 68, 68, 0.2); 
            color: var(--danger); 
        }
        .task-priority.medium, .task-priority.orta { 
            background: rgba(245, 158, 11, 0.2); 
            color: var(--warning); 
        }
        .task-priority.low, .task-priority.dusuk { 
            background: rgba(59, 130, 246, 0.2); 
            color: var(--info); 
        }
        
        .empty-state {
            text-align: center;
            padding: 24px;
            color: var(--vscode-descriptionForeground);
            font-style: italic;
        }
        
        /* Active Project */
        .active-project {
            display: flex;
            align-items: center;
            gap: 12px;
            padding: 16px;
            background: linear-gradient(135deg, rgba(99, 102, 241, 0.1), rgba(139, 92, 246, 0.1));
            border: 1px solid var(--accent-primary);
            border-radius: 8px;
        }
        
        .project-icon {
            font-size: 24px;
        }
        
        .project-name {
            font-weight: 600;
            font-size: 15px;
        }
        
        .project-desc {
            font-size: 12px;
            color: var(--vscode-descriptionForeground);
            margin-top: 4px;
        }
        
        /* Footer */
        .footer {
            text-align: center;
            padding-top: 20px;
            border-top: 1px solid var(--vscode-panel-border);
            color: var(--vscode-descriptionForeground);
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div class="dashboard">
        <!-- Header -->
        <div class="header">
            <span class="header-icon">📊</span>
            <div class="header-content">
                <h1>${t('summary.title')}</h1>
                <div class="subtitle">${t('summary.subtitle')}</div>
            </div>
        </div>
        
        <!-- Stats Grid -->
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-icon">📋</div>
                <div class="stat-value">${totalTasks}</div>
                <div class="stat-label">${t('summary.totalTasks')}</div>
            </div>
            <div class="stat-card">
                <div class="stat-icon">📁</div>
                <div class="stat-value">${totalProjects}</div>
                <div class="stat-label">${t('summary.projects')}</div>
            </div>
            <div class="stat-card">
                <div class="stat-icon">📄</div>
                <div class="stat-value">${totalTemplates}</div>
                <div class="stat-label">${t('summary.templates')}</div>
            </div>
            <div class="stat-card completion-card">
                <div class="stat-icon">🎯</div>
                <div class="stat-value">${completionRate}%</div>
                <div class="stat-label">${t('summary.completion')}</div>
            </div>
        </div>
        
        <!-- Main Content -->
        <div class="content-grid">
            <!-- Status Distribution -->
            <div class="card">
                <div class="card-header">
                    <span class="icon">📈</span>
                    <h2>${t('summary.statusDistribution')}</h2>
                </div>
                <div class="status-bar">
                    <div class="status-segment pending" style="width: ${pendingWidth}%"></div>
                    <div class="status-segment in-progress" style="width: ${inProgressWidth}%"></div>
                    <div class="status-segment completed" style="width: ${completedWidth}%"></div>
                </div>
                <div class="status-legend">
                    <div class="legend-item">
                        <span class="legend-dot pending"></span>
                        <span>${t('summary.pending')}</span>
                        <span class="legend-count">${statusCounts.pending}</span>
                    </div>
                    <div class="legend-item">
                        <span class="legend-dot in-progress"></span>
                        <span>${t('summary.inProgress')}</span>
                        <span class="legend-count">${statusCounts.in_progress}</span>
                    </div>
                    <div class="legend-item">
                        <span class="legend-dot completed"></span>
                        <span>${t('summary.completed')}</span>
                        <span class="legend-count">${statusCounts.completed}</span>
                    </div>
                </div>
            </div>
            
            <!-- Priority Distribution -->
            <div class="card">
                <div class="card-header">
                    <span class="icon">🎚️</span>
                    <h2>${t('summary.priorityDistribution')}</h2>
                </div>
                <div class="priority-bars">
                    <div class="priority-row">
                        <span class="priority-label">🔥 ${t('summary.high')}</span>
                        <div class="priority-bar-container">
                            <div class="priority-bar high" style="width: ${totalTasks > 0 ? Math.max((priorityCounts.high / totalTasks) * 100, 10) : 0}%">
                                ${priorityCounts.high}
                            </div>
                        </div>
                    </div>
                    <div class="priority-row">
                        <span class="priority-label">⚡ ${t('summary.medium')}</span>
                        <div class="priority-bar-container">
                            <div class="priority-bar medium" style="width: ${totalTasks > 0 ? Math.max((priorityCounts.medium / totalTasks) * 100, 10) : 0}%">
                                ${priorityCounts.medium}
                            </div>
                        </div>
                    </div>
                    <div class="priority-row">
                        <span class="priority-label">💤 ${t('summary.low')}</span>
                        <div class="priority-bar-container">
                            <div class="priority-bar low" style="width: ${totalTasks > 0 ? Math.max((priorityCounts.low / totalTasks) * 100, 10) : 0}%">
                                ${priorityCounts.low}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            
            <!-- Due Date Summary -->
            <div class="card">
                <div class="card-header">
                    <span class="icon">📅</span>
                    <h2>${t('summary.byDate')}</h2>
                </div>
                <div class="due-date-grid">
                    <div class="due-card overdue">
                        <div class="due-count">${dueDateSummary.overdue}</div>
                        <div class="due-label">${t('summary.overdue')}</div>
                    </div>
                    <div class="due-card today">
                        <div class="due-count">${dueDateSummary.due_today}</div>
                        <div class="due-label">${t('summary.dueToday')}</div>
                    </div>
                    <div class="due-card week">
                        <div class="due-count">${dueDateSummary.due_this_week}</div>
                        <div class="due-label">${t('summary.dueThisWeek')}</div>
                    </div>
                </div>
            </div>
            
            <!-- Active Project -->
            <div class="card">
                <div class="card-header">
                    <span class="icon">🎯</span>
                    <h2>${t('summary.activeProject')}</h2>
                </div>
                ${activeProject ? `
                <div class="active-project">
                    <span class="project-icon">📂</span>
                    <div>
                        <div class="project-name">${escapeHtml(activeProject.name)}</div>
                        ${activeProject.definition ? `<div class="project-desc">${escapeHtml(activeProject.definition).substring(0, 100)}...</div>` : ''}
                    </div>
                </div>
                ` : `<div class="empty-state">${t('summary.noActiveProject')}</div>`}
            </div>
        </div>
        
        <!-- Task Lists -->
        <div class="content-grid">
            <!-- High Priority Tasks -->
            <div class="card">
                <div class="card-header">
                    <span class="icon">🔥</span>
                    <h2>${t('summary.highPriorityTasks')}</h2>
                </div>
                <div class="task-list">
                    ${renderTaskList(highPriorityTasks, t('summary.noHighPriorityTasks'))}
                </div>
            </div>
            
            <!-- Recent Tasks -->
            <div class="card">
                <div class="card-header">
                    <span class="icon">🕐</span>
                    <h2>${t('summary.recentTasks')}</h2>
                </div>
                <div class="task-list">
                    ${renderTaskList(recentTasks, t('summary.noTasks'))}
                </div>
            </div>
        </div>
        
        <!-- Footer -->
        <div class="footer">
            ${t('summary.footer')} • ${dateStr}
        </div>
    </div>
</body>
</html>`;
}

function getStatusIcon(status: string): string {
  switch (status) {
    case 'completed':
    case 'tamamlandi':
      return '✅';
    case 'in_progress':
    case 'devam_ediyor':
      return '🔄';
    default:
      return '⏳';
  }
}

function getPriorityBadge(priority: string): string {
  switch (priority) {
    case 'high':
    case 'yuksek':
      return t('summary.high');
    case 'medium':
    case 'orta':
      return t('summary.medium');
    default:
      return t('summary.low');
  }
}

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}
