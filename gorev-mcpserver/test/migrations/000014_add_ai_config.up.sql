-- Migration: Add AI provider configuration and usage tracking
-- This enables multi-provider AI support (OpenRouter, Anannas) with project-level configuration

-- AI provider configurations per project
-- Stores encrypted API keys and model preferences for each project
CREATE TABLE IF NOT EXISTS ai_providers (
    id TEXT PRIMARY KEY,
    project_id TEXT NOT NULL,
    provider TEXT NOT NULL CHECK(provider IN ('openrouter', 'anannas')),
    api_key_encrypted TEXT NOT NULL,
    model TEXT NOT NULL,
    temperature REAL DEFAULT 0.7,
    max_tokens INTEGER DEFAULT 2000,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    workspace_id TEXT NOT NULL DEFAULT 'default',
    FOREIGN KEY (project_id) REFERENCES projeler(id) ON DELETE CASCADE,
    UNIQUE(project_id, workspace_id)
);

-- Indexes for efficient lookups
CREATE INDEX IF NOT EXISTS idx_ai_providers_project ON ai_providers(project_id);
CREATE INDEX IF NOT EXISTS idx_ai_providers_workspace ON ai_providers(workspace_id);
CREATE INDEX IF NOT EXISTS idx_ai_providers_provider ON ai_providers(provider);

-- Cached model catalogs from providers
-- Avoids frequent API calls to list available models
CREATE TABLE IF NOT EXISTS ai_models (
    id TEXT PRIMARY KEY,
    provider TEXT NOT NULL,
    model_id TEXT NOT NULL,
    model_name TEXT NOT NULL,
    context_window INTEGER,
    input_price REAL DEFAULT 0,
    output_price REAL DEFAULT 0,
    cached_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, model_id)
);

-- Index for provider-specific model lookups
CREATE INDEX IF NOT EXISTS idx_ai_models_provider ON ai_models(provider);
CREATE INDEX IF NOT EXISTS idx_ai_models_cached_at ON ai_models(cached_at);

-- AI operation logs for cost tracking and audit
-- Tracks all AI API calls with token usage and costs
CREATE TABLE IF NOT EXISTS ai_operations (
    id TEXT PRIMARY KEY,
    project_id TEXT,
    operation TEXT NOT NULL CHECK(operation IN (
        'task_creation',
        'task_decomposition',
        'semantic_search',
        'time_estimation',
        'file_association',
        'project_analysis',
        'prioritization',
        'nl_command',
        'chat',
        'suggest'
    )),
    provider TEXT NOT NULL,
    model TEXT NOT NULL,
    input_tokens INTEGER DEFAULT 0,
    output_tokens INTEGER DEFAULT 0,
    total_cost REAL DEFAULT 0,
    duration_ms INTEGER,
    success BOOLEAN NOT NULL DEFAULT TRUE,
    error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    workspace_id TEXT NOT NULL DEFAULT 'default',
    FOREIGN KEY (project_id) REFERENCES projeler(id) ON DELETE SET NULL
);

-- Indexes for cost tracking and analytics
CREATE INDEX IF NOT EXISTS idx_ai_operations_project ON ai_operations(project_id);
CREATE INDEX IF NOT EXISTS idx_ai_operations_created ON ai_operations(created_at);
CREATE INDEX IF NOT EXISTS idx_ai_operations_provider ON ai_operations(provider);
CREATE INDEX IF NOT EXISTS idx_ai_operations_workspace ON ai_operations(workspace_id);
CREATE INDEX IF NOT EXISTS idx_ai_operations_success ON ai_operations(success);

-- Composite index for monthly cost summaries
CREATE INDEX IF NOT EXISTS idx_ai_operations_project_created ON ai_operations(project_id, created_at);

-- Prompt templates for AI operations
-- Stores system and user prompt templates with i18n support
CREATE TABLE IF NOT EXISTS ai_prompt_templates (
    id TEXT PRIMARY KEY,
    operation TEXT NOT NULL,
    language_code TEXT NOT NULL DEFAULT 'tr' CHECK(language_code IN ('tr', 'en')),
    template_name TEXT NOT NULL,
    system_prompt TEXT NOT NULL,
    user_prompt_template TEXT NOT NULL,
    variables TEXT, -- JSON array of variable names
    active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(operation, language_code, template_name)
);

-- Index for prompt lookups
CREATE INDEX IF NOT EXISTS idx_ai_prompts_operation ON ai_prompt_templates(operation, language_code);
CREATE INDEX IF NOT EXISTS idx_ai_prompts_active ON ai_prompt_templates(active);

-- Seed default prompt templates
-- Turkish prompts
INSERT INTO ai_prompt_templates (id, operation, language_code, template_name, system_prompt, user_prompt_template, variables) VALUES
('prompt_task_creation_tr', 'task_creation', 'tr', 'default',
 'Sen görev yönetimi konusunda uzman bir yapay zeka asistanısın. Kullanıcının doğal dil girdisini analiz edip yapılandırılmış bir görev oluştur.',
 'Kullanıcı girdisi: {{input}}

Mevcut görevler:
{{existing_tasks}}

Yeni görev oluştur. JSON formatında döndür:
{
  "title": "Görev başlığı",
  "description": "Görev açıklaması",
  "priority": "dusuk|orta|yuksek",
  "tags": ["etiket1", "etiket2"],
  "template": "bug|feature|research|refactor|test|doc",
  "estimated_hours": 1.5
}',
 '["input", "existing_tasks"]'),

