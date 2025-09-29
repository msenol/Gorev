import React from 'react';
import { Globe } from 'lucide-react';
import { useLanguage } from '../contexts/LanguageContext';

const LanguageSwitcher: React.FC = () => {
  const { language, setLanguage } = useLanguage();

  return (
    <div className="flex items-center space-x-2">
      <Globe className="h-4 w-4 text-gray-400" />
      <select
        value={language}
        onChange={(e) => setLanguage(e.target.value as 'tr' | 'en')}
        className="px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary-500 bg-white"
        title="Select language"
      >
        <option value="tr">🇹🇷 Türkçe</option>
        <option value="en">🇬🇧 English</option>
      </select>
    </div>
  );
};

export default LanguageSwitcher;