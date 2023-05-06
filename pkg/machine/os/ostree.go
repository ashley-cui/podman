//go:build amd64 || arm64
// +build amd64 arm64

package os

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/containers/image/v5/transports/alltransports"
	"github.com/sirupsen/logrus"
)

// OSTree deals with operations on ostree based os's
type OSTree struct { //nolint:revive
}

// Apply takes an OCI image and does an rpm-ostree rebase on the image
// If no containers-transport is specified,
// apply will first check if the image exists locally, then default to pulling.
// Exec-ing out to rpm-ostree rebase requires sudo, so this means apply cannot
// be called within podman's user namespace if run as rootless.
// This means that we need to export images in containers-storage to oci-dirs
// We also need to do this via an exec, because if we tried to use the ABI functions,
// we would enter the user namespace, the rebase command would fail.
// The pull portion of this function essentially is a work-around for two things:
// 1. rpm-ostree requires you to specify the containers-transport when pulling.
// The pull in podman allows the behavior of os apply to match other podman commands,
// where you only pull if the image does not exist in storage already.
// 2. This works around the root/rootless issue.
// Podman machines are by default set up using a rootless connection.
// rpm-ostree needs to be run as root. If a user wants to use an image in containers-storage,
// rpm-ostree will look at the root storage, and not the user storage, which is unexpected behavior.
// Exporting to an oci-dir works around this, without nagging the user to configure the machine in rootful mode.
func (dist *OSTree) Apply(image string, opts ApplyOptions) error {
	imageWithTransport := image

	transport := alltransports.TransportFromImageName(image)

	switch {
	// no transport was specified
	case transport == nil:
		exists, err := execPodmanImageExists(image)
		if err != nil {
			return err
		}

		if exists {
			fmt.Println("Pulling from", "containers-storage"+":", imageWithTransport)
			dir, err := os.MkdirTemp("", pathSafeString(imageWithTransport))
			if err != nil {
				return err
			}
			if err := os.Chmod(dir, 0755); err != nil {
				return err
			}

			defer func() {
				if err := os.RemoveAll(dir); err != nil {
					logrus.Errorf("failed to remove temporary pull file: %v", err)
				}
			}()

			if err := execPodmanSave(dir, image); err != nil {
				return err
			}

			imageWithTransport = "oci:" + dir
		} else {
			// if image doesn't exist locally, assume that we want to pull and use docker transport
			imageWithTransport = "docker://" + image
		}
	// containers-transport specified
	case transport.Name() == "containers-storage":
		fmt.Println("Pulling from", image)
		dir, err := os.MkdirTemp("", pathSafeString(strings.TrimPrefix(image, "containers-storage"+":")))
		if err != nil {
			return err
		}

		if err := os.Chmod(dir, 0755); err != nil {
			return err
		}

		defer func() {
			if err := os.RemoveAll(dir); err != nil {
				logrus.Errorf("failed to remove temporary pull file: %v", err)
			}
		}()

		if err := execPodmanSave(dir, image); err != nil {
			return err
		}
		imageWithTransport = "oci:" + dir
	}

	staged := false

	knownDeployments, err := ManagedDeployments()
	if err != nil {
		return err
	}
	existingDeployments, err := getDeployments()
	if err != nil {
		return err
	}

	for _, dep := range existingDeployments {
		if dep.staged && dep.Deployment == knownDeployments.Current {
			staged = true
		}
	}

	ostreeCli := []string{"rpm-ostree", "--bypass-driver", "rebase", fmt.Sprintf("ostree-unverified-image:%s", imageWithTransport)}
	cmd := exec.Command("sudo", ostreeCli...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}
	return updateManagedDeployments(staged)
	// return nil
}

// func getDeployments()error {
// 	ostreeCli := []string{"rpm-ostree", "--bypass-driver", "rebase", fmt.Sprintf("ostree-unverified-image:%s", imageWithTransport)}
// 	cmd := exec.Command("sudo", ostreeCli...)
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
// 	err := cmd.Run()
// 	if err != nil {
// 		return err
// 	}
// }

// pathSafeString creates a path-safe name for our tmpdirs
func pathSafeString(str string) string {
	alphanumOnly := regexp.MustCompile(`[^a-zA-Z0-9]+`)

	return alphanumOnly.ReplaceAllString(str, "")
}

