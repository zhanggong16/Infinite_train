package view

type CommonSUCView struct {
	RequestId string      `json:"requestId"`
	Result    interface{} `json:"result"`
}

type CommonGidView struct {
	Gid	string	`json:"gid"`
}