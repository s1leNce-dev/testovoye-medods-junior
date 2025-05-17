package routes

import (
	"medods/handlers/auth"
	"medods/handlers/user"
	"medods/middlewares/authmiddleware"
	"medods/middlewares/cors"

	_ "medods/docs"

	"github.com/gin-gonic/gin"
	swagFiles "github.com/swaggo/files"
	ginSwag "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

const api_v1 = "/api/v1"

func InitRoutes(r *gin.Engine, db *gorm.DB) {
	cors.InitCors(r)

	r.GET(api_v1+"/token", func(c *gin.Context) {
		GetTokensRoute(c, db)
	})

	r.GET(api_v1+"/refresh", func(c *gin.Context) {
		RefreshRoute(c, db)
	})

	userRoutes := r.Group(api_v1 + "/user")
	userRoutes.Use(authmiddleware.AuthMiddleware())
	{
		userRoutes.GET("/uuid", GetUUIDRoute)
		userRoutes.GET("/logout", func(c *gin.Context) {
			LogOutRoute(c, db)
		})
	}

	r.GET("/swagger/*any", ginSwag.WrapHandler(swagFiles.Handler))
}

// GetTokensRoute godoc
// @Summary      Получить access и refresh токены
// @Description  Проверяем, есть ли у этого пользователя
// @Tags         Auth
// @Param        user_id query string true "UUID пользователя"
// @Produce      json
// @Success      200 {object} docsmodels.TokenResponse
// @Failure      400 {object} docsmodels.ErrorResponse
// @Failure      500 {object} docsmodels.ErrorResponse
// @Router       /api/v1/token [get]
func GetTokensRoute(c *gin.Context, db *gorm.DB) {
	auth.GetTokens(c, db)
}

// RefreshRoute godoc
// @Summary      Обновить access токен по refresh токену
// @Description  Access_token - Authorization Header; Refresh_token - Cookie (http-only);
// @Tags         Auth
// @Produce      json
// @Success      200 {object} docsmodels.TokenResponse
// @Failure      400 {object} docsmodels.ErrorResponse
// @Failure      500 {object} docsmodels.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/refresh [get]
func RefreshRoute(c *gin.Context, db *gorm.DB) {
	auth.Refresh(c, db)
}

// GetUUIDRoute godoc
// @Summary      Получить UUID авторизованного пользователя
// @Description  Возвращает текущий user_id из middleware
// @Tags         User
// @Produce      json
// @Success      200 {object} docsmodels.UUIDResponse
// @Failure      401 {object} docsmodels.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/user/uuid [get]
func GetUUIDRoute(c *gin.Context) {
	user.GetUUID(c)
}

// LogOutRoute godoc
// @Summary      Выход пользователя
// @Description  Удаляет refresh-сессию из БД и чистит куку
// @Tags         User
// @Produce      json
// @Success      200 {object} docsmodels.SuccessMessage
// @Failure      500 {object} docsmodels.ErrorResponse
// @Security     BearerAuth
// @Router       /api/v1/user/logout [get]
func LogOutRoute(c *gin.Context, db *gorm.DB) {
	auth.LogOut(c, db)
}
