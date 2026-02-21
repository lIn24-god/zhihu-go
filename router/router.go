package router

import (
	"zhihu-go/internal/handler"
	"zhihu-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Router 结构体，持有所有 handler
type Router struct {
	userHandler *handler.UserHandler
}

// NewRouter 构造函数
func NewRouter(userHandler *handler.UserHandler) *Router {
	return &Router{
		userHandler: userHandler,
	}
}

func (r *Router) SetUp(engine *gin.Engine) {
	//公共路由
	public := engine.Group("/api")
	{
		public.POST("/user/login", r.userHandler.Login)
		public.POST("/user/register", r.userHandler.Register)
		public.GET("/posts/search", handler.SearchPosts)
	}

	//需要认证的路由
	protected := engine.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/user/follow", handler.Follow)
		protected.POST("/user/unfollow", handler.Unfollow)
		protected.GET("/user/followers", handler.GetFollowers)
		protected.GET("/user/followees", handler.GetFollowees)
		protected.PATCH("/user/update", r.userHandler.UpdateProfile)
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
		/*admin := protected.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/mute", r.userHandler.MuteUser)
		}*/
	}
}
