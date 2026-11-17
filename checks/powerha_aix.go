package checks

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/openITCOCKPIT/openitcockpit-agent-go/config"
)

// powerhaUtilities contains the commands of PowerHA SystemMirror, the directory is not in $PATH
const powerhaUtilities = "/usr/es/sbin/cluster/utilities/"

// powerhaClusterInterval is the interval of clmgr query cluster. clmgr is a Perl script which takes about 1.5 seconds
// on a small LPAR, the attributes (name, version, unsynchronized changes) rarely change.
const powerhaClusterInterval = 10 * time.Minute

// CheckPowerHA reports the state of the PowerHA SystemMirror cluster: the cluster manager of the node, the nodes,
// the resource groups with their applications and the repository disks.
// The state of the cluster manager and the repository disks are readable without root, clRGinfo, clshowsrv and
// clmgr require root. Commands which fail are reported in the errors of the result.
// The attributes of the cluster (clmgr) are read every 10 minutes, so unsynchronized changes are reported with a delay.
type CheckPowerHA struct {
	clusterAttributes map[string]string // Last result of clmgr query cluster
	clusterQueried    time.Time
}

// Name will be used in the response as check name
func (c *CheckPowerHA) Name() string {
	return "powerha"
}

// Configure the command or return false if the command was disabled
func (c *CheckPowerHA) Configure(config *config.Configuration) (bool, error) {
	if !config.PowerHA {
		return false, nil
	}
	// Only enable the check if PowerHA SystemMirror is installed
	_, err := os.Stat(powerhaUtilities + "clRGinfo")
	return err == nil, nil
}

// Run the actual check
// if error != nil the check result will be nil
// ctx can be canceled and runs the timeout
// CheckResult will be serialized after the return and should not change until the next call to Run
func (c *CheckPowerHA) Run(ctx context.Context) (interface{}, error) {
	result := &resultPowerHA{
		Nodes:           []*resultPowerHANode{},
		ResourceGroups:  []*resultPowerHAResourceGroup{},
		RepositoryDisks: []*resultPowerHARepositoryDisk{},
		Errors:          []string{},
	}

	output, err := runAixCommand(ctx, "lssrc", "-ls", "clstrmgrES")
	if err != nil {
		result.addError("lssrc -ls clstrmgrES", err)
	} else {
		result.ClusterManagerState = parseClstrmgrState(output)
	}

	output, err = runAixCommand(ctx, "lscluster", "-d")
	if err != nil {
		result.addError("lscluster -d", err)
	} else {
		result.RepositoryDisks = parseLsclusterDisks(output)
	}

	// clRGinfo fails if the cluster services of the node are not running (Cluster IPC error ... ST_INIT)
	output, err = runAixCommand(ctx, powerhaUtilities+"clRGinfo", "-s")
	if err != nil {
		result.addError("clRGinfo -s", err)
	} else {
		result.ResourceGroups = parseClRGinfo(output)
		c.readApplications(ctx, result)
	}

	output, err = runAixCommand(ctx, powerhaUtilities+"clshowsrv", "-v")
	if err != nil {
		result.addError("clshowsrv -v", err)
	} else {
		result.Nodes = parseClshowsrv(output)
	}

	c.readClusterAttributes(ctx, result)
	return result, nil
}

// readClusterAttributes sets the attributes of clmgr query cluster, which is only run every powerhaClusterInterval
func (c *CheckPowerHA) readClusterAttributes(ctx context.Context, result *resultPowerHA) {
	if c.clusterAttributesOutdated() {
		output, err := runAixCommand(ctx, powerhaUtilities+"clmgr", "query", "cluster")
		if err != nil {
			// The next run tries it again
			c.clusterAttributes = nil
			result.addError("clmgr query cluster", err)
			return
		}
		c.clusterAttributes = parseClmgrQuery(output)
		c.clusterQueried = time.Now()
	}
	setPowerHAClusterAttributes(result, c.clusterAttributes)
}

func (c *CheckPowerHA) clusterAttributesOutdated() bool {
	return c.clusterAttributes == nil || time.Since(c.clusterQueried) >= powerhaClusterInterval
}

// readApplications adds the application controllers to the resource groups (clRGinfo -m)
func (c *CheckPowerHA) readApplications(ctx context.Context, result *resultPowerHA) {
	output, err := runAixCommand(ctx, powerhaUtilities+"clRGinfo", "-m")
	if err != nil {
		result.addError("clRGinfo -m", err)
		return
	}
	parseClRGinfoApplications(output, result.ResourceGroups)
}

func (r *resultPowerHA) addError(command string, err error) {
	r.Errors = append(r.Errors, fmt.Sprintf("%s: %s", command, err))
}
