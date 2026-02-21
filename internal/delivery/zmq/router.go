package zmq

import (
	"context"
	"fmt"
	"time"

	"github.com/aetherbus/aetherbus-tachyon/internal/domain"
	"github.com/pebbe/zmq4"
)

// Router manages the ZMQ ROUTER socket for incoming events.
type Router struct {
	bindAddress string
	pubAddress  string
	publisher   domain.EventPublisher
	routerSocket *zmq4.Socket
	pubSocket    *zmq4.Socket
}

// NewRouter creates a new ZMQ Router.
func NewRouter(bindAddress, pubAddress string, publisher domain.EventPublisher) *Router {
	return &Router{
		bindAddress: bindAddress,
		pubAddress:  pubAddress,
		publisher:   publisher,
	}
}

// Start initializes and runs the ZMQ ROUTER socket loop.
func (r *Router) Start(ctx context.Context) error {
	// Create ROUTER socket
	routerSocket, err := zmq4.NewSocket(zmq4.ROUTER)
	if err != nil {
		return fmt.Errorf("failed to create router socket: %w", err)
	}
	r.routerSocket = routerSocket

	// Create PUB socket
	pubSocket, err := zmq4.NewSocket(zmq4.PUB)
	if err != nil {
		return fmt.Errorf("failed to create pub socket: %w", err)
	}
	r.pubSocket = pubSocket

	// Bind sockets
	if err := r.routerSocket.Bind(r.bindAddress); err != nil {
		return fmt.Errorf("failed to bind router socket: %w", err)
	}
	if err := r.pubSocket.Bind(r.pubAddress); err != nil {
		return fmt.Errorf("failed to bind pub socket: %w", err)
	}

	fmt.Println("ZMQ Router started")

	// Run the main loop in a goroutine
	go r.loop(ctx)

	return nil
}

// Stop gracefully closes the ZMQ socket.
func (r *Router) Stop() {
	if r.routerSocket != nil {
		r.routerSocket.Close()
	}
	if r.pubSocket != nil {
		r.pubSocket.Close()
	}
	fmt.Println("ZMQ Router stopped")
}

func (r *Router) loop(ctx context.Context) {
	defer r.Stop()

	poller := zmq4.NewPoller()
	poller.Add(r.routerSocket, zmq4.POLLIN)

	for {
		// Poll for events with a timeout
		sockets, err := poller.Poll(250 * time.Millisecond)
		if err != nil {
			// Break on context cancellation
			if ctx.Err() != nil {
				break
			}
			continue
		}

		if len(sockets) > 0 {
			// Read message from router socket
			msg, err := r.routerSocket.RecvMessage(0)
			if err != nil {
				// Handle error (e.g., log it)
				continue
			}

			// Basic validation
			if len(msg) < 2 {
				// Malformed message
				continue
			}

			// msg[0] is the client identity
			// msg[1] is the event data
			rawEvent := msg[1]

			// In a real implementation, you would deserialize this into your
			// domain.Event struct. For this example, we'll skip that and
			// just pass the raw bytes to the publisher.
			// We'll also invent a topic for routing.
			event := domain.Event{
				Topic:   "user.created", // This should be extracted from the message
				Payload: []byte(rawEvent),
			}

			// Publish to the application logic
			if err := r.publisher.Publish(ctx, event); err == nil {
				// If successful, publish to the PUB socket
				r.pubSocket.SendMessage(event.Topic, string(event.Payload))
			}
		}

		// Check for context cancellation
		if ctx.Err() != nil {
			break
		}
	}
}
