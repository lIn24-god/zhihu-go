package router

import (
	"zhihu-go/internal/handler"
	"zhihu-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Router 结构体，持有所有 handler
type Router struct {
	userHandler    *handler.UserHandler
	postHandler    *handler.PostHandler
	likeHandler    *handler.LikeHandler
	followHandler  *handler.FollowHandler
	commentHandler *handler.CommentHandler
}

// NewRouter 构造函数
func NewRouter(userHandler *handler.UserHandler, postHandler *handler.PostHandler,
	likeHandler *handler.LikeHandler, followHandler *handler.FollowHandler,
	commentHandler *handler.CommentHandler) *Router {
	return &Router{
		userHandler:    userHandler,
		postHandler:    postHandler,
		likeHandler:    likeHandler,
		followHandler:  followHandler,
		commentHandler: commentHandler,
	}
}

func (r *Router) SetUp(engine *gin.Engine) {
	//公共路由
	public := engine.Group("/api")
	{
		public.POST("/user/login", r.userHandler.Login)
		public.POST("/user/register", r.userHandler.Register)
		public.GET("/posts/search", r.postHandler.SearchPosts)
	}

	//需要认证的路由
	protected := engine.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/user/follow", r.followHandler.Follow)
		protected.POST("/user/unfollow", r.followHandler.Unfollow)
		protected.GET("/user/followers", r.followHandler.GetFollowers)
		protected.GET("/user/followees", r.followHandler.GetFollowees)
		protected.PATCH("/user/update", r.userHandler.UpdateProfile)
		protected.POST("/post/create", r.postHandler.CreatePost)
		protected.POST("/comment", r.commentHandler.CreateComment)
		protected.GET("/post/draft", r.postHandler.GetDraft)
		protected.GET("/post/published", r.postHandler.GetPublishedPost)
		protected.DELETE("/post/:id/delete", r.postHandler.DeletePost)
		protected.POST("/post/:id/restore", r.postHandler.RestorePost)
		protected.GET("/post/trash", r.postHandler.GetTrash)
		protected.PATCH("/post/:id/update", r.postHandler.UpdatePost)
		protected.POST("/like", r.likeHandler.CreateLike)

		//管理员路由
		/*admin := protected.Group("/admin")
		admin.Use(middleware.AdminMiddleware())
		{
			admin.POST("/mute", r.userHandler.MuteUser)
		}*/
	}
}
