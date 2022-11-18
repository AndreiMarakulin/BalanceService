package usecase

import (
	"context"
	"errors"
	"fmt"

	"balance-service/internal/model"
)

var (
	ErrUserNotExists    = errors.New("user not found")
	ErrInputError       = errors.New("input error")
	ErrNotEnoughBalance = errors.New("not enough balance")
)

type BalanceRepo interface {
	GetUserBalance(ctx context.Context, userID uint64) (model.Balance, error)
	ProcessIncome(ctx context.Context, transaction model.Transaction) error
}

type BalanceUseCase struct {
	repo BalanceRepo
}

func NewBalanceUseCase(repo BalanceRepo) *BalanceUseCase {
	return &BalanceUseCase{repo: repo}
}

func (uc *BalanceUseCase) GetUserBalance(ctx context.Context, userID uint64) (model.Balance, error) {
	balance, err := uc.repo.GetUserBalance(ctx, userID)
	if err != nil {
		return model.Balance{}, err
	}

	balance.ConvertAmountToMajor()

	return balance, nil
}

func (uc *BalanceUseCase) ProcessIncome(ctx context.Context, transaction model.Transaction) error {
	if transaction.AmountMajor <= 0 {
		return fmt.Errorf("%w: income amount should be greaster then 0", ErrInputError)
	}
	transaction.ConvertAmountToMinor()
	transaction.Type = "transfer"
	return uc.repo.ProcessIncome(ctx, transaction)
}
