import React, { createContext, useContext, useState } from 'react';

type Language = 'tr' | 'en';

interface LanguageContextType {
  language: Language;
  setLanguage: (lang: Language) => void;
  t: (key: string) => string;
}

const translations: Record<Language, Record<string, string>> = {
  tr: {
    'app.title': 'Gorev Web UI',
    'projects': 'Projeler',
    'all_projects': 'Tüm Projeler',
    'templates': 'Şablonlar',
    'create_task': 'Yeni görev oluştur',
    'search_tasks': 'Görevlerde ara...',
    'all_statuses': 'Tüm Durumlar',
    'all_priorities': 'Tüm Öncelikler',
    'pending': 'Beklemede',
    'in_progress': 'Devam Ediyor',
    'completed': 'Tamamlandı',
    'high': 'Yüksek',
    'medium': 'Orta',
    'low': 'Düşük',
    'tasks': 'görev',
    'select_project': 'Bir proje seçin',
    'no_tasks': 'Henüz görev yok',
    'wait_templates': 'Görev oluşturmak için önce template\'lerin yüklenmesini bekleyin.',
    'cancel': 'İptal',
    'create': 'Görev Oluştur',
    'select_template': 'Bir şablon seçin:',
    'change_template': '← Şablon değiştir',
    'task_will_be_added': 'Bu görev',
    'to_project': 'projesine eklenecek.',
    'required': 'zorunlu',
    'select': 'Seçiniz...',
    'creating': 'Oluşturuluyor...',
    'new_task': 'Yeni Görev Oluştur',
    'menu': 'Menü',
    'edit': 'Düzenle',
    'delete': 'Sil',
    'delete_confirm': 'Bu görevi silmek istediğinizden emin misiniz?',
    'more_templates': 'daha fazla şablon...',
    'no_project_tasks': 'Bu projede görev yok',
    'no_tasks_yet': 'Henüz görev yok',
    'select_template_to_create': 'Sol taraftan bir template seçerek yeni görev oluşturabilirsiniz.',
    'wait_templates_to_load': 'Görev oluşturmak için önce template\'lerin yüklenmesini bekleyin.',
  },
  en: {
    'app.title': 'Gorev Web UI',
    'projects': 'Projects',
    'all_projects': 'All Projects',
    'templates': 'Templates',
    'create_task': 'Create new task',
    'search_tasks': 'Search tasks...',
    'all_statuses': 'All Statuses',
    'all_priorities': 'All Priorities',
    'pending': 'Pending',
    'in_progress': 'In Progress',
    'completed': 'Completed',
    'high': 'High',
    'medium': 'Medium',
    'low': 'Low',
    'tasks': 'tasks',
    'select_project': 'Select a project',
    'no_tasks': 'No tasks yet',
    'wait_templates': 'Wait for templates to load before creating tasks.',
    'cancel': 'Cancel',
    'create': 'Create Task',
    'select_template': 'Select a template:',
    'change_template': '← Change template',
    'task_will_be_added': 'This task will be added to',
    'to_project': 'project.',
    'required': 'required',
    'select': 'Select...',
    'creating': 'Creating...',
    'new_task': 'Create New Task',
    'menu': 'Menu',
    'edit': 'Edit',
    'delete': 'Delete',
    'delete_confirm': 'Are you sure you want to delete this task?',
    'more_templates': 'more templates...',
    'no_project_tasks': 'No tasks in this project',
    'no_tasks_yet': 'No tasks yet',
    'select_template_to_create': 'Select a template from the left sidebar to create a new task.',
    'wait_templates_to_load': 'Wait for templates to load before creating tasks.',
  },
};

const LanguageContext = createContext<LanguageContextType | undefined>(undefined);

export const LanguageProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [language, setLanguageState] = useState<Language>(() => {
    const saved = localStorage.getItem('gorev_language');
    return (saved === 'en' || saved === 'tr') ? saved : 'tr';
  });

  const setLanguage = async (lang: Language) => {
    setLanguageState(lang);
    localStorage.setItem('gorev_language', lang);

    // Sync language with MCP server
    try {
      await fetch('http://localhost:5082/api/v1/language', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ language: lang }),
      });
      console.log(`🌍 MCP server language changed to: ${lang}`);
    } catch (error) {
      console.warn('Failed to sync language with MCP server:', error);
      // Don't fail the UI language change if API call fails
    }
  };

  const t = (key: string): string => {
    return translations[language][key] || key;
  };

  return (
    <LanguageContext.Provider value={{ language, setLanguage, t }}>
      {children}
    </LanguageContext.Provider>
  );
};

export const useLanguage = () => {
  const context = useContext(LanguageContext);
  if (!context) {
    throw new Error('useLanguage must be used within LanguageProvider');
  }
  return context;
};