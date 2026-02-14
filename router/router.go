package router

import (
	"zhihu-go/internal/handler"
	"zhihu-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetUpRouter(r *gin.Engine, db *gorm.DB) *gin.Engine {
	//设置中间件
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	//公共路由
	public := r.Group("/api")
	{
		public.POST("/login", handler.Login)
		public.POST("/register", handler.Register)
	}

	//需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/follow", handler.Follow)
		protected.GET("/followers", handler.GetFollowers)
		protected.GET("/followees", handler.GetFollowees)
		protected.PATCH("/update", handler.UpdateProfile)
	}

	return r
}
