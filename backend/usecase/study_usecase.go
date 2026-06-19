package usecase

import (
	"context"
	"fmt"

	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type studyUsecase struct {
	repo domain.StudyRepository
}

func NewStudyUsecase(repo domain.StudyRepository) domain.StudyUsecase {
	return &studyUsecase{repo: repo}
}

// ============ NOTES ============

func (u *studyUsecase) CreateNote(ctx context.Context, note *domain.UserNote) error {
	return u.repo.CreateNote(ctx, note)
}

func (u *studyUsecase) GetNote(ctx context.Context, userID, noteID int) (*domain.UserNote, error) {
	note, err := u.repo.GetNoteByID(ctx, noteID)
	if err != nil {
		return nil, err
	}
	if note == nil {
		return nil, fmt.Errorf("note not found")
	}
	if note.UserID != userID {
		return nil, fmt.Errorf("not authorized to view this note")
	}
	return note, nil
}

func (u *studyUsecase) ListMyNotes(ctx context.Context, userID int) ([]*domain.UserNote, error) {
	return u.repo.ListNotesByUser(ctx, userID)
}

func (u *studyUsecase) UpdateNote(ctx context.Context, userID int, note *domain.UserNote) error {
	existing, err := u.repo.GetNoteByID(ctx, note.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("note not found")
	}
	if existing.UserID != userID {
		return fmt.Errorf("not authorized to update this note")
	}
	return u.repo.UpdateNote(ctx, note)
}

func (u *studyUsecase) DeleteNote(ctx context.Context, userID, noteID int) error {
	existing, err := u.repo.GetNoteByID(ctx, noteID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("note not found")
	}
	if existing.UserID != userID {
		return fmt.Errorf("not authorized to delete this note")
	}
	return u.repo.DeleteNote(ctx, noteID)
}

// ============ FLASHCARD DECKS ============

func (u *studyUsecase) CreateDeck(ctx context.Context, deck *domain.FlashcardDeck) error {
	return u.repo.CreateDeck(ctx, deck)
}

func (u *studyUsecase) GetDeckWithCards(ctx context.Context, userID, deckID int) (*domain.FlashcardDeck, error) {
	deck, err := u.repo.GetDeckByID(ctx, deckID)
	if err != nil {
		return nil, err
	}
	if deck == nil {
		return nil, fmt.Errorf("deck not found")
	}
	if deck.UserID != userID {
		return nil, fmt.Errorf("not authorized to view this deck")
	}
	cards, err := u.repo.GetCardsByDeckID(ctx, deckID)
	if err != nil {
		return nil, err
	}
	deck.Cards = cards
	return deck, nil
}

func (u *studyUsecase) ListMyDecks(ctx context.Context, userID int) ([]*domain.FlashcardDeck, error) {
	return u.repo.ListDecksByUser(ctx, userID)
}

func (u *studyUsecase) UpdateDeck(ctx context.Context, userID int, deck *domain.FlashcardDeck) error {
	existing, err := u.repo.GetDeckByID(ctx, deck.ID)
	if err != nil {
		return err
	}
	if existing == nil || existing.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	return u.repo.UpdateDeck(ctx, deck)
}

func (u *studyUsecase) DeleteDeck(ctx context.Context, userID, deckID int) error {
	existing, err := u.repo.GetDeckByID(ctx, deckID)
	if err != nil {
		return err
	}
	if existing == nil || existing.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	return u.repo.DeleteDeck(ctx, deckID)
}

func (u *studyUsecase) AddCard(ctx context.Context, userID int, card *domain.Flashcard) error {
	deck, err := u.repo.GetDeckByID(ctx, card.DeckID)
	if err != nil {
		return err
	}
	if deck == nil || deck.UserID != userID {
		return fmt.Errorf("not authorized to add cards to this deck")
	}
	return u.repo.CreateCard(ctx, card)
}

func (u *studyUsecase) UpdateCard(ctx context.Context, userID int, card *domain.Flashcard) error {
	existing, err := u.repo.GetCardByID(ctx, card.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("card not found")
	}
	deck, err := u.repo.GetDeckByID(ctx, existing.DeckID)
	if err != nil {
		return err
	}
	if deck == nil || deck.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	return u.repo.UpdateCard(ctx, card)
}

func (u *studyUsecase) DeleteCard(ctx context.Context, userID, cardID int) error {
	existing, err := u.repo.GetCardByID(ctx, cardID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("card not found")
	}
	deck, err := u.repo.GetDeckByID(ctx, existing.DeckID)
	if err != nil {
		return err
	}
	if deck == nil || deck.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	return u.repo.DeleteCard(ctx, cardID)
}

// ============ COURSES ============

func (u *studyUsecase) CreateCourse(ctx context.Context, course *domain.UserCourse) error {
	return u.repo.CreateCourse(ctx, course)
}

func (u *studyUsecase) ListUserCourses(ctx context.Context, userID int) ([]*domain.UserCourse, error) {
	return u.repo.ListCoursesByUser(ctx, userID)
}

func (u *studyUsecase) UpdateCourse(ctx context.Context, userID int, course *domain.UserCourse) error {
	existing, err := u.repo.GetCourseByID(ctx, course.ID)
	if err != nil {
		return err
	}
	if existing == nil || existing.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	return u.repo.UpdateCourse(ctx, course)
}

func (u *studyUsecase) DeleteCourse(ctx context.Context, userID, courseID int) error {
	existing, err := u.repo.GetCourseByID(ctx, courseID)
	if err != nil {
		return err
	}
	if existing == nil || existing.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	return u.repo.DeleteCourse(ctx, courseID)
}

// ============ EXAMS ============

func (u *studyUsecase) CreateExam(ctx context.Context, exam *domain.UserExam) error {
	return u.repo.CreateExam(ctx, exam)
}

func (u *studyUsecase) ListUserExams(ctx context.Context, userID int) ([]*domain.UserExam, error) {
	return u.repo.ListExamsByUser(ctx, userID)
}

func (u *studyUsecase) UpdateExam(ctx context.Context, userID int, exam *domain.UserExam) error {
	existing, err := u.repo.GetExamByID(ctx, exam.ID)
	if err != nil {
		return err
	}
	if existing == nil || existing.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	return u.repo.UpdateExam(ctx, exam)
}

func (u *studyUsecase) DeleteExam(ctx context.Context, userID, examID int) error {
	existing, err := u.repo.GetExamByID(ctx, examID)
	if err != nil {
		return err
	}
	if existing == nil || existing.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	return u.repo.DeleteExam(ctx, examID)
}

// ============ CERTIFICATIONS ============

func (u *studyUsecase) CreateCertification(ctx context.Context, cert *domain.Certification) error {
	return u.repo.CreateCertification(ctx, cert)
}

func (u *studyUsecase) ListUserCertifications(ctx context.Context, userID int) ([]*domain.Certification, error) {
	return u.repo.ListCertificationsByUser(ctx, userID)
}

func (u *studyUsecase) UpdateCertification(ctx context.Context, userID int, cert *domain.Certification) error {
	existing, err := u.repo.GetCertificationByID(ctx, cert.ID)
	if err != nil {
		return err
	}
	if existing == nil || existing.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	return u.repo.UpdateCertification(ctx, cert)
}

func (u *studyUsecase) DeleteCertification(ctx context.Context, userID, certID int) error {
	existing, err := u.repo.GetCertificationByID(ctx, certID)
	if err != nil {
		return err
	}
	if existing == nil || existing.UserID != userID {
		return fmt.Errorf("not authorized")
	}
	return u.repo.DeleteCertification(ctx, certID)
}

// ============ BADGES ============

func (u *studyUsecase) ListAllBadges(ctx context.Context) ([]*domain.Badge, error) {
	return u.repo.ListAllBadges(ctx)
}

func (u *studyUsecase) ListUserBadges(ctx context.Context, userID int) ([]*domain.UserBadge, error) {
	return u.repo.ListUserBadges(ctx, userID)
}
