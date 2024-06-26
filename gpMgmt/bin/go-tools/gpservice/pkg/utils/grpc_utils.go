package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/greenplum-db/gp-common-go-libs/gplog"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var newHealthClient = healthpb.NewHealthClient

// CheckGRPCServerHealth checks the health of a gRPC server by sending a health check request.
// The function retries the health check up to 5 times.
// If the server is running and responds with a SERVING status, the function returns nil.
// If the server is not running or responds with a non-SERVING status, the function returns an error.
func CheckGRPCServerHealth(conn *grpc.ClientConn) error {
	var err error
	var resp *healthpb.HealthCheckResponse

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := newHealthClient(conn)

	for i := 0; i < 5; i++ {
		resp, err = client.Check(ctx, &healthpb.HealthCheckRequest{}, grpc.WaitForReady(true))
		if err != nil {
			return fmt.Errorf("failed to get grpc server %s health: %w", conn.Target(), FormatGrpcError(err))
		}

		if resp.Status == healthpb.HealthCheckResponse_SERVING {
			gplog.Debug("grpc server %s running, status: %s", conn.Target(), resp.Status)
			return nil
		}

		err = fmt.Errorf("grpc server %s not running, status: %s", conn.Target(), resp.Status)
		gplog.Debug("%s. Retrying", err)

		System.Sleep(1 * time.Second)
	}

	return err
}

func SetNewHealthClient(mock healthpb.HealthClient) {
	newHealthClient = func(cc grpc.ClientConnInterface) healthpb.HealthClient {
		return mock
	}
}

func ResetNewHealthClient() {
	newHealthClient = healthpb.NewHealthClient
}
