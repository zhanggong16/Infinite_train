package pipe_filter

import "testing"

func TestNewSplitFilter(t *testing.T) {
	spliter := NewSplitFilter(",")
	coverter := NewToIntFilter()
	sum := NewSumFilter()

	sp := NewStraightPipeLine("p1", spliter, coverter, sum)
	ret, err := sp.Process("1,2,3")
	if err != nil {
		t.Fatal(err)
	}
	if ret != 6 {
		t.Fatalf("The expected is 6, but the actual is %d", ret)
	} else {
		t.Log("finish")
	}
}