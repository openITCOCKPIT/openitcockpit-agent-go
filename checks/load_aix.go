package checks

import (
	"context"
)

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckLoad) Run(_ context.Context) (interface{}, error) {
	buf, err := perfstatCpuTotal()
	if err != nil {
		return nil, err
	}
	load, err := decodeLoadAverage(buf)
	if err != nil {
		return nil, err
	}
	return &resultLoad{
		Load1:  load[0],
		Load5:  load[1],
		Load15: load[2],
	}, nil
}
