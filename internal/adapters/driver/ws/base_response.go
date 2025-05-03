package ws

type EventName string

const (
	EventNameColumnCreated        EventName = "column.created"
	EventNameColumnUpdated        EventName = "column.updated"
	EventNameColumnDeleted        EventName = "column.deleted"
	EventNameTaskCreated          EventName = "task.created"
	EventNameTaskUpdated          EventName = "task.updated"
	EventNameTaskDeleted          EventName = "task.deleted"
	EventNameTaskMoved            EventName = "task.moved"
	EventNameTeamCreated          EventName = "team.created"
	EventNameTeamUpdated          EventName = "team.updated"
	EventNameTeamDeleted          EventName = "team.deleted"
	EventNameInvitationCreated    EventName = "invitation.created"
	EventNameProjectMemberCreated EventName = "project_member.created"
	EventNameUserStatusUpdated    EventName = "user.status.updated"
)

type BaseResponse struct {
	Name EventName   `json:"name"`
	Data interface{} `json:"data"`
}