// execPodmanSave execs out to podman save
func execPodmanSave(dir, image string) error {
	saveArgs := []string{"image", "save", "--format", "oci-dir", "-o", dir, image}

	saveCmd := exec.Command("podman", saveArgs...)
	saveCmd.Stdout = os.Stdout
	saveCmd.Stderr = os.Stderr
	return saveCmd.Run()
}

// execPodmanSave execs out to podman image exists
func execPodmanImageExists(image string) (bool, error) {
	existsArgs := []string{"image", "exists", image}

	existsCmd := exec.Command("podman", existsArgs...)

	if err := existsCmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			switch exitCode := exitError.ExitCode(); exitCode {
			case 1:
				return false, nil
			default:
				return false, errors.New("unable to access local image store")
			}
		}
	}
	return true, nil
}

func (dist *OSTree) Revert(opts RevertOptions) error {
	// fmt.Println(getDeployments())
	knownDeployments, err := ManagedDeployments()
	if err != nil {
		return err
	}
	if knownDeployments == nil || knownDeployments.Last.Checksum == "" {
		return errors.New("cannot revert: no os has been applied")
	}

	existingDeployments, err := getDeployments()
	fmt.Println(existingDeployments)
	if err != nil {
		return err
	}
	var currentDeployment *DeploymentWithState
	var lastDeployment *DeploymentWithState

	for i, dep := range existingDeployments {
		if dep.staged && dep.Deployment == knownDeployments.Current {
			if opts.Staged {
				return revertStaged(knownDeployments)
			} else {
				return errors.New("staged OS operation found, use revert -s to unstage, or reboot to complete OS operation")
			}
		}
		if dep.Deployment.Checksum == knownDeployments.Current.Checksum {
			fmt.Println("currdep here!")
			currentDeployment = &dep
			if !dep.pinned {
				pinDeployment(i)
			}
		} else if dep.Deployment.Checksum == knownDeployments.Last.Checksum {
			lastDeployment = &dep
			if !dep.pinned {
				pinDeployment(i)
			}
		}
	}
	fmt.Println(currentDeployment)
	if currentDeployment == nil || !currentDeployment.booted {
		return errors.New("booted deployment not created using os apply")
	}
	if lastDeployment == nil {
		return errors.New("unable to find previous deployment in cache")
	}

	updatedDeployments := managedDeployments{
		Current:  lastDeployment.Deployment,
		Last:     currentDeployment.Deployment,
		Previous: lastDeployment.Deployment,
	}

	ref := lastDeployment.Checksum
	if lastDeployment.Origin != "" {
		ref = lastDeployment.Origin

	}

	ostreeCli := []string{"rpm-ostree", "--bypass-driver", "rebase", ref}
	cmd := exec.Command("sudo", ostreeCli...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	managedDeploymentsFile, err := managedDeploymentsFile()
	if err != nil {
		return err
	}
	marshalled, err := json.MarshalIndent(updatedDeployments, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(managedDeploymentsFile, marshalled, 0600)

	// if dep.Deployment.checksum == knownDeployments.Current.checksum {

	// 	if dep.staged {
	// 		fmt.Println("unstage deployment, mark unstaged deployment as past, current as current")
	// 		matched = true
	// 		break
	// 	} else if dep.booted {
	// 		matched = true
	// 		fmt.Println("rpmostree rebase, mark current as past, past as current")
	// 		break
	// 	}
	// }

}

func revertStaged(knownDeployments *managedDeployments) error {

	ostreeCli := []string{"rpm-ostree", "cleanup", "-p"}
	cmd := exec.Command("sudo", ostreeCli...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	updatedDeployments := managedDeployments{
		Current: knownDeployments.Last,
		Last:    knownDeployments.Previous,
	}
	managedDeploymentsFile, err := managedDeploymentsFile()
	if err != nil {
		return err
	}
	marshalled, err := json.MarshalIndent(updatedDeployments, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(managedDeploymentsFile, marshalled, 0600)
}

// func revert() error {

// }

func getDeployments() ([]DeploymentWithState, error) {
	cmd := exec.Command("rpm-ostree", "status", "--json")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(cmdOutput.Bytes(), &result)
	if err != nil {
		return nil, err
	}
	// fmt.Println(result["deployments"])
	// deployments := result["deployments"].([]interface{})[0].(map[string]interface{})
	deployments := result["deployments"].([]interface{})
	// fmt.Println(deployments)
	// fmt.Println(len(result["deployments"].([]interface{})))
	allDeployments := []DeploymentWithState{}
	for _, dep := range deployments {
		deploymentDetails := dep.(map[string]interface{})
		entry := DeploymentWithState{}
		val, ok := deploymentDetails["origin"]
		if ok {
			entry.Origin = val.(string)
			entry.Containerimage = false
		} else {
			entry.Containerimage = true
		}
		entry.booted = deploymentDetails["booted"].(bool)
		entry.pinned = deploymentDetails["pinned"].(bool)
		entry.staged = deploymentDetails["staged"].(bool)
		entry.Checksum = deploymentDetails["checksum"].(string)
		allDeployments = append(allDeployments, entry)
	}
	return allDeployments, nil
}

func updateManagedDeployments(staged bool) error {
	managedDeployments := managedDeployments{}
	currentDeployments, err := getDeployments()
	if err != nil {
		return err
	}

	prevDeployments, err := ManagedDeployments()
	//this needs work
	if prevDeployments == nil {
		for i, dep := range currentDeployments {
			if dep.staged {
				managedDeployments.Current = dep.Deployment
			} else if dep.booted {
				managedDeployments.Last = dep.Deployment
				if !dep.pinned {
					pinDeployment(i)
				}
			}
		}
	} else {
		if !staged {
			managedDeployments.Last = prevDeployments.Current
			managedDeployments.Previous = prevDeployments.Last
		}

		for i, dep := range currentDeployments {
			// clean up old pinned deployments
			if dep.pinned && dep.Deployment == prevDeployments.Previous {
				err := unpinDeployment(i)
				if err != nil {
					return err
				}
			}

			if dep.staged {
				managedDeployments.Current = dep.Deployment
			} else if dep.booted {
				managedDeployments.Last = dep.Deployment
				if !dep.pinned {
					pinDeployment(i)
				}
			}
			if dep.Deployment == managedDeployments.Previous && !dep.pinned {
				pinDeployment(i)
			}
		}
	}

	managedDeploymentsFile, err := managedDeploymentsFile()
	if err != nil {
		return err
	}
	marshalled, err := json.MarshalIndent(managedDeployments, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(managedDeploymentsFile, marshalled, 0600)
}

func pinDeployment(index int) error {
	ostreeCli := []string{"ostree", "admin", "pin", strconv.Itoa(index)}
	cmd := exec.Command("sudo", ostreeCli...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err

}

func unpinDeployment(index int) error {
	ostreeCli := []string{"ostree", "admin", "pin", "--unpin", strconv.Itoa(index)}
	cmd := exec.Command("sudo", ostreeCli...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func ManagedDeployments() (*managedDeployments, error) {
	manageddeps := new(managedDeployments)
	managedDeploymentsFile, err := managedDeploymentsFile()
	if err != nil {
		return nil, err
	}
	if _, err := os.Stat(managedDeploymentsFile); os.IsNotExist(err) {
		return nil, nil
	}
	bytes, err := os.ReadFile(managedDeploymentsFile)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, manageddeps)
	if err != nil {
		return nil, err
	}

	return manageddeps, nil

}

func managedDeploymentsFile() (string, error) {
	// 	osDir, err := machine.OSFileDir()
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	managedDeploymentsFile := filepath.Join(osDir, "deployments")
	// 	return managedDeploymentsFile, nil
	return "/home/acui/agh", nil

}

// func CheckIfStaged() ,error {
// 	knownDeployments, err := ManagedDeployments()
// 	if err != nil {
// 		return err
// 	}
// 	existingDeployments, err := getDeployments()
// 	if err != nil {
// 		return err
// 	}

// 	for i, dep := range existingDeployments {
// 		if dep.staged && dep.Deployment == knownDeployments.Current {
// 			return true
// 		}
// 	}
// }

type deps struct {
	entry  string
	commit string
}

type Deployment struct {
	Containerimage bool
	Origin         string
	Checksum       string
}
type DeploymentWithState struct {
	Deployment
	booted bool
	staged bool
	pinned bool
}

type managedDeployments struct {
	Current  Deployment
	Last     Deployment
	Previous Deployment
}
