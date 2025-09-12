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
**Description:** Full-text search with natural language support, multi-filter combinations, and saved filter profiles.

**Features:**
- SQLite FTS5 extension
- Natural language query support (integrates with AI Context Management)
- Fuzzy search support
- Filter profiles (save/load)
- Command Palette integration
- Search history
- Smart suggestions based on AI interactions

#### 3. 🔧 Task Dependencies in TreeView
**Status:** Completed (5 July 2025)  
**Description:** Visualize dependencies within TreeView instead of separate dependency graph.

**Implemented Features:**
- ✅ Enhanced visual indicators with progress bars, priority badges, and due date formatting
- ✅ Smart dependency badges (🔒 locked, 🔓 unlocked, 🔗 blocking)
- ✅ Rich markdown tooltips with progress visualization
- ✅ Tag pill badges with configurable display
- ✅ Configuration options for all visual preferences
- ✅ Created TaskDecorationProvider for managing decorations
- ✅ Improved task descriptions with separator formatting
- ✅ Smart due date formatting (Today, Tomorrow, 3d left, etc.)

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

#### 5. 🤖 AI Context Management & Automation System
**Status:** Pending  
**Description:** Smart context management and automation system optimized for AI users, addressing the critical issue of 0% active tasks and improving AI-human collaboration.

**Core Problems Addressed:**
- **77.5% tasks stuck in "beklemede"** - No automatic state transitions
- **0% tasks in "devam_ediyor"** - AI forgets to update task states
- **Context loss in long conversations** - No persistent task context
- **Manual everything** - No automation or smart suggestions

##### Phase 1: Context Management (Priority 1)
- **New MCP Tools:**
  - `gorev_set_active(task_id)` - Set active task context
  - `gorev_get_active()` - Get current active task
  - `gorev_recent(limit=5)` - Get recent task interactions
  - `gorev_context_summary()` - AI-optimized session summary
- **Database Updates:**
  - New table: `ai_context` (active task, recent tasks, session data)
  - New table: `ai_interactions` (track AI-task interactions)
  - Add columns: `last_ai_interaction`, `estimated_hours`, `actual_hours`
- **Auto-State Management:**
  - Auto-transition to "devam_ediyor" when task accessed
  - Auto-transition to "beklemede" after 30min inactivity
  - Parent task completion check when subtasks complete
  - Git commit integration for status updates

##### Phase 2: Natural Language & Batch Operations
- **NLP Query Interface:**
  - `gorev_nlp_query("bugün üzerinde çalıştığım görevler")`
  - Relative references: "son oluşturduğum görev", "database görevleri"
  - Smart interpretation: "yarın yapalım" → son_tarih: tomorrow
- **Batch Operations:**
  - `gorev_batch_update([{id: "123", status: "completed"}, ...])`
  - `gorev_bulk_create([task1, task2, ...])` 
  - `gorev_bulk_transition(ids[], new_status)`
- **AI Summary Dashboard:**
  - Current sprint overview with blockers
  - Suggested next actions based on priorities
  - Time tracking: estimated vs actual

##### Phase 3: Smart Automation
- **Intelligent Task Creation:**
  - Auto-split into subtasks based on description
  - ML-based time estimation
  - Similar task detection and template suggestion
- **Integration Hooks:**
  - File edit → task progress update
  - Code review → task creation from TODOs
  - PR merge → task completion
- **Predictive Analytics:**
  - Task completion predictions
  - Risk analysis (deadline miss probability)
  - Bottleneck detection

**Implementation Timeline:**
- Phase 1 (2 weeks): Context & auto-state management
- Phase 2 (3 weeks): NLP queries & batch operations  
- Phase 3 (4 weeks): Automation & predictions

**Success Metrics:**
- Active task ratio: >20% (from current 0%)
- Context switches: <5 per AI session
- Bulk operations: >30% of all updates
- Auto-state accuracy: >80%

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
- AI Context Management & Automation System (Phase 1 & 2)
- Advanced Search and Filtering (with NLP support)
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