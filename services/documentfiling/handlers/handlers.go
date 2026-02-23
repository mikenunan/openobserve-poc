package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mikenunan/openobserve-poc/services/documentfiling/models"
	"github.com/mikenunan/openobserve-poc/services/documentfiling/store"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	Store  *store.AffiliationStore
	Logger *slog.Logger
}

// CreateAffiliation handles POST /affiliation — creates a document-party association.
func (h *Handler) CreateAffiliation(w http.ResponseWriter, r *http.Request) {
	var aff models.Affiliation
	if err := json.NewDecoder(r.Body).Decode(&aff); err != nil {
		http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if aff.PartyID == 0 || aff.DocumentName == "" {
		http.Error(w, `{"error":"partyId and documentName are required"}`, http.StatusBadRequest)
		return
	}

	// Cross-service validation: check customer exists
	if err := ValidateCustomerExists(r.Context(), aff.PartyID, h.Logger); err != nil {
		h.Logger.WarnContext(r.Context(), "customer validation failed", "partyId", aff.PartyID, "error", err)
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnprocessableEntity)
		return
	}

	// Cross-service validation: check document exists
	if err := ValidateDocumentExists(r.Context(), aff.DocumentName, h.Logger); err != nil {
		h.Logger.WarnContext(r.Context(), "document validation failed", "documentName", aff.DocumentName, "error", err)
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnprocessableEntity)
		return
	}

	// Attempt to add (returns false if duplicate)
	if !h.Store.Add(aff.PartyID, aff.DocumentName) {
		h.Logger.InfoContext(r.Context(), "duplicate affiliation rejected", "partyId", aff.PartyID, "documentName", aff.DocumentName)
		http.Error(w, `{"error":"affiliation already exists"}`, http.StatusConflict)
		return
	}

	h.Logger.InfoContext(r.Context(), "affiliation created", "partyId", aff.PartyID, "documentName", aff.DocumentName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(aff)
}

// GetDocuments handles GET /documents?partyId=<id> — returns affiliated document names.
func (h *Handler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("partyId")
	if idStr == "" {
		http.Error(w, `{"error":"partyId query parameter is required"}`, http.StatusBadRequest)
		return
	}

	partyID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"partyId must be a number"}`, http.StatusBadRequest)
		return
	}

	affiliations := h.Store.GetByPartyID(partyID)

	h.Logger.InfoContext(r.Context(), "listing affiliations", "partyId", partyID, "count", len(affiliations))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(affiliations)
}

// Health handles GET /health — simple health check.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok","service":"documentfiling"}`))
}
