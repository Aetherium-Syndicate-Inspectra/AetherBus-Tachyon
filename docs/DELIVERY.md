# Delivery Architecture

This document defines transport responsibilities and runtime behavior for message delivery.

## 1) ZeroMQ Role Mapping

### ROUTER

Primary ingress and control channel:

- Producer ingress
- Command/request-reply traffic
- Consumer registration
- ACK/NACK handling

### DEALER (worker side)

- Worker/consumer async handling
- Broker internal parallel processing patterns

### PUB

Low-overhead egress for fanout events:

- Broadcast event stream
- Telemetry-style subscribers

## 2) Transport Rule of Thumb

- Use ROUTER/DEALER for flows that require ACK and retries
- Use PUB/SUB for speed-oriented best-effort fanout

## 3) Delivery Flow

1. Receive ZeroMQ frames
2. Decode metadata and envelope
3. Decompress payload
4. Validate envelope
5. Resolve routes via RouteStore (ART)
6. Dispatch via direct/fanout/bridge mode
7. Track ACK for direct mode

## 4) Backpressure Controls

- Per-consumer `max_inflight`
- Per-topic queue limit
- Global ingress limit
- Slow-consumer detection and degradation policy

Suggested slow-consumer behavior:

- Mark consumer as degraded
- Stop assigning additional direct load
- Keep heartbeat/session alive

## 5) Retry and Dead Letter

Retry should support:

- Fixed delay
- Exponential backoff
- Max retry count

DLQ topic convention:

- `_dlq.<original_topic>`

DLQ metadata fields:

- `dead_letter_reason`
- `dead_letter_count`
- `original_topic`

## 6) Failure Handling Baseline

- Malformed payload/decode error: reject immediately
- No route found: emit system unroutable event
- Disconnect with inflight message: retry or DLQ
- Stale ACK: ignore and log
