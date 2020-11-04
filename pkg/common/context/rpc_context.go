package context

type PingPongRequest struct {
	RequestID	string
}

type ReportHeartBeatRequest struct {
	RequestID	string
	AgentIP		string
}