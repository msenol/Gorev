// Package ai provides centralized prompt management for AI operations
// This file contains all prompt templates in one place (DRY compliance)
package ai

import (
	"context"
	"fmt"
)

// PromptManager manages AI prompts with i18n support
type PromptManager struct {
	// Future: Load from database or embedded filesystem
}

// NewPromptManager creates a new prompt manager
func NewPromptManager() *PromptManager {
	return &PromptManager{}
}

// GetPrompt returns the system and user prompts for a given operation and language
func (pm *PromptManager) GetPrompt(ctx context.Context, operation AIOperation, lang string) (system, user string, err error) {
	// Default to Turkish if no language specified
	if lang == "" {
		lang = "tr"
	}

	// Get prompts based on operation and language
	switch operation {
	case AIOperationChat:
		return pm.getChatPrompt(lang)
	case AIOperationTaskCreation:
		return pm.getTaskCreationPrompt(lang)
	case AIOperationTaskDecomposition:
		return pm.getTaskDecompositionPrompt(lang)
	case AIOperationSemanticSearch:
		return pm.getSemanticSearchPrompt(lang)
	case AIOperationTimeEstimation:
		return pm.getTimeEstimationPrompt(lang)
	case AIOperationProjectAnalytics:
		return pm.getProjectAnalyticsPrompt(lang)
	case AIOperationPrioritization:
		return pm.getPrioritizationPrompt(lang)
	case AIOperationNLCommand:
		return pm.getNLCommandPrompt(lang)
	default:
		return "", "", fmt.Errorf("unknown operation: %s", operation)
	}
}

// getChatPrompt returns prompts for general chat
func (pm *PromptManager) getChatPrompt(lang string) (system, user string, err error) {
	if lang == "tr" {
		system = "Sen Gorev görev yönetim sistemi için yapay zeka asistanısın. Kullanıcının görevlerini yönetmesine, organize olmasına ve verimli kalmasına yardımcı ol."
	} else {
		system = "You are an AI assistant for the Gorev task management system. Help users manage, organize, and stay productive with their tasks."
	}
	return system, "", nil
}

// getTaskCreationPrompt returns prompts for creating tasks from natural language
func (pm *PromptManager) getTaskCreationPrompt(lang string) (system, user string, err error) {
	if lang == "tr" {
		system = "Sen doğal dil girdilerinden yapılandırılmış görevler oluşturma konusunda uzmansın. Kullanıcının girdisini analiz et ve template, öncelik, etiketler gibi uygun özelliklerle görev oluştur."
	} else {
		system = "You are an expert at creating structured tasks from natural language input. Analyze the user's input and create tasks with appropriate properties like template, priority, and tags."
	}
	return system, "", nil
}

// getTaskDecompositionPrompt returns prompts for breaking down tasks
func (pm *PromptManager) getTaskDecompositionPrompt(lang string) (system, user string, err error) {
	if lang == "tr" {
		system = "Sen karmaşık görevleri daha küçük, yönetilebilir alt görevlere ayırma konusunda uzmansın. Her alt görev bağımsız olarak tamamlanabilmeli. JSON formatında döndür."
	} else {
		system = "You are an expert at breaking down complex tasks into smaller, manageable subtasks. Each subtask should be completable independently. Return in JSON format."
	}
	return system, "", nil
}

// getSemanticSearchPrompt returns prompts for semantic search
func (pm *PromptManager) getSemanticSearchPrompt(lang string) (system, user string, err error) {
	if lang == "tr" {
		system = "Sen görevlerin anlamsal aranması konusunda uzmansın. Kullanıcının sorgusunu anlamalı ve ilgili görevleri bulmalısın."
	} else {
		system = "You are an expert at semantic search of tasks. Understand the user's query and find relevant tasks."
	}
	return system, "", nil
}

// getTimeEstimationPrompt returns prompts for estimating task duration
func (pm *PromptManager) getTimeEstimationPrompt(lang string) (system, user string, err error) {
	if lang == "tr" {
		system = "Sen görev süreleri tahmin etme konusunda uzmansın. Geçmiş verilere ve görev karmaşıklığına bakarak gerçekçi tahminler ver. JSON formatında döndür."
	} else {
		system = "You are an expert at estimating task durations. Provide realistic estimates based on historical data and task complexity. Return in JSON format."
	}
	return system, "", nil
}

