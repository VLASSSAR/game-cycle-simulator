package optimizer

import (
	"errors"
	"math/rand"
	"time"

	"game-cycle-simulator/internal/domain"
	"game-cycle-simulator/internal/simulator"
)

type Optimizer struct {
	sim *simulator.Simulator
	rng *rand.Rand
}

func NewOptimizer() *Optimizer {
	return &Optimizer{
		sim: simulator.NewSimulator(),
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

const defaultRandomSearchIterations = 100

func (o *Optimizer) Optimize(
	req domain.OptimizationRequest,
) (domain.OptimizationResult, error) {
	normalizedReq := normalizeOptimizationRequest(req)

	if err := validateOptimizationRequest(normalizedReq); err != nil {
		return domain.OptimizationResult{}, err
	}

	switch normalizedReq.Method {
	case domain.OptimizationMethodRandom:
		return o.optimizeRandomSearch(normalizedReq)
	default:
		return domain.OptimizationResult{}, errors.New("unsupported optimization method")
	}
}

func (o *Optimizer) optimizeRandomSearch(
	req domain.OptimizationRequest,
) (domain.OptimizationResult, error) {
	baselineReq := domain.SimulationRequest{
		ArrivalRateLambda: req.ArrivalRateLambda,
		Mu:                req.Mu,
		ChannelsK:         req.Baseline.ChannelsK,
		AlgorithmFactorA:  req.Baseline.AlgorithmFactorA,
		BetLimit:          req.Baseline.BetLimit,
		FraudFactorF:      req.Baseline.FraudFactorF,
		Simulations:       req.Simulations,
	}

	baselineResult, err := o.sim.Run(baselineReq)
	if err != nil {
		return domain.OptimizationResult{}, err
	}

	bestMetrics := baselineResult.Metrics
	bestParams := domain.BestParams{
		ChannelsK:        baselineReq.ChannelsK,
		AlgorithmFactorA: baselineReq.AlgorithmFactorA,
		BetLimit:         baselineReq.BetLimit,
		FraudFactorF:     baselineReq.FraudFactorF,
	}

	foundFeasible := false

	for i := 0; i < req.Iterations; i++ {
		k := randomIntCandidate(o.rng, req.ChannelCandidates)
		a := randomFloatCandidate(o.rng, req.AlgorithmCandidates)
		l := randomFloatCandidate(o.rng, req.BetLimitCandidates)
		f := randomFloatCandidate(o.rng, req.FraudCandidates)

		simReq := domain.SimulationRequest{
			ArrivalRateLambda: req.ArrivalRateLambda,
			Mu:                req.Mu,
			ChannelsK:         k,
			AlgorithmFactorA:  a,
			BetLimit:          l,
			FraudFactorF:      f,
			Simulations:       req.Simulations,
		}

		simResult, err := o.sim.Run(simReq)
		if err != nil {
			return domain.OptimizationResult{}, err
		}

		metrics := simResult.Metrics

		if !isFeasible(metrics, req.MaxRho, req.MaxProcessingTime) {
			continue
		}

		if !foundFeasible || metrics.RevenuePerUnitTime > bestMetrics.RevenuePerUnitTime {
			foundFeasible = true
			bestMetrics = metrics
			bestParams = domain.BestParams{
				ChannelsK:        k,
				AlgorithmFactorA: a,
				BetLimit:         l,
				FraudFactorF:     f,
			}
		}
	}

	if !foundFeasible {
		return domain.OptimizationResult{}, errors.New("no feasible parameter set found")
	}

	return domain.OptimizationResult{
		Request:    req,
		Method:     req.Method,
		Iterations: req.Iterations,

		BestParams:       bestParams,
		BaselineMetrics:  baselineResult.Metrics,
		OptimizedMetrics: bestMetrics,

		DeltaRevenue:    bestMetrics.RevenuePerUnitTime - baselineResult.Metrics.RevenuePerUnitTime,
		DeltaProcessing: baselineResult.Metrics.AvgProcessingTime - bestMetrics.AvgProcessingTime,
	}, nil
}

func normalizeOptimizationRequest(req domain.OptimizationRequest) domain.OptimizationRequest {
	if req.Method == "" {
		req.Method = domain.DefaultOptimizationMethod()
	}

	if req.Iterations <= 0 {
		req.Iterations = defaultRandomSearchIterations
	}

	return req
}

func isFeasible(
	metrics domain.Metrics,
	maxRho float64,
	maxProcessingTime float64,
) bool {
	if metrics.UtilizationRho >= maxRho {
		return false
	}

	if metrics.AvgProcessingTime > maxProcessingTime {
		return false
	}

	return true
}

func validateOptimizationRequest(req domain.OptimizationRequest) error {
	if !req.Method.IsValid() {
		return errors.New("invalid optimization method")
	}

	if req.Method != domain.OptimizationMethodRandom {
		return errors.New("only random optimization method is currently implemented")
	}

	if req.Iterations <= 0 {
		return errors.New("iterations must be > 0")
	}

	if req.ArrivalRateLambda <= 0 {
		return errors.New("arrival_rate_lambda must be > 0")
	}
	if req.Mu <= 0 {
		return errors.New("mu must be > 0")
	}
	if req.Simulations <= 0 {
		return errors.New("simulations must be > 0")
	}

	if req.Baseline.ChannelsK <= 0 {
		return errors.New("baseline.channels_k must be > 0")
	}
	if req.Baseline.AlgorithmFactorA <= 0 {
		return errors.New("baseline.algorithm_factor_a must be > 0")
	}
	if req.Baseline.BetLimit <= 0 {
		return errors.New("baseline.bet_limit must be > 0")
	}
	if req.Baseline.FraudFactorF < 0 || req.Baseline.FraudFactorF > 1 {
		return errors.New("baseline.fraud_factor must be in [0,1]")
	}

	if len(req.ChannelCandidates) == 0 {
		return errors.New("channel_candidates must not be empty")
	}
	if len(req.AlgorithmCandidates) == 0 {
		return errors.New("algorithm_candidates must not be empty")
	}
	if len(req.BetLimitCandidates) == 0 {
		return errors.New("bet_limit_candidates must not be empty")
	}
	if len(req.FraudCandidates) == 0 {
		return errors.New("fraud_candidates must not be empty")
	}
	if req.MaxRho <= 0 {
		return errors.New("max_rho must be > 0")
	}
	if req.MaxProcessingTime <= 0 {
		return errors.New("max_processing_time must be > 0")
	}

	for _, k := range req.ChannelCandidates {
		if k <= 0 {
			return errors.New("all channel_candidates must be > 0")
		}
	}

	for _, a := range req.AlgorithmCandidates {
		if a <= 0 {
			return errors.New("all algorithm_candidates must be > 0")
		}
	}

	for _, l := range req.BetLimitCandidates {
		if l <= 0 {
			return errors.New("all bet_limit_candidates must be > 0")
		}
	}

	for _, f := range req.FraudCandidates {
		if f < 0 || f > 1 {
			return errors.New("all fraud_candidates must be in [0,1]")
		}
	}

	return nil
}

func randomIntCandidate(rng *rand.Rand, candidates []int) int {
	return candidates[rng.Intn(len(candidates))]
}

func randomFloatCandidate(rng *rand.Rand, candidates []float64) float64 {
	return candidates[rng.Intn(len(candidates))]
}
