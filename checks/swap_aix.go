package checks

import (
	"context"
)

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckSwap) Run(_ context.Context) (interface{}, error) {
	buf, err := perfstatMemoryTotal()
	if err != nil {
		return nil, err
	}
	memory, err := decodeMemoryTotal(buf)
	if err != nil {
		return nil, err
	}
	return aixSwapResult(memory), nil
}
