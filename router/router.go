package router

import (
	"zhihu-go/internal/handler"
	"zhihu-go/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetUpRouter(r *gin.Engine, db *gorm.DB, rdb *redis.Client) *gin.Engine {
	//设置中间件
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Set("rdb", rdb)
		c.Next()
	})

	//公共路由
	public := r.Group("/api")
	{
		public.POST("/user/login", handler.Login)
		public.POST("/user/register", handler.Register)
		public.GET("/posts/search", handler.SearchPosts)
	}

	//需要认证的路由
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/user/follow", handler.Follow)
		protected.POST("/user/unfollow", handler.Unfollow)
		protected.GET("/user/followers", handler.GetFollowers)
		protected.GET("/user/followees", handler.GetFollowees)
		protected.PATCH("/user/update", handler.UpdateProfile)
		protected.POST("/post/create", handler.CreatePost)
		protected.POST("/comment", handler.CreateComment)
		protected.GET("/post/draft", handler.GetDraft)
		protected.GET("/post/published", handler.GetPublishedPost)
		protected.DELETE("/post/:id/delete", handler.DeletePost)
		protected.POST("/post/:id/restore", handler.RestorePost)
		protected.GET("/post/trash", handler.GetTrash)
		protected.PATCH("/post/:id/update", handler.UpdatePost)
		protected.POST("/like", handler.CreateLike)

		//管理员路由
		admin := protected.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/mute", handler.MuteUser)
		}
	}

	return r
}
