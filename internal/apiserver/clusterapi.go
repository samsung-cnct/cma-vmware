package apiserver

import (
	"bytes"
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
	maxApplyTimeout = 30
)

func GetManifests(in *pb.CreateClusterMsg) (string, error) {
	tmpl, err := template.New("cluster-api-provider-ssh-cluster").Parse(ClusterAPIProviderSSHTemplate)
	if err != nil {
		return "", err
	}

	var tmplBuf bytes.Buffer
	err = tmpl.Execute(&tmplBuf, in)
	if err != nil {
		return "", err
	}

	return string(tmplBuf.Bytes()), nil
}

func ApplyManifests(manifests string) error {
	cmdName := "kubectl"
	cmdArgs := []string{"apply", "--"}
	cmdTimeout := time.Duration(maxApplyTimeout) * time.Second
	err := RunCommand(cmdName, cmdArgs, manifests, cmdTimeout)
	if err != nil {
		return err
	}

	return nil
}

// Run command with args and kill if timeout is reached. If streamIn is non-nil it will
// also be passed to the command via stdin.
func RunCommand(name string, args []string, streamIn string, timeout time.Duration) error {
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
		if err := cmd.Process.Kill(); err != nil {
			panic(fmt.Sprintf("Failed to kill command %v, err %v", name, err))
		}
		err = fmt.Errorf("Command %v timed out\n", name)
		break
	case err := <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Command %v returned err %v\n", name, err)
		}
		break
	}
	if err != nil {
		output, e := cmd.CombinedOutput()
		if e != nil {
			return e
		}
		fmt.Fprintf(os.Stderr, "%v", output)
		return err
	}
	fmt.Printf("Command %v completed successfully\n", name)

	return nil
}
