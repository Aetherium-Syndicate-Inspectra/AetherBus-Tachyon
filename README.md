# 🌌 AetherBus Tachyon

**The Ultra-Fast Backbone for Decentralized Intelligence and Hyperscale Data Transmission.**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Tech: Go](https://img.shields.io/badge/Tech-Go-blue)](https://go.dev/)
[![GitHub repo stars](https://img.shields.io/github/stars/Aetherium-Syndicate-Inspectra/AetherBus-Tachyon?style=social)](https://github.com/Aetherium-Syndicate-Inspectra/AetherBus-Tachyon)

## 🚀 Vision

AetherBus is a decentralized, high-speed data transmission node system designed to overcome the limitations of traditional networking. Inspired by space-based laser communication technology (the Space Data Highway), it's built to support hyperscale AI processing and data pipelines.

## 🏛️ Architectural Overview

AetherBus Tachyon is built using a clean, hexagonal architecture to ensure a clear separation of concerns, making the system modular, testable, and maintainable.

*   `cmd/aetherbus-node/`: The main application entry point. It initializes all components and starts the node.
*   `internal/`: Contains the core business logic, separated into layers:
    *   `delivery/`: Adapters for incoming connections (e.g., ZeroMQ). This is the outer-most layer.
    *   `usecase/`: Orchestrates the flow of data and implements the core application logic (e.g., `EventRouter`).
    *   `repository/`: Interfaces and implementations for data persistence (e.g., `ArtRouteStore` for the in-memory routing table).
    *   `domain/`: Core data structures and business rules of the application (e.g., `Message`, `Route`).
*   `pkg/`: (Currently unused) Intended for shared libraries that can be used by other projects.
*   `go.mod` & `go.sum`: Manages project dependencies.

## 🚦 Quick Start

To get a local node running, follow these steps. You need to have Go (version 1.22 or later) installed.

```bash
# 1. Clone the repository
git clone https://github.com/Aetherium-Syndicate-Inspectra/AetherBus-Tachyon.git
cd AetherBus-Tachyon

# 2. Tidy dependencies
# This command ensures your project has all the necessary dependencies.
go mod tidy

# 3. Build the application
# This compiles the source code into a single executable binary.
go build ./cmd/aetherbus-node

# 4. Run the node
# This will start the AetherBus Tachyon node.
./aetherbus-node
```

## 🤝 Contributing

We welcome contributions! Please feel free to open an issue or submit a pull request.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.
