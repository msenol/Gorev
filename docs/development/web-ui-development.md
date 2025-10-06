# Web UI Development Guide

**Last Updated:** 30 September 2025
**Module:** gorev-web (v0.16.0+)
**Tech Stack:** React 18 + TypeScript + Vite + Tailwind CSS

## 📋 Overview

The Gorev Web UI is a modern React application that provides a browser-based interface for task management. It's built with TypeScript, uses Vite for fast development, and is embedded into the Go binary for zero-configuration deployment.

## 🏗️ Architecture

### Module Structure

```
gorev-web/
├── src/
│   ├── api/
│   │   └── client.ts              → React Query API client
│   ├── components/
│   │   ├── Header.tsx             → Top navigation with language switcher
│   │   ├── Sidebar.tsx            → Project list & template shortcuts
│   │   ├── TaskCard.tsx           → Individual task display
│   │   ├── ProjectSelector.tsx    → Project grid view
│   │   ├── CreateTaskModal.tsx    → Multi-step task creation wizard
│   │   └── LanguageSwitcher.tsx   → 🌍 Language toggle (TR/EN)
│   ├── contexts/
│   │   └── LanguageContext.tsx    → Language state & i18n
│   ├── types/
│   │   └── index.ts               → TypeScript type definitions
│   ├── App.tsx                    → Main application component
│   ├── main.tsx                   → Entry point with providers
│   └── index.css                  → Tailwind CSS imports
├── public/                        → Static assets
├── dist/                          → Build output (embedded in Go binary)
├── package.json                   → Dependencies
├── tsconfig.json                  → TypeScript configuration
├── vite.config.ts                 → Vite build configuration
└── tailwind.config.js             → Tailwind CSS configuration
```

### Key Technologies

- **React 18**: Component library with hooks
- **TypeScript**: Type-safe development
- **Vite**: Fast build tool with HMR
- **React Query (@tanstack/react-query)**: Server state management
- **Tailwind CSS**: Utility-first CSS framework
- **Lucide React**: Modern SVG icon library

## 🚀 Development Setup

### Prerequisites

- Node.js 18+ (LTS recommended)
- npm or yarn
- Go 1.23+ (for running API server)

### Quick Start

```bash
# 1. Navigate to web UI directory
cd gorev-web

# 2. Install dependencies
npm install

# 3. Start API server (in separate terminal)
cd ../gorev-mcpserver
./gorev serve --api-port 5082 --debug

# 4. Start development server
cd ../gorev-web
npm run dev
```

Development server will start at `http://localhost:5001` with hot reload enabled.

### Environment Variables

Create `.env.local` for local development:

```env
VITE_API_BASE_URL=http://localhost:5082
```

## 📦 Component Architecture

### 1. LanguageContext & i18n

**Purpose:** Manages UI language (Turkish/English) and syncs with MCP server

**File:** `src/contexts/LanguageContext.tsx`

```typescript
// Usage in components
import { useLanguage } from '../contexts/LanguageContext';

function MyComponent() {
  const { language, setLanguage, t } = useLanguage();

  return (
    <div>
      <h1>{t('app.title')}</h1>
      <button onClick={() => setLanguage('en')}>English</button>
    </div>
  );
}
```

**Features:**

- LocalStorage persistence
- Async API sync to MCP server (`POST /api/v1/language`)
- Translation function `t(key)` for all UI text
- Supports Turkish (tr) and English (en)

**Translation Keys:**

```typescript
// Example translations
{
  'app.title': 'Gorev Web UI',
  'projects': 'Projects',
  'tasks': 'tasks',
  'create_task': 'Create new task',
  // ... 30+ keys
}
```

### 2. API Client (React Query)

**Purpose:** Manages all HTTP requests with caching and automatic refetch

**File:** `src/api/client.ts`

```typescript
// Example: Fetch tasks
const { data: tasks, isLoading } = useQuery({
  queryKey: ['tasks', projectId],
  queryFn: () => api.get(`/api/v1/tasks?proje_id=${projectId}`)
});

// Example: Create task
const mutation = useMutation({
  mutationFn: (taskData) => api.post('/api/v1/tasks/from-template', taskData),
  onSuccess: () => {
    queryClient.invalidateQueries(['tasks']);
  }
});
```

