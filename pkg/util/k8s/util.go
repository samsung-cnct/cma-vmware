package k8sutil

import (
	"os"

	"github.com/juju/loggo"
	log "github.com/samsung-cnct/cma-vmware/pkg/util"
)

var (
	logger loggo.Logger
)

func SetLogger() {
	logger = log.GetModuleLogger("pkg.util.k8sutil", loggo.INFO)
}

func logErrorAndExit(err error) {
	logger.Criticalf("error: %s", err)
	os.Exit(1)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
