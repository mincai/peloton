// Code generated by MockGen. DO NOT EDIT.
// Source: code.uber.internal/infra/peloton/hostmgr/mesos (interfaces: MasterDetector,FrameworkInfoProvider)

package mocks

import (
	context "context"
	reflect "reflect"

	v1 "code.uber.internal/infra/peloton/.gen/mesos/v1"
	gomock "github.com/golang/mock/gomock"
)

// MockMasterDetector is a mock of MasterDetector interface
type MockMasterDetector struct {
	ctrl     *gomock.Controller
	recorder *MockMasterDetectorMockRecorder
}

// MockMasterDetectorMockRecorder is the mock recorder for MockMasterDetector
type MockMasterDetectorMockRecorder struct {
	mock *MockMasterDetector
}

// NewMockMasterDetector creates a new mock instance
func NewMockMasterDetector(ctrl *gomock.Controller) *MockMasterDetector {
	mock := &MockMasterDetector{ctrl: ctrl}
	mock.recorder = &MockMasterDetectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockMasterDetector) EXPECT() *MockMasterDetectorMockRecorder {
	return _m.recorder
}

// HostPort mocks base method
func (_m *MockMasterDetector) HostPort() string {
	ret := _m.ctrl.Call(_m, "HostPort")
	ret0, _ := ret[0].(string)
	return ret0
}

// HostPort indicates an expected call of HostPort
func (_mr *MockMasterDetectorMockRecorder) HostPort() *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "HostPort", reflect.TypeOf((*MockMasterDetector)(nil).HostPort))
}

// MockFrameworkInfoProvider is a mock of FrameworkInfoProvider interface
type MockFrameworkInfoProvider struct {
	ctrl     *gomock.Controller
	recorder *MockFrameworkInfoProviderMockRecorder
}

// MockFrameworkInfoProviderMockRecorder is the mock recorder for MockFrameworkInfoProvider
type MockFrameworkInfoProviderMockRecorder struct {
	mock *MockFrameworkInfoProvider
}

// NewMockFrameworkInfoProvider creates a new mock instance
func NewMockFrameworkInfoProvider(ctrl *gomock.Controller) *MockFrameworkInfoProvider {
	mock := &MockFrameworkInfoProvider{ctrl: ctrl}
	mock.recorder = &MockFrameworkInfoProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockFrameworkInfoProvider) EXPECT() *MockFrameworkInfoProviderMockRecorder {
	return _m.recorder
}

// GetFrameworkID mocks base method
func (_m *MockFrameworkInfoProvider) GetFrameworkID(_param0 context.Context) *v1.FrameworkID {
	ret := _m.ctrl.Call(_m, "GetFrameworkID", _param0)
	ret0, _ := ret[0].(*v1.FrameworkID)
	return ret0
}

// GetFrameworkID indicates an expected call of GetFrameworkID
func (_mr *MockFrameworkInfoProviderMockRecorder) GetFrameworkID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetFrameworkID", reflect.TypeOf((*MockFrameworkInfoProvider)(nil).GetFrameworkID), arg0)
}

// GetMesosStreamID mocks base method
func (_m *MockFrameworkInfoProvider) GetMesosStreamID(_param0 context.Context) string {
	ret := _m.ctrl.Call(_m, "GetMesosStreamID", _param0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetMesosStreamID indicates an expected call of GetMesosStreamID
func (_mr *MockFrameworkInfoProviderMockRecorder) GetMesosStreamID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCallWithMethodType(_mr.mock, "GetMesosStreamID", reflect.TypeOf((*MockFrameworkInfoProvider)(nil).GetMesosStreamID), arg0)
}
