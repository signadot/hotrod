package driver

import "time"

import (
	"github.com/signadot/hotrod/services/location"
)

type DispatchRequest struct {
	PickupLocation  *location.Location `json:"pickupLocation"`
	DropoffLocation *location.Location `json:"dropoffLocation"`
	RequestID       uint               `json:"requestID"`
	RequestedAt     time.Time          `json:"requestedAt"`
}

type Driver struct {
	DriverID    string
	Coordinates string
}
