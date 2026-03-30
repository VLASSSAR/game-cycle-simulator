package simulator

import (
	"errors"
	"math"
	"math/rand"
	"time"

	"game-cycle-simulator/internal/domain"
	"game-cycle-simulator/internal/generator"
	"game-cycle-simulator/internal/model"
)

type Simulator struct {
	arrivalGenerator *generator.ArrivalGenerator
	transitionModel  *model.TransitionModel
	timeModel        *model.TimeModel
	revenueModel     *model.RevenueModel
}

func NewSimulator() *Simulator {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &Simulator{
		arrivalGenerator: generator.NewArrivalGenerator(rng),
		transitionModel:  model.NewTransitionModel(rng),
		timeModel:        model.NewTimeModel(rng),
		revenueModel:     model.NewRevenueModel(rng),
	}
}

func (s *Simulator) Run(req domain.SimulationRequest) (domain.SimulationResult, error) {
	if err := validateRequest(req); err != nil {
		return domain.SimulationResult{}, err
	}

	arrivals := s.arrivalGenerator.GenerateArrivals(req.ArrivalRateLambda, req.Simulations)

	// channelFreeTimes[i] = момент, когда i-й канал освободится
	channelFreeTimes := make([]float64, req.ChannelsK)

	var totalProcessingTime float64
	var totalSystemTime float64
	var totalRevenue float64
	var successfulBets int
	var failedBets int

	// Текущая оценка rho на основе параметров системы.
	// Используется в transition model, где вероятность неуспешного перехода зависит от загрузки.
	baseRho := estimateRho(req.ArrivalRateLambda, req.Mu, req.ChannelsK, req.AlgorithmFactorA)

	for _, arrivalTime := range arrivals {
		channelIdx := findEarliestFreeChannel(channelFreeTimes)

		startServiceTime := math.Max(arrivalTime, channelFreeTimes[channelIdx])

		path := s.transitionModel.SimulatePath(
			baseRho,
			req.AlgorithmFactorA,
			req.FraudFactorF,
		)

		processingTime := s.timeModel.TotalProcessingTime(
			path,
			req.Mu,
			req.AlgorithmFactorA,
			req.FraudFactorF,
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
		totalSystemTime += systemTime
		totalRevenue += revenue
	}

	totalBets := req.Simulations
	if totalBets == 0 {
		return domain.SimulationResult{}, errors.New("no bets simulated")
	}

	avgProcessingTime := totalProcessingTime / float64(totalBets)
	avgSystemTime := totalSystemTime / float64(totalBets)
	successProbability := float64(successfulBets) / float64(totalBets)

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

	result := domain.SimulationResult{
		Request: req,
		Metrics: domain.Metrics{
			TotalBets:          totalBets,
			SuccessfulBets:     successfulBets,
			FailedBets:         failedBets,
			SuccessProbability: successProbability,
			AvgProcessingTime:  avgProcessingTime,
			AvgSystemTime:      avgSystemTime,
			ServiceRateMu:      serviceRateMu,
			UtilizationRho:     utilizationRho,
			MeanRevenuePerBet:  meanRevenuePerBet,
			RevenuePerUnitTime: revenuePerUnitTime,
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

	// ограничим rho сверху для устойчивости численных расчётов
	if rho < 0 {
		return 0
	}
	if rho > 1.5 {
		return 1.5
	}

	return rho
}
