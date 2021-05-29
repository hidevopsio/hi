package grpc

import (
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"github.com/hidevopsio/hiboot/pkg/factory"
	"github.com/hidevopsio/hiboot/pkg/log"
	"github.com/hidevopsio/hiboot/pkg/starter/jaeger"
	"github.com/hidevopsio/hiboot/pkg/utils/reflector"
)

// ClientConnector interface is response for creating grpc client connection
type ClientConnector interface {
	// Connect connect the gRPC client
	ConnectWithName(name string, cb interface{}, prop *ClientProperties) (gRPCCli interface{}, err error)
	Connect(address string) (conn *grpc.ClientConn, err error)
}

type clientConnector struct {
	instantiateFactory factory.InstantiateFactory
	tracer jaeger.Tracer
}

func newClientConnector(instantiateFactory factory.InstantiateFactory, tracer jaeger.Tracer) ClientConnector {
	cc := &clientConnector{
		instantiateFactory: instantiateFactory,
		tracer: tracer,
	}
	return cc
}


func Connect(address string, tracer jaeger.Tracer) (conn *grpc.ClientConn, err error) {
	if tracer == nil {
		conn, err = grpc.Dial(address,
			grpc.WithInsecure(),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             time.Second,
				PermitWithoutStream: true,
			}),
		)
	} else {
		conn, err = grpc.Dial(address,
			grpc.WithInsecure(),
			grpc.WithStreamInterceptor(grpc_middleware.ChainStreamClient(
				grpc_opentracing.StreamClientInterceptor(grpc_opentracing.WithTracer(tracer)),
			)),
			grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
				grpc_opentracing.UnaryClientInterceptor(grpc_opentracing.WithTracer(tracer)),
			)),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                10 * time.Second,
				Timeout:             time.Second,
				PermitWithoutStream: true,
			}),
		)
	}
	return conn, err
}

// ConnectWithName connect to grpc server from client with service name
// name: client name
// clientConstructor: client constructor
// properties: properties for configuring
func (c *clientConnector) ConnectWithName(name string, clientConstructor interface{}, properties *ClientProperties) (gRpcCli interface{}, err error) {
	host := properties.Host
	if host == "" {
		host = name
	}
	address := host + ":" + properties.Port
	conn := c.instantiateFactory.GetInstance(name)
	if conn == nil {
		// connect to grpc server
		conn, err = Connect(address, c.tracer)
		c.instantiateFactory.SetInstance(name, conn)
		if err == nil {
			log.Infof("gRPC client connected to: %v", address)
		}
	}
	if err == nil && clientConstructor != nil {
		// get return type for register instance name
		gRpcCli, err = reflector.CallFunc(clientConstructor, conn)
	}
	return
}


// Connect connect to client connection
func (c *clientConnector) Connect(address string) (conn *grpc.ClientConn, err error) {

	conn, err = Connect(address, c.tracer)
	if err == nil {
		log.Infof("gRPC client connected to: %v", address)
	}

	return
}

