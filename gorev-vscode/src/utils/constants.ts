// Application constants

export const APP_NAME = 'Gorev';
export const APP_ID = 'gorev-vscode';

export const DEFAULT_TIMEOUT = 30000; // 30 seconds

export const ICONS = {
  // Task status icons
  TASK_PENDING: 'circle-outline',
  TASK_IN_PROGRESS: 'sync~spin',
  TASK_COMPLETED: 'check',
  
  // Priority icons
  PRIORITY_HIGH: 'flame',
  PRIORITY_MEDIUM: 'dash',
  PRIORITY_LOW: 'arrow-down',
  
  // General icons
  PROJECT: 'folder',
  PROJECT_ACTIVE: 'folder-opened',
  TEMPLATE: 'file-code',
  TAG: 'tag',
  CALENDAR: 'calendar',
  DEPENDENCY: 'references',
  ERROR: 'error',
  WARNING: 'warning',
  INFO: 'info',
} as const;

export const COLORS = {
  HIGH_PRIORITY: 'gorev.highPriorityForeground',
  MEDIUM_PRIORITY: 'gorev.mediumPriorityForeground',
  LOW_PRIORITY: 'gorev.lowPriorityForeground',
} as const;

export const COMMANDS = {
  // Task commands
  CREATE_TASK: 'gorev.createTask',
  REFRESH_TASKS: 'gorev.refreshTasks',
  SHOW_TASK_DETAIL: 'gorev.showTaskDetail',
  UPDATE_TASK_STATUS: 'gorev.updateTaskStatus',
  DELETE_TASK: 'gorev.deleteTask',
  QUICK_CREATE_TASK: 'gorev.quickCreateTask',
  
  // Project commands
  CREATE_PROJECT: 'gorev.createProject',
  SET_ACTIVE_PROJECT: 'gorev.setActiveProject',
  
  // General commands
  SHOW_SUMMARY: 'gorev.showSummary',
  CONNECT: 'gorev.connect',
  DISCONNECT: 'gorev.disconnect',
} as const;

export const VIEWS = {
  TASKS: 'gorevTasks',
  PROJECTS: 'gorevProjects',
  TEMPLATES: 'gorevTemplates',
} as const;

export const CONTEXT_VALUES = {
  TASK: 'task',
  PROJECT: 'project',
  TEMPLATE: 'template',
  PROJECT_ACTIVE: 'project-active',
} as const;