package model

import (
	"context"
	"database/sql"
	"time"
)

type ITransactionRepository interface {
	Create(ctx context.Context, transaction Transaction) error
	FindAll(ctx context.Context, filter TransactionFilter) ([]*Transaction, error)
	FindById(ctx context.Context, id int64) (*Transaction, error)
	Update(ctx context.Context, transaction Transaction) error
	Delete(ctx context.Context, id int64) error
}

type ITransactionUsecase interface {
	Create(ctx context.Context, in CreateTransactionInput) error
	FindAll(ctx context.Context, filter TransactionFilter) ([]*Transaction, error)
	FindById(ctx context.Context, id int64) (*Transaction, error)
	Update(ctx context.Context, in UpdateTransactionInput) error
	Delete(ctx context.Context, id int64) error
}

type Transaction struct {
	Id              int64        `json:"id"`
	UserId          int64        `json:"user_id"`
	OrderId         int64        `json:"order_id"`
	PaymentMethodId int64        `json:"payment_method_id"`
	Status          string       `json:"status"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	DeletedAt       sql.NullTime `json:"deleted_at"`
}

type TransactionLog struct {
	Id            int64  `json:"id"`
	TransactionId int64  `json:"transaction_id"`
	Status        string `json:"status"`
}

type TransactionFilter struct {
	Order_id int64
	Offset   int32
	Limit    int32
}

type CreateTransactionInput struct {
	UserId          int64
	OrderId         int64
	PaymentMethodId int64
	Status          string
}

type CreateTransactionLogInput struct {
	TransactionId int64
	Status        string
}

type UpdateTransactionInput struct {
	Id              int64
	UserId          int64
	OrderId         int64
	PaymentMethodId int64
	Status          string
	UpdatedAt       time.Time
}
