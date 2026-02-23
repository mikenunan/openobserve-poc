package handlers

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

// documentNameToFile maps document display names to file paths.
var documentNameToFile = map[string]string{
	"A-Bank customer Ts&Cs": "customer-tsandcs.txt",
	"Savings Account Ts&Cs": "savings-tsandcs.txt",
	"FSCS Info Sheet":       "fscs-info.txt",
}

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	DataDir string
	Logger  *slog.Logger
}

// GetDocument handles GET /document?name=<docName> — returns the document text.
// Also handles HEAD /document?name=<docName> — existence check (no body).
func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	docName := r.URL.Query().Get("name")
	if docName == "" {
		http.Error(w, `{"error":"name query parameter is required"}`, http.StatusBadRequest)
		return
	}

	fileName, ok := documentNameToFile[docName]
	if !ok {
		h.Logger.InfoContext(r.Context(), "document not found", "name", docName)
		http.Error(w, `{"error":"document not found"}`, http.StatusNotFound)
		return
	}

	// HEAD requests get just the status code, no body
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	filePath := filepath.Join(h.DataDir, fileName)
	content, err := os.ReadFile(filePath)
	if err != nil {
		h.Logger.ErrorContext(r.Context(), "failed to read document file", "path", filePath, "error", err)
		http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		return
	}

	h.Logger.InfoContext(r.Context(), "serving document", "name", docName, "bytes", len(content))

	w.Header().Set("Content-Type", "application/json")
	// Return as JSON with content field
	w.Write([]byte(`{"name":`))
	w.Write([]byte(`"` + docName + `"`))
	w.Write([]byte(`,"content":`))
	w.Write([]byte(`"`))
	// Escape content for JSON
	w.Write(jsonEscapeBytes(content))
	w.Write([]byte(`"}`))
}

// GetDocumentNames handles GET /documents — returns the list of available document names.
func (h *Handler) GetDocumentNames(w http.ResponseWriter, r *http.Request) {
	names := make([]string, 0, len(documentNameToFile))
	for name := range documentNameToFile {
		names = append(names, name)
	}

	h.Logger.InfoContext(r.Context(), "listing document names", "count", len(names))

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`[`))
	for i, name := range names {
		if i > 0 {
			w.Write([]byte(`,`))
		}
		w.Write([]byte(`"` + name + `"`))
	}
	w.Write([]byte(`]`))
}

// Health handles GET /health — simple health check.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok","service":"documentcatalogue"}`))
}

// jsonEscapeBytes escapes a byte slice for embedding in a JSON string value.
func jsonEscapeBytes(b []byte) []byte {
	var result []byte
	for _, c := range b {
		switch c {
		case '"':
			result = append(result, '\\', '"')
		case '\\':
			result = append(result, '\\', '\\')
		case '\n':
			result = append(result, '\\', 'n')
		case '\r':
			result = append(result, '\\', 'r')
		case '\t':
			result = append(result, '\\', 't')
		default:
			result = append(result, c)
		}
	}
	return result
}
