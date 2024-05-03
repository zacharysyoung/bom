# BOM Sniffer

This package provides a simple, well-assuming io.Reader, that sniffs for one of the five UTF BOMs.

"Well-assuming" in that it has only been tested with in-memory Readers, and a few times on disk—so, maybe insert "toy" into the above description as well.

Get a feel for how to use this by looking at [example_test.go](./example_test.go), or [cmd/main.go](./cmd/main.go).
