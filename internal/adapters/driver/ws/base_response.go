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
	EventNameProjectMemberCreated EventName = "project.member.created"
	EventNameUserStatusUpdated    EventName = "user.status.updated"
	EventNameTeamMembersAdded     EventName = "team.members.added"
	EventNameProjectMemberUpdated EventName = "project.member.updated"
	EventNameProjectDeleted       EventName = "project.deleted"
)

type BaseResponse struct {
	Name EventName   `json:"name"`
	Data interface{} `json:"data"`
}
