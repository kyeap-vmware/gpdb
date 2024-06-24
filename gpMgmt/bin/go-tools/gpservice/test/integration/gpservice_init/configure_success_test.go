package configure

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/greenplum-db/gpdb/gpservice/constants"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
	"github.com/greenplum-db/gpdb/gpservice/test/integration/testutils"
)

func TestConfigureHelp(t *testing.T) {
	Testcases := []struct {
		name        string
		cliParams   []string
		expectedOut []string
	}{
		{
			name:        "service configure shows help with --help",
			cliParams:   []string{"--help"},
			expectedOut: helpTxt,
		},
		{
			name:        "service configure shows help with -h",
			cliParams:   []string{"-h"},
			expectedOut: helpTxt,
		},
	}
	for _, tc := range Testcases {
		t.Run(tc.name, func(t *testing.T) {
			// Running the gp configure command with help options
			result, err := testutils.RunGPServiceInit(false, tc.cliParams...)
			// check for command result
			if err != nil {
				t.Errorf("\nUnexpected error: %#v", err)
			}
			if result.ExitCode != 0 {
				t.Errorf("\nExpected: %v \nGot: %v", 0, result.ExitCode)
			}
			for _, item := range tc.expectedOut {
				if !strings.Contains(result.OutputMsg, item) {
					t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
				}
			}
		})
	}
}

