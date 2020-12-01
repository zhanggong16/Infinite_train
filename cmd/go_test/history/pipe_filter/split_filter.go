package pipe_filter

import (
	"errors"
	"strings"
)

var SplitFilterWong = errors.New("input mistake")

type SplitFilter struct {
	delimiter string
}

func NewSplitFilter(delimiter string) *SplitFilter {
	return &SplitFilter{delimiter}
}

func (sf *SplitFilter) Process(data Request) (Response, error) {
	// 第一步数据验证
	str, ok := data.(string)
	if !ok {
		return nil, SplitFilterWong
	}
	parts := strings.Split(str, sf.delimiter)
	return parts, nil
}

