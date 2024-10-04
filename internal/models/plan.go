package models

type PlanCostShares struct {
	Deductible int    `json:"deductible"`
	Org        string `json:"_org"`
	Copay      int    `json:"copay"`
	ObjectID   string `json:"objectId"`
	ObjectType string `json:"objectType"`
}

type LinkedService struct {
	Org        string `json:"_org"`
	ObjectID   string `json:"objectId"`
	ObjectType string `json:"objectType"`
	Name       string `json:"name"`
}

type PlanserviceCostShares struct {
	Deductible int    `json:"deductible"`
	Org        string `json:"_org"`
	Copay      int    `json:"copay"`
	ObjectID   string `json:"objectId"`
	ObjectType string `json:"objectType"`
}

type LinkedPlanService struct {
	LinkedService         LinkedService         `json:"linkedService"`
	PlanserviceCostShares PlanserviceCostShares `json:"planserviceCostShares"`
	Org                   string                `json:"_org"`
	ObjectID              string                `json:"objectId"`
	ObjectType            string                `json:"objectType"`
}

type Plan struct {
	PlanCostShares     PlanCostShares      `json:"planCostShares"`
	LinkedPlanServices []LinkedPlanService `json:"linkedPlanServices"`
	Org                string              `json:"_org"`
	ObjectID           string              `json:"objectId"`
	ObjectType         string              `json:"objectType"`
	PlanType           string              `json:"planType"`
	CreationDate       string              `json:"creationDate"`
}
