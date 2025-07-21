import * as vscode from 'vscode';
import { MCPClient } from '../mcp/client';
import { Gorev, GorevDurum, GorevOncelik, GorevHiyerarsi } from '../models/gorev';
import { Logger } from '../utils/logger';
import { MarkdownParser } from '../utils/markdownParser';
import * as path from 'path';

/**
 * Rich task detail webview panel
 */
export class TaskDetailPanel {
    private static currentPanel: TaskDetailPanel | undefined;
    private readonly panel: vscode.WebviewPanel;
    private task: Gorev;
    private hierarchyInfo?: GorevHiyerarsi;
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
        
        // If we already have a panel, dispose it and create a new one to force reload
        if (TaskDetailPanel.currentPanel) {
            TaskDetailPanel.currentPanel.dispose();
        }
        
        // Otherwise, create a new panel
        const panel = vscode.window.createWebviewPanel(
            'gorevTaskDetail',
            vscode.l10n.t('taskDetail.taskTitle', task.baslik),
            column || vscode.ViewColumn.One,
            {
                enableScripts: true,
                retainContextWhenHidden: false, // Force reload on hide/show
                localResourceRoots: [
                    vscode.Uri.joinPath(extensionUri, 'media'),
                    vscode.Uri.joinPath(extensionUri, 'node_modules', '@vscode', 'codicons', 'dist')
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
            
            // Get hierarchy information if available
            try {
                const hierarchyResult = await this.mcpClient.callTool('gorev_hiyerarsi_goster', { gorev_id: this.task.id });
                if (hierarchyResult && hierarchyResult.content && hierarchyResult.content[0]) {
                    this.parseHierarchyInfo(hierarchyResult.content[0].text);
                }
            } catch (err) {
                Logger.debug('Hierarchy info not available:', err);
            }
            
            // Update webview content
            Logger.info('[TaskDetailPanel] Setting new HTML content with breadcrumb navigation');
            this.panel.webview.html = this.getHtmlContent();
            this.panel.title = vscode.l10n.t('taskDetail.taskTitle', this.task.baslik);
            Logger.info('[TaskDetailPanel] HTML content updated successfully');
        } catch (error) {
            Logger.error('Failed to update task details:', error);
            vscode.window.showErrorMessage(vscode.l10n.t('taskDetail.loadFailed'));
        }
    }
    
    private parseTaskDetails(content: string) {
        // Parse additional details from markdown content
        // This includes dependencies, tags, dates, etc.
        const parsedTask = MarkdownParser.parseGorevDetay(content);
        
        // Update task with parsed details
        if (parsedTask.etiketler) {
            this.task.etiketler = parsedTask.etiketler;
        }
        if (parsedTask.bagimliliklar) {
            this.task.bagimliliklar = parsedTask.bagimliliklar;
        }
        if (parsedTask.son_tarih) {
            this.task.son_tarih = parsedTask.son_tarih;
        }
        
        // Add debug logging
        Logger.debug('Parsed task details:', {
            id: this.task.id,
            etiketler: this.task.etiketler,
            bagimliliklar: this.task.bagimliliklar
        });
    }
    
    private parseHierarchyInfo(content: string) {
        // Parse hierarchy statistics from the response
        const lines = content.split('\n');
        const stats: any = {};
        
        // Add debug logging
        Logger.debug('Parsing hierarchy info from content:', content);
        
        for (const line of lines) {
            if (line.includes('Toplam Alt Görev:')) {
                const match = line.match(/\*?\*?Toplam Alt Görev:\*?\*?\s*(\d+)/);
                if (match) stats.toplamAltGorev = parseInt(match[1]);
            }
            if (line.includes('Tamamlanan:')) {
                const match = line.match(/\*?\*?Tamamlanan:\*?\*?\s*(\d+)/);
                if (match) stats.tamamlananAlt = parseInt(match[1]);
            }
            if (line.includes('Devam Eden:')) {
                const match = line.match(/\*?\*?Devam Eden:\*?\*?\s*(\d+)/);
                if (match) stats.devamEdenAlt = parseInt(match[1]);
            }
            if (line.includes('Beklemede:')) {
                const match = line.match(/\*?\*?Beklemede:\*?\*?\s*(\d+)/);
                if (match) stats.beklemedeAlt = parseInt(match[1]);
            }
            // More flexible parsing for İlerleme (Progress)
            if (line.includes('İlerleme:') || line.includes('Progress:')) {
                // Try multiple patterns for better compatibility
                const patterns = [
                    /\*?\*?İlerleme:\*?\*?\s*([\d.]+)%/,  // Handles **İlerleme:** format
                    /İlerleme:\s*([\d.]+)%/,
                    /İlerleme:\s*%([\d.]+)/,
                    /Progress:\s*([\d.]+)%/,
                    /İlerleme:\s*([\d.]+)/
                ];
                
                for (const pattern of patterns) {
                    const match = line.match(pattern);
                    if (match) {
                        stats.ilerlemeYuzdesi = parseFloat(match[1]);
                        Logger.debug('Parsed progress percentage:', stats.ilerlemeYuzdesi);
                        break;
                    }
                }
            }
        }
        
        // Calculate progress if not provided but we have the data
        if (stats.ilerlemeYuzdesi === undefined && stats.toplamAltGorev > 0) {
            stats.ilerlemeYuzdesi = Math.round((stats.tamamlananAlt / stats.toplamAltGorev) * 100);
            Logger.debug('Calculated progress percentage:', stats.ilerlemeYuzdesi);
        }
        
        // Ensure ilerleme_yuzdesi is always a valid number
        const progressPercentage = stats.ilerlemeYuzdesi || 0;
        const validPercentage = isNaN(progressPercentage) ? 0 : Math.min(100, Math.max(0, progressPercentage));
        
        this.hierarchyInfo = {
            gorev: this.task,
            ust_gorevler: [],
            toplam_alt_gorev: stats.toplamAltGorev || 0,
            tamamlanan_alt: stats.tamamlananAlt || 0,
            devam_eden_alt: stats.devamEdenAlt || 0,
            beklemede_alt: stats.beklemedeAlt || 0,
            ilerleme_yuzdesi: validPercentage
        };
        
        Logger.debug('Final hierarchy info:', this.hierarchyInfo);
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
        // Add timestamp for cache-busting
        const timestamp = new Date().getTime();
        const styleUri = this.panel.webview.asWebviewUri(
            vscode.Uri.joinPath(this.extensionUri, 'media', 'taskDetail.css')
        ) + `?v=${timestamp}`;
        // VS Code provides codicons through its own CSS
        const vscodeIconsUri = 'https://microsoft.github.io/vscode-codicons/dist/codicon.css';
        
        const nonce = this.getNonce();
        
        return `<!DOCTYPE html>
<html lang="tr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src ${this.panel.webview.cspSource} 'unsafe-inline' https://microsoft.github.io; script-src 'nonce-${nonce}'; font-src ${this.panel.webview.cspSource} https://microsoft.github.io;">
    <link href="${vscodeIconsUri}" rel="stylesheet" />
    <link href="${styleUri}" rel="stylesheet">
    <title>${vscode.l10n.t('taskDetail.title')}</title>
    <script nonce="${nonce}">
        console.log('TaskDetail CSS loaded:', '${styleUri}');
        console.log('Page loaded at:', new Date().toISOString());
        console.log('HTML Version: 2.0 - Two Column Layout with Breadcrumb');
        // Log when DOM is ready
        document.addEventListener('DOMContentLoaded', function() {
            console.log('DOM Content Loaded');
            console.log('Found breadcrumb:', document.querySelector('.breadcrumb-navigation'));
            console.log('Found two-column layout:', document.querySelector('.content-layout'));
        });
    </script>
</head>
<body>
    <div class="main-container">
        <!-- Breadcrumb Navigation -->
        ${this.renderBreadcrumb()}
        
        <!-- Two Column Layout -->
        <div class="content-layout">
            <!-- Main Content -->
            <div class="main-content">
                <!-- Header Section -->
                <div class="header card">
                    <div class="header-content">
                        <div class="status-badge ${this.getStatusClass()}">
                            <i class="codicon ${this.getStatusIcon()}"></i>
                            <span>${this.getStatusLabel()}</span>
                        </div>
                        <h1 class="task-title">
                            <span contenteditable="true" id="taskTitle">${this.escapeHtml(this.task.baslik)}</span>
                        </h1>
                        <div class="task-meta">
                            <span class="priority-badge priority-${this.task.oncelik.toLowerCase()}">
                                <i class="codicon codicon-arrow-up"></i> ${this.getPriorityLabel()}
                            </span>
                            ${this.task.son_tarih ? `
                                <span class="due-date-badge ${this.getDueDateClass()}">
                                    <i class="codicon codicon-calendar"></i> ${this.formatDate(this.task.son_tarih)}
                                </span>
                            ` : ''}
                            ${this.task.proje_id ? `
                                <span class="project-badge">
                                    <i class="codicon codicon-folder"></i> <span id="projectName">${vscode.l10n.t('taskDetail.project')}</span>
                                </span>
                            ` : ''}
                        </div>
                    </div>
                </div>
                
                <!-- Description Section -->
                <div class="description-section card">
                    <div class="section-header">
                        <h3><i class="codicon codicon-note"></i> ${vscode.l10n.t('taskDetail.description')}</h3>
                        <div class="editor-mode-toggle">
                            <button class="mode-btn active" data-mode="edit">
                                <i class="codicon codicon-edit"></i> ${vscode.l10n.t('taskDetail.edit')}
                            </button>
                            <button class="mode-btn" data-mode="preview">
                                <i class="codicon codicon-eye"></i> ${vscode.l10n.t('taskDetail.preview')}
                            </button>
                            <button class="mode-btn" data-mode="split">
                                <i class="codicon codicon-split-horizontal"></i> ${vscode.l10n.t('taskDetail.split')}
                            </button>
                        </div>
                    </div>
                    <div class="markdown-editor enhanced">
                        <div class="editor-toolbar">
                            <div class="toolbar-group">
                                <button id="boldBtn" title="${vscode.l10n.t('taskDetail.bold')}">
                                    <i class="codicon codicon-bold"></i>
                                </button>
                                <button id="italicBtn" title="${vscode.l10n.t('taskDetail.italic')}">
                                    <i class="codicon codicon-italic"></i>
                                </button>
                                <button id="strikeBtn" title="${vscode.l10n.t('taskDetail.strikethrough')}">
                                    <i class="codicon codicon-text-strikethrough"></i>
                                    <span style="font-size: 11px; margin-left: 2px;">S</span>
                                </button>
                            </div>
                            <div class="toolbar-separator"></div>
                            <div class="toolbar-group">
                                <button id="h1Btn" title="${vscode.l10n.t('taskDetail.heading1')}">H1</button>
                                <button id="h2Btn" title="${vscode.l10n.t('taskDetail.heading2')}">H2</button>
                                <button id="h3Btn" title="${vscode.l10n.t('taskDetail.heading3')}">H3</button>
                            </div>
                            <div class="toolbar-separator"></div>
                            <div class="toolbar-group">
                                <button id="linkBtn" title="${vscode.l10n.t('taskDetail.link')}">
                                    <i class="codicon codicon-link"></i>
                                </button>
                                <button id="imageBtn" title="${vscode.l10n.t('taskDetail.image')}">
                                    <i class="codicon codicon-file-media"></i>
                                </button>
                                <button id="codeBtn" title="${vscode.l10n.t('taskDetail.code')}">
                                    <i class="codicon codicon-code"></i>
                                </button>
                                <button id="codeBlockBtn" title="${vscode.l10n.t('taskDetail.codeBlock')}">
                                    <i class="codicon codicon-symbol-namespace"></i>
                                </button>
                            </div>
                            <div class="toolbar-separator"></div>
                            <div class="toolbar-group">
                                <button id="listBtn" title="${vscode.l10n.t('taskDetail.bulletList')}">
                                    <i class="codicon codicon-list-unordered"></i>
                                </button>
                                <button id="orderedListBtn" title="${vscode.l10n.t('taskDetail.numberedList')}">
                                    <i class="codicon codicon-list-ordered"></i>
                                </button>
                                <button id="checklistBtn" title="${vscode.l10n.t('taskDetail.taskList')}">
                                    <i class="codicon codicon-checklist"></i>
                                </button>
                                <button id="tableBtn" title="${vscode.l10n.t('taskDetail.table')}">
                                    <i class="codicon codicon-table"></i>
                                </button>
                            </div>
                            <div class="toolbar-separator"></div>
                            <div class="toolbar-group">
                                <button id="undoBtn" title="${vscode.l10n.t('taskDetail.undo')}">
                                    <i class="codicon codicon-discard"></i>
                                </button>
                                <button id="redoBtn" title="${vscode.l10n.t('taskDetail.redo')}">
                                    <i class="codicon codicon-redo"></i>
                                </button>
                            </div>
                            <div class="toolbar-spacer"></div>
                            <div class="toolbar-status">
                                <span id="saveStatus" class="save-status">
                                    <i class="codicon codicon-check"></i> ${vscode.l10n.t('taskDetail.saved')}
                                </span>
                            </div>
                        </div>
                        <div class="editor-container" id="editorContainer">
                            <div class="editor-pane">
                                <textarea id="descriptionEditor" class="editor-content" placeholder="${vscode.l10n.t('taskDetail.writeSomething')}">${this.escapeHtml(this.task.aciklama || '')}</textarea>
                            </div>
                            <div class="preview-pane" id="previewPane" style="display: none;">
                                <div id="descriptionPreview" class="preview-content"></div>
                            </div>
                        </div>
                    </div>
                </div>
                
                <!-- Tags Section -->
                <div class="tags-section card">
                    <h3><i class="codicon codicon-tag"></i> ${vscode.l10n.t('taskDetail.tags')}</h3>
                    <div class="tags-container">
                        ${this.task.etiketler && this.task.etiketler.length > 0 ? 
                            this.task.etiketler.map((tag: string) => `
                                <span class="tag">
                                    <span class="tag-text">${this.escapeHtml(tag)}</span>
                                    <button class="tag-remove" data-tag="${this.escapeHtml(tag)}">
                                        <i class="codicon codicon-close"></i>
                                    </button>
                                </span>
                            `).join('') : 
                            '<span class="empty-state">' + vscode.l10n.t('taskDetail.noTags') + '</span>'
                        }
                        <button class="tag-add" id="addTagBtn">
                            <i class="codicon codicon-add"></i> ${vscode.l10n.t('taskDetail.tags')}
                        </button>
                    </div>
                </div>
            </div>
            
            <!-- Sidebar -->
            <div class="sidebar">
                <!-- Quick Actions -->
                <div class="quick-actions card">
                    <h3>${vscode.l10n.t('taskDetail.actions')}</h3>
                    <div class="actions-grid">
                        <button class="quick-action-btn" id="updateStatusBtn" title="${vscode.l10n.t('taskDetail.updateStatus')}">
                            <i class="codicon codicon-check"></i>
                            <span>${vscode.l10n.t('taskDetail.updateStatus')}</span>
                        </button>
                        <button class="quick-action-btn" id="editTaskBtn" title="${vscode.l10n.t('taskDetail.edit')}">
                            <i class="codicon codicon-edit"></i>
                            <span>${vscode.l10n.t('taskDetail.edit')}</span>
                        </button>
                        <button class="quick-action-btn" id="duplicateBtn" title="${vscode.l10n.t('taskDetail.duplicateTask')}">
                            <i class="codicon codicon-files"></i>
                            <span>${vscode.l10n.t('taskDetail.duplicateTask')}</span>
                        </button>
                        <button class="quick-action-btn danger" id="deleteTaskBtn" title="${vscode.l10n.t('taskDetail.deleteTask')}">
                            <i class="codicon codicon-trash"></i>
                            <span>${vscode.l10n.t('taskDetail.deleteTask')}</span>
                        </button>
                    </div>
                </div>
                
                <!-- Hierarchy Section -->
                ${this.renderEnhancedHierarchySection()}
                
                <!-- Dependencies Section -->
                ${this.renderDependenciesSection()}
                
                <!-- Activity Section -->
                <div class="activity-section card">
                    <h3><i class="codicon codicon-history"></i> ${vscode.l10n.t('taskDetail.activity')}</h3>
                    <div class="activity-timeline compact">
                        ${this.renderActivityTimeline()}
                    </div>
                </div>
            </div>
        </div>
    </div>
    
    <script nonce="${nonce}">
        const vscode = acquireVsCodeApi();
        const taskId = '${this.task.id}';
        
        // Debug: Check if styles are loaded
        window.addEventListener('load', () => {
            const styles = document.styleSheets;
            console.log('Loaded stylesheets:', styles.length);
            for (let i = 0; i < styles.length; i++) {
                console.log('Stylesheet', i, ':', styles[i].href);
            }
        });
        
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
        
        // Handle create subtask button
        const createSubtaskBtn = document.getElementById('createSubtaskBtn');
        if (createSubtaskBtn) {
            createSubtaskBtn.addEventListener('click', function() {
                vscode.postMessage({ command: 'createSubtask' });
            });
        }
        
        // Add event listeners for buttons
        document.getElementById('updateStatusBtn').addEventListener('click', function() {
            vscode.postMessage({ command: 'updateStatus' });
        });
        
        document.getElementById('editTaskBtn').addEventListener('click', function() {
            vscode.postMessage({ command: 'editTask' });
        });
        
        document.getElementById('deleteTaskBtn').addEventListener('click', function() {
            vscode.postMessage({ command: 'deleteTask' });
        });
        
        document.getElementById('addTagBtn').addEventListener('click', function() {
            const tag = prompt('Yeni etiket:');
            if (tag) {
                vscode.postMessage({ command: 'addTag', tag: tag });
            }
        });
        
        // Editor mode toggle
        document.querySelectorAll('.mode-btn').forEach(btn => {
            btn.addEventListener('click', function() {
                document.querySelectorAll('.mode-btn').forEach(b => b.classList.remove('active'));
                this.classList.add('active');
                handleEditorMode(this.dataset.mode);
            });
        });
        
        // Markdown editor buttons
        const boldBtn = document.getElementById('boldBtn');
        const italicBtn = document.getElementById('italicBtn');
        const strikeBtn = document.getElementById('strikeBtn');
        
        if (boldBtn) {
            boldBtn.addEventListener('click', function() { 
                console.log('Bold button clicked');
                toggleBold(); 
            });
        }
        if (italicBtn) {
            italicBtn.addEventListener('click', function() { 
                console.log('Italic button clicked');
                toggleItalic(); 
            });
        }
        if (strikeBtn) {
            strikeBtn.addEventListener('click', function() { 
                console.log('Strike button clicked');
                toggleStrike(); 
            });
        }
        // Safely add event listeners for all editor buttons
        const editorButtons = [
            { id: 'h1Btn', handler: () => insertHeading(1) },
            { id: 'h2Btn', handler: () => insertHeading(2) },
            { id: 'h3Btn', handler: () => insertHeading(3) },
            { id: 'linkBtn', handler: insertLink },
            { id: 'imageBtn', handler: insertImage },
            { id: 'codeBtn', handler: insertCode },
            { id: 'codeBlockBtn', handler: insertCodeBlock },
            { id: 'listBtn', handler: insertList },
            { id: 'orderedListBtn', handler: insertOrderedList },
            { id: 'checklistBtn', handler: insertChecklist },
            { id: 'tableBtn', handler: insertTable },
            { id: 'undoBtn', handler: performUndo },
            { id: 'redoBtn', handler: performRedo }
        ];
        
        editorButtons.forEach(({ id, handler }) => {
            const btn = document.getElementById(id);
            if (btn) {
                btn.addEventListener('click', handler);
            } else {
                console.warn('Button not found:', id);
            }
        });
        
        // Add dependency button
        const addDepBtn = document.getElementById('addDependencyBtn');
        if (addDepBtn) {
            addDepBtn.addEventListener('click', function() {
                vscode.postMessage({ command: 'addDependency' });
            });
        }
        
        // Link buttons for dependencies
        document.querySelectorAll('.link-button').forEach(btn => {
            btn.addEventListener('click', function() {
                const taskId = this.getAttribute('data-task-id');
                vscode.postMessage({ command: 'openTask', taskId: taskId });
            });
        });
        
        // Markdown editor functions
        function toggleBold() {
            insertMarkdown('**', '**');
        }
        
        function toggleItalic() {
            insertMarkdown('*', '*');
        }
        
        function toggleStrike() {
            insertMarkdown('~~', '~~');
        }
        
        function insertHeading(level) {
            const prefix = '#'.repeat(level) + ' ';
            insertAtLineStart(prefix);
        }
        
        function insertLink() {
            // Use vscode message passing instead of prompt
            const selectedText = getSelectedText();
            if (selectedText) {
                vscode.postMessage({ command: 'insertLink', selectedText: selectedText });
            } else {
                vscode.postMessage({ command: 'insertLink' });
            }
        }
        
        function insertImage() {
            vscode.postMessage({ command: 'insertImage' });
        }
        
        function insertCode() {
            insertMarkdown('\`', '\`');
        }
        
        function insertCodeBlock() {
            vscode.postMessage({ command: 'insertCodeBlock' });
        }
        
        function getSelectedText() {
            const editor = document.getElementById('descriptionEditor');
            const start = editor.selectionStart;
            const end = editor.selectionEnd;
            if (start !== end) {
                return editor.value.substring(start, end);
            }
            return '';
        }
        
        function insertList() {
            insertAtLineStart('- ');
        }
        
        function insertOrderedList() {
            insertAtLineStart('1. ');
        }
        
        function insertChecklist() {
            insertAtLineStart('- [ ] ');
        }
        
        function insertTable() {
            vscode.postMessage({ command: 'insertTable' });
        }
        
        function insertMarkdown(before, after) {
            const editor = document.getElementById('descriptionEditor');
            const start = editor.selectionStart;
            const end = editor.selectionEnd;
            const text = editor.value;
            const selected = text.substring(start, end);
            
            // Save current state before change
            saveUndoState(text);
            
            editor.value = text.substring(0, start) + before + selected + after + text.substring(end);
            editor.focus();
            editor.setSelectionRange(start + before.length, end + before.length);
            
            // Save new state and trigger input
            lastValue = editor.value;
            editor.dispatchEvent(new Event('input'));
        }
        
        function insertText(text) {
            const editor = document.getElementById('descriptionEditor');
            const start = editor.selectionStart;
            const value = editor.value;
            
            // Save current state before change
            saveUndoState(value);
            
            editor.value = value.substring(0, start) + text + value.substring(start);
            editor.focus();
            editor.setSelectionRange(start + text.length, start + text.length);
            
            // Save new state
            lastValue = editor.value;
            
            // Trigger input event for auto-save
            editor.dispatchEvent(new Event('input'));
        }
        
        function insertAtLineStart(text) {
            const editor = document.getElementById('descriptionEditor');
            const start = editor.selectionStart;
            const value = editor.value;
            
            // Save current state before change
            saveUndoState(value);
            
            // Find start of current line
            let lineStart = start;
            while (lineStart > 0 && value[lineStart - 1] !== '\\n') {
                lineStart--;
            }
            
            editor.value = value.substring(0, lineStart) + text + value.substring(lineStart);
            editor.focus();
            const newPos = lineStart + text.length;
            editor.setSelectionRange(newPos, newPos);
            
            // Save new state
            lastValue = editor.value;
            
            // Trigger input event for auto-save
            editor.dispatchEvent(new Event('input'));
        }
        
        function handleEditorMode(mode) {
            const container = document.getElementById('editorContainer');
            const editorPane = container.querySelector('.editor-pane');
            const previewPane = document.getElementById('previewPane');
            const preview = document.getElementById('descriptionPreview');
            const editor = document.getElementById('descriptionEditor');
            
            switch(mode) {
                case 'edit':
                    editorPane.style.display = 'block';
                    editorPane.style.width = '100%';
                    previewPane.style.display = 'none';
                    break;
                case 'preview':
                    editorPane.style.display = 'none';
                    previewPane.style.display = 'block';
                    previewPane.style.width = '100%';
                    preview.innerHTML = convertMarkdownToHtml(editor.value);
                    break;
                case 'split':
                    editorPane.style.display = 'block';
                    editorPane.style.width = '50%';
                    previewPane.style.display = 'block';
                    previewPane.style.width = '50%';
                    preview.innerHTML = convertMarkdownToHtml(editor.value);
                    container.style.display = 'flex';
                    break;
            }
        }
        
        function convertMarkdownToHtml(markdown) {
            let html = markdown;
            
            // Code blocks first (to avoid conflicts)
            html = html.replace(/\`\`\`([^\\n]*)\\n([^\`]+)\`\`\`/g, '<pre><code class="language-$1">$2</code></pre>');
            
            // Headers
            html = html.replace(/^### (.+)$/gm, '<h3>$1</h3>');
            html = html.replace(/^## (.+)$/gm, '<h2>$1</h2>');
            html = html.replace(/^# (.+)$/gm, '<h1>$1</h1>');
            
            // Bold and italic
            html = html.replace(/\\*\\*\\*(.+?)\\*\\*\\*/g, '<strong><em>$1</em></strong>');
            html = html.replace(/\\*\\*(.+?)\\*\\*/g, '<strong>$1</strong>');
            html = html.replace(/\\*(.+?)\\*/g, '<em>$1</em>');
            html = html.replace(/~~(.+?)~~/g, '<del>$1</del>');
            
            // Inline code
            html = html.replace(/\`(.+?)\`/g, '<code>$1</code>');
            
            // Links and images
            html = html.replace(/!\\[([^\\]]+)\\]\\(([^\\)]+)\\)/g, '<img src="$2" alt="$1" />');
            html = html.replace(/\\[([^\\]]+)\\]\\(([^\\)]+)\\)/g, '<a href="$2">$1</a>');
            
            // Lists
            html = html.replace(/^\\* (.+)$/gm, '<li>$1</li>');
            html = html.replace(/^- (.+)$/gm, '<li>$1</li>');
            html = html.replace(/^\\d+\\. (.+)$/gm, '<li>$1</li>');
            
            // Checkboxes
            html = html.replace(/^- \\[x\\] (.+)$/gm, '<li><input type="checkbox" checked disabled> $1</li>');
            html = html.replace(/^- \\[ \\] (.+)$/gm, '<li><input type="checkbox" disabled> $1</li>');
            
            // Wrap consecutive li elements in ul
            html = html.replace(/(<li>.*<\\/li>\\s*)+/g, function(match) {
                return '<ul>' + match + '</ul>';
            });
            
            // Line breaks
            html = html.replace(/\\n\\n/g, '</p><p>');
            html = html.replace(/\\n/g, '<br>');
            
            // Wrap in paragraphs
            if (!html.startsWith('<')) {
                html = '<p>' + html + '</p>';
            }
            
            return html;
        }
        
        // Undo/Redo functionality
        let undoStack = [];
        let redoStack = [];
        let lastValue = document.getElementById('descriptionEditor').value;
        
        function saveUndoState(value) {
            if (value !== lastValue) {
                undoStack.push(lastValue);
                redoStack = []; // Clear redo stack on new change
                lastValue = value;
                // Limit undo stack size
                if (undoStack.length > 50) {
                    undoStack.shift();
                }
            }
        }
        
        function performUndo() {
            const editor = document.getElementById('descriptionEditor');
            if (undoStack.length > 0) {
                redoStack.push(editor.value);
                const previousValue = undoStack.pop();
                editor.value = previousValue;
                lastValue = previousValue;
                editor.dispatchEvent(new Event('input'));
                console.log('Undo performed');
            }
        }
        
        function performRedo() {
            const editor = document.getElementById('descriptionEditor');
            if (redoStack.length > 0) {
                undoStack.push(editor.value);
                const nextValue = redoStack.pop();
                editor.value = nextValue;
                lastValue = nextValue;
                editor.dispatchEvent(new Event('input'));
                console.log('Redo performed');
            }
        }
        
        // Track changes for undo/redo
        let undoTimer;
        document.getElementById('descriptionEditor').addEventListener('input', function() {
            clearTimeout(undoTimer);
            undoTimer = setTimeout(() => {
                saveUndoState(this.value);
            }, 500); // Save state after 500ms of no typing
        });
        
        // Keyboard shortcuts
        document.getElementById('descriptionEditor').addEventListener('keydown', function(e) {
            if ((e.ctrlKey || e.metaKey) && !e.shiftKey && e.key === 'z') {
                e.preventDefault();
                performUndo();
            } else if ((e.ctrlKey || e.metaKey) && (e.shiftKey && e.key === 'Z' || e.key === 'y')) {
                e.preventDefault();
                performRedo();
            } else if ((e.ctrlKey || e.metaKey) && e.key === 'b') {
                e.preventDefault();
                toggleBold();
            } else if ((e.ctrlKey || e.metaKey) && e.key === 'i') {
                e.preventDefault();
                toggleItalic();
            } else if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
                e.preventDefault();
                insertLink();
            }
        });
        
        // Update preview in split mode as user types
        let updatePreviewTimeout;
        document.getElementById('descriptionEditor').addEventListener('input', function() {
            const mode = document.querySelector('.mode-btn.active').dataset.mode;
            if (mode === 'split') {
                clearTimeout(updatePreviewTimeout);
                updatePreviewTimeout = setTimeout(() => {
                    const preview = document.getElementById('descriptionPreview');
                    preview.innerHTML = convertMarkdownToHtml(this.value);
                }, 300);
            }
        });
        
        // Handle messages from VS Code
        window.addEventListener('message', event => {
            const message = event.data;
            switch (message.command) {
                case 'updateTask':
                    // Update UI with new task data
                    location.reload();
                    break;
                case 'insertText':
                    insertText(message.text);
                    if (message.cursorOffset) {
                        const editor = document.getElementById('descriptionEditor');
                        const pos = editor.selectionStart + message.cursorOffset;
                        editor.setSelectionRange(pos, pos);
                    }
                    break;
            }
        });
    </script>
</body>
</html>`;
    }
    
    private renderBreadcrumb(): string {
        // TODO: Implement actual hierarchy fetching
        return `
            <div class="breadcrumb-navigation">
                <a href="#" class="breadcrumb-item">
                    <i class="codicon codicon-home"></i> Ana Sayfa
                </a>
                <i class="codicon codicon-chevron-right"></i>
                <a href="#" class="breadcrumb-item">Projeler</a>
                <i class="codicon codicon-chevron-right"></i>
                <a href="#" class="breadcrumb-item">Frontend</a>
                <i class="codicon codicon-chevron-right"></i>
                <span class="breadcrumb-current">${this.escapeHtml(this.task.baslik)}</span>
            </div>
        `;
    }
    
    private renderEnhancedHierarchySection(): string {
        // Check both hierarchyInfo and task.alt_gorevler
        const hasSubtasks = (this.task.alt_gorevler && this.task.alt_gorevler.length > 0);
        const hasHierarchyInfo = this.hierarchyInfo && this.hierarchyInfo.toplam_alt_gorev > 0;
        const hasHierarchy = hasSubtasks || hasHierarchyInfo || this.task.parent_id;
        
        // Calculate progress from task.alt_gorevler if hierarchyInfo is not available
        let progressInfo = this.hierarchyInfo;
        if (!progressInfo && hasSubtasks) {
            // Count all subtasks recursively
            const counts = this.countAllSubtasks(this.task.alt_gorevler!);
            const totalSubtasks = counts.total;
            const completedSubtasks = counts.completed;
            const inProgressSubtasks = counts.inProgress;
            const pendingSubtasks = counts.pending;
            const progressPercentage = totalSubtasks > 0 ? Math.round((completedSubtasks / totalSubtasks) * 100) : 0;
            
            progressInfo = {
                gorev: this.task,
                ust_gorevler: [],
                toplam_alt_gorev: totalSubtasks,
                tamamlanan_alt: completedSubtasks,
                devam_eden_alt: inProgressSubtasks,
                beklemede_alt: pendingSubtasks,
                ilerleme_yuzdesi: progressPercentage
            };
        }
        
        return `
            <div class="hierarchy-section card">
                <h3><i class="codicon codicon-type-hierarchy"></i> ${vscode.l10n.t('taskDetail.hierarchy')}</h3>
                
                ${hasHierarchy ? `
                    <!-- Progress Overview -->
                    ${progressInfo && progressInfo.toplam_alt_gorev > 0 ? `
                        <div class="progress-overview">
                            <div class="circular-progress">
                                <svg viewBox="0 0 36 36" class="circular-chart">
                                    <path class="circle-bg"
                                        d="M18 2.0845
                                        a 15.9155 15.9155 0 0 1 0 31.831
                                        a 15.9155 15.9155 0 0 1 0 -31.831"
                                    />
                                    <path class="circle"
                                        stroke-dasharray="${progressInfo.ilerleme_yuzdesi || 0}, 100"
                                        d="M18 2.0845
                                        a 15.9155 15.9155 0 0 1 0 31.831
                                        a 15.9155 15.9155 0 0 1 0 -31.831"
                                    />
                                </svg>
                                <div class="percentage-overlay">${Math.round(progressInfo.ilerleme_yuzdesi || 0)}%</div>
                            </div>
                            <div class="progress-details">
                                <div class="stat-item">
                                    <span class="stat-value">${progressInfo.toplam_alt_gorev}</span>
                                    <span class="stat-label">Toplam</span>
                                </div>
                                <div class="stat-item success">
                                    <span class="stat-value">${progressInfo.tamamlanan_alt}</span>
                                    <span class="stat-label">Tamamlandı</span>
                                </div>
                                <div class="stat-item warning">
                                    <span class="stat-value">${progressInfo.devam_eden_alt || 0}</span>
                                    <span class="stat-label">Devam Ediyor</span>
                                </div>
                            </div>
                        </div>
                    ` : ''}
                    
                    <!-- Task Tree -->
                    <div class="task-tree">
                        ${this.renderTaskTree()}
                    </div>
                ` : `
                    <div class="empty-state">
                        <i class="codicon codicon-type-hierarchy"></i>
                        <p>Bu görev henüz bir hiyerarşiye sahip değil</p>
                    </div>
                `}
                
                <div class="hierarchy-actions">
                    <button class="action-button small" id="createSubtaskBtn">
                        <i class="codicon codicon-add"></i> Alt Görev
                    </button>
                    ${this.task.parent_id ? `
                        <button class="action-button small" onclick="vscode.postMessage({command: 'removeParent'})">
                            <i class="codicon codicon-ungroup-by-ref-type"></i> Bağımsız Yap
                        </button>
                    ` : `
                        <button class="action-button small" onclick="vscode.postMessage({command: 'changeParent'})">
                            <i class="codicon codicon-type-hierarchy-sub"></i> Üst Görev Ata
                        </button>
                    `}
                </div>
            </div>
        `;
    }
    
    private renderTaskTree(): string {
        // Show actual task hierarchy
        let treeHtml = '';
        
        // If task has parent, show it
        if (this.task.parent_id) {
            treeHtml += `
                <div class="tree-item parent">
                    <span class="tree-icon"><i class="codicon codicon-chevron-down"></i></span>
                    <span class="tree-content">
                        <i class="codicon codicon-symbol-class"></i> Üst Görev
                    </span>
                </div>
            `;
        }
        
        // Show current task
        treeHtml += `
            <div class="tree-item ${this.task.parent_id ? 'child' : ''} current">
                <span class="tree-icon"></span>
                <span class="tree-content">
                    <i class="codicon codicon-circle-filled"></i> ${this.escapeHtml(this.task.baslik)}
                    <span class="tree-badge">Şu an</span>
                </span>
            </div>
        `;
        
        // Show subtasks if any
        if (this.task.alt_gorevler && this.task.alt_gorevler.length > 0) {
            treeHtml += this.renderSubtasks(this.task.alt_gorevler, 1);
        }
        
        return treeHtml;
    }
    
    private countAllSubtasks(subtasks: Gorev[]): { total: number; completed: number; inProgress: number; pending: number } {
        let counts = { total: 0, completed: 0, inProgress: 0, pending: 0 };
        
        subtasks.forEach(subtask => {
            counts.total++;
            
            switch (subtask.durum) {
                case GorevDurum.Tamamlandi:
                    counts.completed++;
                    break;
                case GorevDurum.DevamEdiyor:
                    counts.inProgress++;
                    break;
                default:
                    counts.pending++;
                    break;
            }
            
            // Recursively count sub-subtasks
            if (subtask.alt_gorevler && subtask.alt_gorevler.length > 0) {
                const subCounts = this.countAllSubtasks(subtask.alt_gorevler);
                counts.total += subCounts.total;
                counts.completed += subCounts.completed;
                counts.inProgress += subCounts.inProgress;
                counts.pending += subCounts.pending;
            }
        });
        
        return counts;
    }
    
    private renderSubtasks(subtasks: Gorev[], level: number): string {
        let html = '';
        const indent = '    '.repeat(level);
        
        subtasks.forEach(subtask => {
            const statusIcon = this.getSubtaskStatusIcon(subtask.durum);
            const statusClass = this.getSubtaskStatusClass(subtask.durum);
            const hasChildren = subtask.alt_gorevler && subtask.alt_gorevler.length > 0;
            
            html += `
                <div class="tree-item child" style="padding-left: ${level * 20}px;">
                    <span class="tree-icon">
                        ${hasChildren ? '<i class="codicon codicon-chevron-right"></i>' : ''}
                    </span>
                    <span class="tree-content">
                        <i class="codicon codicon-symbol-method"></i> ${this.escapeHtml(subtask.baslik)}
                        <span class="tree-status ${statusClass}">${statusIcon}</span>
                    </span>
                </div>
            `;
            
            // Recursively render sub-subtasks
            if (hasChildren) {
                html += this.renderSubtasks(subtask.alt_gorevler!, level + 1);
            }
        });
        
        return html;
    }
    
    private getSubtaskStatusIcon(durum: GorevDurum): string {
        switch (durum) {
            case GorevDurum.Tamamlandi: return '✓';
            case GorevDurum.DevamEdiyor: return '🔄';
            default: return '⏳';
        }
    }
    
    private getSubtaskStatusClass(durum: GorevDurum): string {
        switch (durum) {
            case GorevDurum.Tamamlandi: return 'completed';
            case GorevDurum.DevamEdiyor: return 'in-progress';
            default: return 'pending';
        }
    }
    
    private renderHierarchySection(): string {
        // Keep old method for backward compatibility
        return this.renderEnhancedHierarchySection();
    }
    
    private renderDependenciesSection(): string {
        // Debug: Log dependency information
        console.log('Task dependency info:', {
            bagimli_gorev_sayisi: this.task.bagimli_gorev_sayisi,
            tamamlanmamis_bagimlilik_sayisi: this.task.tamamlanmamis_bagimlilik_sayisi,
            bu_goreve_bagimli_sayisi: this.task.bu_goreve_bagimli_sayisi,
            bagimliliklar: this.task.bagimliliklar,
            taskId: this.task.id,
            taskTitle: this.task.baslik
        });
        
        const hasDependencyInfo = this.task.bagimli_gorev_sayisi || this.task.bu_goreve_bagimli_sayisi || 
                                  (this.task.bagimliliklar && this.task.bagimliliklar.length > 0);
        
        let html = `
            <div class="dependencies-section card">
                <h3><i class="codicon codicon-link"></i> ${vscode.l10n.t('taskDetail.dependencies')}</h3>
        `;
        
        // Summary stats or empty state
        if (this.task.bagimli_gorev_sayisi || this.task.bu_goreve_bagimli_sayisi) {
            html += `
                <div class="dependency-stats">
                    ${this.task.bagimli_gorev_sayisi ? `
                        <div class="stat-item">
                            <i class="codicon codicon-arrow-right"></i>
                            <span class="stat-label">${vscode.l10n.t('taskDetail.thisDependsOn')}</span>
                            <span class="stat-value">${vscode.l10n.t('taskDetail.dependentTasks', this.task.bagimli_gorev_sayisi.toString())}</span>
                            ${this.task.tamamlanmamis_bagimlilik_sayisi ? `
                                <span class="stat-warning">⚠️ ${vscode.l10n.t('taskDetail.incompleteCount', this.task.tamamlanmamis_bagimlilik_sayisi.toString())}</span>
                            ` : '<span class="stat-success">✓ ' + vscode.l10n.t('taskDetail.allCompleted') + '</span>'}
                        </div>
                    ` : ''}
                    
                    ${this.task.bu_goreve_bagimli_sayisi ? `
                        <div class="stat-item">
                            <i class="codicon codicon-arrow-left"></i>
                            <span class="stat-label">${vscode.l10n.t('taskDetail.dependsOnThis')}</span>
                            <span class="stat-value">${vscode.l10n.t('taskDetail.dependentTasks', this.task.bu_goreve_bagimli_sayisi.toString())}</span>
                        </div>
                    ` : ''}
                </div>
            `;
        } else if (!hasDependencyInfo) {
            // Empty state
            html += `
                <div class="empty-state">
                    <i class="codicon codicon-link"></i>
                    <p>${vscode.l10n.t('taskDetail.noDependenciesYet')}</p>
                </div>
            `;
        }
        
        // Dependency list (if available)
        if (this.task.bagimliliklar && this.task.bagimliliklar.length > 0) {
            html += `
                <div class="dependency-list compact">
                    <h4>Bağımlı Olduğu Görevler:</h4>
                    ${this.task.bagimliliklar.map((dep: any) => `
                        <div class="dependency-item">
                            <span class="dep-status ${this.getDepStatusClass(dep.hedef_durum || 'beklemede')}">
                                <i class="codicon ${this.getDepStatusIcon(dep.hedef_durum || 'beklemede')}"></i>
                            </span>
                            <span class="dep-title">${this.escapeHtml(dep.hedef_baslik || 'Görev')}</span>
                            <button class="link-button" onclick="vscode.postMessage({command: 'openTask', taskId: '${dep.hedef_id}'})" title="Görevi Aç">
                                <i class="codicon codicon-arrow-right"></i>
                            </button>
                        </div>
                    `).join('')}
                </div>
            `;
        }
        
        html += `
                <button class="add-button" id="addDependencyBtn" onclick="vscode.postMessage({command: 'addDependency'})">
                    <i class="codicon codicon-add"></i> Bağımlılık Ekle
                </button>
            </div>
        `;
        
        return html;
    }
    
    private getDepStatusClass(durum: string): string {
        if (durum.includes('tamamland')) return 'completed';
        if (durum.includes('devam')) return 'in-progress';
        return 'pending';
    }
    
    private renderActivityTimeline(): string {
        let html = '';
        
        // Oluşturulma aktivitesi
        html += `
            <div class="timeline-item">
                <span class="timeline-icon"><i class="codicon codicon-add"></i></span>
                <div class="timeline-content">
                    <div class="timeline-title">Oluşturuldu</div>
                    <div class="timeline-time">${this.formatRelativeTime(this.task.olusturma_tarih)}</div>
                </div>
            </div>
        `;
        
        // Durum değişiklikleri
        if (this.task.durum === GorevDurum.DevamEdiyor) {
            html += `
                <div class="timeline-item">
                    <span class="timeline-icon"><i class="codicon codicon-debug-start"></i></span>
                    <div class="timeline-content">
                        <div class="timeline-title">Başlatıldı</div>
                        <div class="timeline-time">${this.formatRelativeTime(this.task.guncelleme_tarih)}</div>
                    </div>
                </div>
            `;
        } else if (this.task.durum === GorevDurum.Tamamlandi) {
            // Başlatılma (varsa)
            if (this.task.guncelleme_tarih !== this.task.olusturma_tarih) {
                html += `
                    <div class="timeline-item">
                        <span class="timeline-icon"><i class="codicon codicon-debug-start"></i></span>
                        <div class="timeline-content">
                            <div class="timeline-title">Başlatıldı</div>
                            <div class="timeline-time">-</div>
                        </div>
                    </div>
                `;
            }
            
            // Tamamlanma
            html += `
                <div class="timeline-item">
                    <span class="timeline-icon"><i class="codicon codicon-pass-filled"></i></span>
                    <div class="timeline-content">
                        <div class="timeline-title">Tamamlandı</div>
                        <div class="timeline-time">${this.formatRelativeTime(this.task.guncelleme_tarih)}</div>
                    </div>
                </div>
            `;
        }
        
        // Son güncelleme (farklıysa)
        if (this.task.guncelleme_tarih && 
            this.task.guncelleme_tarih !== this.task.olusturma_tarih &&
            this.task.durum === GorevDurum.Beklemede) {
            html += `
                <div class="timeline-item">
                    <span class="timeline-icon"><i class="codicon codicon-edit"></i></span>
                    <div class="timeline-content">
                        <div class="timeline-title">Güncellendi</div>
                        <div class="timeline-time">${this.formatRelativeTime(this.task.guncelleme_tarih)}</div>
                    </div>
                </div>
            `;
        }
        
        return html;
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
                case 'insertLink':
                    await this.handleInsertLink(message.selectedText);
                    break;
                case 'insertImage':
                    await this.handleInsertImage();
                    break;
                case 'insertCodeBlock':
                    await this.handleInsertCodeBlock();
                    break;
                case 'insertTable':
                    await this.handleInsertTable();
                    break;
                    
                case 'addDependency':
                    await this.showDependencyPicker();
                    break;
                    
                case 'openTask':
                    await this.openTask(message.taskId);
                    break;
                    
                case 'createSubtask':
                    await vscode.commands.executeCommand('gorev.createSubtask', { task: this.task });
                    break;
                    
                case 'changeParent':
                    await vscode.commands.executeCommand('gorev.changeParent', { task: this.task });
                    break;
                    
                case 'removeParent':
                    await vscode.commands.executeCommand('gorev.removeParent', { task: this.task });
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
        // Create a tree item to pass to the command
        const treeItem = {
            task: this.task
        };
        await vscode.commands.executeCommand('gorev.detailedEdit', treeItem);
    }
    
    private async deleteTask() {
        const confirm = await vscode.window.showWarningMessage(
            `"${this.task.baslik}" görevini silmek istediğinizden emin misiniz?`,
            'Evet, Sil',
            'İptal'
        );
        
        if (confirm === 'Evet, Sil') {
            await this.mcpClient.callTool('gorev_sil', {
                id: this.task.id,
                onay: true
            });
            
            this.panel.dispose();
            vscode.window.showInformationMessage('Görev silindi');
            
            // Refresh the tree view
            await vscode.commands.executeCommand('gorev.refreshTasks');
        }
    }
    
    private async addTag(tag: string) {
        const currentTags = this.task.etiketler || [];
        currentTags.push(tag);
        
        await this.updateTaskField('etiketler', currentTags.join(','));
        this.task.etiketler = currentTags;
        this.update();
    }
    
    private async handleInsertLink(selectedText?: string) {
        const url = await vscode.window.showInputBox({
            prompt: 'Link URL\'sini girin',
            placeHolder: 'https://example.com'
        });
        
        if (url) {
            const linkText = selectedText || await vscode.window.showInputBox({
                prompt: 'Link metni',
                placeHolder: 'Link açıklaması',
                value: selectedText || ''
            }) || url;
            
            this.panel.webview.postMessage({
                command: 'insertText',
                text: `[${linkText}](${url})`
            });
        }
    }
    
    private async handleInsertImage() {
        const url = await vscode.window.showInputBox({
            prompt: 'Resim URL\'sini girin',
            placeHolder: 'https://example.com/image.png'
        });
        
        if (url) {
            const altText = await vscode.window.showInputBox({
                prompt: 'Alternatif metin',
                placeHolder: 'Resim açıklaması',
                value: 'Resim'
            }) || 'Resim';
            
            this.panel.webview.postMessage({
                command: 'insertText',
                text: `![${altText}](${url})`
            });
        }
    }
    
    private async handleInsertCodeBlock() {
        const language = await vscode.window.showInputBox({
            prompt: 'Programlama dili (opsiyonel)',
            placeHolder: 'javascript, python, go, vb.'
        }) || '';
        
        this.panel.webview.postMessage({
            command: 'insertText',
            text: `\n\`\`\`${language}\n\n\`\`\`\n`,
            cursorOffset: -5
        });
    }
    
    private async handleInsertTable() {
        const colsStr = await vscode.window.showInputBox({
            prompt: 'Kolon sayısını girin',
            placeHolder: '3',
            value: '3'
        });
        
        if (colsStr) {
            const cols = parseInt(colsStr) || 3;
            let table = '\n| ';
            for (let i = 0; i < cols; i++) {
                table += `Başlık ${i + 1} | `;
            }
            table += '\n| ';
            for (let i = 0; i < cols; i++) {
                table += '--- | ';
            }
            table += '\n| ';
            for (let i = 0; i < cols; i++) {
                table += 'Hücre | ';
            }
            table += '\n';
            
            this.panel.webview.postMessage({
                command: 'insertText',
                text: table
            });
        }
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
    
    private getStatusLabel(): string {
        switch (this.task.durum) {
            case GorevDurum.Tamamlandi: return vscode.l10n.t('taskDetail.status.completed');
            case GorevDurum.DevamEdiyor: return vscode.l10n.t('taskDetail.status.inProgress');
            default: return vscode.l10n.t('taskDetail.status.pending');
        }
    }
    
    private getDepStatusIcon(durum: string): string {
        if (durum.includes('tamamland')) return 'codicon-pass-filled';
        if (durum.includes('devam')) return 'codicon-debug-start';
        return 'codicon-circle-outline';
    }
    
    private getPriorityLabel(): string {
        switch (this.task.oncelik) {
            case GorevOncelik.Yuksek: return vscode.l10n.t('taskDetail.priority.high');
            case GorevOncelik.Orta: return vscode.l10n.t('taskDetail.priority.medium');
            case GorevOncelik.Dusuk: return vscode.l10n.t('taskDetail.priority.low');
            default: return vscode.l10n.t('taskDetail.priority.medium');
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
    
    private formatRelativeTime(dateStr?: string): string {
        if (!dateStr) return '';
        
        const date = new Date(dateStr);
        const now = new Date();
        const diffMs = now.getTime() - date.getTime();
        const diffSecs = Math.floor(diffMs / 1000);
        const diffMins = Math.floor(diffSecs / 60);
        const diffHours = Math.floor(diffMins / 60);
        const diffDays = Math.floor(diffHours / 24);
        
        if (diffSecs < 60) {
            return 'Az önce';
        } else if (diffMins < 60) {
            return `${diffMins} dakika önce`;
        } else if (diffHours < 24) {
            return `${diffHours} saat önce`;
        } else if (diffDays < 7) {
            return `${diffDays} gün önce`;
        } else {
            return this.formatDate(dateStr);
        }
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