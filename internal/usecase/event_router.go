package usecase

import (
	"context"
	"fmt"

	"github.com/aetherbus/aetherbus-tachyon/internal/domain"
	"github.com/aetherbus/aetherbus-tachyon/pkg/errors"
)

// EventRouter is the core application logic for routing events.
// It orchestrates the interaction between the delivery layer and the repository layer.
type EventRouter struct {
	routeStore domain.RouteStore
}

// NewEventRouter creates a new EventRouter.
func NewEventRouter(routeStore domain.RouteStore) *EventRouter {
	return &EventRouter{
		routeStore: routeStore,
	}
}

// Publish implements the domain.EventPublisher interface.
// It receives an event, finds the destination, and "sends" it.
func (r *EventRouter) Publish(ctx context.Context, event domain.Event) error {
	// 1. Match the topic to a destination node ID.
	destNodeID := r.routeStore.Match(event.Topic)
	if destNodeID == "" {
		// Log the warning and return a specific error.
		fmt.Printf("WARN: No route found for topic '%s'\n", event.Topic)
		return errors.ErrNoRouteFound
	}

	// 2. In a real system, you would now use the destNodeID to
	//    look up the node's address and publish the event to it
	//    over the network (e.g., using the ZMQ DEALER-ROUTER pattern).
	//    For this example, we'll just simulate this action.

	// fmt.Printf("Routing event %s for topic %s to node %s\n", event.ID, event.Topic, destNodeID)

	return nil
}
