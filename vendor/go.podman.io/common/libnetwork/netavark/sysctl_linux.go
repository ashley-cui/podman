//go:build linux

package netavark

import (
	"fmt"
	"os"
	"strings"

	"go.podman.io/common/libnetwork/types"
	"go.podman.io/common/pkg/netns"
)

// ifSysctl represents a single sysctl to apply to an interface.
type ifSysctl struct {
	key   string
	value string
}

// applyPerIFSysctls applies per-interface sysctls from PerNetworkOptions.Options.
// It looks for the types.PerIFSysctl key in each network's Options map and applies
// the sysctls inside the container's network namespace.
func applyPerIFSysctls(namespacePath string, networks []types.NamedPerNetworkOptions) error {
	for _, network := range networks {
		sysctlVal, ok := network.Options[types.PerIFSysctl]
		if !ok || sysctlVal == "" {
			continue
		}

		ifName := network.InterfaceName
		if ifName == "" {
			return fmt.Errorf("network %s has per-interface sysctls but no interface name", network.Name)
		}

		sysctls, err := parseIFSysctls(sysctlVal, ifName)
		if err != nil {
			return fmt.Errorf("parsing per-interface sysctls for network %s: %w", network.Name, err)
		}

		// Enter the container's network namespace and apply sysctls.
		err = netns.WithNetNSPath(namespacePath, func(_ netns.NetNS) error {
			for _, sysctl := range sysctls {
				path := "/proc/sys/" + strings.ReplaceAll(sysctl.key, ".", "/")
				if err := os.WriteFile(path, []byte(sysctl.value), 0o644); err != nil {
					return fmt.Errorf("setting sysctl %s=%s: %w", sysctl.key, sysctl.value, err)
				}
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("applying sysctls for network %s interface %s: %w", network.Name, ifName, err)
		}
	}
	return nil
}

// parseIFSysctls parses a comma-separated list of sysctl=value pairs and
// resolves IFNAME placeholders to the actual interface name.
// Input format: "net.ipv4.conf.IFNAME.log_martians=1,net.ipv6.conf.IFNAME.accept_ra=0"
func parseIFSysctls(val, ifName string) ([]ifSysctl, error) {
	var sysctls []ifSysctl
	for _, item := range strings.Split(val, ",") {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		parts := strings.SplitN(item, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid sysctl format: %q (expected key=value)", item)
		}

		key := strings.Replace(parts[0], "IFNAME", ifName, 1)
		value := parts[1]

		sysctls = append(sysctls, ifSysctl{key: key, value: value})
	}
	return sysctls, nil
}
