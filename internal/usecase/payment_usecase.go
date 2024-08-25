package usecase

import (
	"context"

	"github.com/tubagusmf/payment-service-gb1/internal/model"

	"github.com/go-playground/validator"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

type PaymentUsecase struct {
	paymentRepo model.IPaymentRepository
	// workerClient *worker.AsynqClient
}

var v = validator.New()

func NewPaymentUsecase(
	paymentRepo model.IPaymentRepository,
	// workerClient *worker.AsynqClient,
) model.IPaymentUsecase {
	return &PaymentUsecase{
		paymentRepo: paymentRepo,
		// workerClient: workerClient,
	}
}

func (p *PaymentUsecase) FindAll(ctx context.Context, filter model.PaymentFilter) ([]*model.Payment, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx":    ctx,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})

	payments, err := p.paymentRepo.FindAll(ctx, filter)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return payments, nil
}

func (p *PaymentUsecase) FindById(ctx context.Context, id int64) (*model.Payment, error) {
	log := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"id":  id,
	})

	payment, err := p.paymentRepo.FindById(ctx, id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return payment, nil
}

func (p *PaymentUsecase) Create(ctx context.Context, in model.CreatePaymentInput) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":       ctx,
		"name":      in.Name,
		"bank_code": in.BankCode,
	})

	err := p.validateCreatePaymentInput(ctx, in)
	if err != nil {
		log.Error(err)
		return err
	}

	payment := model.Payment{
		Name:     in.Name,
		BankCode: in.BankCode,
	}

	err = p.paymentRepo.Create(ctx, payment)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (p *PaymentUsecase) Update(ctx context.Context, in model.UpdatePaymentInput) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx":       ctx,
		"id":        in.Id,
		"name":      in.Name,
		"bank_code": in.BankCode,
	})

	err := p.validateUpdatePaymentInput(ctx, in)
	if err != nil {
		log.Error(err)
		return err
	}

	newPayment := model.Payment{
		Id:        in.Id,
		Name:      in.Name,
		BankCode:  in.BankCode,
		UpdatedAt: in.UpdatedAt,
	}

	err = p.paymentRepo.Update(ctx, newPayment)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (p *PaymentUsecase) Delete(ctx context.Context, id int64) error {
	log := logrus.WithFields(logrus.Fields{
		"ctx": ctx,
		"id":  id,
	})

	err := p.paymentRepo.Delete(ctx, id)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (p *PaymentUsecase) validateCreatePaymentInput(ctx context.Context, in model.CreatePaymentInput) error {
	err := v.StructCtx(ctx, in)
	if err != nil {
		log.Error(err)
		return model.ErrInvalidInput
	}
	return nil
}

func (p *PaymentUsecase) validateUpdatePaymentInput(ctx context.Context, in model.UpdatePaymentInput) error {
	err := v.StructCtx(ctx, in)
	if err != nil {
		log.Error(err)
		return model.ErrInvalidInput
	}
	return nil
}
