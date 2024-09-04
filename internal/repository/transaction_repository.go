package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/kodinggo/payment-service-gb1/internal/model"
	"github.com/redis/go-redis/v9"
)

type transactionRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewTransactionRepository(db *sql.DB, redis *redis.Client) model.ITransactionRepository {
	return &transactionRepository{
		db:    db,
		redis: redis,
	}
}

func (p *transactionRepository) FindAll(ctx context.Context, filter model.TransactionFilter) ([]*model.Transaction, error) {
	res, err := p.db.QueryContext(ctx, "SELECT id, payment_id, status, created_at, updated_at, deleted_at FROM transactions LIMIT ? OFFSET ?", filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	var transactions []*model.Transaction
	for res.Next() {
		var transaction model.Transaction
		if err := res.Scan(&transaction.Id, &transaction.UserId, &transaction.OrderId, &transaction.PaymentMethodId, &transaction.Status, &transaction.CreatedAt, &transaction.UpdatedAt, &transaction.DeletedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}

func (p *transactionRepository) FindById(ctx context.Context, id int64) (*model.Transaction, error) {
	transactionKey := getTransactionKey(id)

	var transaction model.Transaction
	tr, err := p.redis.Get(ctx, transactionKey).Result()
	if err == nil {
		err := json.Unmarshal([]byte(tr), &transaction)
		if err != nil {
			return nil, err
		}
		return &transaction, nil
	}

	err = p.db.QueryRowContext(ctx, "SELECT id, user_id, order_id, payment_method_id, status FROM transactions WHERE id=?", id).Scan(&transaction.Id, &transaction.UserId, &transaction.OrderId, &transaction.PaymentMethodId, &transaction.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	transactionJson, err := json.Marshal(&transaction)
	if err != nil {
		return nil, err
	}

	err = p.redis.Set(ctx, transactionKey, string(transactionJson), 0).Err()
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (p *transactionRepository) Create(ctx context.Context, transaction model.Transaction) error {
	result, err := p.db.ExecContext(ctx, "INSERT INTO transactions (user_id, order_id, payment_method_id, status) VALUES (?, ?, ?, ?)", transaction.UserId, transaction.OrderId, transaction.PaymentMethodId, transaction.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	_, err = p.db.ExecContext(ctx, "INSERT INTO transaction_logs (transaction_id, status) VALUES (?, ?)", id, transaction.Status)
	if err != nil {
		return err
	}
	return nil
}

func (p *transactionRepository) Update(ctx context.Context, transaction model.Transaction) error {
	_, err := p.db.ExecContext(ctx, "UPDATE transactions SET status=? WHERE id=?", transaction.Status, transaction.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *transactionRepository) Delete(ctx context.Context, id int64) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM transactions WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}

func getTransactionKey(id int64) string {
	return fmt.Sprintf("%s:%d", model.PaymentMethodKey, id)
}
