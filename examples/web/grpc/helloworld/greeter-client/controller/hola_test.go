package controller

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/mock"
	"github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/protobuf"
	"github.com/hidevopsio/hiboot/pkg/app"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc/mockgrpc"
	"github.com/hidevopsio/hiboot/pkg/starter/logging"
)

func TestHolaClient(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHolaClient := mock.NewMockHolaServiceClient(ctrl)
	app.Register("protobuf.holaServiceClient", mockHolaClient)

	testApp := web.NewTestApp(t, newHolaController).SetProperty(logging.Level, logging.LevelDebug).Run(t)

	req := &protobuf.HolaRequest{Name: "Steve"}

	mockHolaClient.EXPECT().SayHola(
		gomock.Any(),
		&mockgrpc.RPCMsg{Message: req},
	).Return(&protobuf.HolaReply{Message: "Hola " + req.Name}, nil)

	testApp.Get("/hola/{name}").
		WithPath("name", req.Name).
		Expect().Status(http.StatusOK).
		Body().Contains(req.Name)
}
