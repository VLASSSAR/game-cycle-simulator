import type { SimulationRequest } from "../types/simulation";

interface SimulationFormProps {
    value: SimulationRequest;
    onChange: (data: SimulationRequest) => void;
    onSubmit: (data: SimulationRequest) => Promise<void>;
    isLoading: boolean;
}

export default function SimulationForm({
                                           value,
                                           onChange,
                                           onSubmit,
                                           isLoading,
                                       }: SimulationFormProps) {
    function handleChange(e: React.ChangeEvent<HTMLInputElement>) {
        const { name, value: inputValue } = e.target;

        onChange({
            ...value,
            [name]: Number(inputValue),
        });
    }

    async function handleSubmit(e: React.FormEvent) {
        e.preventDefault();
        await onSubmit(value);
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
                    value={value.arrival_rate_lambda}
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
                    value={value.mu}
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
                    value={value.channels_k}
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
                    value={value.algorithm_factor_a}
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
                    value={value.bet_limit}
                    onChange={handleChange}
                    style={styles.input}
                />
            </label>

            <label style={styles.label}>
                F (anti-fraud параметр)
                <input
                    type="number"
                    step="0.01"
                    min="0"
                    max="1"
                    name="fraud_factor"
                    value={value.fraud_factor}
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
                    value={value.simulations}
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