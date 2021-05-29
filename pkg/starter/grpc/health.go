package grpc

import (
	"context"
	pb "google.golang.org/grpc/health/grpc_health_v1"
	"github.com/hidevopsio/hiboot/pkg/at"
	"time"
)

// controller
type HealthCheckService struct {
	at.HealthCheckService
	// declare HelloServiceClient
	healthClient pb.HealthClient
}

// NewHealthCheckService is the constructor of healthCheckService
func NewHealthCheckService(healthClient pb.HealthClient) *HealthCheckService {
	return &HealthCheckService{
		healthClient: healthClient,
	}
}

// Status return health check display name grpc
func (c *HealthCheckService) Name() (name string) {
	return Profile
}

// Status return grpc health check status as bool
func (c *HealthCheckService) Status() (up bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.healthClient.Check(ctx, &pb.HealthCheckRequest{})
	if err == nil {
		up = resp.Status == pb.HealthCheckResponse_SERVING
	}
	return
}
