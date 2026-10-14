package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/taisuck/mantong/src/api"
	"github.com/taisuck/mantong/src/config"
	"github.com/taisuck/mantong/src/logging"
	"github.com/taisuck/mantong/src/mongodb"
)

func main() {
	configPath := flag.String("config", "config.yml", "설정 파일 경로")
	flag.Parse()

	if err := run(*configPath); err != nil {
		slog.Error("서버 종료", "err", err)
		os.Exit(1)
	}
}

func run(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	logCloser, err := logging.Setup(cfg.Log.Dir, cfg.Log.Level)
	if err != nil {
		return fmt.Errorf("로그 설정 실패: %w", err)
	}
	defer logCloser.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	client, db, err := mongodb.Connect(ctx, cfg.MongoDB)
	if err != nil {
		return err
	}
	defer func() {
		if err := mongodb.Disconnect(client, cfg.MongoDB.Timeout); err != nil {
			slog.Error("mongodb 연결 종료 실패", "err", err)
		}
	}()
	slog.Info("mongodb 연결 완료", "database", cfg.MongoDB.Database)

	srv := &http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           api.NewRouter(mongodb.NewStockRepository(db)),
		ReadHeaderTimeout: 10 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("웹 서버 시작", "addr", cfg.Server.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("웹 서버 실행 실패: %w", err)
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	slog.Info("웹 서버 종료 중")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("웹 서버 종료 실패: %w", err)
	}
	return nil
}
