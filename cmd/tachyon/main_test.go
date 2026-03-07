package main

import (
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/aetherbus/aetherbus-tachyon/internal/domain"
	"github.com/aetherbus/aetherbus-tachyon/internal/media"
	"github.com/pebbe/zmq4"
)

func TestMainIntegration(t *testing.T) {
	// 1. Setup: Start the main function in a goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		main()
	}()

	// Give the server a moment to start
	time.Sleep(250 * time.Millisecond)

	// 2. Test Execution: Connect a client and send/receive a message
	// Create a ZMQ SUB socket to listen for published events
	subSocket, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		t.Fatalf("Failed to create SUB socket: %v", err)
	}
	defer subSocket.Close()

	if err := subSocket.Connect("tcp://127.0.0.1:5556"); err != nil {
		t.Fatalf("Failed to connect SUB socket: %v", err)
	}
	subSocket.SetSubscribe("user.created")

	// Create a ZMQ DEALER socket to send an event
	dealerSocket, err := zmq4.NewSocket(zmq4.DEALER)
	if err != nil {
		t.Fatalf("Failed to create DEALER socket: %v", err)
	}
	defer dealerSocket.Close()

	if err := dealerSocket.Connect("tcp://127.0.0.1:5555"); err != nil {
		t.Fatalf("Failed to connect DEALER socket: %v", err)
	}

	testEvent := domain.Event{
		ID:              "evt-123",
		Source:          "test-client",
		Data:            map[string]any{"id": "123", "name": "test"},
		DataContentType: "application/json",
		SpecVersion:     "1.0",
	}

	codec := media.NewJSONCodec()
	compressor := media.NewLZ4Compressor()

	encodedEvent, err := codec.Encode(testEvent)
	if err != nil {
		t.Fatalf("Failed to encode event: %v", err)
	}

	compressedEvent, err := compressor.Compress(encodedEvent)
	if err != nil {
		t.Fatalf("Failed to compress event: %v", err)
	}

	// Send message following the router protocol: [Delimiter, Topic, Payload]
	// DEALER automatically prefixes its identity for ROUTER.
	_, err = dealerSocket.SendMessage("", "user.created", compressedEvent)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Receive the message from the SUB socket
	msg, err := subSocket.RecvMessageBytes(0)
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}

	// 3. Assertion: Verify the received message
	if len(msg) != 2 {
		t.Fatalf("Expected 2 parts in message, got %d", len(msg))
	}

	expectedTopic := "user.created"
	if string(msg[0]) != expectedTopic {
		t.Errorf("Expected topic '%s', got '%s'", expectedTopic, string(msg[0]))
	}

	var receivedEvent domain.Event
	if err := codec.Decode(msg[1], &receivedEvent); err != nil {
		t.Fatalf("Failed to decode received payload: %v", err)
	}

	if receivedEvent.ID != testEvent.ID {
		t.Errorf("Expected event ID '%s', got '%s'", testEvent.ID, receivedEvent.ID)
	}

	if receivedEvent.Source != testEvent.Source {
		t.Errorf("Expected source '%s', got '%s'", testEvent.Source, receivedEvent.Source)
	}

	if receivedEvent.Topic != expectedTopic {
		t.Errorf("Expected topic '%s' in event, got '%s'", expectedTopic, receivedEvent.Topic)
	}

	// 4. Teardown: Gracefully shut down the server
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGINT)

	wg.Wait()
}
