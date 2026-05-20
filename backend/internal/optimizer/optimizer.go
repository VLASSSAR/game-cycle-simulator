package optimizer

import (
	"errors"
	"math/rand"
	"sort"
	"time"

	"game-cycle-simulator/internal/domain"
	"game-cycle-simulator/internal/simulator"
)

type Optimizer struct {
	sim *simulator.Simulator
	rng *rand.Rand
}

type candidate struct {
	ChannelsK        int
	AlgorithmFactorA float64
	BetLimit         float64
	FraudFactorF     float64
}

type evaluatedCandidate struct {
	Candidate candidate
	Metrics   domain.Metrics
	Fitness   float64
	Feasible  bool
}

func NewOptimizer() *Optimizer {
	return &Optimizer{
		sim: simulator.NewSimulator(),
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

const (
	defaultRandomSearchIterations = 100

	minGeneticPopulationSize   = 4
	maxGeneticPopulationSize   = 20
	geneticMutationProbability = 0.2

	adaptiveExplorationProbability = 0.25
)

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

	case domain.OptimizationMethodGrid:
		return o.optimizeGridSearch(normalizedReq)

	case domain.OptimizationMethodGenetic:
		return o.optimizeGenetic(normalizedReq)

	case domain.OptimizationMethodAdaptive:
		return o.optimizeAdaptive(normalizedReq)

	default:
		return domain.OptimizationResult{}, errors.New("optimization method is not implemented yet")
	}
}

func (o *Optimizer) optimizeRandomSearch(
	req domain.OptimizationRequest,
) (domain.OptimizationResult, error) {
	baselineResult, err := o.runBaseline(req)
	if err != nil {
		return domain.OptimizationResult{}, err
	}

	bestMetrics := baselineResult.Metrics
	bestParams := domain.BestParams{
		ChannelsK:        req.Baseline.ChannelsK,
		AlgorithmFactorA: req.Baseline.AlgorithmFactorA,
		BetLimit:         req.Baseline.BetLimit,
		FraudFactorF:     req.Baseline.FraudFactorF,
	}

	foundFeasible := isFeasible(baselineResult.Metrics, req.MaxRho, req.MaxProcessingTime)

	for i := 0; i < req.Iterations; i++ {
		c := o.randomCandidate(req)

		simResult, err := o.runCandidate(req, c)
		if err != nil {
			return domain.OptimizationResult{}, err
		}

		metrics := simResult.Metrics

		if !isFeasible(metrics, req.MaxRho, req.MaxProcessingTime) {
			continue
		}

		if !foundFeasible || isBetter(metrics, bestMetrics) {
			foundFeasible = true
			bestMetrics = metrics
			bestParams = bestParamsFromCandidate(c)
		}
	}

	if !foundFeasible {
		return domain.OptimizationResult{}, errors.New("no feasible parameter set found")
	}

	return buildOptimizationResult(
		req,
		req.Method,
		req.Iterations,
		bestParams,
		baselineResult.Metrics,
		bestMetrics,
	), nil
}

