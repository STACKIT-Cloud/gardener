// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gardener/gardener/pkg/component (interfaces: Deployer,Waiter,DeployWaiter,DeployMigrateWaiter)
//
// Generated by this command:
//
//	mockgen -package mock -destination=mocks.go github.com/gardener/gardener/pkg/component Deployer,Waiter,DeployWaiter,DeployMigrateWaiter
//

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	v1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	gomock "go.uber.org/mock/gomock"
)

// MockDeployer is a mock of Deployer interface.
type MockDeployer struct {
	ctrl     *gomock.Controller
	recorder *MockDeployerMockRecorder
}

// MockDeployerMockRecorder is the mock recorder for MockDeployer.
type MockDeployerMockRecorder struct {
	mock *MockDeployer
}

// NewMockDeployer creates a new mock instance.
func NewMockDeployer(ctrl *gomock.Controller) *MockDeployer {
	mock := &MockDeployer{ctrl: ctrl}
	mock.recorder = &MockDeployerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeployer) EXPECT() *MockDeployerMockRecorder {
	return m.recorder
}

// Deploy mocks base method.
func (m *MockDeployer) Deploy(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deploy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deploy indicates an expected call of Deploy.
func (mr *MockDeployerMockRecorder) Deploy(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deploy", reflect.TypeOf((*MockDeployer)(nil).Deploy), arg0)
}

// Destroy mocks base method.
func (m *MockDeployer) Destroy(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy.
func (mr *MockDeployerMockRecorder) Destroy(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockDeployer)(nil).Destroy), arg0)
}

// MockWaiter is a mock of Waiter interface.
type MockWaiter struct {
	ctrl     *gomock.Controller
	recorder *MockWaiterMockRecorder
}

// MockWaiterMockRecorder is the mock recorder for MockWaiter.
type MockWaiterMockRecorder struct {
	mock *MockWaiter
}

// NewMockWaiter creates a new mock instance.
func NewMockWaiter(ctrl *gomock.Controller) *MockWaiter {
	mock := &MockWaiter{ctrl: ctrl}
	mock.recorder = &MockWaiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWaiter) EXPECT() *MockWaiterMockRecorder {
	return m.recorder
}

// Wait mocks base method.
func (m *MockWaiter) Wait(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockWaiterMockRecorder) Wait(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockWaiter)(nil).Wait), arg0)
}

// WaitCleanup mocks base method.
func (m *MockWaiter) WaitCleanup(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitCleanup", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitCleanup indicates an expected call of WaitCleanup.
func (mr *MockWaiterMockRecorder) WaitCleanup(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitCleanup", reflect.TypeOf((*MockWaiter)(nil).WaitCleanup), arg0)
}

// MockDeployWaiter is a mock of DeployWaiter interface.
type MockDeployWaiter struct {
	ctrl     *gomock.Controller
	recorder *MockDeployWaiterMockRecorder
}

// MockDeployWaiterMockRecorder is the mock recorder for MockDeployWaiter.
type MockDeployWaiterMockRecorder struct {
	mock *MockDeployWaiter
}

// NewMockDeployWaiter creates a new mock instance.
func NewMockDeployWaiter(ctrl *gomock.Controller) *MockDeployWaiter {
	mock := &MockDeployWaiter{ctrl: ctrl}
	mock.recorder = &MockDeployWaiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeployWaiter) EXPECT() *MockDeployWaiterMockRecorder {
	return m.recorder
}

// Deploy mocks base method.
func (m *MockDeployWaiter) Deploy(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deploy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deploy indicates an expected call of Deploy.
func (mr *MockDeployWaiterMockRecorder) Deploy(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deploy", reflect.TypeOf((*MockDeployWaiter)(nil).Deploy), arg0)
}

// Destroy mocks base method.
func (m *MockDeployWaiter) Destroy(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy.
func (mr *MockDeployWaiterMockRecorder) Destroy(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockDeployWaiter)(nil).Destroy), arg0)
}

