import React from 'react';
import { FolderOpen, Users } from 'lucide-react';
import type { Project } from '@/types';

interface ProjectSelectorProps {
  projects: Project[];
  onProjectSelect: (project: Project) => void;
  loading: boolean;
}

const ProjectSelector: React.FC<ProjectSelectorProps> = ({
  projects,
  onProjectSelect,
  loading,
}) => {
  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
        {[...Array(3)].map((_, i) => (
          <div key={i} className="card animate-pulse">
            <div className="h-6 bg-gray-200 rounded w-3/4 mb-2"></div>
            <div className="h-4 bg-gray-200 rounded w-1/2 mb-3"></div>
            <div className="h-3 bg-gray-200 rounded w-1/4"></div>
          </div>
        ))}
      </div>
    );
  }

  return (
    <div className="mb-8">
      <h2 className="text-xl font-semibold text-gray-900 mb-4">
        Bir proje seçin
      </h2>
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {projects.map((project) => (
          <button
            key={project.id}
            onClick={() => onProjectSelect(project)}
            className="card text-left hover:shadow-md transition-shadow border-2 border-transparent hover:border-primary-200"
          >
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center">
                <FolderOpen className="h-5 w-5 text-primary-500 mr-2" />
                <h3 className="font-semibold text-gray-900">{project.isim}</h3>
              </div>
              {project.is_active && (
                <div className="w-2 h-2 rounded-full bg-green-500"></div>
              )}
            </div>

            {project.tanim && (
              <p className="text-gray-600 text-sm mb-3 line-clamp-2">
                {project.tanim}
              </p>
            )}

            <div className="flex items-center justify-between text-xs text-gray-500">
              <div className="flex items-center">
                <Users className="h-3 w-3 mr-1" />
                {project.gorev_sayisi || 0} görev
              </div>
              <span>
                {project.olusturma_tarihi ? new Date(project.olusturma_tarihi).toLocaleDateString('tr-TR') : ''}
              </span>
            </div>
          </button>
        ))}
      </div>
    </div>
  );
};

export default ProjectSelector;