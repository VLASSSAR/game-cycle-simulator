package model

import (
	"math/rand"

	"game-cycle-simulator/internal/domain"
)

type RevenueModel struct {
	rng *rand.Rand
}

func NewRevenueModel(rng *rand.Rand) *RevenueModel {
	return &RevenueModel{
		rng: rng,
	}
}

const (
	marginGamma  = 0.08
	errorPenalty = 20.0
)

// RevenueForPath возвращает экономический результат обработки одной ставки.
//
// Логика:
// - если ставка завершилась ошибкой, начисляется штраф;
// - если завершилась успешно, доход зависит от лимита ставки.
//
// Средний размер ставки моделируется как случайная величина
// на интервале [0, L], где L — лимит ставки.
func (m *RevenueModel) RevenueForPath(
	path []domain.State,
	betLimit float64,
) float64 {
	if len(path) == 0 {
		return 0
	}

	lastState := path[len(path)-1]

	if lastState == domain.ErrorState {
		return -errorPenalty
	}

	betSize := m.randomBetSize(betLimit)

	revenue := marginGamma * betSize

	// Добавляем небольшой шум, чтобы модель не была слишком жёсткой.
	noise := m.rng.NormFloat64() * (0.03 * revenue)
	revenue += noise

	if revenue < 0 {
		return 0
	}

	return revenue
}

// ExpectedBetSize возвращает математическое ожидание размера ставки
// при равномерном распределении на [0, L].
func (m *RevenueModel) ExpectedBetSize(betLimit float64) float64 {
	return betLimit / 2.0
}

func (m *RevenueModel) randomBetSize(betLimit float64) float64 {
	return m.rng.Float64() * betLimit
}
