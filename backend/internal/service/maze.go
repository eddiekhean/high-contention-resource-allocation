package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/config"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
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
