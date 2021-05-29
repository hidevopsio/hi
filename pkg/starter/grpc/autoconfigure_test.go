// Copyright 2018 John Deng (hi.devops.io@gmail.com).
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//go:generate protoc -I ./helloworld --go_out=plugins=grpc:./helloworld ./helloworld/helloworld.proto

package grpc_test

import (
	"github.com/golang/mock/gomock"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc/helloworld"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc/mockgrpc"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"testing"
)

// gRpc server
// server is used to implement helloworld.GreeterServer.
type greeterServerService struct{
	helloworld.UnimplementedGreeterServer
}

// newGreeterServerService is the constructor of greeterServerService
func newGreeterServerService() helloworld.GreeterServer {
	return &greeterServerService{}
}

// SayHello implements helloworld.GreeterServer
func (s *greeterServerService) SayHello(ctx context.Context, req *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + req.Name}, nil
}

// gRpc client
type greeterClientService struct {
	greeterClient helloworld.GreeterClient
}

// newGreeterClientService inject greeterClient
func newGreeterClientService(greeterClient helloworld.GreeterClient) *greeterClientService {
	return &greeterClientService{
		greeterClient: greeterClient,
	}
}

// SayHello gRpc client service that call remote server service
func (s *greeterClientService) SayHello(name string) (*helloworld.HelloReply, error) {
	response, err := s.greeterClient.SayHello(context.Background(), &helloworld.HelloRequest{Name: name})
	return response, err
}

func TestGrpcServerAndClient(t *testing.T) {

	app.Register(newGreeterClientService)

	grpc.Server(helloworld.RegisterGreeterServer, newGreeterServerService)
	grpc.Client("greeter-service", helloworld.NewGreeterClient)

	testApp := web.RunTestApplication(t)
	assert.NotEqual(t, nil, testApp)

	applicationContext := testApp.(app.ApplicationContext)

	t.Run("should find gRpc client and call its services", func(t *testing.T) {
		cliSvc := applicationContext.GetInstance(greeterClientService{})
		assert.NotEqual(t, nil, cliSvc)
		if cliSvc != nil {
			greeterCliSvc := cliSvc.(*greeterClientService)
			assert.NotEqual(t, nil, greeterCliSvc.greeterClient)
		}
	})

	t.Run("should connect to gRpc service at runtime", func(t *testing.T) {
		cc := applicationContext.GetInstance(new(grpc.ClientConnector)).(grpc.ClientConnector)
		f := applicationContext.GetInstance(new(factory.InstantiateFactory)).(factory.InstantiateFactory)
		assert.NotEqual(t, nil, cc)
		prop := new(grpc.ClientProperties)
		f.InjectDefaultValue(prop)
		grpcCli, err := cc.ConnectWithName("", helloworld.NewGreeterClient, prop)
		assert.Equal(t, nil, err)
		assert.NotEqual(t, nil, grpcCli)
	})

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockGreeterClient := mockgrpc.NewMockGreeterClient(ctrl)

	t.Run("should get message from mock gRpc client indirectly", func(t *testing.T) {
		cliSvc := newGreeterClientService(mockGreeterClient)
		assert.NotEqual(t, nil, cliSvc)
		if cliSvc != nil {
			req := &helloworld.HelloRequest{Name: "Steve"}
			mockGreeterClient.EXPECT().SayHello(
				gomock.Any(),
				&mockgrpc.RPCMsg{Message: req},
			).Return(&helloworld.HelloReply{Message: "Hello " + req.Name}, nil)
			resp, err := cliSvc.SayHello("Steve")
			assert.Equal(t, nil, err)
			assert.Equal(t, "Hello "+req.Name, resp.Message)
		}
	})

	mockHealthClient := mockgrpc.NewMockHealthClient(ctrl)
	t.Run("should get health status from client", func(t *testing.T) {
		mockHealthClient.EXPECT().Check(
			gomock.Any(),
			gomock.Any(),
			gomock.Any(),
		).Return(&grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil)

		healthCheckService := grpc.NewHealthCheckService(mockHealthClient)
		assert.Equal(t, grpc.Profile, healthCheckService.Name())
		assert.Equal(t, true, healthCheckService.Status())
	})
}
