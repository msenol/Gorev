import React from 'react';
import { Menu, Search, Filter } from 'lucide-react';
import type { Project, TaskFilter, TaskStatus, TaskPriority } from '@/types';
import LanguageSwitcher from './LanguageSwitcher';
import { WorkspaceSwitcher } from './WorkspaceSwitcher';
import { useLanguage } from '../contexts/LanguageContext';

interface HeaderProps {
  selectedProject: Project | null;
  taskFilter: TaskFilter;
  onFilterChange: (filter: TaskFilter) => void;
  onSidebarToggle: () => void;
  totalTasks: number;
}

const Header: React.FC<HeaderProps> = ({
  selectedProject,
  taskFilter,
  onFilterChange,
  onSidebarToggle,
  totalTasks,
}) => {
  const { t } = useLanguage();
  const handleStatusFilter = (status: TaskStatus | undefined) => {
    onFilterChange({
      ...taskFilter,
      status,
    });
  };

  const handlePriorityFilter = (priority: TaskPriority | undefined) => {
    onFilterChange({
      ...taskFilter,
      priority,
    });
  };

  const handleSearchChange = (search: string) => {
    onFilterChange({
      ...taskFilter,
      search: search || undefined,
    });
  };

  return (
    <header className="bg-white shadow-sm border-b border-gray-200">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Left side */}
          <div className="flex items-center">
            <button
              onClick={onSidebarToggle}
              className="p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 lg:hidden"
            >
              <Menu className="h-5 w-5" />
            </button>

            <div className="ml-2 lg:ml-0">
              <h1 className="text-xl font-semibold text-gray-900">
                🚀 {t('app.title')}
              </h1>
              {selectedProject && (
                <p className="text-sm text-gray-500">
                  {selectedProject.name} • {totalTasks} {t('tasks')}
                </p>
              )}
            </div>
          </div>

          {/* Search, Filters and Language Switcher */}
          <div className="flex items-center space-x-4">
            {/* Search */}
            <div className="relative">
              <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 h-4 w-4" />
              <input
                type="text"
                placeholder={t('search_tasks')}
                value={taskFilter.search || ''}
                onChange={(e) => handleSearchChange(e.target.value)}
                className="pl-10 pr-4 py-2 w-64 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                data-testid="search-input"
              />
            </div>

            {/* Filters */}
            <div className="flex items-center space-x-2" data-testid="filter-section">
              <Filter className="h-4 w-4 text-gray-400" />

              {/* Status Filter */}
              <select
                value={taskFilter.status || ''}
                onChange={(e) => handleStatusFilter(e.target.value as TaskStatus || undefined)}
                className="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                data-testid="status-filter"
              >
                <option value="">{t('all_statuses')}</option>
                <option value="beklemede">{t('pending')}</option>
                <option value="devam_ediyor">{t('in_progress')}</option>
                <option value="tamamlandi">{t('completed')}</option>
              </select>

              {/* Priority Filter */}
              <select
                value={taskFilter.priority || ''}
                onChange={(e) => handlePriorityFilter(e.target.value as TaskPriority || undefined)}
                className="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary-500"
                data-testid="priority-filter"
              >
                <option value="">{t('all_priorities')}</option>
                <option value="yuksek">{t('high')}</option>
                <option value="orta">{t('medium')}</option>
                <option value="dusuk">{t('low')}</option>
              </select>
            </div>

            {/* Workspace Switcher */}
            <WorkspaceSwitcher />

            {/* Language Switcher */}
            <LanguageSwitcher />
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;