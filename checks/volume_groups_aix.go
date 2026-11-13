package checks

import (
	"context"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// CheckVolumeGroups reports the AIX LVM volume groups with their logical and physical volumes
type CheckVolumeGroups struct {
}

// Name will be used in the response as check name
func (c *CheckVolumeGroups) Name() string {
	return "volume_groups"
}

// Configure the command or return false if the command was disabled
func (c *CheckVolumeGroups) Configure(config *config.Configuration) (bool, error) {
	return config.Lvm && aixCommandAvailable("lsvg"), nil
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckVolumeGroups) Run(ctx context.Context) (interface{}, error) {
	allOutput, err := runAixCommand(ctx, "lsvg")
	if err != nil {
		return nil, err
	}
	onlineOutput, err := runAixCommand(ctx, "lsvg", "-o")
	if err != nil {
		return nil, err
	}

	online := map[string]bool{}
	for _, name := range parseLsvgList(onlineOutput) {
		online[name] = true
	}

	vgs := []*resultVolumeGroup{}
	for _, name := range parseLsvgList(allOutput) {
		vg := newVolumeGroup(name, online[name])
		if vg.Online {
			readVolumeGroup(ctx, vg)
		}
		vgs = append(vgs, vg)
	}
	return vgs, nil
}

// newVolumeGroup returns a volume group with empty lists, offline volume groups (e.g. the volume groups of
// the other node of a cluster) have no details
func newVolumeGroup(name string, online bool) *resultVolumeGroup {
	return &resultVolumeGroup{
		Name:       name,
		Online:     online,
		LVs:        []*resultLogicalVolume{},
		PVs:        []*resultPhysicalVolume{},
		ReadErrors: []string{},
	}
}

// readVolumeGroup reads the details of an online volume group. Errors are reported per volume group,
// e.g. if the volume group was varied off in the meantime.
func readVolumeGroup(ctx context.Context, vg *resultVolumeGroup) {
	output, err := runAixCommand(ctx, "lsvg", vg.Name)
	if err != nil {
		vg.ReadErrors = append(vg.ReadErrors, err.Error())
	} else {
		parseLsvg(output, vg)
	}

	output, err = runAixCommand(ctx, "lsvg", "-l", vg.Name)
	if err != nil {
		vg.ReadErrors = append(vg.ReadErrors, err.Error())
	} else {
		vg.LVs = parseLsvgLogicalVolumes(output)
	}

	output, err = runAixCommand(ctx, "lsvg", "-p", vg.Name)
	if err != nil {
		vg.ReadErrors = append(vg.ReadErrors, err.Error())
	} else {
		vg.PVs = parseLsvgPhysicalVolumes(output)
	}
}
