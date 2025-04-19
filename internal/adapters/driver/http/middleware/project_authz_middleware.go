package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers"
	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/validation"
	ports "github.com/fatihsen-dev/kanban-backend/internal/core/ports/driver"
	"github.com/fatihsen-dev/kanban-backend/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type MemberType string

const (
	Admin MemberType = "admin"
	Owner MemberType = "owner"
)

type AuthzMiddleware struct {
	projectMemberService ports.ProjectMemberService
	teamService          ports.TeamService
}

func NewAuthzMiddleware(projectMemberService ports.ProjectMemberService, teamService ports.TeamService) *AuthzMiddleware {
	return (&AuthzMiddleware{
		projectMemberService: projectMemberService,
		teamService:          teamService,
	})
}

func (m *AuthzMiddleware) Handle(ctx *gin.Context, projectIdSource string, fieldName string, memberType MemberType) {
	user := ctx.MustGet("user").(*jwt.UserClaims)
	projectID, err := m.getProjectID(ctx, projectIdSource, fieldName)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datatransfers.ResponseAbort(err.Error()))
		return
	}

	fmt.Println("memberType", projectID, memberType)

	if user.IsAdmin {
		ctx.Next()
		return
	}

	ctx.Next()
}

func (m *AuthzMiddleware) getProjectID(ctx *gin.Context, projectIdSource string, fieldName string) (string, error) {
	projectID := ""

	switch projectIdSource {
	case "query":
		switch fieldName {
		case "project_id":
			projectID = ctx.Query(fieldName)
		case "id":
			projectID = ctx.Query(fieldName)
		default:
			return "", errors.New("invalid project id")
		}
	case "param":
		switch fieldName {
		case "project_id":
			projectID = ctx.Param(fieldName)
		case "id":
			projectID = ctx.Param(fieldName)
		default:
			return "", errors.New("invalid project id")
		}
	case "body":
		switch fieldName {
		case "project_id":
			var requestData struct {
				ProjectID string `json:"project_id"`
			}
			if err := ctx.ShouldBindJSON(&requestData); err != nil {
				return "", errors.New("invalid project id")
			}
			projectID = requestData.ProjectID
		case "id":
			var requestData struct {
				ID string `json:"id"`
			}
			if err := ctx.ShouldBindJSON(&requestData); err != nil {
				return "", errors.New("invalid project id")
			}
			projectID = requestData.ID
		default:
			return "", errors.New("invalid project id")
		}
	}

	if projectID == "" {
		return "", errors.New("project id is required")
	}

	if err := validation.ValidateUUID(projectID); err != nil {
		return "", errors.New("invalid project id")
	}

	return projectID, nil
}
