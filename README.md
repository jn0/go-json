# go-json

My own approach to JSON in Go.

The point is to have special types for JSON entities that boil down to
[`JsonValue`](json_values.go#L19) type:

  - `JsonInt`
  - `JsonFloat`
  - `JsonBool`
  - `JsonString`
  - `JsonArray`
  - `JsonObject`

coverage: 88.1% of statements

# EOF #
