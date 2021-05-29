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
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
)

// ServerFactory build grpc servers
type ServerFactory interface {
}

type serverFactory struct {
}

func newServerFactory(instantiateFactory factory.InstantiateFactory, properties *properties, grpcServer *grpc.Server) ServerFactory {
	sf := &serverFactory{}

	// just return if grpc server is not enabled
	if properties.Server.Enabled && grpcServer != nil {
		address := properties.Server.Host + ":" + properties.Server.Port
		lis, err := net.Listen(properties.Server.Network, address)
		if err == nil {
			// register server
			// Register reflection service on gRPC server.
			chn := make(chan bool)
			go func() {
				for _, srv := range grpcServers {
					svc := instantiateFactory.GetInstance(srv.name)
					reflector.CallFunc(srv.cb, grpcServer, svc)
					if err == nil {
						log.Infof("Registered %v on gRPC server", srv.name)
					}
				}
				reflection.Register(grpcServer)
				chn <- true
				grpcServer.Serve(lis)
			}()
			<-chn

			log.Infof("gRPC server listening on: localhost%v", address)
		}

	}

	return sf
}
