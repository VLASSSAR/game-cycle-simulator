package domain

type LoadScenario struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`

	ArrivalRateLambda float64 `json:"arrival_rate_lambda"`
	Mu                float64 `json:"mu"`
	ChannelsK         int     `json:"channels_k"`
	AlgorithmFactorA  float64 `json:"algorithm_factor_a"`
	BetLimit          float64 `json:"bet_limit"`
	FraudFactorF      float64 `json:"fraud_factor"`
	Simulations       int     `json:"simulations"`

	Seed int64 `json:"seed"`

	MaxRho            float64 `json:"max_rho"`
	MaxProcessingTime float64 `json:"max_processing_time"`
}
