package main

import (
	"context"
	"database/sql"
	domainGauge "energosbyt/service/internal/domains/gauge"
	domainInvoices "energosbyt/service/internal/domains/invoices"
	domainReadings "energosbyt/service/internal/domains/readings"
	domainRealty "energosbyt/service/internal/domains/realty"
	domainTariffs "energosbyt/service/internal/domains/tariffs"
	domainUsers "energosbyt/service/internal/domains/users"
	"energosbyt/service/internal/handlers/bind_gauge"
	"energosbyt/service/internal/handlers/get_invoices"
	"energosbyt/service/internal/handlers/get_readings"
	"energosbyt/service/internal/handlers/ping"
	"energosbyt/service/internal/handlers/post_reading"
	"energosbyt/service/internal/handlers/realty"
	"energosbyt/service/internal/handlers/tariffs"
	"energosbyt/service/internal/handlers/users"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/ClickHouse/clickhouse-go/v2"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("application stopped with error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	// 1. Чтение конфигурации
	chHost := getEnv("CLICKHOUSE_HOST", "localhost")
	chPort := getEnv("CLICKHOUSE_PORT", "9000")
	chDB := getEnv("CLICKHOUSE_DATABASE", "energosbyt")
	chUser := getEnv("CLICKHOUSE_USER", "default")
	chPassword := os.Getenv("CLICKHOUSE_PASSWORD")
	appPort := getEnv("APP_PORT", "8080")

	// 2. Подключение к ClickHouse с проверкой по таймауту
	dsn := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s?dial_timeout=10s&read_timeout=20s",
		chUser, chPassword, chHost, chPort, chDB)

	db, err := sql.Open("clickhouse", dsn)
	if err != nil {
		return fmt.Errorf("open clickhouse: %w", err)
	}
	defer db.Close()

	pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		return fmt.Errorf("ping clickhouse: %w", err)
	}

	realtyDomain := domainRealty.NewDomain(db)
	usersDomain := domainUsers.NewDomain(db)
	tariffsDomain := domainTariffs.NewDomain(db)
	gaugeDomain := domainGauge.NewDomain(db, realtyDomain)
	invoicesDomain := domainInvoices.NewDomain(db, realtyDomain)
	readingsDomain := domainReadings.NewDomain(db, usersDomain, gaugeDomain, tariffsDomain, realtyDomain)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", ping.New(logger).ServeHTTP)

	mux.HandleFunc("GET /users/{user_id}", users.New(usersDomain, logger).ServeHTTP)
	mux.HandleFunc("PATCH /users/{user_id}", users.New(usersDomain, logger).ServeHTTP)

	mux.HandleFunc("GET /realty/{realty_id}", realty.New(realtyDomain, logger).ServeHTTP)
	mux.HandleFunc("POST /realty/{realty_id}/readings", post_reading.New(readingsDomain, logger).ServeHTTP)
	mux.HandleFunc("GET /realty/{realty_id}/readings", get_readings.New(readingsDomain, logger).ServeHTTP)
	mux.HandleFunc("GET /realty/{realty_id}/invoices", get_invoices.New(invoicesDomain, logger).ServeHTTP)
	mux.HandleFunc("POST /realty/{realty_id}/gauges", bind_gauge.New(gaugeDomain, logger).ServeHTTP)

	// Тарифы
	mux.HandleFunc("GET /tariffs", tariffs.New(tariffsDomain, logger).ServeHTTP)

	// 5. Настройка HTTP-сервера с защитными таймаутами
	httpServer := &http.Server{
		Addr:         ":" + appPort,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 6. Запуск сервера и плавное завершение работы (Graceful Shutdown)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("starting http server", "port", appPort)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return fmt.Errorf("http server failed: %w", err)
	case <-ctx.Done():
		logger.Info("shutting down server gracefully...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server forced to shutdown: %w", err)
		}
		logger.Info("server stopped")
		return nil
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
