package ws

type Message struct {
	Namespace string                 `json:"namespace" validate:"required,min=1"`
	Action    string                 `json:"action" validate:"required,min=1"`
	Data      map[string]interface{} `json:"data" validate:"required,dive,required"`
}

const (
	// Namespace - Tasks
	NamespaceTasks = "tasks"
	// Actions - Tasks
	ActionCreateTask = "create_task"
	ActionUpdateTask = "update_task"

	// Namespace - Projects
	NamespaceProjects = "projects"
	// Actions - Projects
	ActionCreateProject = "create_project"
	ActionUpdateProject = "update_project"

	// Namespace - Columns
	NamespaceColumns = "columns"
	// Actions - Columns
	ActionCreateColumn = "create_column"
	ActionUpdateColumn = "update_column"
)
