# 🗺️ Gorev Roadmap

This roadmap outlines the planned development path for the Gorev task management system.

## 🔨 Active Development Tasks

### High Priority

#### 1. 🔍 Test Coverage Improvement (Target 95%) - Geçici olarak durdurldu.
**Status:** Pending  
**Description:** Expand test suite for both modules to achieve 95% coverage and implement E2E testing.

**Research Areas:**
- Go testing best practices (testify vs ginkgo)
- VS Code extension testing (Jest vs Mocha vs Vitest)
- E2E testing tools (Playwright vs Cypress)
- Performance testing (k6, Gatling)

**Targets:**
- Backend (Go): 95+ coverage
- Frontend (TypeScript): 90+ coverage
- E2E: 100% critical user journeys

#### 2. ✨ Advanced Search and Filtering System
**Status:** Pending  
**Description:** Full-text search, multi-filter combinations, and saved filter profiles.

**Features:**
- SQLite FTS5 extension
- Fuzzy search support
- Filter profiles (save/load)
- Command Palette integration
- Search history

#### 3. 🔧 Task Dependencies in TreeView
**Status:** Pending  
**Description:** Visualize dependencies within TreeView instead of separate dependency graph.

**Solution:**
- Badge system (🔗 emoji, counter)
- Tooltips showing dependent tasks
- Visual link indicators
- Drag & drop dependency creation

#### 4. ✅ Subtask System with Unlimited Hierarchy
**Status:** Completed (30 June 2025) - v0.8.0  
**Description:** Unlimited depth hierarchical task structure.

**Implemented Features:**
- ✅ Added `parent_id` column with foreign key constraints
- ✅ Recursive CTE queries for efficient hierarchy operations
- ✅ Hierarchical display in task listing with tree structure
- ✅ Circular dependency prevention with validation
- ✅ Parent task progress tracking based on subtask completion
- ✅ Business rules: deletion prevention, completion validation, project inheritance
- ✅ New MCP tools: `gorev_altgorev_olustur`, `gorev_ust_degistir`, `gorev_hiyerarsi_goster`
- ✅ Comprehensive test coverage for all subtask operations
- 🔲 Collapsible hierarchy in VS Code TreeView (future enhancement)

#### 5. 🤖 AI-Powered Task Enrichment System
**Status:** Pending  
**Description:** Intelligent system that enriches tasks with current technology information using user-provided API keys.

**Technical Architecture:**

##### API Key Management
- Secure storage using VS Code Secret Storage API
- Supported services: OpenAI, Anthropic, Context7, GitHub APIs
- Secure API key input in Settings UI

##### MCP Server Updates
- New MCP tool: `gorev_zenginlestir`
- New table: `gorev_zenginlestirmeleri` (cache and history)
- NLP technology extraction
- Smart caching (TTL: 7 days)

##### VS Code Extension Updates
- New settings for AI enablement and auto-enrichment
- API key management panel with connection testing
- "🤖 AI Enrich" button in task details
- Progress indicators and notifications

##### Enrichment Types
1. **Technology Updates**: Framework versions, best practices, migration guides
2. **Resources**: Documentation links, tutorials, Stack Overflow solutions
3. **Risk Warnings**: Deprecated APIs, security vulnerabilities, breaking changes
4. **Alternative Suggestions**: Similar technologies, modern approaches, community recommendations

##### Smart Features
- Contextual understanding with NLP
- Progressive enhancement levels
- Team preference learning
- Offline mode support

**Implementation Phases:**
- Phase 1 (2 weeks): Basic infrastructure, API key management
- Phase 2 (3 weeks): AI integration, NLP analysis
- Phase 3 (2 weeks): Auto-enrichment, analytics

**Expected Benefits:**
- 40% faster task planning
- Current technology tracking
- Early risk detection
- Standardized best practices

### Medium Priority

#### 6. ✨ DevOps Pipeline and Automated Release Process
**Status:** Pending  
**Description:** Comprehensive CI/CD pipeline with quality control, automated testing, and release processes.

**Components:**
- Code quality checks (golangci-lint, ESLint, gosec, Snyk)
- Automated testing (unit, integration, E2E)
- Multi-platform builds (Linux, macOS, Windows)
- Release automation (semantic versioning, changelog)
- VS Code Marketplace deployment
- Docker Hub, Homebrew, Scoop distribution

**Acceptance Criteria:**
- Automated code quality control for every PR
- Cross-platform builds under 15 minutes
- Rollback capability under 5 minutes

#### 7. ✅ Fix Filter State Persistence Issue
**Status:** Completed (30 June 2025)  
**Description:** Users cannot clear filters once applied - requires VS Code restart to reset.

