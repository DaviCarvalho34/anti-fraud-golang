package services

import (
	"fmt"
	"sync"
	"time"
	
	"github.com/anti-fraud-golang/internal/models"
)

// InMemoryProfileStore implementação em memória do ProfileStore
type InMemoryProfileStore struct {
	profiles map[string]*models.UserProfile
	mu       sync.RWMutex
}

// NewInMemoryProfileStore cria uma nova instância
func NewInMemoryProfileStore() *InMemoryProfileStore {
	return &InMemoryProfileStore{
		profiles: make(map[string]*models.UserProfile),
	}
}

// GetUserProfile obtém o perfil de um usuário
func (s *InMemoryProfileStore) GetUserProfile(userID string) (*models.UserProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	profile, exists := s.profiles[userID]
	if !exists {
		return nil, fmt.Errorf("profile not found for user %s", userID)
	}
	
	return profile, nil
}

// UpdateUserProfile atualiza o perfil de um usuário
func (s *InMemoryProfileStore) UpdateUserProfile(profile *models.UserProfile) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.profiles[profile.UserID] = profile
	return nil
}

// CreateSampleProfile cria um perfil de exemplo para testes
func (s *InMemoryProfileStore) CreateSampleProfile(userID string) *models.UserProfile {
	profile := &models.UserProfile{
		UserID:              userID,
		AvgTransactionValue: 500.0,
		TotalTransactions:   150,
		FirstTransactionAt:  time.Now().AddDate(0, -6, 0),
		LastTransactionAt:   time.Now().Add(-24 * time.Hour),
		CommonLocations: []models.Location{
			{
				Country:   "BR",
				City:      "São Paulo",
				Latitude:  -23.5505,
				Longitude: -46.6333,
			},
		},
		CommonMerchants: []string{"Amazon", "Mercado Livre", "Magazine Luiza"},
		FraudHistory:    []models.FraudIncident{},
		TrustedDevices:  []string{"device-123"},
	}
	
	s.UpdateUserProfile(profile)
	return profile
}

// InMemoryBlacklistStore implementação em memória do BlacklistStore
type InMemoryBlacklistStore struct {
	entries map[string]map[string]*models.BlacklistEntry
	mu      sync.RWMutex
}

// NewInMemoryBlacklistStore cria uma nova instância
func NewInMemoryBlacklistStore() *InMemoryBlacklistStore {
	return &InMemoryBlacklistStore{
		entries: make(map[string]map[string]*models.BlacklistEntry),
	}
}

// IsBlacklisted verifica se um valor está na lista negra
func (s *InMemoryBlacklistStore) IsBlacklisted(entryType, value string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	typeEntries, exists := s.entries[entryType]
	if !exists {
		return false, nil
	}
	
	entry, exists := typeEntries[value]
	if !exists {
		return false, nil
	}
	
	// Verifica se a entrada está ativa e não expirou
	if !entry.IsActive {
		return false, nil
	}
	
	if entry.ExpiresAt != nil && time.Now().After(*entry.ExpiresAt) {
		return false, nil
	}
	
	return true, nil
}

// Add adiciona uma entrada na lista negra
func (s *InMemoryBlacklistStore) Add(entry *models.BlacklistEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.entries[entry.Type]; !exists {
		s.entries[entry.Type] = make(map[string]*models.BlacklistEntry)
	}
	
	s.entries[entry.Type][entry.Value] = entry
	return nil
}

// AddSampleBlacklist adiciona entradas de exemplo
func (s *InMemoryBlacklistStore) AddSampleBlacklist() {
	// Adiciona alguns exemplos
	s.Add(&models.BlacklistEntry{
		ID:       "bl-1",
		Type:     "user",
		Value:    "BLOCKED_USER_123",
		Reason:   "Múltiplas tentativas de fraude confirmadas",
		AddedAt:  time.Now(),
		IsActive: true,
	})
	
	s.Add(&models.BlacklistEntry{
		ID:       "bl-2",
		Type:     "card",
		Value:    "4567",
		Reason:   "Cartão reportado como roubado",
		AddedAt:  time.Now(),
		IsActive: true,
	})
	
	s.Add(&models.BlacklistEntry{
		ID:       "bl-3",
		Type:     "ip",
		Value:    "192.168.1.100",
		Reason:   "IP associado a atividade fraudulenta",
		AddedAt:  time.Now(),
		IsActive: true,
	})
}
