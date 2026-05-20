import type { Metrics } from "../types/simulation";

interface MetricsCardProps {
    metrics: Metrics;
}

export default function MetricsCard({ metrics }: MetricsCardProps) {
    return (
        <div style={styles.card}>
            <h2 style={styles.title}>Результаты симуляции</h2>

            <section style={styles.section}>
                <h3 style={styles.subtitle}>Основные показатели</h3>
                <ul style={styles.list}>
                    <li>Всего ставок: {metrics.total_bets}</li>
                    <li>Успешных ставок: {metrics.successful_bets}</li>
                    <li>Неуспешных ставок: {metrics.failed_bets}</li>
                    <li>Вероятность успеха: {metrics.success_probability.toFixed(4)}</li>
                    <li>Вероятность отказа: {metrics.failure_probability.toFixed(4)}</li>
                </ul>
            </section>

            <section style={styles.section}>
                <h3 style={styles.subtitle}>Временные характеристики</h3>
                <ul style={styles.list}>
                    <li>Среднее время обработки: {metrics.avg_processing_time.toFixed(4)}</li>
                    <li>Среднее время ожидания в очереди: {metrics.avg_queue_time.toFixed(4)}</li>
                    <li>Среднее время в системе: {metrics.avg_system_time.toFixed(4)}</li>
                    <li>μ (интенсивность обслуживания): {metrics.service_rate_mu.toFixed(4)}</li>
                    <li>ρ (коэффициент загрузки): {metrics.utilization_rho.toFixed(4)}</li>
                </ul>
            </section>

            <section style={styles.section}>
                <h3 style={styles.subtitle}>Экономические показатели</h3>
                <ul style={styles.list}>
                    <li>Средний доход на ставку: {metrics.mean_revenue_per_bet.toFixed(4)}</li>
                    <li>Доход за единицу времени: {metrics.revenue_per_unit_time.toFixed(4)}</li>
                </ul>
            </section>

            <section style={styles.section}>
                <h3 style={styles.subtitle}>Anti-fraud показатели</h3>
                <ul style={styles.list}>
                    <li>Количество anti-fraud проверок: {metrics.fraud_checks}</li>
                    <li>Подозрительных операций: {metrics.fraud_detected}</li>
                    <li>Фактических fraud-операций: {metrics.actual_fraud_operations}</li>
                    <li>False Positive: {metrics.false_positives}</li>
                    <li>False Negative: {metrics.false_negatives}</li>
                    <li>Средний fraud-score: {metrics.avg_fraud_score.toFixed(4)}</li>
                    <li>Detection Rate: {metrics.fraud_detection_rate.toFixed(4)}</li>
                    <li>False Positive Rate: {metrics.false_positive_rate.toFixed(4)}</li>
                    <li>False Negative Rate: {metrics.false_negative_rate.toFixed(4)}</li>
                </ul>
            </section>
        </div>
    );
}

const styles: Record<string, React.CSSProperties> = {
    card: {
        padding: "20px",
        border: "1px solid #ddd",
        borderRadius: "12px",
        background: "#fff",
        minWidth: "360px",
        color: "#000",
    },
    title: {
        color: "#000",
        margin: 0,
        marginBottom: "16px",
    },
    subtitle: {
        color: "#000",
        margin: 0,
        marginBottom: "8px",
        fontSize: "16px",
    },
    section: {
        marginBottom: "18px",
    },
    list: {
        paddingLeft: "18px",
        lineHeight: 1.8,
        margin: 0,
    },
};