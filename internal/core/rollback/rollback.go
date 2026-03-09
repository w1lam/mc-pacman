// Package rollback holds the rollback handler
package rollback

import "fmt"

type Rollback struct {
	steps []func() error
}

// Add adds a step to the rollback
func (r *Rollback) Add(fn func() error) {
	r.steps = append(r.steps, fn)
}

// Run runs the rollback steps
func (r *Rollback) Run() error {
	var errs []error
	for i := len(r.steps) - 1; i >= 0; i-- {
		if err := r.steps[i](); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("rollback errors: %v", errs)
	}
	return nil
}
