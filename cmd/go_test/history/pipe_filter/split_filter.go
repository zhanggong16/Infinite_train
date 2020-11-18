package pipe_filter

import (
	"errors"
	"strings"
)

type SplitFilter struct {
	delimiter string
}

func NewSplitFilter(delimiter string) *SplitFilter {
	return &SplitFilter{delimiter}
}

var SplitFilterWong = errors.New("input mistake")

func (sf *SplitFilter) Process(data Request) (Response, error) {
	str, ok := data.(string)
	if !ok {
		return nil, SplitFilterWong
	}
	parts := strings.Split(str, sf.delimiter)
	return parts, nil
}

