import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { Gorev, GorevDurum, GorevOncelik } from '../models/gorev';
import { Logger } from '../utils/logger';
import * as path from 'path';

/**
 * Rich task detail webview panel
 */
export class TaskDetailPanel {
    private static currentPanel: TaskDetailPanel | undefined;
    private readonly panel: vscode.WebviewPanel;
    private task: Gorev;
    private disposables: vscode.Disposable[] = [];
    
    private constructor(
        panel: vscode.WebviewPanel,
        private mcpClient: MCPClient,
        task: Gorev,
        private extensionUri: vscode.Uri
    ) {
        this.panel = panel;
        this.task = task;
        
        // Set the webview's initial html content
        this.update();
        
        // Listen for when the panel is disposed
        this.panel.onDidDispose(() => this.dispose(), null, this.disposables);
        
        // Handle messages from the webview
        this.panel.webview.onDidReceiveMessage(
            message => this.handleMessage(message),
            null,
            this.disposables
        );
        
        // Update the content based on view changes
        this.panel.onDidChangeViewState(
            e => {
                if (this.panel.visible) {
                    this.update();
                }
            },
            null,
            this.disposables
        );
    }
    
    public static async createOrShow(
        mcpClient: MCPClient,
        task: Gorev,
        extensionUri: vscode.Uri
    ): Promise<void> {
        const column = vscode.window.activeTextEditor
            ? vscode.window.activeTextEditor.viewColumn
            : undefined;
        
        // If we already have a panel, show it
        if (TaskDetailPanel.currentPanel) {
            TaskDetailPanel.currentPanel.task = task;
            TaskDetailPanel.currentPanel.panel.reveal(column);
            TaskDetailPanel.currentPanel.update();
            return;
        }
        
        // Otherwise, create a new panel
        const panel = vscode.window.createWebviewPanel(
            'gorevTaskDetail',
            `Görev: ${task.baslik}`,
            column || vscode.ViewColumn.One,
            {
                enableScripts: true,
                retainContextWhenHidden: true,
                localResourceRoots: [
                    vscode.Uri.joinPath(extensionUri, 'media'),
                    vscode.Uri.joinPath(extensionUri, 'node_modules')
                ]
            }
        );
        
        // Set icon
        panel.iconPath = {
            light: vscode.Uri.joinPath(extensionUri, 'media', 'task-light.svg'),
            dark: vscode.Uri.joinPath(extensionUri, 'media', 'task-dark.svg')
        };
        
        TaskDetailPanel.currentPanel = new TaskDetailPanel(panel, mcpClient, task, extensionUri);
    }
    
    private async update() {
        try {
            // Get fresh task details from server
            const result = await this.mcpClient.callTool('gorev_detay', { id: this.task.id });
            const content = result.content[0].text;
            
            // Parse task details from markdown
            this.parseTaskDetails(content);
            
            // Update webview content
            this.panel.webview.html = this.getHtmlContent();
            this.panel.title = `Görev: ${this.task.baslik}`;
        } catch (error) {
            Logger.error('Failed to update task details:', error);
            vscode.window.showErrorMessage('Görev detayları yüklenemedi');
        }
    }
    
    private parseTaskDetails(content: string) {
        // Parse additional details from markdown content
        // This includes dependencies, tags, dates, etc.
        const lines = content.split('\n');
        
        for (const line of lines) {
            if (line.includes('Bağımlılıklar:')) {
                // Parse dependencies
                this.task.bagimliliklar = this.parseDependencies(content);
            }
        }
    }
    
    private parseDependencies(content: string): any[] {
        const dependencies: any[] = [];
        const depSection = content.split('## Bağımlılıklar')[1];
        
        if (depSection) {
            const lines = depSection.split('\n');
            for (const line of lines) {
                const match = line.match(/- (.+) \(ID: ([^)]+)\) - (.+)/);
                if (match) {
                    dependencies.push({
                        baslik: match[1],
                        id: match[2],
                        durum: match[3]
                    });
                }
            }
        }
        
