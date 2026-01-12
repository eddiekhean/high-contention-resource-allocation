package maze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/utils"
	"github.com/sirupsen/logrus"
)

type MazeService struct {
	client  *http.Client
	baseURL string
	logger  *logrus.Logger
}

func NewMazeService(cfg *config.Config, logger *logrus.Logger) *MazeService {
	return &MazeService{
		client:  &http.Client{Timeout: 10 * time.Second},
		baseURL: cfg.MazeService.URL,
		logger:  logger,
	}
}

func (s *MazeService) GenerateMaze(req models.MazeConfig) (*models.MazeResponse, error) {
	url := fmt.Sprintf("%s/maze/generate", s.baseURL)

	// Default values if not provided (optional, but good practice)
	if req.Rows == 0 {
		req.Rows = 20
	}
	if req.Cols == 0 {
		req.Cols = 20
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	s.logger.WithFields(logrus.Fields{
		"url": url,
		"req": req,
	}).Info("Calling Maze Generator Service")

	resp, err := s.client.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to call maze service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("maze service returned status: %d", resp.StatusCode)
	}

	var mazeResp models.MazeResponse
	if err := json.NewDecoder(resp.Body).Decode(&mazeResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &mazeResp, nil
}

func (s *MazeService) SolveMaze(
	req *models.MazeSolveRequest,
) (*models.SolveResult, error) {

	if !req.Strategy.IsValid() {
		return nil, utils.ErrInvalidStrategy
	}

	graph := s.buildGraph(req)

	solver, err := NewSolver(req.Strategy)
	if err != nil {
		return nil, err
	}

	result, ok := solver.Solve(req.Start, req.End, graph)
	if !ok {
		return nil, utils.ErrMazeUnreachable
	}

	log.Printf(
		"[MazeSolve] strategy=%s path_len=%d steps=%d",
		req.Strategy,
		len(result.Path),
		len(result.Steps),
	)

	return &result, nil
}

func NewSolver(strategy models.SolveStrategy) (MazeSolver, error) {
	switch strategy {
	case models.StrategyBFS:
		return &BFSSolver{}, nil
	case models.StrategyDFS:
		return &DFSSolver{}, nil
	case models.StrategyAStar:
		return &AStarSolver{}, nil
	case models.StrategyGreedy:
		return &GreedySolver{}, nil
	default:
		return nil, utils.ErrInvalidStrategy
	}
}
