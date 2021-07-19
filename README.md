# go-cipher-writer-example

This is an example about a very specific issue when using Go's standard library [`crypto/cipher.StreamWriter`](https://pkg.go.dev/crypto/cipher#StreamWriter), previously reported in [a similar case involving pgp encryption](https://github.com/golang/go/issues/26578).

For some uses cases like wrapping it with [`compress/gzip.Writer`](https://pkg.go.dev/compress/gzip#Writer), the [`StreamWriter.Write`](https://cs.opensource.google/go/go/+/go1.16.6:src/crypto/cipher/io.go;l=37) (`io.Writer` interface) allocates a new temporary buffer used for encryption, but due to the behavior of `gzip.Writer`, most of them are quite small.

If we use both to write large files (in this example, writing up to 10 million lines using `encoding/csv.CSVWriter`), we face a large memory usage, pressing the GC afterwards.

As demonstrated in this example, there's a way to avoid this problem by reusing same-size buffers (`streamWriter.Write` method in [`writer.go`](./writer.go)).

To confirm the findings, you can run `make bench`, but make sure you can [benchstat](https://pkg.go.dev/golang.org/x/perf/cmd/benchstat) installed on GOPATH.

The results of the benchmark with the change are:

```txt
name            old time/op    new time/op    delta
Writer/1k-4      2.33ms ± 5%    2.28ms ± 3%     ~     (p=0.151 n=5+5)
Writer/10k-4     17.0ms ± 2%    16.3ms ± 1%   -3.86%  (p=0.008 n=5+5)
Writer/100k-4     160ms ± 3%     162ms ±11%     ~     (p=0.548 n=5+5)
Writer/1m-4       1.68s ±10%     1.56s ± 5%   -7.38%  (p=0.032 n=5+5)
Writer/10m-4      16.3s ± 5%     15.9s ± 3%     ~     (p=0.151 n=5+5)

name            old alloc/op   new alloc/op   delta
Writer/1k-4       843kB ± 0%     836kB ± 0%   -0.93%  (p=0.008 n=5+5)
Writer/10k-4      914kB ± 0%     835kB ± 0%   -8.59%  (p=0.008 n=5+5)
Writer/100k-4    1.62MB ± 0%    0.84MB ± 0%  -48.45%  (p=0.008 n=5+5)
Writer/1m-4      8.69MB ± 0%    0.84MB ± 0%  -90.36%  (p=0.008 n=5+5)
Writer/10m-4     79.3MB ± 0%     0.8MB ± 0%  -98.95%  (p=0.016 n=5+4)

name            old allocs/op  new allocs/op  delta
Writer/1k-4        88.0 ± 0%      61.0 ± 0%  -30.68%  (p=0.008 n=5+5)
Writer/10k-4        382 ± 0%        61 ± 0%  -84.03%  (p=0.008 n=5+5)
Writer/100k-4     3.33k ± 0%     0.06k ± 1%  -98.15%  (p=0.008 n=5+5)
Writer/1m-4       32.8k ± 0%      0.1k ± 4%  -99.79%  (p=0.008 n=5+5)
Writer/10m-4       327k ± 0%        0k ± 0%  -99.98%  (p=0.016 n=5+4)
```
