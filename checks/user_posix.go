//go:build linux || darwin || freebsd
// +build linux darwin freebsd

package checks

import (
	"context"
	"runtime"

	"github.com/shirou/gopsutil/v4/host"
)

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckUser) Run(ctx context.Context) (interface{}, error) {
	users, err := host.UsersWithContext(ctx)
	if err != nil {
		return nil, err
	}

	userResults := make([]*resultUser, 0, len(users))

	for _, user := range users {
		var result *resultUser
		if runtime.GOOS == "darwin" {
			result = &resultUser{
				Name:     user.User,
				Terminal: user.Terminal,
				Host:     user.Host,
				//https://github.com/shirou/gopsutil/issues/1989
				Started: 0,
			}
		} else {
			result = &resultUser{
				Name:     user.User,
				Terminal: user.Terminal,
				Host:     user.Host,
				Started:  int64(user.Started),
			}

		}
		userResults = append(userResults, result)
	}
	return userResults, nil
}
