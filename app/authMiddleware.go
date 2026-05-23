package app

import (
	"strings"

	"github.com/amrshaban2005/microservice-api/domain"
	"github.com/amrshaban2005/banking-lib/errs"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	repo domain.AuthRepository
}

func (a AuthMiddleware) authorizationHandler(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		urlParams := make(map[string]string)
		if c.Param("id") != "" {
			urlParams["customer_id"] = c.Param("id")
		}
		if c.Param("account_id") != "" {
			urlParams["account_id"] = c.Param("account_id")
		}

		if authHeader != "" {
			token := getTokenFromHeader(authHeader)
			isAuthorized := a.repo.IsAuthorize(token, permission, urlParams)
			if isAuthorized {
				c.Next()
			} else {
				appError := errs.NewAuthorizationError("Unauthorized")
				c.JSON(appError.Code, gin.H{"error": appError.AsMessage()})
				c.Abort()
				return
			}
		} else {
			appError := errs.NewAuthorizationError("missing token")
			c.JSON(appError.Code, gin.H{"error": appError.AsMessage()})
			c.Abort()
			return
		}

	}
}

func getTokenFromHeader(header string) string {
	/*
	   token is coming in the format as below
	   "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50cyI6W.yI5NTQ3MCIsIjk1NDcyIiw"
	*/
	splitToken := strings.Split(header, "Bearer")
	if len(splitToken) == 2 {
		return strings.TrimSpace(splitToken[1])
	}
	return ""
}
