// pkg/client/publisher.go
package client

import (
	"context"
	"fmt"
)

func (c *tachyonClient) Publish(ctx context.Context, topic string, payload []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return fmt.Errorf("client is closed")
	}

	// The message envelope must match the ROUTER server's expectation:
	// Frame 1: Client Identity (added automatically by DEALER)
	// Frame 2: Empty delimiter (for ROUTER compatibility)
	// Frame 3: Topic
	// Frame 4: Payload
	_, err := c.dealer.SendMessage("", topic, payload)
	if err != nil {
		return fmt.Errorf("failed to publish message to topic '%s': %w", topic, err)
	}

	return nil
}
