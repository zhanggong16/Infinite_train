package pipe_filter

import "testing"

func TestNewSplitFilter(t *testing.T) {
	spliter := NewSplitFilter(",")
	sp := NewStraightPipeLine("p1", spliter)
	ret, err := sp.Process("1,2,3")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ret)
	/*if ret != 6 {
		t.Fatalf("The expected is 6, but the actual is %d", ret)
	}*/
}