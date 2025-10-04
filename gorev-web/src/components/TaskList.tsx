import type { Task } from '@/types';
import TaskCard from './TaskCard';

interface TaskListProps {
  tasks: Task[];
  loading: boolean;
  onTaskUpdate: () => void;
}

const TaskList: React.FC<TaskListProps> = ({
  tasks,
  loading,
  onTaskUpdate,
}) => {
  if (loading) {
    return (
      <div className="space-y-4">
        {[...Array(5)].map((_, i) => (
          <div key={i} className="card animate-pulse">
            <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
            <div className="h-3 bg-gray-200 rounded w-1/2 mb-2"></div>
            <div className="h-3 bg-gray-200 rounded w-1/4"></div>
          </div>
        ))}
      </div>
    );
  }

  // Group tasks by status
  const groupedTasks = tasks.reduce((acc, task) => {
    const status = task.durum;
    if (!acc[status]) {
      acc[status] = [];
    }
    acc[status].push(task);
    return acc;
  }, {} as Record<string, Task[]>);

  const statusConfig = {
    beklemede: {
      title: '⏳ Beklemede',
      color: 'border-yellow-200 bg-yellow-50',
    },
    devam_ediyor: {
      title: '🔄 Devam Ediyor',
      color: 'border-blue-200 bg-blue-50',
    },
    tamamlandi: {
      title: '✅ Tamamlandı',
      color: 'border-green-200 bg-green-50',
    },
  };

  const statusOrder = ['beklemede', 'devam_ediyor', 'tamamlandi'];

  return (
    <div className="space-y-6">
      {statusOrder.map((status) => {
        const statusTasks = groupedTasks[status] || [];
        if (statusTasks.length === 0) return null;

        const config = statusConfig[status as keyof typeof statusConfig];

        return (
          <div key={status} className={`border rounded-lg p-4 ${config.color}`}>
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-gray-900">
                {config.title}
              </h3>
              <span className="text-sm text-gray-600 bg-white px-2 py-1 rounded-full">
                {statusTasks.length} görev
              </span>
            </div>

            <div className="space-y-3">
              {statusTasks.map((task) => (
                <TaskCard
                  key={task.id}
                  task={task}
                  onUpdate={onTaskUpdate}
                />
              ))}
            </div>
          </div>
        );
      })}
    </div>
  );
};

export default TaskList;