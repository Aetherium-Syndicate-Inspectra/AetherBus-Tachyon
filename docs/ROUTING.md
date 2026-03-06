# Routing Design (ART-backed)

This document defines route matching semantics and precedence for AetherBus-Tachyon.

## 1) Route Model

Suggested internal route fields:

```go
type Route struct {
    Pattern       string
    RouteType     RouteType
    DestinationID string
    Priority      int
    Metadata      map[string]string
}
```

## 2) Why Adaptive Radix Tree (ART)

ART is used as the routing core because it offers:

- Fast prefix and key lookup
- Efficient memory usage for sparse route sets
- Strong performance for large routing tables

## 3) Topic Pattern Semantics

- Exact route: full topic string match
- Single wildcard route: `*` for one segment
- Remainder wildcard route: `>` for the trailing remainder of a topic

`>` should be allowed only as the trailing suffix of a pattern to keep matching deterministic.

## 4) Route Resolution Order

When multiple patterns match, resolve in this order:

1. Exact match
2. Single-segment wildcard (`*`)
3. Remainder wildcard (`>`)
4. Route priority (`Priority`, higher first)
5. Stable insertion order fallback

## 5) Route Lifecycle

- Register route
- Validate topic pattern grammar
- Insert into ART index
- Resolve by topic on message ingress
- Remove/update route on unregistration or config change

## 6) Operational Metrics

Recommended routing-specific metrics:

- `route_count`
- `wildcard_route_count`
- `exact_match_hit_rate`
- `wildcard_match_hit_rate`
- `route_lookup_latency_ns`
- `longest_lookup_ns`
