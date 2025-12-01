/**
 * API Data Seeder for E2E Tests
 *
 * Seeds test data using the Gorev REST API.
 * Provides consistent test data across all E2E tests.
 */

import { APIRequestContext } from '@playwright/test';

export interface Project {
  id: string;
  name: string;
  description: string;
  workspace_id: string;
  created_at: string;
  updated_at: string;
}

export interface Task {
  id: string;
  title: string;
  description: string;
  status: 'beklemede' | 'devam_ediyor' | 'tamamlandi' | 'iptal';
  priority: 'dusuk' | 'orta' | 'yuksek';
  project_id?: string;
  parent_id?: string;
  due_date?: string;
  workspace_id: string;
  created_at: string;
  updated_at: string;
}

export interface SeedConfig {
  language: 'tr' | 'en';
  minimal: boolean;
  includeSubtasks: boolean;
  includeDependencies: boolean;
}

export interface SeedResult {
  projects: Project[];
  tasks: Task[];
  subtasks: Task[];
}

const defaultSeedConfig: SeedConfig = {
  language: 'tr',
  minimal: false,
  includeSubtasks: true,
  includeDependencies: true,
};

// Sample data matching the Go seeder fixtures
const sampleProjects = {
  tr: [
    { name: 'Mobil Uygulama', description: 'iOS ve Android uygulama geliştirme projesi' },
    { name: 'Backend API', description: 'REST API geliştirme ve bakım projesi' },
    { name: 'Web Dashboard', description: 'Admin paneli ve kullanıcı arayüzü projesi' },
  ],
  en: [
    { name: 'Mobile App', description: 'iOS and Android app development project' },
    { name: 'Backend API', description: 'REST API development and maintenance project' },
    { name: 'Web Dashboard', description: 'Admin panel and user interface project' },
  ],
};

