package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"zhihu-go/internal/handler"
)

func SetUpRouter(r *gin.Engine, db *gorm.DB) *gin.Engine {
	//设置中间件
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
	r.POST("/follow", handler.Follow)

	return r
}
