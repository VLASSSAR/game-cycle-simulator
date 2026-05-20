import { useState } from "react";
import SimulationForm from "../components/SimulationForm";
import MetricsCard from "../components/MetricsCard";
import OptimizationPanel from "../components/OptimizationPanel";
import { runSimulation, runOptimization } from "../api/simulationApi";
import type {
    SimulationRequest,
    SimulationResponse,
    OptimizationResponse,
    OptimizationRequest,
    OptimizationMethod,
} from "../types/simulation";

export default function HomePage() {
    const [simulationResult, setSimulationResult] = useState<SimulationResponse | null>(null);
    const [optimizationResult, setOptimizationResult] = useState<OptimizationResponse | null>(null);
    const [lastSimulationRequest, setLastSimulationRequest] = useState<SimulationRequest | null>(null);

    const [optimizationMethod, setOptimizationMethod] = useState<OptimizationMethod>("random");
    const [optimizationIterations, setOptimizationIterations] = useState<number>(100);

    const [isSimulating, setIsSimulating] = useState(false);
    const [isOptimizing, setIsOptimizing] = useState(false);
    const [error, setError] = useState<string | null>(null);

    async function handleSimulation(data: SimulationRequest) {
        try {
            setIsSimulating(true);
            setError(null);
            setOptimizationResult(null);

            const response = await runSimulation(data);
            setSimulationResult(response);
            setLastSimulationRequest(data);
        } catch (err) {
            if (err instanceof Error) {
                setError(err.message);
            } else {
                setError("Произошла неизвестная ошибка");
            }
        } finally {
            setIsSimulating(false);
        }
    }

    async function handleOptimization() {
        if (!lastSimulationRequest) {
            setError("Сначала запустите симуляцию, чтобы задать базовые параметры.");
            return;
        }

        if (optimizationIterations <= 0) {
            setError("Число итераций оптимизации должно быть больше 0.");
            return;
        }

        if (optimizationMethod !== "random") {
            setError("Пока реализован только метод Random Search. Остальные методы будут добавлены далее.");
            return;
        }

        try {
            setIsOptimizing(true);
            setError(null);

            const payload: OptimizationRequest = {
                method: optimizationMethod,
                iterations: optimizationIterations,

                arrival_rate_lambda: lastSimulationRequest.arrival_rate_lambda,
                mu: lastSimulationRequest.mu,
                simulations: lastSimulationRequest.simulations,

                baseline: {
                    channels_k: lastSimulationRequest.channels_k,
                    algorithm_factor_a: lastSimulationRequest.algorithm_factor_a,
                    bet_limit: lastSimulationRequest.bet_limit,
                    fraud_factor: lastSimulationRequest.fraud_factor,
                },

                channel_candidates: [
                    Math.max(1, lastSimulationRequest.channels_k - 1),
                    lastSimulationRequest.channels_k,
                    lastSimulationRequest.channels_k + 1,
                    lastSimulationRequest.channels_k + 2,
                ],

                algorithm_candidates: [
                    Math.max(0.5, Number((lastSimulationRequest.algorithm_factor_a - 0.2).toFixed(2))),
                    lastSimulationRequest.algorithm_factor_a,
                    Number((lastSimulationRequest.algorithm_factor_a + 0.2).toFixed(2)),
                    Number((lastSimulationRequest.algorithm_factor_a + 0.5).toFixed(2)),
                ],

                bet_limit_candidates: [
                    Math.max(10, Number((lastSimulationRequest.bet_limit * 0.5).toFixed(2))),
                    Number((lastSimulationRequest.bet_limit * 0.75).toFixed(2)),
                    lastSimulationRequest.bet_limit,
                    Number((lastSimulationRequest.bet_limit * 1.25).toFixed(2)),
                    Number((lastSimulationRequest.bet_limit * 1.5).toFixed(2)),
                ],

                fraud_candidates: [
                    Math.max(0, Number((lastSimulationRequest.fraud_factor - 0.1).toFixed(2))),
                    lastSimulationRequest.fraud_factor,
                    Math.min(1, Number((lastSimulationRequest.fraud_factor + 0.1).toFixed(2))),
                    Math.min(1, Number((lastSimulationRequest.fraud_factor + 0.2).toFixed(2))),
                ],

                max_rho: 0.95,
                max_processing_time: 0.3,
            };

            const response = await runOptimization(payload);
            setOptimizationResult(response);
        } catch (err) {
            if (err instanceof Error) {
                setError(err.message);
            } else {
                setError("Произошла неизвестная ошибка");
            }
        } finally {
            setIsOptimizing(false);
        }
    }

    return (
        <div style={styles.page}>
            <h1 style={styles.title}>Симулятор игрового цикла</h1>
            <p style={styles.subtitle}>
                Имитационное моделирование и оптимизация игрового цикла онлайн-платформы
            </p>

            <div style={styles.layout}>
                <div style={styles.leftColumn}>
                    <SimulationForm onSubmit={handleSimulation} isLoading={isSimulating} />

                    <div style={styles.optimizationSettings}>
                        <h2 style={styles.settingsTitle}>Параметры оптимизации</h2>

                        <label style={styles.label}>
                            Метод оптимизации
                            <select
                                value={optimizationMethod}
                                onChange={(event) =>
                                    setOptimizationMethod(event.target.value as OptimizationMethod)
                                }
                                style={styles.input}
                            >
                                <option value="random">Random Search</option>
                                <option value="grid">Grid Search — будет добавлен далее</option>
                                <option value="genetic">Genetic Algorithm — будет добавлен далее</option>
                                <option value="adaptive">Adaptive Optimizer — будет добавлен далее</option>
                            </select>
                        </label>

                        <label style={styles.label}>
                            Число итераций
                            <input
                                type="number"
                                min="1"
                                step="1"
                                value={optimizationIterations}
                                onChange={(event) =>
                                    setOptimizationIterations(Number(event.target.value))
                                }
                                style={styles.input}
                            />
                        </label>
                    </div>

                    <button
                        onClick={handleOptimization}
                        disabled={isOptimizing || !lastSimulationRequest}
                        style={styles.optimizeButton}
                    >
                        {isOptimizing ? "Оптимизация..." : "Оптимизировать параметры"}
                    </button>
                </div>

                <div style={styles.rightColumn}>
                    {error && <div style={styles.error}>{error}</div>}

                    {simulationResult && (
                        <div style={styles.block}>
                            <MetricsCard metrics={simulationResult.metrics} />
                        </div>
                    )}

                    {optimizationResult && (
                        <div style={styles.block}>
                            <OptimizationPanel result={optimizationResult} />
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

const styles: Record<string, React.CSSProperties> = {
    page: {
        minHeight: "100vh",
        padding: "32px",
        background: "#E6E6FA",
        fontFamily: "Arial, sans-serif",
    },
    title: {
        color: "#000",
        margin: "20px",
    },
    subtitle: {
        color: "#333",
        marginTop: "8px",
        marginBottom: "24px",
    },
    layout: {
        display: "flex",
        gap: "24px",
        alignItems: "flex-start",
        flexWrap: "wrap",
    },
    leftColumn: {
        display: "flex",
        flexDirection: "column",
        gap: "16px",
    },
    rightColumn: {
        display: "flex",
        flexDirection: "column",
        gap: "16px",
        flex: 1,
        minWidth: "520px",
    },
    block: {
        display: "flex",
        flexDirection: "column",
        gap: "12px",
    },
    optimizationSettings: {
        display: "flex",
        flexDirection: "column",
        gap: "12px",
        padding: "20px",
        border: "1px solid #ddd",
        borderRadius: "12px",
        maxWidth: "440px",
        background: "#fff",
        color: "#000",
    },
    settingsTitle: {
        color: "#000",
        margin: 0,
        marginBottom: "8px",
    },
    label: {
        display: "flex",
        flexDirection: "column",
        gap: "6px",
        fontWeight: 500,
        color: "#000",
    },
    input: {
        padding: "10px",
        borderRadius: "8px",
        border: "1px solid #ccc",
        color: "#000",
        background: "#fff",
    },
    optimizeButton: {
        padding: "12px 16px",
        borderRadius: "8px",
        border: "none",
        background: "#4C1D95",
        color: "#fff",
        cursor: "pointer",
        fontWeight: 600,
    },
    error: {
        padding: "12px 16px",
        borderRadius: "8px",
        background: "#fee2e2",
        color: "#991b1b",
        border: "1px solid #fecaca",
    },
};