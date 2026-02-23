package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mikenunan/openobserve-poc/services/partyman/models"
	"github.com/mikenunan/openobserve-poc/services/partyman/store"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	Store  *store.CustomerStore
	Logger *slog.Logger
}

// GetCustomers handles GET /customers — returns all customers.
func (h *Handler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	customers := h.Store.GetAll()

	h.Logger.InfoContext(r.Context(), "listing customers", "count", len(customers))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// GetCustomer handles GET /customer?partyId=<id> — returns a single customer.
// Also handles HEAD /customer?partyId=<id> — existence check (no body).
func (h *Handler) GetCustomer(w http.ResponseWriter, r *http.Request) {
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

	customer, found := h.Store.GetByID(partyID)
	if !found {
		h.Logger.InfoContext(r.Context(), "customer not found", "partyId", partyID)
		http.Error(w, `{"error":"customer not found"}`, http.StatusNotFound)
		return
	}

	// HEAD requests get just the status code, no body
	if r.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	h.Logger.InfoContext(r.Context(), "retrieved customer", "partyId", partyID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// CreateCustomer handles POST /customer — creates a new customer.
func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var c models.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if c.FirstName == "" || c.LastName == "" {
		http.Error(w, `{"error":"firstName and lastName are required"}`, http.StatusBadRequest)
		return
	}

	created := h.Store.Create(c)

	h.Logger.InfoContext(r.Context(), "created customer", "partyId", created.PartyID, "name", created.FirstName+" "+created.LastName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// UpdateCustomer handles PUT /customer — updates an existing customer.
func (h *Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var c models.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, `{"error":"invalid JSON body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if c.PartyID == 0 {
		http.Error(w, `{"error":"partyId is required"}`, http.StatusBadRequest)
		return
	}
	if c.FirstName == "" || c.LastName == "" {
		http.Error(w, `{"error":"firstName and lastName are required"}`, http.StatusBadRequest)
		return
	}

	if !h.Store.Update(c) {
		h.Logger.InfoContext(r.Context(), "customer not found for update", "partyId", c.PartyID)
		http.Error(w, `{"error":"customer not found"}`, http.StatusNotFound)
		return
	}

	h.Logger.InfoContext(r.Context(), "updated customer", "partyId", c.PartyID, "name", c.FirstName+" "+c.LastName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

// Health handles GET /health — simple health check.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok","service":"partyman"}`))
}
