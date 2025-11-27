import axios from 'axios';
import { EventEmitter } from 'events';
import { Logger } from '../utils/logger';
import {
  WorkspaceContext,
  WorkspaceInfo,
  WorkspaceRegistration,
  WorkspaceRegistrationResponse,
  WorkspaceListResponse
} from '../models/workspace';

export interface ApiResponse<T = unknown> {
  data: T;
  success: boolean;
  message?: string;
  total?: number;
}

export interface MCPToolResult {
  content: {
    type: string;
    text: string;
  }[];
}

export interface Task {
  id: string;
  baslik: string;
  aciklama: string;
  durum: 'beklemede' | 'devam_ediyor' | 'tamamlandi';
  oncelik: 'dusuk' | 'orta' | 'yuksek';
  olusturma_tarihi: string;
  guncelleme_tarihi: string;
  son_tarih?: string;
  proje_id?: string;
  proje_name?: string;
  etiketler?: { id: string; isim: string }[];
  // Hierarchy fields
  parent_id?: string;
  alt_gorevler?: Task[];
  seviye?: number;
  // Dependency count fields
  bagimli_gorev_sayisi?: number;
  tamamlanmamis_bagimlilik_sayisi?: number;
  bu_goreve_bagimli_sayisi?: number;
}

// API response with English field names (v0.17.0+)
interface ApiTask {
  id: string;
  title: string;
  description: string;
  status: string;
  priority: string;
  created_at: string;
  updated_at: string;
  due_date?: string;
  project_id?: string;
  proje_id?: string; // Some endpoints still use this
  proje_name?: string;
  parent_id?: string;
  workspace_id?: string;
  tags?: { id: string; name: string }[];
  dependency_count?: number;
  uncompleted_dependency_count?: number;
  dependent_on_this_count?: number;
}

// Map API task (English fields) to frontend Task (Turkish fields)
function mapApiTaskToTask(apiTask: ApiTask): Task {
  return {
    id: apiTask.id,
    baslik: apiTask.title,
    aciklama: apiTask.description,
    durum: apiTask.status as Task['durum'],
    oncelik: apiTask.priority as Task['oncelik'],
    olusturma_tarihi: apiTask.created_at,
    guncelleme_tarihi: apiTask.updated_at,
    son_tarih: apiTask.due_date,
    proje_id: apiTask.project_id || apiTask.proje_id,
    proje_name: apiTask.proje_name,
    parent_id: apiTask.parent_id,
    etiketler: apiTask.tags?.map(t => ({ id: t.id, isim: t.name })),
    bagimli_gorev_sayisi: apiTask.dependency_count,
    tamamlanmamis_bagimlilik_sayisi: apiTask.uncompleted_dependency_count,
    bu_goreve_bagimli_sayisi: apiTask.dependent_on_this_count,
  };
}

export interface Project {
  id: string;
  name: string;
  definition?: string;
  created_at: string;
  updated_at?: string;
  is_active?: boolean;
  task_count: number;
}

export interface Template {
  id: string;
  name: string;
  definition: string;
  alias?: string;
  default_title?: string;
  description_template?: string;
  category: string;
  fields: TemplateField[];
  sample_values?: Record<string, string> | null;
  active: boolean;
  language_code?: string;
  base_template_id?: string;
}

export interface TemplateField {
  name: string;
  type: 'text' | 'select' | 'date';
  required: boolean;
  default?: string;
  options?: string[] | null;
  description?: string;
}

export interface CreateTaskFromTemplateRequest {
  template_id: string;
  degerler: Record<string, string>;
}

export interface TaskFilter {
  durum?: string;
  oncelik?: string;
  proje_id?: string;
  etiket?: string;
  limit?: number;
  offset?: number;
  sirala?: string;
  filtre?: string;
  tum_projeler?: boolean;
}

export interface SubtaskData {
  baslik: string;
  aciklama?: string;
  oncelik?: string;
  son_tarih?: string;
  etiketler?: string;
}

export interface TaskHierarchy {
  gorev: Task;
  alt_gorevler: Task[];
  toplam_alt_gorev: number;
  tamamlanan_alt_gorev: number;
}

export interface DependencyRequest {
  kaynak_id: string;
  baglanti_tipi?: string;
}

export interface ExportRequest {
  output_path: string;
  format?: string;
  include_completed?: boolean;
  include_dependencies?: boolean;
  include_templates?: boolean;
  include_ai_context?: boolean;
  project_filter?: string[];
}

