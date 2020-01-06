// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jamieabc/bitmarkd-broadcast-monitor/nodes/node (interfaces: Remote)

// Package mock_node is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	communication "github.com/jamieabc/bitmarkd-broadcast-monitor/communication"
	network "github.com/jamieabc/bitmarkd-broadcast-monitor/network"
	reflect "reflect"
)

// MockRemote is a mock of Remote interface
type MockRemote struct {
	ctrl     *gomock.Controller
	recorder *MockRemoteMockRecorder
}

// MockRemoteMockRecorder is the mock recorder for MockRemote
type MockRemoteMockRecorder struct {
	mock *MockRemote
}

// NewMockRemote creates a new mock instance
func NewMockRemote(ctrl *gomock.Controller) *MockRemote {
	mock := &MockRemote{ctrl: ctrl}
	mock.recorder = &MockRemoteMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRemote) EXPECT() *MockRemoteMockRecorder {
	return m.recorder
}

// BlockHeader mocks base method
func (m *MockRemote) BlockHeader(arg0 uint64) (*communication.BlockHeaderResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockHeader", arg0)
	ret0, _ := ret[0].(*communication.BlockHeaderResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// BlockHeader indicates an expected call of BlockHeader
func (mr *MockRemoteMockRecorder) BlockHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockHeader", reflect.TypeOf((*MockRemote)(nil).BlockHeader), arg0)
}

// BroadcastReceiver mocks base method
func (m *MockRemote) BroadcastReceiver() network.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BroadcastReceiver")
	ret0, _ := ret[0].(network.Client)
	return ret0
}

// BroadcastReceiver indicates an expected call of BroadcastReceiver
func (mr *MockRemoteMockRecorder) BroadcastReceiver() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BroadcastReceiver", reflect.TypeOf((*MockRemote)(nil).BroadcastReceiver))
}

// Close mocks base method
func (m *MockRemote) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockRemoteMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRemote)(nil).Close))
}

// CommandSender mocks base method
func (m *MockRemote) CommandSender() network.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CommandSender")
	ret0, _ := ret[0].(network.Client)
	return ret0
}

// CommandSender indicates an expected call of CommandSender
func (mr *MockRemoteMockRecorder) CommandSender() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CommandSender", reflect.TypeOf((*MockRemote)(nil).CommandSender))
}

// Height mocks base method
func (m *MockRemote) Height() (*communication.HeightResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Height")
	ret0, _ := ret[0].(*communication.HeightResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Height indicates an expected call of Height
func (mr *MockRemoteMockRecorder) Height() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Height", reflect.TypeOf((*MockRemote)(nil).Height))
}

// Info mocks base method
func (m *MockRemote) Info() (*communication.InfoResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Info")
	ret0, _ := ret[0].(*communication.InfoResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Info indicates an expected call of Info
func (mr *MockRemoteMockRecorder) Info() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockRemote)(nil).Info))
}