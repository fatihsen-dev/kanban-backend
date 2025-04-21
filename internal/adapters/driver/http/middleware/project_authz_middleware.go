package middlewares

import (
	"context"
	"net/http"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers"
	"github.com/fatihsen-dev/kanban-backend/internal/core/domain"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driver"
	"github.com/fatihsen-dev/kanban-backend/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type AccessType string

const (
	Admin  AccessType = "admin"
	Owner  AccessType = "owner"
	Member AccessType = "member"
)

type ProjectAuthzMiddleware struct {
	projectMemberService ports.ProjectMemberService
	teamService          ports.TeamService
}

func NewProjectAuthzMiddleware(projectMemberService ports.ProjectMemberService, teamService ports.TeamService) *ProjectAuthzMiddleware {
	return &ProjectAuthzMiddleware{
		projectMemberService: projectMemberService,
		teamService:          teamService,
	}
}

func (m *ProjectAuthzMiddleware) Handle(accessType AccessType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*jwt.UserClaims)
		projectID := ctx.Param("project_id")

		projectMember, err := CheckMemberAccess(user.ID, projectID, ctx.Request.Context(), m.projectMemberService)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, datatransfers.ResponseAbort("You are not a member of this project"))
			return
		}

		if !CheckMemberRole(projectMember, accessType, ctx) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, datatransfers.ResponseAbort("You are not authorized to access this project"))
			return
		}

		ctx.Next()
	}
}

func CheckMemberAccess(userID string, projectID string, ctx context.Context, projectMemberService ports.ProjectMemberService) (domain.ProjectMember, error) {
	result, err := projectMemberService.GetByUserIDAndProjectID(ctx, userID, projectID)
	if err != nil {
		return domain.ProjectMember{}, err
	}

	return *result, nil
}

func CheckMemberRole(projectMember domain.ProjectMember, accessType AccessType, ctx *gin.Context) bool {
	if projectMember.Role == domain.ProjectMemberRole(Owner) {
		return true
	}

	if projectMember.Role == domain.ProjectAdminRole && accessType == Admin {
		return true
	}

	if accessType == Member && (projectMember.Role == domain.ProjectReadRole && ctx.Request.Method == "GET") {
		return true
	}

	if accessType == Member && (projectMember.Role == domain.ProjectWriteRole && (ctx.Request.Method == "POST" || ctx.Request.Method == "PUT" || ctx.Request.Method == "PATCH")) {
		return true
	}

	return false
}
