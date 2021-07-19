package main

import (
	"bufio"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/csv"
	"io"
)

type streamWriterFactory func(cipher.Stream, io.Writer) io.WriteCloser

type writer struct {
	sw io.WriteCloser
	gw *gzip.Writer
	cw *csv.Writer
}

func (w *writer) Write(record []string) error {
	return w.cw.Write(record)
}

func (w *writer) Close() (err error) {
	// Ensure every CSV line was written
	w.cw.Flush()
	err = w.cw.Error()
	if err != nil {
		return
	}

	// Close gzip writer first, then close the cipher stream writer as
	// instructed by https://pkg.go.dev/compress/gzip#NewWriter
	err = w.gw.Close()
	if err != nil {
		return
	}
	err = w.sw.Close()
	if err != nil {
		return
	}

	return
}

func newWriter(output io.Writer, wrapStream streamWriterFactory, pub *rsa.PublicKey) (w *writer, err error) {
	// Create random symmetric key to be shared.
	key := make([]byte, 32)
	if _, err = io.ReadFull(rand.Reader, key); err != nil {
		return
	}

	// Encrypt symmetric key using asymmetric public key. Only the owner of the
	// respective private key will be able to decrypt the symmetric key and the
	// rest of the contents.
	encryptedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, key, nil)
	if err != nil {
		return
	}
	output.Write(encryptedKey)

	// Create the AES cipher and initialization vector (IV). The IV will be
	// written to the output file in plain because it's not secret.
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	iv := make([]byte, block.BlockSize())
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}
	output.Write(iv)

	// Create the cipher outputting to the writer.
	sw := wrapStream(cipher.NewCTR(block, iv), output)

	// Compress the content using GZip and write it to the cipher stream.
	gz := gzip.NewWriter(sw)

	// Encode everything as CSV
	cw := csv.NewWriter(gz)

	return &writer{
		sw: sw,
		gw: gz,
		cw: cw,
	}, nil
}

func classicStreamWriter(s cipher.Stream, w io.Writer) io.WriteCloser {
	return cipher.StreamWriter{S: s, W: w}
}

func newStreamWriter(s cipher.Stream, w io.Writer) io.WriteCloser {
	m := make(map[int]*[]byte)
	return &streamWriter{S: s, W: w, m: m}
}

type streamWriter struct {
	S cipher.Stream
	W io.Writer
	m map[int]*[]byte
}

func (w *streamWriter) Write(src []byte) (n int, err error) {
	if _, ok := w.m[len(src)]; !ok {
		b := make([]byte, len(src))
		w.m[len(b)] = &b
	}
	c := *w.m[len(src)]

	w.S.XORKeyStream(c, src)
	n, err = w.W.Write(c)
	if n != len(src) && err == nil { // should never happen
		err = io.ErrShortWrite
	}
	return
}

func (w *streamWriter) Close() error {
	if c, ok := w.W.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

type bufferedStreamWriter struct {
	wc io.WriteCloser
	b  *bufio.Writer
}

func buffered(f streamWriterFactory) streamWriterFactory {
	return func(s cipher.Stream, w io.Writer) io.WriteCloser {
		wc := f(s, w)
		b := bufio.NewWriter(w)
		return bufferedStreamWriter{wc, b}
	}
}

func (w bufferedStreamWriter) Write(p []byte) (n int, err error) {
	return w.b.Write(p)
}

func (w bufferedStreamWriter) Close() error {
	if err := w.b.Flush(); err != nil {
		return err
	}
	return w.wc.Close()
}