export interface ImportRequest {
  file_path: string;
  import_mode?: string;
  conflict_resolution?: string;
  dry_run?: boolean;
  preserve_ids?: boolean;
  project_mapping?: Record<string, string>;
}

// Custom error class for API errors
export class ApiError extends Error {
  constructor(
    public statusCode: number,
    public apiError: string,
    public endpoint: string
  ) {
    super(`API Error ${statusCode} at ${endpoint}: ${apiError}`);
    this.name = 'ApiError';

    // Maintain proper stack trace for where our error was thrown (only available on V8)
    if (Error.captureStackTrace) {
      Error.captureStackTrace(this, ApiError);
    }
  }

  isNotFound(): boolean {
    return this.statusCode === 404;
  }

  isBadRequest(): boolean {
    return this.statusCode === 400;
  }

  isServerError(): boolean {
    return this.statusCode >= 500;
  }
}

export class ApiClient extends EventEmitter {
  private axiosInstance: ReturnType<typeof axios.create>;
  private connected = false;
  private baseURL: string;
  private workspaceContext: WorkspaceContext | undefined;

  constructor(baseURL = 'http://localhost:5082') {
    super();
    this.baseURL = baseURL;

    this.axiosInstance = axios.create({
      baseURL: `${baseURL}/api/v1`,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.setupInterceptors();
  }

  isConnected(): boolean {
    return this.connected;
  }

  async connect(): Promise<void> {
    try {
      const response = await this.axiosInstance.get('/health');
      const healthData = response.data as { status: string };
      this.connected = healthData && healthData.status === 'ok';
      if (this.connected) {
        this.emit('connected');
        Logger.info('[ApiClient] Connected to API server');
      }
    } catch (error) {
      this.connected = false;
      Logger.error('[ApiClient] Failed to connect:', error);
      throw error;
    }
  }

  disconnect(): void {
    this.connected = false;
    this.emit('disconnected');
    Logger.info('[ApiClient] Disconnected from API server');
  }

  private setupInterceptors(): void {
    // Request interceptor for logging and workspace header injection
    this.axiosInstance.interceptors.request.use(
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (config: any) => {
        // Inject workspace headers if context is set
        if (this.workspaceContext) {
          config.headers['X-Workspace-Id'] = this.workspaceContext.workspaceId;
          config.headers['X-Workspace-Path'] = this.workspaceContext.workspacePath;
          config.headers['X-Workspace-Name'] = this.workspaceContext.workspaceName;
        }

        Logger.debug(`[ApiClient] Request: ${config.method?.toUpperCase()} ${config.url}`, config.data);
        return config;
      },
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (error: any) => {
        Logger.error('[ApiClient] Request Error:', error);
        return Promise.reject(error);
      }
    );

    // Response interceptor for logging and error handling
    this.axiosInstance.interceptors.response.use(
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (response: any) => {
        Logger.debug(`[ApiClient] Response: ${response.status} ${response.config.url}`, response.data);
        return response;
      },
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      (error: any) => {
        // Convert axios error to ApiError
        if (error.response) {
          const statusCode = error.response.status;
          const endpoint = error.config?.url || 'unknown';
          const errorMessage = (error.response.data as { error?: string; message?: string })?.error ||
                               (error.response.data as { error?: string; message?: string })?.message ||
                               error.message;

          Logger.error(`[ApiClient] API Error ${statusCode} at ${endpoint}:`, errorMessage);

          const apiError = new ApiError(statusCode, errorMessage, endpoint);
          return Promise.reject(apiError);
        }

        // Network error or other error
        Logger.error('[ApiClient] Response Error:', error.message);
        return Promise.reject(error);
      }
    );
  }

  // Health check
  async checkHealth(): Promise<{ status: string }> {
    const response = await this.axiosInstance.get('/health');
    return response.data as { status: string };
  }

  // Tasks API
  async getTasks(filters?: TaskFilter): Promise<ApiResponse<Task[]>> {
    const params = new URLSearchParams();

    if (filters?.durum) params.append('durum', filters.durum);
    if (filters?.oncelik) params.append('oncelik', filters.oncelik);
    if (filters?.proje_id) params.append('proje_id', filters.proje_id);
    if (filters?.etiket) params.append('etiket', filters.etiket);
    if (filters?.limit !== undefined) params.append('limit', filters.limit.toString());
    if (filters?.offset !== undefined) params.append('offset', filters.offset.toString());
    if (filters?.sirala) params.append('sirala', filters.sirala);
    if (filters?.filtre) params.append('filtre', filters.filtre);
    if (filters?.tum_projeler) params.append('tum_projeler', 'true');

    const response = await this.axiosInstance.get(`/tasks?${params.toString()}`);
    const apiResponse = response.data as ApiResponse<ApiTask[]>;

    // Map API response (English fields) to frontend format (Turkish fields)
    return {
      ...apiResponse,
      data: apiResponse.data?.map(mapApiTaskToTask) || []
    };
  }

  async getTask(id: string): Promise<ApiResponse<Task>> {
    const response = await this.axiosInstance.get(`/tasks/${id}`);
    const apiResponse = response.data as ApiResponse<ApiTask>;

    // Map API response (English fields) to frontend format (Turkish fields)
    return {
      ...apiResponse,
      data: apiResponse.data ? mapApiTaskToTask(apiResponse.data) : apiResponse.data as unknown as Task
    };
  }

  async createTaskFromTemplate(request: CreateTaskFromTemplateRequest): Promise<ApiResponse<Task>> {
    const response = await this.axiosInstance.post('/tasks/from-template', request);
    const apiResponse = response.data as ApiResponse<ApiTask>;

    // Map API response (English fields) to frontend format (Turkish fields)
    return {
      ...apiResponse,
      data: apiResponse.data ? mapApiTaskToTask(apiResponse.data) : apiResponse.data as unknown as Task
    };
  }

  async post<T>(endpoint: string, data?: unknown): Promise<ApiResponse<T>> {
    const response = await this.axiosInstance.post(endpoint, data);
    return response.data as ApiResponse<T>;
  }

  async updateTask(id: string, updates: Partial<Task>): Promise<ApiResponse<Task>> {
    const response = await this.axiosInstance.put(`/tasks/${id}`, updates);
    const apiResponse = response.data as ApiResponse<ApiTask>;

    // Map API response (English fields) to frontend format (Turkish fields)
    return {
      ...apiResponse,
      data: apiResponse.data ? mapApiTaskToTask(apiResponse.data) : apiResponse.data as unknown as Task
    };
  }

  async deleteTask(id: string): Promise<ApiResponse<void>> {
    const response = await this.axiosInstance.delete(`/tasks/${id}`);
    return response.data as ApiResponse<void>;
  }

  // Projects API
  async getProjects(): Promise<ApiResponse<Project[]>> {
    const response = await this.axiosInstance.get('/projects');
    return response.data as ApiResponse<Project[]>;
  }

  async getProject(id: string): Promise<ApiResponse<Project>> {
    const response = await this.axiosInstance.get(`/projects/${id}`);
    return response.data as ApiResponse<Project>;
  }

  async createProject(project: { isim: string; tanim?: string }): Promise<ApiResponse<Project>> {
    const response = await this.axiosInstance.post('/projects', project);
    return response.data as ApiResponse<Project>;
  }

  async getProjectTasks(projectId: string, filters?: TaskFilter): Promise<ApiResponse<Task[]>> {
    const params = new URLSearchParams();

    if (filters?.durum) params.append('durum', filters.durum);
    if (filters?.oncelik) params.append('oncelik', filters.oncelik);

    const response = await this.axiosInstance.get(`/projects/${projectId}/tasks?${params.toString()}`);
    const apiResponse = response.data as ApiResponse<ApiTask[]>;

    // Map API response (English fields) to frontend format (Turkish fields)
    return {
      ...apiResponse,
      data: apiResponse.data?.map(mapApiTaskToTask) || []
    };
  }

  async activateProject(id: string): Promise<ApiResponse<Project>> {
    const response = await this.axiosInstance.put(`/projects/${id}/activate`);
    return response.data as ApiResponse<Project>;
  }

  // Templates API
  async getTemplates(kategori?: string): Promise<ApiResponse<Template[]>> {
    const params = kategori ? `?kategori=${encodeURIComponent(kategori)}` : '';
    const response = await this.axiosInstance.get(`/templates${params}`);
    return response.data as ApiResponse<Template[]>;
  }

  // Summary API
  async getSummary(): Promise<ApiResponse<{
    projects: number;
    tasks: number;
    templates: number;
  }>> {
    const response = await this.axiosInstance.get('/summary');
    return response.data as ApiResponse<{ projects: number; tasks: number; templates: number; }>;
  }

  // Export/Import API
  async exportData(request: ExportRequest): Promise<ApiResponse<{
    file_path: string;
    exported_count: number;
  }>> {
    const response = await this.axiosInstance.post('/export', request);
    return response.data as ApiResponse<{ file_path: string; exported_count: number; }>;
  }

  async importData(request: ImportRequest): Promise<ApiResponse<{
    imported_count: number;
    skipped_count: number;
    errors: string[];
  }>> {
    const response = await this.axiosInstance.post('/import', request);
    return response.data as ApiResponse<{ imported_count: number; skipped_count: number; errors: string[]; }>;
  }

  // Subtask API
  async createSubtask(parentId: string, data: SubtaskData): Promise<ApiResponse<Task>> {
    // Map Turkish field names to English for API (v0.17.0+ migration)
    const apiData = {
      title: data.baslik,
      description: data.aciklama,
      priority: data.oncelik,
      due_date: data.son_tarih,
      etiketler: data.etiketler // Backend still uses Turkish for tags
    };
    const response = await this.axiosInstance.post(`/tasks/${parentId}/subtasks`, apiData);
    const apiResponse = response.data as ApiResponse<ApiTask>;

    // Map API response (English fields) to frontend format (Turkish fields)
    return {
      ...apiResponse,
      data: apiResponse.data ? mapApiTaskToTask(apiResponse.data) : apiResponse.data as unknown as Task
    };
  }

  async changeParent(taskId: string, newParentId: string): Promise<ApiResponse<Task>> {
    const response = await this.axiosInstance.put(`/tasks/${taskId}/parent`, {
      new_parent_id: newParentId
    });
    const apiResponse = response.data as ApiResponse<ApiTask>;

    // Map API response (English fields) to frontend format (Turkish fields)
    return {
      ...apiResponse,
      data: apiResponse.data ? mapApiTaskToTask(apiResponse.data) : apiResponse.data as unknown as Task
    };
  }

  async getHierarchy(taskId: string): Promise<ApiResponse<TaskHierarchy>> {
    const response = await this.axiosInstance.get(`/tasks/${taskId}/hierarchy`);
    return response.data as ApiResponse<TaskHierarchy>;
  }

  // Dependency API
  async addDependency(targetId: string, dependency: DependencyRequest): Promise<ApiResponse<void>> {
    const response = await this.axiosInstance.post(`/tasks/${targetId}/dependencies`, dependency);
    return response.data as ApiResponse<void>;
  }

  async removeDependency(targetId: string, sourceId: string): Promise<ApiResponse<void>> {
    const response = await this.axiosInstance.delete(`/tasks/${targetId}/dependencies/${sourceId}`);
    return response.data as ApiResponse<void>;
  }

  // Active Project API
  async getActiveProject(): Promise<ApiResponse<Project | null>> {
    const response = await this.axiosInstance.get('/active-project');
    return response.data as ApiResponse<Project | null>;
  }

  async removeActiveProject(): Promise<ApiResponse<void>> {
    const response = await this.axiosInstance.delete('/active-project');
    return response.data as ApiResponse<void>;
  }

  // Language API
  async getLanguage(): Promise<{ success: boolean; language: string }> {
    const response = await this.axiosInstance.get('/language');
    return response.data as { success: boolean; language: string };
  }

  async setLanguage(language: 'tr' | 'en'): Promise<{ success: boolean; language: string; message: string }> {
    const response = await this.axiosInstance.post('/language', { language });
    return response.data as { success: boolean; language: string; message: string };
  }

  // Convert API responses to MCP-like format for compatibility
  async callTool(name: string, params?: Record<string, unknown>): Promise<MCPToolResult> {
    Logger.info(`[ApiClient] Calling tool: ${name} with params:`, params);

    try {
      let result: ApiResponse<unknown>;

      switch (name) {
        case 'gorev_listele':
          result = await this.getTasks(params as TaskFilter | undefined);
          return this.convertToMCPFormat(result);

        case 'proje_listele':
          result = await this.getProjects();
          return this.convertToMCPFormat(result);

        case 'template_listele':
          result = await this.getTemplates(params?.kategori as string | undefined);
          return this.convertToMCPFormat(result);

        case 'templateden_gorev_olustur':
          result = await this.createTaskFromTemplate(params as unknown as CreateTaskFromTemplateRequest);
          return this.convertToMCPFormat(result);

        case 'gorev_guncelle':
          result = await this.updateTask(
            (params as { id: string }).id,
            { durum: (params as { durum: 'beklemede' | 'devam_ediyor' | 'tamamlandi' }).durum }
          );
          return this.convertToMCPFormat(result);

        case 'gorev_duzenle': {
          // General task update - can update proje_id, durum, oncelik, etc.
          const { id, ...updateFields } = params as { id: string; proje_id?: string; durum?: string; oncelik?: string };
          result = await this.updateTask(id, updateFields as Partial<Task>);
          return this.convertToMCPFormat(result);
        }

        case 'gorev_sil':
          result = await this.deleteTask((params as { id: string }).id);
          return this.convertToMCPFormat(result);

        case 'gorev_ust_gorev_degistir': {
          // Change parent task
          const changeParentParams = params as { gorev_id: string; yeni_ust_gorev_id: string };
          result = await this.changeParent(changeParentParams.gorev_id, changeParentParams.yeni_ust_gorev_id);
          return this.convertToMCPFormat(result);
        }

        case 'gorev_bagimlilik_ekle': {
          // Add dependency - hedef_id is the task that will have the dependency, kaynak_id is the source
          const depParams = params as { kaynak_id: string; hedef_id: string; baglanti_tipi?: string };
          result = await this.addDependency(depParams.hedef_id, {
            kaynak_id: depParams.kaynak_id,
            baglanti_tipi: depParams.baglanti_tipi
          });
          return this.convertToMCPFormat(result);
        }

        case 'alt_gorev_olustur': {
          // Create subtask (legacy parameter names)
          const subtaskParams = params as { ust_gorev_id: string; baslik: string; aciklama?: string; oncelik?: string };
          result = await this.createSubtask(subtaskParams.ust_gorev_id, {
            baslik: subtaskParams.baslik,
            aciklama: subtaskParams.aciklama,
            oncelik: subtaskParams.oncelik
          });
          return this.convertToMCPFormat(result);
        }

        case 'gorev_altgorev_olustur': {
          // Create subtask (current MCP tool name)
          const subtaskParams = params as {
            parent_id: string;
            baslik: string;
            aciklama?: string;
            oncelik?: string;
            son_tarih?: string;
            tags?: string;
          };
          result = await this.createSubtask(subtaskParams.parent_id, {
            baslik: subtaskParams.baslik,
            aciklama: subtaskParams.aciklama,
            oncelik: subtaskParams.oncelik,
            son_tarih: subtaskParams.son_tarih,
            etiketler: subtaskParams.tags
          });
          return this.convertToMCPFormat(result);
        }

        case 'proje_olustur':
          result = await this.createProject(params as { isim: string; tanim?: string });
          return this.convertToMCPFormat(result);

        case 'aktif_proje_ayarla':
          result = await this.activateProject((params as { proje_id: string }).proje_id);
          return this.convertToMCPFormat(result);

        case 'ozet_goster':
          result = await this.getSummary();
          return this.convertToMCPFormat(result);

        case 'gorev_export':
          result = await this.exportData(params as unknown as ExportRequest);
          return this.convertToMCPFormat(result);

        case 'gorev_import':
          result = await this.importData(params as unknown as ImportRequest);
          return this.convertToMCPFormat(result);

        // MCP-only tools (no REST API endpoints)
        case 'gorev_set_active':
        case 'gorev_get_active':
        case 'gorev_nlp_query':
        case 'gorev_context_summary':
        case 'gorev_batch_update':
          Logger.warn(`[ApiClient] Tool '${name}' is MCP-only and not available via REST API`);
          return {
            content: [{
              type: 'text',
              text: `Tool '${name}' is only available via MCP protocol, not REST API. Skipping.`
            }]
          };

        default:
          throw new Error(`Unsupported tool: ${name}`);
      }
    } catch (error) {
      Logger.error(`[ApiClient] Tool call failed for ${name}:`, error);
      throw error;
    }
  }

  private convertToMCPFormat(apiResponse: ApiResponse<unknown>): MCPToolResult {
    // Convert API response to MCP tool result format
    if (apiResponse.success) {
      return {
        content: [
          {
            type: 'text',
            text: this.formatDataAsText(apiResponse.data)
          }
        ]
      };
    } else {
      throw new Error(apiResponse.message || 'API call failed');
    }
  }

  private formatDataAsText(data: unknown): string {
    if (Array.isArray(data)) {
      if (data.length === 0) {
        return 'No data found.';
      }

      // Format different types of arrays
      if (data[0]?.baslik) {
        // Tasks
        return this.formatTasks(data);
      } else if (data[0]?.name && data[0]?.task_count !== undefined) {
        // Projects (use English field names from API)
        return this.formatProjects(data);
      } else if (data[0]?.fields) {
        // Templates (use English field name from API)
        return this.formatTemplates(data);
      }
    }

    return JSON.stringify(data, null, 2);
  }

  // Workspace Management API

  /**
   * Set workspace context for header injection
   */
  setWorkspaceHeaders(context: WorkspaceContext): void {
    this.workspaceContext = context;
    Logger.info(`[ApiClient] Workspace context set: ${context.workspaceName} (${context.workspaceId})`);
  }

  /**
   * Clear workspace context
   */
  clearWorkspaceHeaders(): void {
    this.workspaceContext = undefined;
    Logger.info('[ApiClient] Workspace context cleared');
  }

  /**
   * Get current workspace context
   */
  getWorkspaceContext(): WorkspaceContext | undefined {
    return this.workspaceContext;
  }

  /**
   * Register a workspace with the server
   */
  async registerWorkspace(registration: WorkspaceRegistration): Promise<WorkspaceRegistrationResponse> {
    const response = await this.axiosInstance.post('/workspaces/register', registration);
    return response.data as WorkspaceRegistrationResponse;
  }

  /**
   * List all registered workspaces
   */
  async listWorkspaces(): Promise<WorkspaceListResponse> {
    const response = await this.axiosInstance.get('/workspaces');
    return response.data as WorkspaceListResponse;
  }

  /**
   * Get workspace details by ID
   */
  async getWorkspace(workspaceId: string): Promise<ApiResponse<WorkspaceInfo>> {
    const response = await this.axiosInstance.get(`/workspaces/${workspaceId}`);
    return response.data as ApiResponse<WorkspaceInfo>;
  }

  /**
   * Unregister a workspace
   */
  async unregisterWorkspace(workspaceId: string): Promise<ApiResponse<void>> {
    const response = await this.axiosInstance.delete(`/workspaces/${workspaceId}`);
    return response.data as ApiResponse<void>;
  }

  // Formatting helpers

  private formatTasks(tasks: Task[]): string {
    const grouped = tasks.reduce((acc, task) => {
      if (!acc[task.durum]) acc[task.durum] = [];
      acc[task.durum].push(task);
      return acc;
    }, {} as Record<string, Task[]>);

    const statusLabels = {
      'beklemede': '⏳ Beklemede',
      'devam_ediyor': '🔄 Devam Ediyor',
      'tamamlandi': '✅ Tamamlandı'
    };

    let output = '## Görev Listesi\n\n';

    for (const [status, tasks] of Object.entries(grouped)) {
      const label = statusLabels[status as keyof typeof statusLabels] || status;
      output += `### ${label}\n\n`;

      for (const task of tasks) {
        const priority = task.oncelik === 'yuksek' ? '🔴' : task.oncelik === 'orta' ? '🟡' : '🟢';
        output += `- **${task.baslik}** ${priority}\n`;
        if (task.aciklama) {
          output += `  ${task.aciklama}\n`;
        }
        if (task.proje_name) {
          output += `  📁 ${task.proje_name}\n`;
        }
        output += '\n';
      }
    }

    return output;
  }

  private formatProjects(projects: Project[]): string {
    let output = '## Proje Listesi\n\n';

    for (const project of projects) {
      output += `### ${project.name}\n`;
      if (project.definition) {
        output += `${project.definition}\n`;
      }
      output += `📊 ${project.task_count} görev\n`;
      if (project.is_active) {
        output += '🟢 Aktif proje\n';
      }
      output += '\n';
    }

    return output;
  }

  private formatTemplates(templates: Template[]): string {
    let output = '## 📋 Görev Template\'leri\n\n';

    const grouped = templates.reduce((acc, template) => {
      if (!acc[template.category]) acc[template.category] = [];
      acc[template.category].push(template);
      return acc;
    }, {} as Record<string, Template[]>);

    for (const [category, categoryTemplates] of Object.entries(grouped)) {
      output += `### ${category}\n\n`;

      for (const template of categoryTemplates) {
        output += `#### ${template.name}\n`;
        output += `- **ID:** \`${template.id}\`\n`;
        output += `- **Açıklama:** ${template.definition}\n`;
        if (template.alias) {
          output += `- **Alias:** ${template.alias}\n`;
        }
        output += '\n';
      }
    }

    output += '\n💡 **Kullanım:** `templateden_gorev_olustur` komutunu template ID\'si ve alan değerleriyle kullanın.';

    return output;
  }
}