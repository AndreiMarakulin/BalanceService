package controller

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"avito-internship/internal/model"
	"avito-internship/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type orderRoutes struct {
	uc OrderUseCase
}

type OrderUseCase interface {
	CreateOrder(ctx context.Context, transaction model.Transaction) error
	FinishOrder(ctx context.Context, transaction model.Transaction) error
	CancelOrder(ctx context.Context, transaction model.Transaction) error
}

func NewOrderRoutes(router chi.Router, uc OrderUseCase) {
	br := &orderRoutes{uc: uc}

	router.Route("/order", func(r chi.Router) {
		r.Post("/create", br.createOrder)
		r.Post("/finish", br.finishOrder)
		r.Post("/cancel", br.cancelOrder)
	})
}

func (br *orderRoutes) createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	transaction := model.Transaction{}
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = br.uc.CreateOrder(r.Context(), transaction)
	result := map[string]string{"status": "success"}
	if err != nil {
		switch errors.Unwrap(err) {
		case usecase.ErrUserNotExists, usecase.ErrInputError:
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err)
			return
		case usecase.ErrNotEnoughBalance:
			result["status"] = "Not enough balance"
		default:
			http.Error(w, "server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (br *orderRoutes) finishOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	transaction := model.Transaction{}
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = br.uc.FinishOrder(r.Context(), transaction)
	result := map[string]string{"status": "success"}
	if err != nil {
		switch errors.Unwrap(err) {
		case usecase.ErrUserNotExists, usecase.ErrInputError:
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err)
			return
		case usecase.ErrNotEnoughBalance:
			result["status"] = "Not enough balance"
		default:
			http.Error(w, "server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func (br *orderRoutes) cancelOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	transaction := model.Transaction{}
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = br.uc.CancelOrder(r.Context(), transaction)
	result := map[string]string{"status": "success"}
	if err != nil {
		switch errors.Unwrap(err) {
		case usecase.ErrUserNotExists, usecase.ErrInputError:
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Println(err)
			return
		case usecase.ErrNotEnoughBalance:
			result["status"] = "No order to cancel"
		default:
			http.Error(w, "server error", http.StatusInternalServerError)
			log.Println(err)
			return
		}
	}
	if err = json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
