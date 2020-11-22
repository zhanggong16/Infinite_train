package pipe_match

type Request interface {
}

type Response interface {
}

type Match interface {
	Process(data Request) (Response, error)
}