const sampleTasks = {
  tr: [
    // Mobil Uygulama tasks
    { title: 'Login sayfası 404 hatası', description: 'Production ortamında login sayfasına giderken 404 hatası alınıyor', status: 'beklemede' as const, priority: 'yuksek' as const, projectIndex: 0, templateAlias: 'bug' },
    { title: 'Push notification sistemi', description: 'Firebase Cloud Messaging ile push notification implementasyonu', status: 'devam_ediyor' as const, priority: 'orta' as const, projectIndex: 0, templateAlias: 'feature' },
    { title: 'Dark mode tema', description: 'Karanlık mod tema desteği eklenmesi', status: 'beklemede' as const, priority: 'dusuk' as const, projectIndex: 0, templateAlias: 'feature' },
    { title: 'Bellek sızıntısı', description: 'Uzun süreli kullanımda bellek sızıntısı tespit edildi', status: 'tamamlandi' as const, priority: 'yuksek' as const, projectIndex: 0, templateAlias: 'bug' },
    { title: 'Offline mod desteği', description: 'İnternet bağlantısı olmadan çalışma özelliği', status: 'iptal' as const, priority: 'orta' as const, projectIndex: 0, templateAlias: 'feature' },
    // Backend API tasks
    { title: 'Redis cache entegrasyonu', description: 'API response caching için Redis implementasyonu', status: 'tamamlandi' as const, priority: 'yuksek' as const, projectIndex: 1, templateAlias: 'feature' },
    { title: 'API rate limiting', description: 'DDoS koruması için rate limiting implementasyonu', status: 'beklemede' as const, priority: 'orta' as const, projectIndex: 1, templateAlias: 'feature' },
    { title: 'Timeout hatası düzeltme', description: 'Büyük veri setlerinde timeout hatası alınıyor', status: 'devam_ediyor' as const, priority: 'orta' as const, projectIndex: 1, templateAlias: 'bug' },
    { title: 'GraphQL API desteği', description: 'REST API yanında GraphQL endpoint eklenmesi', status: 'beklemede' as const, priority: 'dusuk' as const, projectIndex: 1, templateAlias: 'feature' },
    { title: 'API dokümantasyonu güncelleme', description: 'Swagger/OpenAPI dokümantasyonunun güncellenmesi', status: 'tamamlandi' as const, priority: 'dusuk' as const, projectIndex: 1, templateAlias: 'feature' },
    // Web Dashboard tasks
    { title: 'Dashboard ana sayfa tasarımı', description: 'Modern ve kullanıcı dostu dashboard tasarımı', status: 'devam_ediyor' as const, priority: 'orta' as const, projectIndex: 2, templateAlias: 'feature' },
    { title: 'Kullanıcı yönetimi modülü', description: 'Admin panelinde kullanıcı CRUD işlemleri', status: 'devam_ediyor' as const, priority: 'yuksek' as const, projectIndex: 2, templateAlias: 'feature', isParent: true },
    { title: 'Tablo sıralama hatası', description: 'Tablo sütunlarında sıralama düzgün çalışmıyor', status: 'tamamlandi' as const, priority: 'orta' as const, projectIndex: 2, templateAlias: 'bug' },
    { title: 'Raporlama modülü', description: 'Detaylı raporlar ve analitik dashboard', status: 'beklemede' as const, priority: 'orta' as const, projectIndex: 2, templateAlias: 'feature' },
    { title: 'React güncellemesi', description: 'React 17\'den React 18\'e güncelleme', status: 'tamamlandi' as const, priority: 'dusuk' as const, projectIndex: 2, templateAlias: 'feature' },
  ],
  en: [
    // Mobile App tasks
    { title: 'Login page 404 error', description: 'Getting 404 error when navigating to login page in production', status: 'beklemede' as const, priority: 'yuksek' as const, projectIndex: 0, templateAlias: 'bug' },
    { title: 'Push notification system', description: 'Firebase Cloud Messaging push notification implementation', status: 'devam_ediyor' as const, priority: 'orta' as const, projectIndex: 0, templateAlias: 'feature' },
    { title: 'Dark mode theme', description: 'Adding dark mode theme support', status: 'beklemede' as const, priority: 'dusuk' as const, projectIndex: 0, templateAlias: 'feature' },
    { title: 'Memory leak', description: 'Memory leak detected during long-term usage', status: 'tamamlandi' as const, priority: 'yuksek' as const, projectIndex: 0, templateAlias: 'bug' },
    { title: 'Offline mode support', description: 'Ability to work without internet connection', status: 'iptal' as const, priority: 'orta' as const, projectIndex: 0, templateAlias: 'feature' },
    // Backend API tasks
    { title: 'Redis cache integration', description: 'Redis implementation for API response caching', status: 'tamamlandi' as const, priority: 'yuksek' as const, projectIndex: 1, templateAlias: 'feature' },
    { title: 'API rate limiting', description: 'Rate limiting implementation for DDoS protection', status: 'beklemede' as const, priority: 'orta' as const, projectIndex: 1, templateAlias: 'feature' },
    { title: 'Timeout error fix', description: 'Getting timeout error with large datasets', status: 'devam_ediyor' as const, priority: 'orta' as const, projectIndex: 1, templateAlias: 'bug' },
    { title: 'GraphQL API support', description: 'Adding GraphQL endpoint alongside REST API', status: 'beklemede' as const, priority: 'dusuk' as const, projectIndex: 1, templateAlias: 'feature' },
    { title: 'API documentation update', description: 'Updating Swagger/OpenAPI documentation', status: 'tamamlandi' as const, priority: 'dusuk' as const, projectIndex: 1, templateAlias: 'feature' },
    // Web Dashboard tasks
    { title: 'Dashboard homepage design', description: 'Modern and user-friendly dashboard design', status: 'devam_ediyor' as const, priority: 'orta' as const, projectIndex: 2, templateAlias: 'feature' },
    { title: 'User management module', description: 'User CRUD operations in admin panel', status: 'devam_ediyor' as const, priority: 'yuksek' as const, projectIndex: 2, templateAlias: 'feature', isParent: true },
    { title: 'Table sorting bug', description: 'Table column sorting not working properly', status: 'tamamlandi' as const, priority: 'orta' as const, projectIndex: 2, templateAlias: 'bug' },
    { title: 'Reporting module', description: 'Detailed reports and analytics dashboard', status: 'beklemede' as const, priority: 'orta' as const, projectIndex: 2, templateAlias: 'feature' },
    { title: 'React upgrade', description: 'Upgrading from React 17 to React 18', status: 'tamamlandi' as const, priority: 'dusuk' as const, projectIndex: 2, templateAlias: 'feature' },
  ],
};

