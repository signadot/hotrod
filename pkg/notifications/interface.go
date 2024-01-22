package notifications

import (
	"context"
	"time"
)

type Notification struct {
	ID        string               `json:"id"`
	Timestamp time.Time            `json:"timestamp"`
	Context   *NotificationContext `json:"context"`
	Body      string               `json:"body"`
}

type NotificationContext struct {
	Request          *RequestContext `json:"request"`
	RoutingKey       string          `json:"routingKey"`
	BaselineWorkload string          `json:"baselineWorkload"`
	SandboxName      string          `json:"sandboxName"`
}

type RequestContext struct {
	ID                uint `json:"id"`
	SessionID         uint `json:"sessionID"`
	PickupLocationID  uint `json:"pickupLocationID"`
	DropoffLocationID uint `json:"dropoffLocationID"`
}

type NotificationList struct {
	Cursor        int            `json:"cursor"`
	Notifications []Notification `json:"notifications"`
}

// Interface exposed by the Driver service.
type Interface interface {
	NotificationContext(reqCtx *RequestContext, routingKey string) *NotificationContext
	List(ctx context.Context, sessionID uint, cursor int) (*NotificationList, error)
	Store(ctx context.Context, notification *Notification) error
}
