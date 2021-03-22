# evtq

Evtq is a system where connectors can register with some metadata:

```go
connectionInfo := &core.ConnectorBuilderInfo{
    Name:         "my-connector",
    EventTypes:   []string{"event-type-1", "event-type-2", "event-type-3"},
    EncodingType: infrastructure.JSONEncoding, // also supports Gob, see infrastructure/encdec.go
    Info:         &infrastructure.LocalConnectorInfo{},
}
```

Then register with the broker. See `main.go` for a set up example.

```go
b := evtq.NewBroker(...)
// In the case of local, this returns a reader and an optional system message
// In the case of tcp, this returns a nil reader and a system message alerting the receiver to listen on a port for events.
// TODO: this is not yet finished
reader, systemMessage, err := b.AddConnectorFromInfo(connectionInfo)

```
