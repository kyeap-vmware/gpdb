package stop

import (
	"os"
	"strings"
	"testing"

	"github.com/greenplum-db/gpdb/gpservice/test/integration/testutils"
)

func TestStopFailsWithoutSvcRunning(t *testing.T) {

	t.Run("stop agents fails when hub is not running", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		_, _ = testutils.RunGpserviceStart()
		_, _ = testutils.RunGpserviceStop("--hub")

		expectedOut := []string{
			"[ERROR]:-failed to stop agent service",
		}

		result, err := testutils.RunGpserviceStop("--agent")
		if err == nil {
			t.Errorf("\nExpected error Got: %#v", err)
		}
		if result.ExitCode != testutils.ExitCode1 {
			t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
		}
		for _, item := range expectedOut {
			if !strings.Contains(result.OutputMsg, item) {
				t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
			}
		}
	})

	t.Run("stop services fails when services are not running", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		_, _ = testutils.RunGpserviceStop()

		expectedOut := "The services may already be stopped. Use `gpservice status` to check the status"

		result, err := testutils.RunGpserviceStop()
		if err == nil {
			t.Errorf("\nExpected error Got: %#v", err)
		}
		if result.ExitCode != testutils.ExitCode1 {
			t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
		}
		if !strings.Contains(result.OutputMsg, expectedOut) {
			t.Errorf("\nExpected string: %#v \nNot found in: %#v", expectedOut, result.OutputMsg)
		}

	})

	t.Run("stop hub fails when hub is not running", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		_, _ = testutils.RunGpserviceStop()

		expectedOut := "The services may already be stopped. Use `gpservice status` to check the status"
		result, err := testutils.RunGpserviceStop("--hub")
		if err == nil {
			t.Errorf("\nExpected error Got: %#v", err)
		}
		if result.ExitCode != testutils.ExitCode1 {
			t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
		}
		if !strings.Contains(result.OutputMsg, expectedOut) {
			t.Errorf("\nExpected string: %#v \nNot found in: %#v", expectedOut, result.OutputMsg)
		}
	})

	t.Run("stop agents fails when services are not running", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		_, _ = testutils.RunGpserviceStop()

		expectedOut := "The services may already be stopped. Use `gpservice status` to check the status"
		result, err := testutils.RunGpserviceStop("--agent")
		if err == nil {
			t.Errorf("\nExpected error Got: %#v", err)
		}
		if result.ExitCode != testutils.ExitCode1 {
			t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
		}
		if !strings.Contains(result.OutputMsg, expectedOut) {
			t.Errorf("\nExpected string: %#v \nNot found in: %#v", expectedOut, result.OutputMsg)
		}
	})
}

func TestStopFailureWithoutConfig(t *testing.T) {
	TestCases := []struct {
		name        string
		cliParams   []string
		expectedOut []string
	}{
		{
			name: "stop services fails when service configuration file is not present",
			expectedOut: []string{
				"could not open service config file", "no such file or directory",
			},
		},
		{
			name: "stop hub fails when service configuration file is not present",
			cliParams: []string{
				"--hub",
			},
			expectedOut: []string{
				"could not open service config file", "no such file or directory",
			},
		},
		{
			name: "stop agents fails when service configuration file is not present",
			cliParams: []string{
				"--agent",
			},
			expectedOut: []string{
				"could not open service config file", "no such file or directory",
			},
		},
	}
	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _ = testutils.RunGpserviceStart()
			_ = testutils.CopyFile(testutils.DefaultConfigurationFile, "/tmp/config.conf")
			_ = os.RemoveAll(testutils.DefaultConfigurationFile)

			result, err := testutils.RunGpserviceStop(tc.cliParams...)
			if err == nil {
				t.Errorf("\nExpected error Got: %#v", err)
			}
			if result.ExitCode != testutils.ExitCode1 {
				t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
			}
			for _, item := range tc.expectedOut {
				if !strings.Contains(result.OutputMsg, item) {
					t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
				}
			}
			_, _ = testutils.RunGpserviceStop("--config-file", "/tmp/config.conf")
		})
	}
}

func TestStopFailsWithoutCertificates(t *testing.T) {
	TestCases := []struct {
		name        string
		cliParams   []string
		expectedOut []string
	}{
		{
			name: "stop services fails when certificates are not present",
			cliParams: []string{
				"--config-file", configCopy,
			},
			expectedOut: []string{
				"error while loading server certificate",
			},
		},
		{
			name: "stop hub fails when certificates are not present",
			cliParams: []string{
				"--hub", "--config-file", configCopy,
			},
			expectedOut: []string{
				"error while loading server certificate",
			},
		},
		{
			name: "stop agents fails when certificates are not present",
			cliParams: []string{
				"--agent", "--config-file", configCopy,
			},
			expectedOut: []string{
				"error while loading server certificate",
			},
		},
	}
	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.InitService(*hostfile, testutils.CertificateParams)
			_, _ = testutils.RunGpserviceStart()
			_ = testutils.CpCfgWithoutCertificates(configCopy)

			result, err := testutils.RunGpserviceStop(tc.cliParams...)
			if err == nil {
				t.Errorf("\nExpected error Got: %#v", err)
			}
			if result.ExitCode != testutils.ExitCode1 {
				t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
			}
			for _, item := range tc.expectedOut {
				if !strings.Contains(result.OutputMsg, item) {
					t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
				}
			}
			_, _ = testutils.RunGpserviceStop()
		})
	}
}

func TestStopFailsWithInvalidOptions(t *testing.T) {
	TestCases := []struct {
		name        string
		cliParams   []string
		expectedOut []string
	}{
		{
			name: "stop services with no value for --config-file will fail",
			cliParams: []string{
				"--config-file",
			},
			expectedOut: []string{
				"flag needs an argument: --config-file",
			},
		},
		{
			name: "stop services with non-existing file for --config-file will fail",
			cliParams: []string{
				"--config-file", "file",
			},
			expectedOut: []string{
				"no such file or directory",
			},
		},
		{
			name: "stop services with empty string for --config-file will fail",
			cliParams: []string{
				"--config-file", "",
			},
			expectedOut: []string{
				"no such file or directory",
			},
		},

		{
			name: "stop command with invalid param shows help",
			cliParams: []string{
				"invalid",
			},
			expectedOut: append([]string{
				"Error: unknown command \"invalid\" for \"gpservice stop\"",
			}, testutils.CommonHelpText...),
		},
	}

	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.InitService(*hostfile, testutils.CertificateParams)
			_, _ = testutils.RunGpserviceStart()

			result, err := testutils.RunGpserviceStop(tc.cliParams...)
			if err == nil {
				t.Errorf("\nExpected error Got: %#v", err)
			}
			if result.ExitCode != testutils.ExitCode1 {
				t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
			}
			for _, item := range tc.expectedOut {
				if !strings.Contains(result.OutputMsg, item) {
					t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
				}
			}
			_, _ = testutils.RunGpserviceStop()
		})
	}
}