package model

import (
	"game-cycle-simulator/internal/domain"
	"math/rand"
)

type TimeModel struct {
	rng *rand.Rand
}

func NewTimeModel(rng *rand.Rand) *TimeModel {

	return &TimeModel{
		rng: rng,
	}

}

const fraudTimeAlpha = 0.5

// StateDuration возвращает случайное время пребывания в состоянии.

//

// В модели используется:

// - базовая скорость обработки mu;

// - алгоритмический коэффициент a;

// - параметр антифрод-фильтрации F.

//

// Чем больше mu и a, тем быстрее обработка.

// Чем больше F, тем больше дополнительных проверок и тем больше время обработки.

func (m *TimeModel) StateDuration(

	state domain.State,
	mu float64,
	algorithmFactorA float64,
	fraudFactorF float64,

) float64 {

	meanTime := m.stateMeanTime(state, mu, algorithmFactorA, fraudFactorF)
	// Экспоненциальная модель времени обработки.
	// Такая форма позволяет учитывать случайный характер длительности операций.
	return m.rng.ExpFloat64() * meanTime

}

// TotalProcessingTime считает суммарное время обработки по всей траектории.

func (m *TimeModel) TotalProcessingTime(

	path []domain.State,
	mu float64,
	algorithmFactorA float64,
	fraudFactorF float64,

) float64 {

	total := 0.0
	for _, state := range path {
		total += m.StateDuration(state, mu, algorithmFactorA, fraudFactorF)
	}
	return total

}

// stateMeanTime задаёт среднее время обработки для одного состояния.

//

// Базовая логика модели:

// baseMean = 1 / (mu * a)

//

// где:

// mu — базовая интенсивность обслуживания;

// a — коэффициент эффективности алгоритмов обработки.

//

// Антифрод-фильтрация увеличивает длительность обработки:

//

// fraudMultiplier = 1 + alpha * F

//

// где:

// F — параметр строгости антифрод-фильтрации;

// alpha — коэффициент влияния антифрода на время обработки.

//

// Итоговое среднее время:

//

// meanTime = baseMean * fraudMultiplier * stateWeight

func (m *TimeModel) stateMeanTime(

	state domain.State,
	mu float64,
	algorithmFactorA float64,
	fraudFactorF float64,

) float64 {

	baseMean := 1.0 / (mu * algorithmFactorA)
	fraudMultiplier := 1.0 + fraudTimeAlpha*fraudFactorF
	totalMean := baseMean * fraudMultiplier
	weight := stateWeight(state)
	return totalMean * weight

}

func stateWeight(state domain.State) float64 {

	switch state {
	case domain.Init:
		return 0.10
	case domain.Reserve:
		return 0.25
	case domain.Engine:
		return 0.35
	case domain.Settlement:
		return 0.20
	case domain.Wallet:
		return 0.08
	case domain.ErrorState:
		return 0.12
	default:
		return 0.10
	}

}
