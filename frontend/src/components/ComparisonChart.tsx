import {
    BarChart,
    Bar,
    CartesianGrid,
    XAxis,
    YAxis,
    Tooltip,
    ResponsiveContainer,
    Legend,
} from "recharts";
import type { Metrics } from "../types/simulation";

interface ComparisonChartProps {
    baseline: Metrics;
    optimized: Metrics;
}

interface ChartDataItem {
    name: string;
    baseline: number;
    optimized: number;
}

interface SingleMetricChartProps {
    title: string;
    data: ChartDataItem[];
    yDomain?: [number, number];
}

function SingleMetricChart({ title, data, yDomain }: SingleMetricChartProps) {
    return (
        <div style={styles.chartWrapper}>
            <h4 style={styles.chartTitle}>{title}</h4>

            <div style={styles.chartBox}>
                <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={data} margin={{ top: 12, right: 20, left: 10, bottom: 10 }}>
                        <CartesianGrid strokeDasharray="3 3" />

                        <XAxis dataKey="name" />
                        <YAxis domain={yDomain} />

                        <Tooltip />
                        <Legend />

                        {/* baseline — серый */}
                        <Bar
                            dataKey="baseline"
                            name="Baseline"
                            fill="#9CA3AF"
                            radius={[6, 6, 0, 0]}
                        />

                        {/* optimized — зелёный */}
                        <Bar
                            dataKey="optimized"
                            name="Optimized"
                            fill="#22C55E"
                            radius={[6, 6, 0, 0]}
                        />
                    </BarChart>
                </ResponsiveContainer>
            </div>
        </div>
    );
}

export default function ComparisonChart({
                                            baseline,
                                            optimized,
                                        }: ComparisonChartProps) {
    const processingData: ChartDataItem[] = [
        {
            name: "Время обработки",
            baseline: Number(baseline.avg_processing_time.toFixed(4)),
            optimized: Number(optimized.avg_processing_time.toFixed(4)),
        },
    ];

    const rhoData: ChartDataItem[] = [
        {
            name: "Загрузка ρ",
            baseline: Number(baseline.utilization_rho.toFixed(4)),
            optimized: Number(optimized.utilization_rho.toFixed(4)),
        },
    ];

    const revenueData: ChartDataItem[] = [
        {
            name: "Доход/ед.времени",
            baseline: Number(baseline.revenue_per_unit_time.toFixed(4)),
            optimized: Number(optimized.revenue_per_unit_time.toFixed(4)),
        },
    ];

    return (
        <div style={styles.wrapper}>
            <h3 style={styles.title}>Графики сравнения режимов</h3>

            <SingleMetricChart
                title="Сравнение среднего времени обработки"
                data={processingData}
            />

            <SingleMetricChart
                title="Сравнение коэффициента загрузки (ρ)"
                data={rhoData}
                yDomain={[0, 1.2]} // важно: фиксируем диапазон
            />

            <SingleMetricChart
                title="Сравнение дохода за единицу времени"
                data={revenueData}
            />
        </div>
    );
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
        marginBottom: "8px",
        color: "#000",
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