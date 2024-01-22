package frontend

import (
	"context"
	"fmt"

	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/services/location"
	"go.uber.org/zap"
)

type dispatcher struct {
	logger   log.Factory
	location location.Interface
}

func newDispatcher(logger log.Factory, location location.Interface) *dispatcher {
	return &dispatcher{
		logger:   logger,
		location: location,
	}
}

func (d *dispatcher) DispatchDriver(ctx context.Context, req *DispatchRequest) error {
	// resolve locations
	pickupLocation, err := d.location.Get(ctx, int(req.PickupLocationID))
	if err != nil {
		return fmt.Errorf("couldn't resolve pickup location, %w", err)
	}

	dropoffLocation, err := d.location.Get(ctx, int(req.DropoffLocationID))
	if err != nil {
		return fmt.Errorf("couldn't resolve dropoff location, %w", err)
	}

	d.logger.For(ctx).Info("Found locations", zap.Any("pickupLocation", pickupLocation),
		zap.Any("dropoffLocation", dropoffLocation))

	return nil
}
