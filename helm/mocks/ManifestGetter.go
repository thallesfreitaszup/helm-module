// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ManifestGetter is an autogenerated mock type for the ManifestGetter type
type ManifestGetter struct {
	mock.Mock
}

// Get provides a mock function with given fields:
func (_m *ManifestGetter) Get() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
