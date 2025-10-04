/**
 * Workspace models for multi-workspace support
 */

export interface WorkspaceInfo {
  id: string;
  name: string;
  path: string;
  database_path: string;
  last_accessed: string;
  created_at: string;
  task_count: number;
}

export interface WorkspaceRegistration {
  path: string;
  name?: string;
}

export interface WorkspaceRegistrationResponse {
  success: boolean;
  workspace_id: string;
  workspace: WorkspaceInfo;
}

export interface WorkspaceListResponse {
  success: boolean;
  workspaces: WorkspaceInfo[];
  total: number;
}

export interface WorkspaceContext {
  workspaceId: string;
  workspacePath: string;
  workspaceName: string;
}

/**
 * Interface for workspace-aware operations
 */
export interface WorkspaceAware {
  getWorkspaceContext(): WorkspaceContext | undefined;
  setWorkspaceContext(context: WorkspaceContext): void;
  clearWorkspaceContext(): void;
}
