package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mikenunan/openobserve-poc/services/documentcatalogue/handlers"
	"github.com/mikenunan/openobserve-poc/services/documentcatalogue/telemetry"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {
	// Structured JSON logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// Initialise OpenTelemetry
	ctx := context.Background()
	shutdownTracer, err := telemetry.Setup(ctx, "documentcatalogue")
	if err != nil {
		slog.Error("failed to initialise telemetry", "error", err)
		os.Exit(1)
	}
	defer shutdownTracer(ctx)

	// Create handler
	dataDir := "data"
	if envDir := os.Getenv("DATA_DIR"); envDir != "" {
		dataDir = envDir
	}

	h := &handlers.Handler{
		DataDir: dataDir,
		Logger:  logger,
	}

	// Register routes with OTel HTTP instrumentation
	mux := http.NewServeMux()
	mux.HandleFunc("GET /document", h.GetDocument)
	mux.HandleFunc("HEAD /document", h.GetDocument)
	mux.HandleFunc("GET /documents", h.GetDocumentNames)
	mux.HandleFunc("GET /health", h.Health)

	// Wrap with OTel HTTP middleware
	otelHandler := otelhttp.NewHandler(mux, "documentcatalogue")

	server := &http.Server{
		Addr:    ":8082",
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

	slog.Info("DocumentCatalogue service starting", "port", 8082)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
	slog.Info("DocumentCatalogue service stopped")
}
