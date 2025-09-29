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
  baslik: string;
  aciklama: string;
  durum: TaskStatus;
  oncelik: TaskPriority;
  proje_id?: string;
  proje_name?: string;
  parent_id?: string;
  son_tarih?: string;
  etiketler?: string[];
  olusturma_tarihi: string;
  guncelleme_tarihi: string;
  // Subtask and dependency info
  alt_gorevler?: Task[];
  has_subtasks?: boolean;
  subtask_count?: number;
  bagimli_gorev_sayisi?: number;
  tamamlanmamis_bagimlilik_sayisi?: number;
}

export type TaskStatus = 'beklemede' | 'devam_ediyor' | 'tamamlandi';
export type TaskPriority = 'dusuk' | 'orta' | 'yuksek';

// Project Types
export interface Project {
  id: string;
  isim: string;
  tanim: string;
  olusturma_tarihi: string;
  gorev_sayisi: number;
  is_active: boolean;
}

// Template Types
export interface Template {
  id: string;
  isim: string;
  tanim: string;
  alias?: string;
  varsayilan_baslik?: string;
  aciklama_template?: string;
  ornek_degerler?: Record<string, string> | null;
  alanlar: TemplateField[];
  kategori: string;
  aktif: boolean;
}

export interface TemplateField {
  isim: string;
  tip: 'text' | 'select' | 'date';
  zorunlu: boolean;
  varsayilan?: string;
  secenekler?: string[];
  aciklama?: string;
}

// Form Types for API requests
export interface CreateTaskFromTemplateRequest {
  template_id: string;
  proje_id: string;
  degerler: Record<string, string>;
}

export interface CreateProjectRequest {
  isim: string;
  tanim?: string;
}

export interface UpdateTaskRequest {
  baslik?: string;
  aciklama?: string;
  durum?: TaskStatus;
  oncelik?: TaskPriority;
  proje_id?: string;
  son_tarih?: string;
  etiketler?: string[];
}

// UI State Types
export interface TaskFilter {
  durum?: TaskStatus;
  oncelik?: TaskPriority;
  proje_id?: string;
  etiket?: string;
  search?: string;
}

export interface AppState {
  selectedProject?: Project;
  taskFilter: TaskFilter;
  sidebarOpen: boolean;
}