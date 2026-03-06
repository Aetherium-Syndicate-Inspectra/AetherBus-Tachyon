# AetherBus Protocol Specification v1 (Draft)

This document defines the canonical wire-level and envelope-level contract for AetherBus-Tachyon.

## 1) Envelope Schema

```json
{
  "spec_version": "abtp/1",
  "message_id": "msg_01HV...",
  "correlation_id": "corr_01HV...",
  "producer_id": "svc.order-api",
  "topic": "orders.created",
  "headers": {
    "content_type": "application/json",
    "codec": "json",
    "compression": "lz4",
    "priority": "normal",
    "ttl_ms": 30000
  },
  "payload": {},
  "timestamp_unix_ms": 1741337000000,
  "reply_to": "orders.reply",
  "partition_key": "customer_123",
  "delivery_mode": "direct",
  "trace_id": "trace_...",
  "auth": {
    "type": "bearer"
  }
}
```

### Required fields

- `spec_version`
- `message_id`
- `producer_id`
- `topic`
- `payload`
- `timestamp_unix_ms`

### Optional fields

- `correlation_id`
- `headers`
- `reply_to`
- `partition_key`
- `delivery_mode`
- `trace_id`
- `auth`

## 2) Delivery Modes

- `direct`: point-to-point delivery (command/request-reply style)
- `fanout`: one-to-many publish/subscribe delivery
- `bridge`: broker-to-broker or broker-to-external transport

## 3) Registration Protocol

### Consumer register request

```json
{
  "type": "consumer.register",
  "consumer_id": "worker.invoice.1",
  "mode": "direct",
  "subscriptions": ["orders.created", "payments.*"],
  "capabilities": {
    "max_inflight": 1024,
    "supports_ack": true,
    "supports_compression": ["lz4"],
    "supports_codec": ["json"]
  }
}
```

### Consumer register acknowledgement

```json
{
  "type": "consumer.registered",
  "consumer_id": "worker.invoice.1",
  "session_id": "sess_01HV...",
  "status": "ok"
}
```

## 4) ACK/NACK Protocol

### ACK

```json
{
  "type": "ack",
  "message_id": "msg_01HV...",
  "consumer_id": "worker.invoice.1",
  "status": "processed",
  "processed_at": 1741337000500
}
```

### NACK

```json
{
  "type": "nack",
  "message_id": "msg_01HV...",
  "consumer_id": "worker.invoice.1",
  "status": "retryable_error",
  "reason": "temporary db timeout"
}
```

Retry policy should be configurable with fixed delay or exponential backoff, max retry count, and dead-letter routing.

## 5) Topic Grammar

Canonical format:

- Dot-separated segments: `domain.entity.action`
- Examples: `orders.created`, `payments.authorized`, `system.node.heartbeat`

Wildcard support:

- `*` matches one segment
- `>` matches remaining segments and is allowed only at the end of a pattern

Examples:

- `orders.*` matches `orders.created`
- `system.*.heartbeat` matches `system.node.heartbeat`
- `agents.>` matches all topics under `agents`

## 6) Wire Framing (ZeroMQ)

Recommended frame layout for latency-oriented transport:

- Frame 0: socket/routing identity
- Frame 1: protocol version
- Frame 2: flags
- Frame 3: topic
- Frame 4: headers
- Frame 5: payload

This framing allows the broker to inspect topic and control metadata without fully decoding payload in all paths.

## 7) Baseline Guarantees

- PUB/SUB path: at-most-once, best-effort fanout
- ROUTER/DEALER path with ACK: at-least-once transport behavior
- Application-level effect: effectively-once when consumers are idempotent
