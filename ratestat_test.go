package ratestat

import (
	"testing"
	"time"
)

func TestRateStat(t *testing.T) {
	rs := New(5, time.Second*2)
	time.Sleep(1 * time.Second)
	rs.Log(1)
	time.Sleep(2 * time.Second)
	rs.Log(2)
	if rs.Value() != 3 {
		t.Error("RS != 3", rs.Value())
	}
	time.Sleep(8 * time.Second)
	if rs.Value() != 2 {
		t.Error("RS != 2", rs.Value())
	}
	time.Sleep(2 * time.Second)
	if rs.Value() != 0 {
		t.Error("RS != 0", rs.Value())
	}
}
