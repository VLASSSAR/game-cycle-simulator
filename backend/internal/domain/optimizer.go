package domain

type OptimizationMethod string

const (
	OptimizationMethodRandom   OptimizationMethod = "random"
	OptimizationMethodGrid     OptimizationMethod = "grid"
	OptimizationMethodGenetic  OptimizationMethod = "genetic"
	OptimizationMethodAdaptive OptimizationMethod = "adaptive"
)

func (m OptimizationMethod) IsValid() bool {
	switch m {
	case OptimizationMethodRandom,
		OptimizationMethodGrid,
		OptimizationMethodGenetic,
		OptimizationMethodAdaptive:
		return true
	default:
		return false
	}
}

func DefaultOptimizationMethod() OptimizationMethod {
	return OptimizationMethodRandom
}
