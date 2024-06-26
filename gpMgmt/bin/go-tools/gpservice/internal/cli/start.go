package cli

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/greenplum-db/gp-common-go-libs/gplog"
	"github.com/greenplum-db/gpdb/gpservice/idl"
	"github.com/greenplum-db/gpdb/gpservice/pkg/gpservice_config"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
)

func StartCmd() *cobra.Command {
	var startHub, startAgent bool

	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start hub and agent services",
		Args:  cobra.NoArgs,
		Example: `Start the hub and agent services
$ gpservice start

To start only the hub service
$ gpservice start --hub

To start only the agent service
$ gpservice start --agent
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := runStartCmd(startHub, startAgent)
			if err != nil {
				utils.LogErrorAndExit(err, 1)
			}
		},
	}

	startCmd.Flags().BoolVar(&startHub, "hub", false, "Start only the hub service")
	startCmd.Flags().BoolVar(&startAgent, "agent", false, "Start only the agent service. Hub service should already be running")

	return startCmd
}

func runStartCmd(startHub, startAgent bool) error {
	if startHub {
		return startHubService(serviceConfig)
	} else if startAgent {
		return startAgentService(serviceConfig)
	} else {
		return StartServices(serviceConfig)
	}
}

func startHubService(conf *gpservice_config.Config) error {
	errPrefix := "failed to start hub service"
	out, err := platform.GetStartHubCommand(conf.ServiceName).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s, %w", errPrefix, out, err)
	}

	credentials, err := conf.Credentials.LoadClientCredentials()
	if err != nil {
		return fmt.Errorf("%s: %w", errPrefix, err)
	}

	address := net.JoinHostPort("localhost", strconv.Itoa(conf.HubPort))
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(credentials))
	if err != nil {
		return fmt.Errorf("%s: %w", errPrefix, err)
	}

	err = utils.CheckGRPCServerHealth(conn)
	if err != nil {
		return fmt.Errorf("%s: %w", errPrefix, err)
	}

	gplog.Info("Hub service started successfully")
	if verbose {
		if verbose {
			status, _ := getHubStatus(conf)
			displayServiceStatus(os.Stdout, status)
		}
	}

	return nil
}

func startAgentService(conf *gpservice_config.Config) error {
	client, err := gpservice_config.ConnectToHub(conf)
	if err != nil {
		return err
	}

	_, err = client.StartAgents(context.Background(), &idl.StartAgentsRequest{})
	if err != nil {
		if utils.IsGrpcServerUnavailableErr(err) {
			return utils.NewHelpErr(err, "Ensure the hub service is running. If it is not, start the service using the 'gpservice start --hub' command.")
		}
		return fmt.Errorf("failed to start agent service: %w", err)
	}

	gplog.Info("Agent service started successfully")
	if verbose {
		status, _ := getAgentStatus(conf)
		displayServiceStatus(os.Stdout, status)
	}

	return nil
}

func StartServices(conf *gpservice_config.Config) error {
	err := startHubService(conf)
	if err != nil {
		return err
	}

	err = startAgentService(conf)
	if err != nil {
		return err
	}

	return nil
}
