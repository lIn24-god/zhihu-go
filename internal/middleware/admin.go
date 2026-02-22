package middleware

import (
	"net/http"
	"zhihu-go/internal/service"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware 检查当前用户是否为管理员
func AdminMiddleware(userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		uintUserID, ok := userID.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid"})
			c.Abort()
			return
		}

		isAdmin, err := userService.IsAdmin(uintUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
			c.Abort()
			return
		}

		if !isAdmin {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		c.Next()
	}
}
