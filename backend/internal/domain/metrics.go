package domain

type Metrics struct {
	TotalBets          int     `json:"total_bets"`
	SuccessfulBets     int     `json:"successful_bets"`
	FailedBets         int     `json:"failed_bets"`
	SuccessProbability float64 `json:"success_probability"`
	FailureProbability float64 `json:"failure_probability"`
	AvgProcessingTime  float64 `json:"avg_processing_time"`
	AvgQueueTime       float64 `json:"avg_queue_time"`
	AvgSystemTime      float64 `json:"avg_system_time"`
	ServiceRateMu      float64 `json:"service_rate_mu"`
	UtilizationRho     float64 `json:"utilization_rho"`
	MeanRevenuePerBet  float64 `json:"mean_revenue_per_bet"`
	RevenuePerUnitTime float64 `json:"revenue_per_unit_time"`
}
