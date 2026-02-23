package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/mikenunan/openobserve-poc/services/partyman/models"
)

// CustomerStore provides thread-safe in-memory storage for customer records.
type CustomerStore struct {
	mu        sync.RWMutex
	customers map[int]models.Customer
	nextID    int
}

// New creates a CustomerStore and loads seed data from the given JSON file path.
func New(seedFilePath string) (*CustomerStore, error) {
	data, err := os.ReadFile(seedFilePath)
	if err != nil {
		return nil, fmt.Errorf("reading seed file %s: %w", seedFilePath, err)
	}

	var seeds []models.Customer
	if err := json.Unmarshal(data, &seeds); err != nil {
		return nil, fmt.Errorf("parsing seed file %s: %w", seedFilePath, err)
	}

	s := &CustomerStore{
		customers: make(map[int]models.Customer, len(seeds)),
	}

	maxID := 0
	for _, c := range seeds {
		s.customers[c.PartyID] = c
		if c.PartyID > maxID {
			maxID = c.PartyID
		}
	}
	s.nextID = maxID + 1

	return s, nil
}

// GetAll returns all customers as a slice.
func (s *CustomerStore) GetAll() []models.Customer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Customer, 0, len(s.customers))
	for _, c := range s.customers {
		result = append(result, c)
	}
	return result
}

// GetByID returns a single customer by PartyID, or false if not found.
func (s *CustomerStore) GetByID(partyID int) (models.Customer, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	c, ok := s.customers[partyID]
	return c, ok
}

// Exists returns true if a customer with the given PartyID exists.
func (s *CustomerStore) Exists(partyID int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.customers[partyID]
	return ok
}

// Create adds a new customer, auto-assigning the next PartyID.
// Returns the created customer (with PartyID populated).
func (s *CustomerStore) Create(c models.Customer) models.Customer {
	s.mu.Lock()
	defer s.mu.Unlock()

	c.PartyID = s.nextID
	s.nextID++
	s.customers[c.PartyID] = c
	return c
}

// Update replaces the customer record for the given PartyID.
// Returns false if the PartyID does not exist.
func (s *CustomerStore) Update(c models.Customer) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.customers[c.PartyID]; !ok {
		return false
	}
	s.customers[c.PartyID] = c
	return true
}
