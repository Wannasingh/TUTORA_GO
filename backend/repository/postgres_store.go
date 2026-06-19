package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Wannasingh/TUTORA_GO/backend/domain"
)

type postgresStoreRepository struct {
	db *pgxpool.Pool
}

func NewPostgresStoreRepository(db *pgxpool.Pool) domain.StoreRepository {
	return &postgresStoreRepository{db: db}
}

func (r *postgresStoreRepository) CreateItem(ctx context.Context, item *domain.StoreItem) error {
	query := `INSERT INTO tutora_app.store_items (seller_id, title, description, category, subject, price_in_coins, file_url, status) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`
	return r.db.QueryRow(ctx, query, item.SellerID, item.Title, item.Description, item.Category, item.Subject, item.PriceInCoins, item.FileURL, item.Status).
		Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
}

func (r *postgresStoreRepository) GetItemByID(ctx context.Context, id int) (*domain.StoreItem, error) {
	query := `SELECT s.id, s.seller_id, s.title, s.description, s.category, s.subject, s.price_in_coins, s.file_url, s.status, s.created_at, s.updated_at, u.name
	          FROM tutora_app.store_items s
	          JOIN tutora_app.users u ON s.seller_id = u.id
	          WHERE s.id = $1`
	item := &domain.StoreItem{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&item.ID, &item.SellerID, &item.Title, &item.Description, &item.Category, &item.Subject, &item.PriceInCoins, &item.FileURL, &item.Status, &item.CreatedAt, &item.UpdatedAt, &item.SellerName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (r *postgresStoreRepository) ListItems(ctx context.Context, category, subject, search string, requesterUserID int) ([]*domain.StoreItem, error) {
	query := `SELECT s.id, s.seller_id, s.title, s.description, s.category, s.subject, s.price_in_coins, s.status, s.created_at, s.updated_at, u.name
	          FROM tutora_app.store_items s
	          JOIN tutora_app.users u ON s.seller_id = u.id
	          WHERE s.status = 'active'`
	
	var args []interface{}
	argCount := 1

	if category != "" {
		query += fmt.Sprintf(" AND s.category = $%d", argCount)
		args = append(args, category)
		argCount++
	}
	if subject != "" {
		query += fmt.Sprintf(" AND s.subject ILIKE $%d", argCount)
		args = append(args, "%"+subject+"%")
		argCount++
	}
	if search != "" {
		query += fmt.Sprintf(" AND (s.title ILIKE $%d OR s.description ILIKE $%d)", argCount, argCount)
		args = append(args, "%"+search+"%")
		argCount++
	}

	var orderClause string
	if requesterUserID > 0 {
		orderClause = fmt.Sprintf(` ORDER BY (
		                     (
		                       (SELECT COUNT(*) FROM tutora_app.item_purchases WHERE item_id = s.id) * 3.0 + 1.0
		                     ) / 
		                     POWER(EXTRACT(EPOCH FROM (NOW() - s.created_at)) / 3600 + 2, 1.5)
		                   ) * (
		                     CASE WHEN (
		                       EXISTS (SELECT 1 FROM tutora_app.tutors WHERE user_id = $%d AND subject ILIKE '%%' || s.subject || '%%')
		                       OR EXISTS (SELECT 1 FROM tutora_app.item_purchases ip2 JOIN tutora_app.store_items s2 ON ip2.item_id = s2.id WHERE ip2.buyer_id = $%d AND s2.subject = s.subject)
		                       OR EXISTS (SELECT 1 FROM tutora_app.post_likes pl JOIN tutora_app.posts p2 ON pl.post_id = p2.id WHERE pl.user_id = $%d AND p2.subject = s.subject)
		                       OR EXISTS (SELECT 1 FROM tutora_app.post_saves ps JOIN tutora_app.posts p3 ON ps.post_id = p3.id WHERE ps.user_id = $%d AND p3.subject = s.subject)
		                     ) THEN 1.5 ELSE 1.0 END
		                   ) DESC, s.created_at DESC`, argCount, argCount, argCount, argCount)
		args = append(args, requesterUserID)
		argCount++
	} else {
		orderClause = ` ORDER BY (
		                    (
		                      (SELECT COUNT(*) FROM tutora_app.item_purchases WHERE item_id = s.id) * 3.0 + 1.0
		                    ) / 
		                    POWER(EXTRACT(EPOCH FROM (NOW() - s.created_at)) / 3600 + 2, 1.5)
		                  ) DESC, s.created_at DESC`
	}

	query += orderClause

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.StoreItem
	for rows.Next() {
		item := &domain.StoreItem{}
		err := rows.Scan(&item.ID, &item.SellerID, &item.Title, &item.Description, &item.Category, &item.Subject, &item.PriceInCoins, &item.Status, &item.CreatedAt, &item.UpdatedAt, &item.SellerName)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *postgresStoreRepository) GetBalance(ctx context.Context, userID int) (int, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM tutora_app.coin_transactions WHERE user_id = $1`
	var balance int
	err := r.db.QueryRow(ctx, query, userID).Scan(&balance)
	return balance, err
}

func (r *postgresStoreRepository) AddCoinTransaction(ctx context.Context, tx *domain.CoinTransaction) error {
	query := `INSERT INTO tutora_app.coin_transactions (user_id, amount, transaction_type, reference_id) 
	          VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, tx.UserID, tx.Amount, tx.TransactionType, tx.ReferenceID).
		Scan(&tx.ID, &tx.CreatedAt)
}

func (r *postgresStoreRepository) GetTransactions(ctx context.Context, userID int) ([]*domain.CoinTransaction, error) {
	query := `SELECT id, user_id, amount, transaction_type, reference_id, created_at 
	          FROM tutora_app.coin_transactions 
	          WHERE user_id = $1 
	          ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []*domain.CoinTransaction
	for rows.Next() {
		tx := &domain.CoinTransaction{}
		err := rows.Scan(&tx.ID, &tx.UserID, &tx.Amount, &tx.TransactionType, &tx.ReferenceID, &tx.CreatedAt)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

func (r *postgresStoreRepository) CreatePurchase(ctx context.Context, purchase *domain.ItemPurchase, sellerID int) error {
	// Execute inside database transaction
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// 1. Check buyer balance
	var balance int
	balanceQuery := `SELECT COALESCE(SUM(amount), 0) FROM tutora_app.coin_transactions WHERE user_id = $1`
	err = tx.QueryRow(ctx, balanceQuery, purchase.BuyerID).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < purchase.CoinsSpent {
		return errors.New("insufficient coin balance")
	}

	// 2. Insert purchase record
	purchaseQuery := `INSERT INTO tutora_app.item_purchases (buyer_id, item_id, coins_spent, coins_platform_fee, coins_seller_amount) 
	                  VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err = tx.QueryRow(ctx, purchaseQuery, purchase.BuyerID, purchase.ItemID, purchase.CoinsSpent, purchase.CoinsPlatformFee, purchase.CoinsSellerAmount).
		Scan(&purchase.ID, &purchase.CreatedAt)
	if err != nil {
		return err
	}

	purchaseRef := fmt.Sprintf("purchase_%d", purchase.ID)

	// 3. Insert negative coin transaction for buyer
	buyerTxQuery := `INSERT INTO tutora_app.coin_transactions (user_id, amount, transaction_type, reference_id) 
	                 VALUES ($1, $2, 'marketplace_buy', $3)`
	_, err = tx.Exec(ctx, buyerTxQuery, purchase.BuyerID, -purchase.CoinsSpent, purchaseRef)
	if err != nil {
		return err
	}

	// 4. Insert positive coin transaction for seller
	sellerTxQuery := `INSERT INTO tutora_app.coin_transactions (user_id, amount, transaction_type, reference_id) 
	                  VALUES ($1, $2, 'marketplace_sale', $3)`
	_, err = tx.Exec(ctx, sellerTxQuery, sellerID, purchase.CoinsSellerAmount, purchaseRef)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *postgresStoreRepository) HasPurchased(ctx context.Context, userID, itemID int) (bool, error) {
	// Check if user is seller of the item
	sellerCheck := `SELECT EXISTS(SELECT 1 FROM tutora_app.store_items WHERE id = $1 AND seller_id = $2)`
	var isSeller bool
	err := r.db.QueryRow(ctx, sellerCheck, itemID, userID).Scan(&isSeller)
	if err != nil {
		return false, err
	}
	if isSeller {
		return true, nil
	}

	// Check if user has purchased the item
	purchaseCheck := `SELECT EXISTS(SELECT 1 FROM tutora_app.item_purchases WHERE item_id = $1 AND buyer_id = $2)`
	var isBuyer bool
	err = r.db.QueryRow(ctx, purchaseCheck, itemID, userID).Scan(&isBuyer)
	return isBuyer, err
}

func (r *postgresStoreRepository) GetPurchasedItems(ctx context.Context, userID int) ([]*domain.StoreItem, error) {
	query := `SELECT s.id, s.seller_id, s.title, s.description, s.category, s.subject, s.price_in_coins, s.status, s.created_at, s.updated_at, u.name
	          FROM tutora_app.item_purchases ip
	          JOIN tutora_app.store_items s ON ip.item_id = s.id
	          JOIN tutora_app.users u ON s.seller_id = u.id
	          WHERE ip.buyer_id = $1
	          ORDER BY ip.created_at DESC`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*domain.StoreItem
	for rows.Next() {
		item := &domain.StoreItem{}
		err := rows.Scan(&item.ID, &item.SellerID, &item.Title, &item.Description, &item.Category, &item.Subject, &item.PriceInCoins, &item.Status, &item.CreatedAt, &item.UpdatedAt, &item.SellerName)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *postgresStoreRepository) CreatePayout(ctx context.Context, payout *domain.Payout) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Check balance
	var balance int
	balanceQuery := `SELECT COALESCE(SUM(amount), 0) FROM tutora_app.coin_transactions WHERE user_id = $1`
	err = tx.QueryRow(ctx, balanceQuery, payout.UserID).Scan(&balance)
	if err != nil {
		return err
	}

	if balance < payout.CoinsDebited {
		return errors.New("insufficient coin balance for payout")
	}

	// Insert payout record
	payoutQuery := `INSERT INTO tutora_app.payouts (user_id, coins_debited, cash_amount_thb, bank_account_details, status) 
	                VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`
	err = tx.QueryRow(ctx, payoutQuery, payout.UserID, payout.CoinsDebited, payout.CashAmountTHB, payout.BankAccountDetails, payout.Status).
		Scan(&payout.ID, &payout.CreatedAt, &payout.UpdatedAt)
	if err != nil {
		return err
	}

	// Debit coins from user balance
	payoutRef := fmt.Sprintf("payout_%d", payout.ID)
	txQuery := `INSERT INTO tutora_app.coin_transactions (user_id, amount, transaction_type, reference_id) 
	            VALUES ($1, $2, 'withdrawal', $3)`
	_, err = tx.Exec(ctx, txQuery, payout.UserID, -payout.CoinsDebited, payoutRef)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *postgresStoreRepository) CreateQAPin(ctx context.Context, pin *domain.ItemQAPin) error {
	query := `INSERT INTO tutora_app.item_qa_pins (item_id, user_id, page_number, coordinate_x, coordinate_y, question_text) 
	          VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, pin.ItemID, pin.UserID, pin.PageNumber, pin.CoordinateX, pin.CoordinateY, pin.QuestionText).
		Scan(&pin.ID, &pin.CreatedAt)
}

func (r *postgresStoreRepository) GetQAPinByID(ctx context.Context, pinID int) (*domain.ItemQAPin, error) {
	query := `SELECT id, item_id, user_id, page_number, coordinate_x, coordinate_y, question_text, created_at 
	          FROM tutora_app.item_qa_pins WHERE id = $1`
	pin := &domain.ItemQAPin{}
	err := r.db.QueryRow(ctx, query, pinID).
		Scan(&pin.ID, &pin.ItemID, &pin.UserID, &pin.PageNumber, &pin.CoordinateX, &pin.CoordinateY, &pin.QuestionText, &pin.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return pin, nil
}

func (r *postgresStoreRepository) GetQAPinsByItemID(ctx context.Context, itemID int) ([]*domain.ItemQAPin, error) {
	query := `SELECT q.id, q.item_id, q.user_id, q.page_number, q.coordinate_x, q.coordinate_y, q.question_text, q.created_at, u.name 
	          FROM tutora_app.item_qa_pins q
	          JOIN tutora_app.users u ON q.user_id = u.id
	          WHERE q.item_id = $1
	          ORDER BY q.page_number ASC, q.created_at ASC`
	rows, err := r.db.Query(ctx, query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pins []*domain.ItemQAPin
	for rows.Next() {
		pin := &domain.ItemQAPin{}
		err := rows.Scan(&pin.ID, &pin.ItemID, &pin.UserID, &pin.PageNumber, &pin.CoordinateX, &pin.CoordinateY, &pin.QuestionText, &pin.CreatedAt, &pin.UserName)
		if err != nil {
			return nil, err
		}
		pins = append(pins, pin)
	}
	return pins, nil
}

func (r *postgresStoreRepository) CreateQAReply(ctx context.Context, reply *domain.ItemQAReply) error {
	query := `INSERT INTO tutora_app.item_qa_replies (pin_id, user_id, reply_text) 
	          VALUES ($1, $2, $3) RETURNING id, created_at`
	return r.db.QueryRow(ctx, query, reply.PinID, reply.UserID, reply.ReplyText).
		Scan(&reply.ID, &reply.CreatedAt)
}

func (r *postgresStoreRepository) GetQARepliesByPinID(ctx context.Context, pinID int) ([]*domain.ItemQAReply, error) {
	query := `SELECT r.id, r.pin_id, r.user_id, r.reply_text, r.created_at, u.name 
	          FROM tutora_app.item_qa_replies r
	          JOIN tutora_app.users u ON r.user_id = u.id
	          WHERE r.pin_id = $1
	          ORDER BY r.created_at ASC`
	rows, err := r.db.Query(ctx, query, pinID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []*domain.ItemQAReply
	for rows.Next() {
		reply := &domain.ItemQAReply{}
		err := rows.Scan(&reply.ID, &reply.PinID, &reply.UserID, &reply.ReplyText, &reply.CreatedAt, &reply.UserName)
		if err != nil {
			return nil, err
		}
		replies = append(replies, reply)
	}
	return replies, nil
}
