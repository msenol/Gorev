import React, { useState } from 'react';
import { X } from 'lucide-react';
import { useMutation } from '@tanstack/react-query';
import { createTaskFromTemplate } from '@/api/client';
import type { Template, Project } from '@/types';

interface CreateTaskModalProps {
  isOpen: boolean;
  onClose: () => void;
  templates: Template[];
  selectedProject: Project | null;
  onTaskCreated: () => void;
}

const CreateTaskModal: React.FC<CreateTaskModalProps> = ({
  isOpen,
  onClose,
  templates,
  selectedProject,
  onTaskCreated,
}) => {
  const [selectedTemplate, setSelectedTemplate] = useState<Template | null>(null);
  const [formData, setFormData] = useState<Record<string, string>>({});

  const createMutation = useMutation({
    mutationFn: createTaskFromTemplate,
    onSuccess: () => {
      onTaskCreated();
      onClose();
      setSelectedTemplate(null);
      setFormData({});
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedTemplate) return;
    if (!selectedProject) {
      alert('Lütfen önce bir proje seçin');
      return;
    }

    createMutation.mutate({
      template_id: selectedTemplate.id,
      proje_id: selectedProject.id,
      values: formData,
    });
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex min-h-screen items-center justify-center p-4">
        <div className="fixed inset-0 bg-black bg-opacity-25" onClick={onClose} />

        <div className="relative bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
          {/* Header */}
          <div className="flex items-center justify-between p-6 border-b border-gray-200">
            <h2 className="text-xl font-semibold text-gray-900">
              Yeni Görev Oluştur
            </h2>
            <button
              onClick={onClose}
              className="p-2 text-gray-400 hover:text-gray-500 hover:bg-gray-100 rounded-md"
            >
              <X className="h-5 w-5" />
            </button>
          </div>

          {/* Content */}
          <div className="p-6">
            {!selectedTemplate ? (
              /* Template Selection */
              <div>
                <h3 className="text-lg font-medium text-gray-900 mb-4">
                  Bir şablon seçin:
                </h3>
                <div className="grid grid-cols-1 gap-3">
                  {templates.map((template) => (
                    <button
                      key={template.id}
                      onClick={() => setSelectedTemplate(template)}
                      className="text-left p-4 border border-gray-200 rounded-lg hover:border-primary-300 hover:bg-primary-50 transition-colors"
                    >
                      <div className="flex items-center justify-between">
                        <h4 className="font-medium text-gray-900">
                          {template.name}
                        </h4>
                        <span className="text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded-full">
                          {template.category}
                        </span>
                      </div>
                      <p className="text-sm text-gray-600 mt-1">
                        {template.definition}
                      </p>
                    </button>
                  ))}
                </div>
              </div>
            ) : (
              /* Template Form */
              <div>
                <div className="mb-4">
                  <button
                    onClick={() => setSelectedTemplate(null)}
                    className="text-sm text-gray-500 hover:text-gray-700"
                  >
                    ← Şablon değiştir
                  </button>
                  <h3 className="text-lg font-medium text-gray-900 mt-2">
                    {selectedTemplate.name}
                  </h3>
                  <p className="text-sm text-gray-600">
                    {selectedTemplate.definition}
                  </p>
                </div>

                <form onSubmit={handleSubmit} className="space-y-4">
                  {selectedTemplate.fields.map((field) => (
                    <div key={field.name}>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        {field.name}
                        {field.required && (
                          <span className="text-red-500 ml-1">*</span>
                        )}
                      </label>

                      {field.type === 'select' ? (
                        <select
                          value={formData[field.name] || field.default || ''}
                          onChange={(e) =>
                            setFormData({ ...formData, [field.name]: e.target.value })
                          }
                          required={field.required}
                          className="input-field"
                        >
                          <option value="">Seçiniz...</option>
                          {field.options?.map((option) => (
                            <option key={option} value={option}>
                              {option}
                            </option>
                          ))}
                        </select>
                      ) : field.type === 'date' ? (
                        <input
                          type="date"
                          value={formData[field.name] || field.default || ''}
                          onChange={(e) =>
                            setFormData({ ...formData, [field.name]: e.target.value })
                          }
                          required={field.required}
                          className="input-field"
                        />
                      ) : (
                        <textarea
                          value={formData[field.name] || field.default || ''}
                          onChange={(e) =>
                            setFormData({ ...formData, [field.name]: e.target.value })
                          }
                          required={field.required}
                          rows={field.name === 'aciklama' ? 4 : 2}
                          className="input-field"
                          placeholder={field.description}
                        />
                      )}
                    </div>
                  ))}

                  {selectedProject && (
                    <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
                      <p className="text-sm text-blue-800">
                        📁 Bu görev <strong>{selectedProject.name}</strong> projesine eklenecek.
                      </p>
                    </div>
                  )}

                  <div className="flex justify-end space-x-3 pt-4">
                    <button
                      type="button"
                      onClick={onClose}
                      className="btn-secondary"
                    >
                      İptal
                    </button>
                    <button
                      type="submit"
                      disabled={createMutation.isPending}
                      className="btn-primary disabled:opacity-50"
                    >
                      {createMutation.isPending ? 'Oluşturuluyor...' : 'Görev Oluştur'}
                    </button>
                  </div>
                </form>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default CreateTaskModal;