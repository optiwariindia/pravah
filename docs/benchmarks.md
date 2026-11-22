# Pravah MVP Benchmarks

## Reliable SSE Delivery

Measured on 2026-11-22 in the local Docker Compose environment.

| Metric | Target | Result |
| --- | ---: | ---: |
| Concurrent SSE clients | 100 | 100 |
| Publication rate | 50/s | 50.77/s |
| Events delivered | 5,000 | 5,000 |
| Dropped connected signals | 0 | 0 |
| Observer p95 latency | < 1 second | 944 ms |

The test connects 100 clients to the same tenant and stream, publishes 50
signals at the configured 50 signals/second rate, and records end-to-end latency
on one observer client. All 100 clients must receive every event.

Reproduce with the command in [DevelopersGuide.md](../DevelopersGuide.md).

## Container Size

The Go 1.25 scratch-based production image measured 26,581,885 bytes
(approximately 26.6 MB) after this test.

## Multi-Tab Connection Sharing

Automated browser lifecycle tests verify:

- simultaneous tabs elect one active SSE leader;
- a follower becomes leader after the leader closes;
- tenant and stream scopes do not share leaders;
- leader signals reach follower tabs.

These tests are run with `npm test -- --run` in `pravah-sdk/`.
