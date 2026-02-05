package services

import (
	"time"
	
	"github.com/anti-fraud-golang/internal/models"
	"github.com/anti-fraud-golang/internal/rules"
	"github.com/anti-fraud-golang/pkg/utils"
)

// FraudDetectionService serviço de detecção de fraude
type FraudDetectionService struct {
	ruleEngine    *rules.RuleEngine
	profileStore  ProfileStore
	blacklistStore BlacklistStore
}

// ProfileStore interface para armazenamento de perfis
type ProfileStore interface {
	GetUserProfile(userID string) (*models.UserProfile, error)
	UpdateUserProfile(profile *models.UserProfile) error
}

// BlacklistStore interface para lista negra
type BlacklistStore interface {
	IsBlacklisted(entryType, value string) (bool, error)
	Add(entry *models.BlacklistEntry) error
}

// NewFraudDetectionService cria uma nova instância do serviço
func NewFraudDetectionService(profileStore ProfileStore, blacklistStore BlacklistStore) *FraudDetectionService {
	return &FraudDetectionService{
		ruleEngine:    rules.NewRuleEngine(),
		profileStore:  profileStore,
		blacklistStore: blacklistStore,
	}
}

// AnalyzeTransaction analisa uma transação para detectar fraude
func (s *FraudDetectionService) AnalyzeTransaction(transaction *models.Transaction) (*models.FraudAnalysisResult, error) {
	startTime := time.Now()
	
	// Define timestamp se não estiver definido
	if transaction.Timestamp.IsZero() {
		transaction.Timestamp = time.Now()
	}
	
	// Verifica lista negra primeiro
	blacklisted, err := s.checkBlacklist(transaction)
	if err != nil {
		return nil, err
	}
	
	if blacklisted {
		return s.createBlockedResult(transaction, "Entidade na lista negra", startTime), nil
	}
	
	// Obtém perfil do usuário
	profile, err := s.profileStore.GetUserProfile(transaction.UserID)
	if err != nil {
		// Se não encontrar perfil, cria um vazio (usuário novo)
		profile = nil
	}
	
	// Avalia todas as regras
	ruleResults := s.ruleEngine.Evaluate(transaction, profile)
	
	// Calcula score total
	totalScore := s.ruleEngine.CalculateTotalScore(ruleResults)
	
	// Determina nível de risco e decisão
	riskLevel := rules.GetRiskLevel(totalScore)
	decision := rules.GetDecision(riskLevel)
	
	// Extrai razões e regras acionadas
	reasons := make([]string, 0)
	rulesTriggered := make([]string, 0)
	
	for _, result := range ruleResults {
		if result.Triggered {
			reasons = append(reasons, result.Description)
			rulesTriggered = append(rulesTriggered, result.RuleName)
		}
	}
	
	// Cria resultado da análise
	analysisResult := &models.FraudAnalysisResult{
		TransactionID:  transaction.ID,
		RiskScore:      totalScore,
		RiskLevel:      riskLevel,
		Decision:       decision,
		Reasons:        reasons,
		RulesTriggered: rulesTriggered,
		AnalyzedAt:     time.Now(),
		ProcessingTime: time.Since(startTime).Milliseconds(),
		Details: map[string]interface{}{
			"user_id":  transaction.UserID,
			"amount":   transaction.Amount,
			"merchant": transaction.Merchant,
		},
	}
	
	return analysisResult, nil
}

// checkBlacklist verifica se algum elemento da transação está na lista negra
func (s *FraudDetectionService) checkBlacklist(transaction *models.Transaction) (bool, error) {
	// Verifica usuário
	blacklisted, err := s.blacklistStore.IsBlacklisted("user", transaction.UserID)
	if err != nil {
		return false, err
	}
	if blacklisted {
		return true, nil
	}
	
	// Verifica cartão
	if transaction.CardLast4 != "" {
		blacklisted, err = s.blacklistStore.IsBlacklisted("card", transaction.CardLast4)
		if err != nil {
			return false, err
		}
		if blacklisted {
			return true, nil
		}
	}
	
	// Verifica IP
	if transaction.Location.IPAddress != "" {
		blacklisted, err = s.blacklistStore.IsBlacklisted("ip", transaction.Location.IPAddress)
		if err != nil {
			return false, err
		}
		if blacklisted {
			return true, nil
		}
	}
	
	// Verifica dispositivo
	if transaction.DeviceInfo.DeviceID != "" {
		blacklisted, err = s.blacklistStore.IsBlacklisted("device", transaction.DeviceInfo.DeviceID)
		if err != nil {
			return false, err
		}
		if blacklisted {
			return true, nil
		}
	}
	
	return false, nil
}

// createBlockedResult cria um resultado de bloqueio
func (s *FraudDetectionService) createBlockedResult(transaction *models.Transaction, reason string, startTime time.Time) *models.FraudAnalysisResult {
	return &models.FraudAnalysisResult{
		TransactionID:  transaction.ID,
		RiskScore:      100,
		RiskLevel:      models.RiskLevelHigh,
		Decision:       models.DecisionBlocked,
		Reasons:        []string{reason},
		RulesTriggered: []string{"Blacklist Check"},
		AnalyzedAt:     time.Now(),
		ProcessingTime: time.Since(startTime).Milliseconds(),
		Details: map[string]interface{}{
			"blocked_reason": reason,
		},
	}
}

// GetTransactionAnalytics retorna analytics de transações
func (s *FraudDetectionService) GetTransactionAnalytics(userID string) (*TransactionAnalytics, error) {
	profile, err := s.profileStore.GetUserProfile(userID)
	if err != nil {
		return nil, err
	}
	
	if profile == nil {
		return &TransactionAnalytics{
			UserID:             userID,
			TotalTransactions:  0,
			FraudCount:        0,
			FraudRate:         0,
		}, nil
	}
	
	fraudCount := len(profile.FraudHistory)
	fraudRate := 0.0
	if profile.TotalTransactions > 0 {
		fraudRate = float64(fraudCount) / float64(profile.TotalTransactions) * 100
	}
	
	return &TransactionAnalytics{
		UserID:             userID,
		TotalTransactions:  profile.TotalTransactions,
		AverageAmount:      profile.AvgTransactionValue,
		FraudCount:        fraudCount,
		FraudRate:         fraudRate,
		LastTransactionAt: profile.LastTransactionAt,
	}, nil
}

// TransactionAnalytics estatísticas de transações
type TransactionAnalytics struct {
	UserID             string    `json:"user_id"`
	TotalTransactions  int       `json:"total_transactions"`
	AverageAmount      float64   `json:"average_amount"`
	FraudCount        int       `json:"fraud_count"`
	FraudRate         float64   `json:"fraud_rate"`
	LastTransactionAt time.Time `json:"last_transaction_at"`
}
