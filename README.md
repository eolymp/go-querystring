# QueryString

A Go library for unmarshaling URL query strings into structs, to help map REST to RPC (gRPC).

It uses a simple, opinionated mapping based on key-value pairs for basic types (`int`, `string`, `bool`), and resorts to JSON
encoding for more complex types.

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
?Offset=1&SearchQuery=jose&Filters={"field":"nickname","equals":"jose"}`
```

## License

MIT License