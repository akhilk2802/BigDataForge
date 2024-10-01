package models

type PlanCostShares struct {
	Deductible int    `json:"deductible"`
	Org        string `json:"_org"`
	Copay      int    `json:"copay"`
	ObjectId   string `json:"objectId"`
	ObjectType string `json:"objectType"`
}

type LinkedService struct {
	Org        string `json:"_org"`
	ObjectId   string `json:"objectId"`
	ObjectType string `json:"objectType"`
	Name       string `json:"name"`
}

type PlanServiceCostShares struct {
	Deductible int    `json:"deductible"`
	Org        string `json:"_org"`
	Copay      int    `json:"copay"`
	ObjectId   string `json:"objectId"`
	ObjectType string `json:"objectType"`
}

type LinkedPlanServices struct {
	LinkedService         LinkedService         `json:"linkedService"`
	PlanServiceCostShares PlanServiceCostShares `json:"planserviceCostShares"`
	Org                   string                `json:"_org"`
	ObjectId              string                `json:"objectId"`
	ObjectType            string                `json:"objectType"`
}

type Plan struct {
	PlanCostShares     PlanCostShares       `json:"planCostShares"`
	LinkedPlanServices []LinkedPlanServices `json:"linkedPlanServices"`
	Org                string               `json:"_org"`
	ObjectId           string               `json:"objectId"`
	ObjectType         string               `json:"objectType"`
	PlanType           string               `json:"planType"`
	CreationDate       string               `json:"creationDate"`
}
