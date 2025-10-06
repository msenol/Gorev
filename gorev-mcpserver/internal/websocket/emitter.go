package websocket

import (
	"time"
)

// EventEmitter interface for emitting database change events
type EventEmitter interface {
	EmitTaskCreated(workspaceID, taskID string, data map[string]interface{})
	EmitTaskUpdated(workspaceID, taskID string, data map[string]interface{})
	EmitTaskDeleted(workspaceID, taskID string)
	EmitProjectCreated(workspaceID, projectID string, data map[string]interface{})
	EmitProjectUpdated(workspaceID, projectID string, data map[string]interface{})
	EmitProjectDeleted(workspaceID, projectID string)
	EmitTemplateChanged(workspaceID string)
	EmitWorkspaceSync(workspaceID string)
}

// HubEventEmitter implements EventEmitter using the WebSocket hub
type HubEventEmitter struct {
	hub *Hub
}

// NewHubEventEmitter creates a new event emitter that broadcasts to the hub
func NewHubEventEmitter(hub *Hub) *HubEventEmitter {
	return &HubEventEmitter{hub: hub}
}

// EmitTaskCreated broadcasts a task creation event
func (e *HubEventEmitter) EmitTaskCreated(workspaceID, taskID string, data map[string]interface{}) {
	event := &ChangeEvent{
		Type:        EventTaskCreated,
		WorkspaceID: workspaceID,
		EntityID:    taskID,
		EntityType:  "task",
		Action:      "created",
		Data:        data,
		Timestamp:   time.Now().Unix(),
	}
	e.hub.BroadcastEvent(event)
}

// EmitTaskUpdated broadcasts a task update event
func (e *HubEventEmitter) EmitTaskUpdated(workspaceID, taskID string, data map[string]interface{}) {
	event := &ChangeEvent{
		Type:        EventTaskUpdated,
		WorkspaceID: workspaceID,
		EntityID:    taskID,
		EntityType:  "task",
		Action:      "updated",
		Data:        data,
		Timestamp:   time.Now().Unix(),
	}
	e.hub.BroadcastEvent(event)
}

// EmitTaskDeleted broadcasts a task deletion event
func (e *HubEventEmitter) EmitTaskDeleted(workspaceID, taskID string) {
	event := &ChangeEvent{
		Type:        EventTaskDeleted,
		WorkspaceID: workspaceID,
		EntityID:    taskID,
		EntityType:  "task",
		Action:      "deleted",
		Timestamp:   time.Now().Unix(),
	}
	e.hub.BroadcastEvent(event)
}

// EmitProjectCreated broadcasts a project creation event
func (e *HubEventEmitter) EmitProjectCreated(workspaceID, projectID string, data map[string]interface{}) {
	event := &ChangeEvent{
		Type:        EventProjectCreated,
		WorkspaceID: workspaceID,
		EntityID:    projectID,
		EntityType:  "project",
		Action:      "created",
		Data:        data,
		Timestamp:   time.Now().Unix(),
	}
	e.hub.BroadcastEvent(event)
}

// EmitProjectUpdated broadcasts a project update event
func (e *HubEventEmitter) EmitProjectUpdated(workspaceID, projectID string, data map[string]interface{}) {
	event := &ChangeEvent{
		Type:        EventProjectUpdated,
		WorkspaceID: workspaceID,
		EntityID:    projectID,
		EntityType:  "project",
		Action:      "updated",
		Data:        data,
		Timestamp:   time.Now().Unix(),
	}
	e.hub.BroadcastEvent(event)
}

// EmitProjectDeleted broadcasts a project deletion event
func (e *HubEventEmitter) EmitProjectDeleted(workspaceID, projectID string) {
	event := &ChangeEvent{
		Type:        EventProjectDeleted,
		WorkspaceID: workspaceID,
		EntityID:    projectID,
		EntityType:  "project",
		Action:      "deleted",
		Timestamp:   time.Now().Unix(),
	}
	e.hub.BroadcastEvent(event)
}

// EmitTemplateChanged broadcasts a template change event
func (e *HubEventEmitter) EmitTemplateChanged(workspaceID string) {
	event := &ChangeEvent{
		Type:        EventTemplateChanged,
		WorkspaceID: workspaceID,
		EntityType:  "template",
		Action:      "changed",
		Timestamp:   time.Now().Unix(),
	}
	e.hub.BroadcastEvent(event)
}

// EmitWorkspaceSync broadcasts a workspace sync event
func (e *HubEventEmitter) EmitWorkspaceSync(workspaceID string) {
	event := &ChangeEvent{
		Type:        EventWorkspaceSync,
		WorkspaceID: workspaceID,
		Action:      "sync",
		Timestamp:   time.Now().Unix(),
	}
	e.hub.BroadcastEvent(event)
}

// NoOpEventEmitter is a no-op implementation of EventEmitter
// Used when WebSocket support is disabled
type NoOpEventEmitter struct{}

// NewNoOpEventEmitter creates a new no-op event emitter
func NewNoOpEventEmitter() *NoOpEventEmitter {
	return &NoOpEventEmitter{}
}

func (e *NoOpEventEmitter) EmitTaskCreated(workspaceID, taskID string, data map[string]interface{})    {}
func (e *NoOpEventEmitter) EmitTaskUpdated(workspaceID, taskID string, data map[string]interface{})    {}
func (e *NoOpEventEmitter) EmitTaskDeleted(workspaceID, taskID string)                                 {}
func (e *NoOpEventEmitter) EmitProjectCreated(workspaceID, projectID string, data map[string]interface{}) {}
func (e *NoOpEventEmitter) EmitProjectUpdated(workspaceID, projectID string, data map[string]interface{}) {}
func (e *NoOpEventEmitter) EmitProjectDeleted(workspaceID, projectID string)                           {}
func (e *NoOpEventEmitter) EmitTemplateChanged(workspaceID string)                                     {}
func (e *NoOpEventEmitter) EmitWorkspaceSync(workspaceID string)                                       {}
