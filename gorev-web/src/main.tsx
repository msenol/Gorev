import React from 'react'
import ReactDOM from 'react-dom/client'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import App from './App.tsx'
import './index.css'
import { LanguageProvider } from './contexts/LanguageContext'
import { WorkspaceProvider } from './contexts/WorkspaceContext'

// Create a client
const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5 minutes
      retry: 1,
    },
  },
})

console.log('main.tsx loaded');

const rootElement = document.getElementById('root');
console.log('Root element:', rootElement);

if (rootElement) {
  try {
    const root = ReactDOM.createRoot(rootElement);
    console.log('Creating React root with QueryClient');
    root.render(
      <React.StrictMode>
        <LanguageProvider>
          <QueryClientProvider client={queryClient}>
            <WorkspaceProvider>
              <App />
            </WorkspaceProvider>
          </QueryClientProvider>
        </LanguageProvider>
      </React.StrictMode>
    );
    console.log('React app with QueryClient mounted successfully');
  } catch (error) {
    console.error('Error mounting React app:', error);
  }
} else {
  console.error('Root element not found');
}