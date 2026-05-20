package service

import "game-cycle-simulator/internal/domain"

type ScenarioService struct{}

func NewScenarioService() *ScenarioService {
	return &ScenarioService{}
}

func (s *ScenarioService) ListScenarios() []domain.LoadScenario {
	return []domain.LoadScenario{
		{
			ID:          "low_load",
			Name:        "Низкая нагрузка",
			Description: "Сценарий спокойной работы системы с небольшим входящим потоком заявок.",

			ArrivalRateLambda: 4,
			Mu:                5,
			ChannelsK:         3,
			AlgorithmFactorA:  1.0,
			BetLimit:          100,
			FraudFactorF:      0.2,
			Simulations:       10000,

			MaxRho:            0.95,
			MaxProcessingTime: 0.3,
		},
		{
			ID:          "medium_load",
			Name:        "Средняя нагрузка",
			Description: "Базовый сценарий функционирования платформы при умеренной интенсивности пользовательских операций.",

			ArrivalRateLambda: 8,
			Mu:                5,
			ChannelsK:         3,
			AlgorithmFactorA:  1.0,
			BetLimit:          100,
			FraudFactorF:      0.2,
			Simulations:       10000,

			MaxRho:            0.95,
			MaxProcessingTime: 0.3,
		},
		{
			ID:          "high_load",
			Name:        "Высокая нагрузка",
			Description: "Сценарий повышенной нагрузки, при котором возрастает время ожидания и коэффициент загрузки системы.",

			ArrivalRateLambda: 12,
			Mu:                5,
			ChannelsK:         3,
			AlgorithmFactorA:  1.0,
			BetLimit:          100,
			FraudFactorF:      0.25,
			Simulations:       10000,

			MaxRho:            0.95,
			MaxProcessingTime: 0.35,
		},
		{
			ID:          "stress_load",
			Name:        "Стресс-нагрузка",
			Description: "Сценарий работы системы вблизи предельной загрузки для анализа устойчивости и масштабируемости.",

			ArrivalRateLambda: 16,
			Mu:                5,
			ChannelsK:         3,
			AlgorithmFactorA:  1.0,
			BetLimit:          100,
			FraudFactorF:      0.3,
			Simulations:       10000,

			MaxRho:            0.98,
			MaxProcessingTime: 0.5,
		},
	}
}
