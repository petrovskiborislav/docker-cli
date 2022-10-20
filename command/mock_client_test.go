// Code generated by mockery v2.14.0. DO NOT EDIT.

package command_test

import (
	context "context"

	docker "github.com/petrovskiborislav/docker-cli/docker"
	mock "github.com/stretchr/testify/mock"
)

// mockClient is an autogenerated mock type for the Client type
type mockClient struct {
	mock.Mock
}

// ServiceDecommissioning provides a mock function with given fields: ctx, container
func (_m *mockClient) ServiceDecommissioning(ctx context.Context, container docker.Container) error {
	ret := _m.Called(ctx, container)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, docker.Container) error); ok {
		r0 = rf(ctx, container)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ServiceProvisioning provides a mock function with given fields: ctx, container
func (_m *mockClient) ServiceProvisioning(ctx context.Context, container docker.Container) error {
	ret := _m.Called(ctx, container)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, docker.Container) error); ok {
		r0 = rf(ctx, container)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTnewMockClient interface {
	mock.TestingT
	Cleanup(func())
}

// newMockClient creates a new instance of mockClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockClient(t mockConstructorTestingTnewMockClient) *mockClient {
	mock := &mockClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