const minimalTasks = {
  tr: [
    { title: 'Test Bug', description: 'Test bug açıklaması', status: 'beklemede' as const, priority: 'orta' as const, templateAlias: 'bug' },
    { title: 'Test Feature', description: 'Test özellik açıklaması', status: 'devam_ediyor' as const, priority: 'yuksek' as const, templateAlias: 'feature' },
    { title: 'Completed Task', description: 'Tamamlanmış görev', status: 'tamamlandi' as const, priority: 'dusuk' as const, templateAlias: 'feature' },
  ],
  en: [
    { title: 'Test Bug', description: 'Test bug description', status: 'beklemede' as const, priority: 'orta' as const, templateAlias: 'bug' },
    { title: 'Test Feature', description: 'Test feature description', status: 'devam_ediyor' as const, priority: 'yuksek' as const, templateAlias: 'feature' },
    { title: 'Completed Task', description: 'Completed task', status: 'tamamlandi' as const, priority: 'dusuk' as const, templateAlias: 'feature' },
  ],
};

export class DataSeeder {
  private request: APIRequestContext;
  private baseUrl: string;
  private config: SeedConfig;

  constructor(request: APIRequestContext, baseUrl: string, config: Partial<SeedConfig> = {}) {
    this.request = request;
    this.baseUrl = baseUrl;
    this.config = { ...defaultSeedConfig, ...config };
  }

  /**
   * Create a project via API
   */
  async createProject(name: string, description: string): Promise<Project> {
    const response = await this.request.post(`${this.baseUrl}/api/v1/projects`, {
      data: { name, description },
    });

    if (!response.ok()) {
      throw new Error(`Failed to create project: ${await response.text()}`);
    }

    const result = await response.json();
    return result.data;
  }

  /**
   * Create a task from template via API
   */
  async createTaskFromTemplate(
    templateAlias: string,
    values: Record<string, string>,
    projectId?: string,
    status?: string,
    priority?: string
  ): Promise<Task> {
    const response = await this.request.post(`${this.baseUrl}/api/v1/tasks/from-template`, {
      data: {
        template_alias: templateAlias,
        values,
        project_id: projectId,
        status,
        priority,
      },
    });

    if (!response.ok()) {
      throw new Error(`Failed to create task: ${await response.text()}`);
    }

    const result = await response.json();
    return result.data;
  }

  /**
   * Create a subtask via API
   */
  async createSubtask(
    parentId: string,
    title: string,
    description: string,
    status: string,
    priority: string
  ): Promise<Task> {
    const response = await this.request.post(`${this.baseUrl}/api/v1/tasks/${parentId}/subtasks`, {
      data: {
        title,
        description,
        status,
        priority,
      },
    });

    if (!response.ok()) {
      throw new Error(`Failed to create subtask: ${await response.text()}`);
    }

    const result = await response.json();
    return result.data;
  }

  /**
   * Set active project via API
   */
  async setActiveProject(projectId: string): Promise<void> {
    const response = await this.request.post(`${this.baseUrl}/api/v1/projects/${projectId}/set-active`);

    if (!response.ok()) {
      throw new Error(`Failed to set active project: ${await response.text()}`);
    }
  }

