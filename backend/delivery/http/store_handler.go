package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
	"github.com/Wannasingh/TUTORA_GO/backend/utils"
)

// Extension methods for HttpHandler to handle store/wallet actions

func (h *HttpHandler) ListItem(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		Title        string `json:"title" binding:"required"`
		Description  string `json:"description" binding:"required"`
		Category     string `json:"category" binding:"required"`
		Subject      string `json:"subject" binding:"required"`
		PriceInCoins int    `json:"price_in_coins" binding:"required,gt=0"`
		FileURL      string `json:"file_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := &domain.StoreItem{
		SellerID:     userID.(int),
		Title:        req.Title,
		Description:  req.Description,
		Category:     req.Category,
		Subject:      req.Subject,
		PriceInCoins: req.PriceInCoins,
		FileURL:      req.FileURL,
	}

	if err := h.postUsecase.(*StoreUsecaseHelper).StoreUsecase.ListItem(c.Request.Context(), item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *HttpHandler) BrowseItems(c *gin.Context) {
	category := c.Query("category")
	subject := c.Query("subject")
	search := c.Query("search")

	// Read optional authorization header to personalize the feed if token is present
	userID := 0
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString)
		if err == nil {
			if sub, ok := claims["sub"].(float64); ok {
				userID = int(sub)
			}
		}
	}

	items, err := h.postUsecase.(*StoreUsecaseHelper).StoreUsecase.BrowseItems(c.Request.Context(), category, subject, search, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *HttpHandler) GetItemDetails(c *gin.Context) {
	idStr := c.Param("id")
	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	// Read optional authorization header to fetch details
	userID := 0
	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString)
		if err == nil {
			if sub, ok := claims["sub"].(float64); ok {
				userID = int(sub)
			}
		}
	}

	item, err := h.postUsecase.(*StoreUsecaseHelper).StoreUsecase.GetItemDetails(c.Request.Context(), userID, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *HttpHandler) PurchaseItem(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		ItemID int `json:"item_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.postUsecase.(*StoreUsecaseHelper).StoreUsecase.PurchaseItem(c.Request.Context(), userID.(int), req.ItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "purchase successful"})
}

func (h *HttpHandler) GetMyLibrary(c *gin.Context) {
	userID, _ := c.Get("userID")

	items, err := h.postUsecase.(*StoreUsecaseHelper).StoreUsecase.GetMyLibrary(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *HttpHandler) GetDownloadURL(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("id")
	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	url, err := h.postUsecase.(*StoreUsecaseHelper).StoreUsecase.GetDownloadURL(c.Request.Context(), userID.(int), itemID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"download_url": url})
}

func (h *HttpHandler) GetWalletInfo(c *gin.Context) {
	userID, _ := c.Get("userID")

	balance, txs, err := h.postUsecase.(*StoreUsecaseHelper).StoreUsecase.GetWalletInfo(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance_coins": balance,
		"transactions":  txs,
	})
}

func (h *HttpHandler) ValidateAppleReceipt(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		Receipt string `json:"receipt" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.postUsecase.(*StoreUsecaseHelper).StoreUsecase.ValidateAppleReceipt(c.Request.Context(), userID.(int), req.Receipt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "receipt validated and coins credited successfully"})
}

func (h *HttpHandler) RequestPayout(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		Coins       int    `json:"coins" binding:"required,gt=0"`
		BankAccount string `json:"bank_account" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.postUsecase.(*StoreUsecaseHelper).StoreUsecase.RequestPayout(c.Request.Context(), userID.(int), req.Coins, req.BankAccount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "payout request submitted successfully"})
}

func (h *HttpHandler) CreateQAPin(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("id")
	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	var req struct {
		PageNumber   int     `json:"page_number" binding:"required,gt=0"`
		CoordinateX  float64 `json:"coordinate_x" binding:"required"`
		CoordinateY  float64 `json:"coordinate_y" binding:"required"`
		QuestionText string  `json:"question_text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pin := &domain.ItemQAPin{
		ItemID:       itemID,
		UserID:       userID.(int),
		PageNumber:   req.PageNumber,
		CoordinateX:  req.CoordinateX,
		CoordinateY:  req.CoordinateY,
		QuestionText: req.QuestionText,
	}

	err = h.postUsecase.(*StoreUsecaseHelper).StoreUsecase.CreateQAPin(c.Request.Context(), pin)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pin)
}

func (h *HttpHandler) GetQAPins(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("id")
	itemID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}

	pins, err := h.postUsecase.(*StoreUsecaseHelper).StoreUsecase.GetQAPins(c.Request.Context(), userID.(int), itemID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pins)
}

func (h *HttpHandler) ReplyToQAPin(c *gin.Context) {
	userID, _ := c.Get("userID")
	idStr := c.Param("pin_id")
	pinID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pin id"})
		return
	}

	var req struct {
		ReplyText string `json:"reply_text" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reply := &domain.ItemQAReply{
		PinID:     pinID,
		UserID:    userID.(int),
		ReplyText: req.ReplyText,
	}

	err = h.postUsecase.(*StoreUsecaseHelper).StoreUsecase.ReplyToQAPin(c.Request.Context(), reply)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reply)
}

// StoreUsecaseHelper wraps PostUsecase and StoreUsecase inside HttpHandler without breaking original function signatures.
type StoreUsecaseHelper struct {
	domain.PostUsecase
	StoreUsecase domain.StoreUsecase
}

func NewStoreUsecaseHelper(pu domain.PostUsecase, su domain.StoreUsecase) domain.PostUsecase {
	return &StoreUsecaseHelper{
		PostUsecase:  pu,
		StoreUsecase: su,
	}
}
