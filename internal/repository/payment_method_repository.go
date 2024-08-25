package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/tubagusmf/payment-service-gb1/internal/model"

	"github.com/redis/go-redis/v9"
)

type PaymentRepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewPaymentRepository(db *sql.DB, redis *redis.Client) model.IPaymentRepository {
	return &PaymentRepository{
		db:    db,
		redis: redis,
	}
}

func (p *PaymentRepository) FindAll(ctx context.Context, filter model.PaymentFilter) ([]*model.Payment, error) {

	res, err := p.db.QueryContext(ctx, "SELECT id, name, bank_code, created_at, updated_at, deleted_at FROM payment_methods LIMIT ? OFFSET ?", filter.Limit, filter.Offset)

	if err != nil {
		return nil, err
	}

	var payments []*model.Payment
	for res.Next() {
		var payment model.Payment
		if err := res.Scan(&payment.Id, &payment.Name, &payment.BankCode, &payment.CreatedAt, &payment.UpdatedAt, &payment.DeletedAt); err != nil {
			return nil, err
		}
		payments = append(payments, &payment)
	}

	return payments, nil
}

func (p *PaymentRepository) FindById(ctx context.Context, id int64) (*model.Payment, error) {
	paymentKey := getPaymentKey(id)

	var payment model.Payment
	pm, err := p.redis.Get(ctx, paymentKey).Result()
	if err == nil {
		err := json.Unmarshal([]byte(pm), &payment)
		if err != nil {
			return nil, err
		}
		return &payment, nil
	}

	res, err := p.db.QueryContext(ctx, "SELECT id, name, bank_code, created_at, updated_at, deleted_at FROM payment_methods WHERE id=?", id)

	if err != nil {
		return nil, err
	}

	for res.Next() {
		if err := res.Scan(&payment.Id, &payment.Name, &payment.BankCode, &payment.CreatedAt, &payment.UpdatedAt, &payment.DeletedAt); err != nil {
			return nil, err
		}
	}

	paymentJson, err := json.Marshal(payment)
	if err != nil {
		return nil, err
	}

	err = p.redis.Set(ctx, paymentKey, string(paymentJson), 0).Err()
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (p *PaymentRepository) Create(ctx context.Context, payment model.Payment) error {
	_, err := p.db.ExecContext(ctx, "INSERT INTO payment_methods (name, bank_code) VALUES (?, ?)", payment.Name, payment.BankCode)
	if err != nil {
		return err
	}

	return nil
}

func (p *PaymentRepository) Update(ctx context.Context, payment model.Payment) error {
	_, err := p.db.ExecContext(ctx, "UPDATE payment_methods SET name=?, bank_code=? WHERE id=?", payment.Name, payment.BankCode, payment.Id)
	if err != nil {
		return err
	}
	return nil
}

func (p *PaymentRepository) Delete(ctx context.Context, id int64) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM payment_methods WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}

// For Key Redis
func getPaymentKey(id int64) string {
	return fmt.Sprintf("%s:%d", model.PaymentMethodKey, id)
}
