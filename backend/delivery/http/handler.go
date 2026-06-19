package http

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Wannasingh/TUTORA_GO/backend/config"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
	"github.com/Wannasingh/TUTORA_GO/backend/utils"
)

type HttpHandler struct {
	userUsecase    domain.UserUsecase
	tutorUsecase   domain.TutorUsecase
	authUsecase    domain.AuthUsecase
	postUsecase    domain.PostUsecase
	storageService utils.StorageService
}

func NewHttpHandler(
	r *gin.Engine,
	uu domain.UserUsecase,
	tu domain.TutorUsecase,
	au domain.AuthUsecase,
	pu domain.PostUsecase,
	storage utils.StorageService,
	cfg *config.Config,
	followH *FollowHandler,
	profileH *ProfileHandler,
	msgH *MessageHandler,
	notifH *NotificationHandler,
	studyH *StudyHandler,
) {
	handler := &HttpHandler{
		userUsecase:    uu,
		tutorUsecase:   tu,
		authUsecase:    au,
		postUsecase:    pu,
		storageService: storage,
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
		api.GET("/posts", handler.ListFeed)
		api.GET("/posts/:id", handler.GetPostDetails)
		api.GET("/store/items", handler.BrowseItems)
		api.GET("/store/items/:id", handler.GetItemDetails)
		api.GET("/badges", studyH.ListAllBadges)

		// Protected Endpoints
		protected := api.Group("")
		protected.Use(AuthMiddleware())
		{
			protected.GET("/users/:id", handler.GetUserProfile)
			protected.POST("/tutors", handler.BecomeTutor)
			protected.DELETE("/users/me", handler.DeleteAccount)

			// Protected Post Actions
			protected.POST("/posts", handler.CreatePost)
			protected.PUT("/posts/:id", handler.UpdatePost)
			protected.DELETE("/posts/:id", handler.DeletePost)
			protected.POST("/posts/:id/like", handler.ToggleLike)
			protected.POST("/posts/:id/save", handler.ToggleSave)
			protected.POST("/posts/:id/comments", handler.AddComment)
			protected.POST("/comments/:id/like", handler.ToggleCommentLike)
			protected.DELETE("/comments/:id", handler.DeleteComment)
			protected.POST("/reports", handler.ReportContent)

			// Protected Upload Action
			protected.POST("/upload", handler.UploadImage)

			// Protected Store & Wallet Actions
			protected.POST("/store/items", handler.ListItem)
			protected.POST("/store/purchase-with-coins", handler.PurchaseItem)
			protected.GET("/store/purchases/my-library", handler.GetMyLibrary)
			protected.GET("/store/purchases/my-library/:id/download", handler.GetDownloadURL)
			protected.GET("/wallet/balance", handler.GetWalletInfo)
			protected.POST("/wallet/iap-validate", handler.ValidateAppleReceipt)
			protected.POST("/store/withdrawals", handler.RequestPayout)

			// Protected Store Interactive Q&A Actions
			protected.POST("/store/items/:id/qa-pins", handler.CreateQAPin)
			protected.GET("/store/items/:id/qa-pins", handler.GetQAPins)
			protected.POST("/store/qa-pins/:pin_id/replies", handler.ReplyToQAPin)

			// ============ NEW: Follow & Social ============
			protected.POST("/users/:id/follow", followH.ToggleFollow)
			protected.GET("/users/:id/followers", followH.GetFollowers)
			protected.GET("/users/:id/following", followH.GetFollowing)
			protected.GET("/users/:id/follow-stats", followH.GetFollowStats)

			// ============ NEW: Tutor Reviews ============
			protected.POST("/tutors/:id/reviews", followH.SubmitReview)
			protected.GET("/tutors/:id/reviews", followH.GetTutorReviews)
			protected.DELETE("/reviews/:id", followH.DeleteReview)

			// ============ NEW: Profile ============
			protected.PUT("/profile", profileH.UpdateProfile)
			protected.GET("/profile/:id", profileH.GetFullProfile)
			protected.GET("/users/:id/posts", profileH.GetUserPosts)
			protected.GET("/users/:id/liked", profileH.GetUserLikedPosts)
			protected.GET("/users/:id/saved", profileH.GetUserSavedPosts)
			protected.GET("/users/:id/reposts", profileH.GetUserRepostedPosts)
			protected.POST("/posts/:id/repost", profileH.ToggleRepost)
			protected.POST("/posts/quote", profileH.CreateQuotePost)

			// ============ NEW: Messaging ============
			protected.POST("/conversations", msgH.StartConversation)
			protected.GET("/conversations", msgH.ListMyConversations)
			protected.GET("/conversations/:id", msgH.GetConversation)
			protected.GET("/conversations/:id/messages", msgH.GetMessages)
			protected.POST("/conversations/:id/messages", msgH.SendMessage)
			protected.POST("/conversations/:id/read", msgH.MarkAsRead)

			// ============ NEW: Notifications ============
			protected.GET("/notifications", notifH.ListNotifications)
			protected.GET("/notifications/unread-count", notifH.GetUnreadCount)
			protected.POST("/notifications/:id/read", notifH.MarkRead)
			protected.POST("/notifications/read-all", notifH.MarkAllRead)
			protected.DELETE("/notifications/:id", notifH.DeleteNotification)

			// ============ NEW: Study Tools ============
			// Notes
			protected.POST("/notes", studyH.CreateNote)
			protected.GET("/notes", studyH.ListNotes)
			protected.GET("/notes/:id", studyH.GetNote)
			protected.PUT("/notes/:id", studyH.UpdateNote)
			protected.DELETE("/notes/:id", studyH.DeleteNote)

			// Flashcard Decks
			protected.POST("/decks", studyH.CreateDeck)
			protected.GET("/decks", studyH.ListDecks)
			protected.GET("/decks/:id", studyH.GetDeck)
			protected.PUT("/decks/:id", studyH.UpdateDeck)
			protected.DELETE("/decks/:id", studyH.DeleteDeck)

			// Flashcards within a deck
			protected.POST("/decks/:id/cards", studyH.AddCard)
			protected.PUT("/cards/:cardId", studyH.UpdateCard)
			protected.DELETE("/cards/:cardId", studyH.DeleteCard)

			// Courses
			protected.POST("/courses", studyH.CreateCourse)
			protected.GET("/courses", studyH.ListCourses)
			protected.PUT("/courses/:id", studyH.UpdateCourse)
			protected.DELETE("/courses/:id", studyH.DeleteCourse)

			// Exams
			protected.POST("/exams", studyH.CreateExam)
			protected.GET("/exams", studyH.ListExams)
			protected.PUT("/exams/:id", studyH.UpdateExam)
			protected.DELETE("/exams/:id", studyH.DeleteExam)

			// Certifications
			protected.POST("/certifications", studyH.CreateCertification)
			protected.GET("/certifications", studyH.ListCertifications)
			protected.PUT("/certifications/:id", studyH.UpdateCertification)
			protected.DELETE("/certifications/:id", studyH.DeleteCertification)

			// Badges
			protected.GET("/users/:id/badges", studyH.ListUserBadges)
		}
	}

	// WebSocket endpoint (outside encrypted API group)
	r.GET("/ws", msgH.HandleWebSocket)
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

func (h *HttpHandler) CreatePost(c *gin.Context) {
	requesterID, _ := c.Get("userID")
	var req domain.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post := &domain.Post{
		UserID:   requesterID.(int),
		Subject:  req.Subject,
		Title:    req.Title,
		Body:     req.Body,
		ImageURL: req.ImageURL,
		VideoURL: req.VideoURL,
	}

	if err := h.postUsecase.CreatePost(c.Request.Context(), post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *HttpHandler) ListFeed(c *gin.Context) {
	subject := c.Query("subject")
	
	var requesterUserID int
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString)
		if err == nil {
			if sub, ok := claims["sub"].(float64); ok {
				requesterUserID = int(sub)
			}
		}
	}

	posts, err := h.postUsecase.ListFeed(c.Request.Context(), subject, requesterUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *HttpHandler) GetPostDetails(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var requesterUserID int
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString)
		if err == nil {
			if sub, ok := claims["sub"].(float64); ok {
				requesterUserID = int(sub)
			}
		}
	}

	post, comments, err := h.postUsecase.GetPostDetails(c.Request.Context(), id, requesterUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"post":     post,
		"comments": comments,
	})
}

