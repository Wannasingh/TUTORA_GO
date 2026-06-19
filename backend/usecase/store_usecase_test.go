package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type mockStoreRepository struct {
	items        map[int]*domain.StoreItem
	purchases    map[string]*domain.ItemPurchase
	transactions map[int][]*domain.CoinTransaction
	payouts      []*domain.Payout
	qaPins       map[int]*domain.ItemQAPin
	qaReplies    map[int][]*domain.ItemQAReply
}

func newMockStoreRepository() *mockStoreRepository {
	return &mockStoreRepository{
		items:        make(map[int]*domain.StoreItem),
		purchases:    make(map[string]*domain.ItemPurchase),
		transactions: make(map[int][]*domain.CoinTransaction),
		payouts:      make([]*domain.Payout, 0),
		qaPins:       make(map[int]*domain.ItemQAPin),
		qaReplies:    make(map[int][]*domain.ItemQAReply),
	}
}

func (m *mockStoreRepository) CreateItem(ctx context.Context, item *domain.StoreItem) error {
	item.ID = len(m.items) + 1
	m.items[item.ID] = item
	return nil
}

func (m *mockStoreRepository) GetItemByID(ctx context.Context, id int) (*domain.StoreItem, error) {
	return m.items[id], nil
}

func (m *mockStoreRepository) ListItems(ctx context.Context, category, subject, search string, requesterUserID int) ([]*domain.StoreItem, error) {
	var result []*domain.StoreItem
	for _, item := range m.items {
		result = append(result, item)
	}
	return result, nil
}

func (m *mockStoreRepository) GetBalance(ctx context.Context, userID int) (int, error) {
	balance := 0
	for _, tx := range m.transactions[userID] {
		balance += tx.Amount
	}
	return balance, nil
}

func (m *mockStoreRepository) AddCoinTransaction(ctx context.Context, tx *domain.CoinTransaction) error {
	m.transactions[tx.UserID] = append(m.transactions[tx.UserID], tx)
	return nil
}

func (m *mockStoreRepository) GetTransactions(ctx context.Context, userID int) ([]*domain.CoinTransaction, error) {
	return m.transactions[userID], nil
}

func (m *mockStoreRepository) CreatePurchase(ctx context.Context, purchase *domain.ItemPurchase, sellerID int) error {
	balance, _ := m.GetBalance(ctx, purchase.BuyerID)
	if balance < purchase.CoinsSpent {
		return errors.New("insufficient coin balance")
	}

	key := fmt.Sprintf("%d-%d", purchase.BuyerID, purchase.ItemID)
	m.purchases[key] = purchase

	// Add debit transaction for buyer
	buyerTx := &domain.CoinTransaction{
		UserID:          purchase.BuyerID,
		Amount:          -purchase.CoinsSpent,
		TransactionType: "marketplace_buy",
	}
	m.transactions[purchase.BuyerID] = append(m.transactions[purchase.BuyerID], buyerTx)

	// Add credit transaction for seller
	sellerTx := &domain.CoinTransaction{
		UserID:          sellerID,
		Amount:          purchase.CoinsSellerAmount,
		TransactionType: "marketplace_sale",
	}
	m.transactions[sellerID] = append(m.transactions[sellerID], sellerTx)

	return nil
}

func (m *mockStoreRepository) HasPurchased(ctx context.Context, userID, itemID int) (bool, error) {
	item := m.items[itemID]
	if item != nil && item.SellerID == userID {
		return true, nil
	}
	key := fmt.Sprintf("%d-%d", userID, itemID)
	_, found := m.purchases[key]
	return found, nil
}

func (m *mockStoreRepository) GetPurchasedItems(ctx context.Context, userID int) ([]*domain.StoreItem, error) {
	var result []*domain.StoreItem
	for key, purchase := range m.purchases {
		if strings.HasPrefix(key, fmt.Sprintf("%d-", userID)) {
			result = append(result, m.items[purchase.ItemID])
		}
	}
	return result, nil
}

