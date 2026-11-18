Implement the Composite pattern in Go (window example):

- Interface `Item` with `Price() float64` and `Display(indent int)`.
- `Glass`, `Frame` and `Window` all implement `Item`.
- `Glass` and `Frame` are leaves with a fixed price each.
- `Window` is the composite: it holds `[]Item` (at least one, including nested `Window`s). `Price()` = sum of its children; `Display()` prints the tree with indentation.
- `main.go` builds a nested example window and prints the tree and total price.
- One package, `go.mod`, a test for `Price()`, short README. Idiomatic Go, no dependencies.