func (o *Optimizer) optimizeGridSearch(
	req domain.OptimizationRequest,
) (domain.OptimizationResult, error) {
	baselineResult, err := o.runBaseline(req)
	if err != nil {
		return domain.OptimizationResult{}, err
	}

	bestMetrics := baselineResult.Metrics
	bestParams := domain.BestParams{
		ChannelsK:        req.Baseline.ChannelsK,
		AlgorithmFactorA: req.Baseline.AlgorithmFactorA,
		BetLimit:         req.Baseline.BetLimit,
		FraudFactorF:     req.Baseline.FraudFactorF,
	}

	foundFeasible := isFeasible(baselineResult.Metrics, req.MaxRho, req.MaxProcessingTime)
	evaluatedCombinations := 0

	for _, k := range req.ChannelCandidates {
		for _, a := range req.AlgorithmCandidates {
			for _, l := range req.BetLimitCandidates {
				for _, f := range req.FraudCandidates {
					evaluatedCombinations++

					c := candidate{
						ChannelsK:        k,
						AlgorithmFactorA: a,
						BetLimit:         l,
						FraudFactorF:     f,
					}

					simResult, err := o.runCandidate(req, c)
					if err != nil {
						return domain.OptimizationResult{}, err
					}

					metrics := simResult.Metrics

					if !isFeasible(metrics, req.MaxRho, req.MaxProcessingTime) {
						continue
					}

					if !foundFeasible || isBetter(metrics, bestMetrics) {
						foundFeasible = true
						bestMetrics = metrics
						bestParams = bestParamsFromCandidate(c)
					}
				}
			}
		}
	}

	if !foundFeasible {
		return domain.OptimizationResult{}, errors.New("no feasible parameter set found")
	}

	return buildOptimizationResult(
		req,
		req.Method,
		evaluatedCombinations,
		bestParams,
		baselineResult.Metrics,
		bestMetrics,
	), nil
}

func (o *Optimizer) optimizeGenetic(
	req domain.OptimizationRequest,
) (domain.OptimizationResult, error) {
	baselineResult, err := o.runBaseline(req)
	if err != nil {
		return domain.OptimizationResult{}, err
	}

	bestMetrics := baselineResult.Metrics
	bestParams := domain.BestParams{
		ChannelsK:        req.Baseline.ChannelsK,
		AlgorithmFactorA: req.Baseline.AlgorithmFactorA,
		BetLimit:         req.Baseline.BetLimit,
		FraudFactorF:     req.Baseline.FraudFactorF,
	}

	foundFeasible := isFeasible(baselineResult.Metrics, req.MaxRho, req.MaxProcessingTime)

	populationSize := geneticPopulationSize(req.Iterations)
	population := make([]candidate, 0, populationSize)

	for i := 0; i < populationSize; i++ {
		population = append(population, o.randomCandidate(req))
	}

	evaluatedTotal := 0

	for evaluatedTotal < req.Iterations {
		evaluatedPopulation := make([]evaluatedCandidate, 0, len(population))

		for _, c := range population {
			if evaluatedTotal >= req.Iterations {
				break
			}

			simResult, err := o.runCandidate(req, c)
			if err != nil {
				return domain.OptimizationResult{}, err
			}

			evaluatedTotal++

			metrics := simResult.Metrics
			feasible := isFeasible(metrics, req.MaxRho, req.MaxProcessingTime)
			fitness := fitnessValue(metrics, feasible)

			evaluated := evaluatedCandidate{
				Candidate: c,
				Metrics:   metrics,
				Fitness:   fitness,
				Feasible:  feasible,
			}

			evaluatedPopulation = append(evaluatedPopulation, evaluated)

			if feasible && (!foundFeasible || isBetter(metrics, bestMetrics)) {
				foundFeasible = true
				bestMetrics = metrics
				bestParams = bestParamsFromCandidate(c)
			}
		}

		if len(evaluatedPopulation) == 0 {
			break
		}

		sort.Slice(evaluatedPopulation, func(i, j int) bool {
			return evaluatedPopulation[i].Fitness > evaluatedPopulation[j].Fitness
		})

		population = o.nextGeneration(req, evaluatedPopulation, populationSize)
	}

	if !foundFeasible {
		return domain.OptimizationResult{}, errors.New("no feasible parameter set found")
	}

	return buildOptimizationResult(
		req,
		req.Method,
		evaluatedTotal,
		bestParams,
		baselineResult.Metrics,
		bestMetrics,
	), nil
}