func (m *mockStoreRepository) CreatePayout(ctx context.Context, payout *domain.Payout) error {
	balance, _ := m.GetBalance(ctx, payout.UserID)
	if balance < payout.CoinsDebited {
		return errors.New("insufficient coin balance for payout")
	}
	m.payouts = append(m.payouts, payout)
	tx := &domain.CoinTransaction{
		UserID:          payout.UserID,
		Amount:          -payout.CoinsDebited,
		TransactionType: "withdrawal",
	}
	m.transactions[payout.UserID] = append(m.transactions[payout.UserID], tx)
	return nil
}

func (m *mockStoreRepository) CreateQAPin(ctx context.Context, pin *domain.ItemQAPin) error {
	pin.ID = len(m.qaPins) + 1
	m.qaPins[pin.ID] = pin
	return nil
}

func (m *mockStoreRepository) GetQAPinsByItemID(ctx context.Context, itemID int) ([]*domain.ItemQAPin, error) {
	var result []*domain.ItemQAPin
	for _, pin := range m.qaPins {
		if pin.ItemID == itemID {
			result = append(result, pin)
		}
	}
	return result, nil
}

func (m *mockStoreRepository) GetQAPinByID(ctx context.Context, pinID int) (*domain.ItemQAPin, error) {
	return m.qaPins[pinID], nil
}

func (m *mockStoreRepository) CreateQAReply(ctx context.Context, reply *domain.ItemQAReply) error {
	reply.ID = len(m.qaReplies) + 1
	m.qaReplies[reply.PinID] = append(m.qaReplies[reply.PinID], reply)
	return nil
}

func (m *mockStoreRepository) GetQARepliesByPinID(ctx context.Context, pinID int) ([]*domain.ItemQAReply, error) {
	return m.qaReplies[pinID], nil
}

func TestStoreUsecase_ListItem(t *testing.T) {
	repo := newMockStoreRepository()
	uc := NewStoreUsecase(repo)

	item := &domain.StoreItem{
		Title:        "Math Guide",
		FileURL:      "https://example.com/math.pdf",
		PriceInCoins: 100,
	}

	err := uc.ListItem(context.Background(), item)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if item.ID != 1 {
		t.Errorf("expected item ID to be 1, got %d", item.ID)
	}
}

func TestStoreUsecase_GetItemDetails_Security(t *testing.T) {
	repo := newMockStoreRepository()
	uc := NewStoreUsecase(repo)

	sellerID := 1
	buyerID := 2
	unauthorizedID := 3

	item := &domain.StoreItem{
		SellerID:     sellerID,
		Title:        "Physics Cheat Sheet",
		FileURL:      "https://example.com/phys.pdf",
		PriceInCoins: 50,
	}
	_ = repo.CreateItem(context.Background(), item)

	// Seller details (should see file URL)
	detailsSeller, err := uc.GetItemDetails(context.Background(), sellerID, item.ID)
	if err != nil || detailsSeller == nil || detailsSeller.FileURL == "" {
		t.Errorf("seller should view file URL, got err: %v", err)
	}

	// Unauthorized details (should not see file URL)
	detailsUnauth, err := uc.GetItemDetails(context.Background(), unauthorizedID, item.ID)
	if err != nil || detailsUnauth == nil || detailsUnauth.FileURL != "" {
		t.Errorf("unauthorized user should not view file URL, got FileURL: '%s'", detailsUnauth.FileURL)
	}

	// Buy and check buyer details
	_ = repo.AddCoinTransaction(context.Background(), &domain.CoinTransaction{UserID: buyerID, Amount: 100})
	err = uc.PurchaseItem(context.Background(), buyerID, item.ID)
	if err != nil {
		t.Fatalf("purchase failed: %v", err)
	}

	detailsBuyer, err := uc.GetItemDetails(context.Background(), buyerID, item.ID)
	if err != nil || detailsBuyer == nil || detailsBuyer.FileURL == "" {
		t.Errorf("buyer should view file URL, got err: %v", err)
	}
}

