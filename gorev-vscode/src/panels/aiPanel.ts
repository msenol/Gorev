import * as vscode from 'vscode';
import { Logger } from '../utils/logger';

/**
 * Webview message types
 */
interface WebviewMessage {
  type: 'configure' | 'chat' | 'analyze' | 'decompose' | 'estimate' | 'search';
  payload: unknown;
}

interface ConfigurePayload {
  provider: string;
  apiKey: string;
  model?: string;
  temperature?: number;
  maxTokens?: number;
}

interface ChatPayload {
  message: string;
}

interface AnalyzePayload {
  projectId: string;
  projectName: string;
  taskCount: number;
  completedCount: number;
  pendingCount: number;
}

interface DecomposePayload {
  taskId: string;
  title: string;
  description: string;
  maxDepth?: number;
}

interface EstimatePayload {
  taskId: string;
  title: string;
  description: string;
  tags?: string[];
}

interface SearchPayload {
  query: string;
  projectId?: string;
  limit?: number;
}

/**
 * AI Panel Webview Provider
 * Provides AI chat interface and configuration UI
 * Hidden if AI is not configured (Rule 15 compliance)
 */
export class AIPanelProvider {
  private static instance: AIPanelProvider;
  private panel: vscode.WebviewPanel | undefined;
  private disposables: vscode.Disposable[] = [];

  private constructor(
    private extensionUri: vscode.Uri,
    private context: vscode.ExtensionContext
  ) {}

  static getInstance(extensionUri: vscode.Uri, context: vscode.ExtensionContext): AIPanelProvider {
    if (!AIPanelProvider.instance) {
      AIPanelProvider.instance = new AIPanelProvider(extensionUri, context);
    }
    return AIPanelProvider.instance;
  }

  /**
   * Show or create the AI panel
   */
  async show(aiConfigured: boolean): Promise<void> {
    if (!aiConfigured) {
      vscode.window.showInformationMessage(
        'AI is not configured. Use the gorev_ai tool with action=configure to set up AI for this project.',
        'Configure Now'
      ).then(selection => {
        if (selection === 'Configure Now') {
          vscode.commands.executeCommand('gorev.ai.configure');
        }
      });
      return;
    }

    if (this.panel) {
      this.panel.reveal();
      return;
    }

    this.panel = vscode.window.createWebviewPanel(
      'gorev.aiPanel',
      'Gorev AI Assistant',
      vscode.ViewColumn.Beside,
      {
        enableScripts: true,
        retainContextWhenHidden: true,
        localResourceRoots: [
          vscode.Uri.joinPath(this.extensionUri, 'media'),
          vscode.Uri.joinPath(this.extensionUri, 'dist')
        ]
      }
    );

    this.panel.webview.html = this.getWebviewContent();

    // Handle messages from webview
    this.panel.webview.onDidReceiveMessage(
      async (message) => {
        await this.handleMessage(message);
      },
      null,
      this.disposables
    );

    // Handle panel disposal
    this.panel.onDidDispose(
      () => {
        this.panel = undefined;
        this.dispose();
      },
      null,
      this.disposables
    );
  }

  /**
   * Handle messages from webview
   */
  private async handleMessage(message: WebviewMessage): Promise<void> {
    Logger.debug('AI Panel received message:', message);

    switch (message.type) {
      case 'configure':
        await this.handleConfigure(message.payload as ConfigurePayload);
        break;
      case 'chat':
        await this.handleChat(message.payload as ChatPayload);
        break;
      case 'analyze':
        await this.handleAnalyze(message.payload as AnalyzePayload);
        break;
      case 'decompose':
        await this.handleDecompose(message.payload as DecomposePayload);
        break;
      case 'estimate':
        await this.handleEstimate(message.payload as EstimatePayload);
        break;
      case 'search':
        await this.handleSearch(message.payload as SearchPayload);
        break;
      default:
        Logger.warn('Unknown message type:', message.type);
    }
  }

