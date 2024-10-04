package unit

import (
	"BigDataForge/internal/models"
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

// Mock Redis Client
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

// Test Create Plan
func TestCreatePlan(t *testing.T) {
	mockRedis := new(MockRedisClient)
	service := services.NewPlanService(mockRedis)

	plan := models.Plan{
		ObjectID: "plan123",
		PlanType: "inNetwork",
	}

	planJSON, _ := json.Marshal(plan)
	mockRedis.On("Set", ctx, "plan:plan123", planJSON, mock.Anything).Return(nil)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	req, _ := http.NewRequest("POST", "/api/v1/plans", bytes.NewBuffer(planJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	service.CreatePlan(c)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// Test Get Plan
func TestGetPlan(t *testing.T) {
	mockRedis := new(MockRedisClient)
	service := services.NewPlanService(mockRedis)

	planID := "plan123"
	plan := models.Plan{
		ObjectID: "plan123",
		PlanType: "inNetwork",
	}

	planJSON, _ := json.Marshal(plan)

	mockRedis.On("Get", ctx, "plan:"+planID).Return(string(planJSON), nil)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	req, _ := http.NewRequest("GET", "/api/v1/plans?id=plan123", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	service.GetPlan(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test Plan Not Found
func TestGetPlanNotFound(t *testing.T) {
	mockRedis := new(MockRedisClient)
	service := services.NewPlanService(mockRedis)

	planID := "nonexistent"

	// Set mock expectation for missing plan
	mockRedis.On("Get", ctx, "plan:"+planID).Return("", redis.Nil)

	gin.SetMode(gin.TestMode)
	r := gin.Default()

	req, _ := http.NewRequest("GET", "/api/v1/plans?id=nonexistent", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	service.GetPlan(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
