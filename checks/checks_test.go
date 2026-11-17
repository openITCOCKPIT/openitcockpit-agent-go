package checks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/moby/moby/client"
)

func TestChecksWithDefault(t *testing.T) {
	checks := getPlatformChecks()

	for _, c := range checks {
		if c.Name() == "" {
			t.Error("Invalid name")
		}
		r, err := c.Run(context.Background())

		if err != nil {
			// Docker is not installed, not running or the socket is not accessible (permission denied)
			if client.IsErrConnectionFailed(err) {
				fmt.Println("Docker is not available on this system: ", err)
				continue
			}

			// e.g. the Ceph check is only enabled by Configure if the ceph command is installed
			if strings.Contains(err.Error(), "executable file not found") {
				fmt.Println(err)
				continue
			}

			// Commands like lvs require root, which is not given if the tests are not executed as root
			if strings.Contains(err.Error(), "non-root") || strings.Contains(err.Error(), "Permission denied") {
				fmt.Println(err)
				continue
			}

			if strings.Contains(err.Error(), "/var/run/utmp: no such file or directory") {
				// https://github.com/shirou/gopsutil/issues/900
				fmt.Println(err)
				fmt.Println("The golang docker image does not have this directory. Are you using the golang docker image?")
				continue
			}
		}

		if err != nil {
			t.Errorf("Test of check %s failed with error: %s", c.Name(), err)
			continue
		}
		if r == nil {
			t.Errorf("Test of check %s returned nil", c.Name())
			continue
		}
		js, err := json.Marshal(r)
		if err != nil {
			t.Errorf("Test of check %s returned result that can't be marshaled: %s", c.Name(), err)
			continue
		}
		fmt.Println(string(js))
	}
}
