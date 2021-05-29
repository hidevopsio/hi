package mockgrpc

import (
	"github.com/hidevopsio/hiboot/pkg/starter/grpc/helloworld"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRpcMsg(t *testing.T) {
	rpcMsg := new(RPCMsg)
	req := &helloworld.HelloRequest{Name: "unit_test"}

	ok := rpcMsg.Matches(nil)
	assert.Equal(t, false, ok)

	rpcMsg.Message = req
	assert.Contains(t, rpcMsg.Message.String(), "unit_test", )

	ok = rpcMsg.Matches(req)
	assert.Equal(t, true, ok)

}
