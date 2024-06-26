package testutils

import (
	"context"

	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// MockHealthClient implements the healthpb.HealthClient interface
// which can be  used to mock the the gRPC health check RPCs
type MockHealthClient struct {
	status healthpb.HealthCheckResponse_ServingStatus
	err    error
}

func NewMockHealthClient(status healthpb.HealthCheckResponse_ServingStatus, err error) *MockHealthClient {
	return &MockHealthClient{
		status: status,
		err:    err,
	}
}

func (m *MockHealthClient) Check(ctx context.Context, in *healthpb.HealthCheckRequest, opts ...grpc.CallOption) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{
		Status: m.status,
	}, m.err
}

func (m *MockHealthClient) Watch(ctx context.Context, in *healthpb.HealthCheckRequest, opts ...grpc.CallOption) (healthpb.Health_WatchClient, error) {
	return nil, m.err
}
