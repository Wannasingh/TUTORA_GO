package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type postgresStudyRepository struct {
	db *pgxpool.Pool
}

func NewPostgresStudyRepository(db *pgxpool.Pool) domain.StudyRepository {
	return &postgresStudyRepository{db: db}
}

// ============ NOTES ============

func (r *postgresStudyRepository) CreateNote(ctx context.Context, note *domain.UserNote) error {
	query := `INSERT INTO tutora_app.user_notes (user_id, title, body, subject) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`
	var ca, ua time.Time
	err := r.db.QueryRow(ctx, query, note.UserID, note.Title, note.Body, note.Subject).Scan(&note.ID, &ca, &ua)
	if err == nil {
		note.CreatedAt = ca.Format(time.RFC3339)
		note.UpdatedAt = ua.Format(time.RFC3339)
	}
	return err
}

func (r *postgresStudyRepository) GetNoteByID(ctx context.Context, id int) (*domain.UserNote, error) {
	query := `SELECT id, user_id, title, body, subject, created_at, updated_at FROM tutora_app.user_notes WHERE id = $1`
	n := &domain.UserNote{}
	var ca, ua time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.Subject, &ca, &ua)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return nil, nil }
		return nil, err
	}
	n.CreatedAt = ca.Format(time.RFC3339)
	n.UpdatedAt = ua.Format(time.RFC3339)
	return n, nil
}

func (r *postgresStudyRepository) ListNotesByUser(ctx context.Context, userID int) ([]*domain.UserNote, error) {
	query := `SELECT id, user_id, title, body, subject, created_at, updated_at FROM tutora_app.user_notes WHERE user_id = $1 ORDER BY updated_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var notes []*domain.UserNote
	for rows.Next() {
		n := &domain.UserNote{}
		var ca, ua time.Time
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Body, &n.Subject, &ca, &ua); err != nil { return nil, err }
		n.CreatedAt = ca.Format(time.RFC3339)
		n.UpdatedAt = ua.Format(time.RFC3339)
		notes = append(notes, n)
	}
	return notes, nil
}

func (r *postgresStudyRepository) UpdateNote(ctx context.Context, note *domain.UserNote) error {
	query := `UPDATE tutora_app.user_notes SET title = $2, body = $3, subject = $4, updated_at = NOW() WHERE id = $1 RETURNING updated_at`
	var ua time.Time
	err := r.db.QueryRow(ctx, query, note.ID, note.Title, note.Body, note.Subject).Scan(&ua)
	if err == nil { note.UpdatedAt = ua.Format(time.RFC3339) }
	return err
}

func (r *postgresStudyRepository) DeleteNote(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tutora_app.user_notes WHERE id = $1`, id)
	return err
}

// ============ FLASHCARD DECKS ============

func (r *postgresStudyRepository) CreateDeck(ctx context.Context, deck *domain.FlashcardDeck) error {
	query := `INSERT INTO tutora_app.flashcard_decks (user_id, title, subject) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	var ca, ua time.Time
	err := r.db.QueryRow(ctx, query, deck.UserID, deck.Title, deck.Subject).Scan(&deck.ID, &ca, &ua)
	if err == nil {
		deck.CreatedAt = ca.Format(time.RFC3339)
		deck.UpdatedAt = ua.Format(time.RFC3339)
	}
	return err
}

func (r *postgresStudyRepository) GetDeckByID(ctx context.Context, id int) (*domain.FlashcardDeck, error) {
	query := `SELECT d.id, d.user_id, d.title, d.subject, d.created_at, d.updated_at,
	                 (SELECT COUNT(*) FROM tutora_app.flashcards WHERE deck_id = d.id) as card_count
	          FROM tutora_app.flashcard_decks d WHERE d.id = $1`
	d := &domain.FlashcardDeck{}
	var ca, ua time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(&d.ID, &d.UserID, &d.Title, &d.Subject, &ca, &ua, &d.CardCount)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return nil, nil }
		return nil, err
	}
	d.CreatedAt = ca.Format(time.RFC3339)
	d.UpdatedAt = ua.Format(time.RFC3339)
	return d, nil
}

func (r *postgresStudyRepository) ListDecksByUser(ctx context.Context, userID int) ([]*domain.FlashcardDeck, error) {
	query := `SELECT d.id, d.user_id, d.title, d.subject, d.created_at, d.updated_at,
	                 (SELECT COUNT(*) FROM tutora_app.flashcards WHERE deck_id = d.id) as card_count
	          FROM tutora_app.flashcard_decks d WHERE d.user_id = $1 ORDER BY d.updated_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var decks []*domain.FlashcardDeck
	for rows.Next() {
		d := &domain.FlashcardDeck{}
		var ca, ua time.Time
		if err := rows.Scan(&d.ID, &d.UserID, &d.Title, &d.Subject, &ca, &ua, &d.CardCount); err != nil { return nil, err }
		d.CreatedAt = ca.Format(time.RFC3339)
		d.UpdatedAt = ua.Format(time.RFC3339)
		decks = append(decks, d)
	}
	return decks, nil
}

