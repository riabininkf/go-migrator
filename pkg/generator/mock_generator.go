// Code generated by mockery v1.0.0. DO NOT EDIT.

package generator

import mock "github.com/stretchr/testify/mock"

// MockGenerator is an autogenerated mock type for the Generator type
type MockGenerator struct {
	mock.Mock
}

// Generate provides a mock function with given fields: name, path
func (_m *MockGenerator) Generate(name string, path string) error {
	ret := _m.Called(name, path)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(name, path)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}