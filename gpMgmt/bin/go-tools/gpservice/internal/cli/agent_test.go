package cli_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/greenplum-db/gp-common-go-libs/testhelper"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
	"github.com/greenplum-db/gpdb/gpservice/testutils"

	"github.com/greenplum-db/gpdb/gpservice/internal/cli"
)

func TestAgentCmd(t *testing.T) {
	t.Run("agent command starts the agent server", func(t *testing.T) {
		config := testutils.CreateDummyServiceConfig(t)
		config.AgentPort = testutils.GetPort(t)
		resetConf := cli.SetConf(config)
		defer resetConf()

		ctx, cancel := context.WithCancel(context.Background())
		go testutils.ExecuteCobraCommandContext(t, ctx, cli.AgentCmd()) // nolint

		testutils.CheckGRPCServerRunning(t, config.AgentPort)
		cancel()
	})

	t.Run("returns error when fails to start the agent server", func(t *testing.T) {
		_, _, logfile := testhelper.SetupTestLogger()

		port, cleanup := testutils.GetAndListenOnPort(t)
		defer cleanup()

		config := testutils.CreateDummyServiceConfig(t)
		config.AgentPort = port
		resetConf := cli.SetConf(config)
		defer resetConf()

		utils.System.OSExit = func(code int) {}
		defer utils.ResetSystemFunctions()

		testutils.ExecuteCobraCommand(t, cli.AgentCmd()) // nolint

		expected := fmt.Sprintf(`\[ERROR\]:-could not listen on port %d:`, port)
		testutils.AssertLogMessage(t, logfile, expected)
	})
}
