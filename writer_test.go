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

func BenchmarkWriterClassic(b *testing.B) {
	key := readTestKey()
	lines := readLines()

	benchmark := func(n int) func(*testing.B) {
		return func(bb *testing.B) {
			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				w, _ := newWriter(io.Discard, classicStreamWriter, key)
				for i := 0; i < n; i++ {
					w.Write(lines.Value.([]string))
					lines = lines.Next()
				}
				w.Close()
			}
		}
	}

	b.Run("1k", benchmark(1000))
	b.Run("10k", benchmark(10000))
	b.Run("100k", benchmark(100000))
	b.Run("1m", benchmark(1000000))
	b.Run("10m", benchmark(10000000))
}

func BenchmarkWriterNew(b *testing.B) {
	key := readTestKey()
	lines := readLines()

	benchmark := func(n int) func(*testing.B) {
		return func(bb *testing.B) {
			bb.ResetTimer()

			for i := 0; i < bb.N; i++ {
				w, _ := newWriter(io.Discard, newStreamWriter, key)
				for i := 0; i < n; i++ {
					w.Write(lines.Value.([]string))
					lines = lines.Next()
				}
				w.Close()
			}
		}
	}

	b.Run("1k", benchmark(1000))
	b.Run("10k", benchmark(10000))
	b.Run("100k", benchmark(100000))
	b.Run("1m", benchmark(1000000))
	b.Run("10m", benchmark(10000000))
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
