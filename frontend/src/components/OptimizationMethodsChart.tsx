import {
    BarChart,
    Bar,
    CartesianGrid,
    XAxis,
    YAxis,
    Tooltip,
    ResponsiveContainer,
} from "recharts";
import type {
    OptimizationComparisonResponse,
    OptimizationComparisonItem,
} from "../types/simulation";

interface OptimizationMethodsChartProps {
    result: OptimizationComparisonResponse;
}

interface ChartDataItem {
    method: string;
    value: number;
}

interface SingleMethodChartProps {
    title: string;
    data: ChartDataItem[];
    yDomain?: [number, number];
}

export default function OptimizationMethodsChart({
                                                     result,
                                                 }: OptimizationMethodsChartProps) {
    const successfulResults = result.results.filter(hasMetrics);

    const revenueData = successfulResults.map((item) => ({
        method: formatMethodName(item.method),
        value: Number(item.optimized_metrics.revenue_per_unit_time.toFixed(4)),
    }));

    const processingTimeData = successfulResults.map((item) => ({
        method: formatMethodName(item.method),
        value: Number(item.optimized_metrics.avg_processing_time.toFixed(4)),
    }));

    const queueTimeData = successfulResults.map((item) => ({
        method: formatMethodName(item.method),
        value: Number(item.optimized_metrics.avg_queue_time.toFixed(4)),
    }));

    const fraudDetectionData = successfulResults.map((item) => ({
        method: formatMethodName(item.method),
        value: Number(item.optimized_metrics.fraud_detection_rate.toFixed(4)),
    }));

    if (successfulResults.length === 0) {
        return (
            <div style={styles.wrapper}>
                <h3 style={styles.title}>Графики сравнения методов</h3>
                <p style={styles.emptyText}>
                    Нет успешных результатов оптимизации для построения графиков.
                </p>
            </div>
        );
    }

    return (
        <div style={styles.wrapper}>
            <h3 style={styles.title}>Графики сравнения методов оптимизации</h3>

            <SingleMethodChart
                title="Доход за единицу времени"
                data={revenueData}
            />

            <SingleMethodChart
                title="Среднее время обработки"
                data={processingTimeData}
            />

            <SingleMethodChart
                title="Среднее время ожидания в очереди"
                data={queueTimeData}
            />

            <SingleMethodChart
                title="Fraud Detection Rate"
                data={fraudDetectionData}
                yDomain={[0, 1]}
            />
        </div>
    );
}

function SingleMethodChart({
                               title,
                               data,
                               yDomain,
                           }: SingleMethodChartProps) {
    return (
        <div style={styles.chartWrapper}>
            <h4 style={styles.chartTitle}>{title}</h4>

            <div style={styles.chartBox}>
                <ResponsiveContainer width="100%" height="100%">
                    <BarChart
                        data={data}
                        margin={{ top: 12, right: 20, left: 10, bottom: 10 }}
                    >
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis dataKey="method" />
                        <YAxis domain={yDomain} />
                        <Tooltip />
                        <Bar
                            dataKey="value"
                            name={title}
                            fill="#4C1D95"
                            radius={[6, 6, 0, 0]}
                        />
                    </BarChart>
                </ResponsiveContainer>
            </div>
        </div>
    );
}

function hasMetrics(
    item: OptimizationComparisonItem
): item is OptimizationComparisonItem & {
    optimized_metrics: NonNullable<OptimizationComparisonItem["optimized_metrics"]>;
} {
    return Boolean(item.optimized_metrics);
}

function formatMethodName(method: string): string {
    switch (method) {
        case "random":
            return "Random";
        case "grid":
            return "Grid";
        case "genetic":
            return "Genetic";
        case "adaptive":
            return "Adaptive";
        default:
            return method;
    }
}

const styles: Record<string, React.CSSProperties> = {
    wrapper: {
        padding: "20px",
        border: "1px solid #ddd",
        borderRadius: "12px",
        background: "#fff",
        color: "#000",
        display: "flex",
        flexDirection: "column",
        gap: "24px",
    },
    title: {
        margin: 0,
        color: "#000",
    },
    emptyText: {
        margin: 0,
        color: "#333",
        lineHeight: 1.6,
    },
    chartWrapper: {
        display: "flex",
        flexDirection: "column",
        gap: "8px",
    },
    chartTitle: {
        margin: 0,
        color: "#000",
    },
    chartBox: {
        width: "100%",
        height: "260px",
    },
};