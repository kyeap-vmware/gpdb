package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/greenplum-db/gp-common-go-libs/gplog"
	"github.com/greenplum-db/gpdb/gpservice/constants"
	. "github.com/greenplum-db/gpdb/gpservice/internal/platform"
	config "github.com/greenplum-db/gpdb/gpservice/pkg/gpservice_config"
	"github.com/greenplum-db/gpdb/gpservice/pkg/greenplum"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
)

var (
	platform       = GetPlatform()
	agentPort      int
	caCertPath     string
	gpHome         string
	hubLogDir      string
	hubPort        int
	hostnames      []string
	hostfilePath   string
	serverCertPath string
	serverKeyPath  string
	serviceName    string
	noTlsFlag      bool

	GetUlimitSsh = GetUlimitSshFn
)

func InitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize gpservice as a systemd service",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if !cmd.Flags().Lookup("no-tls").Changed && !cmd.Flags().Lookup("ca-certificate").Changed &&
				!cmd.Flags().Lookup("server-certificate").Changed && !cmd.Flags().Lookup("server-key").Changed {
				return fmt.Errorf("no security options specified. Use either --no-tls for insecure communication or provide " +
					"the certificates using --ca-certificate, --server-certificate and --server-key for secure communication")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			err := RunConfigure(cmd)
			if err != nil {
				utils.LogErrorAndExit(err, 1)
			}
		},
	}

	v := viper.New()
	v.AutomaticEnv()

	initCmd.Flags().IntVar(&agentPort, "agent-port", constants.DefaultAgentPort, `Port on which the agents should listen`)
	initCmd.Flags().StringVar(&gpHome, "gphome", "/usr/local/greenplum-db", `Path to GPDB installation`)
	viper.BindPFlag("gphome", initCmd.Flags().Lookup("gphome")) // nolint
	initCmd.Flags().IntVar(&hubPort, "hub-port", constants.DefaultHubPort, `Port on which the hub should listen`)
	initCmd.Flags().StringVar(&hubLogDir, "log-dir", greenplum.GetDefaultHubLogDir(), `Path to gpservice log directory`)
	initCmd.Flags().StringVar(&serviceName, "service-name", constants.DefaultServiceName, `Name for the generated systemd service file`)
	initCmd.Flags().StringVar(&caCertPath, "ca-certificate", "", `Path to SSL/TLS CA certificate`)
	initCmd.Flags().StringVar(&serverCertPath, "server-certificate", "", `Path to hub SSL/TLS server certificate`)
	initCmd.Flags().StringVar(&serverKeyPath, "server-key", "", `Path to hub SSL/TLS server private key`)
	initCmd.Flags().StringArrayVar(&hostnames, "host", []string{}, `Segment hostname`)
	initCmd.Flags().StringVar(&hostfilePath, "hostfile", "", `Path to file containing a list of segment hostnames`)
	initCmd.Flags().BoolVar(&noTlsFlag, "no-tls", false, "Set this flag if need to run hub and agents without transport layer security (TLS)")

	initCmd.MarkFlagsMutuallyExclusive("host", "hostfile")
	initCmd.MarkFlagsOneRequired("host", "hostfile")
	initCmd.MarkFlagsMutuallyExclusive("no-tls", "ca-certificate")
	initCmd.MarkFlagsMutuallyExclusive("no-tls", "server-certificate")
	initCmd.MarkFlagsMutuallyExclusive("no-tls", "server-key")
	initCmd.MarkFlagsRequiredTogether("ca-certificate", "server-certificate", "server-key")

	gpHome = v.GetString("gphome")

	return initCmd
}

func InitGpService(configFilepath string, hubPort, agentPort int, hostnames []string, hubLogDir, serviceName,
	gpHome, caCertPath, serverCertPath, serverKeyPath string, noTlsFlag, defaultConfig bool) error {

	credentials := &utils.GpCredentials{}

	if !noTlsFlag {
		credentials = &utils.GpCredentials{
			CACertPath:     caCertPath,
			ServerCertPath: serverCertPath,
			ServerKeyPath:  serverKeyPath,
			TlsEnabled:     true,
		}
		if _, err := os.Stat(caCertPath); errors.Is(err, os.ErrNotExist) {
			gplog.Warn("ca-certificate file %s does not exists. Please make sure file exists before starting services", caCertPath)
		}
		if _, err := os.Stat(serverCertPath); errors.Is(err, os.ErrNotExist) {
			gplog.Warn("server-certificate file %s does not exists. Please make sure file exists before starting services", serverCertPath)
		}
		if _, err := os.Stat(serverKeyPath); errors.Is(err, os.ErrNotExist) {
			gplog.Warn("server-key file %s does not exists. Please make sure file exists before starting services", serverKeyPath)
		}
	} else {
		credentials.TlsEnabled = false
	}
	err := config.Create(configFilepath, hubPort, agentPort, hostnames, hubLogDir, serviceName, gpHome, credentials, false)
	if err != nil {
		return err
	}

	err = platform.CreateServiceDir(hostnames, gpHome)
	if err != nil {
		return err
	}

	err = platform.CreateAndInstallHubServiceFile(gpHome, serviceName, configFilepath)
	if err != nil {
		return err
	}

	err = platform.CreateAndInstallAgentServiceFile(hostnames, gpHome, serviceName, configFilepath)
	if err != nil {
		return err
	}

	err = platform.EnableUserLingering(hostnames, gpHome)
	if err != nil {
		return err
	}

	CheckOpenFilesLimitOnHosts(hostnames)

	return err
}

