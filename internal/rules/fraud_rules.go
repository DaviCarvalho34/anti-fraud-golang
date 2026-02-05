package rules

import (
	"math"
	"time"
	
	"github.com/anti-fraud-golang/internal/models"
)

// HighAmountRule detecta transações de valor alto
type HighAmountRule struct {
	ID          string
	Name        string
	Enabled     bool
	Weight      int
	Threshold   float64
}

func (r *HighAmountRule) GetID() string { 
	if r.ID == "" {
		return "high_amount_rule"
	}
	return r.ID
}
func (r *HighAmountRule) GetName() string { 
	if r.Name == "" {
		return "High Amount Transaction"
	}
	return r.Name
}
func (r *HighAmountRule) GetWeight() int { 
	if r.Weight == 0 {
		return 25
	}
	return r.Weight
}
func (r *HighAmountRule) IsEnabled() bool { 
	return true
}

func (r *HighAmountRule) Evaluate(transaction *models.Transaction, profile *models.UserProfile) RuleResult {
	threshold := 10000.0
	if r.Threshold > 0 {
		threshold = r.Threshold
	}
	
	triggered := transaction.Amount > threshold
	score := 0
	
	if triggered {
		// Score baseado em quanto excede o threshold
		factor := transaction.Amount / threshold
		if factor > 5 {
			score = r.GetWeight()
		} else if factor > 3 {
			score = int(float64(r.GetWeight()) * 0.8)
		} else {
			score = int(float64(r.GetWeight()) * 0.6)
		}
	}
	
	return RuleResult{
		RuleID:      r.GetID(),
		RuleName:    r.GetName(),
		Triggered:   triggered,
		Score:       score,
		Description: "Transação com valor acima do limite normal",
		Details: map[string]interface{}{
			"amount":    transaction.Amount,
			"threshold": threshold,
		},
	}
}

// VelocityRule detecta múltiplas transações em curto período
type VelocityRule struct{}

func (r *VelocityRule) GetID() string { return "velocity_rule" }
func (r *VelocityRule) GetName() string { return "Transaction Velocity" }
func (r *VelocityRule) GetWeight() int { return 20 }
func (r *VelocityRule) IsEnabled() bool { return true }

func (r *VelocityRule) Evaluate(transaction *models.Transaction, profile *models.UserProfile) RuleResult {
	// Simula verificação de velocidade
	// Em produção, isso consultaria um cache/database
	triggered := false
	score := 0
	
	if profile != nil && profile.LastTransactionAt.After(time.Time{}) {
		timeDiff := transaction.Timestamp.Sub(profile.LastTransactionAt)
		
		// Se houver outra transação nos últimos 5 minutos
		if timeDiff < 5*time.Minute {
			triggered = true
			score = r.GetWeight()
		}
	}
	
	return RuleResult{
		RuleID:      r.GetID(),
		RuleName:    r.GetName(),
		Triggered:   triggered,
		Score:       score,
		Description: "Múltiplas transações em curto período",
		Details:     map[string]interface{}{},
	}
}

// GeoVelocityRule detecta mudanças geográficas impossíveis
type GeoVelocityRule struct{}

func (r *GeoVelocityRule) GetID() string { return "geo_velocity_rule" }
func (r *GeoVelocityRule) GetName() string { return "Geographical Velocity" }
func (r *GeoVelocityRule) GetWeight() int { return 30 }
func (r *GeoVelocityRule) IsEnabled() bool { return true }

func (r *GeoVelocityRule) Evaluate(transaction *models.Transaction, profile *models.UserProfile) RuleResult {
	triggered := false
	score := 0
	
	if profile != nil && len(profile.CommonLocations) > 0 {
		lastLocation := profile.CommonLocations[len(profile.CommonLocations)-1]
		
		// Calcula distância entre localizações
		distance := calculateDistance(
			lastLocation.Latitude, lastLocation.Longitude,
			transaction.Location.Latitude, transaction.Location.Longitude,
		)
		
		timeDiff := transaction.Timestamp.Sub(profile.LastTransactionAt).Hours()
		
		// Velocidade em km/h
		if timeDiff > 0 {
			speed := distance / timeDiff
			
			// Se a velocidade for maior que 900 km/h (velocidade de avião)
			if speed > 900 && distance > 100 {
				triggered = true
				score = r.GetWeight()
			}
		}
	}
	
	return RuleResult{
		RuleID:      r.GetID(),
		RuleName:    r.GetName(),
		Triggered:   triggered,
		Score:       score,
		Description: "Mudança geográfica impossível detectada",
		Details:     map[string]interface{}{},
	}
}

