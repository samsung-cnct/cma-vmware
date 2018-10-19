package apiserver

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
	"time"

	pb "github.com/samsung-cnct/cma-vmware/pkg/generated/api"
	"github.com/samsung-cnct/cma-vmware/pkg/util"
)

const (
	kubectlCmd = "kubectl"
	sshCmd     = "ssh"

	maxApplyTimeout   = 30
	maxUpgradeTimeout = 540 // !? TODO: Determine a better value for this.
	upgradeRetrySleep = 15
)

type SSHClusterParams struct {
	Name              string
	PrivateKey        string // These are base64 _and_ PEM encoded Eliptic
	PublicKey         string // Curve (EC) keys used in JSON and YAML.
	K8SVersion        string
	ControlPlaneNodes []SSHMachineParams
	WorkerNodes       []SSHMachineParams
}

type SSHMachineParams struct {
	Username string
	Host     string
	Port     int32
	Password string
}

func TranslateCreateClusterMsg(in *pb.CreateClusterMsg) SSHClusterParams {
	cluster := SSHClusterParams{
		Name:       in.Name,
		K8SVersion: in.K8SVersion,
		PrivateKey: in.PrivateKey,
	}

	for _, m := range in.ControlPlaneNodes {
		cluster.ControlPlaneNodes = append(cluster.ControlPlaneNodes, SSHMachineParams{
			Username: m.Username,
			Password: m.Password,
			Host:     m.Host,
			Port:     m.Port,
		})
	}

	for _, m := range in.WorkerNodes {
		cluster.WorkerNodes = append(cluster.WorkerNodes, SSHMachineParams{
			Username: m.Username,
			Password: m.Password,
			Host:     m.Host,
			Port:     m.Port,
		})
	}

	return cluster
}

func TranslateAdjustClusterMsg(in *pb.AdjustClusterMsg, version string) SSHClusterParams {
	cluster := SSHClusterParams{
		Name:       in.Name,
		K8SVersion: version,
	}

	for _, m := range in.AddNodes {
		cluster.WorkerNodes = append(cluster.WorkerNodes, SSHMachineParams{
			Username: m.Username,
			Password: m.Password,
			Host:     m.Host,
			Port:     m.Port,
		})
	}

	return cluster
}

