package config

import (
	"os"
)

func SignadotSandboxName() string {
	return os.Getenv("SIGNADOT_SANDBOX_NAME")
}

// This sets the consumer group with suffix '-' + <routing key>
// if running in fork.  otherwise, it just returns the argument.
func SignadotConsumerGroup(groupID string) string {
	sandboxRoutingKey := os.Getenv("SIGNADOT_SANDBOX_ROUTING_KEY")
	if sandboxRoutingKey != "" {
		return groupID + "-" + sandboxRoutingKey
	}
	return groupID
}
