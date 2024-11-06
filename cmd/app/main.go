package main

import (
	"context"
	"github.com/zahidhasanpapon/go-clean-architecture/config"
	"github.com/zahidhasanpapon/go-clean-architecture/internal/infra/http"
	"github.com/zahidhasanpapon/go-clean-architecture/pkg/logger/zap"

	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Initialize configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	logCfg := &zap.Config{
		Level:      "debug",
		OutputPath: "stdout",
		Encoding:   "console",
		DevMode:    true,
	}

	// Initialize logger
	l, err := zap.NewZapLogger(logCfg)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}

	// Create new server instance
	srv, _ := http.NewServer(cfg, l)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			l.Fatal(context.Background(), "failed to start server: %v")
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	//l.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		//l.Fatalf("Server forced to shutdown: %v", err)
	}

	//l.Info("Server exited properly")
}
