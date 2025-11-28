/**
 * Mock Server Utilities for UI Testing
 * Provides a lightweight mock API server for testing VS Code extension UI interactions
 */

import express, { Request, Response } from 'express';
import { createServer } from 'http';

export interface MockTask {
  id: string;
  title: string;
  description: string;
  status: 'pending' | 'in_progress' | 'completed';
  priority: 'low' | 'medium' | 'high';
  project_id?: string;
  parent_id?: string;
  due_date?: string;
  tags?: Array<{ id: string; name: string }>;
  created_at: string;
  updated_at: string;
  workspace_id: string;
  dependency_count?: number;
  uncompleted_dependency_count?: number;
  dependent_on_this_count?: number;
}

export interface MockProject {
  id: string;
  name: string;
  description?: string;
  workspace_id: string;
  created_at: string;
  updated_at: string;
  task_count?: number;
}

export class MockServer {
  private app: express.Application;
  private server: any;
  private port: number;
  private tasks: Map<string, MockTask> = new Map();
  private projects: Map<string, MockProject> = new Map();

  constructor(port: number = 0) {
    this.app = express();
    this.port = port;
    this.setupMiddleware();
    this.setupRoutes();
    this.initializeTestData();
  }

  public getPort(): number {
    return this.port;
  }

  private setupMiddleware(): void {
    this.app.use(express.json());
    this.app.use((req, res, next) => {
      // CORS support
      res.header('Access-Control-Allow-Origin', '*');
      res.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
      res.header('Access-Control-Allow-Headers', 'Content-Type, Authorization');
      if (req.method === 'OPTIONS') {
        res.sendStatus(200);
      } else {
        next();
      }
    });
  }

  private setupRoutes(): void {
    // Health check
    this.app.get('/api/v1/health', (req: Request, res: Response) => {
      res.json({ status: 'ok', timestamp: new Date().toISOString() });
    });

    // Tasks API
    this.app.get('/api/v1/tasks', (req: Request, res: Response) => {
      const { project_id, status, limit = 100, offset = 0 } = req.query;
      let tasks = Array.from(this.tasks.values());

      if (project_id) {
        tasks = tasks.filter(t => t.project_id === project_id);
      }
      if (status) {
        tasks = tasks.filter(t => t.status === status);
      }

      const paginated = tasks.slice(Number(offset), Number(offset) + Number(limit));
      res.json({
        success: true,
        data: paginated,
        total: tasks.length,
        limit: Number(limit),
        offset: Number(offset)
      });
    });

    this.app.get('/api/v1/tasks/:id', (req: Request, res: Response) => {
      const task = this.tasks.get(req.params.id);
      if (!task) {
        return res.status(404).json({ success: false, error: 'Task not found' });
      }
      res.json({ success: true, data: task });
    });

    this.app.post('/api/v1/tasks/from-template', (req: Request, res: Response) => {
      const task = this.createTask(req.body);
      this.tasks.set(task.id, task);
      res.json({ success: true, data: task });
    });

    this.app.put('/api/v1/tasks/:id', (req: Request, res: Response) => {
      const existing = this.tasks.get(req.params.id);
      if (!existing) {
        return res.status(404).json({ success: false, error: 'Task not found' });
      }
      const updated = { ...existing, ...req.body, updated_at: new Date().toISOString() };
      this.tasks.set(req.params.id, updated);
      res.json({ success: true, data: updated });
    });

    this.app.delete('/api/v1/tasks/:id', (req: Request, res: Response) => {
      if (!this.tasks.has(req.params.id)) {
        return res.status(404).json({ success: false, error: 'Task not found' });
      }
      this.tasks.delete(req.params.id);
      res.json({ success: true, message: 'Task deleted' });
    });

    // Subtasks
    this.app.get('/api/v1/tasks/:id/subtasks', (req: Request, res: Response) => {
      const parentId = req.params.id;
      const subtasks = Array.from(this.tasks.values()).filter(t => t.parent_id === parentId);
      res.json({ success: true, data: subtasks });
    });

    this.app.post('/api/v1/tasks/:id/subtasks', (req: Request, res: Response) => {
      const task = this.createTask({ ...req.body, parent_id: req.params.id });
      this.tasks.set(task.id, task);
      res.json({ success: true, data: task });
    });

    // Projects API
    this.app.get('/api/v1/projects', (req: Request, res: Response) => {
      const projects = Array.from(this.projects.values()).map(p => ({
        ...p,
        task_count: Array.from(this.tasks.values()).filter(t => t.project_id === p.id).length
      }));
      res.json({ success: true, data: projects });
    });

    this.app.post('/api/v1/projects', (req: Request, res: Response) => {
      const project = this.createProject(req.body);
      this.projects.set(project.id, project);
      res.json({ success: true, data: project });
    });

    this.app.get('/api/v1/projects/:id/tasks', (req: Request, res: Response) => {
      const tasks = Array.from(this.tasks.values()).filter(t => t.project_id === req.params.id);
      res.json({ success: true, data: tasks });
    });

    // Templates
    this.app.get('/api/v1/templates', (req: Request, res: Response) => {
      res.json({
        success: true,
        data: [
          {
            id: 'bug-report',
            name: 'Bug Report',
            description: 'Template for reporting bugs',
            language_code: 'en',
            fields: [
              { key: 'title', label: 'Title', required: true },
              { key: 'description', label: 'Description', required: true },
              { key: 'severity', label: 'Severity', required: true, type: 'select', options: ['low', 'medium', 'high'] }
            ]
          },
          {
            id: 'feature',
            name: 'Feature Request',
            description: 'Template for new features',
            language_code: 'en',
            fields: [
              { key: 'title', label: 'Title', required: true },
              { key: 'description', label: 'Description', required: true },
              { key: 'type', label: 'Type', required: true, type: 'select', options: ['enhancement', 'new_feature', 'improvement'] }
            ]
          }
        ]
      });
    });

    // Summary
    this.app.get('/api/v1/summary', (req: Request, res: Response) => {
      const tasks = Array.from(this.tasks.values());
      res.json({
        success: true,
        data: {
          total_tasks: tasks.length,
          pending_tasks: tasks.filter(t => t.status === 'pending').length,
          in_progress_tasks: tasks.filter(t => t.status === 'in_progress').length,
          completed_tasks: tasks.filter(t => t.status === 'completed').length,
          total_projects: this.projects.size
        }
      });
    });
  }

