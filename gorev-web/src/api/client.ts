import axios from 'axios';
import type {
  ApiResponse,
  Task,
  Project,
  Template,
  CreateTaskFromTemplateRequest,
  CreateProjectRequest,
  UpdateTaskRequest,
  TaskFilter
} from '@/types';

// Create axios instance with base configuration
const api = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add request interceptor for debugging
api.interceptors.request.use(
  (config) => {
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

  if (filters?.durum) params.append('durum', filters.durum);
  if (filters?.oncelik) params.append('oncelik', filters.oncelik);
  if (filters?.proje_id) params.append('proje_id', filters.proje_id);
  if (filters?.etiket) params.append('etiket', filters.etiket);

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

  if (filters?.durum) params.append('durum', filters.durum);
  if (filters?.oncelik) params.append('oncelik', filters.oncelik);

  const { data } = await api.get(`/projects/${projectId}/tasks?${params.toString()}`);
  return data;
};

export const activateProject = async (id: string): Promise<ApiResponse<Project>> => {
  const { data } = await api.put(`/projects/${id}/activate`);
  return data;
};

// Templates API
export const getTemplates = async (kategori?: string): Promise<ApiResponse<Template[]>> => {
  const params = kategori ? `?kategori=${encodeURIComponent(kategori)}` : '';
  const { data } = await api.get(`/templates${params}`);
  return data;
};

// Summary API
export const getSummary = async () => {
  const { data } = await api.get('/summary');
  return data;
};

export default api;