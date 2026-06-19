package domain

import (
	"context"
	"time"
)

type CoinTransaction struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Amount          int       `json:"amount"`
	TransactionType string    `json:"transaction_type"` // 'iap_purchase', 'marketplace_buy', 'marketplace_sale', 'withdrawal'
	ReferenceID     *string   `json:"reference_id,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type StoreItem struct {
	ID           int       `json:"id"`
	SellerID     int       `json:"seller_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Category     string    `json:"category"` // 'notes', 'practice_exams', 'video_course'
	Subject      string    `json:"subject"`
	PriceInCoins int       `json:"price_in_coins"`
	FileURL      string    `json:"file_url,omitempty"` // only visible if purchased or seller
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	SellerName   string    `json:"seller_name,omitempty"`
}

type ItemPurchase struct {
	ID                int       `json:"id"`
	BuyerID           int       `json:"buyer_id"`
	ItemID            int       `json:"item_id"`
	CoinsSpent        int       `json:"coins_spent"`
	CoinsPlatformFee  int       `json:"coins_platform_fee"`
	CoinsSellerAmount int       `json:"coins_seller_amount"`
	CreatedAt         time.Time `json:"created_at"`
}

type Payout struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	CoinsDebited       int       `json:"coins_debited"`
	CashAmountTHB      float64   `json:"cash_amount_thb"`
	BankAccountDetails string    `json:"bank_account_details"`
	Status             string    `json:"status"` // 'requested', 'processed', 'rejected'
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type ItemQAPin struct {
	ID           int            `json:"id"`
	ItemID       int            `json:"item_id"`
	UserID       int            `json:"user_id"`
	PageNumber   int            `json:"page_number"`
	CoordinateX  float64        `json:"coordinate_x"`
	CoordinateY  float64        `json:"coordinate_y"`
	QuestionText string         `json:"question_text"`
	CreatedAt    time.Time      `json:"created_at"`
	UserName     string         `json:"user_name,omitempty"`
	Replies      []*ItemQAReply `json:"replies,omitempty"`
}

type ItemQAReply struct {
	ID        int       `json:"id"`
	PinID     int       `json:"pin_id"`
	UserID    int       `json:"user_id"`
	ReplyText string    `json:"reply_text"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"user_name,omitempty"`
}

type StoreRepository interface {
	// Item operations
	CreateItem(ctx context.Context, item *StoreItem) error
	GetItemByID(ctx context.Context, id int) (*StoreItem, error)
	ListItems(ctx context.Context, category, subject, search string, requesterUserID int) ([]*StoreItem, error)
	
	// Wallet operations
	GetBalance(ctx context.Context, userID int) (int, error)
	AddCoinTransaction(ctx context.Context, tx *CoinTransaction) error
	GetTransactions(ctx context.Context, userID int) ([]*CoinTransaction, error)
	
	// Purchase operations
	CreatePurchase(ctx context.Context, purchase *ItemPurchase, sellerID int) error
	HasPurchased(ctx context.Context, userID, itemID int) (bool, error)
	GetPurchasedItems(ctx context.Context, userID int) ([]*StoreItem, error)
	
	// Payout operations
	CreatePayout(ctx context.Context, payout *Payout) error

	// Interactive Q&A operations
	CreateQAPin(ctx context.Context, pin *ItemQAPin) error
	GetQAPinsByItemID(ctx context.Context, itemID int) ([]*ItemQAPin, error)
	GetQAPinByID(ctx context.Context, pinID int) (*ItemQAPin, error)
	CreateQAReply(ctx context.Context, reply *ItemQAReply) error
	GetQARepliesByPinID(ctx context.Context, pinID int) ([]*ItemQAReply, error)
}

type StoreUsecase interface {
	ListItem(ctx context.Context, item *StoreItem) error
	BrowseItems(ctx context.Context, category, subject, search string, requesterUserID int) ([]*StoreItem, error)
	GetItemDetails(ctx context.Context, userID, itemID int) (*StoreItem, error)
	
	// Wallet and IAP
	GetWalletInfo(ctx context.Context, userID int) (int, []*CoinTransaction, error)
	ValidateAppleReceipt(ctx context.Context, userID int, receipt string) error
	
	// Transactions
	PurchaseItem(ctx context.Context, buyerID, itemID int) error
	GetMyLibrary(ctx context.Context, userID int) ([]*StoreItem, error)
	GetDownloadURL(ctx context.Context, userID, itemID int) (string, error)
	
	// Cash Out
	RequestPayout(ctx context.Context, userID int, coins int, bankAccount string) error

	// Interactive Q&A
	CreateQAPin(ctx context.Context, pin *ItemQAPin) error
	GetQAPins(ctx context.Context, userID, itemID int) ([]*ItemQAPin, error)
	ReplyToQAPin(ctx context.Context, reply *ItemQAReply) error
}
