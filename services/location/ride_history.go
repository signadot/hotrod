package location

import "time"

// RideHistoryEntry represents one requested ride persisted in storage.
type RideHistoryEntry struct {
	SessionID       uint      `json:"sessionID"`
	RequestID       uint      `json:"requestID"`
	PickupLocation  string    `json:"pickupLocation"`
	DropoffLocation string    `json:"dropoffLocation"`
	RequestedAt     time.Time `json:"requestedAt"`
	DriverPlate     string    `json:"driverPlate"`
}
