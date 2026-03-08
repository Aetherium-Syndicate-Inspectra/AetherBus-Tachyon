# Performance Model and Benchmarking

This document defines performance phases, benchmark scope, and baseline targets.

## 1) Data Path Roadmap

### Phase 1 (current practical baseline)

- ZeroMQ transport
- JSON codec
- LZ4 compression
- ART-backed in-memory route resolution

### Phase 2 (broker optimization)

- Binary framing on transport
- Pooled buffers and lower allocations
- Lock minimization in hot path
- Reduced copy and parse overhead

### Phase 3 (extreme path)

- Shared-memory and advanced I/O integration
- Transport specialization by workload
- Cluster/federation-oriented broker mesh

## 2) Benchmark Methodology

Track throughput and latency per scenario:

- Direct mode with ACK
- Fanout mode under varying subscriber counts
- Mixed topic distributions (exact vs wildcard)
- Small/medium/large payload classes

Always report:

- p50 / p95 / p99 latency
- messages/sec throughput
- CPU and memory usage
- allocations/op (hot path)

## 3) Key Metrics

Broker counters:

- `ingress_messages_total`
- `routed_messages_total`
- `unroutable_messages_total`
- `fanout_messages_total`
- `direct_messages_total`

Latency metrics:

- `decode_latency_ns`
- `route_lookup_latency_ns`
- `publish_latency_ns`
- `ack_latency_ns`
- `end_to_end_latency_ns`

Health metrics:

- `connected_consumers`
- `inflight_messages`
- `retry_queue_depth`
- `dlq_depth`

## 4) Initial Target Template

The project should define explicit target SLOs over time, for example:

- p99 direct-route latency target under nominal load
- minimum sustained throughput target
- maximum acceptable unroutable rate

Set concrete numbers once baseline benchmark results are available in CI/perf environments.

## 5) Compression Toggle for Benchmarking

The default runtime wiring (`internal/app.NewRuntime`) uses LZ4 compression.
For benchmark scenarios that require `--compress=false`, use benchmark runtime wiring:

- `cmd/tachyon/main.go` exposes a `--compress` flag (default `true`).
- `internal/app.NewBenchmarkRuntime(...)` selects LZ4 when enabled and a no-op compressor when disabled.

This allows fair comparison between compressed and uncompressed broker paths without introducing cgo/FFI complexity.

