# AetherBus-Tachyon

AetherBus-Tachyon is a high-performance, lightweight event router built with Go and ZeroMQ. It serves as a central hub for routing messages (events) from producers to the appropriate consumers based on predefined topics.

It is designed for scenarios where you need a fast, reliable, and scalable message bus without the overhead of larger message brokers.

## Architecture

Tachyon uses a combination of ZMQ socket patterns to create a robust and flexible routing system:

1.  **`ROUTER` Socket (`tcp://*:5555`)**: This is the main entry point for incoming events from producers. Producers connect to this socket using a `DEALER` socket. The `ROUTER` socket allows Tachyon to receive messages from multiple clients and know the identity of the sender.

2.  **`PUB` (Publish) Socket (`tcp://*:5556`)**: After an event is received and validated, Tachyon publishes the event on this socket. Subscribers (consumers) can connect to this `PUB` socket using a `SUB` socket to receive events for the topics they are interested in.

This model decouples event producers from consumers, allowing for a scalable and resilient microservices architecture.

## Getting Started

### Prerequisites

*   **Go**: Version 1.22 or higher.
*   **Docker**: (Optional) For containerized deployment.
*   **ZMQ Library**: You need to have the ZeroMQ library installed on your system.
    *   **On macOS**: `brew install zeromq`
    *   **On Debian/Ubuntu**: `sudo apt-get install libzmq3-dev`

### Running Locally

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/aetherbus/aetherbus-tachyon.git
    cd aetherbus-tachyon
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Run the server:**
    ```bash
    go run ./cmd/tachyon
    ```
    The server will start and listen for connections on ports 5555 and 5556.

### Running with Docker

1.  **Build the Docker image:**
    ```bash
    docker build -t aetherbus-tachyon .
    ```

2.  **Run the Docker container:**
    ```bash
    docker run -p 5555:5555 -p 5556:5556 --name tachyon-server aetherbus-tachyon
    ```
    This will start the Tachyon server inside a container.

## Configuration

The application is configured via environment variables.

| Variable           | Description                             | Default                  |
| ------------------ | --------------------------------------- | ------------------------ |
| `ZMQ_BIND_ADDRESS` | The address for the `ROUTER` socket.    | `>tcp://127.0.0.1:5555`  |
| `ZMQ_PUB_ADDRESS`  | The address for the `PUB` socket.       | `>tcp://127.0.0.1:5556`  |

*Note: The `>` prefix in the default address is a convention used by `go-zmq` to indicate a `bind` operation.*

When running with `docker run`, you can set these variables using the `-e` flag:
```bash
docker run -e ZMQ_BIND_ADDRESS=">tcp://0.0.0.0:5555" -e ZMQ_PUB_ADDRESS=">tcp://0.0.0.0:5556" -p 5555:5555 -p 5556:5556 aetherbus-tachyon
```

## Testing

The project includes both unit and integration tests.

To run all tests, use the following command:
```bash
go test -v ./...
```

The integration test (`cmd/tachyon/main_test.go`) will start a full server instance, send a message, and verify that the message is published correctly.

## Client Example (Go)

Here is a basic example of how to interact with the Tachyon server.

### Producer (Sending an Event)

This program sends a simple JSON payload to the `user.created` topic.

```go
package main

import (
    "log"
    "time"
    "github.com/pebbe/zmq4"
)

func main() {
    dealer, _ := zmq4.NewSocket(zmq4.DEALER)
    defer dealer.Close()
    dealer.Connect("tcp://localhost:5555")

    // The topic is part of the message payload in many designs,
    // or can be inferred. Here, we send a simple JSON string.
    payload := `{"message": "hello world"}`

    for {
        log.Println("Sending event...")
        // The DEALER socket adds the necessary empty frame for the ROUTER
        dealer.SendMessage(payload)
        time.Sleep(2 * time.Second)
    }
}
```

### Consumer (Subscribing to an Event)

This program subscribes to the `user.created` topic and prints any messages it receives.

```go
package main

import (
    "log"
    "github.com/pebbe/zmq4"
)

func main() {
    subscriber, _ := zmq4.NewSocket(zmq4.SUB)
    defer subscriber.Close()
    subscriber.Connect("tcp://localhost:5556")
    subscriber.SetSubscribe("user.created") // Subscribe to the topic

    log.Println("Listening for events...")
    for {
        msg, err := subscriber.RecvMessage(0)
        if err != nil {
            log.Fatalf("Error receiving message: %v", err)
        }

        // Message format: [topic, payload]
        if len(msg) == 2 {
            log.Printf("Received Event: Topic=%s, Payload=%s
", msg[0], msg[1])
        }
    }
}

```
