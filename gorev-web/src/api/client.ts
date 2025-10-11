import axios from 'axios';
import type {
  ApiResponse,
  Task,
  Project,
  Template,
  CreateTaskFromTemplateRequest,
  CreateProjectRequest,
  UpdateTaskRequest,
  TaskFilter,
  WorkspaceContext,
  WorkspaceListResponse,
  WorkspaceInfo
} from '@/types';

// Create axios instance with base configuration
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Workspace context management
let currentWorkspaceContext: WorkspaceContext | null = null;

export const setWorkspaceContext = (context: WorkspaceContext | null) => {
  currentWorkspaceContext = context;
  console.log('📁 Workspace context set:', context);
};

export const getWorkspaceContext = (): WorkspaceContext | null => {
  return currentWorkspaceContext;
};

// Add request interceptor for workspace headers and debugging
api.interceptors.request.use(
  (config) => {
    // Inject workspace headers if context is set
    if (currentWorkspaceContext) {
      config.headers['X-Workspace-Id'] = currentWorkspaceContext.workspaceId;
      config.headers['X-Workspace-Path'] = currentWorkspaceContext.workspacePath;
      config.headers['X-Workspace-Name'] = currentWorkspaceContext.workspaceName;
    }

    console.log(`🔥 API Request: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error('❌ API Request Error:', error);
    return Promise.reject(error);
  }
);

// Add response interceptor for error handling
api.interceptors.response.use(
  (response) => {
    console.log(`✅ API Response: ${response.status} ${response.config.url}`);
    return response;
  },
  (error) => {
    console.error('❌ API Response Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

// Health Check
export const checkHealth = async () => {
  const { data } = await api.get('/health');
  return data;
};

// Tasks API
export const getTasks = async (filters?: TaskFilter): Promise<ApiResponse<Task[]>> => {
  const params = new URLSearchParams();

  if (filters?.status) params.append('status', filters.status);
  if (filters?.priority) params.append('priority', filters.priority);
  if (filters?.proje_id) params.append('proje_id', filters.proje_id);
  if (filters?.tag) params.append('tag', filters.tag);

  const { data } = await api.get(`/tasks?${params.toString()}`);
  return data;
};

export const getTask = async (id: string): Promise<ApiResponse<Task>> => {
  const { data } = await api.get(`/tasks/${id}`);
  return data;
};

export const createTaskFromTemplate = async (
  request: CreateTaskFromTemplateRequest
): Promise<ApiResponse<Task>> => {
  const { data } = await api.post('/tasks/from-template', request);
  return data;
};

export const updateTask = async (
  id: string,
  updates: UpdateTaskRequest
): Promise<ApiResponse<Task>> => {
  const { data } = await api.put(`/tasks/${id}`, updates);
  return data;
};

export const deleteTask = async (id: string): Promise<ApiResponse<void>> => {
  const { data } = await api.delete(`/tasks/${id}`);
  return data;
};

// Projects API
export const getProjects = async (): Promise<ApiResponse<Project[]>> => {
  const { data } = await api.get('/projects');
  return data;
};

export const getProject = async (id: string): Promise<ApiResponse<Project>> => {
  const { data } = await api.get(`/projects/${id}`);
  return data;
};

export const createProject = async (
  request: CreateProjectRequest
): Promise<ApiResponse<Project>> => {
  const { data } = await api.post('/projects', request);
  return data;
};

export const getProjectTasks = async (
  projectId: string,
  filters?: TaskFilter
): Promise<ApiResponse<Task[]>> => {
  const params = new URLSearchParams();

  if (filters?.status) params.append('status', filters.status);
  if (filters?.priority) params.append('priority', filters.priority);

  const { data } = await api.get(`/projects/${projectId}/tasks?${params.toString()}`);
  return data;
};

export const activateProject = async (id: string): Promise<ApiResponse<Project>> => {
  const { data } = await api.put(`/projects/${id}/activate`);
  return data;
};

// Templates API
export const getTemplates = async (category?: string): Promise<ApiResponse<Template[]>> => {
  const params = category ? `?category=${encodeURIComponent(category)}` : '';
  const { data } = await api.get(`/templates${params}`);
  return data;
};

// Summary API
export const getSummary = async () => {
  const { data } = await api.get('/summary');
  return data;
};

// Workspace API
export const getWorkspaces = async (): Promise<WorkspaceListResponse> => {
  const { data } = await api.get('/workspaces');
  return data;
};

export const getWorkspace = async (id: string): Promise<ApiResponse<WorkspaceInfo>> => {
  const { data } = await api.get(`/workspaces/${id}`);
  return data;
};

export default api;