  /**
   * Clear all data (for test isolation)
   */
  async clearAllData(): Promise<void> {
    // Get all tasks and delete them
    const tasksResponse = await this.request.get(`${this.baseUrl}/api/v1/tasks`);
    if (tasksResponse.ok()) {
      const tasksResult = await tasksResponse.json();
      for (const task of tasksResult.data || []) {
        await this.request.delete(`${this.baseUrl}/api/v1/tasks/${task.id}`);
      }
    }

    // Get all projects and delete them
    const projectsResponse = await this.request.get(`${this.baseUrl}/api/v1/projects`);
    if (projectsResponse.ok()) {
      const projectsResult = await projectsResponse.json();
      for (const project of projectsResult.data || []) {
        await this.request.delete(`${this.baseUrl}/api/v1/projects/${project.id}`);
      }
    }
  }

  /**
   * Seed projects
   */
  async seedProjects(): Promise<Project[]> {
    const projects: Project[] = [];
    const projectData = this.config.minimal
      ? [{ name: this.config.language === 'tr' ? 'Test Projesi' : 'Test Project', description: this.config.language === 'tr' ? 'Test amaçlı proje' : 'Project for testing purposes' }]
      : sampleProjects[this.config.language];

    for (const data of projectData) {
      const project = await this.createProject(data.name, data.description);
      projects.push(project);
    }

    return projects;
  }

  /**
   * Seed tasks
   */
  async seedTasks(projects: Project[]): Promise<Task[]> {
    const tasks: Task[] = [];
    const taskData = this.config.minimal
      ? minimalTasks[this.config.language].map((t) => ({ ...t, projectIndex: 0 }))
      : sampleTasks[this.config.language];

    for (const data of taskData) {
      const projectId = projects[data.projectIndex]?.id;
      const task = await this.createTaskFromTemplate(
        data.templateAlias,
        {
          title: data.title,
          description: data.description,
        },
        projectId,
        data.status,
        data.priority
      );
      tasks.push(task);
    }

    return tasks;
  }

  /**
   * Seed subtasks for parent tasks
   */
  async seedSubtasks(tasks: Task[]): Promise<Task[]> {
    if (this.config.minimal || !this.config.includeSubtasks) {
      return [];
    }

    const subtasks: Task[] = [];
    const parentTask = tasks.find((t) => t.title.includes('Kullanıcı yönetimi') || t.title.includes('User management'));

    if (parentTask) {
      const subtaskData = this.config.language === 'tr'
        ? [
            { title: 'Kayıt formu', description: 'Yeni kullanıcı kayıt formu implementasyonu', status: 'devam_ediyor', priority: 'orta' },
            { title: 'Giriş sistemi', description: 'Kullanıcı giriş ve oturum yönetimi', status: 'beklemede', priority: 'yuksek' },
          ]
        : [
            { title: 'Registration form', description: 'New user registration form implementation', status: 'devam_ediyor', priority: 'orta' },
            { title: 'Login system', description: 'User login and session management', status: 'beklemede', priority: 'yuksek' },
          ];

      for (const data of subtaskData) {
        const subtask = await this.createSubtask(
          parentTask.id,
          data.title,
          data.description,
          data.status,
          data.priority
        );
        subtasks.push(subtask);
      }
    }

    return subtasks;
  }

  /**
   * Seed all data
   */
  async seedAll(): Promise<SeedResult> {
    console.log(`[DataSeeder] Seeding data (language: ${this.config.language}, minimal: ${this.config.minimal})`);

    const projects = await this.seedProjects();
    console.log(`[DataSeeder] Created ${projects.length} projects`);

    const tasks = await this.seedTasks(projects);
    console.log(`[DataSeeder] Created ${tasks.length} tasks`);

    const subtasks = await this.seedSubtasks(tasks);
    console.log(`[DataSeeder] Created ${subtasks.length} subtasks`);

    // Set first project as active
    if (projects.length > 0) {
      await this.setActiveProject(projects[0].id);
      console.log(`[DataSeeder] Set active project: ${projects[0].name}`);
    }

    return { projects, tasks, subtasks };
  }

  /**
   * Get seeder configuration
   */
  getConfig(): SeedConfig {
    return { ...this.config };
  }
}

export default DataSeeder;
