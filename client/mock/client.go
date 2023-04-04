// Code generated by MockGen. DO NOT EDIT.
// Source: client/client.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	types "github.com/obada-foundation/registry/types"
	asset "github.com/obada-foundation/sdkgo/asset"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockClient) Get(DID string) (types.DIDDocument, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", DID)
	ret0, _ := ret[0].(types.DIDDocument)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockClientMockRecorder) Get(DID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockClient)(nil).Get), DID)
}

// GetMetadataHistory mocks base method.
func (m *MockClient) GetMetadataHistory(DID string) (asset.DataArrayVersions, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetadataHistory", DID)
	ret0, _ := ret[0].(asset.DataArrayVersions)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetadataHistory indicates an expected call of GetMetadataHistory.
func (mr *MockClientMockRecorder) GetMetadataHistory(DID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetadataHistory", reflect.TypeOf((*MockClient)(nil).GetMetadataHistory), DID)
}

// Register mocks base method.
func (m *MockClient) Register(newDID types.RegisterDID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", newDID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockClientMockRecorder) Register(newDID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockClient)(nil).Register), newDID)
}

// SaveMetadata mocks base method.
func (m *MockClient) SaveMetadata(DID string, md types.SaveMetadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMetadata", DID, md)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveMetadata indicates an expected call of SaveMetadata.
func (mr *MockClientMockRecorder) SaveMetadata(DID, md interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMetadata", reflect.TypeOf((*MockClient)(nil).SaveMetadata), DID, md)
}
