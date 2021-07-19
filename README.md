# go-cipher-writer-example

This is an example about a very specific issue when using Go's standard library [`crypto/cipher.StreamWriter`](https://pkg.go.dev/crypto/cipher#StreamWriter), previously reported in [a similar case involving pgp encryption](https://github.com/golang/go/issues/26578).

For some uses cases like wrapping it with [`compress/gzip.Writer`](https://pkg.go.dev/compress/gzip#Writer), the [`StreamWriter.Write`](https://cs.opensource.google/go/go/+/go1.16.6:src/crypto/cipher/io.go;l=37) (`io.Writer` interface) allocates a new temporary buffer used for encryption, but due to the behavior of `gzip.Writer`, most of them are quite small.

If we use both to write large files (in this example, writing up to 10 million lines using `encoding/csv.CSVWriter`), we face a large memory usage, pressing the GC afterwards.

As demonstrated in this example, there's a way to avoid this problem by reusing same-size buffers (`streamWriter.Write` method in [`writer.go`](./writer.go)).

To confirm the findings, you can run `make bench`, but make sure you can [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) installed on GOPATH.