// Wait mocks base method.
func (m *MockDeployWaiter) Wait(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockDeployWaiterMockRecorder) Wait(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockDeployWaiter)(nil).Wait), arg0)
}

// WaitCleanup mocks base method.
func (m *MockDeployWaiter) WaitCleanup(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitCleanup", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitCleanup indicates an expected call of WaitCleanup.
func (mr *MockDeployWaiterMockRecorder) WaitCleanup(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitCleanup", reflect.TypeOf((*MockDeployWaiter)(nil).WaitCleanup), arg0)
}

// MockDeployMigrateWaiter is a mock of DeployMigrateWaiter interface.
type MockDeployMigrateWaiter struct {
	ctrl     *gomock.Controller
	recorder *MockDeployMigrateWaiterMockRecorder
}

// MockDeployMigrateWaiterMockRecorder is the mock recorder for MockDeployMigrateWaiter.
type MockDeployMigrateWaiterMockRecorder struct {
	mock *MockDeployMigrateWaiter
}

// NewMockDeployMigrateWaiter creates a new mock instance.
func NewMockDeployMigrateWaiter(ctrl *gomock.Controller) *MockDeployMigrateWaiter {
	mock := &MockDeployMigrateWaiter{ctrl: ctrl}
	mock.recorder = &MockDeployMigrateWaiterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeployMigrateWaiter) EXPECT() *MockDeployMigrateWaiterMockRecorder {
	return m.recorder
}

// Deploy mocks base method.
func (m *MockDeployMigrateWaiter) Deploy(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Deploy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Deploy indicates an expected call of Deploy.
func (mr *MockDeployMigrateWaiterMockRecorder) Deploy(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deploy", reflect.TypeOf((*MockDeployMigrateWaiter)(nil).Deploy), arg0)
}

// Destroy mocks base method.
func (m *MockDeployMigrateWaiter) Destroy(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy.
func (mr *MockDeployMigrateWaiterMockRecorder) Destroy(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockDeployMigrateWaiter)(nil).Destroy), arg0)
}

// Migrate mocks base method.
func (m *MockDeployMigrateWaiter) Migrate(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Migrate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Migrate indicates an expected call of Migrate.
func (mr *MockDeployMigrateWaiterMockRecorder) Migrate(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Migrate", reflect.TypeOf((*MockDeployMigrateWaiter)(nil).Migrate), arg0)
}

// Restore mocks base method.
func (m *MockDeployMigrateWaiter) Restore(arg0 context.Context, arg1 *v1beta1.ShootState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Restore", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Restore indicates an expected call of Restore.
func (mr *MockDeployMigrateWaiterMockRecorder) Restore(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Restore", reflect.TypeOf((*MockDeployMigrateWaiter)(nil).Restore), arg0, arg1)
}

// Wait mocks base method.
func (m *MockDeployMigrateWaiter) Wait(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Wait", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Wait indicates an expected call of Wait.
func (mr *MockDeployMigrateWaiterMockRecorder) Wait(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Wait", reflect.TypeOf((*MockDeployMigrateWaiter)(nil).Wait), arg0)
}

// WaitCleanup mocks base method.
func (m *MockDeployMigrateWaiter) WaitCleanup(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitCleanup", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitCleanup indicates an expected call of WaitCleanup.
func (mr *MockDeployMigrateWaiterMockRecorder) WaitCleanup(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitCleanup", reflect.TypeOf((*MockDeployMigrateWaiter)(nil).WaitCleanup), arg0)
}

// WaitMigrate mocks base method.
func (m *MockDeployMigrateWaiter) WaitMigrate(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WaitMigrate", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WaitMigrate indicates an expected call of WaitMigrate.
func (mr *MockDeployMigrateWaiterMockRecorder) WaitMigrate(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitMigrate", reflect.TypeOf((*MockDeployMigrateWaiter)(nil).WaitMigrate), arg0)
}
