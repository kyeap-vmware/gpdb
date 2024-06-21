package delete

import (
	"strings"
	"testing"

	"github.com/greenplum-db/gpdb/gpservice/internal/platform"
	"github.com/greenplum-db/gpdb/gpservice/test/integration/testutils"
)

func TestDeleteSuccess(t *testing.T) {
	hosts := testutils.GetHostListFromFile(*hostfile)

	t.Run("delete services successfully when services are not running", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		cliParams := []string{"services"}
		expectedOut := []string{"Successfully deleted service configuration file", "Removed hub service file", "Successfully removed agent service file"}

		runDeleteCmdAndCheckOutput(t, cliParams, expectedOut)
		// check if service is running
		for _, svc := range []string{"gpservice_hub", "gpservice_agent"} {
			hostList := hosts[:1]
			if svc == "gpservice_agent" {
				hostList = hosts
			}
			for _, host := range hostList {
				status, _ := testutils.GetSvcStatusOnHost(p.(platform.GpPlatform), svc, host)
				testutils.VerifySvcNotRunning(t, status.OutputMsg)
			}
		}

	})
	t.Run("delete services successfully when services are started", func(t *testing.T) {
		testutils.InitService(*hostfile, testutils.CertificateParams)
		_, _ = testutils.RunGpserviceStart("services")
		cliParams := []string{"services"}
		expectedOut := []string{"Successfully deleted service configuration", "Removed hub service file", "Successfully removed agent service file"}

		runDeleteCmdAndCheckOutput(t, cliParams, expectedOut)
		// check if service is running
		for _, svc := range []string{"gpservice_hub", "gpservice_agent"} {
			hostList := hosts[:1]
			if svc == "gpservice_agent" {
				hostList = hosts
			}
			for _, host := range hostList {
				status, _ := testutils.GetSvcStatusOnHost(p.(platform.GpPlatform), svc, host)
				testutils.VerifySvcNotRunning(t, status.OutputMsg)
			}
		}

	})
	t.Run("delete services after gp configure with --service-name param", func(t *testing.T) {
		_, _ = testutils.RunConfigure(true, []string{
			"--hostfile", *hostfile,
			"--service-name", "dummySvc",
		}...)
		_, _ = testutils.RunGpserviceStart("services")
		cliParams := []string{"services"}
		expectedOut := []string{"-Successfully deleted service configuration", "Removed hub service file", "Successfully removed agent service file"}

		runDeleteCmdAndCheckOutput(t, cliParams, expectedOut)
		// check if service is running
		for _, svc := range []string{"gpservice_hub", "gpservice_agent"} {
			hostList := hosts[:1]
			if svc == "gpservice_agent" {
				hostList = hosts
			}
			for _, host := range hostList {
				status, _ := testutils.GetSvcStatusOnHost(p.(platform.GpPlatform), svc, host)
				testutils.VerifySvcNotRunning(t, status.OutputMsg)
			}
		}
	})
}
func runDeleteCmdAndCheckOutput(t *testing.T, input []string, output []string) {
	result, err := testutils.RunGpServiceDelete(input...)
	// check for command result
	if err != nil {
		t.Errorf("\nUnexpected error: %#v", err)
	}
	if result.ExitCode != 0 {
		t.Errorf("\nExpected: %v \nGot: %v", 0, result.ExitCode)
	}
	for _, item := range output {
		if !strings.Contains(result.OutputMsg, item) {
			t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
		}
	}
}
