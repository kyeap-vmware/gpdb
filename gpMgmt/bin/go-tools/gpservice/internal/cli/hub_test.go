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

func TestHubCmd(t *testing.T) {
	t.Run("hub command starts the hub server", func(t *testing.T) {
		config := testutils.CreateDummyServiceConfig(t)
		config.HubPort = testutils.GetPort(t)
		resetConf := cli.SetConf(config)
		defer resetConf()

		ctx, cancel := context.WithCancel(context.Background())
		go testutils.ExecuteCobraCommandContext(t, ctx, cli.HubCmd()) //nolint

		testutils.CheckGRPCServerRunning(t, config.HubPort)
		cancel()
	})

	t.Run("returns error when fails to start the hub server", func(t *testing.T) {
		_, _, logfile := testhelper.SetupTestLogger()

		port, cleanup := testutils.GetAndListenOnPort(t)
		defer cleanup()

		config := testutils.CreateDummyServiceConfig(t)
		config.HubPort = port
		resetConf := cli.SetConf(config)
		defer resetConf()

		utils.System.OSExit = func(code int) {}
		defer utils.ResetSystemFunctions()

		testutils.ExecuteCobraCommand(t, cli.HubCmd()) // nolint

		expected := fmt.Sprintf(`\[ERROR\]:-could not listen on port %d:`, port)
		testutils.AssertLogMessage(t, logfile, expected)
	})
}
