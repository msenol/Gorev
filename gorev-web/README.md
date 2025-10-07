# Gorev Web UI

Modern React + TypeScript web interface for Gorev task management system.

**Version:** 0.16.3 | **License:** MIT

## 📋 Overview

Gorev Web UI is a browser-based interface that provides full access to your Gorev task management system. It offers an intuitive, modern interface for managing tasks, projects, templates, and dependencies without requiring any IDE or extension installation.

## ✨ Features

### Task Management

- **Create Tasks**: Use templates to create structured tasks
- **View Tasks**: Card-based task visualization with filtering
- **Update Tasks**: Edit task details, status, and priority
- **Delete Tasks**: Remove tasks with confirmation

### Advanced Features

- **Subtask Hierarchy**: View and manage nested subtasks with expand/collapse
- **Dependencies**: Visualize task dependencies with badges
- **Project Organization**: Group tasks by project
- **Status Management**: Track task progress (Beklemede, Devam Ediyor, Tamamlandı)
- **Priority Levels**: High, medium, low priority indicators
- **Date Formatting**: Turkish locale date display
- **Real-time Updates**: React Query for automatic data synchronization

### UI Components

- **Sidebar**: Project list with task counts
- **Task Cards**: Rich task display with metadata
- **Create Modal**: Template-based task creation wizard
- **Project Selector**: Quick project switching

## 🚀 Quick Start

> **⚠️ ÖNEMLI**: Web UI artık Go binary'sine embedded olarak gömülüdür. Kullanıcıların ayrı kurulum yapmasına gerek yoktur!

### Production Kullanım (Kullanıcılar İçin)

Web UI otomatik olarak MCP sunucusuyla birlikte gelir:

```bash
# MCP sunucusunu başlat (daemon auto-starts)
npx @mehmetsenol/gorev-mcp-server@latest

# Or start daemon directly
gorev daemon --detach

# Web UI otomatik olarak şu adreste hazır:
# http://localhost:5082
```

### Development (Sadece Geliştirme İçin)

Bu bölüm sadece web UI üzerinde geliştirme yapmak isteyenler içindir:

#### Prerequisites

- Node.js 18+ and npm
- Gorev MCP server running with API enabled (port 5082)

#### Installation

```bash
# Navigate to web UI directory
cd gorev-web

# Install dependencies
npm install

# Start development server
npm run dev
```

The web UI will be available at `http://localhost:5001`

### Start API Server

The web UI requires the REST API server to be running:

```bash
# In another terminal
cd gorev-mcpserver

# Build the server
make build

# Start with API enabled
./gorev serve --api-port 5082 --debug
```

## 🔧 Available Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Start development server with hot reload (port 5001) |
| `npm run build` | Build for production (outputs to `dist/`) |
| `npm run preview` | Preview production build locally |
| `npm run lint` | Run ESLint for code quality |

## 🏗️ Architecture

### Tech Stack

- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **React Query (TanStack Query)** - Server state management
- **Tailwind CSS** - Utility-first CSS framework
- **Lucide React** - Icon library

### Project Structure

```
gorev-web/
├── src/
│   ├── api/
│   │   └── client.ts           # API client with React Query
│   ├── components/
│   │   ├── TaskCard.tsx        # Task display component
│   │   ├── Sidebar.tsx         # Project navigation
│   │   ├── CreateTaskModal.tsx # Task creation form
│   │   └── ProjectSelector.tsx # Project dropdown
│   ├── types/
│   │   └── index.ts            # TypeScript type definitions
│   ├── App.tsx                 # Main application component
│   ├── App.css                 # Custom styles
│   ├── index.css               # Tailwind imports
│   └── main.tsx                # Application entry point
├── public/                     # Static assets
├── index.html                  # HTML template
├── package.json                # Dependencies and scripts
├── tsconfig.json               # TypeScript configuration
├── vite.config.ts              # Vite configuration
└── tailwind.config.js          # Tailwind CSS configuration
```

## 🔌 API Integration

The web UI communicates with the Gorev MCP server via REST API:

