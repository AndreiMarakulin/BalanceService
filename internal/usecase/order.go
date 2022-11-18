package usecase

import (
	"context"
	"fmt"

	"balance-service/internal/model"
)

type OrderRepo interface {
	CreateOrder(ctx context.Context, transaction model.Transaction) error
	FinishOrder(ctx context.Context, transaction model.Transaction) error
	CancelOrder(ctx context.Context, transaction model.Transaction) error
}

type OrderUseCase struct {
	repo OrderRepo
}

func NewOrderUseCase(repo OrderRepo) *OrderUseCase {
	return &OrderUseCase{repo: repo}
}

func (uc *OrderUseCase) CreateOrder(ctx context.Context, transaction model.Transaction) error {
	if transaction.AmountMajor <= 0 {
		return fmt.Errorf("%w: total amount should be greaster then 0", ErrInputError)
	}
	transaction.ConvertAmountToMinor()
	transaction.Type = "capture"
	return uc.repo.CreateOrder(ctx, transaction)
}

func (uc *OrderUseCase) FinishOrder(ctx context.Context, transaction model.Transaction) error {
	if transaction.AmountMajor <= 0 {
		return fmt.Errorf("%w: total amount should be greaster then 0", ErrInputError)
	}
	transaction.ConvertAmountToMinor()
	transaction.Type = "write-off"
	return uc.repo.FinishOrder(ctx, transaction)
}

func (uc *OrderUseCase) CancelOrder(ctx context.Context, transaction model.Transaction) error {
	if transaction.AmountMajor <= 0 {
		return fmt.Errorf("%w: total amount should be greaster then 0", ErrInputError)
	}
	transaction.ConvertAmountToMinor()
	transaction.Type = "cancel"
	return uc.repo.CancelOrder(ctx, transaction)
}
