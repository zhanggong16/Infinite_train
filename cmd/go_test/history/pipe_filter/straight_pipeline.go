package pipe_filter

type StraightPipeLine struct {
	Name 	string
	Filters *[]Filter
}

func NewStraightPipeLine(name string, filters ...Filter) *StraightPipeLine {
	return &StraightPipeLine{
		Name: name,
		Filters: &filters,
	}
}

func (f *StraightPipeLine) Process(data Request) (Response, error) {
	var ret interface{}
	var err error
	for _, filter := range *f.Filters {
		ret, err = filter.Process(data)
		if err != nil {
			return ret, err
		}
		data = ret
	}
	return ret, err
}