package view

type CommonSUCView struct {
	RequestID string      `json:"requestId"`
	Result    interface{} `json:"result"`
}

type CommonGidView struct {
	GID	string	`json:"gid"`
}