// Renders a Namespace, Cluster, and a Secret. Also renders all Machines.
func RenderClusterManifests(cluster SSHClusterParams) (string, error) {
	tmpl, err := template.New("cluster-api-provider-ssh-cluster").Parse(SSHClusterTemplate)
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

func PrepareNodes(cluster *SSHClusterParams) error {
	private, public, err := util.GenerateSSHKeyPair()
	if err != nil {
		return err
	}

	cluster.PrivateKey = base64.StdEncoding.EncodeToString([]byte(private))
	cluster.PublicKey = base64.StdEncoding.EncodeToString([]byte(public))

	for _, node := range cluster.ControlPlaneNodes {
		err := setupPrivateKeyAccess(node, private, public)
		if err != nil {
			return err
		}
	}

	for _, node := range cluster.WorkerNodes {
		err := setupPrivateKeyAccess(node, private, public)
		if err != nil {
			return err
		}
	}

	return nil
}

func setupPrivateKeyAccess(machine SSHMachineParams, privateKey string, publicKey string) error {
	//TODO: add public key to local known_hosts (?)

	err := util.AddPublicKeyToRemoteNode(
		machine.Host,
		machine.Port,
		machine.Username,
		machine.Password,
		publicKey)
	if err != nil {
		fmt.Printf("ERROR: Failed to add public key to %s@%s:%d\n",
			machine.Username, machine.Host, machine.Port)
		return err
	}

	// Test private key
	testCmd := "echo cma-vmware: $(date) >> ~/.ssh/test-pvka"

	authMethod, err := util.SSHAuthMethPublicKey(privateKey)
	if err != nil {
		fmt.Printf("ERROR: Failed generate a public key for ssh authentication")
		return err
	}

	err = util.ExecuteCommandOnRemoteNode(machine.Host, machine.Port, machine.Username, authMethod, testCmd)
	if err != nil {
		fmt.Printf("ERROR: Failed to execute test command via private key on remote node")
		return err
	}

	return nil
}

// Renders all Machines (both control plane and worker).
func RenderMachineManifests(cluster SSHClusterParams) (string, error) {
	tmpl, err := template.New("cluster-api-provider-ssh-machine").Parse(SSHMachineTemplate)
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

func CreateSSHCluster(in *pb.CreateClusterMsg) error {
	cluster := TranslateCreateClusterMsg(in)
	if cluster.PrivateKey == "" {
		err := PrepareNodes(&cluster)
		if err != nil {
			return err
		}
	}

	manifests, err := RenderClusterManifests(cluster)
	if err != nil {
		return err
	}

	cmdName := kubectlCmd
	cmdArgs := []string{"--help"}
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second

	// Create all cluster resources.
	cmdArgs = []string{"create", "--validate=false", "-f", "-"}
	_, err = RunCommand(cmdName, cmdArgs, manifests, cmdTimeout)
	if err != nil {
		return err
	}

	return nil
}

// Control plane _machines_ must be deleted before the workers to ensure the
// cooresponding _nodes_ can be drained and deleted. The cluster-private-key
// secret and cluster object must be deleted after all machines; otherwise
// the machines can not be deleted.
func DeleteSSHCluster(clusterName string) error {
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
	err = waitForMachinesDeleted(clusterName, false)
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
	err = waitForMachinesDeleted(clusterName, true)
	if err != nil {
		return err
	}

	// Delete the cluster resource
	clusterCmdArgs := []string{"delete", "cluster", clusterName, "-n", clusterName}
	_, err = RunCommand(cmdName, clusterCmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}
	time.Sleep(10 * time.Second)

	// Delete the namespace and anything else in it (e.g. the Cluster, Secrets, etc.)
	cmdArgs = []string{"delete", "namespace", clusterName}
	_, err = RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	return nil
}

// The kubeconfig for each cluster is stored as a Secret.
func GetKubeConfig(clusterName string) ([]byte, error) {
	if clusterName == "" {
		return nil, errors.New("clusterName can not be nil")
	}

	cmdName := kubectlCmd
	cmdArgs := []string{"--help"}
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second

	cmdArgs = []string{"get", "secret", clusterName + "-kubeconfig", "-n", clusterName, "-o", "jsonpath={.data.kubeconfig}"}
	encodedKubeconfig, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return nil, err
	}

	decodedKubeconfig, err := base64.StdEncoding.DecodeString(string(encodedKubeconfig.Bytes()))
	if err != nil {
		return nil, err
	}

	return decodedKubeconfig, nil
}

func ListSSHClusters() ([]string, error) {
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

func removeDuplicates(s []string) []string {
	result := []string{}
	seen := make(map[string]bool)

	for _, x := range s {
		if !seen[x] {
			result = append(result, x)
			seen[x] = true
		}
	}

	return result
}

func AdjustSSHCluster(in *pb.AdjustClusterMsg) error {
	cmdName := kubectlCmd
	cmdArgs := []string{"--help"}
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second

	// Get kubelet version for all machines in cluster namespace.
	cmdArgs = []string{"get", "machines", "-n", in.Name, "-o", "jsonpath={.items[*].spec.versions.kubelet}"}
	allVersions, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	// Since the CMA API only allows a single version to be passed during
	// create and update, all machines should be using the same version.
	// They might not be if there was a failure after the control plane
	// was upgraded but before the workers.
	uniqueVersions := removeDuplicates(strings.Split(string(allVersions.Bytes()), " "))
	if len(uniqueVersions) != 1 {
		return fmt.Errorf("expected exactly one k8s version, found %v", len(uniqueVersions))
	}
	version := uniqueVersions[0]

	// Generate manifests for new machines.
	cluster := TranslateAdjustClusterMsg(in, version)
	manifests, err := RenderMachineManifests(cluster)
	if err != nil {
		return err
	}

	// Create added machines.
	cmdArgs = []string{"create", "--validate=false", "-f", "-"}
	_, err = RunCommand(cmdName, cmdArgs, manifests, cmdTimeout)
	if err != nil {
		return err
	}

	// Delete each removed machine.
	for _, m := range in.RemoveNodes {
		cmdArgs = []string{"delete", "machine", m.Host, "-n", in.Name}
		_, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
		if err != nil {
			return err
		}
	}

	return nil
}

func patchMachineVersions(clusterName, machineName, controlPlaneVersion, kubeletVersion string) error {
	cmdName := kubectlCmd
	cmdArgs := []string{"--help"}
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second

	if controlPlaneVersion != "" {
		cmdArgs = []string{"patch", "machine", machineName, "-n", clusterName, "-p", `{"spec":{"versions":{"controlPlane":"` + controlPlaneVersion + `"}}}`}
		_, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
		if err != nil {
			cmdArgs = []string{"get", "machine", machineName, "-n", clusterName, "-o", "jsonpath={.spec.versions.controlPlane}"}
			observeredVersionBuffer, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
			if err != nil {
				return err
			}

			observeredVersion := string(observeredVersionBuffer.Bytes())
			if observeredVersion != controlPlaneVersion {
				return fmt.Errorf("failed to set controlPlane version (from %s to %s) for machine %s in cluster %s)",
					observeredVersion, controlPlaneVersion, machineName, clusterName)
			}
		}
	}

	cmdArgs = []string{"patch", "machine", machineName, "-n", clusterName, "-p", `{"spec":{"versions":{"kubelet":"` + kubeletVersion + `"}}}`}
	_, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		cmdArgs = []string{"get", "machine", machineName, "-n", clusterName, "-o", "jsonpath={.spec.versions.kubelet}"}
		observeredVersionBuffer, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
		if err != nil {
			return err
		}

		observeredVersion := string(observeredVersionBuffer.Bytes())
		if observeredVersion != kubeletVersion {
			return fmt.Errorf("failed to set kubelet version (from %s to %s) for machine %s in cluster %s)",
				observeredVersion, kubeletVersion, machineName, clusterName)
		}
	}

	return nil
}

// parseVersionFromNodes takes formatted output from get nodes, and the machine
// name 'namespace/machineName'. It returns the semantic version of the node
func parseVersionFromNodes(outbytes []byte, fullMachineName string) (string, error) {
	for _, nodeLine := range strings.Split(string(outbytes), "\n") {
		if strings.Contains(nodeLine, fullMachineName) {
			strs := strings.Split(string(nodeLine), " ")
			if len(strs) != 3 {
				return "", errors.New("parseVersionFromNodes, could not parse version")
			}
			version := semanticVersion(strs[1])
			return version, nil
		}
	}
	return "", errors.New("parseVersionFromNodes, could not find fullmachineName " + fullMachineName)
}

func kubeletVersionMatch(clusterName string, machineName string, expectedVersion string, kubeconfigfn string) (bool, error) {
	if machineName == "" {
		fmt.Printf("ERROR: kubeletVersionMatch, invalid machineName\n")
		return false, nil
	}
	getNodesCmdArgs := []string{"get", "nodes", "-o", "go-template='{{range .items}}{{.metadata.name}} {{.status.nodeInfo.kubeletVersion}} {{.metadata.annotations.machine}}{{\"\\n\"}}{{end}}'", "--kubeconfig", kubeconfigfn}
	cmdName := kubectlCmd
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second
	nodeOutput, err := RunCommand(cmdName, getNodesCmdArgs, "", cmdTimeout)
	if err != nil {
		return false, err
	}
	fullMachineName := clusterName + "/" + machineName
	nodeVersion, err := parseVersionFromNodes(nodeOutput.Bytes(), fullMachineName)
	if err != nil {
		return false, err
	}
	if nodeVersion != expectedVersion {
		return false, nil
	}
	return true, nil
}

// machinesDeleted returns true when there are no machines of the boolean
// specified.  masters should be false when deleting workers.  masters should
// be true when deleting masters.
func machinesDeleted(clusterName string, masters bool) (bool, error) {
	getMastersCmd := []string{"get", "machines", "-n", clusterName, "-l", "controlPlaneVersion"}
	getWorkersCmd := []string{"get", "machines", "-n", clusterName, "-l", `!controlPlaneVersion`}
	cmdName := kubectlCmd
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second
	var machinesFound bytes.Buffer
	var err error
	if masters {
		machinesFound, err = RunCommand(cmdName, getMastersCmd, "", cmdTimeout)
	} else {
		machinesFound, err = RunCommand(cmdName, getWorkersCmd, "", cmdTimeout)
	}
	if err != nil {
		return false, err
	}
	if len(machinesFound.Bytes()) == 0 {
		return true, nil
	}
	return false, nil
}

// waitForMachinesDeleted polls machinesDeleted until timeout
// clusterName is the namespace of the cluster.
// masters should be false when deleting workers in the namespace.
func waitForMachinesDeleted(clusterName string, masters bool) error {
	fmt.Printf("INFO: waitForMachinesDeleted clusterName = %s, masters = %v\n", clusterName, masters)
	done := make(chan error, 1)
	go func() {
		for i := 0; i*upgradeRetrySleep < maxUpgradeTimeout; i++ {
			deleted, err := machinesDeleted(clusterName, masters)
			if err != nil {
				fmt.Printf("WARN: waitForMachinesDeleted, error from machinesDeleted %v\n", err)
				done <- err
				break
			}
			if deleted {
				fmt.Printf("INFO: machines deleted for cluster %s masters %v\n", clusterName, masters)
				done <- nil
				break
			}
			time.Sleep(time.Duration(upgradeRetrySleep) * time.Second)
		}
	}()

	select {
	case <-time.After(maxUpgradeTimeout * time.Second):
		return fmt.Errorf("WARN: timed out waiting for machines deleted in cluster %s.  masters %v", clusterName, masters)
	case err := <-done:
		if err != nil {
			return err
		}
	}
	fmt.Println("INFO: waitForMachinesDeleted returning successfully\n")

	return nil
}

func ClusterExists(clusterName string) (bool, error) {
	cmdName := kubectlCmd
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second
	getClusterCmdArgs := []string{"get", "clusters", "-n", clusterName, "-o", "go-template={{range .items}}{{.metadata.name}}{{\"\\n\"}}{{end}}"}
	machineName, err := RunCommand(cmdName, getClusterCmdArgs, "", cmdTimeout)
	if err != nil {
		return false, err
	}
	if len(machineName.Bytes()) == 0 {
		return false, errors.New("Cluster NotFound")
	}
	return true, nil
}

func GetSSHClusterStatus(in *pb.GetClusterMsg, kubeconfig []byte) (pb.ClusterStatus, error) {
	cmdName := kubectlCmd
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second
	// init return value
	var clusterStatus pb.ClusterStatus = pb.ClusterStatus_STATUS_UNSPECIFIED
	if kubeconfig == nil {
		return pb.ClusterStatus_PROVISIONING, nil
	}

	file, err := ioutil.TempFile("/tmp", in.Name)
	if err != nil {
		return clusterStatus, err
	}
	defer os.Remove(file.Name())
	kubeconfigfn := file.Name()
	err = ioutil.WriteFile(kubeconfigfn, kubeconfig, 0644)
	if err != nil {
		return clusterStatus, err
	}
	fmt.Printf("INFO: temporarily writing kubeconfig to %s\n", kubeconfigfn)

	// Get a list of all machines.
	getMachineCmdArgs := []string{"get", "machines", "-n", in.Name, "-o", "go-template={{range .items}}{{.metadata.name}} {{.spec.versions.kubelet}}{{\"\\n\"}}{{end}}"}
	machineNames, err := RunCommand(cmdName, getMachineCmdArgs, "", cmdTimeout)
	if err != nil {
		return clusterStatus, err
	}

	for _, machineName := range strings.Split(string(machineNames.Bytes()), "\n") {
		// compare spec version with running node version
		machineInfo := strings.Split(machineName, " ")
		if len(machineInfo) == 2 {
			matchingVersions, err := kubeletVersionMatch(in.Name, machineInfo[0], machineInfo[1], kubeconfigfn)
			if err != nil {
				fmt.Printf("ERROR: GetSSHClusterStatus, kubelet version match error %v\n", err)
			}
			if !matchingVersions {
				clusterStatus = pb.ClusterStatus_RECONCILING
				break
			}
		}
		clusterStatus = pb.ClusterStatus_RUNNING
	}
	return clusterStatus, nil
}

// Waits for the node associated with the machine namespace/name to report
// the expected kubelet version.
func waitForKubeletVersion(clusterName string, machineName string, expectedVersion string, kubeconfigfn string) error {
	fmt.Printf("INFO: waitforKubeletVersion clusterName = %s, machineName = %s, expectedVersion = %s\n", clusterName, machineName, expectedVersion)
	done := make(chan error, 1)
	go func() {
		for i := 0; i*upgradeRetrySleep < maxUpgradeTimeout; i++ {

			matchingVersions, err := kubeletVersionMatch(clusterName, machineName, expectedVersion, kubeconfigfn)
			if err != nil {
				fmt.Printf("WARN: waitForKubeletVersion, kubelet version match error %v\n", err)
			}
			if matchingVersions {
				fmt.Printf("INFO: machine versions matched for machine %s\n", machineName)
				done <- nil
				break
			}

			time.Sleep(time.Duration(upgradeRetrySleep) * time.Second)
		}
	}()

	select {
	case <-time.After(maxUpgradeTimeout * time.Second):
		return fmt.Errorf("WARN: timed out waiting for machine %v to upgrade to kubelet verson %v", machineName, expectedVersion)
	case err := <-done:
		if err != nil {
			return err
		}
	}
	fmt.Println("INFO: waitForKubeletVersion returning successfully\n")

	return nil
}

// Upgrade (or downgrade) all nodes in the cluster named clusterName to the
// version specified by k8sVersion.
func UpgradeSSHCluster(clusterName, k8sVersion string, kubeconfig []byte) error {
	if clusterName == "" {
		return errors.New("UpgradeSSHCluster, clusterName can not be nil")
	}
	if kubeconfig == nil {
		return errors.New("UpgradeSSHCluster, kubeconfig can not be nil")
	}
	// create temp kubeconfig file
	file, err := ioutil.TempFile("/tmp", clusterName)
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	kubeconfigfn := file.Name()
	err = ioutil.WriteFile(kubeconfigfn, kubeconfig, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("INFO: temporarily writing kubeconfig to %s\n", kubeconfigfn)

	cmdName := kubectlCmd
	cmdArgs := []string{"--help"}
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second

	// Get a list of all machines.
	cmdArgs = []string{"get", "machines", "-n", clusterName, "-o", "jsonpath={.items[*].metadata.name}"}
	machineNames, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
	if err != nil {
		return err
	}

	// Update control plane.
	var controlPlaneMachines, workerMachines []string
	for _, machineName := range strings.Split(string(machineNames.Bytes()), " ") {
		// Determine which machines are masters by looking for non-empty
		// controlPlane fields.
		cmdArgs = []string{"get", "machine", machineName, "-n", clusterName, "-o", "jsonpath={.spec.versions.controlPlane}"}
		controlPlaneVersionBuffer, err := RunCommand(cmdName, cmdArgs, "", cmdTimeout)
		if err != nil {
			return err
		}

		controlPlaneVersion := string(controlPlaneVersionBuffer.Bytes())
		if controlPlaneVersion != "" {
			controlPlaneMachines = append(controlPlaneMachines, machineName)

			err = patchMachineVersions(clusterName, machineName, k8sVersion, k8sVersion)
			if err != nil {
				return err
			}

			// Wait for node to be updated before proceeding to the
			// next one. This ensures the control plane is available
			// while the workers are upgraded.
			err = waitForKubeletVersion(clusterName, machineName, k8sVersion, kubeconfigfn)
			if err != nil {
				return err
			}
		} else {
			// Remeber worker machines for later.
			workerMachines = append(workerMachines, machineName)
		}
	}

	// Update workers.
	for _, machineName := range workerMachines {
		err = patchMachineVersions(clusterName, machineName, "", k8sVersion)
		if err != nil {
			return err
		}

		// Wait for node to be updated before proceeding to the next
		// one.
		err = waitForKubeletVersion(clusterName, machineName, k8sVersion, kubeconfigfn)
		if err != nil {
			return err
		}
	}

	return nil
}

func semanticVersion(version string) string {
	semVersion := strings.TrimPrefix(version, "v")
	semVersion = strings.TrimSpace(semVersion)
	return semVersion
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

		return streamOut, fmt.Errorf("Command %v timed out", name)
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
