package start

import (
	"os"
	"strings"
	"testing"

	"github.com/greenplum-db/gpdb/gpservice/test/integration/testutils"
)

func TestStartFailWithoutConfig(t *testing.T) {
	t.Run("starting services without configuration file will fail", func(t *testing.T) {
		_ = os.RemoveAll(testutils.DefaultConfigurationFile)
		expectedOut := []string{
			"could not open service config file",
			"no such file or directory",
		}
		// start services
		result, err := testutils.RunGpserviceStart()
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

	t.Run("starting hub without service file", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		testutils.DisableandDeleteServiceFiles(p)

		expectedOut := "failed to start hub service"

		// start hub
		result, err := testutils.RunGpserviceStart("--hub")
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

	t.Run("starting hub without certificates", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		err := testutils.CpCfgWithoutCertificates(configCopy)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expectedOut := "error while loading server certificate"

		result, err := testutils.RunGpserviceStart("--hub", "--config-file", configCopy)
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

	t.Run("starting agents without certificates", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		_ = testutils.CpCfgWithoutCertificates(configCopy)

		expectedOut := "error while loading server certificate"

		// start agents
		result, err := testutils.RunGpserviceStart("--agent", "--config-file", configCopy)
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

	t.Run("starting services without certificates", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		_ = testutils.CpCfgWithoutCertificates(configCopy)

		expectedOut := "error while loading server certificate"
		// start agents
		result, err := testutils.RunGpserviceStart("--config-file", configCopy)
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

func TestStartGlobalFlagsFailures(t *testing.T) {
	failTestCases := []struct {
		name        string
		cliParams   []string
		expectedOut string
	}{
		{
			name: "starting services with no value for --config-file will fail",
			cliParams: []string{
				"--config-file",
			},
			expectedOut: "flag needs an argument: --config-file",
		},
		{
			name: "starting services with non-existing file for --config-file will fail",
			cliParams: []string{
				"--config-file", "file",
			},
			expectedOut: "no such file or directory",
		},
		{
			name: "starting services with empty string for --config-file will fail",
			cliParams: []string{
				"--config-file", "",
			},
			expectedOut: "no such file or directory",
		},
		{
			name: "starting agents without starting hub will fail",
			cliParams: []string{
				"--agent",
			},
			expectedOut: "Ensure the hub service is running",
		},
	}

	for _, tc := range failTestCases {
		t.Run(tc.name, func(t *testing.T) {
			testutils.DisableandDeleteServiceFiles(p)
			testutils.InitService(*hostfile, testutils.CertificateParams)
			_, _ = testutils.RunGpserviceStop()

			result, err := testutils.RunGpserviceStart(tc.cliParams...)
			if err == nil {
				t.Errorf("\nExpected error Got: %#v", err)
			}
			if result.ExitCode != testutils.ExitCode1 {
				t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
			}
			if !strings.Contains(result.OutputMsg, tc.expectedOut) {
				t.Errorf("\nExpected string: %#v \nNot found in: %#v", tc.expectedOut, result.OutputMsg)
			}

		})
		_, _ = testutils.RunGpserviceStop()
	}
}
func TestStartInvalidFlagsFailures(t *testing.T) {
	t.Run("checking status of services command with invalid param shows help", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		cliParams := []string{
			"invalid",
		}

		expectedOut := append([]string{
			"Error: unknown command \"invalid\" for \"gpservice start\"",
		}, testutils.CommonHelpText...)

		result, err := testutils.RunGpserviceStart(cliParams...)
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
}

func TestStartInvalidCertificatePathFailures(t *testing.T) {
	t.Run("checking status of services command with invalid param shows help", func(t *testing.T) {
		cliParams := []string{
			"--host", testutils.DefaultHost,
			"--server-certificate", "/invalid/path/to/server-cert.pem",
			"--server-key", "/invalid/path/to/server-key.pem",
			"--ca-certificate", "/invalid/path/to/ca-cert.pem",
		}

		result, err := testutils.RunGPServiceInit(false, cliParams...)
		if err != nil {
			t.Errorf("\nUnexpected error: %#v", err)
		}

		expectedOut := "[ERROR]:-failed to start hub service: error while loading server certificate"

		result, err = testutils.RunGpserviceStart()
		if err == nil {
			t.Errorf("\nExpected error Got: %#v", err)
		}
		if result.ExitCode != testutils.ExitCode1 {
			t.Errorf("\nExpected: %#v \nGot: %v", testutils.ExitCode1, result.ExitCode)
		}
		if !strings.Contains(result.OutputMsg, expectedOut) {
			t.Fatalf("got %q, want %q", result.OutputMsg, expectedOut)
		}

	})
}
