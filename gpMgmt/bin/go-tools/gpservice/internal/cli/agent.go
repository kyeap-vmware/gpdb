package cli

import (
	"github.com/greenplum-db/gpdb/gpservice/internal/agent"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
	"github.com/spf13/cobra"
)

func AgentCmd() *cobra.Command {
	agentCmd := &cobra.Command{
		Use:    "agent",
		Short:  "Start gpservice as an agent process",
		Long:   "Start gpservice as an agent process",
		Hidden: true, // Should not be invoked by the user
		Run: func(cmd *cobra.Command, args []string) {
			err := runAgentCmd()
			if err != nil {
				utils.LogErrorAndExit(err, 1)
			}
		},
	}

	return agentCmd
}

func runAgentCmd() error {
	agentConf := agent.Config{
		Port:        conf.AgentPort,
		ServiceName: conf.ServiceName,
		GpHome:      conf.GpHome,
		Credentials: conf.Credentials,
		LogDir:      conf.LogDir,
	}
	a := agent.New(agentConf)

	err := a.Start()
	if err != nil {
		return err
	}

	return nil
}
