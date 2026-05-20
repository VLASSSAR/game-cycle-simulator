package fraud

import (
	"math/rand"

	"game-cycle-simulator/internal/domain"
)

// OperationProfile описывает сгенерированный профиль одной игровой операции.
//
// В рамках имитационной модели мы генерируем:
// - признаки пользовательской активности;
// - истинный статус операции: мошенническая или нет.
//
// Истинный статус нужен не для принятия решения,
// а для оценки качества антифрод-фильтра:
// false positive и false negative.
type OperationProfile struct {
	Features domain.FraudFeatures
	IsFraud  bool
}

// FeatureGenerator генерирует признаки пользовательской активности
// для имитационной модели антифрод-анализа.
type FeatureGenerator struct {
	rng *rand.Rand
}

// NewFeatureGenerator создает генератор признаков.
func NewFeatureGenerator(rng *rand.Rand) *FeatureGenerator {
	return &FeatureGenerator{
		rng: rng,
	}
}

// Generate формирует профиль одной игровой операции.
//
// fraudProbability задает вероятность того, что операция
// является мошеннической в рамках имитационной модели.
func (g *FeatureGenerator) Generate(fraudProbability float64) OperationProfile {
	fraudProbability = clamp01(fraudProbability)

	isFraud := g.rng.Float64() < fraudProbability

	if isFraud {
		return OperationProfile{
			Features: g.generateFraudFeatures(),
			IsFraud:  true,
		}
	}

	return OperationProfile{
		Features: g.generateNormalFeatures(),
		IsFraud:  false,
	}
}

// generateNormalFeatures генерирует признаки обычной пользовательской активности.
func (g *FeatureGenerator) generateNormalFeatures() domain.FraudFeatures {
	return domain.FraudFeatures{
		BetFrequency:       g.randomInRange(0.05, 0.45),
		AvgBetSize:         g.randomInRange(0.05, 0.55),
		SessionCount:       g.randomInRange(0.05, 0.50),
		GeoRisk:            g.randomInRange(0.00, 0.35),
		DeviceRisk:         g.randomInRange(0.00, 0.35),
		HistoryRisk:        g.randomInRange(0.00, 0.30),
		LimitViolationRisk: g.randomInRange(0.00, 0.25),
	}
}

// generateFraudFeatures генерирует признаки подозрительной активности.
//
// Для мошеннических операций признаки в среднем имеют более высокие значения,
// что повышает итоговый fraud-score.
func (g *FeatureGenerator) generateFraudFeatures() domain.FraudFeatures {
	return domain.FraudFeatures{
		BetFrequency:       g.randomInRange(0.55, 1.00),
		AvgBetSize:         g.randomInRange(0.45, 1.00),
		SessionCount:       g.randomInRange(0.45, 1.00),
		GeoRisk:            g.randomInRange(0.40, 1.00),
		DeviceRisk:         g.randomInRange(0.40, 1.00),
		HistoryRisk:        g.randomInRange(0.45, 1.00),
		LimitViolationRisk: g.randomInRange(0.40, 1.00),
	}
}

func (g *FeatureGenerator) randomInRange(minValue float64, maxValue float64) float64 {
	if maxValue <= minValue {
		return minValue
	}

	return minValue + g.rng.Float64()*(maxValue-minValue)
}
