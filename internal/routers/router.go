package routers

import (
	"net/http"

	"configuration-management/internal/routers/private/v1/apps"

	"configuration-management/global"
	"configuration-management/internal/routers/private/v1/card"
	"configuration-management/internal/routers/private/v1/configuration"
	"configuration-management/internal/routers/private/v1/user"
	"configuration-management/pkg/app"
	"configuration-management/pkg/logger"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())

	// Swagger Docs
	r.GET("/TIPCRFNFJJ/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Private router
	privateGroup := r.Group("/private/v1")
	privateGroup.Use(jwtAuthMiddleware())

	// Public router
	publicGroup := r.Group("/public/v1")
	publicGroup.Use(corsMiddleware())

	{
		// Configuration
		configHandler := configuration.NewHandler()
		privateGroup.GET("/configuration", configHandler.GetConfiguration)
		privateGroup.POST("/configuration", configHandler.CreateConfiguration)
	}

	{
		// Card
		cardHandler := card.NewHandler()

		// public
		cardPublicGroup := publicGroup.Group("")
		cardPublicGroup.POST("/identity", cardHandler.Identity)
		cardPublicGroup.POST("/activate", cardHandler.Activate)

		// private
		privateGroup.GET("/card/:value", cardHandler.GetCardByValue)
		privateGroup.GET("/cards", cardHandler.GetCards)
		privateGroup.GET("/export-cards", cardHandler.Export)
		privateGroup.POST("/card", cardHandler.CreateCard)
		privateGroup.PUT("/card", cardHandler.UpdateCard)
		privateGroup.DELETE("/card/:value", cardHandler.DeleteCard)
		privateGroup.POST("/cards", cardHandler.CreateCards)
		privateGroup.DELETE("/cards", cardHandler.DeleteCardsByValues)
		privateGroup.PUT("/batch-update-card-status", cardHandler.BatchUpdateStatus)
		privateGroup.GET("/batch-query", cardHandler.BatchQuery)
		privateGroup.GET("/get-card-count-by-status", cardHandler.GetCardCountByStatus)
		privateGroup.PUT("/set-expired-at", cardHandler.SetCardExpiredAt)
	}

	{
		// App
		appsHandler := apps.NewHandler()
		privateGroup.GET("/apps", appsHandler.QueryAppList)
		privateGroup.POST("/app", appsHandler.CreateApp)
		privateGroup.PUT("/app", appsHandler.UpdateApp)
		privateGroup.GET("/app-options", appsHandler.QueryAppOptions)
	}

	{
		// User
		userHandler := user.NewHandler()
		// public
		publicGroup.POST("/login", userHandler.Login)
		publicGroup.POST("/logout", userHandler.Logout)

		// private
		privateGroup.GET("/get-user-info", userHandler.GetUserInfo)
		userPrivateGroup := privateGroup.Group("")
		userPrivateGroup.Use(userManageAuthMiddleware())
		userPrivateGroup.GET("/users", userHandler.QueryUserList)
		userPrivateGroup.POST("/user", userHandler.CreateUser)
		userPrivateGroup.PUT("/user", userHandler.UpdateUser)
		userPrivateGroup.POST("/reset-password", userHandler.ResetPassword)
	}

	return r
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, V-Token")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Type, Content-Length, Authorization, V-Token")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func jwtAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("V-Token")
		if token == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, app.ResponseContent{
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		// 检查 token 是否已经失效
		// TODO: 从数据库中获取已经失效的 token，写入 InvalidTokenCache
		if _, ok := global.InvalidTokenCache.Get(token); ok {
			global.Logger.WithFields(logger.Fields{
				"token": token,
			}).Error("token is invalid")
			context.AbortWithStatusJSON(http.StatusUnauthorized, app.ResponseContent{
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		userInfo, err := app.GetUserInfoFromToken(token)
		if err != nil {
			global.Logger.WithFields(logger.Fields{
				"token": token,
			}).Error("get user info from token failed", err)
			context.AbortWithStatusJSON(http.StatusUnauthorized, app.ResponseContent{
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		// 将 userInfo 信息写入上下文
		context.Set(app.UserInfoKey, userInfo)
	}
}

// 用户管理权限校验中间件
func userManageAuthMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		userInfo := app.GetUserInfoFromContext(context)
		if userInfo.UserId == "" {
			context.AbortWithStatusJSON(http.StatusUnauthorized, app.ResponseContent{
				StatusCode: http.StatusUnauthorized,
			})
			return
		}

		// 检查是否有权限管理用户
		if !userInfo.IsRoot() {
			global.Logger.WithFields(logger.Fields{
				"user_info": userInfo,
			}).Error("没有权限管理用户")
			context.AbortWithStatusJSON(http.StatusUnauthorized, app.ResponseContent{
				StatusCode: http.StatusUnauthorized,
			})
			return
		}
	}
}
