package cli

import (
	"context"
	"fmt"

	"github.com/greenplum-db/gp-common-go-libs/gplog"
	"github.com/greenplum-db/gpdb/gpservice/idl"
	"github.com/greenplum-db/gpdb/gpservice/pkg/gpservice_config"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	StopServices = StopServicesFunc
)

func StopCmd() *cobra.Command {
	var stopHub, stopAgent bool

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stop hub and agent services",
		Args:  cobra.NoArgs,
		Example: `Stop the hub and agent services
$ gpservice stop

To stop only the hub service
$ gpservice stop --hub

To stop only the agent service
$ gpservice stop --agent
`,
		Run: func(cmd *cobra.Command, args []string) {
			err := runStopCmd(stopHub, stopAgent)
			if err != nil {
				utils.LogErrorAndExit(err, 1)
			}
		},
	}

	stopCmd.Flags().BoolVar(&stopHub, "hub", false, "Stop only the hub service")
	stopCmd.Flags().BoolVar(&stopAgent, "agent", false, "Stop only the agent service. Hub service should already be running")

	return stopCmd
}

func runStopCmd(stopHub, stopAgent bool) error {
	if stopHub {
		return stopHubService(conf)
	} else if stopAgent {
		return stopAgentService(conf)
	} else {
		return StopServices(conf)
	}
}

func stopHubService(conf *gpservice_config.Config) error {
	client, err := gpservice_config.ConnectToHub(conf)
	if err != nil {
		return err
	}

	_, err = client.Stop(context.Background(), &idl.StopHubRequest{})
	// Ignore a "hub already stopped" error
	if err != nil {
		if utils.IsGrpcServerUnavailableErr(err) {
			return utils.NewHelpErr(err, "The services may already be stopped. Use `gpservice status` to check the status.")
		}
		return fmt.Errorf("failed to stop hub service: %w", err)
	}

	gplog.Info("Hub service stopped successfully")
	return nil
}

func stopAgentService(conf *gpservice_config.Config) error {
	client, err := gpservice_config.ConnectToHub(conf)
	if err != nil {
		return err
	}

	_, err = client.StopAgents(context.Background(), &idl.StopAgentsRequest{})
	if err != nil {
		if utils.IsGrpcServerUnavailableErr(err) {
			return utils.NewHelpErr(fmt.Errorf("failed to stop agent service: %w", err), "The services may already be stopped. Use `gpservice status` to check the status.")
		}
		return fmt.Errorf("failed to stop agent service: %w", err)
	}

	gplog.Info("Agent service stopped successfully")
	return nil
}

func StopServicesFunc(conf *gpservice_config.Config) error {
	err := stopAgentService(conf)
	if err != nil {
		return err
	}

	err = stopHubService(conf)
	if err != nil {
		return err
	}

	return nil
}
