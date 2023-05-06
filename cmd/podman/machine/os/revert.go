//go:build amd64 || arm64
// +build amd64 arm64

package os

import (
	"fmt"

	"github.com/containers/podman/v4/cmd/podman/common"
	"github.com/containers/podman/v4/cmd/podman/machine"
	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/containers/podman/v4/cmd/podman/validate"
	os "github.com/containers/podman/v4/pkg/machine/os"
	"github.com/spf13/cobra"
)

var (
	revertCmd = &cobra.Command{
		Use:               "revert [options] IMAGE [NAME]",
		Short:             "revert an OCI image to a Podman Machine's OS",
		Long:              "revert custom layers from a containerized Fedora CoreOS OCI image on top of an existing VM",
		PersistentPreRunE: validate.NoOp,
		RunE:              revert,
		ValidArgsFunction: common.AutocompleteImages,
		Example:           `podman machine os revert myimage`,
	}
)
var staged bool

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Command: revertCmd,
		Parent:  machine.OSCmd,
	})
	flags := revertCmd.Flags()

	stagedFlagName := "staged"
	flags.BoolVarP(&staged, stagedFlagName, "s", false, "Revert staged os action")
}

func revert(cmd *cobra.Command, args []string) error {
	// vmName := ""
	// if len(args) == 2 {
	// 	vmName = args[1]
	// }
	// managerOpts := ManagerOpts{
	// 	VMName:  vmName,
	// 	CLIArgs: args,
	// 	Restart: restart,
	// }
	// // osManager, err := NewOSManager(managerOpts)
	// if err != nil {
	// 	return err
	// }

	vmName := ""
	if len(args) == 2 {
		vmName = args[1]
	}
	managerOpts := ManagerOpts{
		VMName:  vmName,
		CLIArgs: args,
		Restart: restart,
	}
	osManager, err := NewOSManager(managerOpts)
	if err != nil {
		return err
	}
	fmt.Println(staged)
	revertOpts := os.RevertOptions{Staged: staged}
	return osManager.Revert(revertOpts)
}
