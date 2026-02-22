package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mikenunan/openobserve-poc/services/partyman/handlers"
	"github.com/mikenunan/openobserve-poc/services/partyman/store"
	"github.com/mikenunan/openobserve-poc/services/partyman/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	// Structured JSON logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Initialise OpenTelemetry
	ctx := context.Background()
	shutdownTracer, err := telemetry.Setup(ctx, "partyman")
	if err != nil {
		slog.Error("failed to initialise telemetry", "error", err)
		os.Exit(1)
	}
	defer shutdownTracer(ctx)

	// Load seed data
	dataPath := "data/customers.json"
	if envPath := os.Getenv("SEED_DATA_PATH"); envPath != "" {
		dataPath = envPath
	}

	customerStore, err := store.New(dataPath)
	if err != nil {
		slog.Error("failed to load seed data", "path", dataPath, "error", err)
		os.Exit(1)
	}
	slog.Info("seed data loaded", "path", dataPath)

	// Create handler
	h := &handlers.Handler{
		Store:  customerStore,
		Logger: logger,
	}

	// Register routes with OTel HTTP instrumentation
	mux := http.NewServeMux()
	mux.HandleFunc("GET /customers", h.GetCustomers)
	mux.HandleFunc("GET /customer", h.GetCustomer)
	mux.HandleFunc("HEAD /customer", h.GetCustomer)
	mux.HandleFunc("POST /customer", h.CreateCustomer)
	mux.HandleFunc("PUT /customer", h.UpdateCustomer)
	mux.HandleFunc("GET /health", h.Health)

	// Wrap with OTel HTTP middleware
	otelHandler := otelhttp.NewHandler(mux, "partyman")

	server := &http.Server{
		Addr:    ":8081",
		Handler: otelHandler,
	}

	// Graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		slog.Info("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(shutdownCtx)
	}()

	slog.Info("PartyMan service starting", "port", 8081)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
	slog.Info("PartyMan service stopped")
}
