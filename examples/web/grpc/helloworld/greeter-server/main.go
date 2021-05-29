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

// if protoc report command not found error, should install proto and protc-gen-go
// find protoc install instruction on http://google.github.io/proto-lens/installing-protoc.html
// go get -u -v github.com/golang/protobuf/{proto,protoc-gen-go}
///go:generate protoc --proto_path=../protobuf --go_out=plugins=grpc:../protobuf --go_opt=paths=source_relative ../protobuf/helloworld.proto
//go:generate protoc -I ../protobuf --go_out=plugins=grpc:../protobuf ../protobuf/helloworld.proto

package main

import (
	"github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/protobuf"
	"github.com/hidevopsio/hiboot/pkg/app/web"
	_ "github.com/hidevopsio/hiboot/pkg/starter/actuator"
	"github.com/hidevopsio/hiboot/pkg/starter/grpc"
	"golang.org/x/net/context"
)

// server is used to implement protobuf.GreeterServer.
type helloServiceServerImpl struct {
}

func newHelloServiceServer() protobuf.HelloServiceServer {
	return &helloServiceServerImpl{}
}

// SayHello implements helloworld.GreeterServer
func (s *helloServiceServerImpl) SayHello(ctx context.Context, request *protobuf.HelloRequest) (*protobuf.HelloReply, error) {
	// response to client
	return &protobuf.HelloReply{Message: "Hello " + request.Name}, nil
}

// server is used to implement protobuf.GreeterServer.
type holaServiceServerImpl struct {
}

func newHolaServiceServer() protobuf.HolaServiceServer {
	return &holaServiceServerImpl{}
}

// SayHello implements helloworld.GreeterServer
func (s *holaServiceServerImpl) SayHola(ctx context.Context, request *protobuf.HolaRequest) (*protobuf.HolaReply, error) {
	// response to client
	return &protobuf.HolaReply{Message: "Hola " + request.Name}, nil
}

func init() {
	// must: register grpc server
	// please note that holaServiceServerImpl must implement protobuf.HelloServiceServer, or it won't be registered.
	grpc.Server(protobuf.RegisterHelloServiceServer, newHelloServiceServer)
	grpc.Server(protobuf.RegisterHolaServiceServer, newHolaServiceServer)
}

func main() {
	// create new web application and run it
	web.NewApplication().Run()
}