// getProjectAnalyticsPrompt returns prompts for project analysis
func (pm *PromptManager) getProjectAnalyticsPrompt(lang string) (system, user string, err error) {
	if lang == "tr" {
		system = "Sen proje analizi konusunda uzmansın. Kritik yol analizi, risk değerlendirmesi ve kaynak tahsisi yapabilirsin. JSON formatında döndür."
	} else {
		system = "You are an expert at project analytics. Perform critical path analysis, risk assessment, and resource allocation. Return in JSON format."
	}
	return system, "", nil
}

// getPrioritizationPrompt returns prompts for task prioritization
func (pm *PromptManager) getPrioritizationPrompt(lang string) (system, user string, err error) {
	if lang == "tr" {
		system = "Sen görev önceliklendirme konusunda uzmansın. Görevleri önem, aciliyet ve bağımlılıklara göre önceliklendir."
	} else {
		system = "You are an expert at task prioritization. Prioritize tasks based on importance, urgency, and dependencies."
	}
	return system, "", nil
}

// getNLCommandPrompt returns prompts for natural language command parsing
func (pm *PromptManager) getNLCommandPrompt(lang string) (system, user string, err error) {
	if lang == "tr" {
		system = "Sen doğal dil komutlarını yapılandırılmış işlem komutlarına dönüştürme konusunda uzmansın."
	} else {
		system = "You are an expert at converting natural language commands into structured operation commands."
	}
	return system, "", nil
}

// BuildTaskDecompositionPrompt builds the user prompt for task decomposition
func (pm *PromptManager) BuildTaskDecompositionPrompt(lang string, title, description string) string {
	if lang == "tr" {
		return fmt.Sprintf(`Ana görev:
Başlık: %s
Açıklama: %s

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
}`, title, description)
	}
	return fmt.Sprintf(`Main task:
Title: %s
Description: %s

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
}`, title, description)
}

// BuildTimeEstimationPrompt builds the user prompt for time estimation
func (pm *PromptManager) BuildTimeEstimationPrompt(lang string, title, description string, tags []string) string {
	tagsStr := ""
	if len(tags) > 0 {
		tagsStr = fmt.Sprintf("%v", tags)
	}

	if lang == "tr" {
		return fmt.Sprintf(`Görev:
Başlık: %s
Açıklama: %s
Etiketler: %s

JSON formatında döndür:
{
  "estimated_hours": 4.5,
  "confidence": 0.8,
  "reasoning": "Tahmin nedeni"
}`, title, description, tagsStr)
	}
	return fmt.Sprintf(`Task:
Title: %s
Description: %s
Tags: %s

Return in JSON format:
{
  "estimated_hours": 4.5,
  "confidence": 0.8,
  "reasoning": "Reasoning for estimate"
}`, title, description, tagsStr)
}

// BuildProjectAnalyticsPrompt builds the user prompt for project analytics
func (pm *PromptManager) BuildProjectAnalyticsPrompt(lang string, projectName string, taskCount, completedCount, pendingCount int) string {
	if lang == "tr" {
		return fmt.Sprintf(`Proje Analizi: %s
Toplam Görev: %d
Tamamlanan: %d
Bekleyen: %d

JSON formatında döndür:
{
  "critical_path": ["task_id_1", "task_id_2"],
  "risk_level": "high",
  "risk_factors": ["factor1", "factor2"],
  "recommendations": ["recommendation1"],
  "bottlenecks": ["bottleneck1"]
}`, projectName, taskCount, completedCount, pendingCount)
	}
	return fmt.Sprintf(`Project Analytics: %s
Total Tasks: %d
Completed: %d
Pending: %d

Return in JSON format:
{
  "critical_path": ["task_id_1", "task_id_2"],
  "risk_level": "high",
  "risk_factors": ["factor1", "factor2"],
  "recommendations": ["recommendation1"],
  "bottlenecks": ["bottleneck1"]
}`, projectName, taskCount, completedCount, pendingCount)
}
