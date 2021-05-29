package main

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	protobuf2 "github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/protobuf"
)

var mu sync.Mutex

func TestRunMain(t *testing.T) {
	mu.Lock()
	go main()
	mu.Unlock()
}

func TestHelloServer(t *testing.T) {

	serviceServer := newHelloServiceServer()

	name := "Steve"
	expected := "Hello " + name
	req := &protobuf2.HelloRequest{Name: name}
	res, err := serviceServer.SayHello(context.Background(), req)

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, res.Message)
}

func TestHolaServer(t *testing.T) {

	serviceServer := newHolaServiceServer()

	name := "Steve"
	expected := "Hola " + name
	req := &protobuf2.HolaRequest{Name: name}
	res, err := serviceServer.SayHola(context.Background(), req)

	assert.Equal(t, nil, err)
	assert.Equal(t, expected, res.Message)
}
