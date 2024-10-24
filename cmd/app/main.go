package main

import (
	"context"
	"errors"
	"github.com/varonikp/opencontact-task/internal/repository/mysqlrepo"
	"github.com/varonikp/opencontact-task/internal/services"
	"github.com/varonikp/opencontact-task/internal/transport/httpserver"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := openDatabase(os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("init database: %v", err)
	}

	err = runMigrations(os.Getenv("DSN"), os.Getenv("MIGRATIONS_PATH"))
	if err != nil {
		log.Fatalf("up migrations: %v", err)
	}

	ratesRepo := mysqlrepo.NewCurrencyRepository(db)

	httpClient := http.DefaultClient

	ratesService := services.NewRatesService(httpClient, ratesRepo)
	if err := ratesService.Start(ctx); err != nil {
		log.Fatalf("failed start rates service: %v", err)
	}

	ratesHandler := httpserver.NewRatesHandler(ctx, ratesRepo)

	srv := httpserver.New(ratesHandler, os.Getenv("HTTP_ADDR"))

	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(stopped)
	}()

	log.Printf("Starting HTTP server on %s", os.Getenv("HTTP_ADDR"))

	// start HTTP server
	if err := srv.Start(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("http server error: %v", err)
	}

	<-stopped

	log.Printf("Goodbye.")

}
