package integration

import (
	"BigDataForge/internal/models"
	"BigDataForge/internal/routes"
	"BigDataForge/internal/services"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Redis Client for integration testing
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return &redis.StatusCmd{}
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	result := new(redis.StringCmd)
	if v := args.String(0); v != "" {
		result.SetVal(v)
	} else {
		result.SetErr(args.Error(1))
	}
	return result
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	return new(redis.IntCmd)
}

// Integration test for POST /plans
func TestIntegrationCreatePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockRedis := new(MockRedisClient)
	services := services.NewPlanService(mockRedis)

	routes.SetupRoutes(r, mockRedis)

	plan := models.Plan{
		ObjectID: "plan123",
		PlanType: "inNetwork",
	}
	planJSON, _ := json.Marshal(plan)

	mockRedis.On("Set", ctx, "plan:plan123", planJSON, mock.Anything).Return(nil)

	// Create HTTP request
	req, _ := http.NewRequest("POST", "/api/v1/plans", bytes.NewBuffer(planJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// Integration test for GET /plans?id=<id>
func TestIntegrationGetPlan(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockRedis := new(MockRedisClient)
	services := services.NewPlanService(mockRedis)

	routes.SetupRoutes(r, mockRedis)

	planID := "plan123"
	plan := models.Plan{
		ObjectID: "plan123",
		PlanType: "inNetwork",
	}
	planJSON, _ := json.Marshal(plan)

	mockRedis.On("Get", ctx, "plan:"+planID).Return(string(planJSON), nil)

	req, _ := http.NewRequest("GET", "/api/v1/plans?id=plan123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