func (o *Optimizer) optimizeAdaptive(
	req domain.OptimizationRequest,
) (domain.OptimizationResult, error) {
	baselineResult, err := o.runBaseline(req)
	if err != nil {
		return domain.OptimizationResult{}, err
	}

	currentCandidate := candidate{
		ChannelsK:        req.Baseline.ChannelsK,
		AlgorithmFactorA: req.Baseline.AlgorithmFactorA,
		BetLimit:         req.Baseline.BetLimit,
		FraudFactorF:     req.Baseline.FraudFactorF,
	}

	bestCandidate := currentCandidate
	bestMetrics := baselineResult.Metrics
	foundFeasible := isFeasible(baselineResult.Metrics, req.MaxRho, req.MaxProcessingTime)

	evaluatedTotal := 0

	for evaluatedTotal < req.Iterations {
		nextCandidate := o.proposeAdaptiveCandidate(req, currentCandidate)

		simResult, err := o.runCandidate(req, nextCandidate)
		if err != nil {
			return domain.OptimizationResult{}, err
		}

		evaluatedTotal++

		metrics := simResult.Metrics
		feasible := isFeasible(metrics, req.MaxRho, req.MaxProcessingTime)

		if feasible && (!foundFeasible || isBetter(metrics, bestMetrics)) {
			foundFeasible = true
			bestMetrics = metrics
			bestCandidate = nextCandidate
			currentCandidate = nextCandidate
			continue
		}

		// Если решение допустимое, но не лучшее, иногда переходим к нему,
		// чтобы алгоритм мог исследовать пространство параметров шире.
		if feasible && o.rng.Float64() < adaptiveExplorationProbability {
			currentCandidate = nextCandidate
		}
	}

	if !foundFeasible {
		return domain.OptimizationResult{}, errors.New("no feasible parameter set found")
	}

	return buildOptimizationResult(
		req,
		req.Method,
		evaluatedTotal,
		bestParamsFromCandidate(bestCandidate),
		baselineResult.Metrics,
		bestMetrics,
	), nil
}

func (o *Optimizer) nextGeneration(
	req domain.OptimizationRequest,
	evaluatedPopulation []evaluatedCandidate,
	populationSize int,
) []candidate {
	next := make([]candidate, 0, populationSize)

	eliteCount := minInt(2, len(evaluatedPopulation))
	for i := 0; i < eliteCount; i++ {
		next = append(next, evaluatedPopulation[i].Candidate)
	}

	for len(next) < populationSize {
		parentA := o.selectParent(evaluatedPopulation)
		parentB := o.selectParent(evaluatedPopulation)

		child := o.crossover(parentA.Candidate, parentB.Candidate)
		child = o.mutate(req, child)

		next = append(next, child)
	}

	return next
}

func (o *Optimizer) selectParent(
	evaluatedPopulation []evaluatedCandidate,
) evaluatedCandidate {
	best := evaluatedPopulation[o.rng.Intn(len(evaluatedPopulation))]

	tournamentSize := minInt(3, len(evaluatedPopulation))

	for i := 1; i < tournamentSize; i++ {
		candidateIndex := o.rng.Intn(len(evaluatedPopulation))
		current := evaluatedPopulation[candidateIndex]

		if current.Fitness > best.Fitness {
			best = current
		}
	}

	return best
}

func (o *Optimizer) crossover(
	parentA candidate,
	parentB candidate,
) candidate {
	child := candidate{}

	if o.rng.Float64() < 0.5 {
		child.ChannelsK = parentA.ChannelsK
	} else {
		child.ChannelsK = parentB.ChannelsK
	}

	if o.rng.Float64() < 0.5 {
		child.AlgorithmFactorA = parentA.AlgorithmFactorA
	} else {
		child.AlgorithmFactorA = parentB.AlgorithmFactorA
	}

	if o.rng.Float64() < 0.5 {
		child.BetLimit = parentA.BetLimit
	} else {
		child.BetLimit = parentB.BetLimit
	}

	if o.rng.Float64() < 0.5 {
		child.FraudFactorF = parentA.FraudFactorF
	} else {
		child.FraudFactorF = parentB.FraudFactorF
	}

	return child
}

