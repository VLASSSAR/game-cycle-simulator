import { useEffect, useState } from "react";
import SimulationForm from "../components/SimulationForm";
import MetricsCard from "../components/MetricsCard";
import OptimizationPanel from "../components/OptimizationPanel";
import OptimizationComparisonPanel from "../components/OptimizationComparisonPanel";
import {
    fetchScenarios,
    runSimulation,
    runOptimization,
    runOptimizationComparison,
} from "../api/simulationApi";
import type {
    LoadScenario,
    SimulationRequest,
    SimulationResponse,
    OptimizationResponse,
    OptimizationRequest,
    OptimizationMethod,
    OptimizationComparisonResponse,
} from "../types/simulation";

const defaultSimulationRequest: SimulationRequest = {
    arrival_rate_lambda: 8,
    mu: 5,
    channels_k: 3,
    algorithm_factor_a: 1,
    bet_limit: 100,
    fraud_factor: 0.2,
    simulations: 10000,
};

export default function HomePage() {
    const [simulationFormData, setSimulationFormData] = useState<SimulationRequest>(
        defaultSimulationRequest
    );

    const [scenarios, setScenarios] = useState<LoadScenario[]>([]);
    const [selectedScenarioId, setSelectedScenarioId] = useState<string>("custom");

    const [simulationResult, setSimulationResult] = useState<SimulationResponse | null>(null);
    const [optimizationResult, setOptimizationResult] = useState<OptimizationResponse | null>(null);
    const [comparisonResult, setComparisonResult] = useState<OptimizationComparisonResponse | null>(null);
    const [lastSimulationRequest, setLastSimulationRequest] = useState<SimulationRequest | null>(null);

    const [optimizationMethod, setOptimizationMethod] = useState<OptimizationMethod>("random");
    const [optimizationIterations, setOptimizationIterations] = useState<number>(100);

    const [maxRho, setMaxRho] = useState<number>(0.95);
    const [maxProcessingTime, setMaxProcessingTime] = useState<number>(0.3);

    const [isLoadingScenarios, setIsLoadingScenarios] = useState(false);
    const [isSimulating, setIsSimulating] = useState(false);
    const [isOptimizing, setIsOptimizing] = useState(false);
    const [isComparing, setIsComparing] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        async function loadScenarios() {
            try {
                setIsLoadingScenarios(true);

                const response = await fetchScenarios();
                setScenarios(response);
            } catch (err) {
                if (err instanceof Error) {
                    setError(err.message);
                } else {
                    setError("Не удалось загрузить сценарии нагрузки");
                }
            } finally {
                setIsLoadingScenarios(false);
            }
        }

        loadScenarios();
    }, []);

    function handleScenarioChange(scenarioId: string) {
        setSelectedScenarioId(scenarioId);

        if (scenarioId === "custom") {
            return;
        }

        const scenario = scenarios.find((item) => item.id === scenarioId);

        if (!scenario) {
            return;
        }

        setSimulationFormData({
            arrival_rate_lambda: scenario.arrival_rate_lambda,
            mu: scenario.mu,
            channels_k: scenario.channels_k,
            algorithm_factor_a: scenario.algorithm_factor_a,
            bet_limit: scenario.bet_limit,
            fraud_factor: scenario.fraud_factor,
            simulations: scenario.simulations,
        });

        setMaxRho(scenario.max_rho);
        setMaxProcessingTime(scenario.max_processing_time);

        setSimulationResult(null);
        setOptimizationResult(null);
        setComparisonResult(null);
        setLastSimulationRequest(null);
    }

    function handleSimulationFormChange(data: SimulationRequest) {
        setSimulationFormData(data);
        setSelectedScenarioId("custom");
    }

    async function handleSimulation(data: SimulationRequest) {
        try {
            setIsSimulating(true);
            setError(null);
            setOptimizationResult(null);
            setComparisonResult(null);

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

        try {
            setIsOptimizing(true);
            setError(null);
            setComparisonResult(null);

            const payload = buildOptimizationPayload(
                lastSimulationRequest,
                optimizationMethod,
                optimizationIterations,
                maxRho,
                maxProcessingTime
            );

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

    async function handleOptimizationComparison() {
        if (!lastSimulationRequest) {
            setError("Сначала запустите симуляцию, чтобы задать базовые параметры.");
            return;
        }

        if (optimizationIterations <= 0) {
            setError("Число итераций / проверок должно быть больше 0.");
            return;
        }

        try {
            setIsComparing(true);
            setError(null);
            setOptimizationResult(null);

            const payload = buildOptimizationPayload(
                lastSimulationRequest,
                optimizationMethod,
                optimizationIterations,
                maxRho,
                maxProcessingTime
            );

            const response = await runOptimizationComparison(payload);
            setComparisonResult(response);
        } catch (err) {
            if (err instanceof Error) {
                setError(err.message);
            } else {
                setError("Произошла неизвестная ошибка");
            }
        } finally {
            setIsComparing(false);
        }
    }

    const selectedScenario = scenarios.find((item) => item.id === selectedScenarioId);

    return (
        <div style={styles.page}>
            <h1 style={styles.title}>Симулятор игрового цикла</h1>
            <p style={styles.subtitle}>
                Имитационное моделирование и оптимизация игрового цикла онлайн-платформы
            </p>

            <div style={styles.layout}>
                <div style={styles.leftColumn}>
                    <div style={styles.scenarioCard}>
                        <h2 style={styles.settingsTitle}>Сценарий нагрузки</h2>

                        <label style={styles.label}>
                            Выбор сценария
                            <select
                                value={selectedScenarioId}
                                onChange={(event) => handleScenarioChange(event.target.value)}
                                style={styles.input}
                            >
                                <option value="custom">Пользовательские параметры</option>
                                {scenarios.map((scenario) => (
                                    <option key={scenario.id} value={scenario.id}>
                                        {scenario.name}
                                    </option>
                                ))}
                            </select>
                        </label>

                        {isLoadingScenarios && (
                            <p style={styles.hint}>Загрузка сценариев...</p>
                        )}

                        {selectedScenario && (
                            <p style={styles.hint}>{selectedScenario.description}</p>
                        )}
                    </div>

                    <SimulationForm
                        value={simulationFormData}
                        onChange={handleSimulationFormChange}
                        onSubmit={handleSimulation}
                        isLoading={isSimulating}
                    />

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
                                <option value="grid">Grid Search</option>
                                <option value="genetic">Genetic Algorithm</option>
                                <option value="adaptive">Adaptive Optimizer</option>
                            </select>
                        </label>

                        <label style={styles.label}>
                            Число итераций / проверок
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

                        <label style={styles.label}>
                            Максимальный ρ
                            <input
                                type="number"
                                min="0"
                                step="0.01"
                                value={maxRho}
                                onChange={(event) =>
                                    setMaxRho(Number(event.target.value))
                                }
                                style={styles.input}
                            />
                        </label>

                        <label style={styles.label}>
                            Максимальное время обработки
                            <input
                                type="number"
                                min="0"
                                step="0.01"
                                value={maxProcessingTime}
                                onChange={(event) =>
                                    setMaxProcessingTime(Number(event.target.value))
                                }
                                style={styles.input}
                            />
                        </label>
                    </div>

                    <button
                        onClick={handleOptimization}
                        disabled={isOptimizing || isComparing || !lastSimulationRequest}
                        style={styles.optimizeButton}
                    >
                        {isOptimizing ? "Оптимизация..." : "Оптимизировать параметры"}
                    </button>

                    <button
                        onClick={handleOptimizationComparison}
                        disabled={isOptimizing || isComparing || !lastSimulationRequest}
                        style={styles.compareButton}
                    >
                        {isComparing ? "Сравнение..." : "Сравнить методы оптимизации"}
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

                    {comparisonResult && (
                        <div style={styles.block}>
                            <OptimizationComparisonPanel result={comparisonResult} />
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

function buildOptimizationPayload(
    request: SimulationRequest,
    method: OptimizationMethod,
    iterations: number,
    maxRho: number,
    maxProcessingTime: number
): OptimizationRequest {
    return {
        method,
        iterations,

        arrival_rate_lambda: request.arrival_rate_lambda,
        mu: request.mu,
        simulations: request.simulations,

        baseline: {
            channels_k: request.channels_k,
            algorithm_factor_a: request.algorithm_factor_a,
            bet_limit: request.bet_limit,
            fraud_factor: request.fraud_factor,
        },

        channel_candidates: [
            Math.max(1, request.channels_k - 1),
            request.channels_k,
            request.channels_k + 1,
            request.channels_k + 2,
        ],

        algorithm_candidates: [
            Math.max(0.5, Number((request.algorithm_factor_a - 0.2).toFixed(2))),
            request.algorithm_factor_a,
            Number((request.algorithm_factor_a + 0.2).toFixed(2)),
            Number((request.algorithm_factor_a + 0.5).toFixed(2)),
        ],

        bet_limit_candidates: [
            Math.max(10, Number((request.bet_limit * 0.5).toFixed(2))),
            Number((request.bet_limit * 0.75).toFixed(2)),
            request.bet_limit,
            Number((request.bet_limit * 1.25).toFixed(2)),
            Number((request.bet_limit * 1.5).toFixed(2)),
        ],

        fraud_candidates: [
            Math.max(0, Number((request.fraud_factor - 0.1).toFixed(2))),
            request.fraud_factor,
            Math.min(1, Number((request.fraud_factor + 0.1).toFixed(2))),
            Math.min(1, Number((request.fraud_factor + 0.2).toFixed(2))),
        ],

        max_rho: maxRho,
        max_processing_time: maxProcessingTime,
    };
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
    scenarioCard: {
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
    hint: {
        margin: 0,
        color: "#333",
        lineHeight: 1.5,
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
    compareButton: {
        padding: "12px 16px",
        borderRadius: "8px",
        border: "none",
        background: "#047857",
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