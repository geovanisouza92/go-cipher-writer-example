package main

import (
	"bytes"
	"container/ring"
	"crypto/rsa"
	"crypto/x509"
	"encoding/csv"
	"encoding/pem"
	"io"
	"io/ioutil"
	"testing"
)

func BenchmarkWriterClassicNoBuffer(b *testing.B) {
	b.Run("1k", benchmark(classicStreamWriter, 1000))
	b.Run("10k", benchmark(classicStreamWriter, 10000))
	b.Run("100k", benchmark(classicStreamWriter, 100000))
	b.Run("1m", benchmark(classicStreamWriter, 1000000))
	b.Run("10m", benchmark(classicStreamWriter, 10000000))
}

func BenchmarkWriterClassicBufferred(b *testing.B) {
	b.Run("1k", benchmark(buffered(classicStreamWriter), 1000))
	b.Run("10k", benchmark(buffered(classicStreamWriter), 10000))
	b.Run("100k", benchmark(buffered(classicStreamWriter), 100000))
	b.Run("1m", benchmark(buffered(classicStreamWriter), 1000000))
	b.Run("10m", benchmark(buffered(classicStreamWriter), 10000000))
}

func BenchmarkWriterNewNoBuffer(b *testing.B) {
	b.Run("1k", benchmark(newStreamWriter, 1000))
	b.Run("10k", benchmark(newStreamWriter, 10000))
	b.Run("100k", benchmark(newStreamWriter, 100000))
	b.Run("1m", benchmark(newStreamWriter, 1000000))
	b.Run("10m", benchmark(newStreamWriter, 10000000))
}

func BenchmarkWriterNewBufferred(b *testing.B) {
	b.Run("1k", benchmark(buffered(newStreamWriter), 1000))
	b.Run("10k", benchmark(buffered(newStreamWriter), 10000))
	b.Run("100k", benchmark(buffered(newStreamWriter), 100000))
	b.Run("1m", benchmark(buffered(newStreamWriter), 1000000))
	b.Run("10m", benchmark(buffered(newStreamWriter), 10000000))
}

func benchmark(f streamWriterFactory, n int) func(*testing.B) {
	key := readTestKey()
	lines := readLines()

	return func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			w, _ := newWriter(io.Discard, f, key)
			for i := 0; i < n; i++ {
				w.Write(lines.Value.([]string))
				lines = lines.Next()
			}
			w.Close()
		}
	}
}

func readTestKey() *rsa.PublicKey {
	b, _ := ioutil.ReadFile("private-key.pem")
	der, _ := pem.Decode(b)
	priv, _ := x509.ParsePKCS1PrivateKey(der.Bytes)
	return &priv.PublicKey
}

func readLines() *ring.Ring {
	b, _ := ioutil.ReadFile("sample.csv")
	lines, _ := csv.NewReader(bytes.NewReader(b)).ReadAll()
	r := ring.New(len(lines))
	for _, line := range lines {
		r.Value = line
		r = r.Next()
	}
	return r
}
