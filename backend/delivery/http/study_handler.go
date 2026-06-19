package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type StudyHandler struct {
	studyUsecase domain.StudyUsecase
}

func NewStudyHandler(su domain.StudyUsecase) *StudyHandler {
	return &StudyHandler{studyUsecase: su}
}

// ============ NOTES ============

func (h *StudyHandler) CreateNote(c *gin.Context) {
	userID, _ := c.Get("userID")
	var note domain.UserNote
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note.UserID = userID.(int)
	if err := h.studyUsecase.CreateNote(c.Request.Context(), &note); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, note)
}

func (h *StudyHandler) ListNotes(c *gin.Context) {
	userID, _ := c.Get("userID")
	notes, err := h.studyUsecase.ListMyNotes(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notes)
}

func (h *StudyHandler) GetNote(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	note, err := h.studyUsecase.GetNote(c.Request.Context(), userID.(int), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, note)
}

func (h *StudyHandler) UpdateNote(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	var note domain.UserNote
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	note.ID = id
	if err := h.studyUsecase.UpdateNote(c.Request.Context(), userID.(int), &note); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "note updated"})
}

func (h *StudyHandler) DeleteNote(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.studyUsecase.DeleteNote(c.Request.Context(), userID.(int), id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "note deleted"})
}

// ============ FLASHCARD DECKS ============

func (h *StudyHandler) CreateDeck(c *gin.Context) {
	userID, _ := c.Get("userID")
	var deck domain.FlashcardDeck
	if err := c.ShouldBindJSON(&deck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deck.UserID = userID.(int)
	if err := h.studyUsecase.CreateDeck(c.Request.Context(), &deck); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, deck)
}

func (h *StudyHandler) ListDecks(c *gin.Context) {
	userID, _ := c.Get("userID")
	decks, err := h.studyUsecase.ListMyDecks(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, decks)
}

func (h *StudyHandler) GetDeck(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	deck, err := h.studyUsecase.GetDeckWithCards(c.Request.Context(), userID.(int), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, deck)
}

func (h *StudyHandler) UpdateDeck(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	var deck domain.FlashcardDeck
	if err := c.ShouldBindJSON(&deck); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deck.ID = id
	if err := h.studyUsecase.UpdateDeck(c.Request.Context(), userID.(int), &deck); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deck updated"})
}

func (h *StudyHandler) DeleteDeck(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.studyUsecase.DeleteDeck(c.Request.Context(), userID.(int), id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deck deleted"})
}

// ============ FLASHCARDS ============

func (h *StudyHandler) AddCard(c *gin.Context) {
	userID, _ := c.Get("userID")
	deckID, _ := strconv.Atoi(c.Param("id"))
	var card domain.Flashcard
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card.DeckID = deckID
	if err := h.studyUsecase.AddCard(c.Request.Context(), userID.(int), &card); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, card)
}

func (h *StudyHandler) UpdateCard(c *gin.Context) {
	userID, _ := c.Get("userID")
	cardID, _ := strconv.Atoi(c.Param("cardId"))
	var card domain.Flashcard
	if err := c.ShouldBindJSON(&card); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	card.ID = cardID
	if err := h.studyUsecase.UpdateCard(c.Request.Context(), userID.(int), &card); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "card updated"})
}

func (h *StudyHandler) DeleteCard(c *gin.Context) {
	userID, _ := c.Get("userID")
	cardID, _ := strconv.Atoi(c.Param("cardId"))
	if err := h.studyUsecase.DeleteCard(c.Request.Context(), userID.(int), cardID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "card deleted"})
}

// ============ COURSES ============

func (h *StudyHandler) CreateCourse(c *gin.Context) {
	userID, _ := c.Get("userID")
	var course domain.UserCourse
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	course.UserID = userID.(int)
	if err := h.studyUsecase.CreateCourse(c.Request.Context(), &course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, course)
}

func (h *StudyHandler) ListCourses(c *gin.Context) {
	userID, _ := c.Get("userID")
	courses, err := h.studyUsecase.ListUserCourses(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, courses)
}

func (h *StudyHandler) UpdateCourse(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	var course domain.UserCourse
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	course.ID = id
	if err := h.studyUsecase.UpdateCourse(c.Request.Context(), userID.(int), &course); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "course updated"})
}

func (h *StudyHandler) DeleteCourse(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.studyUsecase.DeleteCourse(c.Request.Context(), userID.(int), id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "course deleted"})
}

// ============ EXAMS ============

func (h *StudyHandler) CreateExam(c *gin.Context) {
	userID, _ := c.Get("userID")
	var exam domain.UserExam
	if err := c.ShouldBindJSON(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exam.UserID = userID.(int)
	if err := h.studyUsecase.CreateExam(c.Request.Context(), &exam); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, exam)
}

func (h *StudyHandler) ListExams(c *gin.Context) {
	userID, _ := c.Get("userID")
	exams, err := h.studyUsecase.ListUserExams(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, exams)
}

func (h *StudyHandler) UpdateExam(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	var exam domain.UserExam
	if err := c.ShouldBindJSON(&exam); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	exam.ID = id
	if err := h.studyUsecase.UpdateExam(c.Request.Context(), userID.(int), &exam); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "exam updated"})
}

func (h *StudyHandler) DeleteExam(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.studyUsecase.DeleteExam(c.Request.Context(), userID.(int), id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "exam deleted"})
}

// ============ CERTIFICATIONS ============

func (h *StudyHandler) CreateCertification(c *gin.Context) {
	userID, _ := c.Get("userID")
	var cert domain.Certification
	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cert.UserID = userID.(int)
	if err := h.studyUsecase.CreateCertification(c.Request.Context(), &cert); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cert)
}

func (h *StudyHandler) ListCertifications(c *gin.Context) {
	userID, _ := c.Get("userID")
	certs, err := h.studyUsecase.ListUserCertifications(c.Request.Context(), userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, certs)
}

func (h *StudyHandler) UpdateCertification(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	var cert domain.Certification
	if err := c.ShouldBindJSON(&cert); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cert.ID = id
	if err := h.studyUsecase.UpdateCertification(c.Request.Context(), userID.(int), &cert); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "certification updated"})
}

func (h *StudyHandler) DeleteCertification(c *gin.Context) {
	userID, _ := c.Get("userID")
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.studyUsecase.DeleteCertification(c.Request.Context(), userID.(int), id); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "certification deleted"})
}

// ============ BADGES ============

func (h *StudyHandler) ListAllBadges(c *gin.Context) {
	badges, err := h.studyUsecase.ListAllBadges(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, badges)
}

func (h *StudyHandler) ListUserBadges(c *gin.Context) {
	targetID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	badges, err := h.studyUsecase.ListUserBadges(c.Request.Context(), targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, badges)
}
