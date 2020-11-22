package pipe_match

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
	// 入参类型转换
	str, ok := data.(string)
	if !ok {
		return nil, SplitFilterWong
	}
	// 处理
	parts := strings.Split(str, sf.delimiter)
	return parts, nil
}

