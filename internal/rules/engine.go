package rules

import (
	"github.com/anti-fraud-golang/internal/models"
)

// RuleEngine motor de regras para detecção de fraude
type RuleEngine struct {
	rules []FraudRule
}

// FraudRule interface para regras de fraude
type FraudRule interface {
	Evaluate(transaction *models.Transaction, profile *models.UserProfile) RuleResult
	GetID() string
	GetName() string
	GetWeight() int
	IsEnabled() bool
}

// RuleResult resultado da avaliação de uma regra
type RuleResult struct {
	RuleID      string
	RuleName    string
	Triggered   bool
	Score       int
	Description string
	Details     map[string]interface{}
}

// NewRuleEngine cria uma nova instância do motor de regras
func NewRuleEngine() *RuleEngine {
	engine := &RuleEngine{
		rules: make([]FraudRule, 0),
	}
	
	// Registra todas as regras
	engine.RegisterRule(&HighAmountRule{})
	engine.RegisterRule(&VelocityRule{})
	engine.RegisterRule(&GeoVelocityRule{})
	engine.RegisterRule(&UnusualHourRule{})
	engine.RegisterRule(&NewUserRule{})
	engine.RegisterRule(&RoundAmountRule{})
	engine.RegisterRule(&MultipleFailedAttemptsRule{})
	
	return engine
}

// RegisterRule registra uma nova regra
func (e *RuleEngine) RegisterRule(rule FraudRule) {
	e.rules = append(e.rules, rule)
}

// Evaluate avalia todas as regras contra uma transação
func (e *RuleEngine) Evaluate(transaction *models.Transaction, profile *models.UserProfile) []RuleResult {
	results := make([]RuleResult, 0)
	
	for _, rule := range e.rules {
		if !rule.IsEnabled() {
			continue
		}
		
		result := rule.Evaluate(transaction, profile)
		if result.Triggered {
			results = append(results, result)
		}
	}
	
	return results
}

// CalculateTotalScore calcula a pontuação total de risco
func (e *RuleEngine) CalculateTotalScore(results []RuleResult) int {
	totalScore := 0
	for _, result := range results {
		totalScore += result.Score
	}
	
	// Limita o score entre 0 e 100
	if totalScore > 100 {
		totalScore = 100
	}
	
	return totalScore
}

// GetRiskLevel determina o nível de risco baseado no score
func GetRiskLevel(score int) models.RiskLevel {
	if score <= 30 {
		return models.RiskLevelLow
	} else if score <= 70 {
		return models.RiskLevelMedium
	}
	return models.RiskLevelHigh
}

// GetDecision determina a decisão baseada no nível de risco
func GetDecision(riskLevel models.RiskLevel) models.Decision {
	switch riskLevel {
	case models.RiskLevelLow:
		return models.DecisionApproved
	case models.RiskLevelMedium:
		return models.DecisionReview
	case models.RiskLevelHigh:
		return models.DecisionBlocked
	default:
		return models.DecisionReview
	}
}
