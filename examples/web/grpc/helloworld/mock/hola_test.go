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

package mock

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	protobuf2 "github.com/hidevopsio/hiboot/examples/web/grpc/helloworld/protobuf"
	mockproto "github.com/hidevopsio/hiboot/pkg/starter/grpc/mockgrpc"
	"testing"
	"time"
)

func TestMockHola(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockHolaClient := NewMockHolaServiceClient(ctrl)
	t.Run("should get message from mock gRpc client directly", func(t *testing.T) {
		req := &protobuf2.HolaRequest{Name: "unit_test"}
		opt := &grpc.HeaderCallOption{}
		mockHolaClient.EXPECT().SayHola(
			gomock.Any(),
			&mockproto.RPCMsg{Message: req},
			opt,
		).Return(&protobuf2.HolaReply{Message: "Mocked Interface"}, nil)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := mockHolaClient.SayHola(ctx, req, opt)
		assert.Equal(t, nil, err)
		assert.Equal(t, "Mocked Interface", r.Message)
	})
}
