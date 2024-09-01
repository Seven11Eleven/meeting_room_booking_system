// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/Seven11Eleven/meeting_room_booking_system/internal/domain/models"
	mock "github.com/stretchr/testify/mock"
)

// ReservationService is an autogenerated mock type for the ReservationService type
type ReservationService struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, reservation
func (_m *ReservationService) Create(ctx context.Context, reservation *models.Reservation) error {
	ret := _m.Called(ctx, reservation)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Reservation) error); ok {
		r0 = rf(ctx, reservation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteReservation provides a mock function with given fields: ctx, reservation
func (_m *ReservationService) DeleteReservation(ctx context.Context, reservation *models.Reservation) error {
	ret := _m.Called(ctx, reservation)

	if len(ret) == 0 {
		panic("no return value specified for DeleteReservation")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Reservation) error); ok {
		r0 = rf(ctx, reservation)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetByRoomID provides a mock function with given fields: ctx, roomID
func (_m *ReservationService) GetByRoomID(ctx context.Context, roomID string) (*models.RoomReservations, error) {
	ret := _m.Called(ctx, roomID)

	if len(ret) == 0 {
		panic("no return value specified for GetByRoomID")
	}

	var r0 *models.RoomReservations
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.RoomReservations, error)); ok {
		return rf(ctx, roomID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.RoomReservations); ok {
		r0 = rf(ctx, roomID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.RoomReservations)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, roomID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewReservationService creates a new instance of ReservationService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewReservationService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ReservationService {
	mock := &ReservationService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