func (r *postgresStudyRepository) UpdateDeck(ctx context.Context, deck *domain.FlashcardDeck) error {
	query := `UPDATE tutora_app.flashcard_decks SET title = $2, subject = $3, updated_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, deck.ID, deck.Title, deck.Subject)
	return err
}

func (r *postgresStudyRepository) DeleteDeck(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tutora_app.flashcard_decks WHERE id = $1`, id)
	return err
}

// ============ FLASHCARDS ============

func (r *postgresStudyRepository) CreateCard(ctx context.Context, card *domain.Flashcard) error {
	query := `INSERT INTO tutora_app.flashcards (deck_id, front_text, back_text, image_url, sort_order) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	var ca time.Time
	err := r.db.QueryRow(ctx, query, card.DeckID, card.FrontText, card.BackText, card.ImageURL, card.SortOrder).Scan(&card.ID, &ca)
	if err == nil { card.CreatedAt = ca.Format(time.RFC3339) }
	return err
}

func (r *postgresStudyRepository) GetCardByID(ctx context.Context, id int) (*domain.Flashcard, error) {
	query := `SELECT id, deck_id, front_text, back_text, image_url, sort_order, created_at FROM tutora_app.flashcards WHERE id = $1`
	c := &domain.Flashcard{}
	var ca time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.DeckID, &c.FrontText, &c.BackText, &c.ImageURL, &c.SortOrder, &ca)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return nil, nil }
		return nil, err
	}
	c.CreatedAt = ca.Format(time.RFC3339)
	return c, nil
}

func (r *postgresStudyRepository) GetCardsByDeckID(ctx context.Context, deckID int) ([]*domain.Flashcard, error) {
	query := `SELECT id, deck_id, front_text, back_text, image_url, sort_order, created_at FROM tutora_app.flashcards WHERE deck_id = $1 ORDER BY sort_order, id`
	rows, err := r.db.Query(ctx, query, deckID)
	if err != nil { return nil, err }
	defer rows.Close()
	var cards []*domain.Flashcard
	for rows.Next() {
		c := &domain.Flashcard{}
		var ca time.Time
		if err := rows.Scan(&c.ID, &c.DeckID, &c.FrontText, &c.BackText, &c.ImageURL, &c.SortOrder, &ca); err != nil { return nil, err }
		c.CreatedAt = ca.Format(time.RFC3339)
		cards = append(cards, c)
	}
	return cards, nil
}

func (r *postgresStudyRepository) UpdateCard(ctx context.Context, card *domain.Flashcard) error {
	query := `UPDATE tutora_app.flashcards SET front_text = $2, back_text = $3, image_url = $4, sort_order = $5 WHERE id = $1`
	_, err := r.db.Exec(ctx, query, card.ID, card.FrontText, card.BackText, card.ImageURL, card.SortOrder)
	return err
}

func (r *postgresStudyRepository) DeleteCard(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tutora_app.flashcards WHERE id = $1`, id)
	return err
}

// ============ COURSES ============

func (r *postgresStudyRepository) CreateCourse(ctx context.Context, course *domain.UserCourse) error {
	query := `INSERT INTO tutora_app.user_courses (user_id, title, institution, status, started_at, completed_at)
	          VALUES ($1, $2, $3, $4, $5::date, $6::date) RETURNING id, created_at`
	var ca time.Time
	err := r.db.QueryRow(ctx, query, course.UserID, course.Title, course.Institution, course.Status, course.StartedAt, course.CompletedAt).Scan(&course.ID, &ca)
	if err == nil { course.CreatedAt = ca.Format(time.RFC3339) }
	return err
}

