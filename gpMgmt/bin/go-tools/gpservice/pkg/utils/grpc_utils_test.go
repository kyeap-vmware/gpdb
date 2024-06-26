package utils_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/greenplum-db/gp-common-go-libs/testhelper"
	"github.com/greenplum-db/gpdb/gpservice/pkg/utils"
	"github.com/greenplum-db/gpdb/gpservice/testutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func TestCheckGRPCServerHealth(t *testing.T) {
	conn, err := grpc.NewClient("cdw:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	utils.System.Sleep = func(d time.Duration) {}
	defer utils.ResetSystemFunctions()

	t.Run("does not error out when the server is running", func(t *testing.T) {
		_, _, logfile := testhelper.SetupTestLogger()

		status := grpc_health_v1.HealthCheckResponse_SERVING
		utils.SetNewHealthClient(testutils.NewMockHealthClient(status, nil))
		defer utils.ResetNewHealthClient()

		err := utils.CheckGRPCServerHealth(conn)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		testutils.AssertLogMessage(t, logfile, fmt.Sprintf(`\[DEBUG\]:-grpc server %s running, status: %s`, conn.Target(), status))
	})

	t.Run("when the check RPC fails to execute", func(t *testing.T) {
		expectedErr := errors.New("error")
		utils.SetNewHealthClient(testutils.NewMockHealthClient(grpc_health_v1.HealthCheckResponse_SERVING, expectedErr))
		defer utils.ResetNewHealthClient()

		err := utils.CheckGRPCServerHealth(conn)
		if !errors.Is(err, expectedErr) {
			t.Fatalf("got %#v, want %#v", err, expectedErr)
		}

		expectedErrPrefix := fmt.Sprintf("failed to get grpc server %s health:", conn.Target())
		if !strings.HasPrefix(err.Error(), expectedErrPrefix) {
			t.Fatalf("got %v, want prefix %s", err, expectedErrPrefix)
		}
	})

	t.Run("retries when the server is not currently serving", func(t *testing.T) {
		_, _, logfile := testhelper.SetupTestLogger()

		status := grpc_health_v1.HealthCheckResponse_NOT_SERVING
		utils.SetNewHealthClient(testutils.NewMockHealthClient(status, nil))
		defer utils.ResetNewHealthClient()

		err := utils.CheckGRPCServerHealth(conn)
		expectedErrStr := fmt.Sprintf("grpc server %s not running, status: %s", conn.Target(), status)
		if err.Error() != expectedErrStr {
			t.Fatalf("got %v, want %s", err, expectedErrStr)
		}

		testutils.AssertLogMessageCount(t, logfile, fmt.Sprintf(`[DEBUG]:-grpc server %s not running, status: %s`, conn.Target(), status), 5)
	})
}
