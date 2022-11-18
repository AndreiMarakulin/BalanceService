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

type OrderRepo struct {
	pool *pgxpool.Pool
}

func NewOrderRepo(pool *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{pool: pool}
}

func (r *OrderRepo) CreateOrder(ctx context.Context, transaction model.Transaction) error {
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

	checkBalanceQuery := `SELECT balance FROM balance WHERE user_id = $1;`
	var currentBalance int
	if err = r.pool.QueryRow(ctx, checkBalanceQuery, transaction.UserID).Scan(&currentBalance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: userID=%d", usecase.ErrUserNotExists, transaction.UserID)
		}
		return fmt.Errorf("check user balance error: %w", err)
	}
	if currentBalance < transaction.AmountMinor {
		return fmt.Errorf("%w: currentBalance=%.2f, needed=$%.2f",
			usecase.ErrNotEnoughBalance, float64(currentBalance)/model.MinorUnitsInMajor, transaction.AmountMajor)
	}

	balanceQuery := `UPDATE balance SET balance = balance.balance - $2 WHERE user_id=$1;`
	holdBalanceQuery := `UPDATE balance SET hold = balance.hold + $2 WHERE user_id=$1;`
	transactionQuery := `INSERT INTO transaction (user_id, type, service_id, order_id, total) 
						VALUES ($1, $2, $3, $4, $5);`

	_, err = tx.Exec(ctx, balanceQuery, transaction.UserID, transaction.AmountMinor)
	if err != nil {
		return fmt.Errorf("update balance error: %w", err)
	}
	_, err = tx.Exec(ctx, holdBalanceQuery, transaction.UserID, transaction.AmountMinor)
	if err != nil {
		return fmt.Errorf("update hold balance error: %w", err)
	}
	_, err = tx.Exec(ctx, transactionQuery, transaction.UserID, transaction.Type,
		transaction.ServiceID, transaction.OrderID, -transaction.AmountMinor)
	if err != nil {
		return fmt.Errorf("insert transaction error: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit trancaction error: %w", err)
	}
	return nil
}

func (r *OrderRepo) FinishOrder(ctx context.Context, transaction model.Transaction) error {
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

	checkHoldBalanceQuery := `SELECT hold FROM balance WHERE user_id = $1;`
	var currentHoldBalance int
	if err = r.pool.QueryRow(ctx, checkHoldBalanceQuery, transaction.UserID).Scan(&currentHoldBalance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: userID=%d", usecase.ErrUserNotExists, transaction.UserID)
		}
		return fmt.Errorf("check user hold balance error: %w", err)
	}
	if currentHoldBalance < transaction.AmountMinor {
		return fmt.Errorf("%w: currentHoldBalance=%.2f, needed=$%.2f",
			usecase.ErrNotEnoughBalance, float64(currentHoldBalance)/model.MinorUnitsInMajor, transaction.AmountMajor)
	}

	balanceQuery := `UPDATE balance SET hold = balance.hold - $2 WHERE user_id=$1;`
	transactionQuery := `INSERT INTO transaction (user_id, type, service_id, order_id, total) 
						VALUES ($1, $2, $3, $4, $5);`

	_, err = tx.Exec(ctx, balanceQuery, transaction.UserID, transaction.AmountMinor)
	if err != nil {
		return fmt.Errorf("update hold balance error: %w", err)
	}
	_, err = tx.Exec(ctx, transactionQuery, transaction.UserID, transaction.Type,
		transaction.ServiceID, transaction.OrderID, -transaction.AmountMinor)
	if err != nil {
		return fmt.Errorf("insert transaction error: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit trancaction error: %w", err)
	}
	return nil
}

func (r *OrderRepo) CancelOrder(ctx context.Context, transaction model.Transaction) error {
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

	checkHoldBalanceQuery := `SELECT hold FROM balance WHERE user_id = $1;`
	var currentHoldBalance int
	if err = r.pool.QueryRow(ctx, checkHoldBalanceQuery, transaction.UserID).Scan(&currentHoldBalance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%w: userID=%d", usecase.ErrUserNotExists, transaction.UserID)
		}
		return fmt.Errorf("check user hold balance error: %w", err)
	}
	if currentHoldBalance < transaction.AmountMinor {
		return fmt.Errorf("%w: currentHoldBalance=%.2f, needed=$%.2f",
			usecase.ErrNotEnoughBalance, float64(currentHoldBalance)/model.MinorUnitsInMajor, transaction.AmountMajor)
	}

	balanceQuery := `UPDATE balance SET balance = balance.balance + $2 WHERE user_id=$1;`
	holdBalanceQuery := `UPDATE balance SET hold = balance.hold - $2 WHERE user_id=$1;`
	transactionQuery := `INSERT INTO transaction (user_id, type, service_id, order_id, total) 
						VALUES ($1, $2, $3, $4, $5);`

	_, err = tx.Exec(ctx, balanceQuery, transaction.UserID, transaction.AmountMinor)
	if err != nil {
		return fmt.Errorf("update balance error: %w", err)
	}
	_, err = tx.Exec(ctx, holdBalanceQuery, transaction.UserID, transaction.AmountMinor)
	if err != nil {
		return fmt.Errorf("update hold balance error: %w", err)
	}
	_, err = tx.Exec(ctx, transactionQuery, transaction.UserID, transaction.Type,
		transaction.ServiceID, transaction.OrderID, -transaction.AmountMinor)
	if err != nil {
		return fmt.Errorf("insert transaction error: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit trancaction error: %w", err)
	}
	return nil
}
