package utils_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/greenplum-db/gp-common-go-libs/testhelper"
	"github.com/greenplum-db/gpdb/gpservice/constants"
	"github.com/greenplum-db/gpdb/gpservice/internal/cli"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
	"github.com/greenplum-db/gpdb/gpservice/testutils"
)

func TestCreateDir(t *testing.T) {
	_, _, logfile := testhelper.SetupTestLogger()
	rootCmd := cli.RootCommand()
	ld, _ := utils.NewLockDir(rootCmd)

	t.Run("should not create directory if its already present", func(t *testing.T) {

		dirPath := filepath.Join(constants.LockDirPath, strings.ReplaceAll(rootCmd.CommandPath(), " ", "_"))
		err := os.Mkdir(dirPath, 0700)
		if err != nil && !os.IsExist(err) {
			t.Fatalf("Unexpected error: %v", err)
		}

		defer os.RemoveAll(dirPath)

		err = ld.Create()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		expectedErr := fmt.Sprintf("Directory %s already present. Skipping.", dirPath)
		testutils.AssertLogMessage(t, logfile, expectedErr)

	})

	t.Run("should create directory and file", func(t *testing.T) {

		err := ld.Create()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		dirPath := filepath.Join("/tmp", strings.ReplaceAll(rootCmd.CommandPath(), " ", "_"), "pid")
		_, err = utils.System.Stat(dirPath)
		if err != nil && !os.IsNotExist(err) {
			t.Fatalf("Failed to create directory: %v", err)
		}
		os.RemoveAll(dirPath)
	})
}

func TestCreateProcessIdFile(t *testing.T) {

	_, _, logfile := testhelper.SetupTestLogger()
	rootCmd := cli.RootCommand()
	lockDir, _ := utils.NewLockDir(rootCmd)
	t.Run("should successfully write current process id to file", func(t *testing.T) {

		err := lockDir.CreateProcessIdFile()
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		dirPath := filepath.Join(constants.LockDirPath, strings.ReplaceAll(rootCmd.CommandPath(), " ", "_"))
		processIdFile := filepath.Join(dirPath, constants.LockDirPidFile)
		raw, err := os.ReadFile(processIdFile)
		if err != nil {
			t.Fatalf("Failed to read process-id file: %v", err)
		}

		pid, err := strconv.Atoi(string(raw))
		if err != nil {
			t.Fatalf("String to int coversion failed: %v", err)
		}

		ExpectedPid := os.Getpid()
		if pid != ExpectedPid {
			t.Fatalf("got pid %d, want pid %d", pid, ExpectedPid)
		}
	})

	t.Run("Check for open failure", func(t *testing.T) {

		expectedErr := errors.New("error")
		utils.System.OpenFile = func(name string, flag int, perm fs.FileMode) (*os.File, error) {
			return nil, expectedErr
		}
		defer utils.ResetSystemFunctions()

		err := lockDir.CreateProcessIdFile()
		if !errors.Is(err, expectedErr) {
			t.Fatalf("got %#v, want %#v", err, expectedErr)
		}

		dirPath := filepath.Join(constants.LockDirPath, strings.ReplaceAll(rootCmd.CommandPath(), " ", "_"))
		processIdFile := filepath.Join(dirPath, constants.LockDirPidFile)
		expErr := fmt.Sprintf("Failed to open file %s", processIdFile)
		testutils.AssertLogMessage(t, logfile, expErr)
	})

}