  /**
   * Handle AI configuration
   */
  private async handleConfigure(payload: ConfigurePayload): Promise<void> {
    try {
      const result = await vscode.commands.executeCommand(
        'gorev.mcp.call',
        'gorev_ai',
        `action=configure provider=${payload.provider} api_key=${payload.apiKey} model=${payload.model || ''}`
      );

      this.panel?.webview.postMessage({
        type: 'configureResult',
        success: true,
        result
      });
    } catch (error) {
      Logger.error('AI configuration failed:', error);
      this.panel?.webview.postMessage({
        type: 'configureResult',
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  }

  /**
   * Handle AI chat
   */
  private async handleChat(payload: ChatPayload): Promise<void> {
    try {
      const result = await vscode.commands.executeCommand(
        'gorev.mcp.call',
        'gorev_ai',
        `action=chat message=${encodeURIComponent(payload.message)}`
      );

      this.panel?.webview.postMessage({
        type: 'chatResponse',
        success: true,
        result
      });
    } catch (error) {
      Logger.error('AI chat failed:', error);
      this.panel?.webview.postMessage({
        type: 'chatResponse',
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  }

  /**
   * Handle project analysis
   */
  private async handleAnalyze(payload: AnalyzePayload): Promise<void> {
    try {
      const result = await vscode.commands.executeCommand(
        'gorev.mcp.call',
        'gorev_ai',
        `action=analyze project_id=${payload.projectId || ''}`
      );

      this.panel?.webview.postMessage({
        type: 'analyzeResult',
        success: true,
        result
      });
    } catch (error) {
      Logger.error('AI analysis failed:', error);
      this.panel?.webview.postMessage({
        type: 'analyzeResult',
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  }

  /**
   * Handle task decomposition
   */
  private async handleDecompose(payload: DecomposePayload): Promise<void> {
    try {
      const params = new URLSearchParams({
        action: 'decompose',
        task_id: payload.taskId,
        title: payload.title,
        description: payload.description,
        max_depth: String(payload.maxDepth ?? 3)
      });

      const result = await vscode.commands.executeCommand(
        'gorev.mcp.call',
        'gorev_ai',
        params.toString()
      );

      this.panel?.webview.postMessage({
        type: 'decomposeResult',
        success: true,
        result
      });
    } catch (error) {
      Logger.error('Task decomposition failed:', error);
      this.panel?.webview.postMessage({
        type: 'decomposeResult',
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  }

  /**
   * Handle time estimation
   */
  private async handleEstimate(payload: EstimatePayload): Promise<void> {
    try {
      const tags = payload.tags?.join(',') || '';
      const params = new URLSearchParams({
        action: 'estimate',
        task_id: payload.taskId,
        title: payload.title,
        description: payload.description,
        tags: tags
      });

      const result = await vscode.commands.executeCommand(
        'gorev.mcp.call',
        'gorev_ai',
        params.toString()
      );

      this.panel?.webview.postMessage({
        type: 'estimateResult',
        success: true,
        result
      });
    } catch (error) {
      Logger.error('Time estimation failed:', error);
      this.panel?.webview.postMessage({
        type: 'estimateResult',
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  }

  /**
   * Handle semantic search
   */
  private async handleSearch(payload: SearchPayload): Promise<void> {
    try {
      const result = await vscode.commands.executeCommand(
        'gorev.mcp.call',
        'gorev_ai',
        `action=search query=${encodeURIComponent(payload.query)} limit=${payload.limit || 10}`
      );

      this.panel?.webview.postMessage({
        type: 'searchResult',
        success: true,
        result
      });
    } catch (error) {
      Logger.error('AI search failed:', error);
      this.panel?.webview.postMessage({
        type: 'searchResult',
        success: false,
        error: error instanceof Error ? error.message : 'Unknown error'
      });
    }
  }

  /**
   * Get webview HTML content
   */
  private getWebviewContent(): string {
    const nonce = this.getNonce();

    return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="Content-Security-Policy" content="
    default-src 'none';
    script-src 'nonce-${nonce}' 'unsafe-inline';
    style-src 'nonce-${nonce}' 'unsafe-inline';
  ">
  <title>Gorev AI Assistant</title>
  <style>
    body {
      font-family: var(--vscode-font-family);
      color: var(--vscode-foreground);
      background-color: var(--vscode-editor-background);
      padding: 16px;
      margin: 0;
    }
    .container {
      max-width: 800px;
      margin: 0 auto;
    }
    .header {
      display: flex;
      align-items: center;
      margin-bottom: 20px;
      padding-bottom: 10px;
      border-bottom: 1px solid var(--vscode-panel-border);
    }
    .header h1 {
      font-size: 18px;
      margin: 0;
      flex: 1;
    }
    .tabs {
      display: flex;
      gap: 10px;
      margin-bottom: 20px;
    }
    .tab {
      padding: 8px 16px;
      background: var(--vscode-button-secondaryBackground);
      border: none;
      color: var(--vscode-button-secondaryForeground);
      cursor: pointer;
      border-radius: 4px;
    }
    .tab.active {
      background: var(--vscode-button-background);
      color: var(--vscode-button-foreground);
    }
    .tab-content {
      display: none;
    }
    .tab-content.active {
      display: block;
    }
    .form-group {
      margin-bottom: 16px;
    }
    label {
      display: block;
      margin-bottom: 8px;
      font-weight: 500;
    }
    input, select, textarea {
      width: 100%;
      padding: 8px;
      background: var(--vscode-input-background);
      border: 1px solid var(--vscode-input-border);
      color: var(--vscode-input-foreground);
      border-radius: 4px;
      box-sizing: border-box;
    }
    button {
      padding: 8px 16px;
      background: var(--vscode-button-background);
      color: var(--vscode-button-foreground);
      border: none;
      border-radius: 4px;
      cursor: pointer;
    }
    button:hover {
      background: var(--vscode-button-hoverBackground);
    }
    .chat-container {
      height: 400px;
      display: flex;
      flex-direction: column;
      border: 1px solid var(--vscode-panel-border);
      border-radius: 4px;
      overflow: hidden;
    }
    .chat-messages {
      flex: 1;
      overflow-y: auto;
      padding: 16px;
    }
    .message {
      margin-bottom: 12px;
      padding: 8px 12px;
      border-radius: 8px;
      max-width: 80%;
    }
    .message.user {
      background: var(--vscode-button-background);
      margin-left: auto;
    }
    .message.assistant {
      background: var(--vscode-editor-selectionBackground);
    }
    .chat-input {
      display: flex;
      gap: 8px;
      padding: 8px;
      border-top: 1px solid var(--vscode-panel-border);
    }
    .chat-input input {
      flex: 1;
    }
    .result-box {
      margin-top: 16px;
      padding: 12px;
      background: var(--vscode-textBlockQuote-background);
      border-left: 4px solid var(--vscode-textBlockQuote-border);
      border-radius: 4px;
    }
    .error {
      color: var(--vscode-errorForeground);
    }
    .loading {
      opacity: 0.6;
      pointer-events: none;
    }
  </style>
</head>
<body>
  <div class="container">
    <div class="header">
      <h1>🤖 Gorev AI Assistant</h1>
    </div>

    <div class="tabs">
      <button class="tab active" onclick="switchTab('chat')">Chat</button>
      <button class="tab" onclick="switchTab('analyze')">Analyze</button>
      <button class="tab" onclick="switchTab('decompose')">Decompose</button>
      <button class="tab" onclick="switchTab('estimate')">Estimate</button>
      <button class="tab" onclick="switchTab('search')">Search</button>
      <button class="tab" onclick="switchTab('configure')">Configure</button>
    </div>

    <!-- Chat Tab -->
    <div id="chat-tab" class="tab-content active">
      <div class="chat-container">
        <div class="chat-messages" id="chatMessages">
          <div class="message assistant">
            Hello! I'm your AI assistant. How can I help you manage your tasks today?
          </div>
        </div>
        <div class="chat-input">
          <input type="text" id="chatInput" placeholder="Type your message..." />
          <button onclick="sendChat()">Send</button>
        </div>
      </div>
    </div>

    <!-- Analyze Tab -->
    <div id="analyze-tab" class="tab-content">
      <div class="form-group">
        <label>Project ID (optional, uses active project)</label>
        <input type="text" id="analyzeProjectId" placeholder="Leave empty for active project">
      </div>
      <button onclick="analyzeProject()">Analyze Project</button>
      <div id="analyzeResult"></div>
    </div>

    <!-- Decompose Tab -->
    <div id="decompose-tab" class="tab-content">
      <div class="form-group">
        <label>Task Title</label>
        <input type="text" id="decomposeTitle" placeholder="Enter task title">
      </div>
      <div class="form-group">
        <label>Task Description</label>
        <textarea id="decomposeDescription" rows="4" placeholder="Enter task description"></textarea>
      </div>
      <div class="form-group">
        <label>Max Depth</label>
        <input type="number" id="decomposeMaxDepth" value="3" min="1" max="5">
      </div>
      <button onclick="decomposeTask()">Decompose Task</button>
      <div id="decomposeResult"></div>
    </div>

    <!-- Estimate Tab -->
    <div id="estimate-tab" class="tab-content">
      <div class="form-group">
        <label>Task Title</label>
        <input type="text" id="estimateTitle" placeholder="Enter task title">
      </div>
      <div class="form-group">
        <label>Task Description</label>
        <textarea id="estimateDescription" rows="4" placeholder="Enter task description"></textarea>
      </div>
      <div class="form-group">
        <label>Tags (comma-separated)</label>
        <input type="text" id="estimateTags" placeholder="backend, api, feature">
      </div>
      <button onclick="estimateTime()">Estimate Time</button>
      <div id="estimateResult"></div>
    </div>

    <!-- Search Tab -->
    <div id="search-tab" class="tab-content">
      <div class="form-group">
        <label>Search Query</label>
        <input type="text" id="searchQuery" placeholder="Enter your search query">
      </div>
      <div class="form-group">
        <label>Limit</label>
        <input type="number" id="searchLimit" value="10" min="1" max="50">
      </div>
      <button onclick="semanticSearch()">Search</button>
      <div id="searchResult"></div>
    </div>

    <!-- Configure Tab -->
    <div id="configure-tab" class="tab-content">
      <div class="form-group">
        <label>AI Provider</label>
        <select id="configProvider">
          <option value="openrouter">OpenRouter</option>
          <option value="anannas">Anannas</option>
        </select>
      </div>
      <div class="form-group">
        <label>API Key</label>
        <input type="password" id="configApiKey" placeholder="Enter your API key">
      </div>
      <div class="form-group">
        <label>Model (optional)</label>
        <input type="text" id="configModel" placeholder="e.g., openai/gpt-4o-mini">
      </div>
      <button onclick="configureAI()">Configure AI</button>
      <div id="configureResult"></div>
    </div>
  </div>

  <script nonce="${nonce}">
    const vscode = acquireVsCodeApi();

    function switchTab(tabName) {
      // Update tab buttons
      document.querySelectorAll('.tab').forEach(tab => tab.classList.remove('active'));
      event.target.classList.add('active');

      // Update tab content
      document.querySelectorAll('.tab-content').forEach(content => content.classList.remove('active'));
      document.getElementById(tabName + '-tab').classList.add('active');
    }

    function addMessage(content: string, isUser: boolean) {
      const messagesDiv = document.getElementById('chatMessages');
      const messageDiv = document.createElement('div');
      messageDiv.className = 'message ' + (isUser ? 'user' : 'assistant');
      messageDiv.textContent = content;
      messagesDiv.appendChild(messageDiv);
      messagesDiv.scrollTop = messagesDiv.scrollHeight;
    }

    async function sendChat() {
      const input = document.getElementById('chatInput') as HTMLInputElement;
      const message = input.value.trim();
      if (!message) return;

      addMessage(message, true);
      input.value = '';

      try {
        const response = await vscode.postMessage({
          type: 'chat',
          payload: { message }
        });

        if (response.success) {
          addMessage(response.result, false);
        } else {
          addMessage('Error: ' + response.error, false);
        }
      } catch (error) {
        addMessage('Error: ' + error, false);
      }
    }

    async function analyzeProject() {
      const projectId = (document.getElementById('analyzeProjectId') as HTMLInputElement).value;
      const resultDiv = document.getElementById('analyzeResult');
      resultDiv.innerHTML = '<p>Analyzing...</p>';

      try {
        const response = await vscode.postMessage({
          type: 'analyze',
          payload: { projectId }
        });

        if (response.success) {
          resultDiv.innerHTML = '<pre>' + response.result + '</pre>';
        } else {
          resultDiv.innerHTML = '<p class="error">Error: ' + response.error + '</p>';
        }
      } catch (error) {
        resultDiv.innerHTML = '<p class="error">Error: ' + error + '</p>';
      }
    }

    async function decomposeTask() {
      const title = (document.getElementById('decomposeTitle') as HTMLInputElement).value;
      const description = (document.getElementById('decomposeDescription') as HTMLTextAreaElement).value;
      const maxDepth = (document.getElementById('decomposeMaxDepth') as HTMLInputElement).value;
      const resultDiv = document.getElementById('decomposeResult');
      resultDiv.innerHTML = '<p>Decomposing task...</p>';

      try {
        const response = await vscode.postMessage({
          type: 'decompose',
          payload: { title, description, maxDepth }
        });

        if (response.success) {
          resultDiv.innerHTML = '<pre>' + response.result + '</pre>';
        } else {
          resultDiv.innerHTML = '<p class="error">Error: ' + response.error + '</p>';
        }
      } catch (error) {
        resultDiv.innerHTML = '<p class="error">Error: ' + error + '</p>';
      }
    }

    async function estimateTime() {
      const title = (document.getElementById('estimateTitle') as HTMLInputElement).value;
      const description = (document.getElementById('estimateDescription') as HTMLTextAreaElement).value;
      const tags = (document.getElementById('estimateTags') as HTMLInputElement).value;
      const resultDiv = document.getElementById('estimateResult');
      resultDiv.innerHTML = '<p>Estimating time...</p>';

      try {
        const response = await vscode.postMessage({
          type: 'estimate',
          payload: { title, description, tags: tags.split(',').map(t => t.trim()).filter(t => t) }
        });

        if (response.success) {
          resultDiv.innerHTML = '<pre>' + response.result + '</pre>';
        } else {
          resultDiv.innerHTML = '<p class="error">Error: ' + response.error + '</p>';
        }
      } catch (error) {
        resultDiv.innerHTML = '<p class="error">Error: ' + error + '</p>';
      }
    }

    async function semanticSearch() {
      const query = (document.getElementById('searchQuery') as HTMLInputElement).value;
      const limit = (document.getElementById('searchLimit') as HTMLInputElement).value;
      const resultDiv = document.getElementById('searchResult');
      resultDiv.innerHTML = '<p>Searching...</p>';

      try {
        const response = await vscode.postMessage({
          type: 'search',
          payload: { query, limit: parseInt(limit) }
        });

        if (response.success) {
          resultDiv.innerHTML = '<pre>' + response.result + '</pre>';
        } else {
          resultDiv.innerHTML = '<p class="error">Error: ' + response.error + '</p>';
        }
      } catch (error) {
        resultDiv.innerHTML = '<p class="error">Error: ' + error + '</p>';
      }
    }

    async function configureAI() {
      const provider = (document.getElementById('configProvider') as HTMLSelectElement).value;
      const apiKey = (document.getElementById('configApiKey') as HTMLInputElement).value;
      const model = (document.getElementById('configModel') as HTMLInputElement).value;
      const resultDiv = document.getElementById('configureResult');
      resultDiv.innerHTML = '<p>Configuring AI...</p>';

      try {
        const response = await vscode.postMessage({
          type: 'configure',
          payload: { provider, apiKey, model }
        });

        if (response.success) {
          resultDiv.innerHTML = '<p>✅ AI configured successfully!</p>';
        } else {
          resultDiv.innerHTML = '<p class="error">Error: ' + response.error + '</p>';
        }
      } catch (error) {
        resultDiv.innerHTML = '<p class="error">Error: ' + error + '</p>';
      }
    }

    // Handle Enter key in chat input
    document.getElementById('chatInput')?.addEventListener('keypress', (e) => {
      if (e.key === 'Enter') {
        sendChat();
      }
    });
  </script>
</body>
</html>`;
  }

  /**
   * Generate a nonce for CSP
   */
  private getNonce(): string {
    let text = '';
    const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    for (let i = 0; i < 32; i++) {
      text += possible.charAt(Math.floor(Math.random() * possible.length));
    }
    return text;
  }

  /**
   * Dispose resources
   */
  dispose(): void {
    this.panel?.dispose();
    while (this.disposables.length) {
      const disposable = this.disposables.pop();
      if (disposable) {
        disposable.dispose();
      }
    }
  }
}
