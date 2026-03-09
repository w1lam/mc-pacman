package rollback

import (
	"errors"
	"slices"
	"testing"
)

func TestRollback_RunsInReverse(t *testing.T) {
	var order []int
	var rb Rollback

	rb.Add(func() error { order = append(order, 1); return nil })
	rb.Add(func() error { order = append(order, 2); return nil })
	rb.Add(func() error { order = append(order, 3); return nil })

	rb.Run()

	want := []int{3, 2, 1}
	if !slices.Equal(order, want) {
		t.Errorf("got %v want %v", order, want)
	}
}

func TestRollback_ContinuesOnError(t *testing.T) {
	ran := 0
	var rb Rollback

	rb.Add(func() error { ran++; return nil })
	rb.Add(func() error { return errors.New("step failed") })
	rb.Add(func() error { ran++; return nil })

	err := rb.Run()

	if ran != 2 {
		t.Errorf("expected 2 steps to run, got %d", ran)
	}

	if err == nil {
		t.Error("exåected error, got nil")
	}
}

func TestRollback_EmptyDoesNothing(t *testing.T) {
	var rb Rollback
	if err := rb.Run(); err != nil {
		t.Errorf("empty rollback should return nil, got %v", err)
	}
}
