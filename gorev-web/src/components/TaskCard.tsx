import React, { useState } from 'react';
import { Clock, Calendar, ChevronDown, ChevronRight, Link2, AlertCircle } from 'lucide-react';
import { useMutation } from '@tanstack/react-query';
import { updateTask, deleteTask } from '@/api/client';
import type { Task, TaskStatus, TaskPriority } from '@/types';

interface TaskCardProps {
  task: Task;
  onUpdate: () => void;
}

const TaskCard: React.FC<TaskCardProps> = ({ task, onUpdate }) => {
  const [showActions, setShowActions] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [showSubtasks, setShowSubtasks] = useState(false);

  const updateMutation = useMutation({
    mutationFn: ({ id, updates }: { id: string; updates: any }) =>
      updateTask(id, updates),
    onSuccess: () => {
      onUpdate();
      setIsEditing(false);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: deleteTask,
    onSuccess: () => {
      onUpdate();
    },
  });

  const handleStatusChange = (status: TaskStatus) => {
    updateMutation.mutate({
      id: task.id,
      updates: { status },
    });
  };

  const handleDelete = () => {
    if (confirm('Bu görevi silmek istediğinizden emin misiniz?')) {
      deleteMutation.mutate(task.id);
    }
  };

  const getPriorityColor = (priority: TaskPriority) => {
    switch (priority) {
      case 'yuksek':
        return 'bg-red-100 text-red-800 border-red-200';
      case 'orta':
        return 'bg-orange-100 text-orange-800 border-orange-200';
      case 'dusuk':
        return 'bg-gray-100 text-gray-800 border-gray-200';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  const getStatusColor = (status: TaskStatus) => {
    switch (status) {
      case 'beklemede':
        return 'bg-yellow-100 text-yellow-800';
      case 'devam_ediyor':
        return 'bg-blue-100 text-blue-800';
      case 'tamamlandi':
        return 'bg-green-100 text-green-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const formatDate = (dateString: string) => {
    if (!dateString) return '';
    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) return ''; // Invalid date check
      return date.toLocaleDateString('tr-TR', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
      });
    } catch {
      return '';
    }
  };

  return (
    <div className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow" data-testid="task-item" data-task-id={task.id}>
      <div className="flex items-start justify-between">
        <div className="flex-1">
          {/* Title */}
          <h4 className="text-lg font-medium text-gray-900 mb-2 line-clamp-2" data-testid="task-title">
            {task.title}
          </h4>

          {/* Description */}
          {task.description && (
            <p className="text-gray-600 text-sm mb-3 line-clamp-3" data-testid="task-description">
              {task.description}
            </p>
          )}

          {/* Meta Information */}
          <div className="flex items-center space-x-4 text-xs text-gray-500 mb-3">
            <div className="flex items-center">
              <Clock className="h-3 w-3 mr-1" />
              {formatDate(task.created_at)}
            </div>

            {task.due_date && (
              <div className="flex items-center">
                <Calendar className="h-3 w-3 mr-1" />
                {formatDate(task.due_date)}
              </div>
            )}

            {task.proje_name && (
              <div className="flex items-center">
                <div className="w-2 h-2 rounded-full bg-primary-500 mr-1"></div>
                {task.proje_name}
              </div>
            )}
          </div>

          {/* Tags */}
          {task.tags && task.tags.length > 0 && (
            <div className="flex flex-wrap gap-1 mb-3">
              {task.tags.map((tag) => (
                <span
                  key={tag.id}
                  className="inline-flex items-center px-2 py-1 rounded-full text-xs bg-gray-100 text-gray-700"
                >
                  {tag.name}
                </span>
              ))}
            </div>
          )}

          {/* Dependencies */}
          {((task.dependency_count ?? 0) > 0 || (task.uncompleted_dependency_count ?? 0) > 0) && (
            <div className="flex items-center gap-2 mb-3 text-xs">
              {(task.dependency_count ?? 0) > 0 && (
                <div className="flex items-center text-gray-600">
                  <Link2 className="h-3 w-3 mr-1" />
                  <span>{task.dependency_count} bağımlılık</span>
                </div>
              )}
              {(task.uncompleted_dependency_count ?? 0) > 0 && (
                <div className="flex items-center text-orange-600">
                  <AlertCircle className="h-3 w-3 mr-1" />
                  <span>{task.uncompleted_dependency_count} bekliyor</span>
                </div>
              )}
            </div>
          )}

          {/* Subtasks */}
          {task.subtasks && task.subtasks.length > 0 && (
            <div className="mb-3" data-testid="subtasks-section">
              <button
                onClick={() => setShowSubtasks(!showSubtasks)}
                className="flex items-center text-sm text-gray-600 hover:text-gray-900"
                data-testid="expand-button"
              >
                {showSubtasks ? (
                  <ChevronDown className="h-4 w-4 mr-1" />
                ) : (
                  <ChevronRight className="h-4 w-4 mr-1" />
                )}
                <span data-testid="subtask-count">{task.subtasks.length} alt görev</span>
              </button>

              {showSubtasks && (
                <div className="mt-2 ml-5 space-y-2 border-l-2 border-gray-200 pl-3" data-testid="subtasks-list">
                  {task.subtasks.map((subtask) => (
                    <div key={subtask.id} className="text-sm" data-testid="subtask-item" data-subtask-id={subtask.id}>
                      <div className="flex items-center gap-2">
                        <span className={`inline-block w-2 h-2 rounded-full ${
                          subtask.status === 'tamamlandi'
                            ? 'bg-green-500'
                            : subtask.status === 'devam_ediyor'
                            ? 'bg-blue-500'
                            : 'bg-gray-300'
                        }`} data-testid="subtask-status-indicator"></span>
                        <span className={subtask.status === 'tamamlandi' ? 'line-through text-gray-500' : 'text-gray-700'} data-testid="subtask-title">
                          {subtask.title}
                        </span>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          )}

          {/* Status and Priority */}
          <div className="flex items-center space-x-2">
            <select
              value={task.status}
              onChange={(e) => handleStatusChange(e.target.value as TaskStatus)}
              className={`px-2 py-1 rounded-full text-xs font-medium border ${getStatusColor(
                task.status
              )} focus:outline-none focus:ring-2 focus:ring-primary-500`}
              disabled={updateMutation.isPending}
              data-testid="status-select"
            >
              <option value="beklemede">Beklemede</option>
              <option value="devam_ediyor">Devam Ediyor</option>
              <option value="tamamlandi">Tamamlandı</option>
            </select>

            <span
              className={`px-2 py-1 rounded-full text-xs font-medium border ${getPriorityColor(
                task.priority
              )}`}
              data-testid="priority-badge"
            >
              {task.priority === 'yuksek'
                ? 'Yüksek'
                : task.priority === 'orta'
                ? 'Orta'
                : 'Düşük'}
            </span>
          </div>
        </div>

        {/* Actions */}
        <div className="relative ml-4">
          <button
            onClick={() => setShowActions(!showActions)}
            className="p-1 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100"
            title="Menü"
            data-testid="task-menu-button"
          >
            ⋮
          </button>

          {showActions && (
            <div className="absolute right-0 top-8 z-10 w-48 bg-white border border-gray-200 rounded-md shadow-lg" data-testid="context-menu">
              <button
                onClick={() => {
                  setIsEditing(true);
                  setShowActions(false);
                }}
                className="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 flex items-center"
                data-testid="menu-item-edit"
              >
                ✏️ Düzenle
              </button>
              <button
                onClick={handleDelete}
                disabled={deleteMutation.isPending}
                className="w-full text-left px-4 py-2 text-sm text-red-700 hover:bg-red-50 flex items-center"
                data-testid="menu-item-delete"
              >
                🗑️ Sil
              </button>
            </div>
          )}
        </div>
      </div>

      {/* Edit Modal */}
      {isEditing && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" data-testid="edit-task-modal">
          <div className="bg-white rounded-lg p-6 w-full max-w-md" data-testid="edit-task-dialog">
            <h3 className="text-lg font-semibold mb-4">Görevi Düzenle</h3>
            <div className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Başlık
                </label>
                <input
                  type="text"
                  value={task.title}
                  onChange={(e) => {
                    // Basit başlık güncelleme
                    updateMutation.mutate({
                      id: task.id,
                      updates: { title: e.target.value }
                    });
                  }}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                  data-testid="input-title"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Durum
                </label>
                <select
                  value={task.status}
                  onChange={(e) => handleStatusChange(e.target.value as TaskStatus)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
                  data-testid="edit-status-select"
                >
                  <option value="beklemede">Beklemede</option>
                  <option value="devam_ediyor">Devam Ediyor</option>
                  <option value="tamamlandi">Tamamlandı</option>
                </select>
              </div>
            </div>
            <div className="flex justify-end space-x-2 mt-6">
              <button
                onClick={() => setIsEditing(false)}
                className="px-4 py-2 text-gray-700 bg-gray-200 rounded-md hover:bg-gray-300"
                data-testid="cancel-button"
              >
                İptal
              </button>
              <button
                onClick={() => setIsEditing(false)}
                className="px-4 py-2 text-white bg-primary-600 rounded-md hover:bg-primary-700"
                data-testid="save-button"
              >
                Kaydet
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Loading States */}
      {(updateMutation.isPending || deleteMutation.isPending) && (
        <div className="absolute inset-0 bg-white bg-opacity-50 flex items-center justify-center rounded-lg">
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-primary-500"></div>
        </div>
      )}
    </div>
  );
};

export default TaskCard;