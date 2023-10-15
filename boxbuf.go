package main

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
)

//
// maxChunkSize determines the amount of data written to an EncWriter before a
// new chunk is written.

const maxChunkSize = 16384 // 16kb

// EncWriter is an io.Writer that can be used to encrypt data with a secret key.
// EncWriter uses golang.org/x/crypto/nacl/secretbox to perform symmetric
// encryption.
type EncWriter struct {
	out       io.Writer
	buf       []byte
	secretKey [32]byte
}

// DecReader is an io.Reader that can be used to decrypt data using a secret
// key. DecWriter uses golang.org/x/crypto/nacl/secretbox to perform symmetric
// decryption.
type DecReader struct {
	in    io.Reader
	buf   []byte
	index int

	secretKey [32]byte
}

// NewWriter creates a new EncWriter using the provided secretKey to encrypt
// data as needed to out.
func NewWriter(secretKey [32]byte, out io.Writer) *EncWriter {
	return &EncWriter{secretKey: secretKey, out: out}
}

// NewReader creates a new DecReader using secretKey to decrypt the data as
// needed from in.
func NewReader(secretKey [32]byte, in io.Reader) *DecReader {
	return &DecReader{secretKey: secretKey, in: in}
}

// Write writes the entirety of p to the underlying io.Writer, encrypting the
// data with the public key and chunking as needed.
func (w *EncWriter) Write(p []byte) (int, error) {
	totalWritten := 0
	for len(p) > 0 {
		chunkSize := maxChunkSize
		if len(p) < maxChunkSize {
			chunkSize = len(p)
		}
		err := w.writeChunk(p[:chunkSize])
		if err != nil {
			return totalWritten, err
		}
		totalWritten += chunkSize
		p = p[chunkSize:]
	}
	return totalWritten, nil
}

// writeChunk writes a chunk using EncWriter's buf and resets the buffer.
func (w *EncWriter) writeChunk(p []byte) error {
	var nonce [24]byte
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return err
	}
	aead, err := chacha20poly1305.NewX(w.secretKey[:])
	if err != nil {
		return err
	}
	encryptedData := aead.Seal(nil, nonce[:], p, nil)
	_, err = w.out.Write(nonce[:])
	if err != nil {
		return err
	}
	err = binary.Write(w.out, binary.LittleEndian, uint64(len(encryptedData)))
	if err != nil {
		return err
	}
	_, err = w.out.Write(encryptedData)
	return err
}

// Read reads from the underlying io.Reader, decrypting bytes as needed, until
// len(p) byte have been read or the underlying stream is exhausted.
func (b *DecReader) Read(p []byte) (int, error) {
	read := 0
	for i := range p {
		if b.index == 0 {
			err := b.nextChunk()
			if err != nil {
				return read, err
			}
		}
		p[i] = b.buf[b.index]
		b.index++
		read++
		if b.index >= len(b.buf) {
			b.index = 0
		}
	}
	return read, nil
}

// nextChunk reads the next chunk into DecReader's buf.
func (b *DecReader) nextChunk() error {
	var nonce [24]byte
	_, err := io.ReadFull(b.in, nonce[:])
	if err != nil {
		return err
	}
	var chunkSize uint64
	err = binary.Read(b.in, binary.LittleEndian, &chunkSize)
	if err != nil {
		return err
	}
	if chunkSize > maxChunkSize+16 {
		return errors.New("chunk too large")
	}
	chunkData := make([]byte, chunkSize)
	_, err = io.ReadFull(b.in, chunkData)
	if err != nil {
		return err
	}
	aead, err := chacha20poly1305.NewX(b.secretKey[:])
	if err != nil {
		return err
	}
	b.buf, err = aead.Open(nil, nonce[:], chunkData, nil)
	return err
}