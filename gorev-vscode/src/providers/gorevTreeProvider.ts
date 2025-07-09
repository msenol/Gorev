import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { Gorev, GorevDurum, GorevOncelik } from '../models/gorev';
import { ICONS, COLORS, CONTEXT_VALUES } from '../utils/constants';
import { Logger } from '../utils/logger';

export class GorevTreeProvider implements vscode.TreeDataProvider<GorevTreeItem> {
  private _onDidChangeTreeData = new vscode.EventEmitter<GorevTreeItem | undefined | null | void>();
  readonly onDidChangeTreeData = this._onDidChangeTreeData.event;
  
  private tasks: Gorev[] = [];

  constructor(private mcpClient: MCPClient) {}

  getTreeItem(element: GorevTreeItem): vscode.TreeItem {
    return element;
  }

  async getChildren(element?: GorevTreeItem): Promise<GorevTreeItem[]> {
    if (!this.mcpClient.isConnected()) {
      return [];
    }

    if (!element) {
      // Root level - return all tasks
      try {
        await this.loadTasks();
        return this.tasks.map((task) => new GorevTreeItem(task));
      } catch (error) {
        Logger.error('Failed to load tasks:', error);
        return [];
      }
    }

    return [];
  }

  async refresh(): Promise<void> {
    await this.loadTasks();
    this._onDidChangeTreeData.fire();
  }

  private async loadTasks(): Promise<void> {
    try {
      this.tasks = [];
      let offset = 0;
      const pageSize = 100;
      let hasMoreTasks = true;
      
      // Fetch all tasks with pagination
      while (hasMoreTasks) {
        console.log('[GorevTreeProvider] Fetching tasks with offset:', offset);
        
        const result = await this.mcpClient.callTool('gorev_listele', {
          tum_projeler: true,
          limit: pageSize,
          offset: offset
        });
        
        if (!result || !result.content || !result.content[0]) {
          break;
        }
        
        const responseText = result.content[0].text;
        
        // Check for pagination info
        const paginationMatch = responseText.match(/Görevler \((\d+)-(\d+) \/ (\d+)\)/);
        if (paginationMatch) {
          const [_, start, end, total] = paginationMatch;
          console.log(`[GorevTreeProvider] Pagination: ${start}-${end} / ${total}`);
          
          if (parseInt(end) >= parseInt(total)) {
            hasMoreTasks = false;
          }
        } else {
          hasMoreTasks = false;
        }
        
        // Parse tasks from this page
        const pageTasks = this.parseTasksFromContent(responseText);
        this.tasks.push(...pageTasks);
        
        // Update offset
        offset += pageSize;
        
        // Safety check
        if (offset > 1000) break;
      }
      
      console.log('[GorevTreeProvider] Total tasks fetched:', this.tasks.length);
    } catch (error) {
      Logger.error('Failed to load tasks:', error);
      throw error;
    }
  }

  private parseTasksFromContent(content: string): Gorev[] {
    const tasks: Gorev[] = [];
    
    // Check for empty task list
    if (content.includes('Henüz görev bulunmuyor')) {
      return tasks;
    }
    
    // Split content into lines for processing
    const lines = content.split('\n');
    let currentTask: Partial<Gorev> | null = null;
    let descriptionLines: string[] = [];
    
    for (let i = 0; i < lines.length; i++) {
      const line = lines[i].trim();
      
      // Skip empty lines and headers
      if (!line || line.startsWith('##')) {
        continue;
      }
      
      // Check if this is a task line (starts with "- [")
      const taskMatch = line.match(/^- \[([^\]]+)\] (.+) \(([^)]+) öncelik\)$/);
      if (taskMatch) {
        // Save previous task if exists
        if (currentTask && currentTask.id) {
          if (descriptionLines.length > 0) {
            currentTask.aciklama = descriptionLines.join('\n').trim();
          }
          tasks.push(currentTask as Gorev);
        }
        
        // Start new task
        const [, durum, baslik, oncelik] = taskMatch;
        currentTask = {
          baslik,
          durum: this.mapDurum(durum),
          oncelik: this.mapOncelik(oncelik),
          etiketler: [],
        };
        descriptionLines = [];
        continue;
      }
      
      // Check for ID line
      if (line.startsWith('ID:') && currentTask) {
        currentTask.id = line.substring(3).trim();
        continue;
      }
      
      // Check for project line
      if (line.startsWith('Proje:') && currentTask) {
        // Just skip the project name, as we already have proje_id in the task
        continue;
      }
      
