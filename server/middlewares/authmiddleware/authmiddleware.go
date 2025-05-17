package authmiddleware

import (
	jwtutils "medods/utils/jwt"
	"medods/utils/snippets"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			snippets.HandleErrorJSONAnswer(c, 400, "empty access_token", "empty access_token", "[JWT]")
			c.Abort()
			return
		}

		accessClaims, err := jwtutils.ValidateAccessToken(accessToken)
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				snippets.HandleErrorJSONAnswer(c, 401, "invalid token signature", "invalid token signature", "[JWT]")
				c.Abort()
				return
			}
			snippets.HandleErrorJSONAnswer(c, 400, "token error", "token error", "[JWT]")
			c.Abort()
			return
		}

		c.Set("user_id", accessClaims.UserID)
		c.Set("session_id", accessClaims.SessionID)
		c.Next()
	}
}