func TestConfigureSuccess(t *testing.T) {
	hosts := testutils.GetHostListFromFile(*hostfile)
	agentFile := fmt.Sprintf("%s/%s_%s.%s", defaultServiceDir, constants.DefaultServiceName, "agent", serviceExt)
	hubFile := fmt.Sprintf("%s/%s_%s.%s", defaultServiceDir, constants.DefaultServiceName, "hub", serviceExt)

	testServiceCreation := func(t *testing.T, cliParams []string) {
		runConfigureAndCheckOutput(t, cliParams)

		expectedTestConfig := defaultGPConf
		expectedTestConfig.Hostnames = hosts
		actualTestConfig := testutils.ParseConfig(testutils.DefaultConfigurationFile)
		if !reflect.DeepEqual(expectedTestConfig, actualTestConfig) {
			t.Fatalf("Got:%+v Expected %+v", expectedTestConfig, actualTestConfig)
		}

		testutils.FilesExistOnHub(t, hubFile)
		testutils.FilesExistsOnAgents(t, agentFile, hosts)
	}

	t.Run("configure service with --host option", func(t *testing.T) {
		var cliParams []string
		for _, h := range hosts {
			cliParams = append(cliParams, "--host", h)
		}
		testServiceCreation(t, cliParams)

	})

	t.Run("configure service with --hostfile option", func(t *testing.T) {
		cliParams := []string{"--hostfile", *hostfile}
		testServiceCreation(t, cliParams)

	})

	t.Run("configure service with host and agent_port option", func(t *testing.T) {
		cliParams := []string{
			"--hostfile", *hostfile,
			"--agent-port", "8001",
		}
		runConfigureAndCheckOutput(t, cliParams)

		expectedTestConfig := defaultGPConf
		expectedTestConfig.Hostnames = hosts
		actualTestConfig := testutils.ParseConfig(testutils.DefaultConfigurationFile)
		expectedTestConfig.AgentPort = 8001
		if !reflect.DeepEqual(expectedTestConfig, actualTestConfig) {
			t.Fatalf("Got:%+v Expected %+v", expectedTestConfig, actualTestConfig)
		}
	})

	t.Run("configure service with host and hub_port option", func(t *testing.T) {
		cliParams := []string{
			"--hostfile", *hostfile,
			"--hub-port", "8001",
		}
		runConfigureAndCheckOutput(t, cliParams)

		expectedTestConfig := defaultGPConf
		expectedTestConfig.Hostnames = hosts
		actualTestConfig := testutils.ParseConfig(testutils.DefaultConfigurationFile)
		expectedTestConfig.HubPort = 8001
		if !reflect.DeepEqual(expectedTestConfig, actualTestConfig) {
			t.Fatalf("Got:%+v Expected %+v", expectedTestConfig, actualTestConfig)
		}
	})

	t.Run("configure service with server and client certificates", func(t *testing.T) {
		cliParams := []string{
			"--ca-certificate", "/tmp/certificates/ca-cert.pem",
			"--server-certificate", "/tmp/certificates/server-cert.pem",
			"--server-key", "/tmp/certificates/server-key.pem",
			"--hostfile", *hostfile,
		}
		testServiceCreation(t, cliParams)
	})

	t.Run("configure service with verbose option", func(t *testing.T) {
		cliParams := []string{
			"--hostfile", *hostfile,
			"--verbose",
		}
		testServiceCreation(t, cliParams)
	})

	t.Run("configure service with config-file option", func(t *testing.T) {
		configFile := testutils.DefaultConfigurationFile
		cliParams := []string{
			"--hostfile", *hostfile,
			"--config-file", configFile,
		}
		testServiceCreation(t, cliParams)
	})

	t.Run("configure service with changing gpHome value", func(t *testing.T) {
		cliParams := []string{
			"--hostfile", *hostfile,
			"--gphome", os.Getenv("GPHOME"),
		}
		testServiceCreation(t, cliParams)
	})

	t.Run("configure service with log_dir option", func(t *testing.T) {
		logDir := "/tmp/log"
		_ = os.MkdirAll(logDir, 0777)

		cliParams := []string{
			"--hostfile", *hostfile,
			"--log-dir", logDir,
		}
		runConfigureAndCheckOutput(t, cliParams)

		expectedTestConfig := defaultGPConf
		expectedTestConfig.Hostnames = hosts
		actualTestConfig := testutils.ParseConfig(testutils.DefaultConfigurationFile)
		expectedTestConfig.LogDir = logDir
		if !reflect.DeepEqual(expectedTestConfig, actualTestConfig) {
			t.Fatalf("Got:%+v Expected %+v", expectedTestConfig, actualTestConfig)
		}

	})

	t.Run("configure service with service-name option", func(t *testing.T) {
		svcName := "dummySvc"
		cliParams := []string{
			"--hostfile", *hostfile,
			"--service-name", svcName,
		}
		runConfigureAndCheckOutput(t, cliParams)

		expectedTestConfig := defaultGPConf
		expectedTestConfig.Hostnames = hosts
		actualTestConfig := testutils.ParseConfig(testutils.DefaultConfigurationFile)
		expectedTestConfig.ServiceName = svcName
		if !reflect.DeepEqual(expectedTestConfig, actualTestConfig) {
			t.Fatalf("Got:%+v Expected %+v", expectedTestConfig, actualTestConfig)
		}
	})

	t.Run("configure service with --no-tls option", func(t *testing.T) {
		cliParams := []string{
			"--hostfile", *hostfile,
			"--no-tls",
		}
		result, err := testutils.RunGPServiceInit(false, cliParams...)
		if err != nil {
			t.Errorf("\nUnexpected error: %#v", err)
		}
		if result.ExitCode != 0 {
			t.Errorf("\nExpected: %v \nGot: %v", 0, result.ExitCode)
		}
		for _, item := range expectedOutput {
			if !strings.Contains(result.OutputMsg, item) {
				t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
			}
		}
		expectedTestConfig := defaultGPConf
		expectedTestConfig.Hostnames = hosts
		actualTestConfig := testutils.ParseConfig(testutils.DefaultConfigurationFile)
		expectedTestConfig.Credentials = &utils.GpCredentials{
			TlsEnabled: false,
		}
		if !reflect.DeepEqual(expectedTestConfig, actualTestConfig) {
			t.Fatalf("Got:%+v Expected %+v", expectedTestConfig, actualTestConfig)
		}
	})

}

func runConfigureAndCheckOutput(t *testing.T, input []string) {
	// Running the gp configure command with input params
	result, err := testutils.RunGPServiceInit(true, input...)
	// check for command result
	if err != nil {
		t.Errorf("\nUnexpected error: %#v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("\nExpected: %v \nGot: %v", 0, result.ExitCode)
	}
	for _, item := range expectedOutput {
		if !strings.Contains(result.OutputMsg, item) {
			t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
		}
	}
}
