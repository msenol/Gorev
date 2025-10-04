import React, { useState } from 'react';
import { X, FolderOpen, Plus, FileText } from 'lucide-react';
import type { Project, Template } from '@/types';
import CreateTaskModal from './CreateTaskModal';

interface SidebarProps {
  isOpen: boolean;
  onToggle: () => void;
  projects: Project[];
  selectedProject: Project | null;
  onProjectSelect: (project: Project | null) => void;
  templates: Template[];
  onTaskCreated: () => void;
}

const Sidebar: React.FC<SidebarProps> = ({
  isOpen,
  onToggle,
  projects,
  selectedProject,
  onProjectSelect,
  templates,
  onTaskCreated,
}) => {
  const [createTaskModalOpen, setCreateTaskModalOpen] = useState(false);

  if (!isOpen) {
    return (
      <div className="w-16 bg-white border-r border-gray-200 flex flex-col items-center py-4">
        <button
          onClick={onToggle}
          className="p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100"
        >
          <FolderOpen className="h-5 w-5" />
        </button>
      </div>
    );
  }

  return (
    <>
      <div className="w-80 bg-white border-r border-gray-200 flex flex-col">
        {/* Header */}
        <div className="p-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <h2 className="text-lg font-semibold text-gray-900">Projeler</h2>
            <button
              onClick={onToggle}
              className="p-1 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100"
            >
              <X className="h-4 w-4" />
            </button>
          </div>
        </div>

        {/* Projects */}
        <div className="flex-1 overflow-auto p-4">
          <div className="space-y-2">
            {/* All Projects */}
            <button
              onClick={() => onProjectSelect(null)}
              className={`w-full text-left px-3 py-2 rounded-lg transition-colors ${
                !selectedProject
                  ? 'bg-primary-50 text-primary-700 border border-primary-200'
                  : 'text-gray-700 hover:bg-gray-100'
              }`}
            >
              <div className="flex items-center">
                <FolderOpen className="h-4 w-4 mr-2" />
                <span className="font-medium">Tüm Projeler</span>
              </div>
              <div className="text-xs text-gray-500 ml-6">
                {projects.reduce((total, p) => total + (p.gorev_sayisi || 0), 0)} görev
              </div>
            </button>

            {/* Individual Projects */}
            {projects.map((project) => (
              <button
                key={project.id}
                onClick={() => onProjectSelect(project)}
                className={`w-full text-left px-3 py-2 rounded-lg transition-colors ${
                  selectedProject?.id === project.id
                    ? 'bg-primary-50 text-primary-700 border border-primary-200'
                    : 'text-gray-700 hover:bg-gray-100'
                }`}
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <div className="w-3 h-3 rounded-full bg-primary-500 mr-2"></div>
                    <span className="font-medium truncate">{project.isim}</span>
                  </div>
                  {project.is_active && (
                    <div className="w-2 h-2 rounded-full bg-green-500"></div>
                  )}
                </div>
                <div className="text-xs text-gray-500 ml-5">
                  {project.gorev_sayisi || 0} görev
                </div>
              </button>
            ))}
          </div>

          {/* Templates Section */}
          <div className="mt-8">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-gray-900">Şablonlar</h3>
              <button
                onClick={() => setCreateTaskModalOpen(true)}
                className="p-1 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100"
                title="Yeni görev oluştur"
              >
                <Plus className="h-4 w-4" />
              </button>
            </div>

            <div className="space-y-1">
              {templates.slice(0, 8).map((template) => (
                <button
                  key={template.id}
                  onClick={() => setCreateTaskModalOpen(true)}
                  className="w-full text-left px-3 py-2 rounded-md text-sm text-gray-600 hover:bg-gray-100 transition-colors"
                >
                  <div className="flex items-center">
                    <FileText className="h-3 w-3 mr-2" />
                    <span className="truncate">{template.isim}</span>
                  </div>
                </button>
              ))}
            </div>

            {templates.length > 8 && (
              <button
                onClick={() => setCreateTaskModalOpen(true)}
                className="w-full text-left px-3 py-2 rounded-md text-sm text-gray-500 hover:bg-gray-100 transition-colors"
              >
                +{templates.length - 8} daha fazla şablon...
              </button>
            )}
          </div>
        </div>
      </div>

      {/* Create Task Modal */}
      <CreateTaskModal
        isOpen={createTaskModalOpen}
        onClose={() => setCreateTaskModalOpen(false)}
        templates={templates}
        selectedProject={selectedProject}
        onTaskCreated={onTaskCreated}
      />
    </>
  );
};

export default Sidebar;