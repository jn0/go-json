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

Benchmark gives

    goos: linux
    goarch: amd64
    BenchmarkAll-4   	2000000000	         0.16 ns/op
    PASS
    ok  	_/home/jno/src/go-json	4.935s

on a *Intel(R) Core(TM) i5-6600 CPU @ 3.30GHz* box.

coverage: 88.1% of statements

# EOF #
