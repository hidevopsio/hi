// Code generated by MockGen. DO NOT EDIT.
// Source: google.golang.org/grpc/health/grpc_health_v1 (interfaces: HealthClient)

// Package mockgrpc is a generated GoMock package.
package mockgrpc

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	grpc "google.golang.org/grpc"
	grpc_health_v1 "google.golang.org/grpc/health/grpc_health_v1"
	reflect "reflect"
)

// MockHealthClient is a mock of HealthClient interface
type MockHealthClient struct {
	ctrl     *gomock.Controller
	recorder *MockHealthClientMockRecorder
}

// MockHealthClientMockRecorder is the mock recorder for MockHealthClient
type MockHealthClientMockRecorder struct {
	mock *MockHealthClient
}

// NewMockHealthClient creates a new mock instance
func NewMockHealthClient(ctrl *gomock.Controller) *MockHealthClient {
	mock := &MockHealthClient{ctrl: ctrl}
	mock.recorder = &MockHealthClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockHealthClient) EXPECT() *MockHealthClientMockRecorder {
	return m.recorder
}

// Check mocks base method
func (m *MockHealthClient) Check(arg0 context.Context, arg1 *grpc_health_v1.HealthCheckRequest, arg2 ...grpc.CallOption) (*grpc_health_v1.HealthCheckResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Check", varargs...)
	ret0, _ := ret[0].(*grpc_health_v1.HealthCheckResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Check indicates an expected call of Check
func (mr *MockHealthClientMockRecorder) Check(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockHealthClient)(nil).Check), varargs...)
}

// Watch mocks base method
func (m *MockHealthClient) Watch(arg0 context.Context, arg1 *grpc_health_v1.HealthCheckRequest, arg2 ...grpc.CallOption) (grpc_health_v1.Health_WatchClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Watch", varargs...)
	ret0, _ := ret[0].(grpc_health_v1.Health_WatchClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch
func (mr *MockHealthClientMockRecorder) Watch(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockHealthClient)(nil).Watch), varargs...)
}
