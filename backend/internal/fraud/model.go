package fraud

import "game-cycle-simulator/internal/domain"

// Model реализует антифрод-модель,
// основанную на расчете интегрального показателя риска fraud-score.
type Model struct {
	config domain.FraudConfig
}

// NewModel создает антифрод-модель с заданной конфигурацией.
func NewModel(config domain.FraudConfig) *Model {
	return &Model{
		config: normalizeConfig(config),
	}
}

// NewDefaultModel создает антифрод-модель с настройками по умолчанию.
func NewDefaultModel() *Model {
	return NewModel(DefaultConfig())
}

// DefaultConfig возвращает базовую конфигурацию антифрод-модели.
func DefaultConfig() domain.FraudConfig {
	return domain.FraudConfig{
		Threshold:  0.6,
		Strictness: 1.0,
		Weights: domain.FraudWeights{
			BetFrequencyWeight:       0.20,
			AvgBetSizeWeight:         0.15,
			SessionCountWeight:       0.10,
			GeoRiskWeight:            0.15,
			DeviceRiskWeight:         0.15,
			HistoryRiskWeight:        0.15,
			LimitViolationRiskWeight: 0.10,
		},
	}
}

// Evaluate выполняет антифрод-проверку операции.
//
// isFraud — истинный статус операции в рамках имитационной модели.
// Он нужен не для принятия решения, а для оценки качества антифрод-фильтра:
// false positive и false negative.
func (m *Model) Evaluate(
	features domain.FraudFeatures,
	isFraud bool,
) domain.FraudResult {
	score := m.Score(features)
	isSuspicious := score >= m.config.Threshold

	return domain.FraudResult{
		Score:           score,
		IsSuspicious:    isSuspicious,
		IsFraud:         isFraud,
		IsFalsePositive: isSuspicious && !isFraud,
		IsFalseNegative: !isSuspicious && isFraud,
	}
}

// Config возвращает текущую конфигурацию модели.
func (m *Model) Config() domain.FraudConfig {
	return m.config
}

func normalizeConfig(config domain.FraudConfig) domain.FraudConfig {
	defaultConfig := DefaultConfig()

	if config.Threshold <= 0 || config.Threshold > 1 {
		config.Threshold = defaultConfig.Threshold
	}

	if config.Strictness <= 0 {
		config.Strictness = defaultConfig.Strictness
	}

	if weightsSum(config.Weights) <= 0 {
		config.Weights = defaultConfig.Weights
	}

	return config
}

func weightsSum(weights domain.FraudWeights) float64 {
	return weights.BetFrequencyWeight +
		weights.AvgBetSizeWeight +
		weights.SessionCountWeight +
		weights.GeoRiskWeight +
		weights.DeviceRiskWeight +
		weights.HistoryRiskWeight +
		weights.LimitViolationRiskWeight
}
