package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type storeUsecase struct {
	storeRepo domain.StoreRepository
}

func NewStoreUsecase(repo domain.StoreRepository) domain.StoreUsecase {
	return &storeUsecase{storeRepo: repo}
}

func (u *storeUsecase) ListItem(ctx context.Context, item *domain.StoreItem) error {
	if item.PriceInCoins <= 0 {
		return errors.New("price must be greater than 0 coins")
	}
	if item.Title == "" || item.FileURL == "" {
		return errors.New("title and file URL are required")
	}
	item.Status = "active"
	return u.storeRepo.CreateItem(ctx, item)
}

func (u *storeUsecase) BrowseItems(ctx context.Context, category, subject, search string, requesterUserID int) ([]*domain.StoreItem, error) {
	return u.storeRepo.ListItems(ctx, category, subject, search, requesterUserID)
}

func (u *storeUsecase) GetItemDetails(ctx context.Context, userID, itemID int) (*domain.StoreItem, error) {
	item, err := u.storeRepo.GetItemByID(ctx, itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}

	hasPurchased, err := u.storeRepo.HasPurchased(ctx, userID, itemID)
	if err != nil {
		return nil, err
	}

	// Protect file URL from unauthorized users by copying the item structure first
	itemCopy := *item
	if !hasPurchased {
		itemCopy.FileURL = ""
	}
	return &itemCopy, nil
}

func (u *storeUsecase) GetWalletInfo(ctx context.Context, userID int) (int, []*domain.CoinTransaction, error) {
	balance, err := u.storeRepo.GetBalance(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	txs, err := u.storeRepo.GetTransactions(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	return balance, txs, nil
}

func (u *storeUsecase) ValidateAppleReceipt(ctx context.Context, userID int, receipt string) error {
	if receipt == "" {
		return errors.New("receipt is required")
	}

	// Mock validating receipt with Apple App Store API.
	// In production, we would call https://buy.itunes.apple.com/verifyReceipt
	var coins int
	var refID string

	if strings.Contains(receipt, "coin_50") {
		coins = 50
		refID = "apple_ref_mock_50_" + receipt[len(receipt)-6:]
	} else if strings.Contains(receipt, "coin_100") {
		coins = 100
		refID = "apple_ref_mock_100_" + receipt[len(receipt)-6:]
	} else if strings.Contains(receipt, "coin_500") {
		coins = 500
		refID = "apple_ref_mock_500_" + receipt[len(receipt)-6:]
	} else {
		// Default fallback for general testing
		coins = 100
		refID = "apple_ref_mock_gen_" + fmt.Sprintf("%d", len(receipt))
	}

	tx := &domain.CoinTransaction{
		UserID:          userID,
		Amount:          coins,
		TransactionType: "iap_purchase",
		ReferenceID:     &refID,
	}

	return u.storeRepo.AddCoinTransaction(ctx, tx)
}

func (u *storeUsecase) PurchaseItem(ctx context.Context, buyerID, itemID int) error {
	item, err := u.storeRepo.GetItemByID(ctx, itemID)
	if err != nil {
		return err
	}
	if item == nil {
		return errors.New("item not found")
	}
	if item.SellerID == buyerID {
		return errors.New("cannot purchase your own item")
	}

	alreadyPurchased, err := u.storeRepo.HasPurchased(ctx, buyerID, itemID)
	if err != nil {
		return err
	}
	if alreadyPurchased {
		return errors.New("you have already purchased this item")
	}

	// Calculate 15% platform commission
	fee := int(float64(item.PriceInCoins) * 0.15)
	sellerAmt := item.PriceInCoins - fee

	purchase := &domain.ItemPurchase{
		BuyerID:           buyerID,
		ItemID:            itemID,
		CoinsSpent:        item.PriceInCoins,
		CoinsPlatformFee:  fee,
		CoinsSellerAmount: sellerAmt,
	}

	return u.storeRepo.CreatePurchase(ctx, purchase, item.SellerID)
}

func (u *storeUsecase) GetMyLibrary(ctx context.Context, userID int) ([]*domain.StoreItem, error) {
	return u.storeRepo.GetPurchasedItems(ctx, userID)
}

func (u *storeUsecase) GetDownloadURL(ctx context.Context, userID, itemID int) (string, error) {
	hasPurchased, err := u.storeRepo.HasPurchased(ctx, userID, itemID)
	if err != nil {
		return "", err
	}
	if !hasPurchased {
		return "", errors.New("access denied: purchase required")
	}

	item, err := u.storeRepo.GetItemByID(ctx, itemID)
	if err != nil {
		return "", err
	}
	if item == nil {
		return "", errors.New("item not found")
	}

	return item.FileURL, nil
}

func (u *storeUsecase) RequestPayout(ctx context.Context, userID int, coins int, bankAccount string) error {
	if coins <= 0 {
		return errors.New("withdrawal amount must be greater than 0")
	}
	if bankAccount == "" {
		return errors.New("bank account details are required")
	}

	// 1 Coin = 0.70 THB (accounting for Apple's 30% IAP commission)
	cashAmt := float64(coins) * 0.70

	payout := &domain.Payout{
		UserID:             userID,
		CoinsDebited:       coins,
		CashAmountTHB:      cashAmt,
		BankAccountDetails: bankAccount,
		Status:             "requested",
	}

	return u.storeRepo.CreatePayout(ctx, payout)
}

func (u *storeUsecase) CreateQAPin(ctx context.Context, pin *domain.ItemQAPin) error {
	hasPurchased, err := u.storeRepo.HasPurchased(ctx, pin.UserID, pin.ItemID)
	if err != nil {
		return err
	}
	if !hasPurchased {
		return errors.New("access denied: purchase required to ask questions")
	}
	if pin.QuestionText == "" {
		return errors.New("question text cannot be empty")
	}
	return u.storeRepo.CreateQAPin(ctx, pin)
}

func (u *storeUsecase) GetQAPins(ctx context.Context, userID, itemID int) ([]*domain.ItemQAPin, error) {
	hasPurchased, err := u.storeRepo.HasPurchased(ctx, userID, itemID)
	if err != nil {
		return nil, err
	}
	if !hasPurchased {
		return nil, errors.New("access denied: purchase required to view Q&A")
	}

	pins, err := u.storeRepo.GetQAPinsByItemID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	// Fetch replies for each pin
	for _, pin := range pins {
		replies, err := u.storeRepo.GetQARepliesByPinID(ctx, pin.ID)
		if err != nil {
			return nil, err
		}
		pin.Replies = replies
	}

	return pins, nil
}

func (u *storeUsecase) ReplyToQAPin(ctx context.Context, reply *domain.ItemQAReply) error {
	if reply.ReplyText == "" {
		return errors.New("reply text cannot be empty")
	}

	pin, err := u.storeRepo.GetQAPinByID(ctx, reply.PinID)
	if err != nil {
		return err
	}
	if pin == nil {
		return errors.New("Q&A Pin not found")
	}

	hasPurchased, err := u.storeRepo.HasPurchased(ctx, reply.UserID, pin.ItemID)
	if err != nil {
		return err
	}
	if !hasPurchased {
		return errors.New("access denied: purchase required to reply to questions")
	}

	return u.storeRepo.CreateQAReply(ctx, reply)
}
