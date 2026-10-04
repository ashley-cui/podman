//go:build freebsd

package netavark

import (
	"go.podman.io/common/libnetwork/types"
)

// applyPerIFSysctls is a no-op on FreeBSD (no /proc/sys/).
func applyPerIFSysctls(_ string, _ []types.NamedPerNetworkOptions) error {
	return nil
}
