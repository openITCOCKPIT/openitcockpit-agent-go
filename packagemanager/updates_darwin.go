package packagemanager

import (
	"context"
	"fmt"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/utils"
)

type MacOSUpdatesManager struct{}

func (m MacOSUpdatesManager) getAvailableUpdates(ctx context.Context) (string, error) {
	timeout := time.Duration(120 * time.Second)
	result, err := utils.RunCommand(ctx, utils.CommandArgs{
		Command: "softwareupdate --list --no-scan",
		Timeout: timeout,
	})

	if err != nil {
		return "", err
	}

	if result.RC != 0 {
		return "", fmt.Errorf("error fetching macOS Updates via softwareupdate, exit code %d: %s", result.RC, result.Stdout)
	}

	return result.Stdout, nil
}

func (m MacOSUpdatesManager) ListAvailableUpdates(ctx context.Context) ([]MacosUpdate, error) {
	output, err := m.getAvailableUpdates(ctx)
	if err != nil {
		return nil, err
	}

	return parseSoftwareUpdateOutput(output)
}
