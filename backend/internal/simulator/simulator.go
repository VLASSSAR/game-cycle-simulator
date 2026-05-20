package simulator

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"game-cycle-simulator/internal/domain"
	antifraud "game-cycle-simulator/internal/fraud"
	"game-cycle-simulator/internal/generator"
	"game-cycle-simulator/internal/model"
)

type Simulator struct {
	arrivalGenerator      *generator.ArrivalGenerator
	transitionModel       *model.TransitionModel
	timeModel             *model.TimeModel
	revenueModel          *model.RevenueModel
	fraudFeatureGenerator *antifraud.FeatureGenerator
}

func NewSimulator() *Simulator {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &Simulator{
		arrivalGenerator:      generator.NewArrivalGenerator(rng),
		transitionModel:       model.NewTransitionModel(rng),
		timeModel:             model.NewTimeModel(rng),
		revenueModel:          model.NewRevenueModel(rng),
		fraudFeatureGenerator: antifraud.NewFeatureGenerator(rng),
	}
}

func (s *Simulator) Run(req domain.SimulationRequest) (domain.SimulationResult, error) {
	if err := validateRequest(req); err != nil {
		return domain.SimulationResult{}, err
	}

	arrivals := s.arrivalGenerator.GenerateArrivals(req.ArrivalRateLambda, req.Simulations)

	// channelFreeTimes[i] = момент, когда i-й канал освободится.
	channelFreeTimes := make([]float64, req.ChannelsK)

	fraudModel := antifraud.NewModel(buildFraudConfig(req.FraudFactorF))

	var totalProcessingTime float64
	var totalQueueTime float64
	var totalSystemTime float64
	var totalRevenue float64
	var totalFraudScore float64

	var successfulBets int
	var failedBets int

	var fraudChecks int
	var fraudDetected int
	var actualFraudOperations int
	var falsePositives int
	var falseNegatives int
	var truePositives int

	// Текущая оценка rho на основе параметров системы.
	// Используется в transition model, где вероятность неуспешного перехода зависит от загрузки.
	baseRho := estimateRho(req.ArrivalRateLambda, req.Mu, req.ChannelsK, req.AlgorithmFactorA)

	for _, arrivalTime := range arrivals {
		channelIdx := findEarliestFreeChannel(channelFreeTimes)

		startServiceTime := math.Max(arrivalTime, channelFreeTimes[channelIdx])
		queueTime := startServiceTime - arrivalTime

		operationProfile := s.fraudFeatureGenerator.Generate(baseFraudProbability())
		fraudResult := fraudModel.Evaluate(operationProfile.Features, operationProfile.IsFraud)

		fraudChecks++
		totalFraudScore += fraudResult.Score

		if fraudResult.IsSuspicious {
			fraudDetected++
		}

		if fraudResult.IsFraud {
			actualFraudOperations++
		}

		if fraudResult.IsSuspicious && fraudResult.IsFraud {
			truePositives++
		}

		if fraudResult.IsFalsePositive {
			falsePositives++
		}

		if fraudResult.IsFalseNegative {
			falseNegatives++
		}

		operationFraudFactor := calculateOperationFraudFactor(
			req.FraudFactorF,
			fraudResult.Score,
			fraudResult.IsSuspicious,
		)

		path := s.transitionModel.SimulatePath(
			baseRho,
			req.AlgorithmFactorA,
			operationFraudFactor,
		)

		processingTime := s.timeModel.TotalProcessingTime(
			path,
			req.Mu,
			req.AlgorithmFactorA,
			operationFraudFactor,
		)

		finishTime := startServiceTime + processingTime
		systemTime := finishTime - arrivalTime

		channelFreeTimes[channelIdx] = finishTime

		revenue := s.revenueModel.RevenueForPath(
			path,
			req.BetLimit,
		)

		if len(path) > 0 && path[len(path)-1] == domain.Wallet {
			successfulBets++
		} else {
			failedBets++
		}

		totalProcessingTime += processingTime
		totalQueueTime += queueTime
		totalSystemTime += systemTime
		totalRevenue += revenue
	}

	totalBets := req.Simulations
	if totalBets == 0 {
		return domain.SimulationResult{}, errors.New("no bets simulated")
	}

	avgProcessingTime := totalProcessingTime / float64(totalBets)
	avgQueueTime := totalQueueTime / float64(totalBets)
	avgSystemTime := totalSystemTime / float64(totalBets)

	successProbability := float64(successfulBets) / float64(totalBets)
	failureProbability := float64(failedBets) / float64(totalBets)

	var serviceRateMu float64
	if avgProcessingTime > 0 {
		serviceRateMu = 1.0 / avgProcessingTime
	}

	var utilizationRho float64
	if serviceRateMu > 0 {
		utilizationRho = req.ArrivalRateLambda / (float64(req.ChannelsK) * serviceRateMu)
	}

	meanRevenuePerBet := totalRevenue / float64(totalBets)

	simulationHorizon := maxFloatSlice(channelFreeTimes)
	var revenuePerUnitTime float64
	if simulationHorizon > 0 {
		revenuePerUnitTime = totalRevenue / simulationHorizon
	}

	var avgFraudScore float64
	if fraudChecks > 0 {
		avgFraudScore = totalFraudScore / float64(fraudChecks)
	}

	var fraudDetectionRate float64
	if actualFraudOperations > 0 {
		fraudDetectionRate = float64(truePositives) / float64(actualFraudOperations)
	}

	normalOperations := totalBets - actualFraudOperations

	var falsePositiveRate float64
	if normalOperations > 0 {
		falsePositiveRate = float64(falsePositives) / float64(normalOperations)
	}

	var falseNegativeRate float64
	if actualFraudOperations > 0 {
		falseNegativeRate = float64(falseNegatives) / float64(actualFraudOperations)
	}

	result := domain.SimulationResult{
		Request: req,
		Metrics: domain.Metrics{
			TotalBets:          totalBets,
			SuccessfulBets:     successfulBets,
			FailedBets:         failedBets,
			SuccessProbability: successProbability,
			FailureProbability: failureProbability,

			AvgProcessingTime: avgProcessingTime,
			AvgQueueTime:      avgQueueTime,
			AvgSystemTime:     avgSystemTime,

			ServiceRateMu:  serviceRateMu,
			UtilizationRho: utilizationRho,

			MeanRevenuePerBet:  meanRevenuePerBet,
			RevenuePerUnitTime: revenuePerUnitTime,

			FraudChecks:           fraudChecks,
			FraudDetected:         fraudDetected,
			ActualFraudOperations: actualFraudOperations,
			FalsePositives:        falsePositives,
			FalseNegatives:        falseNegatives,
			AvgFraudScore:         avgFraudScore,
			FraudDetectionRate:    fraudDetectionRate,
			FalsePositiveRate:     falsePositiveRate,
			FalseNegativeRate:     falseNegativeRate,
		},
	}

	return result, nil
}

