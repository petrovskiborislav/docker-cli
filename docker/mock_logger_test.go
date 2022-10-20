// Code generated by mockery v2.14.0. DO NOT EDIT.

package docker_test

import mock "github.com/stretchr/testify/mock"

// mockLogger is an autogenerated mock type for the Logger type
type mockLogger struct {
	mock.Mock
}

// Error provides a mock function with given fields: msgFormat, params
func (_m *mockLogger) Error(msgFormat string, params ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msgFormat)
	_ca = append(_ca, params...)
	_m.Called(_ca...)
}

// Info provides a mock function with given fields: msgFormat, params
func (_m *mockLogger) Info(msgFormat string, params ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msgFormat)
	_ca = append(_ca, params...)
	_m.Called(_ca...)
}

// Warn provides a mock function with given fields: msgFormat, params
func (_m *mockLogger) Warn(msgFormat string, params ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msgFormat)
	_ca = append(_ca, params...)
	_m.Called(_ca...)
}

type mockConstructorTestingTnewMockLogger interface {
	mock.TestingT
	Cleanup(func())
}

// newMockLogger creates a new instance of mockLogger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockLogger(t mockConstructorTestingTnewMockLogger) *mockLogger {
	mock := &mockLogger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
