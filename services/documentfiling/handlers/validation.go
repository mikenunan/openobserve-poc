package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// tracedClient is an HTTP client instrumented with OTel so outgoing requests
// carry the traceparent header and create child spans.
var tracedClient = &http.Client{
	Transport: otelhttp.NewTransport(http.DefaultTransport),
}

// ValidateCustomerExists calls HEAD /customer?partyId=<id> on PartyMan.
// Returns nil if the customer exists, or an error if not found or on failure.
func ValidateCustomerExists(ctx context.Context, partyID int, logger *slog.Logger) error {
	partymanURL := os.Getenv("PARTYMAN_URL")
	if partymanURL == "" {
		partymanURL = "http://localhost:8081"
	}

	reqURL := fmt.Sprintf("%s/customer?partyId=%d", partymanURL, partyID)
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, reqURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := tracedClient.Do(req)
	if err != nil {
		logger.ErrorContext(ctx, "failed to reach PartyMan", "url", reqURL, "error", err)
		return fmt.Errorf("PartyMan unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("customer with partyId %d does not exist", partyID)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("PartyMan returned unexpected status %d", resp.StatusCode)
	}

	logger.InfoContext(ctx, "customer validated", "partyId", partyID)
	return nil
}

// ValidateDocumentExists calls HEAD /document?name=<name> on DocumentCatalogue.
// Returns nil if the document exists, or an error if not found or on failure.
func ValidateDocumentExists(ctx context.Context, documentName string, logger *slog.Logger) error {
	catalogueURL := os.Getenv("DOCUMENTCATALOGUE_URL")
	if catalogueURL == "" {
		catalogueURL = "http://localhost:8082"
	}

	reqURL := fmt.Sprintf("%s/document?name=%s", catalogueURL, url.QueryEscape(documentName))
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, reqURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := tracedClient.Do(req)
	if err != nil {
		logger.ErrorContext(ctx, "failed to reach DocumentCatalogue", "url", reqURL, "error", err)
		return fmt.Errorf("DocumentCatalogue unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("document '%s' does not exist", documentName)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("DocumentCatalogue returned unexpected status %d", resp.StatusCode)
	}

	logger.InfoContext(ctx, "document validated", "documentName", documentName)
	return nil
}
