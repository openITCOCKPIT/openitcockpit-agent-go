package checks

import (
	"context"
	"net"
)

// On FreeBSD, Go's net package provides access to network interface details.
// However, advanced properties like link speed or duplex mode are not available via net.Interfaces().
func (c *CheckNet) Run(ctx context.Context) (interface{}, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	results := make(map[string]*resultNet)
	for _, iface := range interfaces {
		active := (iface.Flags & net.FlagUp) != 0

		// FreeBSD: Only basic properties (status, MTU) are accessible here.
		results[iface.Name] = &resultNet{
			Isup: active,
			MTU:  int64(iface.MTU),
		}
	}

	return results, nil
}
