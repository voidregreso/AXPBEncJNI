package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/blake2b"
)

const (
	defaultArgonTime   = 4
	defaultArgonMemory = 4e6
	saltSize           = 32
	keyLen             = 32
	macLen             = 32
)

type fileHeader struct {
	Salt        [saltSize]byte
	ArgonTime   uint32
	ArgonMemory uint32
	ArgonLanes  uint8
	Tag         [64]byte
}

var errBadMAC = errors.New("authentication failed")

func generateKey(passphrase []byte) ([]byte, fileHeader, error) {
	var salt [saltSize]byte
	if _, err := rand.Read(salt[:]); err != nil {
		return nil, fileHeader{}, err
	}
	header := fileHeader{
		Salt:        salt,
		ArgonTime:   defaultArgonTime,
		ArgonMemory: defaultArgonMemory,
		ArgonLanes:  uint8(runtime.NumCPU() * 2),
	}
	return argon2.IDKey(passphrase, header.Salt[:], header.ArgonTime, header.ArgonMemory, header.ArgonLanes, keyLen+macLen), header, nil
}

func authenticateHeader(passphrase []byte, input *os.File) (fileHeader, error) {
	header := fileHeader{}
	if _, err := input.Seek(0, 0); err != nil {
		return header, err
	}
	if err := binary.Read(input, binary.LittleEndian, &header); err != nil {
		return header, err
	}
	return header, nil
}

func verifyMAC(macKey []byte, input *os.File, header fileHeader) error {
	var mac [64]byte
	hash, err := blake2b.New512(macKey)
	if err != nil {
		return err
	}
	if _, err := io.Copy(hash, input); err != nil {
		return err
	}
	copy(mac[:], hash.Sum(nil))
	if subtle.ConstantTimeCompare(mac[:], header.Tag[:]) != 1 {
		return errBadMAC
	}
	return nil
}

func openOutputFile(finalOutput string) (*os.File, error) {
	output, err := os.Create(finalOutput + ".temp")
	if err != nil {
		return nil, err
	}
	return output, nil
}

func closeAndRenameOutput(output *os.File, finalOutput string) error {
	if err := output.Sync(); err != nil {
		return err
	}
	if err := output.Close(); err != nil {
		return err
	}
	if err := os.Rename(output.Name(), finalOutput); err != nil {
		return err
	}
	return nil
}

func decryptFile(passphrase []byte, input *os.File, finalOutput string) error {
	header, err := authenticateHeader(passphrase, input)
	if err != nil {
		return err
	}

	ciphertextOffset, err := input.Seek(0, 1)
	if err != nil {
		return err
	}

	skb := argon2.IDKey(passphrase, header.Salt[:], header.ArgonTime, header.ArgonMemory, header.ArgonLanes, keyLen+macLen)
	var sk [keyLen]byte
	var macKey [keyLen]byte
	copy(sk[:], skb[:keyLen])
	copy(macKey[:], skb[keyLen:])

	if err := verifyMAC(macKey[:], input, header); err != nil {
		return err
	}

	if _, err := input.Seek(ciphertextOffset, 0); err != nil {
		return err
	}

	output, err := openOutputFile(finalOutput)
	if err != nil {
		return err
	}
	defer os.Remove(output.Name())

	inputReader := NewReader(sk, input)
	_, err = io.Copy(output, inputReader)
	if err != nil {
		return err
	}

	return closeAndRenameOutput(output, finalOutput)
}

func encryptFile(passphrase []byte, input *os.File, finalOutput string) error {
	output, err := openOutputFile(finalOutput)
	if err != nil {
		return err
	}
	defer os.Remove(output.Name())

	skb, header, err := generateKey(passphrase)
	if err != nil {
		return fmt.Errorf("could not generate secret key")
	}

	if err := binary.Write(output, binary.LittleEndian, header); err != nil {
		return err
	}

	var sk [keyLen]byte
	var macKey [keyLen]byte
	copy(sk[:], skb[:keyLen])
	copy(macKey[:], skb[keyLen:])

	hash, err := blake2b.New512(macKey[:])
	if err != nil {
		return err
	}
	encWriter := NewWriter(sk, io.MultiWriter(hash, output))
	_, err = io.Copy(encWriter, input)
	if err != nil {
		return err
	}
	var mac [64]byte
	copy(mac[:], hash.Sum(nil))
	header.Tag = mac

	if _, err := output.Seek(0, 0); err != nil {
		return err
	}
	if err := binary.Write(output, binary.LittleEndian, header); err != nil {
		return err
	}

	return closeAndRenameOutput(output, finalOutput)
}