package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dnwe/otelsarama"
	"github.com/signadot/hotrod/pkg/baggageutils"
	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/notifications"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/IBM/sarama"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	tracer       trace.Tracer
	logger       log.Factory
	driverStore  *driverStore
	bestETA      *bestETA
	notification notifications.Interface
	ready        chan bool
}

func newConsumer(tracerProvider trace.TracerProvider, logger log.Factory) *Consumer {
	tracer := tracerProvider.Tracer("driver")
	return &Consumer{
		tracer:       tracer,
		logger:       logger,
		driverStore:  newDriverStore(tracer, logger),
		bestETA:      newBestETA(tracerProvider, tracer, logger),
		notification: notifications.NewNotificationHandler(tracerProvider, logger),
		ready:        make(chan bool),
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		session.MarkMessage(message, "")
		consumer.processDispatchRequest(message)
	}

	return nil
}

func (consumer *Consumer) processDispatchRequest(msg *sarama.ConsumerMessage) {
	// Extract tracing info from message
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(msg))

	// extract the request context
	reqContext, err := baggageutils.ExtractRequestContext(ctx)
	if err != nil {
		consumer.logger.For(ctx).Error("cannot extract request context from baggage", zap.Error(err))
		return
	}
	if reqContext == nil {
		consumer.logger.For(ctx).Error("empty request context from baggage, ignoring request")
		return
	}

	ctx, span := consumer.tracer.Start(ctx, "ProcessDispatchRequest", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		semconv.MessagingOperationProcess,
		attribute.
			Key("session.id").
			Int(int(reqContext.SessionID)),
	)
	defer span.End()

	// parse the message body
	var dispatchReq DispatchRequest
	err = json.Unmarshal(msg.Value, &dispatchReq)
	if err != nil {
		consumer.logger.For(ctx).Error("error decoding message body")
		span.SetStatus(codes.Error, err.Error())
		return
	}

	// send a notification
	notificationCtx := consumer.notification.NotificationContext(
		reqContext, baggageutils.GetRoutingKey(ctx))
	consumer.notification.Store(ctx, &notifications.Notification{
		ID:        fmt.Sprintf("req-%d-finding-driver", reqContext.ID),
		Timestamp: time.Now(),
		Context:   notificationCtx,
		Body:      "Finding an available driver",
	})

	// find availanble drivers
	driverIDs := consumer.driverStore.FindDriverIDs(ctx)

	// populate current driver's locations
	drivers, err := consumer.driverStore.GetDriversLocation(ctx, driverIDs)
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		return
	}

	// get the driver with the best ETA
	bestDriver, err := consumer.bestETA.Get(ctx, &dispatchReq, drivers)
	if err != nil {
		consumer.logger.For(ctx).Error("error calculating best route")
		span.SetStatus(codes.Error, err.Error())
		return
	}

	// send a notification
	consumer.notification.Store(ctx, &notifications.Notification{
		ID:        fmt.Sprintf("req-%d-dispatched-driver", reqContext.ID),
		Timestamp: time.Now(),
		Context:   notificationCtx,
		Body:      fmt.Sprintf("Driver %s arriving in %s", bestDriver.DriverID, bestDriver.ETA.String()),
	})
}
