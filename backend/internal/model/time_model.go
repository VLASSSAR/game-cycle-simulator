package model

import (
	"math/rand"

	"game-cycle-simulator/internal/domain"
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
// В модели используется базовая скорость обработки Mu,
// алгоритмический коэффициент a и параметр антифрод-фильтрации F.
//
// Чем больше Mu и a, тем быстрее обработка.
// Чем больше F, тем больше дополнительные проверки и тем больше время.
func (m *TimeModel) StateDuration(
	state domain.State,
	mu float64,
	algorithmFactorA float64,
	fraudFactorF float64,
) float64 {
	meanTime := m.stateMeanTime(state, mu, algorithmFactorA, fraudFactorF)

	// Экспоненциальная модель времени обработки.
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

// stateMeanTime задаёт среднее время для одного состояния.
//
// Базовая формула:
// T_proc = 1 / (mu * a * (1 + alpha * F))
//
// Для state-based модели мы используем веса по состояниям,
// чтобы сохранить различие между этапами игрового цикла.
func (m *TimeModel) stateMeanTime(
	state domain.State,
	mu float64,
	algorithmFactorA float64,
	fraudFactorF float64,
) float64 {
	totalMean := 1.0 / (mu * algorithmFactorA * (1.0 + fraudTimeAlpha*fraudFactorF))

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
