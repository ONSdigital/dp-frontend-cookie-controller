// Code generated by MockGen. DO NOT EDIT.
// Source: handlers/handlers.go

// Package handlers is a generated GoMock package.
package handlers

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockClientError is a mock of ClientError interface
type MockClientError struct {
	ctrl     *gomock.Controller
	recorder *MockClientErrorMockRecorder
}

// MockClientErrorMockRecorder is the mock recorder for MockClientError
type MockClientErrorMockRecorder struct {
	mock *MockClientError
}

// NewMockClientError creates a new mock instance
func NewMockClientError(ctrl *gomock.Controller) *MockClientError {
	mock := &MockClientError{ctrl: ctrl}
	mock.recorder = &MockClientErrorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClientError) EXPECT() *MockClientErrorMockRecorder {
	return m.recorder
}

// Error mocks base method
func (m *MockClientError) Error() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Error")
	ret0, _ := ret[0].(string)
	return ret0
}

// Error indicates an expected call of Error
func (mr *MockClientErrorMockRecorder) Error() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockClientError)(nil).Error))
}

// Code mocks base method
func (m *MockClientError) Code() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Code")
	ret0, _ := ret[0].(int)
	return ret0
}

// Code indicates an expected call of Code
func (mr *MockClientErrorMockRecorder) Code() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Code", reflect.TypeOf((*MockClientError)(nil).Code))
}

// MockRenderClient is a mock of RenderClient interface
type MockRenderClient struct {
	ctrl     *gomock.Controller
	recorder *MockRenderClientMockRecorder
}

// MockRenderClientMockRecorder is the mock recorder for MockRenderClient
type MockRenderClientMockRecorder struct {
	mock *MockRenderClient
}

// NewMockRenderClient creates a new mock instance
func NewMockRenderClient(ctrl *gomock.Controller) *MockRenderClient {
	mock := &MockRenderClient{ctrl: ctrl}
	mock.recorder = &MockRenderClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRenderClient) EXPECT() *MockRenderClientMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockRenderClient) Do(arg0 string, arg1 []byte) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0, arg1)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do
func (mr *MockRenderClientMockRecorder) Do(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockRenderClient)(nil).Do), arg0, arg1)
}