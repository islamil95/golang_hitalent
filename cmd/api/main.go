package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/islamil95/golang_hitalent/internal/config"
	"github.com/islamil95/golang_hitalent/internal/handler"
	"github.com/islamil95/golang_hitalent/internal/middleware"
	"github.com/islamil95/golang_hitalent/internal/repository"
	"github.com/islamil95/golang_hitalent/internal/service"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	cfg, err := config.Load()
	if err != nil {
		log.Error("load config", "err", err)
		os.Exit(1)
	}

	db, err := waitForDB(cfg.DSN, log)
	if err != nil {
		log.Error("db connection", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := runMigrations(db, log); err != nil {
		log.Error("migrations", "err", err)
		os.Exit(1)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // отключаем вывод GORM (record not found, SQL) — только наши логи
	})
	if err != nil {
		log.Error("gorm open", "err", err)
		os.Exit(1)
	}

	depRepo := repository.NewDepartmentRepository(gormDB)
	empRepo := repository.NewEmployeeRepository(gormDB)
	depSvc := service.NewDepartmentService(depRepo, empRepo)
	empSvc := service.NewEmployeeService(depRepo, empRepo)

	depH := handler.NewDepartmentHandler(depSvc, log)
	empH := handler.NewEmployeeHandler(empSvc, log)
	router := handler.NewRouter(depH, empH)

	h := middleware.Recovery(log, middleware.Logging(log, router))

	srv := &http.Server{
		Addr:    cfg.Addr(),
		Handler: h,
	}
	go func() {
		log.Info("server listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen", "err", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("shutdown", "err", err)
	}
}

const (
	dbRetryAttempts = 30
	dbRetryInterval = 2 * time.Second
)

// waitForDB подключается к БД с повторами, чтобы дождаться полного старта контейнера PostgreSQL.
func waitForDB(dsn string, log *slog.Logger) (*sql.DB, error) {
	var db *sql.DB
	var err error
	for i := 0; i < dbRetryAttempts; i++ {
		db, err = sql.Open("pgx", dsn)
		if err != nil {
			log.Warn("db open attempt failed", "attempt", i+1, "err", err)
			time.Sleep(dbRetryInterval)
			continue
		}
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(5)
		if err = db.Ping(); err != nil {
			_ = db.Close()
			log.Warn("db ping failed, waiting for container", "attempt", i+1, "err", err)
			time.Sleep(dbRetryInterval)
			continue
		}
		log.Info("database connected")
		return db, nil
	}
	return nil, err
}

func runMigrations(db *sql.DB, log *slog.Logger) error {
	goose.SetBaseFS(nil)
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	log.Info("running migrations")
	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}
	log.Info("migrations completed")
	return nil
}
