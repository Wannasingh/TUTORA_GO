package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Wannasingh/TUTORA_GO/backend/config"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type HttpHandler struct {
	userUsecase  domain.UserUsecase
	tutorUsecase domain.TutorUsecase
	authUsecase  domain.AuthUsecase
}

func NewHttpHandler(r *gin.Engine, uu domain.UserUsecase, tu domain.TutorUsecase, au domain.AuthUsecase, cfg *config.Config) {
	handler := &HttpHandler{
		userUsecase:  uu,
		tutorUsecase: tu,
		authUsecase:  au,
	}

	api := r.Group("/api")
	api.Use(DecryptionMiddleware(cfg))
	api.Use(EncryptionMiddleware(cfg))
	{
		// Public Auth Endpoints
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.RegisterWithEmail)
			auth.POST("/login", handler.LoginWithEmail)
			auth.POST("/google", handler.LoginWithGoogle)
			auth.POST("/apple", handler.LoginWithApple)
		}

		// Public Resources
		api.GET("/tutors/:id", handler.GetTutorProfile)
		api.GET("/tutors", handler.SearchTutors)

		// Protected Endpoints
		protected := api.Group("")
		protected.Use(AuthMiddleware())
		{
			protected.GET("/users/:id", handler.GetUserProfile)
			protected.POST("/tutors", handler.BecomeTutor)
			protected.DELETE("/users/me", handler.DeleteAccount)
		}
	}
}

func (h *HttpHandler) RegisterWithEmail(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUsecase.RegisterWithEmail(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *HttpHandler) LoginWithEmail(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUsecase.LoginWithEmail(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) LoginWithGoogle(c *gin.Context) {
	var req domain.OAuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUsecase.LoginWithGoogle(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) LoginWithApple(c *gin.Context) {
	var req domain.OAuthLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.authUsecase.LoginWithApple(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *HttpHandler) GetUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Verify requester matches the requested profile (basic guard)
	requesterID, _ := c.Get("userID")
	if requesterID.(int) != id {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	user, err := h.userUsecase.GetProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *HttpHandler) BecomeTutor(c *gin.Context) {
	var tutor domain.Tutor
	if err := c.ShouldBindJSON(&tutor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify user_id matches request
	requesterID, _ := c.Get("userID")
	if requesterID.(int) != tutor.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "cannot register tutor profile for another user"})
		return
	}

	if err := h.tutorUsecase.BecomeTutor(c.Request.Context(), &tutor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tutor)
}

func (h *HttpHandler) GetTutorProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tutor id"})
		return
	}

	tutor, err := h.tutorUsecase.GetTutorProfile(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if tutor == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "tutor profile not found"})
		return
	}

	c.JSON(http.StatusOK, tutor)
}

func (h *HttpHandler) SearchTutors(c *gin.Context) {
	subject := c.Query("subject")

	tutors, err := h.tutorUsecase.SearchTutors(c.Request.Context(), subject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tutors)
}

func (h *HttpHandler) DeleteAccount(c *gin.Context) {
	requesterID, _ := c.Get("userID")

	if err := h.userUsecase.DeleteAccount(c.Request.Context(), requesterID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "account successfully deleted"})
}
