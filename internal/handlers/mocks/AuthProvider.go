// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// AuthProvider is an autogenerated mock type for the AuthProvider type
type AuthProvider struct {
	mock.Mock
}

// SignUp provides a mock function with given fields: ctx, login, password
func (_m *AuthProvider) SignUp(ctx context.Context, login string, password string) (string, error) {
	ret := _m.Called(ctx, login, password)

	if len(ret) == 0 {
		panic("no return value specified for SignUp")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (string, error)); ok {
		return rf(ctx, login, password)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) string); ok {
		r0 = rf(ctx, login, password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, login, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAuthProvider creates a new instance of AuthProvider. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAuthProvider(t interface {
	mock.TestingT
	Cleanup(func())
}) *AuthProvider {
	mock := &AuthProvider{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}