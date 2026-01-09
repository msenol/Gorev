-- Rollback: Remove AI provider configuration and usage tracking

-- Drop prompt templates table
DROP TABLE IF EXISTS ai_prompt_templates;

-- Drop AI operations logs table
DROP TABLE IF EXISTS ai_operations;

-- Drop cached models table
DROP TABLE IF EXISTS ai_models;

-- Drop AI providers configuration table
DROP TABLE IF EXISTS ai_providers;
