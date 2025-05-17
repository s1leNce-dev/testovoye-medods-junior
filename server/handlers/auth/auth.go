package auth

import (
	"encoding/base64"
	"fmt"
	"medods/models"
	"medods/utils/encryption"
	jwtutils "medods/utils/jwt"
	"medods/utils/refresh"
	"medods/utils/snippets"
	"medods/webhook"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetTokens(c *gin.Context, db *gorm.DB) {
	var rt models.RefreshTokenSessions

	userIDQueryReq := c.Query("user_id")
	userID, err := uuid.Parse(userIDQueryReq)
	if err != nil {
		snippets.HandleErrorJSONAnswer(c, 400, "invalid user_id format", err.Error(), "[QUERY]")
		return
	}

	if err := db.Where("user_id = ?", userID).First(&rt).Error; err != nil {
		snippets.HandleErrorJSONAnswer(c, 500, "not found record", err.Error(), "[SQL]")
		return
	}

	//

	newSessionUUID := uuid.New()

	newAccessToken, err := jwtutils.AccessTokenGen(userID.String(), newSessionUUID.String())
	if err != nil {
		snippets.HandleErrorJSONAnswer(c, 500, "failed to generate access token", err.Error(), "[JWT]")
		return
	}

	newRefreshToken, err := refresh.GenerateRefreshToken()
	if err != nil {
		snippets.HandleErrorJSONAnswer(c, 500, "failed to generate refresh token", err.Error(), "[RT]")
		return
	}

	hashedRefreshToken, err := encryption.Hash(newRefreshToken)
	if err != nil {
		snippets.HandleErrorJSONAnswer(c, 500, "failed to hash refresh token", err.Error(), "[CRYPTO]")
		return
	}

	rt.SessionID = newSessionUUID
	rt.Token = hashedRefreshToken

	rt.UserID = userID

	rt.UserAgent = c.Request.UserAgent()
	rt.ClientIP = c.ClientIP()

	rt.ExpiresAt = time.Now().Add(refresh.MaxAgeRT)
	rt.CreatedAt = time.Now()
	rt.UpdatedAt = time.Now()

	if err := db.Create(&rt).Error; err != nil {
		snippets.HandleErrorJSONAnswer(c, 500, "failed to create new session", err.Error(), "[SQL]")
		return
	}

	//

	base64RefreshToken := base64.StdEncoding.EncodeToString([]byte(newRefreshToken))

	c.SetCookie(
		"refresh_token",
		base64RefreshToken,
		int(refresh.MaxAgeRT),
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(200, gin.H{
		"access_token": newAccessToken,
	})
}

func Refresh(c *gin.Context, db *gorm.DB) {
	var rt models.RefreshTokenSessions

	base64CookieRefreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		snippets.HandleErrorJSONAnswer(c, 400, "invalid refresh token", err.Error(), "[COOKIE]")
		return
	}

	if base64CookieRefreshToken == "" {
		snippets.HandleErrorJSONAnswer(c, 400, "invalid refresh token", "empty cookie", "[COOKIE]")
		return
	}

	fmt.Println("refresh_token:", base64CookieRefreshToken)

	decodedRefreshToken, err := base64.StdEncoding.DecodeString(base64CookieRefreshToken)
	if err != nil {
		snippets.HandleErrorJSONAnswer(c, 500, "failed to decode refresh_token", err.Error(), "[BASE64]")
		return
	}

	//

	accessToken := c.GetHeader("Authorization")
	if accessToken == "" {
		snippets.HandleErrorJSONAnswer(c, 400, "empty access_token", "empty access_token", "[JWT]")
		return
	}

	accessClaims, err := jwtutils.ValidateAccessTokenForRefreshToken(accessToken)
	if err != nil {
		snippets.HandleErrorJSONAnswer(c, 400, "invalid access_token", err.Error(), "[JWT]")
		return
	}

	fmt.Println("access_token:", accessClaims)

	////

	if err := db.Where("session_id = ?", accessClaims.SessionID).First(&rt).Error; err != nil {
		snippets.HandleErrorJSONAnswer(c, 500, "session not found", err.Error(), "[SQL]")
		return
	}

	if rt.ExpiresAt.Before(time.Now()) {
		snippets.HandleErrorJSONAnswer(c, 400, "refresh_token is expired", "refresh_token is expired", "[RT]")
		return
	}

	if c.Request.UserAgent() != rt.UserAgent {
		snippets.HandleErrorJSONAnswer(c, 400, "foreign user-agent", "foreign user-agent", "[RT]")
		return
	}
	if c.ClientIP() != rt.ClientIP {
		go webhook.SendWebhook(fmt.Sprintf("Unknown IP attempted to get access: %s", c.ClientIP()))
	}

	if err := encryption.VerifyHashedValue(rt.Token, string(decodedRefreshToken)); err != nil {
		snippets.HandleErrorJSONAnswer(c, 500, "failed to verify hash", err.Error(), "[CRYPTO]")
		return
	}

	////

	newSessionUUID := uuid.New()

	newAccessToken, err := jwtutils.AccessTokenGen(rt.UserID.String(), newSessionUUID.String())
	if err != nil {
		snippets.HandleErrorJSONAnswer(c, 500, "failed to generate access token", err.Error(), "[JWT]")
		return
	}

	newRefreshToken, err := refresh.GenerateRefreshToken()
	if err != nil {
		snippets.HandleErrorJSONAnswer(c, 500, "failed to generate refresh token", err.Error(), "[RT]")
		return
	}

	hashedRefreshToken, err := encryption.Hash(newRefreshToken)
	if err != nil {
		snippets.HandleErrorJSONAnswer(c, 500, "failed to hash refresh token", err.Error(), "[CRYPTO]")
		return
	}

	rt.SessionID = newSessionUUID
	rt.Token = hashedRefreshToken

	rt.UserAgent = c.Request.UserAgent()
	rt.ClientIP = c.ClientIP()

	rt.ExpiresAt = time.Now().Add(refresh.MaxAgeRT)
	rt.UpdatedAt = time.Now()

	db.Save(&rt)

	////

	base64RefreshToken := base64.StdEncoding.EncodeToString([]byte(newRefreshToken))

	c.SetCookie(
		"refresh_token",
		base64RefreshToken,
		int(refresh.MaxAgeRT),
		"/",
		"localhost",
		false,
		true,
	)

	c.JSON(200, gin.H{
		"access_token": newAccessToken,
	})
}

// user_routes
func LogOut(c *gin.Context, db *gorm.DB) {
	var rt models.RefreshTokenSessions

	sID, _ := c.Get("session_id")
	sessionID := sID.(string)

	if err := db.Where("session_id = ?", sessionID).First(&rt).Delete(&rt).Error; err != nil {
		snippets.HandleErrorJSONAnswer(c, 500, "failed to delete session", err.Error(), "[SQL]")
		return
	}

	c.JSON(200, gin.H{
		"message": "success",
	})
}
