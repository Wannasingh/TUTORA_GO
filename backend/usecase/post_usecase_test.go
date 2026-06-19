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

func (m *MockPostRepository) GetComments(ctx context.Context, postID int, requesterUserID int) ([]*domain.Comment, error) {
	return m.comments[postID], nil
}

func (m *MockPostRepository) LikeComment(ctx context.Context, commentID int, userID int) (bool, error) {
	return true, nil
}

func (m *MockPostRepository) DeleteComment(ctx context.Context, id int) error {
	return nil
}

func (m *MockPostRepository) GetCommentByID(ctx context.Context, id int) (*domain.Comment, error) {
	for _, commentList := range m.comments {
		for _, c := range commentList {
			if c.ID == id {
				return c, nil
			}
		}
	}
	return nil, nil
}

func (m *MockPostRepository) CreateReport(ctx context.Context, report *domain.Report) error {
	return nil
}

func (m *MockPostRepository) Update(ctx context.Context, post *domain.Post) error {
	m.posts[post.ID] = post
	return nil
}

func (m *MockPostRepository) Delete(ctx context.Context, id int) error {
	delete(m.posts, id)
	return nil
}

func (m *MockPostRepository) Repost(ctx context.Context, postID int, userID int) (bool, error) {
	return true, nil
}

func (m *MockPostRepository) GetUserPosts(ctx context.Context, userID, requesterUserID int) ([]*domain.Post, error) {
	return nil, nil
}

func (m *MockPostRepository) GetUserLikedPosts(ctx context.Context, userID, requesterUserID int) ([]*domain.Post, error) {
	return nil, nil
}

func (m *MockPostRepository) GetUserSavedPosts(ctx context.Context, userID, requesterUserID int) ([]*domain.Post, error) {
	return nil, nil
}