func TestStoreUsecase_ValidateAppleReceipt(t *testing.T) {
	repo := newMockStoreRepository()
	uc := NewStoreUsecase(repo)
	userID := 10

	err := uc.ValidateAppleReceipt(context.Background(), userID, "mock_receipt_coin_100_abcdef")
	if err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	bal, _ := repo.GetBalance(context.Background(), userID)
	if bal != 100 {
		t.Errorf("expected balance to be 100, got %d", bal)
	}
}

func TestStoreUsecase_RequestPayout(t *testing.T) {
	repo := newMockStoreRepository()
	uc := NewStoreUsecase(repo)
	userID := 15

	// Pre-credit user with coins
	_ = repo.AddCoinTransaction(context.Background(), &domain.CoinTransaction{UserID: userID, Amount: 200})

	err := uc.RequestPayout(context.Background(), userID, 100, "SCB 123-456-7890")
	if err != nil {
		t.Fatalf("payout request failed: %v", err)
	}

	bal, _ := repo.GetBalance(context.Background(), userID)
	if bal != 100 {
		t.Errorf("expected balance to be 100 after debit, got %d", bal)
	}
}

func TestStoreUsecase_InteractiveQA(t *testing.T) {
	repo := newMockStoreRepository()
	uc := NewStoreUsecase(repo)

	sellerID := 1
	buyerID := 2
	thiefID := 3

	item := &domain.StoreItem{
		SellerID:     sellerID,
		Title:        "Chemistry Notebook",
		FileURL:      "https://example.com/chem.pdf",
		PriceInCoins: 80,
	}
	_ = repo.CreateItem(context.Background(), item)

	// Buyer purchases it
	_ = repo.AddCoinTransaction(context.Background(), &domain.CoinTransaction{UserID: buyerID, Amount: 100})
	_ = uc.PurchaseItem(context.Background(), buyerID, item.ID)

	// Thief tries to ask Q&A (should fail)
	badPin := &domain.ItemQAPin{
		ItemID:       item.ID,
		UserID:       thiefID,
		PageNumber:   1,
		QuestionText: "Is this on the exam?",
	}
	err := uc.CreateQAPin(context.Background(), badPin)
	if err == nil {
		t.Error("expected error for non-purchased QA Pin, got nil")
	}

	// Buyer creates Q&A Pin (should succeed)
	goodPin := &domain.ItemQAPin{
		ItemID:       item.ID,
		UserID:       buyerID,
		PageNumber:   1,
		QuestionText: "Is this formula balanced?",
	}
	err = uc.CreateQAPin(context.Background(), goodPin)
	if err != nil {
		t.Fatalf("failed to create QA Pin: %v", err)
	}

	if goodPin.ID != 1 {
		t.Errorf("expected Pin ID to be 1, got %d", goodPin.ID)
	}

	// Seller replies to the Q&A Pin (should succeed since they are the seller)
	reply := &domain.ItemQAReply{
		PinID:     goodPin.ID,
		UserID:    sellerID,
		ReplyText: "Yes, double checked it.",
	}
	err = uc.ReplyToQAPin(context.Background(), reply)
	if err != nil {
		t.Fatalf("failed to reply: %v", err)
	}

	// Retrieve pins and replies
	pins, err := uc.GetQAPins(context.Background(), buyerID, item.ID)
	if err != nil {
		t.Fatalf("failed to fetch pins: %v", err)
	}

	if len(pins) != 1 {
		t.Errorf("expected 1 pin, got %d", len(pins))
	}

	if len(pins[0].Replies) != 1 {
		t.Errorf("expected 1 reply, got %d", len(pins[0].Replies))
	}
}