func (o *Optimizer) mutate(
	req domain.OptimizationRequest,
	c candidate,
) candidate {
	if o.rng.Float64() < geneticMutationProbability {
		c.ChannelsK = randomIntCandidate(o.rng, req.ChannelCandidates)
	}

	if o.rng.Float64() < geneticMutationProbability {
		c.AlgorithmFactorA = randomFloatCandidate(o.rng, req.AlgorithmCandidates)
	}

	if o.rng.Float64() < geneticMutationProbability {
		c.BetLimit = randomFloatCandidate(o.rng, req.BetLimitCandidates)
	}

	if o.rng.Float64() < geneticMutationProbability {
		c.FraudFactorF = randomFloatCandidate(o.rng, req.FraudCandidates)
	}

	return c
}

func (o *Optimizer) proposeAdaptiveCandidate(
	req domain.OptimizationRequest,
	current candidate,
) candidate {
	// Иногда выполняем полностью случайный шаг,
	// чтобы выйти из локального экстремума.
	if o.rng.Float64() < adaptiveExplorationProbability {
		return o.randomCandidate(req)
	}

	next := current

	// На каждом шаге изменяем один или несколько параметров
	// на соседние значения из допустимых candidate-массивов.
	if o.rng.Float64() < 0.5 {
		next.ChannelsK = o.neighborIntCandidate(req.ChannelCandidates, current.ChannelsK)
	}

	if o.rng.Float64() < 0.5 {
		next.AlgorithmFactorA = o.neighborFloatCandidate(req.AlgorithmCandidates, current.AlgorithmFactorA)
	}

	if o.rng.Float64() < 0.5 {
		next.BetLimit = o.neighborFloatCandidate(req.BetLimitCandidates, current.BetLimit)
	}

	if o.rng.Float64() < 0.5 {
		next.FraudFactorF = o.neighborFloatCandidate(req.FraudCandidates, current.FraudFactorF)
	}

	return next
}

func (o *Optimizer) runBaseline(
	req domain.OptimizationRequest,
) (domain.SimulationResult, error) {
	baselineReq := domain.SimulationRequest{
		ArrivalRateLambda: req.ArrivalRateLambda,
		Mu:                req.Mu,
		ChannelsK:         req.Baseline.ChannelsK,
		AlgorithmFactorA:  req.Baseline.AlgorithmFactorA,
		BetLimit:          req.Baseline.BetLimit,
		FraudFactorF:      req.Baseline.FraudFactorF,
		Simulations:       req.Simulations,
	}

	return o.sim.Run(baselineReq)
}

func (o *Optimizer) runCandidate(
	req domain.OptimizationRequest,
	c candidate,
) (domain.SimulationResult, error) {
	simReq := domain.SimulationRequest{
		ArrivalRateLambda: req.ArrivalRateLambda,
		Mu:                req.Mu,
		ChannelsK:         c.ChannelsK,
		AlgorithmFactorA:  c.AlgorithmFactorA,
		BetLimit:          c.BetLimit,
		FraudFactorF:      c.FraudFactorF,
		Simulations:       req.Simulations,
	}

	return o.sim.Run(simReq)
}

func (o *Optimizer) randomCandidate(
	req domain.OptimizationRequest,
) candidate {
	return candidate{
		ChannelsK:        randomIntCandidate(o.rng, req.ChannelCandidates),
		AlgorithmFactorA: randomFloatCandidate(o.rng, req.AlgorithmCandidates),
		BetLimit:         randomFloatCandidate(o.rng, req.BetLimitCandidates),
		FraudFactorF:     randomFloatCandidate(o.rng, req.FraudCandidates),
	}
}

