package grpc

import (
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/utils/mapstruct"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
)

// ClientFactory build grpc clients
type ClientFactory interface {
}

type clientFactory struct {
}

func newClientFactory(instantiateFactory factory.InstantiateFactory, properties *properties, cc ClientConnector) ClientFactory {
	cf := &clientFactory{}

	clientProps := properties.Client
	var gRPCCli interface{}
	for _, cli := range grpcClients {
		prop := new(ClientProperties)
		err := mapstruct.Decode(prop, clientProps[cli.name])
		if err == nil {
			gRPCCli, err = cc.ConnectWithName(cli.name, cli.cb, prop)
			if err == nil {
				clientInstanceName := reflector.GetLowerCamelFullName(gRPCCli)
				// register grpc client
				instantiateFactory.SetInstance(gRPCCli)

				log.Infof("Registered gRPC client %v", clientInstanceName)
			}
		}
	}

	return cf
}
