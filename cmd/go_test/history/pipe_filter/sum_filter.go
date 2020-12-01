package pipe_filter

import "errors"

var SumFilterWong = errors.New("input mistake")

type sumFilter struct {
}

func NewSumFilter() *sumFilter {
	return &sumFilter{}
}

func (sf *sumFilter) Process(data Request) (Response, error) {
	elems, ok := data.([]int)
	if !ok {
		return nil, SumFilterWong
	}
	ret := 0
	for _, v := range elems {
		ret += v
	}
	return ret, nil
}
