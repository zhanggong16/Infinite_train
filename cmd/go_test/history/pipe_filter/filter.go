package pipe_filter

// request is the input of the filter
type Request interface {
}

// response is the output of the filter
type Response interface {
}

// interface define
type Filter interface {
	Process(data Request) (Response, error)
}
