package usecase

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/kodinggo/payment-service-gb1/internal/model"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

type transactionUsecase struct {
	transactionRepo model.ITransactionRepository
	js              jetstream.JetStream
}

var errNotFound = errors.New("transaction not found")

func NewTransactionUsecase(
	transactionRepo model.ITransactionRepository,
) model.ITransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
	}
}

func (t *transactionUsecase) FindAll(ctx context.Context, filter model.TransactionFilter) ([]*model.Transaction, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":    ctx,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})

	transactions, err := t.transactionRepo.FindAll(ctx, filter)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return transactions, nil
}

func (t *transactionUsecase) FindById(ctx context.Context, id int64) (*model.Transaction, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"id":  id,
	})

	transaction, err := t.transactionRepo.FindById(ctx, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if transaction == nil {
		return nil, errNotFound
	}

	return transaction, nil
}

func (t *transactionUsecase) Create(ctx context.Context, in model.CreateTransactionInput) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":               ctx,
		"user_id":           in.UserId,
		"order_id":          in.OrderId,
		"payment_method_id": in.PaymentMethodId,
		"status":            in.Status,
	})

	transaction := model.Transaction{
		UserId:          in.UserId,
		OrderId:         in.OrderId,
		PaymentMethodId: in.PaymentMethodId,
		Status:          in.Status,
	}

	err := t.transactionRepo.Create(ctx, transaction)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (t *transactionUsecase) Update(ctx context.Context, in model.UpdateTransactionInput) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":               ctx,
		"id":                in.Id,
		"user_id":           in.UserId,
		"order_id":          in.OrderId,
		"payment_method_id": in.PaymentMethodId,
		"status":            in.Status,
	})

	newTransaction := model.Transaction{
		Id:              in.Id,
		UserId:          in.UserId,
		OrderId:         in.OrderId,
		PaymentMethodId: in.PaymentMethodId,
		Status:          in.Status,
	}

	err := t.transactionRepo.Update(ctx, newTransaction)
	if err != nil {
		log.Error(err)
		return err
	}

	data, err := json.Marshal(newTransaction)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = t.js.Publish(ctx, "PAYMENT.updated", data)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (t *transactionUsecase) Delete(ctx context.Context, id int64) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"id":  id,
	})

	err := t.transactionRepo.Delete(ctx, id)
	if err != nil {
		log.Error(err)
		return err
	}

	data, err := json.Marshal(id)
	if err != nil {
		log.Error(err)
		return err
	}

	_, err = t.js.Publish(ctx, "PAYMENT.deleted", data)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
