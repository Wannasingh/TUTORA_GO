package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/haru/bytestutor/backend/domain"
)

type HttpHandler struct {
	userUsecase  domain.UserUsecase
	tutorUsecase domain.TutorUsecase
}

func NewHttpHandler(r *gin.Engine, uu domain.UserUsecase, tu domain.TutorUsecase) {
	handler := &HttpHandler{
		userUsecase:  uu,
		tutorUsecase: tu,
	}

	api := r.Group("/api")
	{
		api.POST("/users", handler.RegisterUser)
		api.GET("/users/:id", handler.GetUserProfile)
		api.POST("/tutors", handler.BecomeTutor)
		api.GET("/tutors/:id", handler.GetTutorProfile)
		api.GET("/tutors", handler.SearchTutors)
	}
}

func (h *HttpHandler) RegisterUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userUsecase.Register(c.Request.Context(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *HttpHandler) GetUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
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
