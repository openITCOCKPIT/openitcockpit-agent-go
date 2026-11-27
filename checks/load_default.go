//go:build !aix
// +build !aix

package checks

import (
	"context"

	"github.com/shirou/gopsutil/v4/load"
)

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckLoad) Run(ctx context.Context) (interface{}, error) {
	l, err := load.AvgWithContext(ctx)

	if err != nil {
		return nil, err
	}
	return &resultLoad{
		Load1:  l.Load1,
		Load5:  l.Load5,
		Load15: l.Load15,
	}, nil
}
