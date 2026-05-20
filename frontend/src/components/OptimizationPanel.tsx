import type { OptimizationResponse } from "../types/simulation";
import ComparisonChart from "./ComparisonChart";
import FraudComparisonChart from "./FraudComparisonChart";

interface OptimizationPanelProps {
    result: OptimizationResponse;
}

export default function OptimizationPanel({ result }: OptimizationPanelProps) {
    const {
        best_params,
        baseline_metrics,
        optimized_metrics,
        delta_revenue,
        delta_processing,
    } = result;

    return (
        <div style={styles.card}>
            <h2 style={styles.title}>Результаты оптимизации</h2>

            <div style={styles.section}>
                <h3 style={styles.subtitle}>Лучшие параметры</h3>
                <ul style={styles.list}>
                    <li>k (число каналов): {best_params.channels_k}</li>
                    <li>a (алгоритмический коэффициент): {best_params.algorithm_factor_a}</li>
                    <li>L (лимит ставки): {best_params.bet_limit}</li>
                    <li>F (anti-fraud параметр): {best_params.fraud_factor}</li>
                </ul>
            </div>

            <div style={styles.section}>
                <h3 style={styles.subtitle}>Сравнение режимов</h3>

                <table style={styles.table}>
                    <thead>
                    <tr>
                        <th style={styles.th}>Метрика</th>
                        <th style={styles.th}>Baseline</th>
                        <th style={styles.th}>Optimized</th>
                    </tr>
                    </thead>

                    <tbody>
                    <tr>
                        <td style={styles.td}>Среднее время обработки</td>
                        <td style={styles.td}>{baseline_metrics.avg_processing_time.toFixed(4)}</td>
                        <td style={styles.td}>{optimized_metrics.avg_processing_time.toFixed(4)}</td>
                    </tr>

                    <tr>
                        <td style={styles.td}>Среднее время ожидания в очереди</td>
                        <td style={styles.td}>{baseline_metrics.avg_queue_time.toFixed(4)}</td>
                        <td style={styles.td}>{optimized_metrics.avg_queue_time.toFixed(4)}</td>
                    </tr>

                    <tr>
                        <td style={styles.td}>Среднее время в системе</td>
                        <td style={styles.td}>{baseline_metrics.avg_system_time.toFixed(4)}</td>
                        <td style={styles.td}>{optimized_metrics.avg_system_time.toFixed(4)}</td>
                    </tr>

                    <tr>
                        <td style={styles.td}>Вероятность успеха</td>
                        <td style={styles.td}>{baseline_metrics.success_probability.toFixed(4)}</td>
                        <td style={styles.td}>{optimized_metrics.success_probability.toFixed(4)}</td>
                    </tr>

                    <tr>
                        <td style={styles.td}>Вероятность отказа</td>
                        <td style={styles.td}>{baseline_metrics.failure_probability.toFixed(4)}</td>
                        <td style={styles.td}>{optimized_metrics.failure_probability.toFixed(4)}</td>
                    </tr>

                    <tr>
                        <td style={styles.td}>ρ (коэффициент загрузки)</td>
                        <td style={styles.td}>{baseline_metrics.utilization_rho.toFixed(4)}</td>
                        <td style={styles.td}>{optimized_metrics.utilization_rho.toFixed(4)}</td>
                    </tr>

                    <tr>
                        <td style={styles.td}>Средний доход на ставку</td>
                        <td style={styles.td}>{baseline_metrics.mean_revenue_per_bet.toFixed(4)}</td>
                        <td style={styles.td}>{optimized_metrics.mean_revenue_per_bet.toFixed(4)}</td>
                    </tr>

                    <tr>
                        <td style={styles.td}>Доход за единицу времени</td>
                        <td style={styles.td}>{baseline_metrics.revenue_per_unit_time.toFixed(4)}</td>
                        <td style={styles.td}>{optimized_metrics.revenue_per_unit_time.toFixed(4)}</td>
                    </tr>

                    <tr>
                        <td style={styles.td}>Средний fraud-score</td>
                        <td style={styles.td}>{baseline_metrics.avg_fraud_score.toFixed(4)}</td>
                        <td style={styles.td}>{optimized_metrics.avg_fraud_score.toFixed(4)}</td>
                    </tr>

                    <tr>
                        <td style={styles.td}>Detection Rate</td>
                        <td style={styles.td}>{baseline_metrics.fraud_detection_rate.toFixed(4)}</td>
                        <td style={styles.td}>{optimized_metrics.fraud_detection_rate.toFixed(4)}</td>
                    </tr>

                    <tr>
                        <td style={styles.td}>False Positive Rate</td>
                        <td style={styles.td}>{baseline_metrics.false_positive_rate.toFixed(4)}</td>
                        <td style={styles.td}>{optimized_metrics.false_positive_rate.toFixed(4)}</td>
                    </tr>

                    <tr>
                        <td style={styles.td}>False Negative Rate</td>
                        <td style={styles.td}>{baseline_metrics.false_negative_rate.toFixed(4)}</td>
                        <td style={styles.td}>{optimized_metrics.false_negative_rate.toFixed(4)}</td>
                    </tr>
                    </tbody>
                </table>
            </div>

            <div style={styles.section}>
                <h3 style={styles.subtitle}>Прирост</h3>
                <ul style={styles.list}>
                    <li>Δ дохода: {delta_revenue.toFixed(4)}</li>
                    <li>Δ времени обработки: {delta_processing.toFixed(4)}</li>
                </ul>
            </div>

            <ComparisonChart
                baseline={baseline_metrics}
                optimized={optimized_metrics}
            />

            <FraudComparisonChart
                baseline={baseline_metrics}
                optimized={optimized_metrics}
            />
        </div>
    );
}

const styles: Record<string, React.CSSProperties> = {
    card: {
        padding: "20px",
        border: "1px solid #ddd",
        borderRadius: "12px",
        background: "#fff",
        minWidth: "520px",
        color: "#000",
    },
    title: {
        color: "#000",
        marginBottom: "16px",
    },
    subtitle: {
        color: "#000",
        marginBottom: "8px",
    },
    section: {
        marginBottom: "20px",
    },
    list: {
        paddingLeft: "18px",
        lineHeight: 1.8,
    },
    table: {
        width: "100%",
        borderCollapse: "collapse",
    },
    th: {
        border: "1px solid #ccc",
        padding: "8px",
        background: "#f3f4f6",
        textAlign: "left",
    },
    td: {
        border: "1px solid #ccc",
        padding: "8px",
    },
};