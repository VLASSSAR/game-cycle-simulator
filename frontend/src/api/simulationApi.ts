import type {
    LoadScenario,
    SimulationRequest,
    SimulationResponse,
    OptimizationRequest,
    OptimizationResponse,
    OptimizationComparisonResponse,
} from "../types/simulation";

const API_BASE_URL = "https://game-cycle-simulator.onrender.com";

export async function fetchScenarios(): Promise<LoadScenario[]> {
    const response = await fetch(`${API_BASE_URL}/api/scenarios`, {
        method: "GET",
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Scenarios request failed");
    }

    return response.json();
}

export async function runSimulation(
    payload: SimulationRequest
): Promise<SimulationResponse> {
    const response = await fetch(`${API_BASE_URL}/api/simulate`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Simulation request failed");
    }

    return response.json();
}

export async function runOptimization(
    payload: OptimizationRequest
): Promise<OptimizationResponse> {
    const response = await fetch(`${API_BASE_URL}/api/optimize`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Optimization request failed");
    }

    return response.json();
}

export async function runOptimizationComparison(
    payload: OptimizationRequest
): Promise<OptimizationComparisonResponse> {
    const response = await fetch(`${API_BASE_URL}/api/optimize/compare`, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
    });

    if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Optimization comparison request failed");
    }

    return response.json();
}