package models

import "time"

// BlacklistEntry entrada da lista negra
type BlacklistEntry struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"` // "card", "user", "device", "ip"
	Value      string    `json:"value"`
	Reason     string    `json:"reason"`
	AddedAt    time.Time `json:"added_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	IsActive   bool      `json:"is_active"`
}

// VelocityCheck verifica velocidade de transações
type VelocityCheck struct {
	UserID             string    `json:"user_id"`
	TransactionCount   int       `json:"transaction_count"`
	TotalAmount        float64   `json:"total_amount"`
	TimeWindow         int       `json:"time_window_minutes"`
	LastTransactionAt  time.Time `json:"last_transaction_at"`
}

// GeoVelocity analisa velocidade geográfica
type GeoVelocity struct {
	UserID           string    `json:"user_id"`
	PreviousLocation Location  `json:"previous_location"`
	CurrentLocation  Location  `json:"current_location"`
	Distance         float64   `json:"distance_km"`
	TimeDiff         int       `json:"time_diff_minutes"`
	Speed            float64   `json:"speed_kmh"`
	IsPossible       bool      `json:"is_possible"`
}

// AlertConfig configuração de alertas
type AlertConfig struct {
	ID                  string  `json:"id"`
	MaxAmountThreshold  float64 `json:"max_amount_threshold"`
	VelocityThreshold   int     `json:"velocity_threshold"`
	GeoVelocityLimit    float64 `json:"geo_velocity_limit_kmh"`
	NightHourStart      int     `json:"night_hour_start"`
	NightHourEnd        int     `json:"night_hour_end"`
	HighRiskCountries   []string `json:"high_risk_countries"`
}

// TransactionPattern padrão de transação
type TransactionPattern struct {
	PatternID      string    `json:"pattern_id"`
	PatternType    string    `json:"pattern_type"` // "suspicious", "normal", "fraud"
	Description    string    `json:"description"`
	Occurrences    int       `json:"occurrences"`
	FirstDetected  time.Time `json:"first_detected"`
	LastDetected   time.Time `json:"last_detected"`
	Confidence     float64   `json:"confidence"` // 0-1
}