func (m *MockPostRepository) GetUserRepostedPosts(ctx context.Context, userID, requesterUserID int) ([]*domain.Post, error) {
	return nil, nil
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

func TestGetPostDetails_NestedComments(t *testing.T) {
	postID := 5
	existingPost := &domain.Post{ID: postID, Title: "Nesting test"}
	parentID := 1
	c1 := &domain.Comment{ID: parentID, PostID: postID, UserID: 10, Body: "Parent comment"}
	c2 := &domain.Comment{ID: 2, PostID: postID, UserID: 11, Body: "Child reply", ParentID: &parentID}

	repo := &MockPostRepository{
		posts: map[int]*domain.Post{
			postID: existingPost,
		},
		comments: map[int][]*domain.Comment{
			postID: {c1, c2},
		},
	}
	u := NewPostUsecase(repo)

	_, rootComments, err := u.GetPostDetails(context.Background(), postID, 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(rootComments) != 1 {
		t.Errorf("expected only 1 root comment, got %d", len(rootComments))
	}

	if len(rootComments[0].Replies) != 1 {
		t.Errorf("expected 1 reply under root comment, got %d", len(rootComments[0].Replies))
	}

	if rootComments[0].Replies[0].ID != 2 {
		t.Errorf("expected child comment ID 2, got %d", rootComments[0].Replies[0].ID)
	}
}

func TestAddPostComment_NestedValidation(t *testing.T) {
	repo := &MockPostRepository{
		comments: make(map[int][]*domain.Comment),
	}
	u := NewPostUsecase(repo)

	// Add parent comment first
	parent := &domain.Comment{ID: 10, PostID: 2, UserID: 3, Body: "Original"}
	repo.comments[2] = append(repo.comments[2], parent)

	// Try replying to a non-existent parent comment
	nonExistentParentID := 99
	replyBad := &domain.Comment{
		PostID:   2,
		UserID:   4,
		Body:     "Bad reply",
		ParentID: &nonExistentParentID,
	}
	err := u.AddPostComment(context.Background(), replyBad)
	if err == nil {
		t.Error("expected error for non-existent parent reply, got nil")
	}

	// Reply to valid parent comment (should succeed)
	validParentID := 10
	replyGood := &domain.Comment{
		PostID:   2,
		UserID:   4,
		Body:     "Good reply",
		ParentID: &validParentID,
	}
	err = u.AddPostComment(context.Background(), replyGood)
	if err != nil {
		t.Fatalf("failed to add valid reply: %v", err)
	}
}

func TestDeleteComment_Authorization(t *testing.T) {
	repo := &MockPostRepository{
		comments: make(map[int][]*domain.Comment),
	}
	u := NewPostUsecase(repo)

	c := &domain.Comment{ID: 50, PostID: 1, UserID: 10, Body: "My comment"}
	repo.comments[1] = append(repo.comments[1], c)

	// Unauthorized delete (should fail)
	err := u.DeleteComment(context.Background(), 99, 50)
	if err == nil {
		t.Error("expected error for unauthorized deletion, got nil")
	}

	// Authorized delete (should succeed)
	err = u.DeleteComment(context.Background(), 10, 50)
	if err != nil {
		t.Fatalf("expected delete to succeed, got %v", err)
	}
}

func TestUpdatePost_Success(t *testing.T) {
	existingPost := &domain.Post{
		ID:     5,
		UserID: 10,
		Title:  "Original Title",
		Body:   "Original Body",
	}
	repo := &MockPostRepository{
		posts: map[int]*domain.Post{
			5: existingPost,
		},
	}
	u := NewPostUsecase(repo)

	updatedPost := &domain.Post{
		ID:     5,
		UserID: 10, // match author
		Title:  "Updated Title",
		Body:   "Updated Body",
	}

	err := u.UpdatePost(context.Background(), 10, updatedPost)
	if err != nil {
		t.Fatalf("expected update to succeed, got %v", err)
	}

	if repo.posts[5].Title != "Updated Title" {
		t.Errorf("expected title to be 'Updated Title', got %s", repo.posts[5].Title)
	}
}

func TestUpdatePost_Unauthorized(t *testing.T) {
	existingPost := &domain.Post{
		ID:     5,
		UserID: 10,
		Title:  "Original Title",
		Body:   "Original Body",
	}
	repo := &MockPostRepository{
		posts: map[int]*domain.Post{
			5: existingPost,
		},
	}
	u := NewPostUsecase(repo)

	updatedPost := &domain.Post{
		ID:     5,
		UserID: 10,
		Title:  "Updated Title",
		Body:   "Updated Body",
	}

	// User 99 tries to update User 10's post
	err := u.UpdatePost(context.Background(), 99, updatedPost)
	if err == nil {
		t.Error("expected error for unauthorized update, got nil")
	}
}

func TestDeletePost_Success(t *testing.T) {
	existingPost := &domain.Post{
		ID:     5,
		UserID: 10,
		Title:  "Original Title",
	}
	repo := &MockPostRepository{
		posts: map[int]*domain.Post{
			5: existingPost,
		},
	}
	u := NewPostUsecase(repo)

	err := u.DeletePost(context.Background(), 10, 5)
	if err != nil {
		t.Fatalf("expected deletion to succeed, got %v", err)
	}

	if _, ok := repo.posts[5]; ok {
		t.Error("expected post to be deleted from repo")
	}
}

func TestDeletePost_Unauthorized(t *testing.T) {
	existingPost := &domain.Post{
		ID:     5,
		UserID: 10,
		Title:  "Original Title",
	}
	repo := &MockPostRepository{
		posts: map[int]*domain.Post{
			5: existingPost,
		},
	}
	u := NewPostUsecase(repo)

	// User 99 tries to delete User 10's post
	err := u.DeletePost(context.Background(), 99, 5)
	if err == nil {
		t.Error("expected error for unauthorized deletion, got nil")
	}

	if _, ok := repo.posts[5]; !ok {
		t.Error("expected post to remain in repo")
	}
}


