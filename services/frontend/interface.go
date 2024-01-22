package frontend

type DispatchRequest struct {
	SessionID         uint
	RequestID         uint
	PickupLocationID  uint
	DropoffLocationID uint
}