### API Configuration

Default API endpoint: `http://localhost:5082/api/v1`

To change the API URL, edit `src/api/client.ts`:

```typescript
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:5082/api/v1';
```

Or set environment variable:

```bash
VITE_API_URL=http://your-server:port/api/v1 npm run dev
```

### API Endpoints Used

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/projects` | GET | Fetch all projects |
| `/projects/:id/tasks` | GET | Fetch project tasks |
| `/tasks/from-template` | POST | Create task from template |
| `/tasks/:id` | PUT | Update task |
| `/tasks/:id` | DELETE | Delete task |
| `/templates` | GET | Fetch task templates |

## 🎨 Customization

### Styling

The UI uses Tailwind CSS for styling. To customize:

1. Edit `tailwind.config.js` for theme configuration
2. Modify `src/App.css` for custom component styles
3. Update `src/index.css` for global styles

### Colors

Default color scheme uses Tailwind's primary colors. Modify in `tailwind.config.js`:

```javascript
theme: {
  extend: {
    colors: {
      primary: {
        // Your custom color palette
      }
    }
  }
}
```

## 🐛 Troubleshooting

### Common Issues

**Issue: "Failed to fetch" errors**

- Solution: Ensure MCP server is running with `--api-port 5082`
- Check CORS settings in server configuration

**Issue: Dates showing "Invalid Date"**

- Solution: Already fixed in v0.16.0 with null safety checks

**Issue: Task counts showing "NaN"**

- Solution: Already fixed in v0.16.0 with fallback values

**Issue: Port 5001 already in use**

- Solution: Change port in `vite.config.ts`:

  ```typescript
  server: {
    port: 3000, // Your preferred port
  }
  ```

### Debug Mode

Enable React Query DevTools for debugging:

```typescript
// In src/App.tsx
import { ReactQueryDevtools } from '@tanstack/react-query-devtools'

// Add to component tree
<ReactQueryDevtools initialIsOpen={false} />
```

## 📦 Production Build

### Build for Production

```bash
# Create optimized production build
npm run build

# Output directory: dist/
```

### Deployment Options

**Static Hosting** (Netlify, Vercel, GitHub Pages):

```bash
npm run build
# Deploy dist/ folder
```

**Docker**:

```dockerfile
FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build
CMD ["npm", "run", "preview"]
```

**Nginx**:

```nginx
server {
    listen 80;
    root /path/to/dist;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://localhost:5082;
    }
}
```

## 🔒 Security Considerations

- **API Authentication**: Currently no authentication (add JWT/OAuth in production)
- **CORS**: Ensure proper CORS configuration on server
- **Input Validation**: All inputs validated on server side
- **XSS Protection**: React's built-in XSS protection enabled

## 🌐 Browser Support

- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- Opera 76+

## 📝 Development Workflow

1. Start API server: `cd gorev-mcpserver && ./gorev serve --api-port 5082`
2. Start web UI: `cd gorev-web && npm run dev`
3. Make changes (hot reload enabled)
4. Test features manually
5. Build and verify: `npm run build && npm run preview`

## 🤝 Contributing

Contributions welcome! Please see [CONTRIBUTING.md](../CONTRIBUTING.md) for guidelines.

### Development Guidelines

- Follow React best practices
- Use TypeScript for type safety
- Write semantic HTML
- Keep components small and focused
- Use React Query for all API calls
- Follow existing code style

## 📄 License

MIT License - See [LICENSE](../LICENSE) for details

## 🔗 Related Documentation

- [Main README](../README.md) - Project overview
- [MCP Tools Reference](../docs/legacy/tr/mcp-araclari.md) - MCP tool documentation
- [API Reference](../docs/api/rest-api-reference.md) - REST API documentation (coming soon)
- [Development Guide](../docs/development/TASKS.md) - Development roadmap

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/msenol/gorev/issues)
- **Discussions**: [GitHub Discussions](https://github.com/msenol/gorev/discussions)
- **Documentation**: [docs/](../docs/)

---

**Built with ❤️ using React, TypeScript, and Vite**
