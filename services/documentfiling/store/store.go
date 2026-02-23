package store

import (
	"sync"

	"github.com/mikenunan/openobserve-poc/services/documentfiling/models"
)

// AffiliationStore provides thread-safe in-memory storage for document affiliations.
type AffiliationStore struct {
	mu           sync.RWMutex
	affiliations map[int][]string // partyID -> []documentName
}

// New creates an AffiliationStore with optional pre-seeded affiliations.
func New(seeds []models.Affiliation) *AffiliationStore {
	s := &AffiliationStore{
		affiliations: make(map[int][]string),
	}
	for _, seed := range seeds {
		s.affiliations[seed.PartyID] = append(s.affiliations[seed.PartyID], seed.DocumentName)
	}
	return s
}

// Exists checks if a specific affiliation already exists.
func (s *AffiliationStore) Exists(partyID int, documentName string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	docs, ok := s.affiliations[partyID]
	if !ok {
		return false
	}
	for _, d := range docs {
		if d == documentName {
			return true
		}
	}
	return false
}

// Add creates a new affiliation. Returns false if it already exists.
func (s *AffiliationStore) Add(partyID int, documentName string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check for duplicate
	for _, d := range s.affiliations[partyID] {
		if d == documentName {
			return false
		}
	}

	s.affiliations[partyID] = append(s.affiliations[partyID], documentName)
	return true
}

// GetByPartyID returns all document names affiliated with the given PartyID.
func (s *AffiliationStore) GetByPartyID(partyID int) []models.Affiliation {
	s.mu.RLock()
	defer s.mu.RUnlock()

	docs := s.affiliations[partyID]
	result := make([]models.Affiliation, 0, len(docs))
	for _, d := range docs {
		result = append(result, models.Affiliation{
			PartyID:      partyID,
			DocumentName: d,
		})
	}
	return result
}
