package checks

import (
	"context"
	"errors"
)

// networkPortsPlatform is not implemented on AIX yet, the check is not enabled on AIX
type networkPortsPlatform struct{}

var errNetworkPortsNotImplemented = errors.New("network ports are not implemented on AIX")

func (p *networkPortsPlatform) sockets(_ context.Context) ([]*networkSocket, error) {
	return nil, errNetworkPortsNotImplemented
}

func (p *networkPortsPlatform) ephemeralPortRange(_ context.Context) (int, int, error) {
	return 0, 0, errNetworkPortsNotImplemented
}

func (p *networkPortsPlatform) tcpErrors() *resultTcpErrors {
	return nil
}
