package controllers

import (
	models "BigDataForge/internal/models"
	services "BigDataForge/internal/services"
	validators "BigDataForge/internal/validators"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
)

func CreatePlan(service *services.PlanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !validators.ValidatePlanSchema(c) {
			return
		}
		var plan models.Plan
		if err := c.BindJSON(&plan); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
			return
		}

		if err := service.CreatePlan(plan.ObjectId, plan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plan"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Plan created"})
	}
}

func GetPlan(service *services.PlanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		ifModifiedSince := c.GetHeader("If-Modified-Since")
		var lastModified time.Time
		if ifModifiedSince != "" {
			lastModified, _ = time.Parse(time.RFC1123, ifModifiedSince)
		}

		plan, modified, err := service.GetPlan(id, lastModified)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Plan not found"})
			return
		}
		if !modified {
			c.Status(http.StatusNotModified)
			return
		}
		c.JSON(http.StatusOK, gin.H{"plan": plan})
	}
}

func DeletePlan(service *services.PlanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := service.DeletePlan(id); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Plan not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Plan deleted"})
	}
}
