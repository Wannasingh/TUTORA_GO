package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

func (h *HttpHandler) ToggleCommentLike(c *gin.Context) {
	idStr := c.Param("id")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	requesterID, _ := c.Get("userID")
	liked, err := h.postUsecase.ToggleCommentLike(c.Request.Context(), commentID, requesterID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"liked": liked})
}

func (h *HttpHandler) DeleteComment(c *gin.Context) {
	idStr := c.Param("id")
	commentID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment id"})
		return
	}

	requesterID, _ := c.Get("userID")
	err = h.postUsecase.DeleteComment(c.Request.Context(), requesterID.(int), commentID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted successfully"})
}

func (h *HttpHandler) ReportContent(c *gin.Context) {
	requesterID, _ := c.Get("userID")

	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   int    `json:"target_id" binding:"required"`
		Reason     string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	report := &domain.Report{
		ReporterID: requesterID.(int),
		TargetType: req.TargetType,
		TargetID:   req.TargetID,
		Reason:     req.Reason,
	}

	err := h.postUsecase.ReportContent(c.Request.Context(), report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, report)
}

func (h *HttpHandler) UpdatePost(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req domain.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &domain.Post{
		ID:       postID,
		UserID:   userID.(int),
		Subject:  req.Subject,
		Title:    req.Title,
		Body:     req.Body,
		ImageURL: req.ImageURL,
		VideoURL: req.VideoURL,
	}

	err = h.postUsecase.UpdatePost(c.Request.Context(), userID.(int), post)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *HttpHandler) DeletePost(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	err = h.postUsecase.DeletePost(c.Request.Context(), userID.(int), postID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "post deleted successfully"})
}
