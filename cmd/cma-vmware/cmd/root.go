package cmd

import (
	"flag"
	"fmt"
	"github.com/juju/loggo"
	"github.com/soheilhy/cmux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/samsung-cnct/cma-vmware/pkg/apiserver"
	"github.com/samsung-cnct/cma-vmware/pkg/util"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	logger loggo.Logger

	RootCmd = &cobra.Command{
		Use:   "cma-vmware",
		Short: "The CMA VMWare Helper",
		Long: `The CMA VMWare Helper

Running this by itself will invoke the webserver to run.
See subcommands for additional features`,
		Run: func(cmd *cobra.Command, args []string) {
			runWebServer()
		},
		TraverseChildren: true,
	}
)

func runWebServer() {
	logger := util.GetModuleLogger("cmd.cmavmware", loggo.INFO)

	// get flags
	portNumber := viper.GetInt("port")

	var wg sync.WaitGroup
	stop := make(chan struct{})

	logger.Infof("Creating Web Server")
	tcpMux := createWebServer(&apiserver.ServerOptions{PortNumber: portNumber})
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Infof("Starting to serve requests on port %d", portNumber)
		tcpMux.Serve()
	}()

	<-stop
	logger.Infof("Waiting for controllers to shut down gracefully")
	wg.Wait()
}

func createWebServer(options *apiserver.ServerOptions) cmux.CMux {
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", options.PortNumber))
	if err != nil {
		panic(err)
	}
	tcpMux := cmux.New(conn)

	apiserver.AddServersToMux(tcpMux, options)

	return tcpMux
}

func init() {

	viper.SetEnvPrefix("cmavmware")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)

	// using standard library "flag" package
	RootCmd.Flags().Int("port", 9020, "Port to listen on")
	RootCmd.PersistentFlags().String("kubeconfig", "", "Location of kubeconfig file")
	RootCmd.PersistentFlags().String("kubectl", "kubectl", "Location of kubectl file")
	RootCmd.PersistentFlags().String("kubernetes-namespace", "default", "What namespace to operate on")

	viper.BindPFlag("port", RootCmd.Flags().Lookup("port"))
	viper.BindPFlag("kubeconfig", RootCmd.PersistentFlags().Lookup("kubeconfig"))
	viper.BindPFlag("kubectl", RootCmd.PersistentFlags().Lookup("kubectl"))
	viper.BindPFlag("kubernetes-namespace", RootCmd.PersistentFlags().Lookup("kubernetes-namespace"))

	viper.AutomaticEnv()
	RootCmd.Flags().AddGoFlagSet(flag.CommandLine)

}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
