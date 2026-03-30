import type { Metrics } from "../types/simulation";

interface MetricsCardProps {
    metrics: Metrics;
}

export default function MetricsCard({ metrics }: MetricsCardProps) {
    return (
        <div style={styles.card}>
            <h2 style={{color: "#000"}}>Результаты симуляции</h2>
            <ul style={styles.list}>
                <li>Всего ставок: {metrics.total_bets}</li>
                <li>Успешных ставок: {metrics.successful_bets}</li>
                <li>Неуспешных ставок: {metrics.failed_bets}</li>
                <li>Вероятность успеха: {metrics.success_probability.toFixed(4)}</li>
                <li>Среднее время обработки: {metrics.avg_processing_time.toFixed(4)}</li>
                <li>Среднее время в системе: {metrics.avg_system_time.toFixed(4)}</li>
                <li>μ (интенсивность обслуживания): {metrics.service_rate_mu.toFixed(4)}</li>
                <li>ρ (коэффициент загрузки): {metrics.utilization_rho.toFixed(4)}</li>
                <li>Средний доход на ставку: {metrics.mean_revenue_per_bet.toFixed(4)}</li>
                <li>Доход за единицу времени: {metrics.revenue_per_unit_time.toFixed(4)}</li>
            </ul>
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
    },
    list: {
        paddingLeft: "18px",
        lineHeight: 1.8,
    },
};