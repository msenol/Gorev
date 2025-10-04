import { useState, useEffect } from 'react';
import { useQuery, useQueryClient } from '@tanstack/react-query';
import Header from './components/Header';
import Sidebar from './components/Sidebar';
import TaskList from './components/TaskList';
import ProjectSelector from './components/ProjectSelector';
import { getProjects, getTasks, getTemplates } from './api/client';
import type { Project, TaskFilter } from './types';
import { useLanguage } from './contexts/LanguageContext';
import { useWorkspace } from './contexts/WorkspaceContext';

function App() {
  const { t } = useLanguage();
  const { currentWorkspace, isLoading: workspaceLoading, error: workspaceError } = useWorkspace();
  const queryClient = useQueryClient();
  const [selectedProject, setSelectedProject] = useState<Project | null>(null);
  const [taskFilter, setTaskFilter] = useState<TaskFilter>({});
  const [sidebarOpen, setSidebarOpen] = useState(true);

  // Invalidate all queries when workspace changes
  useEffect(() => {
    if (currentWorkspace) {
      queryClient.invalidateQueries();
    }
  }, [currentWorkspace?.workspaceId]);

  // Load initial data (only after workspace is ready)
  const { data: projectsData, isLoading: projectsLoading } = useQuery({
    queryKey: ['projects', currentWorkspace?.workspaceId],
    queryFn: getProjects,
    enabled: !!currentWorkspace, // Only run when workspace is set
  });

  const { data: tasksData, isLoading: tasksLoading, refetch: refetchTasks } = useQuery({
    queryKey: ['tasks', currentWorkspace?.workspaceId, selectedProject?.id, taskFilter],
    queryFn: () => getTasks({
      ...taskFilter,
      proje_id: selectedProject?.id,
    }),
    enabled: !!currentWorkspace, // Only run when workspace is set
  });

  const { data: templatesData } = useQuery({
    queryKey: ['templates', currentWorkspace?.workspaceId],
    queryFn: () => getTemplates(),
    enabled: !!currentWorkspace, // Only run when workspace is set
  });

  const projects = projectsData?.data || [];
  const tasks = tasksData?.data || [];
  const templates = templatesData?.data || [];

  // Show loading state while workspace is loading
  if (workspaceLoading) {
    return (
      <div className="h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="text-6xl mb-4">⏳</div>
          <h2 className="text-xl font-semibold text-gray-700">Loading workspace...</h2>
        </div>
      </div>
    );
  }

  // Show error if workspace failed to load
  if (workspaceError) {
    return (
      <div className="h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center max-w-md">
          <div className="text-6xl mb-4">⚠️</div>
          <h2 className="text-xl font-semibold text-gray-900 mb-2">Workspace Error</h2>
          <p className="text-gray-600">{workspaceError}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-screen flex bg-gray-50">
      {/* Sidebar */}
      <Sidebar
        isOpen={sidebarOpen}
        onToggle={() => setSidebarOpen(!sidebarOpen)}
        projects={projects}
        selectedProject={selectedProject}
        onProjectSelect={setSelectedProject}
        templates={templates}
        onTaskCreated={() => refetchTasks()}
      />

      {/* Main Content */}
      <div className="flex-1 flex flex-col overflow-hidden">
        {/* Header */}
        <Header
          selectedProject={selectedProject}
          taskFilter={taskFilter}
          onFilterChange={setTaskFilter}
          onSidebarToggle={() => setSidebarOpen(!sidebarOpen)}
          totalTasks={tasks.length}
        />

        {/* Content Area */}
        <main className="flex-1 overflow-auto">
          <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
            {/* Project Selector */}
            {!selectedProject && projects.length > 0 && (
              <ProjectSelector
                projects={projects}
                onProjectSelect={setSelectedProject}
                loading={projectsLoading}
              />
            )}

            {/* Task List */}
            <TaskList
              tasks={tasks}
              loading={tasksLoading}
              onTaskUpdate={() => refetchTasks()}
            />

            {/* Empty State */}
            {!tasksLoading && tasks.length === 0 && (
              <div className="text-center py-12">
                <div className="text-gray-400 text-6xl mb-4">📝</div>
                <h3 className="text-lg font-medium text-gray-900 mb-2">
                  {selectedProject ? t('no_project_tasks') : t('no_tasks_yet')}
                </h3>
                <p className="text-gray-500 mb-6">
                  {templates.length > 0
                    ? t('select_template_to_create')
                    : t('wait_templates_to_load')}
                </p>
              </div>
            )}
          </div>
        </main>
      </div>
    </div>
  );
}

export default App;