**React Query Configuration:**

- Stale time: 30 seconds
- Cache time: 5 minutes
- Retry: 1 attempt
- Refetch on window focus

### 3. TaskCard Component

**Purpose:** Displays individual task with metadata, status, and actions

**Features:**

- Subtask count badge
- Dependency indicator (🔗 count + ⚠️ incomplete)
- Inline status dropdown
- Priority badge (color-coded)
- Context menu (⋮) for edit/delete
- Expand/collapse for long descriptions

**Props:**

```typescript
interface TaskCardProps {
  task: Task;
  onStatusChange: (taskId: string, status: string) => void;
  onDelete: (taskId: string) => void;
}
```

### 4. CreateTaskModal Component

**Purpose:** Multi-step wizard for creating tasks from templates

**Steps:**

1. Template selection (grid view with categories)
2. Field filling (dynamic form based on template)
3. Preview (review before creation)
4. Confirmation (success message + quick actions)

**State Management:**

```typescript
const [step, setStep] = useState(1);
const [selectedTemplate, setSelectedTemplate] = useState<Template | null>(null);
const [formData, setFormData] = useState<Record<string, string>>({});
```

### 5. LanguageSwitcher Component

**Purpose:** Globe icon dropdown for language selection

**Features:**

- 🌍 Globe icon (Lucide)
- Dropdown with flags (🇹🇷/🇬🇧)
- Auto-sync with MCP server
- LocalStorage persistence

**Implementation:**

```typescript
<select
  value={language}
  onChange={(e) => setLanguage(e.target.value as 'tr' | 'en')}
  className="px-3 py-2 border rounded-md"
>
  <option value="tr">🇹🇷 Türkçe</option>
  <option value="en">🇬🇧 English</option>
</select>
```

## 🔧 API Integration

### Endpoints Used

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/v1/tasks` | GET | List tasks with filters |
| `/api/v1/tasks/:id` | GET | Get task details |
| `/api/v1/tasks/:id` | PUT | Update task status |
| `/api/v1/tasks/:id` | DELETE | Delete task |
| `/api/v1/tasks/from-template` | POST | Create task from template |
| `/api/v1/projects` | GET | List all projects |
| `/api/v1/projects/:id/tasks` | GET | Get project tasks |
| `/api/v1/templates` | GET | List available templates |
| `/api/v1/language` | GET | Get current language |
| `/api/v1/language` | POST | Change language |

### Response Format

All API responses follow this structure:

```json
{
  "success": true,
  "data": { /* response data */ },
  "total": 42,
  "message": "Success message"
}
```

Error responses:

```json
{
  "success": false,
  "error": "Error message"
}
```

## 🎨 Styling Guidelines

### Tailwind CSS Classes

**Common Patterns:**

```css
/* Card */
.card { @apply bg-white rounded-lg shadow p-4 }

/* Button Primary */
.btn-primary { @apply bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 }

/* Badge */
.badge { @apply inline-block px-2 py-1 text-xs rounded }
.badge-high { @apply bg-red-100 text-red-800 }
.badge-medium { @apply bg-yellow-100 text-yellow-800 }
.badge-low { @apply bg-green-100 text-green-800 }
```

### Color Scheme

- **Primary:** Blue (#3B82F6)
- **Success:** Green (#10B981)
- **Warning:** Yellow (#F59E0B)
- **Error:** Red (#EF4444)
- **Gray Scale:** 50-900 (Tailwind defaults)

### Responsive Breakpoints

```typescript
// Mobile first approach
sm: '640px',   // Small devices
md: '768px',   // Tablets
lg: '1024px',  // Desktops
xl: '1280px',  // Large screens
```

## 🏗️ Build & Deployment

### Development Build

```bash
npm run dev
# Starts Vite dev server at http://localhost:5001
```

### Production Build

```bash
npm run build
# Output: dist/ folder
# Automatically copied to gorev-mcpserver/cmd/gorev/web/dist/
```

### Build Configuration

**vite.config.ts:**

```typescript
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5001,
    proxy: {
      '/api': 'http://localhost:5082'
    }
  },
  build: {
    outDir: '../gorev-mcpserver/cmd/gorev/web/dist',
    emptyOutDir: true
  }
})
```

### Embedding in Go Binary

The built web UI is embedded using Go's `embed` package:

```go
// cmd/gorev/web_embed.go
//go:embed web/dist/*
var webFS embed.FS
```

This allows zero-configuration deployment - users just run `npx @mehmetsenol/gorev-mcp-server serve` and the Web UI is automatically available at `http://localhost:5082`.

