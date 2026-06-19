package domain

import "context"

// Personal Notes
type UserNote struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Title     string  `json:"title"`
	Body      *string `json:"body,omitempty"`
	Subject   *string `json:"subject,omitempty"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// Flashcard Decks
type FlashcardDeck struct {
	ID        int          `json:"id"`
	UserID    int          `json:"user_id"`
	Title     string       `json:"title"`
	Subject   *string      `json:"subject,omitempty"`
	CardCount int          `json:"card_count"`
	Cards     []*Flashcard `json:"cards,omitempty"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}

type Flashcard struct {
	ID        int     `json:"id"`
	DeckID    int     `json:"deck_id"`
	FrontText string  `json:"front_text"`
	BackText  string  `json:"back_text"`
	ImageURL  *string `json:"image_url,omitempty"`
	SortOrder int     `json:"sort_order"`
	CreatedAt string  `json:"created_at"`
}

// Course Tracking
type UserCourse struct {
	ID          int     `json:"id"`
	UserID      int     `json:"user_id"`
	Title       string  `json:"title"`
	Institution *string `json:"institution,omitempty"`
	Status      string  `json:"status"`
	StartedAt   *string `json:"started_at,omitempty"`
	CompletedAt *string `json:"completed_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

// Exam Records
type UserExam struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	Title     string  `json:"title"`
	Subject   *string `json:"subject,omitempty"`
	Score     *string `json:"score,omitempty"`
	MaxScore  *string `json:"max_score,omitempty"`
	ExamDate  *string `json:"exam_date,omitempty"`
	CreatedAt string  `json:"created_at"`
}

// Certifications
type Certification struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	Title         string  `json:"title"`
	Issuer        string  `json:"issuer"`
	DateEarned    *string `json:"date_earned,omitempty"`
	ExpiryDate    *string `json:"expiry_date,omitempty"`
	ImageURL      *string `json:"image_url,omitempty"`
	CredentialURL *string `json:"credential_url,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

// Badges
type Badge struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IconURL     *string `json:"icon_url,omitempty"`
	Criteria    string  `json:"criteria"`
}

type UserBadge struct {
	Badge    Badge  `json:"badge"`
	EarnedAt string `json:"earned_at"`
}

type StudyRepository interface {
	// Notes
	CreateNote(ctx context.Context, note *UserNote) error
	GetNoteByID(ctx context.Context, id int) (*UserNote, error)
	ListNotesByUser(ctx context.Context, userID int) ([]*UserNote, error)
	UpdateNote(ctx context.Context, note *UserNote) error
	DeleteNote(ctx context.Context, id int) error

	// Flashcard Decks
	CreateDeck(ctx context.Context, deck *FlashcardDeck) error
	GetDeckByID(ctx context.Context, id int) (*FlashcardDeck, error)
	ListDecksByUser(ctx context.Context, userID int) ([]*FlashcardDeck, error)
	UpdateDeck(ctx context.Context, deck *FlashcardDeck) error
	DeleteDeck(ctx context.Context, id int) error

	// Flashcards
	CreateCard(ctx context.Context, card *Flashcard) error
	GetCardByID(ctx context.Context, id int) (*Flashcard, error)
	GetCardsByDeckID(ctx context.Context, deckID int) ([]*Flashcard, error)
	UpdateCard(ctx context.Context, card *Flashcard) error
	DeleteCard(ctx context.Context, id int) error

	// Courses
	CreateCourse(ctx context.Context, course *UserCourse) error
	GetCourseByID(ctx context.Context, id int) (*UserCourse, error)
	ListCoursesByUser(ctx context.Context, userID int) ([]*UserCourse, error)
	UpdateCourse(ctx context.Context, course *UserCourse) error
	DeleteCourse(ctx context.Context, id int) error

	// Exams
	CreateExam(ctx context.Context, exam *UserExam) error
	GetExamByID(ctx context.Context, id int) (*UserExam, error)
	ListExamsByUser(ctx context.Context, userID int) ([]*UserExam, error)
	UpdateExam(ctx context.Context, exam *UserExam) error
	DeleteExam(ctx context.Context, id int) error

	// Certifications
	CreateCertification(ctx context.Context, cert *Certification) error
	GetCertificationByID(ctx context.Context, id int) (*Certification, error)
	ListCertificationsByUser(ctx context.Context, userID int) ([]*Certification, error)
	UpdateCertification(ctx context.Context, cert *Certification) error
	DeleteCertification(ctx context.Context, id int) error

	// Badges
	ListAllBadges(ctx context.Context) ([]*Badge, error)
	ListUserBadges(ctx context.Context, userID int) ([]*UserBadge, error)
	AwardBadge(ctx context.Context, userID, badgeID int) error
}

type StudyUsecase interface {
	// Notes
	CreateNote(ctx context.Context, note *UserNote) error
	GetNote(ctx context.Context, userID, noteID int) (*UserNote, error)
	ListMyNotes(ctx context.Context, userID int) ([]*UserNote, error)
	UpdateNote(ctx context.Context, userID int, note *UserNote) error
	DeleteNote(ctx context.Context, userID, noteID int) error

	// Flashcard Decks & Cards
	CreateDeck(ctx context.Context, deck *FlashcardDeck) error
	GetDeckWithCards(ctx context.Context, userID, deckID int) (*FlashcardDeck, error)
	ListMyDecks(ctx context.Context, userID int) ([]*FlashcardDeck, error)
	UpdateDeck(ctx context.Context, userID int, deck *FlashcardDeck) error
	DeleteDeck(ctx context.Context, userID, deckID int) error
	AddCard(ctx context.Context, userID int, card *Flashcard) error
	UpdateCard(ctx context.Context, userID int, card *Flashcard) error
	DeleteCard(ctx context.Context, userID, cardID int) error

	// Courses
	CreateCourse(ctx context.Context, course *UserCourse) error
	ListUserCourses(ctx context.Context, userID int) ([]*UserCourse, error)
	UpdateCourse(ctx context.Context, userID int, course *UserCourse) error
	DeleteCourse(ctx context.Context, userID, courseID int) error

	// Exams
	CreateExam(ctx context.Context, exam *UserExam) error
	ListUserExams(ctx context.Context, userID int) ([]*UserExam, error)
	UpdateExam(ctx context.Context, userID int, exam *UserExam) error
	DeleteExam(ctx context.Context, userID, examID int) error

	// Certifications
	CreateCertification(ctx context.Context, cert *Certification) error
	ListUserCertifications(ctx context.Context, userID int) ([]*Certification, error)
	UpdateCertification(ctx context.Context, userID int, cert *Certification) error
	DeleteCertification(ctx context.Context, userID, certID int) error

	// Badges
	ListAllBadges(ctx context.Context) ([]*Badge, error)
	ListUserBadges(ctx context.Context, userID int) ([]*UserBadge, error)
}
