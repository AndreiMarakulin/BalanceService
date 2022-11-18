package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"balance-service/internal/model"
	"balance-service/internal/usecase"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BalanceRepo struct {
	pool *pgxpool.Pool
}

func NewBalanceRepo(pool *pgxpool.Pool) *BalanceRepo {
	return &BalanceRepo{pool: pool}
}

func (r *BalanceRepo) GetUserBalance(ctx context.Context, userID uint64) (model.Balance, error) {
	balance := model.Balance{}
	query := "SELECT user_id, balance FROM balance WHERE user_id = $1;"
	err := r.pool.QueryRow(ctx, query, userID).Scan(&balance.UserID, &balance.AmountMinor)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Balance{}, fmt.Errorf("%w: userID=%d", usecase.ErrUserNotExists, userID)
		}
		return model.Balance{}, fmt.Errorf("getUserBalance error: %w", err)
	}
	return balance, nil
}

func (r *BalanceRepo) ProcessIncome(ctx context.Context, transaction model.Transaction) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return fmt.Errorf("start trancaction error: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
		if err != nil && err != pgx.ErrTxClosed {
			log.Printf("rollback trancaction error: %s", err)
		}
	}(tx, ctx)

	balanceQuery := `INSERT INTO balance (user_id, balance) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET balance = balance.balance + excluded.balance;`
	transactionQuery := `INSERT INTO transaction (user_id, type, total) VALUES ($1, $2, $3);`

	_, err = tx.Exec(ctx, balanceQuery, transaction.UserID, transaction.AmountMinor)
	if err != nil {
		return fmt.Errorf("update balance error: %w", err)
	}
	_, err = tx.Exec(ctx, transactionQuery, transaction.UserID, transaction.Type, transaction.AmountMinor)
	if err != nil {
		return fmt.Errorf("insert transaction error: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit trancaction error: %w", err)
	}
	return nil
}
