package handler

import (
	"net/http"
	"zhihu-go/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"strconv"

	"zhihu-go/internal/dto"
)

func CreatePost(c *gin.Context) {
	var request dto.PostRequest

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	if err := service.CreatePost(db, &request, uintUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"massage": "Create post successfully"})
}

// GetDraft 获取用户草稿
func GetDraft(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	result, err := service.GetPost(db, uintUserID, "draft")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get draft"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": result})
}

// GetPublishedPost 获取用户已发布文章
func GetPublishedPost(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	result, err := service.GetPost(db, uintUserID, "published")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get published post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": result})
}

// SearchPosts 处理文章搜索请求
func SearchPosts(c *gin.Context) {
	//获取查询参数
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing keyword"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	//参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	db := c.MustGet("db").(*gorm.DB)

	results, total, err := service.SearchPosts(db, keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     results,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}
