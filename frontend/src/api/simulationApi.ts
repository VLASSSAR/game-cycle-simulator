import type {
    SimulationRequest,
    SimulationResponse,
    OptimizationRequest,
    OptimizationResponse,
} from "../types/simulation";

const API_BASE_URL = "https://game-cycle-simulator.onrender.com";

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