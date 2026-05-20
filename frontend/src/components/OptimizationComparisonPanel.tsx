import type { OptimizationComparisonResponse } from "../types/simulation";

interface OptimizationComparisonPanelProps {
    result: OptimizationComparisonResponse;
}

export default function OptimizationComparisonPanel({
                                                        result,
                                                    }: OptimizationComparisonPanelProps) {
    return (
        <div style={styles.card}>
            <h2 style={styles.title}>Сравнение методов оптимизации</h2>

            <div style={styles.section}>
                <h3 style={styles.subtitle}>Лучший метод</h3>
                <p style={styles.text}>{formatMethodName(result.best_method)}</p>
            </div>

            <div style={styles.section}>
                <h3 style={styles.subtitle}>Baseline</h3>
                <ul style={styles.list}>
                    <li>
                        Доход за единицу времени:{" "}
                        {result.baseline_metrics.revenue_per_unit_time.toFixed(4)}
                    </li>
                    <li>
                        Среднее время обработки:{" "}
                        {result.baseline_metrics.avg_processing_time.toFixed(4)}
                    </li>
                    <li>
                        Среднее время ожидания в очереди:{" "}
                        {result.baseline_metrics.avg_queue_time.toFixed(4)}
                    </li>
                    <li>
                        Вероятность отказа:{" "}
                        {result.baseline_metrics.failure_probability.toFixed(4)}
                    </li>
                </ul>
            </div>

            <div style={styles.section}>
                <h3 style={styles.subtitle}>Результаты методов</h3>

                <table style={styles.table}>
                    <thead>
                    <tr>
                        <th style={styles.th}>Метод</th>
                        <th style={styles.th}>Итерации / проверки</th>
                        <th style={styles.th}>Доход / ед. времени</th>
                        <th style={styles.th}>Δ дохода</th>
                        <th style={styles.th}>Среднее время</th>
                        <th style={styles.th}>Очередь</th>
                        <th style={styles.th}>ρ</th>
                        <th style={styles.th}>Fraud Detection</th>
                        <th style={styles.th}>Статус</th>
                    </tr>
                    </thead>

                    <tbody>
                    {result.results.map((item) => {
                        const metrics = item.optimized_metrics;

                        return (
                            <tr key={item.method}>
                                <td style={styles.td}>{formatMethodName(item.method)}</td>
                                <td style={styles.td}>{item.iterations || "—"}</td>
                                <td style={styles.td}>
                                    {metrics
                                        ? metrics.revenue_per_unit_time.toFixed(4)
                                        : "—"}
                                </td>
                                <td style={styles.td}>
                                    {typeof item.delta_revenue === "number"
                                        ? item.delta_revenue.toFixed(4)
                                        : "—"}
                                </td>
                                <td style={styles.td}>
                                    {metrics
                                        ? metrics.avg_processing_time.toFixed(4)
                                        : "—"}
                                </td>
                                <td style={styles.td}>
                                    {metrics
                                        ? metrics.avg_queue_time.toFixed(4)
                                        : "—"}
                                </td>
                                <td style={styles.td}>
                                    {metrics
                                        ? metrics.utilization_rho.toFixed(4)
                                        : "—"}
                                </td>
                                <td style={styles.td}>
                                    {metrics
                                        ? metrics.fraud_detection_rate.toFixed(4)
                                        : "—"}
                                </td>
                                <td style={styles.td}>
                                    {item.error ? item.error : "OK"}
                                </td>
                            </tr>
                        );
                    })}
                    </tbody>
                </table>
            </div>

            <div style={styles.section}>
                <h3 style={styles.subtitle}>Интерпретация</h3>
                <p style={styles.text}>
                    Таблица позволяет сравнить методы оптимизации по доходности,
                    времени обработки, времени ожидания в очереди, коэффициенту загрузки
                    и качеству anti-fraud фильтрации. Эти результаты можно использовать
                    в главе 4 для анализа эффективности разных алгоритмов.
                </p>
            </div>
        </div>
    );
}

function formatMethodName(method: string): string {
    switch (method) {
        case "random":
            return "Random Search";
        case "grid":
            return "Grid Search";
        case "genetic":
            return "Genetic Algorithm";
        case "adaptive":
            return "Adaptive Optimizer";
        default:
            return method;
    }
}

const styles: Record<string, React.CSSProperties> = {
    card: {
        padding: "20px",
        border: "1px solid #ddd",
        borderRadius: "12px",
        background: "#fff",
        minWidth: "720px",
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
    text: {
        margin: 0,
        lineHeight: 1.6,
        color: "#000",
    },
    list: {
        paddingLeft: "18px",
        lineHeight: 1.8,
        margin: 0,
    },
    table: {
        width: "100%",
        borderCollapse: "collapse",
        fontSize: "14px",
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