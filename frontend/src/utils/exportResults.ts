import type {
    OptimizationComparisonResponse,
    OptimizationComparisonItem,
} from "../types/simulation";

export function exportOptimizationComparisonJson(
    result: OptimizationComparisonResponse
) {
    downloadTextFile(
        "optimization-comparison.json",
        JSON.stringify(result, null, 2),
        "application/json"
    );
}

export function exportOptimizationComparisonCsv(
    result: OptimizationComparisonResponse
) {
    const rows = [
        [
            "method",
            "iterations",
            "revenue_per_unit_time",
            "delta_revenue",
            "avg_processing_time",
            "avg_queue_time",
            "avg_system_time",
            "utilization_rho",
            "failure_probability",
            "avg_fraud_score",
            "fraud_detection_rate",
            "false_positive_rate",
            "false_negative_rate",
            "channels_k",
            "algorithm_factor_a",
            "bet_limit",
            "fraud_factor",
            "status",
        ],
        ...result.results.map((item) => comparisonItemToCsvRow(item)),
    ];

    const csv = rows.map((row) => row.map(escapeCsvValue).join(";")).join("\n");

    downloadTextFile(
        "optimization-comparison.csv",
        csv,
        "text/csv;charset=utf-8"
    );
}

function comparisonItemToCsvRow(item: OptimizationComparisonItem): string[] {
    const metrics = item.optimized_metrics;
    const params = item.best_params;

    return [
        item.method,
        String(item.iterations || ""),
        formatNumber(metrics?.revenue_per_unit_time),
        formatNumber(item.delta_revenue),
        formatNumber(metrics?.avg_processing_time),
        formatNumber(metrics?.avg_queue_time),
        formatNumber(metrics?.avg_system_time),
        formatNumber(metrics?.utilization_rho),
        formatNumber(metrics?.failure_probability),
        formatNumber(metrics?.avg_fraud_score),
        formatNumber(metrics?.fraud_detection_rate),
        formatNumber(metrics?.false_positive_rate),
        formatNumber(metrics?.false_negative_rate),
        params ? String(params.channels_k) : "",
        formatNumber(params?.algorithm_factor_a),
        formatNumber(params?.bet_limit),
        formatNumber(params?.fraud_factor),
        item.error ? item.error : "OK",
    ];
}

function formatNumber(value: number | undefined): string {
    if (typeof value !== "number" || Number.isNaN(value)) {
        return "";
    }

    return value.toFixed(6);
}

function escapeCsvValue(value: string): string {
    const preparedValue = value ?? "";

    if (
        preparedValue.includes(";") ||
        preparedValue.includes("\"") ||
        preparedValue.includes("\n")
    ) {
        return `"${preparedValue.replaceAll("\"", "\"\"")}"`;
    }

    return preparedValue;
}

function downloadTextFile(
    fileName: string,
    content: string,
    mimeType: string
) {
    const blob = new Blob([content], { type: mimeType });
    const url = URL.createObjectURL(blob);

    const link = document.createElement("a");
    link.href = url;
    link.download = fileName;

    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);

    URL.revokeObjectURL(url);
}