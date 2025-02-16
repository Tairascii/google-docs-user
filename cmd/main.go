package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Tairascii/google-docs-user/internal/app"
	grpcServer "github.com/Tairascii/google-docs-user/internal/app/grpc"
	"github.com/Tairascii/google-docs-user/internal/app/handler"
	"github.com/Tairascii/google-docs-user/internal/app/service/user"
	"github.com/Tairascii/google-docs-user/internal/app/service/user/repo"
	"github.com/Tairascii/google-docs-user/internal/app/usecase"
	"github.com/Tairascii/google-docs-user/internal/db"
	"github.com/jmoiron/sqlx"
)

// TODO move to app
func main() {
	cfg, err := app.LoadConfigs()
	if err != nil {
		panic(err)
	}

	dbSettings := db.Settings{
		Host:          cfg.Repo.Host,
		Port:          cfg.Repo.Port,
		User:          cfg.Repo.User,
		Password:      cfg.Repo.Password,
		DbName:        cfg.Repo.DBName,
		Schema:        cfg.Repo.Schema,
		AppName:       cfg.Repo.AppName,
		MaxIdleConns:  cfg.Repo.MaxIdleConns,
		MaxOpenConns:  cfg.Repo.MaxOpenConns,
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
	userUC := usecase.NewUserUseCase(userSrv)

	useCase := app.UseCase{Auth: authUC, User: userUC}
	DI := &app.DI{UseCase: useCase}
	handlers := handler.NewHandler(DI)
	grpcSrv, err := grpcServer.NewGrpcServer(cfg.GrpcServer.Port, DI)
	if err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		ReadTimeout:  cfg.Server.Timeout.Read,
		WriteTimeout: cfg.Server.Timeout.Write,
		IdleTimeout:  cfg.Server.Timeout.Idle,
		Handler:      handlers.InitHandlers(),
	}

	go func() {
		if err := grpcSrv.Start(); err != nil {
			log.Fatalf("failed to start grpc server: %v", err)
		}
	}()

	go func() {
		if err := srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("something went wrong while runing server %s", err.Error())
		}
	}()

	log.Printf("listening on port: %s", cfg.Server.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	<-quit

	log.Println("shutting down server")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("something went wrong while shutting down server %s", err.Error())
	}
}
