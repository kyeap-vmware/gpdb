package cli

import (
	"github.com/spf13/cobra"

	"github.com/greenplum-db/gpdb/gpservice/internal/hub"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
)

func HubCmd() *cobra.Command {
	hubCmd := &cobra.Command{
		Use:    "hub",
		Short:  "Start gpservice as an agent process",
		Long:   "Start gpservice as an agent process",
		Hidden: true, // Should not be invoked by the user
		Run: func(cmd *cobra.Command, args []string) {
			err := runHubCmd()
			if err != nil {
				utils.LogErrorAndExit(err, 1)
			}
		},
	}

	return hubCmd
}

func runHubCmd() error {
	h := hub.New(serviceConfig)
	err := h.Start()
	if err != nil {
		return err
	}

	return nil
}