func (r *postgresStudyRepository) GetCourseByID(ctx context.Context, id int) (*domain.UserCourse, error) {
	query := `SELECT id, user_id, title, institution, status, started_at::text, completed_at::text, created_at FROM tutora_app.user_courses WHERE id = $1`
	c := &domain.UserCourse{}
	var ca time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.UserID, &c.Title, &c.Institution, &c.Status, &c.StartedAt, &c.CompletedAt, &ca)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return nil, nil }
		return nil, err
	}
	c.CreatedAt = ca.Format(time.RFC3339)
	return c, nil
}

func (r *postgresStudyRepository) ListCoursesByUser(ctx context.Context, userID int) ([]*domain.UserCourse, error) {
	query := `SELECT id, user_id, title, institution, status, started_at::text, completed_at::text, created_at FROM tutora_app.user_courses WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var courses []*domain.UserCourse
	for rows.Next() {
		c := &domain.UserCourse{}
		var ca time.Time
		if err := rows.Scan(&c.ID, &c.UserID, &c.Title, &c.Institution, &c.Status, &c.StartedAt, &c.CompletedAt, &ca); err != nil { return nil, err }
		c.CreatedAt = ca.Format(time.RFC3339)
		courses = append(courses, c)
	}
	return courses, nil
}

func (r *postgresStudyRepository) UpdateCourse(ctx context.Context, course *domain.UserCourse) error {
	query := `UPDATE tutora_app.user_courses SET title=$2, institution=$3, status=$4, started_at=$5::date, completed_at=$6::date WHERE id=$1`
	_, err := r.db.Exec(ctx, query, course.ID, course.Title, course.Institution, course.Status, course.StartedAt, course.CompletedAt)
	return err
}

func (r *postgresStudyRepository) DeleteCourse(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tutora_app.user_courses WHERE id = $1`, id)
	return err
}

// ============ EXAMS ============

func (r *postgresStudyRepository) CreateExam(ctx context.Context, exam *domain.UserExam) error {
	query := `INSERT INTO tutora_app.user_exams (user_id, title, subject, score, max_score, exam_date) VALUES ($1,$2,$3,$4,$5,$6::date) RETURNING id, created_at`
	var ca time.Time
	err := r.db.QueryRow(ctx, query, exam.UserID, exam.Title, exam.Subject, exam.Score, exam.MaxScore, exam.ExamDate).Scan(&exam.ID, &ca)
	if err == nil { exam.CreatedAt = ca.Format(time.RFC3339) }
	return err
}

func (r *postgresStudyRepository) GetExamByID(ctx context.Context, id int) (*domain.UserExam, error) {
	query := `SELECT id, user_id, title, subject, score, max_score, exam_date::text, created_at FROM tutora_app.user_exams WHERE id = $1`
	e := &domain.UserExam{}
	var ca time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(&e.ID, &e.UserID, &e.Title, &e.Subject, &e.Score, &e.MaxScore, &e.ExamDate, &ca)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return nil, nil }
		return nil, err
	}
	e.CreatedAt = ca.Format(time.RFC3339)
	return e, nil
}

func (r *postgresStudyRepository) ListExamsByUser(ctx context.Context, userID int) ([]*domain.UserExam, error) {
	query := `SELECT id, user_id, title, subject, score, max_score, exam_date::text, created_at FROM tutora_app.user_exams WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var exams []*domain.UserExam
	for rows.Next() {
		e := &domain.UserExam{}
		var ca time.Time
		if err := rows.Scan(&e.ID, &e.UserID, &e.Title, &e.Subject, &e.Score, &e.MaxScore, &e.ExamDate, &ca); err != nil { return nil, err }
		e.CreatedAt = ca.Format(time.RFC3339)
		exams = append(exams, e)
	}
	return exams, nil
}

func (r *postgresStudyRepository) UpdateExam(ctx context.Context, exam *domain.UserExam) error {
	query := `UPDATE tutora_app.user_exams SET title=$2, subject=$3, score=$4, max_score=$5, exam_date=$6::date WHERE id=$1`
	_, err := r.db.Exec(ctx, query, exam.ID, exam.Title, exam.Subject, exam.Score, exam.MaxScore, exam.ExamDate)
	return err
}

func (r *postgresStudyRepository) DeleteExam(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tutora_app.user_exams WHERE id = $1`, id)
	return err
}

// ============ CERTIFICATIONS ============

