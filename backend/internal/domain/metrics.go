package domain

type Metrics struct {
	TotalBets          int     `json:"total_bets"`
	SuccessfulBets     int     `json:"successful_bets"`
	FailedBets         int     `json:"failed_bets"`
	SuccessProbability float64 `json:"success_probability"`
	FailureProbability float64 `json:"failure_probability"`

	AvgProcessingTime float64 `json:"avg_processing_time"`
	AvgQueueTime      float64 `json:"avg_queue_time"`
	AvgSystemTime     float64 `json:"avg_system_time"`

	ServiceRateMu  float64 `json:"service_rate_mu"`
	UtilizationRho float64 `json:"utilization_rho"`

	MeanRevenuePerBet  float64 `json:"mean_revenue_per_bet"`
	RevenuePerUnitTime float64 `json:"revenue_per_unit_time"`

	FraudChecks           int     `json:"fraud_checks"`
	FraudDetected         int     `json:"fraud_detected"`
	ActualFraudOperations int     `json:"actual_fraud_operations"`
	FalsePositives        int     `json:"false_positives"`
	FalseNegatives        int     `json:"false_negatives"`
	AvgFraudScore         float64 `json:"avg_fraud_score"`
	FraudDetectionRate    float64 `json:"fraud_detection_rate"`
	FalsePositiveRate     float64 `json:"false_positive_rate"`
	FalseNegativeRate     float64 `json:"false_negative_rate"`
}
