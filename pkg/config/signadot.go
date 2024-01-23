package config

import (
	"os"
)

func SignadotBaselineName() string {
	return os.Getenv("SIGNADOT_BASELINE_NAME")
}

func SignadotSandboxName() string {
	return os.Getenv("SIGNADOT_SANDBOX_NAME")
}

func SignadotSandboxRoutingKey() string {
	return os.Getenv("SIGNADOT_SANDBOX_ROUTING_KEY")
}

// This sets the consumer group with suffix '-' + <routing key>
// if running in fork.  otherwise, it just returns the argument.
func SignadotConsumerGroup(groupID string) string {
	sandboxRoutingKey := SignadotSandboxRoutingKey()
	if sandboxRoutingKey != "" {
		return groupID + "-" + sandboxRoutingKey
	}
	return groupID
}
