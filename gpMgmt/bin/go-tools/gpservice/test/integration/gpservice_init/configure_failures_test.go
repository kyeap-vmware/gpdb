package configure

import (
	"os"
	"strings"
	"testing"

	"github.com/greenplum-db/gpdb/gpservice/test/integration/testutils"
)

func TestConfigureFailure(t *testing.T) {
	// creating empty hostfile for test
	_ = os.WriteFile(mockHostFile, []byte(""), 0644)

	var ConfigureFailTestCases = []struct {
		name        string
		cliParams   []string
		expectedOut []string
	}{
		{
			name:      "configure service with empty value for --host option",
			cliParams: []string{"--host", ""},
			expectedOut: []string{
				"please provide a valid input host name",
			},
		},
		{
			name:      "configure service with no value for --host option",
			cliParams: []string{"--host"},
			expectedOut: []string{
				"flag needs an argument: --host",
			},
		},
		{
			name:      "configure service with empty file for --hostfile option",
			cliParams: []string{"--hostfile", "hostlist"},
			expectedOut: []string{
				"no host name found, please provide a valid input host name",
			},
		},
		{
			name:      "configure service with no value for --hostfile option",
			cliParams: []string{"--hostfile"},
			expectedOut: []string{
				"flag needs an argument: --hostfile",
			},
		},
		{
			name:      "configure service with non-existing host for --host option",
			cliParams: []string{"--host", "host"},
			expectedOut: []string{
				"could not copy gp.conf file to segment hosts",
			},
		},
		{
			name:      "configure service with one valid host and invalid host for --host option",
			cliParams: []string{"--host", testutils.DefaultHost, "--host", "invalid"},
			expectedOut: []string{
				"could not copy gp.conf file to segment hosts",
			},
		},
		{
			name: "configure service without any option",
			expectedOut: []string{
				"at least one hostname must be provided using either --host or --hostfile",
			},
		},
		{
			name:      "configure service with invalid option",
			cliParams: []string{"--invalid"},
			expectedOut: []string{
				"unknown flag: --invalid",
			},
		},
		{
			name: "configure service with both host and hostfile options",
			cliParams: []string{"--host", testutils.DefaultHost,
				"--hostfile", "abc"},
			expectedOut: []string{
				"if any flags in the group [host hostfile] are set none of the others can be; [host hostfile] were all set",
			},
		},
		{
			name: "configure service with string value for --agent-port option",
			cliParams: []string{"--host", testutils.DefaultHost,
				"--agent-port", "abc"},
			expectedOut: []string{
				"invalid argument",
			},
		},
		{
			name: "configure service with string value for --hub-port option",
			cliParams: []string{"--host", testutils.DefaultHost,
				"--hub-port", "abc"},
			expectedOut: []string{
				"invalid argument",
			},
		},
		{
			name: "configure service with no value for --agent-port option",
			cliParams: []string{"--host", testutils.DefaultHost,
				"--agent-port"},
			expectedOut: []string{
				"flag needs an argument: --agent-port",
			},
		},
		{
			name: "configure service with no value for --hub-port option",
			cliParams: []string{"--host", testutils.DefaultHost,
				"--hub-port"},
			expectedOut: []string{
				"flag needs an argument: --hub-port",
			},
		},
		{
			name: "configure service with no value for log-dir option",
			cliParams: []string{
				"--host", testutils.DefaultHost,
				"--log-dir",
			},
			expectedOut: []string{
				"flag needs an argument: --log-dir",
			},
		},
		{
			name: "configure fails when value for both --agent-port and --hub-port are same",
			cliParams: []string{
				"--host", testutils.DefaultHost,
				"--agent-port", "2000",
				"--hub-port", "2000",
			},
			expectedOut: []string{
				"[ERROR]:-hub port and agent port must be different",
			},
		},
		{
			name: "configure service fails when --gpHome value is invalid",
			cliParams: []string{
				"--host", testutils.DefaultHost,
				"--gphome", "invalid",
			},
			expectedOut: []string{
				"could not create configuration file invalid/gp.conf",
			},
		},
		{
			name: "configure service fails when --gpHome value is empty",
			cliParams: []string{
				"--host", testutils.DefaultHost,
				"--gphome", "",
			},
			expectedOut: []string{
				"not a valid gpHome found",
			},
		},
		{
			name: "configure service fails when no value given for --gpHome",
			cliParams: []string{
				"--host", testutils.DefaultHost,
				"--gphome",
			},
			expectedOut: []string{
				"flag needs an argument: --gphome",
			},
		},
		{
			name: "configure fails when value for both --no-tls and --server-certificate params are used",
			cliParams: append([]string{
				"--host", testutils.DefaultHost,
				"--no-tls",
			}, testutils.CertificateParams...),
			expectedOut: []string{
				"[ERROR]:-cannot specify --no-tls flag and specify certificates together",
			},
		},
	}

	for _, tc := range ConfigureFailTestCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := testutils.RunGPServiceInit(true, tc.cliParams...)
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

			// check for panic error
			if strings.Contains(result.OutputMsg, "panic") {
				t.Errorf("\nUnexpected string: %#v \nFound in: %#v", "panic", result.OutputMsg)
			}
		})
	}
}

