package cli

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/greenplum-db/gp-common-go-libs/gplog"
	"github.com/greenplum-db/gpdb/gpservice/constants"
	"github.com/greenplum-db/gpdb/gpservice/pkg/gpservice_config"
	"github.com/spf13/cobra"
)

var (
	ConfigFilePath     string
	verbose            bool
	IsConfigured       bool
	IsGpserviceRunning bool
	Conf               *gpservice_config.Config
)

func RootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:  "gpctl",
		Long: "gpctl is a utility to manage a Greenplum Database System",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			IsConfigured = CheckGpServiceIsConfigured()
			if !IsConfigured {
				initializeLogger(cmd, "~/gpAdminLogs")
				IsGpserviceRunning = false
				return
			}

			Conf, err = gpservice_config.Read(ConfigFilePath)
			if err != nil {
				return err
			}
			IsGpserviceRunning = IsGPServiceIsRunning(Conf)
			initializeLogger(cmd, Conf.LogDir)

			return
		}}

	root.PersistentFlags().StringVar(&ConfigFilePath, "service-config-file", filepath.Join(os.Getenv("GPHOME"), constants.ConfigFileName), `Path to gpservice configuration file`)
	root.PersistentFlags().BoolVar(&verbose, "verbose", false, `Provide verbose output`)

	root.CompletionOptions.DisableDefaultCmd = true

	root.AddCommand(
		initCmd(),
	)

	return root
}

func initializeLogger(cmd *cobra.Command, logdir string) {
	// CommandPath lists the names of the called command and all of its parent commands, so this
	// turns e.g. "gp stop hub" into "gp_stop_hub" to generate a unique log file name for each command.
	logName := strings.ReplaceAll(cmd.CommandPath(), " ", "_")
	gplog.InitializeLogging(logName, logdir)

	if verbose {
		gplog.SetVerbosity(gplog.LOGVERBOSE)
	}
}