func validateRequest(req domain.SimulationRequest) error {
	if req.ArrivalRateLambda <= 0 {
		return errors.New("arrival_rate_lambda must be > 0")
	}
	if req.Mu <= 0 {
		return errors.New("mu must be > 0")
	}
	if req.ChannelsK <= 0 {
		return errors.New("channels_k must be > 0")
	}
	if req.AlgorithmFactorA <= 0 {
		return errors.New("algorithm_factor_a must be > 0")
	}
	if req.BetLimit <= 0 {
		return errors.New("bet_limit must be > 0")
	}
	if req.FraudFactorF < 0 || req.FraudFactorF > 1 {
		return errors.New("fraud_factor must be in [0,1]")
	}
	if req.Simulations <= 0 {
		return errors.New("simulations must be > 0")
	}
	return nil
}

func findEarliestFreeChannel(channelFreeTimes []float64) int {
	minIdx := 0
	for i := 1; i < len(channelFreeTimes); i++ {
		if channelFreeTimes[i] < channelFreeTimes[minIdx] {
			minIdx = i
		}
	}
	return minIdx
}

func maxFloatSlice(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	maxVal := values[0]
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

// estimateRho даёт стартовую оценку коэффициента загрузки,
// используемую при моделировании переходов.
func estimateRho(
	lambda float64,
	mu float64,
	channelsK int,
	algorithmFactorA float64,
) float64 {
	effectiveMu := mu * algorithmFactorA
	if effectiveMu <= 0 || channelsK <= 0 {
		return 1.0
	}

	rho := lambda / (float64(channelsK) * effectiveMu)

	// Ограничим rho сверху для устойчивости численных расчётов.
	if rho < 0 {
		return 0
	}
	if rho > 1.5 {
		return 1.5
	}

	return rho
}

func buildFraudConfig(fraudFactorF float64) domain.FraudConfig {
	config := antifraud.DefaultConfig()

	// Чем выше F, тем строже антифрод-фильтрация:
	// - fraud-score усиливается через Strictness;
	// - порог подозрительности немного снижается.
	config.Strictness = 1.0 + 0.5*fraudFactorF
	config.Threshold = clampFloat(0.65-0.20*fraudFactorF, 0.35, 0.90)

	return config
}

func baseFraudProbability() float64 {
	// В рамках имитационной модели считаем, что небольшая доля операций
	// является мошеннической независимо от строгости антифрод-фильтрации.
	return 0.08
}

func calculateOperationFraudFactor(
	baseFraudFactor float64,
	fraudScore float64,
	isSuspicious bool,
) float64 {
	factor := 0.6*baseFraudFactor + 0.4*fraudScore

	if isSuspicious {
		factor += 0.1
	}

	return clampFloat(factor, 0.0, 1.0)
}

func clampFloat(value float64, minValue float64, maxValue float64) float64 {
	return math.Max(minValue, math.Min(maxValue, value))
}
