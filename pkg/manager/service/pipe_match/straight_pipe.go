package pipe_match

type StraightPipeLine struct {
	Name 	string
	Match *[]Match
}

func NewStraightPipeLine(name string, match ...Match) *StraightPipeLine {
	return &StraightPipeLine{
		Name: name,
		Match: &match,
	}
}

func (f *StraightPipeLine) Process(data Request) (Response, error) {
	var ret interface{}
	var err error
	for _, match := range *f.Match {
		ret, err = match.Process(data)
		if err != nil {
			return ret, err
		}
		data = ret
	}
	return ret, err
}
