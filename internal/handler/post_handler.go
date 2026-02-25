package handler

import (
	"net/http"
	"zhihu-go/internal/service"
	"zhihu-go/pkg/response"

	"strconv"

	"github.com/gin-gonic/gin"

	"zhihu-go/internal/dto"
)

// PostHandler 结构体定义
type PostHandler struct {
	postService service.PostService
	userService service.UserService
}

// NewPostHandler 构造函数
func NewPostHandler(postService service.PostService, userService service.UserService) *PostHandler {
	return &PostHandler{
		postService: postService,
		userService: userService,
	}
}

// CreatePost 文章发布
func (h *PostHandler) CreatePost(c *gin.Context) {
	var request dto.PostRequest

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// 从 gin.Context 获取请求的 context
	ctx := c.Request.Context()

	post, err := h.postService.CreatePost(ctx, uintUserID, &request)
	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"post_id": post.ID,
		"message": "Post created successfully",
	})
}

// GetDraft 获取用户草稿
func (h *PostHandler) GetDraft(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	ctx := c.Request.Context()

	resp, err := h.postService.GetPost(ctx, uintUserID, "draft")
	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// GetPublishedPost 获取用户已发布文章
func (h *PostHandler) GetPublishedPost(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	ctx := c.Request.Context()

	resp, err := h.postService.GetPost(ctx, uintUserID, "published")
	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, resp)
}

// SearchPosts 处理文章搜索请求
func (h *PostHandler) SearchPosts(c *gin.Context) {
	//获取查询参数
	keyword := c.Query("q")
	if keyword == "" {
		response.Error(c, http.StatusBadRequest, "missing keyword")
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

	ctx := c.Request.Context()

	results, total, err := h.postService.SearchPosts(ctx, keyword, page, pageSize)
	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"data":     results,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

// DeletePost 删除文章
func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	//从url中获取postID
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid post id")
		return
	}

	ctx := c.Request.Context()

	if err := h.postService.SoftDeletePost(ctx, uint(postID), uintUserID); err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "delete successfully"})
}

// RestorePost 恢复文章
func (h *PostHandler) RestorePost(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	//从url中获取postID
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid post id")
		return
	}

	ctx := c.Request.Context()

	if err := h.postService.RestorePost(ctx, uint(postID), uintUserID); err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "restore successfully"})
}

// GetTrash 获取用户回收站中的文章
func (h *PostHandler) GetTrash(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	ctx := c.Request.Context()

	trash, err := h.postService.GetUserTrash(ctx, uintUserID)
	if err != nil {
		HandleError(c, err)
		return
	}

	var resp []dto.PostResponse
	for _, f := range trash {
		resp = append(resp, dto.PostResponse{
			Title:    f.Title,
			AuthorID: f.AuthorID,
			Content:  f.Content,
		})
	}
	response.Success(c, resp)
}

// UpdatePost 修改文章信息
func (h *PostHandler) UpdatePost(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	//从url中获取postID
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid post id")
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	// 从 gin.Context 获取请求的 context
	ctx := c.Request.Context()

	if err := h.postService.UpdatePost(ctx, uintUserID, uint(postID), req); err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "update successfully"})
}
