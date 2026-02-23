package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mikenunan/openobserve-poc/services/documentfiling/handlers"
	"github.com/mikenunan/openobserve-poc/services/documentfiling/models"
	"github.com/mikenunan/openobserve-poc/services/documentfiling/store"
	"github.com/mikenunan/openobserve-poc/services/documentfiling/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	// Structured JSON logging (bootstrap — will switch to OTel bridge below)
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	// Initialise OpenTelemetry (traces + logs + metrics)
	ctx := context.Background()
	logger, shutdownOTel, err := telemetry.Setup(ctx, "documentfiling")
	if err != nil {
		slog.Error("failed to initialise telemetry", "error", err)
		os.Exit(1)
	}
	defer shutdownOTel(ctx)

	// Switch to OTel-backed logger so all subsequent slog calls are exported
	slog.SetDefault(logger)

	// Pre-seeded demo affiliations
	seeds := []models.Affiliation{
		{PartyID: 12345, DocumentName: "A-Bank customer Ts&Cs"},
		{PartyID: 12345, DocumentName: "FSCS Info Sheet"},
		{PartyID: 12347, DocumentName: "Savings Account Ts&Cs"},
	}
	affiliationStore := store.New(seeds)
	slog.Info("pre-seeded affiliations loaded", "count", len(seeds))

	// Create handler
	h := &handlers.Handler{
		Store:  affiliationStore,
		Logger: logger,
	}

	// Register routes with OTel HTTP instrumentation
	mux := http.NewServeMux()
	mux.HandleFunc("POST /affiliation", h.CreateAffiliation)
	mux.HandleFunc("GET /documents", h.GetDocuments)
	mux.HandleFunc("GET /health", h.Health)

	// Wrap with OTel HTTP middleware
	otelHandler := otelhttp.NewHandler(mux, "documentfiling")

	server := &http.Server{
		Addr:    ":8083",
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

	slog.Info("DocumentFiling service starting", "port", 8083)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
	slog.Info("DocumentFiling service stopped")
}