func (h *HttpHandler) ToggleLike(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	requesterID, _ := c.Get("userID")
	liked, err := h.postUsecase.ToggleLike(c.Request.Context(), postID, requesterID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"liked": liked})
}

func (h *HttpHandler) ToggleSave(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	requesterID, _ := c.Get("userID")
	saved, err := h.postUsecase.ToggleSave(c.Request.Context(), postID, requesterID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"saved": saved})
}

func (h *HttpHandler) AddComment(c *gin.Context) {
	idStr := c.Param("id")
	postID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	requesterID, _ := c.Get("userID")
	var req domain.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment := &domain.Comment{
		PostID:   postID,
		UserID:   requesterID.(int),
		Body:     req.Body,
		ImageURL: req.ImageURL,
		ParentID: req.ParentID,
	}

	if err := h.postUsecase.AddPostComment(c.Request.Context(), comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *HttpHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file 'image' is required"})
		return
	}
	defer file.Close()

	// Verify it's an image
	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/gif" && contentType != "image/webp" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "only jpeg, png, gif, and webp images are allowed"})
		return
	}

	// Compress and optimize image to JPEG to save storage space
	optimizedReader, optimizedContentType, err := utils.OptimizeImage(file, contentType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("failed to parse/compress image: %v", err)})
		return
	}

	url, err := h.storageService.UploadFile(c.Request.Context(), header.Filename, optimizedReader, optimizedContentType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