  private initializeTestData(): void {
    // Create test project
    const project = this.createProject({
      name: 'Test Project',
      description: 'A test project for UI testing',
      workspace_id: 'test-workspace'
    });
    this.projects.set(project.id, project);

    // Create test tasks
    const task1 = this.createTask({
      title: 'Setup Test Environment',
      description: 'Configure testing infrastructure',
      status: 'completed',
      priority: 'high',
      project_id: project.id,
      workspace_id: 'test-workspace'
    });
    this.tasks.set(task1.id, task1);

    const task2 = this.createTask({
      title: 'Implement UI Tests',
      description: 'Create comprehensive UI testing suite',
      status: 'in_progress',
      priority: 'high',
      project_id: project.id,
      workspace_id: 'test-workspace'
    });
    this.tasks.set(task2.id, task2);

    const task3 = this.createTask({
      title: 'Write Documentation',
      description: 'Document testing approach and results',
      status: 'pending',
      priority: 'medium',
      project_id: project.id,
      workspace_id: 'test-workspace'
    });
    this.tasks.set(task3.id, task3);

    // Create subtask
    const subtask = this.createTask({
      title: 'Create Playwright Tests',
      description: 'Set up Playwright for UI testing',
      status: 'in_progress',
      priority: 'high',
      parent_id: task2.id,
      workspace_id: 'test-workspace'
    });
    this.tasks.set(subtask.id, subtask);
  }

  private createTask(data: Partial<MockTask>): MockTask {
    return {
      id: `task-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      title: data.title || 'Untitled Task',
      description: data.description || '',
      status: data.status || 'pending',
      priority: data.priority || 'medium',
      project_id: data.project_id,
      parent_id: data.parent_id,
      due_date: data.due_date,
      tags: data.tags || [],
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      workspace_id: data.workspace_id || 'test-workspace',
      dependency_count: 0,
      uncompleted_dependency_count: 0,
      dependent_on_this_count: 0
    };
  }

  private createProject(data: Partial<MockProject>): MockProject {
    return {
      id: `project-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      name: data.name || 'Untitled Project',
      description: data.description || '',
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString(),
      workspace_id: data.workspace_id || 'test-workspace'
    };
  }

  public async start(): Promise<void> {
    return new Promise((resolve) => {
      this.server = createServer(this.app);
      this.server.listen(this.port, () => {
        // Capture the actual port that was bound (in case port 0 was used)
        const address = this.server.address();
        this.port = typeof address === 'object' && address ? address.port : this.port;
        console.log(`[MockServer] Listening on port ${this.port}`);
        resolve();
      });
    });
  }

  public async stop(): Promise<void> {
    return new Promise((resolve) => {
      if (this.server) {
        this.server.close(() => {
          console.log('[MockServer] Stopped');
          resolve();
        });
      } else {
        resolve();
      }
    });
  }

  public getTasks(): MockTask[] {
    return Array.from(this.tasks.values());
  }

  public getProjects(): MockProject[] {
    return Array.from(this.projects.values());
  }
}

export default MockServer;