func (r *postgresStudyRepository) CreateCertification(ctx context.Context, cert *domain.Certification) error {
	query := `INSERT INTO tutora_app.certifications (user_id, title, issuer, date_earned, expiry_date, image_url, credential_url)
	          VALUES ($1,$2,$3,$4::date,$5::date,$6,$7) RETURNING id, created_at`
	var ca time.Time
	err := r.db.QueryRow(ctx, query, cert.UserID, cert.Title, cert.Issuer, cert.DateEarned, cert.ExpiryDate, cert.ImageURL, cert.CredentialURL).Scan(&cert.ID, &ca)
	if err == nil { cert.CreatedAt = ca.Format(time.RFC3339) }
	return err
}

func (r *postgresStudyRepository) GetCertificationByID(ctx context.Context, id int) (*domain.Certification, error) {
	query := `SELECT id, user_id, title, issuer, date_earned::text, expiry_date::text, image_url, credential_url, created_at FROM tutora_app.certifications WHERE id = $1`
	c := &domain.Certification{}
	var ca time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(&c.ID, &c.UserID, &c.Title, &c.Issuer, &c.DateEarned, &c.ExpiryDate, &c.ImageURL, &c.CredentialURL, &ca)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return nil, nil }
		return nil, err
	}
	c.CreatedAt = ca.Format(time.RFC3339)
	return c, nil
}

func (r *postgresStudyRepository) ListCertificationsByUser(ctx context.Context, userID int) ([]*domain.Certification, error) {
	query := `SELECT id, user_id, title, issuer, date_earned::text, expiry_date::text, image_url, credential_url, created_at FROM tutora_app.certifications WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var certs []*domain.Certification
	for rows.Next() {
		c := &domain.Certification{}
		var ca time.Time
		if err := rows.Scan(&c.ID, &c.UserID, &c.Title, &c.Issuer, &c.DateEarned, &c.ExpiryDate, &c.ImageURL, &c.CredentialURL, &ca); err != nil { return nil, err }
		c.CreatedAt = ca.Format(time.RFC3339)
		certs = append(certs, c)
	}
	return certs, nil
}

func (r *postgresStudyRepository) UpdateCertification(ctx context.Context, cert *domain.Certification) error {
	query := `UPDATE tutora_app.certifications SET title=$2, issuer=$3, date_earned=$4::date, expiry_date=$5::date, image_url=$6, credential_url=$7 WHERE id=$1`
	_, err := r.db.Exec(ctx, query, cert.ID, cert.Title, cert.Issuer, cert.DateEarned, cert.ExpiryDate, cert.ImageURL, cert.CredentialURL)
	return err
}

func (r *postgresStudyRepository) DeleteCertification(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM tutora_app.certifications WHERE id = $1`, id)
	return err
}

// ============ BADGES ============

func (r *postgresStudyRepository) ListAllBadges(ctx context.Context) ([]*domain.Badge, error) {
	query := `SELECT id, name, description, icon_url, criteria FROM tutora_app.badges ORDER BY id`
	rows, err := r.db.Query(ctx, query)
	if err != nil { return nil, err }
	defer rows.Close()
	var badges []*domain.Badge
	for rows.Next() {
		b := &domain.Badge{}
		if err := rows.Scan(&b.ID, &b.Name, &b.Description, &b.IconURL, &b.Criteria); err != nil { return nil, err }
		badges = append(badges, b)
	}
	return badges, nil
}

func (r *postgresStudyRepository) ListUserBadges(ctx context.Context, userID int) ([]*domain.UserBadge, error) {
	query := `SELECT b.id, b.name, b.description, b.icon_url, b.criteria, ub.earned_at
	          FROM tutora_app.user_badges ub
	          JOIN tutora_app.badges b ON b.id = ub.badge_id
	          WHERE ub.user_id = $1
	          ORDER BY ub.earned_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil { return nil, err }
	defer rows.Close()
	var badges []*domain.UserBadge
	for rows.Next() {
		ub := &domain.UserBadge{}
		var earnedAt time.Time
		if err := rows.Scan(&ub.Badge.ID, &ub.Badge.Name, &ub.Badge.Description, &ub.Badge.IconURL, &ub.Badge.Criteria, &earnedAt); err != nil { return nil, err }
		ub.EarnedAt = earnedAt.Format(time.RFC3339)
		badges = append(badges, ub)
	}
	return badges, nil
}

func (r *postgresStudyRepository) AwardBadge(ctx context.Context, userID, badgeID int) error {
	_, err := r.db.Exec(ctx, `INSERT INTO tutora_app.user_badges (user_id, badge_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, userID, badgeID)
	return err
}
