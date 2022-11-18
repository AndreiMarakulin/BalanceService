package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"avito-internship/config"
	"avito-internship/internal/controller"
	"avito-internship/internal/repository"
	"avito-internship/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.New()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.PG.URL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		log.Fatalf("Databse connection error: %v\n", err)
	}

	balanceUseCase := usecase.NewBalanceUseCase(repository.NewBalanceRepo(pool))
	orderUseCase := usecase.NewOrderUseCase(repository.NewOrderRepo(pool))
	router := chi.NewRouter()
	controller.NewBalanceRoutes(router, balanceUseCase)
	controller.NewOrderRoutes(router, orderUseCase)

	httpServer := http.Server{
		Addr:              net.JoinHostPort("", cfg.HTTP.Port),
		Handler:           router,
		ReadHeaderTimeout: 2 * time.Second,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalf("HTTP ListenAndServe error: %v\n", err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	_ = httpServer.Shutdown(ctx)

	pool.Close()
}
