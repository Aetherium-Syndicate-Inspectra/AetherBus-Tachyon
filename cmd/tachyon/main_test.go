// cmd/tachyon/main_test.go
package main

import (
	"context"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

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
	// Use the same context to ensure cleanup
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

	// Send a test message
	testPayload := `{"id":"123","name":"test"}`
	_, err = dealerSocket.SendMessage(testPayload)
	if err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// Receive the message from the SUB socket
	msg, err := subSocket.RecvMessage(0)
	if err != nil {
		t.Fatalf("Failed to receive message: %v", err)
	}

	// 3. Assertion: Verify the received message
	if len(msg) != 2 {
		t.Fatalf("Expected 2 parts in message, got %d", len(msg))
	}

	expectedTopic := "user.created"
	if msg[0] != expectedTopic {
		t.Errorf("Expected topic '%s', got '%s'", expectedTopic, msg[0])
	}

	if msg[1] != testPayload {
		t.Errorf("Expected payload '%s', got '%s'", testPayload, msg[1])
	}

	// 4. Teardown: Gracefully shut down the server
	// Send a SIGINT signal to the current process to trigger shutdown
	// This simulates Ctrl+C
	p, _ := os.FindProcess(os.Getpid())
	p.Signal(syscall.SIGINT)

	// Wait for the main function to exit
	wg.Wait()
}
