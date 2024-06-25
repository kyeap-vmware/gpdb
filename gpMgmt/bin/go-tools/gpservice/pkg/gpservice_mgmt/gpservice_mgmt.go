package gpservice_mgmt

import (
	"github.com/greenplum-db/gpdb/gpservice/internal/cli"
	"github.com/greenplum-db/gpdb/gpservice/pkg/gpservice_config"
)

func StartServices(conf *gpservice_config.Config) error {
	return cli.StartServices(conf)
}

func StopServices(conf *gpservice_config.Config) error {
	return cli.StopServices(conf)
}

func DeleteServices(conf *gpservice_config.Config, confFile string) error {
	return cli.DeleteServices(conf, confFile)
}

func InitialiseGpService(configFilepath string, hubPort, agentPort int, hostnames []string, hubLogDir, serviceName,
	gpHome, caCertPath, serverCertPath, serverKeyPath string, NoTls, defaultConfig bool) error {

	return cli.InitGpService(configFilepath, hubPort, agentPort, hostnames, hubLogDir, serviceName,
		gpHome, caCertPath, serverCertPath, serverKeyPath, NoTls, defaultConfig)
}