        return dependencies;
    }
    
    private getHtmlContent(): string {
        const scriptUri = this.panel.webview.asWebviewUri(
            vscode.Uri.joinPath(this.extensionUri, 'media', 'taskDetail.js')
        );
        const styleUri = this.panel.webview.asWebviewUri(
            vscode.Uri.joinPath(this.extensionUri, 'media', 'taskDetail.css')
        );
        const codiconsUri = this.panel.webview.asWebviewUri(
            vscode.Uri.joinPath(this.extensionUri, 'node_modules', '@vscode/codicons', 'dist', 'codicon.css')
        );
        
        const nonce = this.getNonce();
        
        return `<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src ${this.panel.webview.cspSource} 'unsafe-inline'; script-src 'nonce-${nonce}'; font-src ${this.panel.webview.cspSource};">
    <link href="${codiconsUri}" rel="stylesheet" />
    <link href="${styleUri}" rel="stylesheet">
    <title>Görev Detayı</title>
</head>
<body>
    <div class="container">
        <!-- Header Section -->
        <div class="header">
            <div class="header-content">
                <h1 class="task-title">
                    <span class="status-icon ${this.getStatusClass()}" title="${this.task.durum}">
                        <i class="codicon ${this.getStatusIcon()}"></i>
                    </span>
                    <span contenteditable="true" id="taskTitle">${this.escapeHtml(this.task.baslik)}</span>
                </h1>
                <div class="task-meta">
                    <span class="priority priority-${this.task.oncelik.toLowerCase()}">
                        <i class="codicon codicon-arrow-up"></i> ${this.getPriorityLabel()}
                    </span>
                    ${this.task.son_tarih ? `
                        <span class="due-date ${this.getDueDateClass()}">
                            <i class="codicon codicon-calendar"></i> ${this.formatDate(this.task.son_tarih)}
                        </span>
                    ` : ''}
                    ${this.task.proje_id ? `
                        <span class="project">
                            <i class="codicon codicon-folder"></i> <span id="projectName">Proje</span>
                        </span>
                    ` : ''}
                </div>
            </div>
            <div class="header-actions">
                <button class="action-button" onclick="updateStatus()" title="Durum Güncelle">
                    <i class="codicon codicon-check"></i>
                </button>
                <button class="action-button" onclick="editTask()" title="Düzenle">
                    <i class="codicon codicon-edit"></i>
                </button>
                <button class="action-button danger" onclick="deleteTask()" title="Sil">
                    <i class="codicon codicon-trash"></i>
                </button>
            </div>
        </div>
        
        <!-- Tags Section -->
        ${this.task.etiketler && this.task.etiketler.length > 0 ? `
            <div class="tags-section">
                <h3><i class="codicon codicon-tag"></i> Etiketler</h3>
                <div class="tags">
                    ${this.task.etiketler.map((tag: string) => `
                        <span class="tag">${this.escapeHtml(tag)}</span>
                    `).join('')}
                    <button class="tag add-tag" onclick="addTag()">
                        <i class="codicon codicon-add"></i> Ekle
                    </button>
                </div>
            </div>
        ` : ''}
        
        <!-- Description Section -->
        <div class="description-section">
            <h3><i class="codicon codicon-note"></i> Açıklama</h3>
            <div class="markdown-editor">
                <div class="editor-toolbar">
                    <button onclick="toggleBold()" title="Kalın"><i class="codicon codicon-bold"></i></button>
                    <button onclick="toggleItalic()" title="İtalik"><i class="codicon codicon-italic"></i></button>
                    <button onclick="insertLink()" title="Link"><i class="codicon codicon-link"></i></button>
                    <button onclick="insertCode()" title="Kod"><i class="codicon codicon-code"></i></button>
                    <button onclick="insertList()" title="Liste"><i class="codicon codicon-list-unordered"></i></button>
                    <span class="separator"></span>
                    <button onclick="togglePreview()" title="Önizleme">
                        <i class="codicon codicon-preview"></i> Önizleme
                    </button>
                </div>
                <textarea id="descriptionEditor" class="editor-content">${this.escapeHtml(this.task.aciklama || '')}</textarea>
                <div id="descriptionPreview" class="preview-content" style="display: none;"></div>
            </div>
        </div>
        
        <!-- Dependencies Section -->
        ${this.task.bagimliliklar && this.task.bagimliliklar.length > 0 ? `
            <div class="dependencies-section">
                <h3><i class="codicon codicon-link"></i> Bağımlılıklar</h3>
                <div class="dependency-graph">
                    ${this.renderDependencyGraph()}
                </div>
                <div class="dependency-list">
                    ${this.task.bagimliliklar.map((dep: any) => `
                        <div class="dependency-item">
                            <span class="dep-status status-${dep.durum}">
                                <i class="codicon ${this.getDepStatusIcon(dep.durum)}"></i>
                            </span>
                            <span class="dep-title">${this.escapeHtml(dep.baslik)}</span>
                            <button class="link-button" onclick="openTask('${dep.id}')">
                                <i class="codicon codicon-arrow-right"></i>
                            </button>
                        </div>
                    `).join('')}
                </div>
                <button class="add-dependency" onclick="addDependency()">
                    <i class="codicon codicon-add"></i> Bağımlılık Ekle
                </button>
            </div>
        ` : ''}
        
        <!-- Activity Section -->
        <div class="activity-section">
            <h3><i class="codicon codicon-history"></i> Aktivite</h3>
            <div class="activity-timeline">
                <div class="timeline-item">
                    <span class="timeline-icon"><i class="codicon codicon-add"></i></span>
                    <div class="timeline-content">
                        <div class="timeline-title">Görev oluşturuldu</div>
                        <div class="timeline-time">${this.formatDate(this.task.olusturma_tarih)}</div>
                    </div>
                </div>
                ${this.task.guncelleme_tarih ? `
                    <div class="timeline-item">
                        <span class="timeline-icon"><i class="codicon codicon-edit"></i></span>
                        <div class="timeline-content">
                            <div class="timeline-title">Son güncelleme</div>
                            <div class="timeline-time">${this.formatDate(this.task.guncelleme_tarih)}</div>
                        </div>
                    </div>
                ` : ''}
            </div>
        </div>
    </div>
    
    <script nonce="${nonce}">
        const vscode = acquireVsCodeApi();
        const taskId = '${this.task.id}';
        
        // Handle title editing
        document.getElementById('taskTitle').addEventListener('blur', function() {
            const newTitle = this.textContent.trim();
            if (newTitle !== '${this.escapeHtml(this.task.baslik)}') {
                vscode.postMessage({
                    command: 'updateTitle',
                    title: newTitle
                });
            }
        });
        
        // Handle description editing
        let descriptionTimeout;
        document.getElementById('descriptionEditor').addEventListener('input', function() {
            clearTimeout(descriptionTimeout);
            descriptionTimeout = setTimeout(() => {
                vscode.postMessage({
                    command: 'updateDescription',
                    description: this.value
                });
            }, 1000); // Auto-save after 1 second of inactivity
        });
        
        // Command handlers
        function updateStatus() {
            vscode.postMessage({ command: 'updateStatus' });
        }
        
        function editTask() {
            vscode.postMessage({ command: 'editTask' });
        }
        
        function deleteTask() {
            if (confirm('Bu görevi silmek istediğinizden emin misiniz?')) {
                vscode.postMessage({ command: 'deleteTask' });
            }
        }
        
        function addTag() {
            const tag = prompt('Yeni etiket:');
            if (tag) {
                vscode.postMessage({ command: 'addTag', tag: tag });
            }
        }
        
        function addDependency() {
            vscode.postMessage({ command: 'addDependency' });
        }
        
        function openTask(taskId) {
            vscode.postMessage({ command: 'openTask', taskId: taskId });
        }
        
        // Markdown editor functions
        function toggleBold() {
            insertMarkdown('**', '**');
        }
        
        function toggleItalic() {
            insertMarkdown('*', '*');
        }
        
        function insertLink() {
            const url = prompt('URL:');
            if (url) {
                const text = prompt('Link metni:') || url;
                insertText('[' + text + '](' + url + ')');
            }
        }
        
        function insertCode() {
            insertMarkdown('\`', '\`');
        }
        
        function insertList() {
            insertText('\\n- ');
        }
        
        function insertMarkdown(before, after) {
            const editor = document.getElementById('descriptionEditor');
            const start = editor.selectionStart;
            const end = editor.selectionEnd;
            const text = editor.value;
            const selected = text.substring(start, end);
            
            editor.value = text.substring(0, start) + before + selected + after + text.substring(end);
            editor.focus();
            editor.setSelectionRange(start + before.length, end + before.length);
        }
        
        function insertText(text) {
            const editor = document.getElementById('descriptionEditor');
            const start = editor.selectionStart;
            const value = editor.value;
            
            editor.value = value.substring(0, start) + text + value.substring(start);
            editor.focus();
            editor.setSelectionRange(start + text.length, start + text.length);
        }
        
        let previewMode = false;
        function togglePreview() {
            previewMode = !previewMode;
            const editor = document.getElementById('descriptionEditor');
            const preview = document.getElementById('descriptionPreview');
            
            if (previewMode) {
                editor.style.display = 'none';
                preview.style.display = 'block';
                // Simple markdown to HTML conversion
                preview.innerHTML = convertMarkdownToHtml(editor.value);
            } else {
                editor.style.display = 'block';
                preview.style.display = 'none';
            }
        }
        
        function convertMarkdownToHtml(markdown) {
            return markdown
                .replace(/\\*\\*(.+?)\\*\\*/g, '<strong>$1</strong>')
                .replace(/\\*(.+?)\\*/g, '<em>$1</em>')
                .replace(/\`(.+?)\`/g, '<code>$1</code>')
                .replace(/\\[(.+?)\\]\\((.+?)\\)/g, '<a href="$2">$1</a>')
                .replace(/\\n/g, '<br>');
        }
    </script>
</body>
</html>`;
    }
    
    private renderDependencyGraph(): string {
        // Simple dependency visualization
        return `
            <svg class="dep-graph" viewBox="0 0 400 200">
                <defs>
                    <marker id="arrowhead" markerWidth="10" markerHeight="7" 
                            refX="9" refY="3.5" orient="auto">
                        <polygon points="0 0, 10 3.5, 0 7" fill="#666" />
                    </marker>
                </defs>
                
                <!-- Current task -->
                <rect x="150" y="80" width="100" height="40" rx="5" 
                      class="current-task-node" />
                <text x="200" y="105" text-anchor="middle" class="node-text">
                    ${this.escapeHtml(this.task.baslik.substring(0, 10))}...
                </text>
                
                ${this.task.bagimliliklar?.map((dep: any, index: number) => `
                    <!-- Dependency ${index + 1} -->
                    <rect x="${50 + (index * 120)}" y="20" width="100" height="40" rx="5" 
                          class="dep-node status-${dep.durum}" />
                    <text x="${100 + (index * 120)}" y="45" text-anchor="middle" class="node-text">
                        ${this.escapeHtml(dep.baslik.substring(0, 10))}...
                    </text>
                    <line x1="${100 + (index * 120)}" y1="60" x2="200" y2="80" 
                          stroke="#666" stroke-width="2" marker-end="url(#arrowhead)" />
                `).join('') || ''}
            </svg>
        `;
    }
    
    private async handleMessage(message: any) {
        try {
            switch (message.command) {
                case 'updateTitle':
                    await this.updateTaskField('baslik', message.title);
                    break;
                    
                case 'updateDescription':
                    await this.updateTaskField('aciklama', message.description);
                    break;
                    
                case 'updateStatus':
                    await this.showStatusPicker();
                    break;
                    
                case 'editTask':
                    await this.showEditDialog();
                    break;
                    
                case 'deleteTask':
                    await this.deleteTask();
                    break;
                    
                case 'addTag':
                    await this.addTag(message.tag);
                    break;
                    
                case 'addDependency':
                    await this.showDependencyPicker();
                    break;
                    
                case 'openTask':
                    await this.openTask(message.taskId);
                    break;
            }
        } catch (error) {
            Logger.error('Error handling webview message:', error);
            vscode.window.showErrorMessage(`İşlem başarısız: ${error}`);
        }
    }
    
    private async updateTaskField(field: string, value: any) {
        try {
            const params: any = { id: this.task.id };
            params[field] = value;
            
            await this.mcpClient.callTool('gorev_duzenle', params);
            (this.task as any)[field] = value;
            
            vscode.window.showInformationMessage('Görev güncellendi');
        } catch (error) {
            throw error;
        }
    }
    
    private async showStatusPicker() {
        const items = [
            { label: 'Beklemede', value: GorevDurum.Beklemede },
            { label: 'Devam Ediyor', value: GorevDurum.DevamEdiyor },
            { label: 'Tamamlandı', value: GorevDurum.Tamamlandi }
        ];
        
        const selected = await vscode.window.showQuickPick(items, {
            placeHolder: 'Yeni durum seçin'
        });
        
        if (selected) {
            await this.mcpClient.callTool('gorev_guncelle', {
                id: this.task.id,
                durum: selected.value
            });
            this.task.durum = selected.value;
            this.update();
        }
    }
    
    private async showEditDialog() {
        // Open multi-step input dialog for editing all fields
        vscode.commands.executeCommand('gorev.detailedEdit', this.task);
    }
    
    private async deleteTask() {
        await this.mcpClient.callTool('gorev_sil', {
            id: this.task.id,
            onay: true
        });
        
        this.panel.dispose();
        vscode.window.showInformationMessage('Görev silindi');
    }
    
    private async addTag(tag: string) {
        const currentTags = this.task.etiketler || [];
        currentTags.push(tag);
        
        await this.updateTaskField('etiketler', currentTags.join(','));
        this.task.etiketler = currentTags;
        this.update();
    }
    
    private async showDependencyPicker() {
        // Show task picker for adding dependency
        vscode.commands.executeCommand('gorev.addDependency', this.task.id);
    }
    
    private async openTask(taskId: string) {
        // Get task details and open in new panel
        const result = await this.mcpClient.callTool('gorev_detay', { id: taskId });
        // Parse task from result
        const task: Gorev = {
            id: taskId,
            baslik: 'Task', // Will be parsed from result
            durum: GorevDurum.Beklemede,
            oncelik: GorevOncelik.Orta,
            // ... parse other fields
        } as Gorev;
        
        await TaskDetailPanel.createOrShow(this.mcpClient, task, this.extensionUri);
    }
    
    private getStatusClass(): string {
        switch (this.task.durum) {
            case GorevDurum.Tamamlandi: return 'status-completed';
            case GorevDurum.DevamEdiyor: return 'status-in-progress';
            default: return 'status-pending';
        }
    }
    
    private getStatusIcon(): string {
        switch (this.task.durum) {
            case GorevDurum.Tamamlandi: return 'codicon-pass-filled';
            case GorevDurum.DevamEdiyor: return 'codicon-debug-start';
            default: return 'codicon-circle-outline';
        }
    }
    
    private getDepStatusIcon(durum: string): string {
        if (durum.includes('tamamland')) return 'codicon-pass-filled';
        if (durum.includes('devam')) return 'codicon-debug-start';
        return 'codicon-circle-outline';
    }
    
    private getPriorityLabel(): string {
        switch (this.task.oncelik) {
            case GorevOncelik.Yuksek: return 'Yüksek';
            case GorevOncelik.Orta: return 'Orta';
            case GorevOncelik.Dusuk: return 'Düşük';
            default: return 'Orta';
        }
    }
    
    private getDueDateClass(): string {
        if (!this.task.son_tarih) return '';
        
        const due = new Date(this.task.son_tarih);
        const today = new Date();
        today.setHours(0, 0, 0, 0);
        
        if (due < today && this.task.durum !== GorevDurum.Tamamlandi) {
            return 'overdue';
        } else if (due.toDateString() === today.toDateString()) {
            return 'due-today';
        }
        
        return '';
    }
    
    private formatDate(dateStr?: string): string {
        if (!dateStr) return '';
        
        const date = new Date(dateStr);
        const options: Intl.DateTimeFormatOptions = {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        };
        
        return date.toLocaleDateString('tr-TR', options);
    }
    
    private escapeHtml(text: string): string {
        const map: Record<string, string> = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#039;'
        };
        
        return text.replace(/[&<>"']/g, m => map[m]);
    }
    
    private getNonce(): string {
        let text = '';
        const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
        for (let i = 0; i < 32; i++) {
            text += possible.charAt(Math.floor(Math.random() * possible.length));
        }
        return text;
    }
    
    public dispose() {
        TaskDetailPanel.currentPanel = undefined;
        
        // Clean up resources
        this.panel.dispose();
        
        while (this.disposables.length) {
            const x = this.disposables.pop();
            if (x) {
                x.dispose();
            }
        }
    }
}