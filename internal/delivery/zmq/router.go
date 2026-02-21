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
	codec       domain.Codec
	compressor  domain.Compressor
}

// NewRouter creates a new ZMQ Router.
func NewRouter(bindAddress, pubAddress string, publisher domain.EventPublisher, codec domain.Codec, compressor domain.Compressor) *Router {
	return &Router{
		bindAddress: bindAddress,
		pubAddress:  pubAddress,
		publisher:   publisher,
		codec:       codec,
		compressor:  compressor,
	}
}

// Start initializes and runs the ZMQ ROUTER socket loop.
func (r *Router) Start(ctx context.Context) error {
	// ... (socket creation and binding code remains the same)
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
			if ctx.Err() != nil {
				break
			}
			continue
		}

		if len(sockets) > 0 {
			msg, err := r.routerSocket.RecvMessageBytes(0)
            // Expect: [ClientID, Delimiter, Topic, Payload]
			if err != nil || len(msg) < 4 {
				continue // Malformed or error
			}

			clientID := msg[0]
            // msg[1] is the empty delimiter
			topic := string(msg[2])
			rawEvent := msg[3]

			// 1. Decompress
			decompressedEvent, err := r.compressor.Decompress(rawEvent)
			if err != nil {
				// Log error: failed to decompress
				continue
			}

			// 2. Decode
			var event domain.Event
			if err := r.codec.Decode(decompressedEvent, &event); err != nil {
				// Log error: failed to decode
				continue
			}
            
            // Assign the topic from the message frame
            event.Topic = topic

			envelope := domain.Envelope{
				ClientID: clientID,
				Event:    event,
			}

			// 3. Publish to the application logic
			if err := r.publisher.Publish(ctx, envelope); err == nil {
				// If successful, publish to the PUB socket for subscribers
				r.pubSocket.SendMessage(event.Topic, decompressedEvent)
			}
		}

		if ctx.Err() != nil {
			break
		}
	}
}
