package model

import (
	"math"
	"math/rand"

	"game-cycle-simulator/internal/domain"
)

type TransitionModel struct {
	rng *rand.Rand
}

func NewTransitionModel(rng *rand.Rand) *TransitionModel {
	return &TransitionModel{
		rng: rng,
	}
}

// NextState определяет следующее состояние игрового цикла
// в зависимости от текущего состояния, загрузки rho,
// алгоритмического коэффициента a и параметра антифрода F.
func (m *TransitionModel) NextState(
	current domain.State,
	rho float64,
	algorithmFactorA float64,
	fraudFactorF float64,
) domain.State {
	switch current {
	case domain.Init:
		return domain.Reserve

	case domain.Reserve:
		if m.isFailure(m.reserveFailProb(rho, algorithmFactorA, fraudFactorF)) {
			return domain.ErrorState
		}
		return domain.Engine

	case domain.Engine:
		if m.isFailure(m.engineFailProb(rho, algorithmFactorA, fraudFactorF)) {
			return domain.ErrorState
		}
		return domain.Settlement

	case domain.Settlement:
		if m.isFailure(m.settlementFailProb(rho, algorithmFactorA, fraudFactorF)) {
			return domain.ErrorState
		}
		return domain.Wallet

	default:
		return current
	}
}

// SimulatePath моделирует траекторию одной ставки
// от состояния Init до терминального состояния.
func (m *TransitionModel) SimulatePath(
	rho float64,
	algorithmFactorA float64,
	fraudFactorF float64,
) []domain.State {
	path := []domain.State{domain.Init}
	current := domain.Init

	for !domain.IsTerminalState(current) {
		current = m.NextState(current, rho, algorithmFactorA, fraudFactorF)
		path = append(path, current)
	}

	return path
}

func (m *TransitionModel) isFailure(prob float64) bool {
	return m.rng.Float64() < prob
}

// Ниже задаются вероятности неуспешного перехода.
// Они состоят из трёх факторов:
// 1. базовая вероятность ошибки на этапе,
// 2. рост ошибки при увеличении rho,
// 3. снижение ошибки при росте algorithmFactorA,
// 4. рост отклонений при увеличении fraudFactorF.

func (m *TransitionModel) reserveFailProb(
	rho float64,
	a float64,
	f float64,
) float64 {
	base := 0.01
	alpha := 0.05 // чувствительность к загрузке
	beta := 0.004 // улучшение за счёт алгоритмов
	gamma := 0.03 // строгость антифрода

	return clamp(base+alpha*rho-beta*(a-1.0)+gamma*f, 0.0, 0.9)
}

func (m *TransitionModel) engineFailProb(
	rho float64,
	a float64,
	f float64,
) float64 {
	base := 0.015
	alpha := 0.06
	beta := 0.005
	gamma := 0.01 // антифрод влияет слабее на engine

	return clamp(base+alpha*rho-beta*(a-1.0)+gamma*f, 0.0, 0.9)
}

func (m *TransitionModel) settlementFailProb(
	rho float64,
	a float64,
	f float64,
) float64 {
	base := 0.012
	alpha := 0.05
	beta := 0.004
	gamma := 0.04 // settlement чувствителен к дополнительным проверкам

	return clamp(base+alpha*rho-beta*(a-1.0)+gamma*f, 0.0, 0.9)
}

func clamp(x, minVal, maxVal float64) float64 {
	return math.Max(minVal, math.Min(maxVal, x))
}
