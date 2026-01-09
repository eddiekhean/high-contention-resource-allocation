package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/domain"
	"github.com/eddiekhean/high-contention-resource-allocation-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Manual mock
type mockImageService struct {
	matchFunc func(dhashStr string) (*domain.Image, int, bool, error)
}

func (m *mockImageService) Match(ctx context.Context, dhashStr string) (*domain.Image, int, bool, error) {
	return m.matchFunc(dhashStr)
}

func (m *mockImageService) CreateImage(ctx context.Context, img *domain.Image) error {
	return nil
}

func TestMatchImage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := logrus.New()
	logger.SetOutput(bytes.NewBuffer(nil)) // suppress logs during test

	t.Run("Success Matched", func(t *testing.T) {
		mockSvc := &mockImageService{
			matchFunc: func(dhashStr string) (*domain.Image, int, bool, error) {
				return &domain.Image{ID: 1, URL: "http://example.com/1"}, 5, true, nil
			},
		}
		h := NewMazeHandler(mockSvc, logger)
		r := gin.Default()
		r.POST("/match-image", h.MatchImage)

		reqBody := models.MatchMazeImageRequest{DHash: "1234567812345678"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/match-image", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp models.MazeMatchResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if !resp.Matched {
			t.Error("expected matched to be true")
		}
		if resp.Maze == nil || resp.Maze.ID != 1 {
			t.Errorf("expected maze id 1, got %v", resp.Maze)
		}
		if resp.Distance != 5 {
			t.Errorf("expected distance 5, got %d", resp.Distance)
		}
	})

	t.Run("Success Not Matched", func(t *testing.T) {
		mockSvc := &mockImageService{
			matchFunc: func(dhashStr string) (*domain.Image, int, bool, error) {
				return nil, 0, false, nil
			},
		}
		h := NewMazeHandler(mockSvc, logger)
		r := gin.Default()
		r.POST("/match-image", h.MatchImage)

		reqBody := models.MatchMazeImageRequest{DHash: "1234567812345678"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/match-image", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp models.MazeMatchResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Matched {
			t.Error("expected matched to be false")
		}
	})

	t.Run("Error Internal", func(t *testing.T) {
		mockSvc := &mockImageService{
			matchFunc: func(dhashStr string) (*domain.Image, int, bool, error) {
				return nil, 0, false, errors.New("db error")
			},
		}
		h := NewMazeHandler(mockSvc, logger)
		r := gin.Default()
		r.POST("/match-image", h.MatchImage)

		reqBody := models.MatchMazeImageRequest{DHash: "1234567812345678"}
		body, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/match-image", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}