// UnusualHourRule detecta transações em horários incomuns
type UnusualHourRule struct{}

func (r *UnusualHourRule) GetID() string { return "unusual_hour_rule" }
func (r *UnusualHourRule) GetName() string { return "Unusual Hour Transaction" }
func (r *UnusualHourRule) GetWeight() int { return 10 }
func (r *UnusualHourRule) IsEnabled() bool { return true }

func (r *UnusualHourRule) Evaluate(transaction *models.Transaction, profile *models.UserProfile) RuleResult {
	hour := transaction.Timestamp.Hour()
	
	// Considera 23h às 5h como horário suspeito
	triggered := hour >= 23 || hour <= 5
	score := 0
	
	if triggered {
		score = r.GetWeight()
	}
	
	return RuleResult{
		RuleID:      r.GetID(),
		RuleName:    r.GetName(),
		Triggered:   triggered,
		Score:       score,
		Description: "Transação realizada em horário incomum",
		Details: map[string]interface{}{
			"hour": hour,
		},
	}
}

// NewUserRule detecta usuários novos com transações altas
type NewUserRule struct{}

func (r *NewUserRule) GetID() string { return "new_user_rule" }
func (r *NewUserRule) GetName() string { return "New User High Transaction" }
func (r *NewUserRule) GetWeight() int { return 15 }
func (r *NewUserRule) IsEnabled() bool { return true }

func (r *NewUserRule) Evaluate(transaction *models.Transaction, profile *models.UserProfile) RuleResult {
	triggered := false
	score := 0
	
	if profile != nil {
		// Se o usuário tem menos de 7 dias e faz transação alta
		accountAge := time.Since(profile.FirstTransactionAt).Hours() / 24
		
		if accountAge < 7 && transaction.Amount > 5000 {
			triggered = true
			score = r.GetWeight()
		}
	} else {
		// Usuário completamente novo
		if transaction.Amount > 3000 {
			triggered = true
			score = r.GetWeight()
		}
	}
	
	return RuleResult{
		RuleID:      r.GetID(),
		RuleName:    r.GetName(),
		Triggered:   triggered,
		Score:       score,
		Description: "Novo usuário com transação de valor elevado",
		Details:     map[string]interface{}{},
	}
}

// RoundAmountRule detecta valores redondos suspeitos
type RoundAmountRule struct{}

func (r *RoundAmountRule) GetID() string { return "round_amount_rule" }
func (r *RoundAmountRule) GetName() string { return "Suspicious Round Amount" }
func (r *RoundAmountRule) GetWeight() int { return 5 }
func (r *RoundAmountRule) IsEnabled() bool { return true }

func (r *RoundAmountRule) Evaluate(transaction *models.Transaction, profile *models.UserProfile) RuleResult {
	// Verifica se é um valor redondo e alto
	amount := transaction.Amount
	triggered := false
	score := 0
	
	// Verifica se é múltiplo de 1000 e maior que 5000
	if amount >= 5000 && math.Mod(amount, 1000) == 0 {
		triggered = true
		score = r.GetWeight()
	}
	
	return RuleResult{
		RuleID:      r.GetID(),
		RuleName:    r.GetName(),
		Triggered:   triggered,
		Score:       score,
		Description: "Valor redondo suspeito",
		Details: map[string]interface{}{
			"amount": amount,
		},
	}
}

// MultipleFailedAttemptsRule detecta múltiplas tentativas falhadas
type MultipleFailedAttemptsRule struct{}

func (r *MultipleFailedAttemptsRule) GetID() string { return "multiple_failed_attempts_rule" }
func (r *MultipleFailedAttemptsRule) GetName() string { return "Multiple Failed Attempts" }
func (r *MultipleFailedAttemptsRule) GetWeight() int { return 25 }
func (r *MultipleFailedAttemptsRule) IsEnabled() bool { return true }

func (r *MultipleFailedAttemptsRule) Evaluate(transaction *models.Transaction, profile *models.UserProfile) RuleResult {
	// Esta regra seria implementada com histórico de tentativas
	// Por agora, retorna não-triggered
	return RuleResult{
		RuleID:      r.GetID(),
		RuleName:    r.GetName(),
		Triggered:   false,
		Score:       0,
		Description: "Múltiplas tentativas falhadas detectadas",
		Details:     map[string]interface{}{},
	}
}

// calculateDistance calcula a distância entre dois pontos geográficos em km
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0 // km
	
	dLat := degToRad(lat2 - lat1)
	dLon := degToRad(lon2 - lon1)
	
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degToRad(lat1))*math.Cos(degToRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	
	return earthRadius * c
}

func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
