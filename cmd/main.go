package main

import (
	"context"
	"github.com/Tairascii/google-docs-user/internal/app"
	"github.com/Tairascii/google-docs-user/internal/app/handler"
	"github.com/Tairascii/google-docs-user/internal/app/service/user"
	"github.com/Tairascii/google-docs-user/internal/app/service/user/repo"
	"github.com/Tairascii/google-docs-user/internal/app/usecase"
	"github.com/Tairascii/google-docs-user/internal/db"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//TODO add envs
	dbSettings := db.Settings{
		Host:          "localhost",
		Port:          "5432",
		User:          "admin",
		Password:      "12345",
		DbName:        "google_doc_users",
		Schema:        "google_doc_users_schema",
		AppName:       "google_doc_users",
		MaxIdleConns:  2,
		MaxOpenConns:  5,
		MigrateSchema: true,
	}

	sqlxDb, err := db.Connect(dbSettings)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
		return
	}
	defer func(sqlxDb *sqlx.DB) {
		if err := sqlxDb.Close(); err != nil {
			log.Fatalf("failed to close connection to db: %v", err)
		}
	}(sqlxDb)

	userRepo := repo.New(sqlxDb)
	userSrv := user.New(userRepo)

	authUC := usecase.NewAuthUseCase(userSrv)

	useCase := app.UseCase{Auth: authUC}
	DI := &app.DI{UseCase: useCase}
	handlers := handler.NewHandler(DI)

	srv := &http.Server{
		Addr:         ":8000", // TODO add .env
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
		Handler:      handlers.InitHandlers(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("Something went wrong while runing server %s", err.Error())
		}
	}()

	log.Println("Listening on port 8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-quit

	log.Println("Shutting down server")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Something went wrong while shutting down server %s", err.Error())
	}
}
