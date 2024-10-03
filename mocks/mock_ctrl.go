// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/hdl/http/http.go
//
// Generated by this command:
//
//	mockgen -source=./internal/hdl/http/http.go -destination=mocks/mock_ctrl.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/JMURv/effectiveMobile/pkg/model"
	gomock "go.uber.org/mock/gomock"
)

// MockCtrl is a mock of Ctrl interface.
type MockCtrl struct {
	ctrl     *gomock.Controller
	recorder *MockCtrlMockRecorder
}

// MockCtrlMockRecorder is the mock recorder for MockCtrl.
type MockCtrlMockRecorder struct {
	mock *MockCtrl
}

// NewMockCtrl creates a new mock instance.
func NewMockCtrl(ctrl *gomock.Controller) *MockCtrl {
	mock := &MockCtrl{ctrl: ctrl}
	mock.recorder = &MockCtrlMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCtrl) EXPECT() *MockCtrlMockRecorder {
	return m.recorder
}

// CreateSong mocks base method.
func (m *MockCtrl) CreateSong(ctx context.Context, req *model.Song) (uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSong", ctx, req)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSong indicates an expected call of CreateSong.
func (mr *MockCtrlMockRecorder) CreateSong(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSong", reflect.TypeOf((*MockCtrl)(nil).CreateSong), ctx, req)
}

// DeleteSong mocks base method.
func (m *MockCtrl) DeleteSong(ctx context.Context, id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSong", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSong indicates an expected call of DeleteSong.
func (mr *MockCtrlMockRecorder) DeleteSong(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSong", reflect.TypeOf((*MockCtrl)(nil).DeleteSong), ctx, id)
}

// GetSong mocks base method.
func (m *MockCtrl) GetSong(ctx context.Context, id uint64, page, size int) (*model.PaginatedSongs, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSong", ctx, id, page, size)
	ret0, _ := ret[0].(*model.PaginatedSongs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSong indicates an expected call of GetSong.
func (mr *MockCtrlMockRecorder) GetSong(ctx, id, page, size any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSong", reflect.TypeOf((*MockCtrl)(nil).GetSong), ctx, id, page, size)
}

// ListSongs mocks base method.
func (m *MockCtrl) ListSongs(ctx context.Context, page, size int, filters map[string]any) (*model.PaginatedSongs, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSongs", ctx, page, size, filters)
	ret0, _ := ret[0].(*model.PaginatedSongs)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSongs indicates an expected call of ListSongs.
func (mr *MockCtrlMockRecorder) ListSongs(ctx, page, size, filters any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSongs", reflect.TypeOf((*MockCtrl)(nil).ListSongs), ctx, page, size, filters)
}

// UpdateSong mocks base method.
func (m *MockCtrl) UpdateSong(ctx context.Context, req *model.Song) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSong", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateSong indicates an expected call of UpdateSong.
func (mr *MockCtrlMockRecorder) UpdateSong(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSong", reflect.TypeOf((*MockCtrl)(nil).UpdateSong), ctx, req)
}
