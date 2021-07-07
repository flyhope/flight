# flight
golang single-flight and multi-flight

## QA

### Q: Why not single-flight

If a resource sporadic slow, single flight is too danger. all request wait a slow single flight.

## see also

singleflight：https://pkg.go.dev/golang.org/x/sync/singleflight
