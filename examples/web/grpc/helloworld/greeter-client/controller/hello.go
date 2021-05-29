package controller

import (
	"golang.org/x/net/context"
	"github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/protobuf"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/at"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
)

// controller
type helloController struct {
	// embedded at.RestController
	at.RestController
	at.RequestMapping `value:"/hello"`
	// declare HelloServiceClient
	helloServiceClient protobuf.HelloServiceClient
}

// Init inject helloServiceClient
func newHelloController(helloServiceClient protobuf.HelloServiceClient) *helloController {
	return &helloController{
		helloServiceClient: helloServiceClient,
	}
}

// GET /greeter/name/{name}
func (c *helloController) GetByName(_ struct{at.GetMapping `value:"/{name}"`}, name string) (response string) {

	// call grpc server method
	// pass context.Background() for the sake of simplicity
	result, err := c.helloServiceClient.SayHello(context.Background(), &protobuf.HelloRequest{Name: name})

	// got response
	if err == nil {
		response = result.Message
	}
	return
}

func init() {

	// must: register grpc client, the name greeter-client should configured in application.yml
	// see config/application-grpc.yml
	//
	// grpc:
	//   client:
	// 	   hello-world-service:   # client name
	//       host: localhost # server host
	//       port: 7575      # server port
	//
	grpc.Client("hello-world-service",
		protobuf.NewHelloServiceClient)

	// must: register Rest Controller
	app.Register(newHelloController)
}
