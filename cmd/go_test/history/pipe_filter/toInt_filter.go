package pipe_filter

import (
	"errors"
	"strconv"
)

var ToIntFilterWong = errors.New("input mistake")

type toIntFilter struct {
}

func NewToIntFilter() *toIntFilter {
	return &toIntFilter{}
}

func (tlf *toIntFilter) Process(data Request) (Response, error) {
	parts, ok := data.([]string)
	if !ok {
		return nil, ToIntFilterWong
	}
	ret := []int{}
	for _, part := range parts {
		s, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		ret = append(ret, s)
	}
	return ret, nil
}