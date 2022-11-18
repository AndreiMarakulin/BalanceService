package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"balance-service/internal/model"
	"balance-service/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type balanceRoutes struct {
	uc BalanceUseCase
}

type BalanceUseCase interface {
	GetUserBalance(ctx context.Context, userID uint64) (model.Balance, error)
	ProcessIncome(ctx context.Context, transaction model.Transaction) error
}

func NewBalanceRoutes(router chi.Router, uc BalanceUseCase) {
	br := &balanceRoutes{uc: uc}

	router.Route("/balance", func(r chi.Router) {
		r.Post("/income", br.processIncome)
		r.Get("/{userID}", br.getUserBalance)
	})
}

func (br *balanceRoutes) getUserBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID, err := strconv.ParseUint(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse user id '%s'", chi.URLParam(r, "userID")), http.StatusBadRequest)

		return
	}

	var balance model.Balance
	if balance, err = br.uc.GetUserBalance(r.Context(), userID); err != nil {
		if errors.Unwrap(err) == usecase.ErrUserNotExists {
			http.Error(w, err.Error(), http.StatusNotFound)
			log.Println(err)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if err = json.NewEncoder(w).Encode(balance); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (br *balanceRoutes) processIncome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	transaction := model.Transaction{}
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to parse body data '%s'", r.Body), http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = br.uc.ProcessIncome(r.Context(), transaction)
	if err != nil {
		if errors.Unwrap(err) == usecase.ErrInputError {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err)
			return
		}
		http.Error(w, "server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if err = json.NewEncoder(w).Encode(map[string]string{"status": "success"}); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
