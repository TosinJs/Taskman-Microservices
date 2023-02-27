package authMiddleware

import (
	"net/http"
	"strings"
	"tosinjs/reminder-service/internal/entity/responseEntity"
	"tosinjs/reminder-service/internal/service/authService"

	"github.com/gin-gonic/gin"
)

type authMiddleware struct {
	authSVC authService.AuthService
}

func New(authSVC authService.AuthService) authMiddleware {
	return authMiddleware{
		authSVC: authSVC,
	}
}

func (a authMiddleware) VerifyJWT() gin.HandlerFunc {
	fn := func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responseEntity.BuildErrorResponseObject(
				"Unauthorized", c.FullPath(),
			))
			return
		}

		authArray := strings.Split(authHeader, " ")
		if len(authArray) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, responseEntity.BuildErrorResponseObject(
				"Unauthorized", c.FullPath(),
			))
			return
		}

		authPayload, svcErr := a.authSVC.ValidateJWT(authArray[1])

		if svcErr != nil {
			c.AbortWithStatusJSON(svcErr.StatusCode, responseEntity.BuildServiceErrorResponseObject(
				svcErr, c.FullPath(),
			))
			return
		}

		c.Set("userId", authPayload.Id)
	}
	return fn
}
