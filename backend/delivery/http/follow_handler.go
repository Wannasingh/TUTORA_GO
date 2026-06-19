package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type FollowHandler struct {
	followUsecase domain.FollowUsecase
	reviewUsecase domain.ReviewUsecase
}

func NewFollowHandler(fu domain.FollowUsecase, ru domain.ReviewUsecase) *FollowHandler {
	return &FollowHandler{followUsecase: fu, reviewUsecase: ru}
}

func (h *FollowHandler) ToggleFollow(c *gin.Context) {
	userID, _ := c.Get("userID")
	targetIDStr := c.Param("id")
	targetID, err := strconv.Atoi(targetIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	followed, err := h.followUsecase.ToggleFollow(c.Request.Context(), userID.(int), targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"following": followed})
}

func (h *FollowHandler) GetFollowers(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	users, err := h.followUsecase.ListFollowers(c.Request.Context(), targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *FollowHandler) GetFollowing(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	users, err := h.followUsecase.ListFollowing(c.Request.Context(), targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *FollowHandler) GetFollowStats(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	requesterID := 0
	if rid, exists := c.Get("userID"); exists {
		requesterID = rid.(int)
	}

	stats, err := h.followUsecase.GetFollowStats(c.Request.Context(), targetID, requesterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// ============ REVIEWS ============

func (h *FollowHandler) SubmitReview(c *gin.Context) {
	userID, _ := c.Get("userID")
	tutorID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tutor id"})
		return
	}

	var req struct {
		Rating float64 `json:"rating" binding:"required"`
		Body   *string `json:"body,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review := &domain.TutorReview{
		ReviewerID: userID.(int),
		TutorID:    tutorID,
		Rating:     req.Rating,
		Body:       req.Body,
	}

	if err := h.reviewUsecase.SubmitReview(c.Request.Context(), review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func (h *FollowHandler) GetTutorReviews(c *gin.Context) {
	tutorID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tutor id"})
		return
	}

	reviews, err := h.reviewUsecase.GetTutorReviews(c.Request.Context(), tutorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (h *FollowHandler) DeleteReview(c *gin.Context) {
	userID, _ := c.Get("userID")
	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	if err := h.reviewUsecase.DeleteReview(c.Request.Context(), userID.(int), reviewID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review deleted successfully"})
}
