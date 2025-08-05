package frontend

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/dnwe/otelsarama"
	"github.com/signadot/hotrod/pkg/kafka"
	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/tracing"
	"github.com/signadot/hotrod/services/driver"
	"github.com/signadot/hotrod/services/location"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type dispatcher struct {
	tracer   trace.Tracer
	logger   log.Factory
	location location.Interface
	producer sarama.SyncProducer
}

func newDispatcher(tracerProvider trace.TracerProvider, logger log.Factory,
	location location.Interface) *dispatcher {

	logger = logger.With(zap.String("component", "dispatcher"))

	// create a new tracer provider for kafka
	kafkaTracerProvider := tracing.InitOTEL("kafka", logger)
	producerTicker := time.NewTicker(100 * time.Millisecond)
	defer producerTicker.Stop()
	var (
		producer sarama.SyncProducer
		err      error
	)
	for {
		// get a kafka producer
		producer, err = kafka.GetSyncProducer("hotrod-frontend", kafkaTracerProvider)
		if err != nil {
			logger.Bg().Error("error getting kafka producer (will retry)")
			<-producerTicker.C
		} else {
			break
		}
	}

	return &dispatcher{
		tracer:   tracerProvider.Tracer("dispatcher"),
		logger:   logger,
		location: location,
		producer: producer,
	}
}

func (d *dispatcher) ResolveLocations(ctx context.Context,
	req *DispatchRequest) (*location.Location, *location.Location, error) {
	ctx, span := d.tracer.Start(ctx, "ResolveLocations", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.
			Key("session.id").
			Int(int(req.SessionID)),
	)
	defer span.End()

	// resolve locations
	pickupLoc, err := d.location.Get(ctx, int(req.PickupLocationID))
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't resolve pickup location, %w", err)
	}

	dropoffLoc, err := d.location.Get(ctx, int(req.DropoffLocationID))
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't resolve dropoff location, %w", err)
	}

	d.logger.For(ctx).Info("Found locations", zap.Any("pickupLocation", pickupLoc),
		zap.Any("dropoffLocation", dropoffLoc))
	return pickupLoc, dropoffLoc, nil
}

func (d *dispatcher) DispatchDriver(ctx context.Context, req *DispatchRequest,
	pickupLoc, dropoffLoc *location.Location) error {
	ctx, span := d.tracer.Start(ctx, "DispatchDriver", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.
			Key("session.id").
			Int(int(req.SessionID)),
	)
	defer span.End()

	driverRequest := driver.DispatchRequest{
		PickupLocation:  pickupLoc,
		DropoffLocation: dropoffLoc,
	}
	msgBody, err := json.Marshal(driverRequest)
	if err != nil {
		return fmt.Errorf("couldn't encode driver request, %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: kafka.DispatchDriverTopic,
		Key:   sarama.StringEncoder("message"),
		Value: sarama.StringEncoder(string(msgBody)),
	}
	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(msg))

	partition, offset, err := d.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("couldn't publish message to kafka, %w", err)
	}

	d.logger.For(ctx).Info("Published message to kafka",
		zap.Any("pickupLocation", pickupLoc),
		zap.Any("dropoffLocation", dropoffLoc),
		zap.Any("partition", partition),
		zap.Any("offset", offset),
	)
	return nil
}