('prompt_task_decomposition_tr', 'task_decomposition', 'tr', 'default',
 'Sen karmaşık görevleri daha küçük, yönetilebilir alt görevlere ayırma konusunda uzmansın. Her alt görev bağımsız olarak tamamlanabilmeli.',
 'Ana görev:
Başlık: {{title}}
Açıklama: {{description}}

Bu görevi mantıksal alt görevlere ayır. JSON formatında döndür:
{
  "subtasks": [
    {
      "title": "Alt görev başlığı",
      "description": "Alt görev açıklaması",
      "estimated_hours": 1.0,
      "dependencies": []
    }
  ]
}',
 '["title", "description"]'),

('prompt_semantic_search_tr', 'semantic_search', 'tr', 'default',
 'Sen görevler arasında semantik arama yapma konusunda uzmansın. Kullanıcının sorgusunu anlayıp ilgili görevleri bul.',
 'Kullanıcı sorgusu: {{query}}

Mevcut görevler:
{{tasks}}

JSON formatında döndür:
{
  "results": [
    {
      "task_id": "gorev_123",
      "relevance_score": 0.95,
      "reason": "Eşleşme nedeni"
    }
  ]
}',
 '["query", "tasks"]'),

('prompt_time_estimation_tr', 'time_estimation', 'tr', 'default',
 'Sen görev süreleri tahmin etme konusunda uzmansın. Geçmiş verilere ve görev karmaşıklığına bakarak gerçekçi tahminler ver.',
 'Görev:
Başlık: {{title}}
Açıklama: {{description}}
Etiketler: {{tags}}

Benzer geçmiş görevler:
{{historical_tasks}}

JSON formatında döndür:
{
  "estimated_hours": 4.5,
  "confidence": 0.8,
  "reasoning": "Tahmin nedeni"
}',
 '["title", "description", "tags", "historical_tasks"]'),

('prompt_project_analysis_tr', 'project_analytics', 'tr', 'default',
 'Sen proje analizi ve risk değerlendirmesi konusunda uzmansın. Projenin durumunu analiz edip önemli içgörüler sun.',
 'Proje: {{project_name}}
Görev sayısı: {{task_count}}
Tamamlanan: {{completed_count}}
Devam eden: {{in_progress_count}}
Bekleyen: {{pending_count}}
Gecikmiş: {{overdue_count}}

JSON formatında döndür:
{
  "health_score": 0.75,
  "risks": ["Risk 1", "Risk 2"],
  "recommendations": ["Öneri 1", "Öneri 2"],
  "critical_path": ["task_id1", "task_id2"]
}',
 '["project_name", "task_count", "completed_count", "in_progress_count", "pending_count", "overdue_count"]');

-- English prompts
INSERT INTO ai_prompt_templates (id, operation, language_code, template_name, system_prompt, user_prompt_template, variables) VALUES
('prompt_task_creation_en', 'task_creation', 'en', 'default',
 'You are an expert task management AI assistant. Analyze user input and create a structured task.',
 'User input: {{input}}

Existing tasks:
{{existing_tasks}}

Create a new task. Return in JSON format:
{
  "title": "Task title",
  "description": "Task description",
  "priority": "low|medium|high",
  "tags": ["tag1", "tag2"],
  "template": "bug|feature|research|refactor|test|doc",
  "estimated_hours": 1.5
}',
 '["input", "existing_tasks"]'),

('prompt_task_decomposition_en', 'task_decomposition', 'en', 'default',
 'You are an expert at breaking down complex tasks into smaller, manageable subtasks. Each subtask should be independently completable.',
 'Main task:
Title: {{title}}
Description: {{description}}

Break this task into logical subtasks. Return in JSON format:
{
  "subtasks": [
    {
      "title": "Subtask title",
      "description": "Subtask description",
      "estimated_hours": 1.0,
      "dependencies": []
    }
  ]
}',
 '["title", "description"]'),

('prompt_semantic_search_en', 'semantic_search', 'en', 'default',
 'You are an expert at semantic search across tasks. Understand user queries and find relevant tasks.',
 'User query: {{query}}

Available tasks:
{{tasks}}

Return in JSON format:
{
  "results": [
    {
      "task_id": "gorev_123",
      "relevance_score": 0.95,
      "reason": "Reason for match"
    }
  ]
}',
 '["query", "tasks"]'),

('prompt_time_estimation_en', 'time_estimation', 'en', 'default',
 'You are an expert at estimating task durations. Provide realistic estimates based on historical data and task complexity.',
 'Task:
Title: {{title}}
Description: {{description}}
Tags: {{tags}}

Similar historical tasks:
{{historical_tasks}}

Return in JSON format:
{
  "estimated_hours": 4.5,
  "confidence": 0.8,
  "reasoning": "Reason for estimate"
}',
 '["title", "description", "tags", "historical_tasks"]'),

('prompt_project_analysis_en', 'project_analytics', 'en', 'default',
 'You are an expert at project analysis and risk assessment. Analyze project status and provide key insights.',
 'Project: {{project_name}}
Task count: {{task_count}}
Completed: {{completed_count}}
In progress: {{in_progress_count}}
Pending: {{pending_count}}
Overdue: {{overdue_count}}

Return in JSON format:
{
  "health_score": 0.75,
  "risks": ["Risk 1", "Risk 2"],
  "recommendations": ["Recommendation 1", "Recommendation 2"],
  "critical_path": ["task_id1", "task_id2"]
}',
 '["project_name", "task_count", "completed_count", "in_progress_count", "pending_count", "overdue_count"]');
