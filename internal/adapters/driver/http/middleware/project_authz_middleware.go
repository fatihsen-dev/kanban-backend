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

type AuthzMiddleware struct {
	projectIdSource string
	fieldName       string
	teamService     ports.TeamService
}

func NewAuthzMiddleware(projectIdSource string, fieldName string, teamService ports.TeamService) *AuthzMiddleware {
	return (&AuthzMiddleware{
		projectIdSource: projectIdSource,
		fieldName:       fieldName,
		teamService:     teamService,
	})
}

func (m *AuthzMiddleware) Handle(ctx *gin.Context) {
	user := ctx.MustGet("user").(*jwt.UserClaims)
	projectID, err := m.getProjectID(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, datatransfers.ResponseAbort(err.Error()))
		return
	}

	fmt.Println("projectID", projectID)

	if user.IsAdmin {
		ctx.Next()
		return
	}

	ctx.Next()
}

func (m *AuthzMiddleware) getProjectID(ctx *gin.Context) (string, error) {
	projectID := ""

	switch m.projectIdSource {
	case "query":
		switch m.fieldName {
		case "project_id":
			projectID = ctx.Query(m.fieldName)
		case "id":
			projectID = ctx.Param(m.fieldName)
		default:
			return "", errors.New("invalid project id")
		}
	case "body":
		switch m.fieldName {
		case "project_id":
			var requestData struct {
				ProjectID string `json:"project_id"`
			}
			if err := ctx.ShouldBindJSON(&requestData); err != nil {
				return "", errors.New("invalid project id")
			}
			projectID = requestData.ProjectID
		case "id":
			projectID = ctx.Param(m.fieldName)
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
