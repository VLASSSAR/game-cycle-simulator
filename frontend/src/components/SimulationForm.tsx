import { useState } from "react";
import type { SimulationRequest } from "../types/simulation";

interface SimulationFormProps {
    onSubmit: (data: SimulationRequest) => Promise<void>;
    isLoading: boolean;
}

export default function SimulationForm({
                                           onSubmit,
                                           isLoading,
                                       }: SimulationFormProps) {
    const [formData, setFormData] = useState<SimulationRequest>({
        arrival_rate_lambda: 8,
        mu: 5,
        channels_k: 3,
        algorithm_factor_a: 1,
        bet_limit: 100,
        fraud_factor: 0.2,
        simulations: 10000,
    });

    function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
        const { name, value } = e.target;

        setFormData((prev) => ({
            ...prev,
            [name]: Number(value),
        }));
    }

    async function handleSubmit(e: React.FormEvent) {
        e.preventDefault();
        await onSubmit(formData);
    }

    return (
        <form onSubmit={handleSubmit} style={styles.form}>
            <h2 style={styles.title}>Параметры симуляции</h2>

            <label style={styles.label}>
                λ (интенсивность потока ставок)
                <input
                    type="number"
                    step="0.1"
                    name="arrival_rate_lambda"
                    value={formData.arrival_rate_lambda}
                    onChange={handleChange}
                    style={styles.input}
                />
            </label>

            <label style={styles.label}>
                μ (базовая скорость обработки)
                <input
                    type="number"
                    step="0.1"
                    name="mu"
                    value={formData.mu}
                    onChange={handleChange}
                    style={styles.input}
                />
            </label>

            <label style={styles.label}>
                k (число каналов)
                <input
                    type="number"
                    step="1"
                    name="channels_k"
                    value={formData.channels_k}
                    onChange={handleChange}
                    style={styles.input}
                />
            </label>

            <label style={styles.label}>
                a (алгоритмический коэффициент)
                <input
                    type="number"
                    step="0.1"
                    name="algorithm_factor_a"
                    value={formData.algorithm_factor_a}
                    onChange={handleChange}
                    style={styles.input}
                />
            </label>

            <label style={styles.label}>
                L (лимит ставки)
                <input
                    type="number"
                    step="1"
                    name="bet_limit"
                    value={formData.bet_limit}
                    onChange={handleChange}
                    style={styles.input}
                />
            </label>

            <label style={styles.label}>
                F (антифрод-параметр)
                <input
                    type="number"
                    step="0.01"
                    min="0"
                    max="1"
                    name="fraud_factor"
                    value={formData.fraud_factor}
                    onChange={handleChange}
                    style={styles.input}
                />
            </label>

            <label style={styles.label}>
                Число симуляций
                <input
                    type="number"
                    step="1"
                    name="simulations"
                    value={formData.simulations}
                    onChange={handleChange}
                    style={styles.input}
                />
            </label>

            <button type="submit" disabled={isLoading} style={styles.button}>
                {isLoading ? "Выполняется..." : "Запустить симуляцию"}
            </button>
        </form>
    );
}

const styles: Record<string, React.CSSProperties> = {
    form: {
        display: "flex",
        flexDirection: "column",
        gap: "12px",
        padding: "20px",
        border: "1px solid #ddd",
        borderRadius: "12px",
        maxWidth: "440px",
        background: "#fff",
    },
    title: {
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
    button: {
        padding: "12px",
        borderRadius: "8px",
        border: "none",
        background: "#111827",
        color: "#fff",
        cursor: "pointer",
        fontWeight: 600,
    },
};