# Docker API v1.46: DriverOpts sysctls Implementation

## Overview

Implementing support for `com.docker.network.endpoint.sysctls` in Podman's Docker compat API (v1.46+).

This allows Docker clients to set per-interface sysctls via `EndpointSettings.DriverOpts` in both:
- Container create: `/containers/create` with `NetworkingConfig.EndpointsConfig[network].DriverOpts`
- Network connect: `/networks/{id}/connect` with `EndpointConfig.DriverOpts`

## Background

### Why Docker Added This

**Upstream PR:** [moby/moby#47686 - Per-interface sysctls](https://github.com/moby/moby/pull/47686)

Docker 26.0 removed a prestart hook, which broke `--sysctl` flags that reference interface names (e.g., `--sysctl net.ipv4.conf.eth0.log_martians=1`). The interfaces don't exist when the OCI runtime applies sysctls from the spec, causing "No such file or directory" errors.

Docker's workaround in v1.46: `com.docker.network.endpoint.sysctls` in `DriverOpts`, using `IFNAME` as a placeholder that gets resolved after interfaces are created.

**Related PRs:**
- [moby/moby#48746](https://github.com/moby/moby/pull/48746) - Migration handling from v28.0.0
- [moby/moby#48971](https://github.com/moby/moby/pull/48971) - IPv6 address handling with endpoint sysctls

**Example Docker usage:**
```bash
docker run --network mynet:driver-opt=com.docker.network.endpoint.sysctls="net.ipv4.conf.IFNAME.log_martians=1,net.ipv6.conf.IFNAME.accept_ra=0" ...
```

The value is a comma-separated list of `sysctl=value` pairs where `IFNAME` is replaced with the actual interface name (e.g., `eth0`).

### Why Podman Doesn't Have Docker's Timing Bug

Podman creates the network namespace and interfaces **before** the OCI runtime starts:

**`libpod/container_internal_linux.go:87-93`:**
```go
if c.config.CreateNetNS && noNetNS && !c.config.PostConfigureNetNS {
    c.reservedPorts, createNetNSErr = c.bindPorts()
    if createNetNSErr != nil { return }
    netNS, networkStatus, createNetNSErr = c.runtime.createNetNS(c)
}
```

This happens in `prepare()` before the OCI runtime's prestart. So `--sysctl net.ipv4.conf.eth0.log_martians=1` works fine in Podman.

**Testing proves this:**
```bash
podman run --rm --network testnet --sysctl net.ipv4.conf.eth0.log_martians=1 alpine cat /proc/sys/net/ipv4/conf/eth0/log_martians
# Outputs: 1
```

However, Podman does NOT resolve `IFNAME` placeholders:
```bash
podman run --rm --network testnet --sysctl net.ipv4.conf.IFNAME.log_martians=1 alpine cat /proc/sys/net/ipv4/conf/eth0/log_martians
# Error: open /proc/sys/net/ipv4/conf/IFNAME/log_martians: no such file or directory
```

So for Docker compat API compatibility, we need to support `com.docker.network.endpoint.sysctls` with `IFNAME` resolution.

## Architecture Decision

**Sysctl application happens in container-libs (libnetwork), NOT in Podman.**

**Why:**
- Podman's compat handler only translates Docker's format to libnetwork's format
- Libnetwork's `Setup()` knows the namespace path, interface names, and runs after netavark creates interfaces
- Keeps the logic reusable across different consumers of libnetwork
- Network connect uses the same code path

**Flow:**
1. Compat API handler reads `DriverOpts["com.docker.network.endpoint.sysctls"]`
2. Handler passes raw value (with `IFNAME` unresolved) to `PerNetworkOptions.Options[types.PerIFSysctl]`
3. Libnetwork's `Setup()` (after netavark returns) reads that option
4. Libnetwork resolves `IFNAME` using the known `InterfaceName` from `PerNetworkOptions`
5. Libnetwork enters container netns and writes to `/proc/sys/`

## Changes Made So Far

### Container-libs (`/home/acui/go/src/github.com/podman-container-tools/container-libs`)

#### 1. Added constant: `common/libnetwork/types/const.go`

```go
// PerIFSysctl is the key in PerNetworkOptions.Options for per-interface
// sysctls. The value is a comma-separated list of sysctl=value pairs
// where IFNAME is used as a placeholder for the actual interface name.
PerIFSysctl = "per_if_sysctl"
```

**Line 51-54** in `types/const.go`.

#### 2. Created Linux implementation: `common/libnetwork/netavark/sysctl_linux.go`

**New file**, ~63 lines. Key functions:
- `applyPerIFSysctls(namespacePath string, networks []types.NamedPerNetworkOptions) error`
  - Iterates networks looking for `PerIFSysctl` in `Options`
  - Parses comma-separated sysctls, resolves `IFNAME` to `InterfaceName`
  - Enters netns via `netns.WithNetNSPath`
  - Writes to `/proc/sys/` files

- `parseIFSysctls(val, ifName string) ([]ifSysctl, error)`
  - Parses `"key=value,key2=value2"` format
  - Replaces first occurrence of `IFNAME` with actual interface name

**Prior art:** `common/libnetwork/slirp4netns/slirp4netns.go:345` — same pattern of `os.WriteFile` to `/proc/sys/` inside a netns.

#### 3. Created FreeBSD stub: `common/libnetwork/netavark/sysctl_freebsd.go`

```go
func applyPerIFSysctls(_ string, _ []types.NamedPerNetworkOptions) error {
    return nil
}
```

No-op on FreeBSD (no `/proc/sys/`).

#### 4. Wired into Setup: `common/libnetwork/netavark/run.go:161`

Added before the final return in `Setup()`:
```go
if err := applyPerIFSysctls(namespacePath, options.Networks); err != nil {
    return nil, err
}

return result, err
```

This runs **after** netavark setup completes (line 145 validates results), so interfaces exist in the container netns.

**Build status:** ✅ Compiles cleanly (`go build ./common/libnetwork/netavark/`)

### Podman (`/home/acui/go/src/github.com/podman-container-tools/podman`)

Current working tree has **old prototype code that needs revision/removal**. Summary from `git diff HEAD --stat`:

```
libpod/define/container.go                   |  4 ++++   ❌ REMOVE
libpod/networking_common.go                  |  6 ++++++   ❌ REMOVE
libpod/networking_freebsd.go                 |  4 ++++   ❌ REMOVE
libpod/networking_linux.go                   | 22 +++++++++++++++++++   ❌ REMOVE
pkg/api/handlers/compat/containers_create.go | 26 ++++++++++++++++++++   ⚠️ REVISE
pkg/api/handlers/compat/networks.go          | 10 ++++++++++  ⚠️ REVISE
```

## What Still Needs to Be Done

### 1. Clean up Podman repo

**Remove these changes entirely:**
- `libpod/define/container.go` — Remove `EndpointSysctlsOption` constant (lines 58-60)
- `libpod/networking_linux.go` — Remove `applyEndpointSysctls` function (lines 202-221)
- `libpod/networking_common.go` — Remove sysctl application call (lines 622-626)
- `libpod/networking_freebsd.go` — Remove `applyEndpointSysctls` stub (lines 228-230)

### 2. Revise `pkg/api/handlers/compat/containers_create.go`

**Current code (WRONG):** Lines 387-406 parse the sysctls and merge into `cc.HostConfig.Sysctls`

**Should be:** Pass the raw value to `PerNetworkOptions.Options[types.PerIFSysctl]`

**Replacement code:**
```go
// Docker API v1.46+: handle per-interface sysctls from DriverOpts.
// Pass through to libnetwork which will resolve IFNAME and apply them.
if sysctlVal, ok := endpoint.DriverOpts["com.docker.network.endpoint.sysctls"]; ok && sysctlVal != "" {
    if netOpts.Options == nil {
        netOpts.Options = make(map[string]string)
    }
    netOpts.Options[types.PerIFSysctl] = sysctlVal
}
```

**Location:** Inside the `for netName, endpoint := range endpointsConfig` loop, after the `StaticMAC` handling (around line 380), before the bridge mode check.

**Remove:** The `ifIndex` counter (line 344) and the whole sysctl parsing block (lines 387-406).

### 3. Revise `pkg/api/handlers/compat/networks.go`

**Current code (WRONG):** Lines 418-426 use `define.EndpointSysctlsOption`

**Should be:** Use `types.PerIFSysctl`

**Replacement code:**
```go
// Docker API v1.46+: pass per-interface sysctls from DriverOpts
// through to libnetwork which will resolve IFNAME and apply them.
if sysctlVal, ok := netConnect.EndpointConfig.DriverOpts["com.docker.network.endpoint.sysctls"]; ok && sysctlVal != "" {
    if netOpts.Options == nil {
        netOpts.Options = make(map[string]string)
    }
    netOpts.Options[types.PerIFSysctl] = sysctlVal
}
```

**Location:** After the `StaticMAC` handling (line 416), before the `runtime.ConnectContainerToNetwork` call.

### 4. Vendor container-libs into Podman

```bash
cd /home/acui/go/src/github.com/podman-container-tools/podman
# Update go.mod to point to local container-libs or latest commit
go mod edit -replace go.podman.io/common=../container-libs/common
# OR after container-libs changes are merged:
# go get go.podman.io/common@<commit-hash>
go mod tidy
go mod vendor
```

### 5. Test the implementation

#### Test 1: Container create with endpoint sysctls

```bash
# Create network
podman network create testnet

# Start container with per-interface sysctl via API
curl -X POST http://localhost:8080/v1.46/containers/create \
  -H 'Content-Type: application/json' \
  -d '{
    "Image": "alpine",
    "Cmd": ["cat", "/proc/sys/net/ipv4/conf/eth0/log_martians"],
    "NetworkingConfig": {
      "EndpointsConfig": {
        "testnet": {
          "DriverOpts": {
            "com.docker.network.endpoint.sysctls": "net.ipv4.conf.IFNAME.log_martians=1"
          }
        }
      }
    }
  }'

# Start the container and check output
# Should output: 1
```

#### Test 2: Network connect with endpoint sysctls

```bash
# Create container without sysctl
cid=$(podman run -d --network none alpine sleep 1000)

# Connect via API with sysctl
curl -X POST "http://localhost:8080/v1.46/networks/testnet/connect" \
  -H 'Content-Type: application/json' \
  -d "{
    \"Container\": \"$cid\",
    \"EndpointConfig\": {
      \"DriverOpts\": {
        \"com.docker.network.endpoint.sysctls\": \"net.ipv4.conf.IFNAME.log_martians=1,net.ipv6.conf.IFNAME.accept_ra=0\"
      }
    }
  }"

# Verify sysctls were applied
podman exec $cid cat /proc/sys/net/ipv4/conf/eth0/log_martians
# Should output: 1
podman exec $cid cat /proc/sys/net/ipv6/conf/eth0/accept_ra
# Should output: 0
```

#### Test 3: Multiple networks with different sysctls

```bash
podman network create net1
podman network create net2

curl -X POST http://localhost:8080/v1.46/containers/create \
  -H 'Content-Type: application/json' \
  -d '{
    "Image": "alpine",
    "Cmd": ["sh", "-c", "cat /proc/sys/net/ipv4/conf/eth0/log_martians && cat /proc/sys/net/ipv4/conf/eth1/log_martians"],
    "NetworkingConfig": {
      "EndpointsConfig": {
        "net1": {
          "DriverOpts": {
            "com.docker.network.endpoint.sysctls": "net.ipv4.conf.IFNAME.log_martians=1"
          }
        },
        "net2": {
          "DriverOpts": {
            "com.docker.network.endpoint.sysctls": "net.ipv4.conf.IFNAME.log_martians=0"
          }
        }
      }
    }
  }'

# Should output:
# 1
# 0
```

#### Test 4: Invalid sysctl format

```bash
curl -X POST http://localhost:8080/v1.46/containers/create \
  -H 'Content-Type: application/json' \
  -d '{
    "Image": "alpine",
    "NetworkingConfig": {
      "EndpointsConfig": {
        "testnet": {
          "DriverOpts": {
            "com.docker.network.endpoint.sysctls": "invalid-no-equals"
          }
        }
      }
    }
  }'

# Should return error: "invalid sysctl format: \"invalid-no-equals\""
```

#### Test 5: Verify regular --sysctl still works

```bash
# Existing podman functionality should be unaffected
podman run --rm --network testnet --sysctl net.ipv4.conf.eth0.log_martians=1 \
  alpine cat /proc/sys/net/ipv4/conf/eth0/log_martians
# Should output: 1
```

## File Reference

### Container-libs files
```
/home/acui/go/src/github.com/podman-container-tools/container-libs/
├── common/libnetwork/types/const.go              (Modified: added PerIFSysctl)
├── common/libnetwork/netavark/sysctl_linux.go    (New: sysctl application)
├── common/libnetwork/netavark/sysctl_freebsd.go  (New: stub)
└── common/libnetwork/netavark/run.go             (Modified: calls applyPerIFSysctls)
```

### Podman files to modify
```
/home/acui/go/src/github.com/podman-container-tools/podman/
├── pkg/api/handlers/compat/containers_create.go  (Revise: lines 344, 387-406)
├── pkg/api/handlers/compat/networks.go           (Revise: lines 418-426)
├── libpod/define/container.go                    (Remove: lines 58-60)
├── libpod/networking_linux.go                    (Remove: lines 202-221)
├── libpod/networking_common.go                   (Remove: lines 622-626)
└── libpod/networking_freebsd.go                  (Remove: lines 228-230)
```

## Relevant Code Context

### How sysctls are normally handled

**Podman's `--sysctl` flag flow:**
1. CLI parses `--sysctl key=value`
2. `pkg/specgen/generate/security_linux.go:254` calls `g.AddLinuxSysctl(key, val)`
3. Sysctl goes into OCI spec
4. OCI runtime (crun/runc) applies it at container start by writing to `/proc/sys/` inside the container

This works for interface-specific sysctls because Podman creates the network namespace (with interfaces) in `container_internal_linux.go:prepare()` before the OCI runtime starts.

### PerNetworkOptions.Options

`PerNetworkOptions.Options` is a `map[string]string` that:
- Gets serialized to JSON and passed to netavark
- Is NOT currently read by libnetwork's Go code (only passed through to netavark)
- Our new `applyPerIFSysctls` is the first libnetwork Go code to read it

Keys are defined in `types/const.go`:
- `VLANOption = "vlan"`
- `MTUOption = "mtu"`
- `NoDefaultRoute = "no_default_route"`
- **`PerIFSysctl = "per_if_sysctl"`** ← our addition

### Prior art: slirp4netns sysctl handling

`container-libs/common/libnetwork/slirp4netns/slirp4netns.go:345`:
```go
err = os.WriteFile(ipv6ConfDefaultAcceptDadSysctl, []byte("0"), 0o644)
```

Same pattern: directly write to `/proc/sys/` files inside a network namespace. Our implementation follows this approach.

## Docker API Reference

**Endpoint:** `POST /v1.46/containers/create`

**Field:** `NetworkingConfig.EndpointsConfig[network].DriverOpts`

**Type:** `map[string]string`

**Key:** `"com.docker.network.endpoint.sysctls"`

**Value format:** Comma-separated list of `sysctl=value` pairs, e.g.:
```
"net.ipv4.conf.IFNAME.log_martians=1,net.ipv6.conf.IFNAME.accept_ra=0"
```

`IFNAME` is a literal placeholder string that must be replaced with the actual interface name (e.g., `eth0`, `eth1`).

## Questions & Answers

**Q: Why not apply sysctls in netavark (the Rust binary)?**
A: Netavark has internal sysctl infrastructure for its own needs, but is not designed to handle user-specified per-interface sysctls. Libnetwork (the Go library) is the right abstraction layer.

**Q: Why not merge into HostConfig.Sysctls like the current prototype does?**
A: That would apply them via the OCI runtime, which works for container create but NOT for network connect (container is already running, OCI spec is immutable). Libnetwork's `Setup()` is called for both create and connect.

**Q: Does Podman need this feature since it doesn't have Docker's timing bug?**
A: For compat API parity, yes. Docker clients expect this to work. Also, `IFNAME` is a nicer UX than hardcoding `eth0` when you don't know the interface index.

**Q: What about FreeBSD?**
A: No `/proc/sys/` on FreeBSD, so the stub returns `nil` (no-op). Sysctls on FreeBSD would need a different mechanism (sysctl syscall), but that's out of scope for Docker compat API.

**Q: Could this leak into other containers?**
A: No. Sysctls are written inside the container's network namespace via `netns.WithNetNSPath(namespacePath, ...)`. They're scoped to that container.

## Git Branches

- **Container-libs work:** Currently in working tree at `/home/acui/go/src/github.com/podman-container-tools/container-libs`
  - No branch created yet
  - Changes: 3 new/modified files, all compile cleanly

- **Podman work:** Currently on branch `driveropts` at `/home/acui/go/src/github.com/podman-container-tools/podman`
  - Branch exists but has old prototype code
  - Needs cleanup per section "What Still Needs to Be Done"

## Next Steps Summary

1. ✅ Container-libs changes are done and compile
2. ⬜ Clean up podman repo (remove old prototype code)
3. ⬜ Revise compat handlers to use `types.PerIFSysctl`
4. ⬜ Vendor updated container-libs into podman
5. ⬜ Test all 5 test cases above
6. ⬜ Commit container-libs changes
7. ⬜ Commit podman changes
8. ⬜ Open PRs

---

**Contact for questions:** This work was started in a Claude Code session. Check git history or transcript at `/home/acui/.claude/projects/-home-acui-go-src-github-com-podman-container-tools-podman/04050074-f8a5-4928-bbeb-1f5fd84e3c52.jsonl` for full context.
