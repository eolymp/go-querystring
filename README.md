# QueryString

A Go library for unmarshaling URL query strings into structs. It provides an opinionated and simple approach designed primarily for mapping Eolymp REST API requests to gRPC protobuf messages.

Simple types are mapped directly to query parameter values, while complex types are JSON-encoded as single parameter values. Slices are mapped to repeated query parameters with the same key.

For example, a simple list request might be represented using the following structure:

```go
type ListUsersRequest struct {
	Offset       int
	SearchQuery  string
	Filters      []FilterExpression
}
```

This library would map this request to the query:

```
?Offset=1&SearchQuery=jose&Filters={"field":"nickname","equals":"jose"}
```


## License

MIT License