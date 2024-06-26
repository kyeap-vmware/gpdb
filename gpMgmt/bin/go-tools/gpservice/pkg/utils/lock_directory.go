package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/greenplum-db/gp-common-go-libs/gplog"
	"github.com/greenplum-db/gpdb/gpservice/constants"
	"github.com/spf13/cobra"
)

type LockDir struct {
	dirPath string
}

func NewLockDir(cmd *cobra.Command) (*LockDir, error) {

	lockDir := &LockDir{
		dirPath: filepath.Join(constants.LockDirPath, strings.ReplaceAll(cmd.CommandPath(), " ", "_")),
	}

	if lockDir.IsProcessRunning() {
		return nil, fmt.Errorf("lockfile %s/pid indicates that an instance of %s is already running with PID %d. If this is not the case, remove the lockfile directory at %s",
			lockDir.dirPath, cmd.CommandPath(), os.Getpid(), lockDir.dirPath)
	}

	if err := lockDir.Create(); err != nil {
		return nil, fmt.Errorf("could not create lock directory: %w", err)
	}
	return lockDir, nil
}

func (ld *LockDir) IsProcessRunning() bool {

	pidFilePath := filepath.Join(ld.dirPath, constants.LockDirPidFile)

	raw, err := os.ReadFile(pidFilePath)
	if err != nil {
		gplog.Debug("Readfile failed err: %v", err)
		return false
	}

	pid, err := strconv.Atoi(string(raw))
	if err != nil {
		gplog.Debug("strconv failed err: %v", err)
		return false
	}

	//If pid is found then verify if the process is actually running
	if CheckPid(int(pid)) {
		return true
	}

	return false
}

func (ld *LockDir) Create() error {
	err := os.Mkdir(ld.dirPath, 0700)
	if os.IsExist(err) {
		gplog.Debug("Directory %s already present. Skipping.", ld.dirPath)
	}

	err = ld.CreateProcessIdFile()
	if err != nil {
		gplog.Debug("CreateProcessIdFile failed, error %v.", err)
		return err
	}

	return nil
}

func (ld *LockDir) Remove() error {
	return System.RemoveAll(ld.dirPath)
}

// This func creates the ProcessId file in the lock Directory
func (ld *LockDir) CreateProcessIdFile() error {

	//Get the pid of the running process
	pidFilePath := filepath.Join(ld.dirPath, constants.LockDirPidFile)

	file, err := System.OpenFile(pidFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		gplog.Debug("Failed to open file %s, error %v.", pidFilePath, err)
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(fmt.Sprint(os.Getpid())))
	if err != nil {
		gplog.Debug("Failed to write to file %s, error %v.", pidFilePath, err)
		return err
	}
	return nil

}