## 🧪 Testing

### Unit Tests (Coming Soon)

```bash
npm run test        # Run tests with Vitest
npm run test:ui     # Open Vitest UI
npm run coverage    # Generate coverage report
```

### Component Testing with Storybook (Planned)

```bash
npm run storybook   # Start Storybook dev server
npm run build-storybook  # Build static Storybook
```

## 🐛 Debugging

### React DevTools

Install React DevTools browser extension for component inspection.

### Network Debugging

Use browser DevTools Network tab to inspect API calls:

- Look for `http://localhost:5082/api/v1/*` requests
- Check request/response payloads
- Verify CORS headers

### Console Logging

Language changes are logged:

```javascript
console.log(`🌍 MCP server language changed to: ${lang}`);
```

### Common Issues

**1. API Connection Failed**

```
Error: Failed to fetch
```

**Solution:** Ensure API server is running at `http://localhost:5082`

**2. CORS Error**

```
Access to fetch at 'http://localhost:5082' from origin 'http://localhost:5001' has been blocked by CORS policy
```

**Solution:** Check CORS middleware in `internal/api/server.go`

**3. Build Errors**

```
Error: Cannot find module '@tanstack/react-query'
```

**Solution:** Run `npm install` to install dependencies

## 📚 Best Practices

### 1. Component Guidelines

- **Single Responsibility:** Each component should have one clear purpose
- **Prop Validation:** Use TypeScript interfaces for all props
- **Error Boundaries:** Wrap components in error boundaries for graceful failures
- **Loading States:** Always show loading indicators during async operations

### 2. State Management

- **Server State:** Use React Query for API data
- **UI State:** Use React hooks (useState, useReducer)
- **Global State:** Use Context API (e.g., LanguageContext)
- **Form State:** Use controlled components with validation

### 3. Performance

- **Memoization:** Use `React.memo()` for expensive components
- **Code Splitting:** Use `React.lazy()` for route-based splitting
- **Image Optimization:** Use WebP format, lazy loading
- **Bundle Size:** Keep `dist/` under 500KB gzipped

### 4. Accessibility

- **Semantic HTML:** Use proper HTML5 elements
- **ARIA Labels:** Add `aria-label` for icon buttons
- **Keyboard Navigation:** Ensure all interactive elements are keyboard accessible
- **Color Contrast:** Follow WCAG 2.1 AA standards

## 🔄 Continuous Integration

### Pre-commit Hooks (Planned)

```bash
# .husky/pre-commit
npm run lint
npm run type-check
npm run test
```

### CI/CD Pipeline (Planned)

```yaml
# .github/workflows/web-ui.yml
- name: Build Web UI
  run: |
    cd gorev-web
    npm ci
    npm run build
    npm run test
```

## 📖 Additional Resources

- [React Documentation](https://react.dev)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Vite Guide](https://vitejs.dev/guide/)
- [React Query Docs](https://tanstack.com/query/latest/docs/framework/react/overview)
- [Tailwind CSS](https://tailwindcss.com/docs)

## 🤝 Contributing

When contributing to the Web UI:

1. Follow existing component patterns
2. Add TypeScript types for all new code
3. Update this guide with new components/features
4. Test on both Chrome and Firefox
5. Ensure mobile responsiveness
6. Update translation keys for both TR and EN

---

**Questions?** Check the [main README](../../README.md) or open an issue on GitHub.
