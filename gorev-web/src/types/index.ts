// API Response Types
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  total?: number;
  message?: string;
}

// Task Types
export interface Task {
  id: string;
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  proje_id?: string;
  proje_name?: string;
  parent_id?: string;
  due_date?: string;
  tags?: Array<{ id: string; name: string }>;
  created_at: string;
  updated_at: string;
  // Subtask and dependency info
  subtasks?: Task[];
  has_subtasks?: boolean;
  subtask_count?: number;
  dependency_count?: number;
  uncompleted_dependency_count?: number;
}

export type TaskStatus = 'beklemede' | 'devam_ediyor' | 'tamamlandi';
export type TaskPriority = 'dusuk' | 'orta' | 'yuksek';

// Project Types
export interface Project {
  id: string;
  name: string;
  definition: string;
  created_at: string;
  task_count: number;
  is_active: boolean;
}

// Template Types
export interface Template {
  id: string;
  name: string;
  definition: string;
  alias?: string;
  default_title?: string;
  description_template?: string;
  sample_values?: Record<string, string> | null;
  fields: TemplateField[];
  category: string;
  active: boolean;
}

export interface TemplateField {
  name: string;
  type: 'text' | 'select' | 'date';
  required: boolean;
  default?: string;
  options?: string[];
  description?: string;
}

// Form Types for API requests
export interface CreateTaskFromTemplateRequest {
  template_id: string;
  proje_id: string;
  values: Record<string, string>;
}

export interface CreateProjectRequest {
  name: string;
  definition?: string;
}

export interface UpdateTaskRequest {
  title?: string;
  description?: string;
  status?: TaskStatus;
  priority?: TaskPriority;
  proje_id?: string;
  due_date?: string;
  tags?: string[];
}

// UI State Types
export interface TaskFilter {
  status?: TaskStatus;
  priority?: TaskPriority;
  proje_id?: string;
  tag?: string;
  search?: string;
}

export interface AppState {
  selectedProject?: Project;
  taskFilter: TaskFilter;
  sidebarOpen: boolean;
}

// Workspace Types
export interface WorkspaceInfo {
  id: string;
  name: string;
  path: string;
  database_path: string;
  last_accessed: string;
  created_at: string;
  task_count: number;
}

export interface WorkspaceContext {
  workspaceId: string;
  workspaceName: string;
  workspacePath: string;
}

export interface WorkspaceListResponse {
  success: boolean;
  workspaces: WorkspaceInfo[];
  total: number;
}