// TODO: Below cases needs to be corrected if the error messages are changed for individual params
func TestConfigureCertificateFailure(t *testing.T) {
	// creating empty hostfile for test
	_ = os.WriteFile(mockHostFile, []byte(""), 0644)

	var ConfigureFailTestCases = []struct {
		name        string
		cliParams   []string
		expectedOut []string
	}{

		{
			name: "configure fails when --server-certificate  parameter is missing",
			cliParams: []string{
				"--host", testutils.DefaultHost,
				"--ca-certificate", "/tmp/certificates/ca-cert.pem",
				"--server-key", "/tmp/certificates/server-key.pem",
			},
			expectedOut: []string{
				"[ERROR]:-one of the following flags is missing. Please specify --server-key, --server-certificate & --ca-certificate flags",
			},
		},
		{
			name: "configure fails when --server-key parameter is missing",
			cliParams: []string{
				"--host", testutils.DefaultHost,
				"--ca-certificate", "/tmp/certificates/ca-cert.pem",
				"--server-certificate", "/tmp/certificates/server-cert.pem",
			},
			expectedOut: []string{
				"[ERROR]:-one of the following flags is missing. Please specify --server-key, --server-certificate & --ca-certificate flags",
			},
		},
		{
			name: "configure fails when --ca-certificate parameter is missing",
			cliParams: []string{
				"--host", testutils.DefaultHost,
				"--server-certificate", "/tmp/certificates/server-cert.pem",
				"--server-key", "/tmp/certificates/server-key.pem",
			},
			expectedOut: []string{
				"[ERROR]:-one of the following flags is missing. Please specify --server-key, --server-certificate & --ca-certificate flags",
			},
		},
	}
	for _, tc := range ConfigureFailTestCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := testutils.RunGPServiceInit(false, tc.cliParams...)
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
			if strings.Contains(result.OutputMsg, "panic") {
				t.Errorf("\nUnexpected string: %#v \nFound in: %#v", "panic", result.OutputMsg)
			}
		})
	}
}
func TestConfigureInvalidCertificateFailure(t *testing.T) {
	// creating empty hostfile for test
	_ = os.WriteFile(mockHostFile, []byte(""), 0644)

	t.Run("configure service with --no-tls option", func(t *testing.T) {
		cliParams := []string{
			"--host", testutils.DefaultHost,
			"--server-certificate", "/invalid/path/to/server-cert.pem",
			"--server-key", "/invalid/path/to/server-key.pem",
			"--ca-certificate", "/invalid/path/to/ca-cert.pem",
		}
		expectedOut := []string{
			"Please make sure file exists before starting services",
		}
		result, err := testutils.RunGPServiceInit(false, cliParams...)
		if err != nil {
			t.Errorf("\nUnexpected error: %#v", err)
		}
		if result.ExitCode != 0 {
			t.Errorf("\nExpected: %v \nGot: %v", 0, result.ExitCode)
		}
		for _, item := range expectedOut {
			if !strings.Contains(result.OutputMsg, item) {
				t.Errorf("\nExpected string: %#v \nNot found in: %#v", item, result.OutputMsg)
			}
		}
	})
}
