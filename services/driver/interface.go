package driver

import (
	"github.com/signadot/hotrod/services/location"
)

type DispatchRequest struct {
	PickupLocation  *location.Location `json:"pickupLocation"`
	DropoffLocation *location.Location `json:"dropoffLocation"`
}

type Driver struct {
	DriverID    string
	Coordinates string
}
