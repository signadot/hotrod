package baggageutils

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/signadot/hotrod/pkg/notifications"
	"go.opentelemetry.io/otel/baggage"
)

const (
	reqContextBaggageKey      = "reqContext"
	OpenTelemetryBaggageKeyV3 = "sd-routing-key"
)

func GetRoutingKey(ctx context.Context) string {
	bag := baggage.FromContext(ctx)
	return bag.Member(OpenTelemetryBaggageKeyV3).Value()
}

func InjectRequestContext(ctx context.Context,
	reqContext *notifications.RequestContext) (*baggage.Baggage, error) {
	reqContextData, err := json.Marshal(reqContext)
	if err != nil {
		return nil, err
	}

	m, err := baggage.NewMember(reqContextBaggageKey, base64.StdEncoding.EncodeToString(reqContextData))
	if err != nil {
		return nil, err
	}

	bag := baggage.FromContext(ctx)
	bag, err = bag.SetMember(m)
	if err != nil {
		return nil, err
	}
	return &bag, err
}

func ExtractRequestContext(ctx context.Context) (*notifications.RequestContext, error) {
	bag := baggage.FromContext(ctx)
	reqContextEncoded := bag.Member(reqContextBaggageKey).Value()
	if reqContextEncoded == "" {
		return nil, nil
	}

	reqContextBytes, err := base64.StdEncoding.DecodeString(reqContextEncoded)
	if err != nil {
		return nil, err
	}

	var reqContext notifications.RequestContext
	err = json.Unmarshal(reqContextBytes, &reqContext)
	if err != nil {
		return nil, err
	}
	return &reqContext, nil
}
