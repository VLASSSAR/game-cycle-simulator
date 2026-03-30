package domain

type SimulationRequest struct {
	ArrivalRateLambda float64 `json:"arrival_rate_lambda"`
	Mu                float64 `json:"mu"`
	ChannelsK         int     `json:"channels_k"`
	AlgorithmFactorA  float64 `json:"algorithm_factor_a"`
	BetLimit          float64 `json:"bet_limit"`
	FraudFactorF      float64 `json:"fraud_factor"`
	Simulations       int     `json:"simulations"`
}

type OptimizationRequest struct {
	ArrivalRateLambda   float64   `json:"arrival_rate_lambda"`
	Mu                  float64   `json:"mu"`
	Simulations         int       `json:"simulations"`
	ChannelCandidates   []int     `json:"channel_candidates"`
	AlgorithmCandidates []float64 `json:"algorithm_candidates"`
	BetLimitCandidates  []float64 `json:"bet_limit_candidates"`
	FraudCandidates     []float64 `json:"fraud_candidates"`
	MaxRho              float64   `json:"max_rho"`
	MaxProcessingTime   float64   `json:"max_processing_time"`
}
