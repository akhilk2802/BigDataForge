package validators

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xeipuuv/gojsonschema"
)

const planSchema = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "planCostShares": {
      "type": "object",
      "properties": {
        "deductible": { "type": "integer" },
        "copay": { "type": "integer" },
        "objectId": { "type": "string" },
        "objectType": { "type": "string" }
      },
      "required": ["deductible", "copay", "objectId", "objectType"]
    }
  },
  "required": ["planCostShares"]
}
`

func ValidatePlanSchema(c *gin.Context) bool {
	jsonData, _ := c.GetRawData()
	schemaLoader := gojsonschema.NewStringLoader(planSchema)
	documentLoader := gojsonschema.NewStringLoader(string(jsonData))

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil || !result.Valid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return false
	}
	return true
}
