package model

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

const PaymentMethodKey = "payment"

var (
	ErrInvalidInput = errors.New("invalid input")
)

type IPaymentRepository interface {
	FindAll(ctx context.Context, filter PaymentFilter) ([]*Payment, error)
	FindById(ctx context.Context, id int64) (*Payment, error)
	Create(ctx context.Context, payment Payment) error
	Update(ctx context.Context, payment Payment) error
	Delete(ctx context.Context, id int64) error
}

type IPaymentUsecase interface {
	FindAll(ctx context.Context, filter PaymentFilter) ([]*Payment, error)
	FindById(ctx context.Context, id int64) (*Payment, error)
	Create(ctx context.Context, in CreatePaymentInput) error
	Update(ctx context.Context, in UpdatePaymentInput) error
	Delete(ctx context.Context, id int64) error
}

type Payment struct {
	Id        int64        `json:"id"`
	Name      string       `json:"name"`
	BankCode  string       `json:"bank_code"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type PaymentFilter struct {
	Offset int32
	Limit  int32
}

type CreatePaymentInput struct {
	Name     string `json:"name" validate:"required"`
	BankCode string `json:"bank_code" validate:"required"`
}

type UpdatePaymentInput struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	BankCode  string    `json:"bank_code"`
	UpdatedAt time.Time `json:"updated_at"`
}