**Problem:**
- When filters are applied to task list, there's no reliable way to clear them
- "Clear Filters" button or command may be missing or non-functional
- Filter state persists across sessions inappropriately
- Workaround: VS Code restart required (`Developer: Reload Window`)

**Solution Requirements:**
- Add visible "Clear All Filters" button to filter toolbar
- Implement `Gorev: Clear Filters` command in Command Palette
- Ensure filter state is properly reset in TreeView provider
- Add keyboard shortcut (e.g., `Ctrl+Alt+R`) for quick filter reset
- Fix filter persistence logic in workspace settings

#### 8. 🔧 VS Code Extension UI/UX Improvements and Accessibility
**Status:** Pending  
**Description:** Improve user experience and make extension accessible to all users.

**Improvements:**
- ARIA attributes and screen reader support
- Keyboard navigation (Tab order, shortcuts)
- WCAG 2.1 AA compliance
- Dark/High contrast theme support
- Loading states and error handling
- Touch gesture support

#### 9. 🔧 Performance Optimizations and Scalability
**Status:** Pending  
**Description:** Performance optimization for large datasets (10K+ tasks scenarios).

**Optimizations:**
- Database indexing (composite indexes)
- In-memory LRU cache
- Lazy loading & pagination
- Virtual scrolling TreeView
- Query optimization (prepared statements)

#### 10. ✨ External Service Integrations
**Status:** Pending  
**Description:** GitHub, Jira, Slack integrations.

**GitHub:** Issue sync, PR connections, webhooks  
**Jira:** Issue import/export, status sync  
**Slack:** Notifications, slash commands, interactive messages

#### 11. ✨ Task Statistics and Reporting Dashboard
**Status:** Pending  
**Description:** Analytics and visualization features.

**Dashboard Components:**
- Project progress metrics
- Burndown/burnup charts
- Time analysis and trends
- Interactive charts (Chart.js/D3.js)
- PDF/Excel export

### Low Priority

#### 12. ✨ Multi-user System and Authorization Infrastructure
**Status:** Pending  
**Description:** Transform Gorev into a multi-user system with authentication and authorization.

**Features:**
- JWT token-based authentication
- OAuth 2.0 provider support (Google, GitHub, Microsoft)
- RBAC (Admin, Project Manager, Team Member, Guest)
- Team and project membership system
- Audit trail and compliance

## 🚀 Long-term Goals (v1.0.0)

### MCP Server
- **Multi-user Support**: User management and authorization
- **Cloud Sync**: Cloud synchronization
- **API Gateway**: REST/GraphQL API
- **Plugin System**: Extensible architecture
- **AI Integration**: Task suggestions and automatic categorization

### VS Code Extension
- **Collaboration Features**: Real-time collaboration
- **Mobile Companion App**: Mobile application
- **Voice Commands**: Voice control
- **AI Assistant**: Task management assistant
- **Custom Themes**: Customizable themes

### Ecosystem
- **CLI Tool**: Standalone CLI application
- **Web Dashboard**: Web-based management panel
- **Browser Extension**: Chrome/Firefox extensions
- **Integrations**: Jira, GitHub, GitLab, Trello integrations
- **API SDK**: JavaScript, Python, Go SDKs

## 📅 Timeline

### Q3 2025 (July - September)
- AI-Powered Task Enrichment System
- Advanced Search and Filtering
- Test Coverage Improvement

### Q4 2025 (October - December)
- Performance Optimizations
- DevOps Pipeline
- External Service Integrations

### Q1 2026 (January - March)
- Multi-user System
- Task Statistics Dashboard
- UI/UX Improvements

### Q2 2026 (April - June)
- ~~Subtask System~~ ✅ Completed early (June 2025)
- Dependency Visualization
- Preparation for v1.0.0

## 🎯 Success Metrics

- **User Adoption**: 1K+ active users by end of 2025
- **Performance**: Support for 10K+ tasks per user
- **Quality**: 95%+ test coverage maintained
- **Ecosystem**: 5+ third-party integrations
- **Community**: 100+ GitHub stars, active contributor base

## 🤝 Contributing

We welcome contributions! Check our [development guide](docs/development/) for details on:
- Setting up development environment
- Testing procedures
- Code style guidelines
- Pull request process

## 📬 Feedback

Your feedback shapes our roadmap! Please share:
- Feature requests at [GitHub Issues](https://github.com/msenol/gorev/issues)
- Use cases and success stories
- Integration needs
- Performance feedback

---

*This roadmap is a living document and may be adjusted based on user feedback, technical discoveries, and changing priorities.*