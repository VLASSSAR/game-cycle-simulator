export interface SimulationRequest {
    arrival_rate_lambda: number;
    mu: number;
    channels_k: number;
    algorithm_factor_a: number;
    bet_limit: number;
    fraud_factor: number;
    simulations: number;
}

export interface Metrics {
    total_bets: number;
    successful_bets: number;
    failed_bets: number;
    success_probability: number;
    failure_probability: number;

    avg_processing_time: number;
    avg_queue_time: number;
    avg_system_time: number;

    service_rate_mu: number;
    utilization_rho: number;

    mean_revenue_per_bet: number;
    revenue_per_unit_time: number;

    fraud_checks: number;
    fraud_detected: number;
    actual_fraud_operations: number;
    false_positives: number;
    false_negatives: number;
    avg_fraud_score: number;
    fraud_detection_rate: number;
    false_positive_rate: number;
    false_negative_rate: number;
}

export interface SimulationResponse {
    request: SimulationRequest;
    metrics: Metrics;
}

export interface BaselineParams {
    channels_k: number;
    algorithm_factor_a: number;
    bet_limit: number;
    fraud_factor: number;
}

export interface OptimizationRequest {
    arrival_rate_lambda: number;
    mu: number;
    simulations: number;
    baseline: BaselineParams;
    channel_candidates: number[];
    algorithm_candidates: number[];
    bet_limit_candidates: number[];
    fraud_candidates: number[];
    max_rho: number;
    max_processing_time: number;
}

export interface BestParams {
    channels_k: number;
    algorithm_factor_a: number;
    bet_limit: number;
    fraud_factor: number;
}

export interface OptimizationResponse {
    request: OptimizationRequest;
    best_params: BestParams;
    baseline_metrics: Metrics;
    optimized_metrics: Metrics;
    delta_revenue: number;
    delta_processing: number;
}