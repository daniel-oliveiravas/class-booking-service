// Code generated by mockery v2.23.1. DO NOT EDIT.

package mocks

import (
	context "context"

	classes "github.com/daniel-oliveiravas/class-booking-service/business/classes"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// Add provides a mock function with given fields: ctx, class
func (_m *Repository) Add(ctx context.Context, class classes.Class) (classes.Class, error) {
	ret := _m.Called(ctx, class)

	var r0 classes.Class
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, classes.Class) (classes.Class, error)); ok {
		return rf(ctx, class)
	}
	if rf, ok := ret.Get(0).(func(context.Context, classes.Class) classes.Class); ok {
		r0 = rf(ctx, class)
	} else {
		r0 = ret.Get(0).(classes.Class)
	}

	if rf, ok := ret.Get(1).(func(context.Context, classes.Class) error); ok {
		r1 = rf(ctx, class)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, classID
func (_m *Repository) Delete(ctx context.Context, classID string) error {
	ret := _m.Called(ctx, classID)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, classID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByID provides a mock function with given fields: ctx, classID
func (_m *Repository) GetByID(ctx context.Context, classID string) (classes.Class, error) {
	ret := _m.Called(ctx, classID)

	var r0 classes.Class
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (classes.Class, error)); ok {
		return rf(ctx, classID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) classes.Class); ok {
		r0 = rf(ctx, classID)
	} else {
		r0 = ret.Get(0).(classes.Class)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, classID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsNotFoundErr provides a mock function with given fields: err
func (_m *Repository) IsNotFoundErr(err error) bool {
	ret := _m.Called(err)

	var r0 bool
	if rf, ok := ret.Get(0).(func(error) bool); ok {
		r0 = rf(err)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// List provides a mock function with given fields: ctx, limit, offset
func (_m *Repository) List(ctx context.Context, limit int, offset int) ([]classes.Class, error) {
	ret := _m.Called(ctx, limit, offset)

	var r0 []classes.Class
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int, int) ([]classes.Class, error)); ok {
		return rf(ctx, limit, offset)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int) []classes.Class); ok {
		r0 = rf(ctx, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]classes.Class)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int) error); ok {
		r1 = rf(ctx, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, classID, updateClass
func (_m *Repository) Update(ctx context.Context, classID string, updateClass classes.UpdateClass) (classes.Class, error) {
	ret := _m.Called(ctx, classID, updateClass)

	var r0 classes.Class
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, classes.UpdateClass) (classes.Class, error)); ok {
		return rf(ctx, classID, updateClass)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, classes.UpdateClass) classes.Class); ok {
		r0 = rf(ctx, classID, updateClass)
	} else {
		r0 = ret.Get(0).(classes.Class)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, classes.UpdateClass) error); ok {
		r1 = rf(ctx, classID, updateClass)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}