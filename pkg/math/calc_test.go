package math

import "testing"

func TestAvg(t *testing.T) {
	v := Avg(10, 20)
	if v != 15 {
		t.Error("Expected 15, got", v)
	}
}
