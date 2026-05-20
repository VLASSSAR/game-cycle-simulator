package fraud

import (
	"math"

	"game-cycle-simulator/internal/domain"
)

// Score рассчитывает интегральный показатель риска fraud-score.
//
// В рамках модели предполагается, что признаки нормированы
// в диапазоне [0, 1]. Итоговое значение также ограничивается
// диапазоном [0, 1].
func (m *Model) Score(features domain.FraudFeatures) float64 {
	weights := m.config.Weights
	sumWeights := weightsSum(weights)

	if sumWeights <= 0 {
		return 0
	}

	weightedScore :=
		weights.BetFrequencyWeight*clamp01(features.BetFrequency) +
			weights.AvgBetSizeWeight*clamp01(features.AvgBetSize) +
			weights.SessionCountWeight*clamp01(features.SessionCount) +
			weights.GeoRiskWeight*clamp01(features.GeoRisk) +
			weights.DeviceRiskWeight*clamp01(features.DeviceRisk) +
			weights.HistoryRiskWeight*clamp01(features.HistoryRisk) +
			weights.LimitViolationRiskWeight*clamp01(features.LimitViolationRisk)

	normalizedScore := weightedScore / sumWeights
	strictScore := normalizedScore * m.config.Strictness

	return clamp01(strictScore)
}

func clamp01(value float64) float64 {
	return math.Max(0, math.Min(1, value))
}