func buildOptimizationResult(
	req domain.OptimizationRequest,
	method domain.OptimizationMethod,
	iterations int,
	bestParams domain.BestParams,
	baselineMetrics domain.Metrics,
	optimizedMetrics domain.Metrics,
) domain.OptimizationResult {
	return domain.OptimizationResult{
		Request:    req,
		Method:     method,
		Iterations: iterations,

		BestParams:       bestParams,
		BaselineMetrics:  baselineMetrics,
		OptimizedMetrics: optimizedMetrics,

		DeltaRevenue:    optimizedMetrics.RevenuePerUnitTime - baselineMetrics.RevenuePerUnitTime,
		DeltaProcessing: baselineMetrics.AvgProcessingTime - optimizedMetrics.AvgProcessingTime,
	}
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

func fitnessValue(
	metrics domain.Metrics,
	feasible bool,
) float64 {
	if !feasible {
		return -1e18
	}

	return metrics.RevenuePerUnitTime
}

func isBetter(
	candidate domain.Metrics,
	currentBest domain.Metrics,
) bool {
	return candidate.RevenuePerUnitTime > currentBest.RevenuePerUnitTime
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

func geneticPopulationSize(iterations int) int {
	if iterations < minGeneticPopulationSize {
		return minGeneticPopulationSize
	}

	proposed := iterations / 5

	if proposed < minGeneticPopulationSize {
		return minGeneticPopulationSize
	}

	if proposed > maxGeneticPopulationSize {
		return maxGeneticPopulationSize
	}

	return proposed
}

func bestParamsFromCandidate(c candidate) domain.BestParams {
	return domain.BestParams{
		ChannelsK:        c.ChannelsK,
		AlgorithmFactorA: c.AlgorithmFactorA,
		BetLimit:         c.BetLimit,
		FraudFactorF:     c.FraudFactorF,
	}
}

func (o *Optimizer) neighborIntCandidate(
	candidates []int,
	current int,
) int {
	if len(candidates) == 0 {
		return current
	}

	currentIndex := nearestIntIndex(candidates, current)

	direction := -1
	if o.rng.Float64() < 0.5 {
		direction = 1
	}

	nextIndex := currentIndex + direction

	if nextIndex < 0 {
		nextIndex = 0
	}

	if nextIndex >= len(candidates) {
		nextIndex = len(candidates) - 1
	}

	return candidates[nextIndex]
}

func (o *Optimizer) neighborFloatCandidate(
	candidates []float64,
	current float64,
) float64 {
	if len(candidates) == 0 {
		return current
	}

	currentIndex := nearestFloatIndex(candidates, current)

	direction := -1
	if o.rng.Float64() < 0.5 {
		direction = 1
	}

	nextIndex := currentIndex + direction

	if nextIndex < 0 {
		nextIndex = 0
	}

	if nextIndex >= len(candidates) {
		nextIndex = len(candidates) - 1
	}

	return candidates[nextIndex]
}

func nearestIntIndex(
	candidates []int,
	current int,
) int {
	bestIndex := 0
	bestDistance := absInt(candidates[0] - current)

	for i := 1; i < len(candidates); i++ {
		distance := absInt(candidates[i] - current)
		if distance < bestDistance {
			bestDistance = distance
			bestIndex = i
		}
	}

	return bestIndex
}

func nearestFloatIndex(
	candidates []float64,
	current float64,
) int {
	bestIndex := 0
	bestDistance := absFloat(candidates[0] - current)

	for i := 1; i < len(candidates); i++ {
		distance := absFloat(candidates[i] - current)
		if distance < bestDistance {
			bestDistance = distance
			bestIndex = i
		}
	}

	return bestIndex
}

func absInt(value int) int {
	if value < 0 {
		return -value
	}

	return value
}

func absFloat(value float64) float64 {
	if value < 0 {
		return -value
	}

	return value
}

func randomIntCandidate(rng *rand.Rand, candidates []int) int {
	return candidates[rng.Intn(len(candidates))]
}

func randomFloatCandidate(rng *rand.Rand, candidates []float64) float64 {
	return candidates[rng.Intn(len(candidates))]
}

func minInt(a int, b int) int {
	if a < b {
		return a
	}

	return b
}
