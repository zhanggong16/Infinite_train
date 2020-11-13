package request

type InstanceCreateRequestBody struct {
	GID		string	`json:"gid" validate:"required"`
	Name	string	`json:"name" validate:"required,InstanceName"`
}

//ManagerActionBody ...
type InstanceActionBody struct {
	Actions *InstanceAction `json:"actions" validate:"required"`
}

//ManagerAction ...
type InstanceAction struct {
	Method string      		`json:"method" validate:"required"`
	Params interface{} 		`json:"params" validate:"omitempty"`
}

//ManagerCommonContext is used to delivery context
type ManagerCommonContext struct {
	ID            string
	CommonContext *CommonContext
}

//ChangeInstanceNameParam ...
type ChangeClusterNameParam struct {
	NewName string `json:"new_name" validate:"required,InstanceName"`
}