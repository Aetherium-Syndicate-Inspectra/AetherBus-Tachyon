package domain

import "context"

type RouteStore interface {
	AddRoute(topic string, destNodeID string)
	Match(topic string) []string
}

type EventPublisher interface {
	Publish(ctx context.Context, event Event) error
}
