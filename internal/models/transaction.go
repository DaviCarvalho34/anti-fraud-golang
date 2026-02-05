package models

import "time"

// Transaction representa uma transação financeira
type Transaction struct {
	ID          string     `json:"transaction_id"`
	UserID      string     `json:"user_id"`
	Amount      float64    `json:"amount"`
	Currency    string     `json:"currency"`
	Merchant    string     `json:"merchant"`
	Location    Location   `json:"location"`
	DeviceInfo  DeviceInfo `json:"device_info,omitempty"`
	Timestamp   time.Time  `json:"timestamp"`
	CardLast4   string     `json:"card_last4,omitempty"`
	CardType    string     `json:"card_type,omitempty"`
	Description string     `json:"description,omitempty"`
}

// Location representa a localização geográfica
type Location struct {
	Country   string  `json:"country"`
	City      string  `json:"city"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IPAddress string  `json:"ip_address,omitempty"`
}

// DeviceInfo contém informações do dispositivo
type DeviceInfo struct {
	DeviceID     string `json:"device_id"`
	DeviceType   string `json:"device_type"`
	OS           string `json:"os"`
	Browser      string `json:"browser,omitempty"`
	UserAgent    string `json:"user_agent,omitempty"`
	Fingerprint  string `json:"fingerprint,omitempty"`
}

// FraudAnalysisResult resultado da análise de fraude
type FraudAnalysisResult struct {
	TransactionID   string              `json:"transaction_id"`
	RiskScore       int                 `json:"risk_score"`
	RiskLevel       RiskLevel           `json:"risk_level"`
	Decision        Decision            `json:"decision"`
	Reasons         []string            `json:"reasons"`
	RulesTriggered  []string            `json:"rules_triggered"`
	Details         map[string]interface{} `json:"details,omitempty"`
	AnalyzedAt      time.Time           `json:"analyzed_at"`
	ProcessingTime  int64               `json:"processing_time_ms"`
}

// RiskLevel níveis de risco
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "LOW"
	RiskLevelMedium RiskLevel = "MEDIUM"
	RiskLevelHigh   RiskLevel = "HIGH"
)

// Decision decisão sobre a transação
type Decision string

const (
	DecisionApproved Decision = "APPROVED"
	DecisionReview   Decision = "REVIEW"
	DecisionBlocked  Decision = "BLOCKED"
)

// UserProfile perfil do usuário com histórico
type UserProfile struct {
	UserID              string               `json:"user_id"`
	AvgTransactionValue float64              `json:"avg_transaction_value"`
	TotalTransactions   int                  `json:"total_transactions"`
	FirstTransactionAt  time.Time            `json:"first_transaction_at"`
	LastTransactionAt   time.Time            `json:"last_transaction_at"`
	CommonLocations     []Location           `json:"common_locations"`
	CommonMerchants     []string             `json:"common_merchants"`
	FraudHistory        []FraudIncident      `json:"fraud_history"`
	TrustedDevices      []string             `json:"trusted_devices"`
}

// FraudIncident representa um incidente de fraude
type FraudIncident struct {
	IncidentID     string    `json:"incident_id"`
	TransactionID  string    `json:"transaction_id"`
	DetectedAt     time.Time `json:"detected_at"`
	ConfirmedFraud bool      `json:"confirmed_fraud"`
	Amount         float64   `json:"amount"`
	Description    string    `json:"description"`
}

// Rule representa uma regra de detecção de fraude
type Rule struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Enabled     bool     `json:"enabled"`
	Priority    int      `json:"priority"`
	ScoreWeight int      `json:"score_weight"`
}
