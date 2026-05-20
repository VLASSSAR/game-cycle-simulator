package domain

// FraudFeatures описывает признаки пользовательской активности,
// которые используются для расчета интегрального показателя риска.
//
// Эти признаки соответствуют модели антифрод-анализа,
// описанной в разделе 2.5 магистерской диссертации.
type FraudFeatures struct {
	// BetFrequency — частота игровых операций пользователя.
	BetFrequency float64 `json:"bet_frequency"`

	// AvgBetSize — средний размер ставки пользователя.
	AvgBetSize float64 `json:"avg_bet_size"`

	// SessionCount — количество игровых сессий.
	SessionCount float64 `json:"session_count"`

	// GeoRisk — риск, связанный с географическими параметрами подключения.
	GeoRisk float64 `json:"geo_risk"`

	// DeviceRisk — риск, связанный с параметрами устройства.
	DeviceRisk float64 `json:"device_risk"`

	// HistoryRisk — риск, связанный с историей активности пользователя.
	HistoryRisk float64 `json:"history_risk"`

	// LimitViolationRisk — риск, связанный с превышением лимитов ставок.
	LimitViolationRisk float64 `json:"limit_violation_risk"`
}

// FraudWeights задает веса признаков в модели fraud-score.
type FraudWeights struct {
	BetFrequencyWeight       float64 `json:"bet_frequency_weight"`
	AvgBetSizeWeight         float64 `json:"avg_bet_size_weight"`
	SessionCountWeight       float64 `json:"session_count_weight"`
	GeoRiskWeight            float64 `json:"geo_risk_weight"`
	DeviceRiskWeight         float64 `json:"device_risk_weight"`
	HistoryRiskWeight        float64 `json:"history_risk_weight"`
	LimitViolationRiskWeight float64 `json:"limit_violation_risk_weight"`
}

// FraudConfig задает параметры антифрод-модели.
type FraudConfig struct {
	// Threshold — пороговое значение fraud-score,
	// выше которого операция считается подозрительной.
	Threshold float64 `json:"threshold"`

	// Strictness — общий коэффициент строгости антифрод-фильтрации.
	Strictness float64 `json:"strictness"`

	// Weights — веса признаков пользовательской активности.
	Weights FraudWeights `json:"weights"`
}

// FraudResult описывает результат антифрод-проверки одной операции.
type FraudResult struct {
	// Score — рассчитанный интегральный показатель риска.
	Score float64 `json:"score"`

	// IsSuspicious — признак подозрительной операции.
	IsSuspicious bool `json:"is_suspicious"`

	// IsFraud — истинный признак мошеннической операции
	// в рамках имитационной модели.
	IsFraud bool `json:"is_fraud"`

	// IsFalsePositive — добросовестная операция ошибочно признана подозрительной.
	IsFalsePositive bool `json:"is_false_positive"`

	// IsFalseNegative — мошенническая операция не была обнаружена.
	IsFalseNegative bool `json:"is_false_negative"`
}
