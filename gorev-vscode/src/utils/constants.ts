// Application constants

export const APP_NAME = 'Gorev';
export const APP_ID = 'gorev-vscode';

export const DEFAULT_TIMEOUT = 5000; // 5 seconds

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
  CREATE_SUBTASK: 'gorev.createSubtask',
  CHANGE_PARENT: 'gorev.changeParent',
  REMOVE_PARENT: 'gorev.removeParent',
  ADD_DEPENDENCY: 'gorev.addDependency',
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
  
  // Enhanced TreeView commands
  SELECT_TASK: 'gorev.selectTask',
  SET_GROUPING: 'gorev.setGrouping',
  SET_SORTING: 'gorev.setSorting',
  FILTER_TASKS: 'gorev.filterTasks',
  CLEAR_FILTER: 'gorev.clearFilter',
  TOGGLE_SHOW_COMPLETED: 'gorev.toggleShowCompleted',
  SELECT_ALL: 'gorev.selectAll',
  DESELECT_ALL: 'gorev.deselectAll',
  BULK_UPDATE_STATUS: 'gorev.bulkUpdateStatus',
  BULK_DELETE: 'gorev.bulkDelete',
  
  // Inline edit commands
  EDIT_TASK_TITLE: 'gorev.editTaskTitle',
  QUICK_STATUS_CHANGE: 'gorev.quickStatusChange',
  QUICK_PRIORITY_CHANGE: 'gorev.quickPriorityChange',
  QUICK_DATE_CHANGE: 'gorev.quickDateChange',
  DETAILED_EDIT: 'gorev.detailedEdit',
  
  // Template commands
  OPEN_TEMPLATE_WIZARD: 'gorev.openTemplateWizard',
  CREATE_FROM_TEMPLATE: 'gorev.createFromTemplate',
  QUICK_CREATE_FROM_TEMPLATE: 'gorev.quickCreateFromTemplate',
  REFRESH_TEMPLATES: 'gorev.refreshTemplates',
  INIT_DEFAULT_TEMPLATES: 'gorev.initDefaultTemplates',
  SHOW_TEMPLATE_DETAILS: 'gorev.showTemplateDetails',
  EXPORT_TEMPLATE: 'gorev.exportTemplate',
  
  // Data export/import commands
  EXPORT_DATA: 'gorev.exportData',
  IMPORT_DATA: 'gorev.importData',
  EXPORT_CURRENT_VIEW: 'gorev.exportCurrentView',
  QUICK_EXPORT: 'gorev.quickExport',
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