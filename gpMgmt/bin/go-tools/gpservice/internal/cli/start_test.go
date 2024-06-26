package cli_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/greenplum-db/gp-common-go-libs/testhelper"
	"github.com/greenplum-db/gpdb/gpservice/idl"
	"github.com/greenplum-db/gpdb/gpservice/idl/mock_idl"
	"github.com/greenplum-db/gpdb/gpservice/internal/cli"
	"github.com/greenplum-db/gpdb/gpservice/pkg/gpservice_config"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
	"github.com/greenplum-db/gpdb/gpservice/testutils"
	"github.com/greenplum-db/gpdb/gpservice/testutils/exectest"
)

func TestStartCmd(t *testing.T) {
	t.Run("starts only the hub service", func(t *testing.T) {
		_, _, logfile := testhelper.SetupTestLogger()

		resetConf := cli.SetConf(testutils.CreateDummyServiceConfig(t))
		defer resetConf()

		utils.System.ExecCommand = exectest.NewCommand(exectest.Success)
		utils.System.GetHostName = func() (name string, err error) {
			return "cdw", nil
		}
		defer utils.ResetSystemFunctions()

		utils.SetNewHealthClient(testutils.NewMockHealthClient(grpc_health_v1.HealthCheckResponse_SERVING, nil))
		defer utils.ResetNewHealthClient()

		_, err := testutils.ExecuteCobraCommand(t, cli.StartCmd(), "--hub")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Hub service started successfully`)
		testutils.AssertLogMessageNotPresent(t, logfile, `\[INFO\]:-Agent service started successfully`)
	})

	t.Run("starts only the agent service", func(t *testing.T) {
		_, _, logfile := testhelper.SetupTestLogger()

		resetConf := cli.SetConf(testutils.CreateDummyServiceConfig(t))
		defer resetConf()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock_idl.NewMockHubClient(ctrl)
		client.EXPECT().StartAgents(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(&idl.StartAgentsReply{}, nil)
		gpservice_config.SetConnectToHub(client)
		defer gpservice_config.ResetConfigFunctions()

		_, err := testutils.ExecuteCobraCommand(t, cli.StartCmd(), "--agent")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Agent service started successfully`)
		testutils.AssertLogMessageNotPresent(t, logfile, `\[INFO\]:-Hub service started successfully`)
	})

	t.Run("starts both hub and agent", func(t *testing.T) {
		_, _, logfile := testhelper.SetupTestLogger()

		resetConf := cli.SetConf(testutils.CreateDummyServiceConfig(t))
		defer resetConf()

		utils.System.ExecCommand = exectest.NewCommand(exectest.Success)
		defer utils.ResetSystemFunctions()

		utils.SetNewHealthClient(testutils.NewMockHealthClient(grpc_health_v1.HealthCheckResponse_SERVING, nil))
		defer utils.ResetNewHealthClient()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock_idl.NewMockHubClient(ctrl)
		client.EXPECT().StartAgents(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(&idl.StartAgentsReply{}, nil)
		gpservice_config.SetConnectToHub(client)
		defer gpservice_config.ResetConfigFunctions()

		_, err := testutils.ExecuteCobraCommand(t, cli.StartCmd())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Agent service started successfully`)
		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Hub service started successfully`)
	})

	t.Run("returns error when fails to start the hub service", func(t *testing.T) {
		_, _, logfile := testhelper.SetupTestLogger()

		resetConf := cli.SetConf(testutils.CreateDummyServiceConfig(t))
		defer resetConf()

		utils.System.ExecCommand = exectest.NewCommand(exectest.Failure)
		utils.System.OSExit = func(code int) {}
		defer utils.ResetSystemFunctions()

		testutils.ExecuteCobraCommand(t, cli.StartCmd()) // nolint

		expected := `\[ERROR\]:-failed to start hub service:`
		testutils.AssertLogMessage(t, logfile, expected)
	})

	t.Run("returns error when fails to start the agent service", func(t *testing.T) {
		_, _, logfile := testhelper.SetupTestLogger()

		resetConf := cli.SetConf(testutils.CreateDummyServiceConfig(t))
		defer resetConf()

		utils.System.ExecCommand = exectest.NewCommand(exectest.Success)
		utils.System.OSExit = func(code int) {}
		defer utils.ResetSystemFunctions()

		utils.SetNewHealthClient(testutils.NewMockHealthClient(grpc_health_v1.HealthCheckResponse_SERVING, nil))
		defer utils.ResetNewHealthClient()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		expectedErr := errors.New("error")
		client := mock_idl.NewMockHubClient(ctrl)
		client.EXPECT().StartAgents(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(&idl.StartAgentsReply{}, expectedErr)
		gpservice_config.SetConnectToHub(client)
		defer gpservice_config.ResetConfigFunctions()

		testutils.ExecuteCobraCommand(t, cli.StartCmd()) // nolint

		expected := fmt.Sprintf(`\[ERROR\]:-failed to start agent service: %v`, expectedErr)
		testutils.AssertLogMessage(t, logfile, expected)
	})

	t.Run("shows the status when run in verbose mode", func(t *testing.T) {
		testhelper.SetupTestLogger()

		resetConf := cli.SetConf(testutils.CreateDummyServiceConfig(t))
		defer resetConf()

		utils.System.ExecCommand = exectest.NewCommand(exectest.Success)
		utils.System.GetHostName = func() (name string, err error) {
			return "cdw", nil
		}
		defer utils.ResetSystemFunctions()

		utils.SetNewHealthClient(testutils.NewMockHealthClient(grpc_health_v1.HealthCheckResponse_SERVING, nil))
		defer utils.ResetNewHealthClient()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock_idl.NewMockHubClient(ctrl)
		client.EXPECT().StartAgents(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(&idl.StartAgentsReply{}, nil)
		client.EXPECT().StatusAgents(
			gomock.Any(),
			gomock.Any(),
		).Return(&idl.StatusAgentsReply{
			Statuses: []*idl.ServiceStatus{
				{Role: "Agent", Host: "sdw1", Status: "running", Uptime: "5H", Pid: 123},
			},
		}, nil)
		gpservice_config.SetConnectToHub(client)
		defer gpservice_config.ResetConfigFunctions()

		buffer, writer, resetStdout := testutils.CaptureStdout(t)
		defer resetStdout()

		rootCmd := cli.RootCommand()
		rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {}
		_, err := testutils.ExecuteCobraCommand(t, rootCmd, "start", "--verbose")
		writer.Close()
		stdout := <-buffer

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedStdout := `ROLE      HOST      STATUS       PID       UPTIME
Hub       cdw       not running  0         
ROLE      HOST      STATUS    PID       UPTIME
Agent     sdw1      running   123       5H
`
		if stdout != expectedStdout {
			t.Fatalf("got %s, want %s", stdout, expectedStdout)
		}
	})
}
