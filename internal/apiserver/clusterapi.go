package apiserver

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	pb "github.com/samsung-cnct/cma-vmware/pkg/generated/api"
)

const (
	kubectlCmd = "kubectl"

	maxApplyTimeout = 30
)

type ClusterShim struct {
	Name              string
	PrivateKey        string
	ControlPlaneNodes []MachineShim
	WorkerNodes       []MachineShim
}

type MachineShim struct {
	Username            string
	Host                string
	Port                int
	Password            string
	KubeletVersion      string
	ControlPlaneVersion string
}

func TranslateAPI(in *pb.CreateClusterMsg) ClusterShim {
	cluster := ClusterShim{
		Name:       in.Name,
		PrivateKey: in.PrivateKey,
	}

	for _, m := range in.ControlPlaneNodes {
		cluster.ControlPlaneNodes = append(cluster.ControlPlaneNodes, MachineShim{
			Username:            m.Username,
			Password:            m.Password,
			Host:                m.Host,
			Port:                int(m.Port),
			KubeletVersion:      in.K8SVersion,
			ControlPlaneVersion: in.K8SVersion,
		})
	}

	for _, m := range in.WorkerNodes {
		cluster.WorkerNodes = append(cluster.WorkerNodes, MachineShim{
			Username:       m.Username,
			Password:       m.Password,
			Host:           m.Host,
			Port:           int(m.Port),
			KubeletVersion: in.K8SVersion,
		})
	}

	return cluster
}

func GetManifests(cluster ClusterShim) (string, error) {
	tmpl, err := template.New("cluster-api-provider-ssh-cluster").Parse(ClusterAPIProviderSSHTemplate)
	if err != nil {
		return "", err
	}

	var tmplBuf bytes.Buffer
	err = tmpl.Execute(&tmplBuf, cluster)
	if err != nil {
		return "", err
	}

	return string(tmplBuf.Bytes()), nil
}

func ApplyManifests(cluster ClusterShim) error {
	manifests, err := GetManifests(cluster)
	if err != nil {
		return err
	}

	cmdName := kubectlCmd
	cmdArgs := []string{"create", "--validate=false", "-f", "-"}
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second
	_, err = RunCommand(cmdName, cmdArgs, manifests, cmdTimeout)
	if err != nil {
		return err
	}

	return nil
}

// Control plane _machines_ must be deleted before the workers to ensure the
// cooresponding _nodes_ can be drained and deleted. The cluster-private-key
// secret and cluster object must be deleted after all machines; otherwise
// they can not be deleted.
func DeleteManifests(clusterName string) error {
	if clusterName == "" {
		return errors.New("clusterName can not be nil")
	}

	cmdName := kubectlCmd
	cmdArgs := []string{"--help"}
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second

	// Delete workers. Control plane nodes have a non-empty value for the label key controlPlaneVersion.
	cmdArgs = []string{"delete", "machines", "-n", clusterName, "-l", `!controlPlaneVersion`}
	_, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	// Wait for workers to finish being deleted.
	cmdArgs = []string{"wait", "--for=delete", "machines", "-n", clusterName, "-l", `!controlPlaneVersion`}
	_, err = RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	// Delete control plane.
	cmdArgs = []string{"delete", "machines", "-n", clusterName, "-l", "controlPlaneVersion"}
	_, err = RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	// Wait for control plane to finish being deleted.
	cmdArgs = []string{"wait", "--for=delete", "machines", "-n", clusterName, "-l", "controlPlaneVersion"}
	_, err = RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	// Delete the namespace and anything else in it (e.g. the Cluster, Secrets, etc.)
	cmdArgs = []string{"delete", "namespace", clusterName}
	_, err = RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	return nil
}

func GetKubeConfig(clusterName string) (string, error) {
	if clusterName == "" {
		return "", errors.New("clusterName can not be nil")
	}

	cmdName := kubectlCmd
	cmdArgs := []string{"--help"}
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second

	cmdArgs = []string{"get", "secret", clusterName + "-kubeconfig", "-n", clusterName, "-o", "jsonpath={.data.kubeconfig}"}
	encodedKubeconfig, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return "", err
	}

	decodedKubeconfig, err := base64.StdEncoding.DecodeString(string(encodedKubeconfig.Bytes()))
	if err != nil {
		return "", err
	}

	return string(decodedKubeconfig), nil
}

func ListClusters() ([]string, error) {
	cmdName := kubectlCmd
	cmdArgs := []string{"--help"}
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second

	cmdArgs = []string{"get", "clusters", "--all-namespaces", "-o", "jsonpath={.items[*].metadata.name}"}
	stdout, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(stdout.Bytes()), " "), nil
}

// Upgrade (or downgrade) all nodes in the cluster named clusterName to the
// version specified by k8sVersion.
func Upgrade(clusterName, k8sVersion string) error {
	if clusterName == "" {
		return errors.New("clusterName can not be nil")
	}

	cmdName := kubectlCmd
	cmdArgs := []string{"--help"}
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second

	// Get a list of all machines.
	cmdArgs = []string{"get", "machines", "-n", clusterName, "-o", "jsonpath={.items[*].metadata.name}"}
	machineNames, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	// Update each one.
	for _, name := range strings.Split(string(machineNames.Bytes()), " ") {
		// Determine which machines are masters by looking for non-empty
		// controlPlane fields.
		cmdArgs = []string{"get", "machine", name, "-n", clusterName, "-o", "jsonpath={.items[*].spec.versions.controlPlane}"}
		controlPlaneVersion, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
		if err != nil {
			return err
		}

		if string(controlPlaneVersion.Bytes()) != "" {
			cmdArgs = []string{"patch", "machine", name, "-n", clusterName, "-p", `{"spec":{"versions":{"controlPlane":"` + k8sVersion + `"}}}`}
			_, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
			if err != nil {
				return err
			}

		}
		cmdArgs = []string{"patch", "machine", name, "-n", clusterName, "-p", `{"spec":{"versions":{"kubelet":"` + k8sVersion + `"}}}`}
		_, err = RunCommand(cmdName, cmdArgs, "", cmdTimeout)
		if err != nil {
			return err
		}
	}

	return nil
}

// Run command with args and kill if timeout is reached. If streamIn is not empty it will
// also be passed to the command via stdin.
func RunCommand(name string, args []string, streamIn string, timeout time.Duration) (bytes.Buffer, error) {
	var streamOut, streamErr bytes.Buffer

	fmt.Printf("Running command \"%v %v\"\n", name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)

	if streamIn != "" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return streamOut, err
		}

		go func() {
			defer stdin.Close()
			io.WriteString(stdin, streamIn)
		}()
	}
	cmd.Stdout = &streamOut
	cmd.Stderr = &streamErr

	err := cmd.Start()
	if err != nil {
		return streamOut, err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		// We do not print stdout because it may contain secrets.
		fmt.Fprintf(os.Stderr, "Command %v stderr: %v\n", name, string(streamErr.Bytes()))

		if err := cmd.Process.Kill(); err != nil {
			panic(fmt.Sprintf("Failed to kill command %v, err %v", name, err))
		}

		return streamOut, fmt.Errorf("Command %v timed out\n", name)
	case err := <-done:
		// We do not print stdout because it may contain secrets.
		fmt.Fprintf(os.Stderr, "Command %v stderr: %v\n", name, string(streamErr.Bytes()))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Command %v returned err %v\n", name, err)
			return streamOut, err
		}
	}

	fmt.Printf("Command %v completed successfully\n", name)

	return streamOut, nil
}
