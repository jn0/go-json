# go-json

My own approach to JSON in Go.

See [example](json_test.go#L21) in [`json_test.go`](json_test.go) for an idea.

There is a full (I hope) parser, but the initial idea was to *generate* JSON
in some uniform way from different sources (yepp, sorta system monitor agent).

The point is to have special types for JSON entities that boil down to
[`JsonValue`](json_values.go#L19) type:

  - `JsonInt` (no, Ints aren't "sorta floats")
  - `JsonFloat`
  - `JsonBool`
  - `JsonString` (all escapes, including \uXXXX are ok on input)
  - `JsonArray`
  - `JsonObject`
  - no separate type for `null`, anyone can be.

Any `JsonValue` has `.Json()` method to get a `string` representation of that
value suitable to send over, say, HTTP POST method.

The `JsonArray` can be `.Append()`ed and `JsonObject` has `.Insert()` method.

Any other `JsonValue` considered *immutable* (one can *replace* it with `.Set()`
method). The `.Set()` method accepts a "compatible" value or a `string`. The
"compatibility" means that you can use either `float32` or `float64` as value
for `JsonFloat` and so on. The `string`s are `.Parse()`d, while not `.Set()` into a `JsonString`.

The `.Value()` returns "unJSONed" version of that `JsonValue` (type cast still needed).

The `.Equal()` compares the two `JsonValue`s to be equal and `.IsNull()` compares to zero-value.

[Benchmark](json_test.go#L14) gives

    goos: linux
    goarch: amd64
    BenchmarkAll-4   	2000000000	         0.17 ns/op
    PASS
    ok  	_/home/jno/src/go-json	4.935s

on a *Intel(R) Core(TM) i5-6600 CPU @ 3.30GHz* box.

coverage: 100% of statements

# EOF #
