package apiserver

import (
	"bytes"
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
	Name       string
	PrivateKey string
	Machines   []MachineShim
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
		cluster.Machines = append(cluster.Machines, MachineShim{
			Username:            m.Username,
			Password:            m.Password,
			Host:                m.Host,
			Port:                int(m.Port),
			KubeletVersion:      in.K8SVersion,
			ControlPlaneVersion: in.K8SVersion,
		})
	}

	for _, m := range in.WorkerNodes {
		cluster.Machines = append(cluster.Machines, MachineShim{
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
	cmdArgs := []string{"apply", "--validate=false", "-f", "-"}
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second
	err = RunCommand(cmdName, cmdArgs, manifests, cmdTimeout)
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
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second

	// Delete workers. Control plane nodes have a non-empty value for the label key controlPlane.
	cmdArgs := []string{"delete", "machines", "-n", clusterName, "-l", "controlPlane="}
	err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	// Wait for workers to finish being deleted.
	cmdArgs = []string{"wait", "--for=delete", "machines", "-n", clusterName, "-l", "controlPlane="}
	err = RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	// Delete control plane.
	cmdArgs = []string{"delete", "machines", "-n", clusterName, "--all"}
	err = RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	// Wait for control plane to finish being deleted.
	cmdArgs = []string{"wait", "--for=delete", "machines", "-n", clusterName}
	err = RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	// Delete the namespace and anything else in it (e.g. the Cluster, Secrets, etc.)
	cmdArgs = []string{"delete", "namespace", clusterName, "--all"}
	err = RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	// Wait for the namespace to finish being deleted.
	cmdArgs = []string{"wait", "--for=delete", "ns", clusterName}
	err = RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	return nil
}

// Run command with args and kill if timeout is reached. If streamIn is not empty it will
// also be passed to the command via stdin.
func RunCommand(name string, args []string, streamIn string, timeout time.Duration) error {
	var streamOut, streamErr bytes.Buffer

	fmt.Printf("Running command \"%v %v\"\n", name, strings.Join(args, " "))

	cmd := exec.Command(name, args...)

	if streamIn != "" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return err
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
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		fmt.Fprintf(os.Stderr, "Command %v stdout: %v\n", name, string(streamOut.Bytes()))
		fmt.Fprintf(os.Stderr, "Command %v stderr: %v\n", name, string(streamErr.Bytes()))

		if err := cmd.Process.Kill(); err != nil {
			panic(fmt.Sprintf("Failed to kill command %v, err %v", name, err))
		}

		return fmt.Errorf("Command %v timed out\n", name)
	case err := <-done:
		fmt.Fprintf(os.Stderr, "Command %v stdout: %v\n", name, string(streamOut.Bytes()))
		fmt.Fprintf(os.Stderr, "Command %v stderr: %v\n", name, string(streamErr.Bytes()))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Command %v returned err %v\n", name, err)
			return err
		}
	}

	fmt.Printf("Command %v completed successfully\n", name)

	return nil
}