func RunConfigure(cmd *cobra.Command) error {
	if gpHome == "" {
		return fmt.Errorf("no value found for gphome, please provide a valid gphome using --gphome or set the GPHOME env variable")
	}

	if _, err := os.Stat(gpHome); os.IsNotExist(err) {
		return fmt.Errorf("gphome path %s does not exist", gpHome)
	}

	// Regenerate default flag values if a custom GPHOME or username is passed
	if !cmd.Flags().Lookup("config-file").Changed {
		configFilepath = filepath.Join(gpHome, constants.ConfigFileName)
	}

	if agentPort == hubPort {
		return errors.New("hub port and agent port must be different")
	}

	// Convert file/directory paths to absolute path before writing to service configuration file
	err := resolveAbsolutePaths()
	if err != nil {
		return err
	}

	hostnames, err := getHostnames()
	if err != nil {
		return err
	}


	err = InitGpService(configFilepath, hubPort, agentPort, hostnames, hubLogDir, serviceName,
		gpHome, caCertPath, serverCertPath, serverKeyPath, noTlsFlag, false)

	return nil
}

/*
CheckOpenFilesLimitOnHosts checks for open files limit by calling ulimit command
Executes gpssh command to get the ulimit from remote hosts using go routine
Prints a warning if ulimit is lower.
This function depends on gpssh. Use only in the configure command.
*/
func CheckOpenFilesLimitOnHosts(hostnames []string) {
	// check Ulimit on local host
	ulimit, err := utils.ExecuteAndGetUlimit()
	if err != nil {
		gplog.Warn(err.Error())
	} else if ulimit < constants.OsOpenFiles {
		gplog.Warn("Open files limit for coordinator host. Value set to %d, expected:%d. For proper functioning make sure"+
			" limit is set properly for system and services before starting gp services.",
			ulimit, constants.OsOpenFiles)
	}
	var wg sync.WaitGroup
	//Check ulimit on other hosts
	channel := make(chan Response)
	for _, host := range hostnames {
		wg.Add(1)
		go GetUlimitSsh(host, channel, &wg)
	}
	go func() {
		wg.Wait()
		close(channel)
	}()
	for hostlimits := range channel {
		if hostlimits.Ulimit < constants.OsOpenFiles {
			gplog.Warn("Open files limit for host: %s is set to %d, expected:%d. For proper functioning make sure"+
				" limit is set properly for system and services before starting gp services.",
				hostlimits.Hostname, hostlimits.Ulimit, constants.OsOpenFiles)
		}
	}
}
func GetUlimitSshFn(hostname string, channel chan Response, wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := utils.System.ExecCommand(filepath.Join(gpHome, "bin", constants.GpSSH), "-h", hostname, "-e", "ulimit -n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		gplog.Warn("error executing command to fetch open files limit on host:%s, %v", hostname, err)
		return
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) < 2 {
		gplog.Warn("unexpected output when fetching open files limit on host:%s, gpssh output:%s", hostname, lines)
		return
	}
	values := strings.Split(lines[1], " ")
	if len(values) < 2 {
		gplog.Warn("unexpected output when parsing open files limit output for host:%s, gpssh output:%s", hostname, lines)
		return
	}
	ulimit, err := strconv.Atoi(values[1])
	if err != nil {
		gplog.Warn("unexpected output when converting open files limit value for host:%s, value:%s", hostname, values[1])
		return
	}
	channel <- Response{Hostname: hostname, Ulimit: ulimit}
}

type Response struct {
	Hostname string
	Ulimit   int
}

func resolveAbsolutePaths() error {
	paths := []*string{&caCertPath, &serverCertPath, &serverKeyPath, &hubLogDir, &gpHome}
	for _, path := range paths {
		p, err := filepath.Abs(*path)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute path for %s: %w", *path, err)
		}
		*path = p
	}

	return nil
}

func getHostnames() ([]string, error) {
	var result []string

	if len(hostnames) == 0 {
		contents, err := utils.System.ReadFile(hostfilePath)
		if err != nil {
			return []string{}, fmt.Errorf("failed to read hostfile %s: %w", hostfilePath, err)
		}

		hostnames = strings.Fields(string(contents))
	}

	for _, host := range hostnames {
		if host != "" {
			result = append(result, host)
		}
	}

	if len(result) < 1 {
		return []string{}, fmt.Errorf("no host name found, please provide a valid input host name using either --host or --hostfile")
	}

	return result, nil
}
