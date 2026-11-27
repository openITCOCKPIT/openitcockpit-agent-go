package checks

import (
	"context"
	"fmt"
)

// networkPortsPlatform reads the sockets with netstat. The processes of the listening sockets are not reported,
// as AIX only reports them with rmsock for each socket (root required).
type networkPortsPlatform struct{}

// sockets returns the TCP and UDP sockets of IPv4 and IPv6 (netstat -an, the format is the one of macOS).
// "tcp" and "udp" without 4 or 6 are dual stack sockets. Sockets without a local port are not bound yet.
func (p *networkPortsPlatform) sockets(ctx context.Context) ([]*networkSocket, error) {
	output, err := runAixCommand(ctx, "netstat", "-an", "-f", "inet")
	if err != nil {
		return nil, err
	}
	var sockets []*networkSocket
	for _, socket := range parseDarwinNetstat([]byte(output)) {
		if socket.LocalPort != 0 {
			sockets = append(sockets, socket)
		}
	}
	return sockets, nil
}

// ephemeralPortRange returns the range of the ports for outgoing TCP connections (no -o tcp_ephemeral_low and high)
func (p *networkPortsPlatform) ephemeralPortRange(ctx context.Context) (int, int, error) {
	first, err := aixNetworkOption(ctx, "tcp_ephemeral_low")
	if err != nil {
		return 0, 0, err
	}
	last, err := aixNetworkOption(ctx, "tcp_ephemeral_high")
	if err != nil {
		return 0, 0, err
	}
	return first, last, nil
}

func (p *networkPortsPlatform) tcpErrors() *resultTcpErrors {
	return nil
}

// aixNetworkOption reads a network option with "no -o", which only displays the value
func aixNetworkOption(ctx context.Context, option string) (int, error) {
	output, err := runAixCommand(ctx, "no", "-o", option)
	if err != nil {
		return 0, err
	}
	value, ok := parseAixNetworkOption(output)
	if !ok {
		return 0, fmt.Errorf("no -o %s: unexpected output %q", option, output)
	}
	return value, nil
}
