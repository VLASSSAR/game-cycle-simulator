package domain

type SimulationResult struct {
	Request SimulationRequest `json:"request"`
	Metrics Metrics           `json:"metrics"`
}

type BestParams struct {
	ChannelsK        int     `json:"channels_k"`
	AlgorithmFactorA float64 `json:"algorithm_factor_a"`
	BetLimit         float64 `json:"bet_limit"`
	FraudFactorF     float64 `json:"fraud_factor"`
}

type OptimizationResult struct {
	Request OptimizationRequest `json:"request"`

	Method     OptimizationMethod `json:"method"`
	Iterations int                `json:"iterations"`

	BestParams       BestParams `json:"best_params"`
	BaselineMetrics  Metrics    `json:"baseline_metrics"`
	OptimizedMetrics Metrics    `json:"optimized_metrics"`

	DeltaRevenue    float64 `json:"delta_revenue"`
	DeltaProcessing float64 `json:"delta_processing"`
}
