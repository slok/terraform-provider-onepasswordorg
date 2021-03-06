// Code generated by mockery v2.9.4. DO NOT EDIT.

package onepasswordclimock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// OpCli is an autogenerated mock type for the OpCli type
type OpCli struct {
	mock.Mock
}

// RunOpCmd provides a mock function with given fields: ctx, args
func (_m *OpCli) RunOpCmd(ctx context.Context, args []string) (string, string, error) {
	ret := _m.Called(ctx, args)

	var r0 string
	if rf, ok := ret.Get(0).(func(context.Context, []string) string); ok {
		r0 = rf(ctx, args)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, []string) string); ok {
		r1 = rf(ctx, args)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, []string) error); ok {
		r2 = rf(ctx, args)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}
