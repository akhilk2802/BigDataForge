package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"BigDataForge/internal/models"
	"BigDataForge/internal/routes"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Redis Client for integration testing
type MockRedisClient struct {
	mock.Mock
	redis.Cmdable // Embed the redis.Cmdable interface
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return &redis.StatusCmd{}
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	result := new(redis.StringCmd)
	result.SetVal("mocked_value")
	return result
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	return &redis.IntCmd{}
}

// Integration test for POST /plans
func TestIntegrationCreatePlan(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockRedis := new(MockRedisClient)

	// Pass mockRedis to SetupRoutes (now it implements Cmdable)
	routes.SetupRoutes(r, mockRedis)

	plan := models.Plan{
		ObjectID: "plan123",
		PlanType: "inNetwork",
	}
	planJSON, _ := json.Marshal(plan)

	// Mock Redis Set command
	mockRedis.On("Set", ctx, "plan:plan123", planJSON, mock.Anything).Return(nil)

	// Create HTTP request
	req, _ := http.NewRequest("POST", "/api/v1/plans", bytes.NewBuffer(planJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check if the response is HTTP 201 Created
	assert.Equal(t, http.StatusCreated, w.Code)
}

// Integration test for GET /plans?id=<id>
func TestIntegrationGetPlan(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	mockRedis := new(MockRedisClient)

	// Pass mockRedis to SetupRoutes (now it implements Cmdable)
	routes.SetupRoutes(r, mockRedis)

	planID := "plan123"
	plan := models.Plan{
		ObjectID: "plan123",
		PlanType: "inNetwork",
	}
	planJSON, _ := json.Marshal(plan)

	// Mock Redis Get command
	mockRedis.On("Get", ctx, "plan:"+planID).Return(string(planJSON), nil)

	req, _ := http.NewRequest("GET", "/api/v1/plans?id=plan123", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check if the response is HTTP 200 OK
	assert.Equal(t, http.StatusOK, w.Code)
}
