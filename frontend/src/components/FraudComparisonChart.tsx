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

interface FraudComparisonChartProps {
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

                        <Bar
                            dataKey="baseline"
                            name="Baseline"
                            fill="#9CA3AF"
                            radius={[6, 6, 0, 0]}
                        />

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

export default function FraudComparisonChart({
                                                 baseline,
                                                 optimized,
                                             }: FraudComparisonChartProps) {
    const fraudScoreData: ChartDataItem[] = [
        {
            name: "Fraud-score",
            baseline: Number(baseline.avg_fraud_score.toFixed(4)),
            optimized: Number(optimized.avg_fraud_score.toFixed(4)),
        },
    ];

    const detectionRateData: ChartDataItem[] = [
        {
            name: "Detection Rate",
            baseline: Number(baseline.fraud_detection_rate.toFixed(4)),
            optimized: Number(optimized.fraud_detection_rate.toFixed(4)),
        },
    ];

    const falsePositiveData: ChartDataItem[] = [
        {
            name: "False Positive Rate",
            baseline: Number(baseline.false_positive_rate.toFixed(4)),
            optimized: Number(optimized.false_positive_rate.toFixed(4)),
        },
    ];

    const falseNegativeData: ChartDataItem[] = [
        {
            name: "False Negative Rate",
            baseline: Number(baseline.false_negative_rate.toFixed(4)),
            optimized: Number(optimized.false_negative_rate.toFixed(4)),
        },
    ];

    return (
        <div style={styles.wrapper}>
            <h3 style={styles.title}>Anti-fraud графики</h3>

            <SingleMetricChart
                title="Сравнение среднего fraud-score"
                data={fraudScoreData}
                yDomain={[0, 1]}
            />

            <SingleMetricChart
                title="Сравнение Detection Rate"
                data={detectionRateData}
                yDomain={[0, 1]}
            />

            <SingleMetricChart
                title="Сравнение False Positive Rate"
                data={falsePositiveData}
                yDomain={[0, 1]}
            />

            <SingleMetricChart
                title="Сравнение False Negative Rate"
                data={falseNegativeData}
                yDomain={[0, 1]}
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