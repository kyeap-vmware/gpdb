package cli_test

import (
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/greenplum-db/gp-common-go-libs/testhelper"
	"github.com/greenplum-db/gpdb/gpservice/constants"
	"github.com/greenplum-db/gpdb/gpservice/internal/agent"
	"github.com/greenplum-db/gpdb/gpservice/internal/cli"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
	"github.com/greenplum-db/gpdb/gpservice/testutils"
	"github.com/greenplum-db/gpdb/gpservice/testutils/exectest"
)

func init() {
	exectest.RegisterMains(ValidGpsshOpt, InValidGpsshOpt, InValidGpsshOpt2, InValidGpsshOpt3)
}

func ValidGpsshOpt() {
	os.Stdout.WriteString("test output\nTEST 1234")
}

func InValidGpsshOpt() {
	os.Stdout.WriteString("test output")
}

func InValidGpsshOpt2() {
	os.Stdout.WriteString("test output\nTEST")
}

func InValidGpsshOpt3() {
	os.Stdout.WriteString("test output\nTEST TEST")
}

func TestMain(m *testing.M) {
	os.Exit(exectest.Run(m))
}

func TestInitCmd(t *testing.T) {
	t.Run("successfully configures the services with --no-tls flag", func(t *testing.T) {
		_, _, logfile := testhelper.SetupTestLogger()

		resetConf := cli.SetConf(testutils.CreateDummyServiceConfig(t))
		defer resetConf()

		platform := &testutils.MockPlatform{}
		agent.SetPlatform(platform)
		defer agent.ResetPlatform()

		utils.System.ExecCommand = exectest.NewCommand(exectest.Success)
		utils.System.OpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
			_, writer, _ := os.Pipe()

			return writer, nil
		}
		defer utils.ResetSystemFunctions()

		_, err := testutils.ExecuteCobraCommand(t, cli.RootCommand(), "init", "--no-tls", "--host", "localhost")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Created service file directory .* on all hosts`)
		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Wrote hub service file to .* on coordinator host`)
		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Wrote agent service file to .* on segment hosts`)

	})
	t.Run("successfully configures the services", func(t *testing.T) {
		_, _, logfile := testhelper.SetupTestLogger()

		resetConf := cli.SetConf(testutils.CreateDummyServiceConfig(t))
		defer resetConf()

		platform := &testutils.MockPlatform{}
		agent.SetPlatform(platform)
		defer agent.ResetPlatform()

		utils.System.ExecCommand = exectest.NewCommand(exectest.Success)
		utils.System.OpenFile = func(name string, flag int, perm os.FileMode) (*os.File, error) {
			_, writer, _ := os.Pipe()

			return writer, nil
		}
		defer utils.ResetSystemFunctions()

		args := []string{"init", "--host", "localhost", "--server-key", "tmp", "--server-certificate", "tmp", "--ca-certificate", "tmp"}

		_, err := testutils.ExecuteCobraCommand(t, cli.RootCommand(), args...)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Created service file directory .* on all hosts`)
		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Wrote hub service file to .* on coordinator host`)
		testutils.AssertLogMessage(t, logfile, `\[INFO\]:-Wrote agent service file to .* on segment hosts`)

	})
}

func TestGetUlimitSshFn(t *testing.T) {
	_, _, logile := testhelper.SetupTestLogger()
	t.Run("logs error when gpssh command execution fails", func(t *testing.T) {
		testStr := "error executing command to fetch open files limit on host"
		defer utils.ResetSystemFunctions()
		utils.System.ExecCommand = exectest.NewCommand(exectest.Failure)
		var wg sync.WaitGroup
		wg.Add(1)
		channel := make(chan cli.Response)
		cli.GetUlimitSshFn("sdw1", channel, &wg)
		testutils.AssertLogMessage(t, logile, testStr)
	})
	t.Run("logs error when gpssh output has fewer lines", func(t *testing.T) {
		testStr := "unexpected output when fetching open files limit on host"
		defer utils.ResetSystemFunctions()
		utils.System.ExecCommand = exectest.NewCommand(InValidGpsshOpt)
		var wg sync.WaitGroup
		wg.Add(1)
		channel := make(chan cli.Response)
		cli.GetUlimitSshFn("sdw1", channel, &wg)
		testutils.AssertLogMessage(t, logile, testStr)
	})
	t.Run("logs error when gpssh output has unexpected format", func(t *testing.T) {
		testStr := "unexpected output when parsing open files limit output for host"
		defer utils.ResetSystemFunctions()
		utils.System.ExecCommand = exectest.NewCommand(InValidGpsshOpt2)
		var wg sync.WaitGroup
		wg.Add(1)
		channel := make(chan cli.Response)
		cli.GetUlimitSshFn("sdw1", channel, &wg)
		testutils.AssertLogMessage(t, logile, testStr)
	})
	t.Run("logs error when gpssh output fails to convert to integer", func(t *testing.T) {
		testStr := "unexpected output when converting open files limit value for host"
		defer utils.ResetSystemFunctions()
		utils.System.ExecCommand = exectest.NewCommand(InValidGpsshOpt3)
		var wg sync.WaitGroup
		wg.Add(1)
		channel := make(chan cli.Response)
		cli.GetUlimitSshFn("sdw1", channel, &wg)
		testutils.AssertLogMessage(t, logile, testStr)
	})
	t.Run("logs error when gpssh output fails to convert to integer", func(t *testing.T) {
		defer utils.ResetSystemFunctions()
		utils.System.ExecCommand = exectest.NewCommand(ValidGpsshOpt)
		var wg sync.WaitGroup
		wg.Add(1)
		channel := make(chan cli.Response)
		go cli.GetUlimitSshFn("sdw1", channel, &wg)
		go func() {
			wg.Wait()
			close(channel)
		}()

		for result := range channel {
			if result.Ulimit != 1234 {
				t.Fatalf("got ulimit:%d, expected:1234", result.Ulimit)
			}
		}
	})
}
func TestCheckOpenFilesLimitOnHosts(t *testing.T) {
	_, _, logile := testhelper.SetupTestLogger()
	t.Run("prints warning if fails to execute Ulimit command", func(t *testing.T) {
		testStr := "test error string"
		defer utils.ResetSystemFunctions()
		utils.ExecuteAndGetUlimit = func() (int, error) {
			return -1, fmt.Errorf(testStr)
		}
		cli.CheckOpenFilesLimitOnHosts(nil)
		testutils.AssertLogMessage(t, logile, testStr)
	})
	t.Run("prints warning  Ulimit is lower than required on coordinator", func(t *testing.T) {
		testStr := "For proper functioning make sure limit is set properly for system and services before starting gp services."
		defer utils.ResetSystemFunctions()
		utils.ExecuteAndGetUlimit = func() (int, error) {
			return constants.OsOpenFiles - 1, nil
		}
		cli.CheckOpenFilesLimitOnHosts(nil)
		testutils.AssertLogMessage(t, logile, testStr)
	})
	t.Run("prints warning  Ulimit is lower than required on remote host", func(t *testing.T) {
		testStr := "For proper functioning make sure limit is set properly for system and services before starting gp services."
		defer utils.ResetSystemFunctions()
		defer func() { cli.GetUlimitSsh = cli.GetUlimitSshFn }()
		utils.ExecuteAndGetUlimit = func() (int, error) {
			return constants.OsOpenFiles + 1, nil
		}
		cli.GetUlimitSsh = func(hostname string, channel chan cli.Response, wg *sync.WaitGroup) {
			defer wg.Done()
			channel <- cli.Response{Hostname: "localhost", Ulimit: constants.OsOpenFiles - 1}
		}
		cli.CheckOpenFilesLimitOnHosts([]string{"localhost"})
		testutils.AssertLogMessage(t, logile, testStr)
	})
}
