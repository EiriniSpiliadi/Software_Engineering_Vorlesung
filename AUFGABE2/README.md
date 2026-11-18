# Composite pattern in Go: window example

`Item` is the common interface (`Price`, `Display`). `Glass` and `Frame` are
leaves with fixed prices; `Window` is the composite holding one or more
`Item`s, including nested `Window`s. Its `Price()` sums its children and
`Display()` prints the tree with indentation.

`NewWindow(first, rest...)` guarantees at least one child.

## Run

    go run .
    go test ./...

No dependencies.
