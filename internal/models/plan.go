package models

type SearchRequest struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PlanCostShares struct {
	PlanJoin   map[string]interface{} `json:"plan_join,omitempty"`
	Deductible int                    `json:"deductible"`
	Org        string                 `json:"_org"`
	Copay      int                    `json:"copay"`
	ObjectID   string                 `json:"objectId"`
	ObjectType string                 `json:"objectType"`
}

type LinkedService struct {
	PlanJoin   map[string]interface{} `json:"plan_join,omitempty"`
	Org        string                 `json:"_org"`
	ObjectID   string                 `json:"objectId"`
	ObjectType string                 `json:"objectType"`
	Name       string                 `json:"name"`
}

type PlanserviceCostShares struct {
	PlanJoin   map[string]interface{} `json:"plan_join,omitempty"`
	Deductible int                    `json:"deductible"`
	Org        string                 `json:"_org"`
	Copay      int                    `json:"copay"`
	ObjectID   string                 `json:"objectId"`
	ObjectType string                 `json:"objectType"`
}

type LinkedPlanService struct {
	PlanJoin              map[string]interface{} `json:"plan_join,omitempty"`
	LinkedService         LinkedService          `json:"linkedService"`
	PlanserviceCostShares PlanserviceCostShares  `json:"planserviceCostShares"`
	Org                   string                 `json:"_org"`
	ObjectID              string                 `json:"objectId"`
	ObjectType            string                 `json:"objectType"`
}

type Plan struct {
	PlanJoin           map[string]interface{} `json:"plan_join,omitempty"`
	PlanCostShares     PlanCostShares         `json:"planCostShares"`
	LinkedPlanServices []LinkedPlanService    `json:"linkedPlanServices"`
	Org                string                 `json:"_org"`
	ObjectID           string                 `json:"objectId"`
	ObjectType         string                 `json:"objectType"`
	PlanType           string                 `json:"planType"`
	CreationDate       string                 `json:"creationDate"`
}
