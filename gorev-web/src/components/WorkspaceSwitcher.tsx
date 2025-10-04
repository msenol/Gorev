/**
 * WorkspaceSwitcher Component
 *
 * Dropdown component for switching between registered workspaces.
 * Displays current workspace and allows user to select different workspace.
 */

import React, { useState } from 'react';
import { useWorkspace } from '@/contexts/WorkspaceContext';

interface WorkspaceSwitcherProps {
  className?: string;
}

export const WorkspaceSwitcher: React.FC<WorkspaceSwitcherProps> = ({ className = '' }) => {
  const { currentWorkspace, availableWorkspaces, selectWorkspace, isLoading, error, refreshWorkspaces } = useWorkspace();
  const [isOpen, setIsOpen] = useState(false);

  const handleWorkspaceSelect = (workspaceId: string) => {
    selectWorkspace(workspaceId);
    setIsOpen(false);
  };

  const handleRefresh = async (e: React.MouseEvent) => {
    e.stopPropagation();
    await refreshWorkspaces();
  };

  if (isLoading) {
    return (
      <div className={`workspace-switcher loading ${className}`}>
        <div className="workspace-display">
          <svg className="workspace-icon spin" width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M8 1.5V4M8 12V14.5M3.5 8H1M15 8H12.5M4.22 4.22L2.5 2.5M13.5 13.5L11.78 11.78M11.78 4.22L13.5 2.5M2.5 13.5L4.22 11.78"
                  stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
          </svg>
          <span className="workspace-name">Loading workspaces...</span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`workspace-switcher error ${className}`}>
        <div className="workspace-display error">
          <svg className="workspace-icon" width="16" height="16" viewBox="0 0 16 16" fill="none">
            <path d="M8 1L1 14h14L8 1z" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
            <path d="M8 6v3M8 11h.01" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
          </svg>
          <span className="workspace-name">{error}</span>
        </div>
      </div>
    );
  }

  if (!currentWorkspace) {
    return null;
  }

  return (
    <div className={`workspace-switcher ${className}`}>
      <button
        className="workspace-display"
        onClick={() => setIsOpen(!isOpen)}
        aria-expanded={isOpen}
        aria-haspopup="true"
      >
        <svg className="workspace-icon" width="16" height="16" viewBox="0 0 16 16" fill="none">
          <path d="M2 4.5L8 1.5L14 4.5V11.5L8 14.5L2 11.5V4.5Z" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
          <path d="M8 8V14.5M8 8L2 4.5M8 8L14 4.5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
        </svg>
        <span className="workspace-name">{currentWorkspace.workspaceName}</span>
        <svg className={`dropdown-icon ${isOpen ? 'open' : ''}`} width="12" height="12" viewBox="0 0 12 12" fill="none">
          <path d="M3 4.5L6 7.5L9 4.5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
        </svg>
      </button>

      {isOpen && (
        <>
          <div className="workspace-dropdown-backdrop" onClick={() => setIsOpen(false)} />
          <div className="workspace-dropdown">
            <div className="workspace-dropdown-header">
              <span>Select Workspace</span>
              <button
                className="refresh-button"
                onClick={handleRefresh}
                title="Refresh workspaces"
              >
                <svg width="14" height="14" viewBox="0 0 16 16" fill="none">
                  <path d="M14 8c0 3.314-2.686 6-6 6s-6-2.686-6-6 2.686-6 6-6c1.657 0 3.157.671 4.243 1.757L14 6"
                        stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
                  <path d="M14 2v4h-4" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round"/>
                </svg>
              </button>
            </div>
            <div className="workspace-list">
              {availableWorkspaces.map(workspace => (
                <button
                  key={workspace.id}
                  className={`workspace-item ${workspace.id === currentWorkspace.workspaceId ? 'active' : ''}`}
                  onClick={() => handleWorkspaceSelect(workspace.id)}
                >
                  <div className="workspace-item-content">
                    <span className="workspace-item-name">{workspace.name}</span>
                    <span className="workspace-item-path">{workspace.path}</span>
                  </div>
                  <div className="workspace-item-meta">
                    <span className="workspace-item-tasks">{workspace.task_count} tasks</span>
                  </div>
                  {workspace.id === currentWorkspace.workspaceId && (
                    <svg className="check-icon" width="16" height="16" viewBox="0 0 16 16" fill="none">
                      <path d="M3 8l3 3 7-7" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"/>
                    </svg>
                  )}
                </button>
              ))}
            </div>
          </div>
        </>
      )}

      <style>{`
        .workspace-switcher {
          position: relative;
          font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', sans-serif;
        }

        .workspace-display {
          display: flex;
          align-items: center;
          gap: 8px;
          padding: 6px 12px;
          background: #f5f5f5;
          border: 1px solid #e0e0e0;
          border-radius: 6px;
          cursor: pointer;
          transition: all 0.2s ease;
          font-size: 14px;
          color: #333;
        }

        .workspace-display:hover {
          background: #ececec;
          border-color: #d0d0d0;
        }

        .workspace-display.error {
          background: #fff4f4;
          border-color: #ffcccc;
          color: #cc0000;
          cursor: default;
        }

        .workspace-icon {
          flex-shrink: 0;
          color: #666;
        }

        .workspace-display.error .workspace-icon {
          color: #cc0000;
        }

        .workspace-icon.spin {
          animation: spin 1s linear infinite;
        }

        @keyframes spin {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }

        .workspace-name {
          flex: 1;
          font-weight: 500;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .dropdown-icon {
          flex-shrink: 0;
          transition: transform 0.2s ease;
        }

        .dropdown-icon.open {
          transform: rotate(180deg);
        }

        .workspace-dropdown-backdrop {
          position: fixed;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          z-index: 999;
        }

        .workspace-dropdown {
          position: absolute;
          top: calc(100% + 4px);
          left: 0;
          min-width: 320px;
          max-width: 480px;
          background: white;
          border: 1px solid #e0e0e0;
          border-radius: 8px;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
          z-index: 1000;
          overflow: hidden;
        }

        .workspace-dropdown-header {
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: 12px 16px;
          background: #f8f8f8;
          border-bottom: 1px solid #e0e0e0;
          font-weight: 600;
          font-size: 13px;
          color: #666;
        }

        .refresh-button {
          display: flex;
          align-items: center;
          justify-content: center;
          padding: 4px;
          background: transparent;
          border: none;
          border-radius: 4px;
          cursor: pointer;
          color: #666;
          transition: all 0.2s ease;
        }

        .refresh-button:hover {
          background: #e0e0e0;
          color: #333;
        }

        .workspace-list {
          max-height: 400px;
          overflow-y: auto;
        }

        .workspace-item {
          display: flex;
          align-items: center;
          gap: 12px;
          width: 100%;
          padding: 12px 16px;
          background: white;
          border: none;
          border-bottom: 1px solid #f0f0f0;
          cursor: pointer;
          transition: background 0.2s ease;
          text-align: left;
        }

        .workspace-item:hover {
          background: #f8f8f8;
        }

        .workspace-item.active {
          background: #e3f2fd;
        }

        .workspace-item:last-child {
          border-bottom: none;
        }

        .workspace-item-content {
          flex: 1;
          min-width: 0;
        }

        .workspace-item-name {
          display: block;
          font-weight: 500;
          font-size: 14px;
          color: #333;
          margin-bottom: 2px;
        }

        .workspace-item-path {
          display: block;
          font-size: 12px;
          color: #999;
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }

        .workspace-item-meta {
          flex-shrink: 0;
          font-size: 11px;
          color: #999;
        }

        .check-icon {
          flex-shrink: 0;
          color: #1976d2;
        }
      `}</style>
    </div>
  );
};
