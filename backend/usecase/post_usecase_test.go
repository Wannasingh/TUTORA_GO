package usecase

import (
	"context"
	"testing"

	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type MockPostRepository struct {
	posts       map[int]*domain.Post
	comments    map[int][]*domain.Comment
	likes       map[string]bool
	saves       map[string]bool
	createFunc  func(post *domain.Post) error
	getByIDFunc func(id int, reqID int) (*domain.Post, error)
}

func (m *MockPostRepository) Create(ctx context.Context, post *domain.Post) error {
	if m.createFunc != nil {
		return m.createFunc(post)
	}
	post.ID = len(m.posts) + 1
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) GetByID(ctx context.Context, id int, requesterUserID int) (*domain.Post, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(id, requesterUserID)
	}
	return m.posts[id], nil
}

func (m *MockPostRepository) List(ctx context.Context, subject string, requesterUserID int) ([]*domain.Post, error) {
	var list []*domain.Post
	for _, p := range m.posts {
		list = append(list, p)
	}
	return list, nil
}

func (m *MockPostRepository) Like(ctx context.Context, postID int, userID int) (bool, error) {
	return true, nil
}

func (m *MockPostRepository) Save(ctx context.Context, postID int, userID int) (bool, error) {
	return true, nil
}

func (m *MockPostRepository) AddComment(ctx context.Context, comment *domain.Comment) error {
	comment.ID = len(m.comments[comment.PostID]) + 1
	m.comments[comment.PostID] = append(m.comments[comment.PostID], comment)
	return nil
}

func (m *MockPostRepository) GetComments(ctx context.Context, postID int) ([]*domain.Comment, error) {
	return m.comments[postID], nil
}

func TestCreatePost_Success(t *testing.T) {
	repo := &MockPostRepository{posts: make(map[int]*domain.Post)}
	u := NewPostUsecase(repo)

	post := &domain.Post{
		UserID:  1,
		Subject: "Math",
		Title:   "Limits",
		Body:    "Body",
	}

	err := u.CreatePost(context.Background(), post)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if post.ID != 1 {
		t.Errorf("expected ID to be 1, got %d", post.ID)
	}
}

func TestGetPostDetails_Success(t *testing.T) {
	existingPost := &domain.Post{
		ID:      2,
		UserID:  10,
		Subject: "Physics",
		Title:   "Friction",
		Body:    "Friction details",
	}
	repo := &MockPostRepository{
		posts: map[int]*domain.Post{
			2: existingPost,
		},
		comments: map[int][]*domain.Comment{
			2: {
				{ID: 1, PostID: 2, UserID: 3, Body: "Great explainer!"},
			},
		},
	}
	u := NewPostUsecase(repo)

	post, comments, err := u.GetPostDetails(context.Background(), 2, 3)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if post.ID != 2 {
		t.Errorf("expected post ID 2, got %d", post.ID)
	}

	if len(comments) != 1 || comments[0].Body != "Great explainer!" {
		t.Errorf("unexpected comment details: %v", comments)
	}
}

func TestToggleLike_Success(t *testing.T) {
	repo := &MockPostRepository{}
	u := NewPostUsecase(repo)

	liked, err := u.ToggleLike(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !liked {
		t.Error("expected liked status to be true")
	}
}

func TestToggleSave_Success(t *testing.T) {
	repo := &MockPostRepository{}
	u := NewPostUsecase(repo)

	saved, err := u.ToggleSave(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !saved {
		t.Error("expected saved status to be true")
	}
}
