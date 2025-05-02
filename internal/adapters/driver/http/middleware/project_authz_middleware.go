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

type MemberType string

const (
	Admin  MemberType = "admin"
	Owner  MemberType = "owner"
	Member MemberType = "member"
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

func (m *ProjectAuthzMiddleware) Handle(memberType MemberType) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := ctx.MustGet("user").(*jwt.UserClaims)
		projectID := ctx.Param("project_id")

		projectMember, err := CheckAccess(user.ID, projectID, ctx.Request.Context(), m.projectMemberService)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, datatransfers.ResponseAbort("You are not a member of this project"))
			return
		}

		ctx.Set("project_member", projectMember)

		if !CheckRole(projectMember.Role, memberType, ctx) {

			if projectMember.TeamID != nil {
				team, err := m.teamService.GetTeamByID(ctx.Request.Context(), *projectMember.TeamID)
				if err != nil {
					ctx.AbortWithStatusJSON(http.StatusForbidden, datatransfers.ResponseAbort("You are not authorized to access this project"))
					return
				}

				ctx.Set("team", team)

				if !CheckRole(team.Role, memberType, ctx) {
					ctx.AbortWithStatusJSON(http.StatusForbidden, datatransfers.ResponseAbort("You are not authorized to access this project"))
					return
				}
			}

			ctx.AbortWithStatusJSON(http.StatusForbidden, datatransfers.ResponseAbort("You are not authorized to access this project"))

			return
		}

		ctx.Next()
	}
}

func CheckAccess(userID string, projectID string, ctx context.Context, projectMemberService ports.ProjectMemberService) (domain.ProjectMember, error) {
	result, err := projectMemberService.GetByUserIDAndProjectID(ctx, userID, projectID)
	if err != nil {
		return domain.ProjectMember{}, err
	}

	return *result, nil
}

func CheckRole(role domain.AccessRole, memberType MemberType, ctx *gin.Context) bool {
	switch {
	case role == domain.AccessOwnerRole:
		return true

	case role == domain.AccessAdminRole && memberType == Admin:
		return true

	case memberType == Member:
		isReadMethod := ctx.Request.Method == "GET"
		isWriteMethod := ctx.Request.Method == "POST" ||
			ctx.Request.Method == "PUT" ||
			ctx.Request.Method == "PATCH" ||
			ctx.Request.Method == "DELETE"

		switch {
		case (role == domain.AccessReadRole || role == domain.AccessWriteRole) && isReadMethod:
			return true
		case role == domain.AccessWriteRole && isWriteMethod:
			return true
		}
	}

	return false
}
