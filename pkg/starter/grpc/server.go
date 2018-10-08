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

package grpc

import (
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

// ServerFactory build grpc servers
type ServerFactory interface {
}

type serverFactory struct {
}

func newServerFactory(properties properties, grpcServer *grpc.Server) ServerFactory {
	sf := &serverFactory{}
	sf.buildServers(properties, grpcServer)

	return sf
}

func (f *serverFactory) buildServers(properties properties, grpcServer *grpc.Server) {
	// just return if grpc server is not enabled
	if !properties.Server.Enabled || grpcServer == nil {
		return
	}

	address := properties.Server.Host + ":" + properties.Server.Port
	lis, err := net.Listen(properties.Server.Network, address)
	if err != nil {
		log.Fatalf("failed to listen: %v, %v", address, err)
	}

	// register server
	// Register reflection service on gRPC server.
	chn := make(chan bool)
	go func() {
		for _, srv := range grpcServers {
			reflector.CallFunc(srv.cb, grpcServer, srv.svc)
			svcName, err := reflector.GetName(srv.svc)
			if err == nil {
				log.Infof("Registered %v on gRPC server", svcName)
			}
		}
		reflection.Register(grpcServer)
		chn <- true
		if err := grpcServer.Serve(lis); err != nil {
			log.Print("failed to serve: %v", err)
		}
		log.Print("gRPC server exit\n")
	}()
	<-chn

	log.Infof("gRPC server listening on: localhost%v", address)
}
