package controller

import (
	"github.com/golang/mock/gomock"
	"github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/mock"
	"github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/protobuf"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc/mockgrpc"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
	"net/http"
	"testing"
)

func TestHelloClient(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHelloClient := mock.NewMockHelloServiceClient(ctrl)
	app.Register("protobuf.helloServiceClient", mockHelloClient)

	testApp := web.NewTestApp(t).SetProperty(logging.Level, logging.LevelDebug).Run(t)

	req := &protobuf.HelloRequest{Name: "Steve"}

	mockHelloClient.EXPECT().SayHello(
		gomock.Any(),
		&mockgrpc.RPCMsg{Message: req},
	).Return(&protobuf.HelloReply{Message: "Hello " + req.Name}, nil)

	testApp.Get("/hello/{name}").
		WithPath("name", req.Name).
		Expect().Status(http.StatusOK).
		Body().Contains(req.Name)
}
