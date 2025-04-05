package ws

import (
	"encoding/json"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/validation"
)

func HandleMessage(hub *Hub, client *Client, message []byte) {
	var msg Message

	json.Unmarshal(message, &msg)
	err := validation.ValidateMessage(msg)
	if err != nil {
		client.send <- NewErrorResponse(err.Error())
		return
	}

	var groupID string

	switch msg.Namespace {
	case NamespaceProjects:
		if projectID, ok := msg.Data["project_id"].(string); ok {
			groupID = projectID
			hub.AddToGroup(projectID, client)
		}
		switch msg.Action {
		case ActionCreateProject:
			CreateProject(hub, msg.Data, groupID)
		}
	case NamespaceTasks:
		if projectID, ok := msg.Data["project_id"].(string); ok {
			groupID = projectID
			hub.AddToGroup(projectID, client)
		}
		switch msg.Action {
		case ActionCreateTask:
			CreateTask(hub, msg.Data, groupID)
		}
	case NamespaceColumns:
		if projectID, ok := msg.Data["project_id"].(string); ok {
			groupID = projectID
			hub.AddToGroup(projectID, client)
		}
		switch msg.Action {
		case ActionCreateColumn:
			CreateColumn(hub, msg.Data, groupID)
		}
	default:
		client.send <- NewAbortResponse("unknown namespace")
	}
}

func CreateProject(hub *Hub, data map[string]interface{}, groupID string) {
	// your logic here
	hub.Broadcast(groupID, NewSuccessResponse("success", data))
}

func CreateTask(hub *Hub, data map[string]interface{}, groupID string) {
	// your logic here
	hub.Broadcast(groupID, NewSuccessResponse("success", data))
}

func CreateColumn(hub *Hub, data map[string]interface{}, groupID string) {
	// your logic here
	hub.Broadcast(groupID, NewSuccessResponse("success", data))
}
