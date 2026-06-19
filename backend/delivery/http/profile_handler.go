package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type ProfileHandler struct {
	userUsecase domain.UserUsecase
	postUsecase domain.PostUsecase
}

func NewProfileHandler(uu domain.UserUsecase, pu domain.PostUsecase) *ProfileHandler {
	return &ProfileHandler{userUsecase: uu, postUsecase: pu}
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req domain.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userUsecase.UpdateProfile(c.Request.Context(), userID.(int), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *ProfileHandler) GetFullProfile(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	requesterID := 0
	if rid, exists := c.Get("userID"); exists {
		requesterID = rid.(int)
	}

	profile, err := h.userUsecase.GetFullProfile(c.Request.Context(), targetID, requesterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) GetUserPosts(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	requesterID := 0
	if rid, exists := c.Get("userID"); exists {
		requesterID = rid.(int)
	}

	posts, err := h.postUsecase.GetUserPosts(c.Request.Context(), targetID, requesterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *ProfileHandler) GetUserLikedPosts(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	requesterID := 0
	if rid, exists := c.Get("userID"); exists {
		requesterID = rid.(int)
	}

	posts, err := h.postUsecase.GetUserLikedPosts(c.Request.Context(), targetID, requesterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *ProfileHandler) GetUserSavedPosts(c *gin.Context) {
	userID, _ := c.Get("userID")
	posts, err := h.postUsecase.GetUserSavedPosts(c.Request.Context(), userID.(int), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *ProfileHandler) GetUserRepostedPosts(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	requesterID := 0
	if rid, exists := c.Get("userID"); exists {
		requesterID = rid.(int)
	}

	posts, err := h.postUsecase.GetUserRepostedPosts(c.Request.Context(), targetID, requesterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, posts)
}

func (h *ProfileHandler) ToggleRepost(c *gin.Context) {
	userID, _ := c.Get("userID")
	postID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	reposted, err := h.postUsecase.ToggleRepost(c.Request.Context(), postID, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reposted": reposted})
}

func (h *ProfileHandler) CreateQuotePost(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req domain.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &domain.Post{
		UserID:         userID.(int),
		Subject:        req.Subject,
		Title:          req.Title,
		Body:           req.Body,
		ImageURL:       req.ImageURL,
		VideoURL:       req.VideoURL,
		OriginalPostID: req.OriginalPostID,
	}

	if err := h.postUsecase.CreateQuotePost(c.Request.Context(), post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}
