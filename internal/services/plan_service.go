package services

import (
	"encoding/json"
	"BigDataForge/internal/models"
	"BigDataForge/internal/storage"
	"time"
)

type PlanService struct {
	Redis *storage.RedisStore
}

func NewPlanService(redisStore *storage.RedisStore) *PlanService {
	return &PlanService{Redis: redisStore}
}

func (ps *PlanService) CreatePlan(id string, plan models.Plan) error {
	planData, err := json.Marshal(plan)
	if err != nil {
		return err
	}
	return ps.Redis.SetPlan(id, string(planData))
}

func (ps *PlanService) GetPlan(id string, ifModifiedSince time.Time) (models.Plan, bool, error) {
	data, err := ps.Redis.GetPlan(id)
	if err != nil {
		return models.Plan{}, false, err
	}
	var plan models.Plan
	json.Unmarshal([]byte(data), &plan)

	creationDate, err := time.Parse("02-01-2006", plan.CreationDate)
	if err != nil {
		return models.Plan{}, false, err
	}
	if !ifModifiedSince.IsZero() && !creationDate.After(ifModifiedSince) {
		return plan, false, nil // Not modified
	}
	return plan, true, nil
}

func (ps *PlanService) DeletePlan(id string) error {
	return ps.Redis.DeletePlan(id)
}
