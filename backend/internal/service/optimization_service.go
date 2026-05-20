package service

import (
	"game-cycle-simulator/internal/domain"
	"game-cycle-simulator/internal/optimizer"
)

type OptimizationService struct {
	opt *optimizer.Optimizer
}

func NewOptimizationService() *OptimizationService {
	return &OptimizationService{
		opt: optimizer.NewOptimizer(),
	}
}

func (s *OptimizationService) RunOptimization(
	req domain.OptimizationRequest,
) (domain.OptimizationResult, error) {
	return s.opt.Optimize(req)
}

func (s *OptimizationService) CompareOptimizationMethods(
	req domain.OptimizationRequest,
) (domain.OptimizationComparisonResult, error) {
	return s.opt.Compare(req)
}
