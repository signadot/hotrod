package notifications

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/signadot/hotrod/pkg/config"
	"github.com/signadot/hotrod/pkg/log"
	"github.com/signadot/hotrod/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type notificationHandler struct {
	tracer trace.Tracer
	logger log.Factory
	rdb    *redis.Client
}

func NewNotificationHandler(tracerProvider trace.TracerProvider, logger log.Factory) Interface {
	logger = logger.With(zap.String("component", "notifications"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.GetRedisAddr(),
		Password: config.GetRedisPassword(),
		DB:       0, // use default DB
	})
	// create a new tracer provider for redis
	redisTracerProvider := tracing.InitOTEL("redis", logger)
	if err := redisotel.InstrumentTracing(rdb,
		redisotel.WithTracerProvider(redisTracerProvider)); err != nil {
		panic(err)
	}

	return &notificationHandler{
		tracer: tracerProvider.Tracer("notifications"),
		logger: logger,
		rdb:    rdb,
	}
}

func (h *notificationHandler) NotificationContext(reqCtx *RequestContext,
	routingKey string) *NotificationContext {
	return &NotificationContext{
		Request:          reqCtx,
		RoutingKey:       routingKey,
		BaselineWorkload: config.SignadotBaselineName(),
		SandboxName:      config.SignadotSandboxName(),
	}
}

func (h *notificationHandler) Store(ctx context.Context, notification *Notification) error {
	ctx, span := h.tracer.Start(ctx, "StoreNotification", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.
			Key("session.id").
			Int(int(notification.Context.Request.SessionID)),
		attribute.
			Key("notification.id").
			String(notification.ID),
	)
	defer span.End()

	// get the session key
	sessionKey := getSessionKey(notification.Context.Request.SessionID)

	// apply the changes inside a transaction
	return h.rdb.Watch(ctx, func(tx *redis.Tx) error {
		// get stored notifications for the current session
		r := tx.Get(ctx, sessionKey)
		if r.Err() != nil && r.Err() != redis.Nil {
			return r.Err()
		}

		// parse notifications
		var notifications []Notification
		if r.Err() != redis.Nil {
			err := json.Unmarshal([]byte(r.Val()), &notifications)
			if err != nil {
				return err
			}
		}

		// append the new notification (if it doesn't exist)
		for _, n := range notifications {
			if n.ID == notification.ID {
				// we already have this notification, just return
				return nil
			}
		}
		notifications = append(notifications, *notification)
		jsonData, err := json.Marshal(notifications)
		if err != nil {
			return err
		}

		// run the transaction
		h.logger.For(ctx).Info("Storing notification", zap.Any("notification", *notification))
		pipe := tx.Pipeline()
		pipe.SetEx(ctx, sessionKey, jsonData, 30*time.Second)
		_, err = pipe.Exec(context.Background())
		return err
	}, sessionKey)
}

func (h *notificationHandler) List(ctx context.Context, sessionID uint, cursor int) (*NotificationList, error) {
	ctx, span := h.tracer.Start(ctx, "ListNotifications", trace.WithSpanKind(trace.SpanKindClient))
	span.SetAttributes(
		attribute.
			Key("session.id").
			Int(int(sessionID)),
		attribute.
			Key("cursor").
			Int(int(cursor)),
	)
	defer span.End()

	// get the session key
	sessionKey := getSessionKey(sessionID)

	r := h.rdb.Get(ctx, sessionKey)
	if r.Err() != nil && r.Err() != redis.Nil {
		return nil, r.Err()
	}

	// parse notifications
	var notifications []Notification
	if r.Err() == redis.Nil {
		notifications = make([]Notification, 0)
	} else {
		err := json.Unmarshal([]byte(r.Val()), &notifications)
		if err != nil {
			return nil, err
		}
	}

	// check cursor
	if cursor > len(notifications) {
		// reset cursor
		cursor = -1
	}

	// populate the response
	resp := NotificationList{
		Cursor:        cursor,
		Notifications: []Notification{},
	}
	for i, n := range notifications {
		if i <= int(cursor) {
			continue
		}
		resp.Cursor = i
		resp.Notifications = append(resp.Notifications, n)
	}
	return &resp, nil
}

func getSessionKey(sessionID uint) string {
	return fmt.Sprintf("session:%d", sessionID)
}