      // Check for due date in description
      if (line.includes('Son tarih:') && currentTask) {
        const dateMatch = line.match(/Son tarih: (\d{4}-\d{2}-\d{2})/);
        if (dateMatch) {
          currentTask.son_tarih = dateMatch[1];
        }
        continue;
      }
      
      // Check for tags in description
      if (line.includes('Etiketler:') && currentTask) {
        const tagsMatch = line.match(/Etiketler: (.+)/);
        if (tagsMatch) {
          currentTask.etiketler = tagsMatch[1].split(',').map(tag => tag.trim());
        }
        continue;
      }
      
      // If we're in a task and the line is indented or part of description
      if (currentTask && line && !line.startsWith('-')) {
        descriptionLines.push(line);
      }
    }
    
    // Don't forget the last task
    if (currentTask && currentTask.id) {
      if (descriptionLines.length > 0) {
        currentTask.aciklama = descriptionLines.join('\n').trim();
      }
      tasks.push(currentTask as Gorev);
    }
    
    return tasks;
  }
  
  private mapDurum(durum: string): GorevDurum {
    const durumMap: { [key: string]: GorevDurum } = {
      'beklemede': GorevDurum.Beklemede,
      'devam_ediyor': GorevDurum.DevamEdiyor,
      'tamamlandi': GorevDurum.Tamamlandi,
    };
    
    return durumMap[durum.toLowerCase()] || GorevDurum.Beklemede;
  }
  
  private mapOncelik(oncelik: string): GorevOncelik {
    const oncelikMap: { [key: string]: GorevOncelik } = {
      'dusuk': GorevOncelik.Dusuk,
      'orta': GorevOncelik.Orta,
      'yuksek': GorevOncelik.Yuksek,
    };
    
    return oncelikMap[oncelik.toLowerCase()] || GorevOncelik.Orta;
  }
}

export class GorevTreeItem extends vscode.TreeItem {
  constructor(
    public readonly task: Gorev,
    public readonly collapsibleState: vscode.TreeItemCollapsibleState = vscode.TreeItemCollapsibleState.None
  ) {
    super(task.baslik, collapsibleState);
    
    this.tooltip = this.getTooltip();
    this.description = this.getDescription();
    this.iconPath = this.getIcon();
    this.contextValue = CONTEXT_VALUES.TASK;
  }

  private getTooltip(): string {
    const lines = [
      this.task.baslik,
      `Status: ${this.task.durum}`,
      `Priority: ${this.task.oncelik}`,
    ];
    
    if (this.task.son_tarih) {
      lines.push(`Due: ${this.task.son_tarih}`);
    }
    
    if (this.task.aciklama) {
      lines.push('', this.task.aciklama);
    }
    
    return lines.join('\n');
  }

  private getDescription(): string {
    const parts = [];
    
    if (this.task.son_tarih) {
      const dueDate = new Date(this.task.son_tarih);
      const today = new Date();
      today.setHours(0, 0, 0, 0);
      
      if (dueDate < today) {
        parts.push('⚠️ Overdue');
      } else if (dueDate.getTime() - today.getTime() <= 7 * 24 * 60 * 60 * 1000) {
        parts.push('⏰ Due soon');
      }
    }
    
    if (this.task.etiketler && this.task.etiketler.length > 0) {
      parts.push(`[${this.task.etiketler.join(', ')}]`);
    }
    
    return parts.join(' ');
  }

  private getIcon(): vscode.ThemeIcon {
    let iconId: string;
    
    switch (this.task.durum) {
      case GorevDurum.Tamamlandi:
        iconId = ICONS.TASK_COMPLETED;
        break;
      case GorevDurum.DevamEdiyor:
        iconId = ICONS.TASK_IN_PROGRESS;
        break;
      default:
        iconId = ICONS.TASK_PENDING;
    }
    
    const icon = new vscode.ThemeIcon(iconId);
    
    // Apply color based on priority
    let color: string | undefined;
    switch (this.task.oncelik) {
      case GorevOncelik.Yuksek:
        color = COLORS.HIGH_PRIORITY;
        break;
      case GorevOncelik.Orta:
        color = COLORS.MEDIUM_PRIORITY;
        break;
      case GorevOncelik.Dusuk:
        color = COLORS.LOW_PRIORITY;
        break;
    }
    
    if (color) {
      return new vscode.ThemeIcon(iconId, new vscode.ThemeColor(color));
    }
    
    return icon;
  }
}