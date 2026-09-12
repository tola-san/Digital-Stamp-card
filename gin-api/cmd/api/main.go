package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"gin-api/internal/config"
	"gin-api/internal/database"
	httpdelivery "gin-api/internal/delivery/http"
	"gin-api/internal/repository"
	"gin-api/internal/service"
)

// Entry point for running code gin start up
func main() {

	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("load configuration: %v", err)
	}

	//
	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancelStartup()

	db, err := database.Open(startupCtx, cfg.DatabaseURL, cfg.DatabaseCACertFile)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}

	defer db.Close()

	if err := database.ApplyMigrations(startupCtx, db); err != nil {
		log.Fatalf("apply database migrations: %v", err)
	}
	created, err := database.SeedStaff(startupCtx, db, database.StaffSeed{
		Name:     cfg.SeedStaffName,
		Email:    cfg.SeedStaffEmail,
		Password: cfg.SeedStaffPassword,
	})

	if err != nil {
		log.Fatalf("seed staff account: %v", err)
	}

	if created {
		log.Print("initial staff account created; remove SEED_STAFF_* environment variables")
	}

	customerRepository := repository.NewCustomerRepository(db.DB)
	customerSessionRepository := repository.NewCustomerSessionRepository(db.DB)
	staffRepository := repository.NewStaffRepository(db.DB)
	staffSessionRepository := repository.NewStaffSessionRepository(db.DB)
	transactionRepository := repository.NewTransactionRepository(db.DB)
	customerService := service.NewCustomerService(customerRepository, transactionRepository)
	customerSessionService := service.NewCustomerSessionService(customerRepository, customerSessionRepository)
	staffService := service.NewStaffSessionService(staffRepository, staffSessionRepository)

	router := httpdelivery.NewRouter(httpdelivery.RouterDependencies{
		HealthChecker:          db,
		FrontendURL:            cfg.FrontendURL,
		CustomerService:        customerService,
		CustomerSessionService: customerSessionService,
		StaffService:           staffService,
		CookieSecure:           cfg.CookieSecure,
	})
	server := &http.Server{Addr: ":" + cfg.Port, Handler: router, ReadHeaderTimeout: 5 * time.Second}
	serverErrors := make(chan error, 1)

	go func() {
		log.Printf("Digital Stamp API listening on %s", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	shutdownSignal, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	select {

	//
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("start server: %v", err)
		}
	case <-shutdownSignal.Done():
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
		}
	}
}
