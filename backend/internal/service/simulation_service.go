package service

import (
	"game-cycle-simulator/internal/domain"
	"game-cycle-simulator/internal/simulator"
)

type SimulationService struct {
	sim *simulator.Simulator
}

func NewSimulationService() *SimulationService {
	return &SimulationService{
		sim: simulator.NewSimulator(),
	}
}

func (s *SimulationService) RunSimulation(
	req domain.SimulationRequest,
) (domain.SimulationResult, error) {
	return s.sim.Run(req)
}
