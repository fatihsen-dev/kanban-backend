package middlewares

import (
	"net/http"
	"strings"

	"github.com/fatihsen-dev/kanban-backend/internal/adapters/driver/http/datatransfers"
	"github.com/fatihsen-dev/kanban-backend/pkg/jwt"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	isAdmin bool
}

func NewAuthMiddleware(isAdmin bool) *AuthMiddleware {
	return (&AuthMiddleware{
		isAdmin: isAdmin,
	})
}

func (m *AuthMiddleware) Handle(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.ResponseAbort("missing authorization header"))
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.ResponseAbort("invalid header format"))
		return
	}

	if headerParts[0] != "Bearer" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.ResponseAbort("token must content bearer"))
		return
	}

	user, err := jwt.VerifyToken(headerParts[1])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.ResponseAbort("invalid token"))
		return
	}

	if user.IsAdmin != m.isAdmin && !user.IsAdmin {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, datatransfers.ResponseAbort("you don't have access for this action"))
		return
	}

	ctx.Set("user", user)
	ctx.Next()
}
