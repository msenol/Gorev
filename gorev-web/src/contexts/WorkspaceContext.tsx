/**
 * WorkspaceContext
 *
 * Provides workspace context and management for the entire Web UI application.
 * Handles workspace selection, switching, and automatic workspace header injection.
 */

import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { WorkspaceContext as IWorkspaceContext, WorkspaceInfo } from '@/types';
import { getWorkspaces, setWorkspaceContext as setApiWorkspaceContext } from '@/api/client';

interface WorkspaceContextValue {
  currentWorkspace: IWorkspaceContext | null;
  availableWorkspaces: WorkspaceInfo[];
  isLoading: boolean;
  error: string | null;
  selectWorkspace: (workspaceId: string) => void;
  refreshWorkspaces: () => Promise<void>;
}

const WorkspaceContext = createContext<WorkspaceContextValue | undefined>(undefined);

interface WorkspaceProviderProps {
  children: ReactNode;
}

export const WorkspaceProvider: React.FC<WorkspaceProviderProps> = ({ children }) => {
  const [currentWorkspace, setCurrentWorkspace] = useState<IWorkspaceContext | null>(null);
  const [availableWorkspaces, setAvailableWorkspaces] = useState<WorkspaceInfo[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Load workspaces from server
  const loadWorkspaces = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const response = await getWorkspaces();

      if (response.success && response.workspaces.length > 0) {
        setAvailableWorkspaces(response.workspaces);
      } else {
        setError('No workspaces registered. Please use VS Code extension or MCP client to register a workspace.');
      }
    } catch (err) {
      console.error('Failed to load workspaces:', err);
      setError('Failed to connect to server. Please make sure the server is running.');
    } finally {
      setIsLoading(false);
    }
  };

  // Select a workspace by ID
  const selectWorkspace = (workspaceId: string) => {
    const workspace = availableWorkspaces.find(w => w.id === workspaceId);

    if (workspace) {
      const context: IWorkspaceContext = {
        workspaceId: workspace.id,
        workspaceName: workspace.name,
        workspacePath: workspace.path
      };

      setCurrentWorkspace(context);
      setApiWorkspaceContext(context);

      // Store in localStorage for persistence
      localStorage.setItem('gorev:selected-workspace', workspaceId);
    }
  };

  // Refresh workspaces from server
  const refreshWorkspaces = async () => {
    await loadWorkspaces();
  };

  // Load workspaces on mount
  useEffect(() => {
    loadWorkspaces();
  }, []);

  // Try to restore selected workspace from localStorage
  useEffect(() => {
    if (availableWorkspaces.length > 0 && !currentWorkspace) {
      const savedWorkspaceId = localStorage.getItem('gorev:selected-workspace');

      if (savedWorkspaceId && availableWorkspaces.find(w => w.id === savedWorkspaceId)) {
        selectWorkspace(savedWorkspaceId);
      } else {
        // Auto-select first workspace if no saved preference
        selectWorkspace(availableWorkspaces[0].id);
      }
    }
  }, [availableWorkspaces]);

  const value: WorkspaceContextValue = {
    currentWorkspace,
    availableWorkspaces,
    isLoading,
    error,
    selectWorkspace,
    refreshWorkspaces
  };

  return (
    <WorkspaceContext.Provider value={value}>
      {children}
    </WorkspaceContext.Provider>
  );
};

export const useWorkspace = (): WorkspaceContextValue => {
  const context = useContext(WorkspaceContext);

  if (context === undefined) {
    throw new Error('useWorkspace must be used within a WorkspaceProvider');
  }

  return context;
};
