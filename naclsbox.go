package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/nacl/secretbox"
)

var (
	encrypt = flag.Bool("e", false, "Encrypt files")
	decrypt = flag.Bool("d", false, "Decrypt files")
	key     = flag.String("k", "", "Decryption key ")
)

func decodeKey(k string) (*[32]byte, error) {
	dk, err := hex.DecodeString(k)
	if err != nil {
		return nil, err
	}
	if len(dk) != 32 {
		return nil, errors.New("invalid key length")
	}

	var sizedKey [32]byte
	copy(sizedKey[:], dk)
	return &sizedKey, nil
}

func newKey() (*[32]byte, error) {
	var k [32]byte

	_, err := io.ReadFull(rand.Reader, k[:])
	if err != nil {
		return nil, err
	}

	return &k, nil
}

func unpackAndDecrypt(key *[32]byte, packed []byte) ([]byte, error) {
	if len(packed) < 24 {
		return nil, errors.New("not an sbox file")
	}

	var nonce [24]byte
	copy(nonce[:], packed[0:24])

	decrypted, ok := secretbox.Open(nil, packed[24:], &nonce, key)
	if !ok {
		return nil, fmt.Errorf("could not decrypt")
	}
	return decrypted, nil
}

func encryptAndPack(data []byte, key *[32]byte) ([]byte, error) {
	var nonce [24]byte

	// random nonce
	_, err := io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return nil, err
	}

	// encrypt data
	blob := secretbox.Seal(nil, data, &nonce, key)

	// pack all the things
	packed := make([]byte, len(nonce)+len(blob))
	copy(packed[0:], nonce[:])
	copy(packed[24:], blob)

	return packed, nil
}

func _main() error {
	flag.Parse()

	// Must provide -e OR -d not both.
	if !*encrypt && !*decrypt {
		flag.PrintDefaults()
		return errors.New("must provide -e or -d\n")

	}
	if *encrypt && *decrypt {
		flag.PrintDefaults()
		return errors.New("must provide only -e or -d\n")
	}

	// -d requires a key
	if *decrypt && *key == "" {
		flag.PrintDefaults()
		return errors.New("must provide only -k when decrypting\n")
	}

	if *encrypt {
		k, err := newKey()
		if err != nil {
			return err
		}
		fmt.Printf("encryption key: %x\n", *k)

		for _, filename := range flag.Args() {
			f, err := os.Open(filename)
			if err != nil {
				return err
			}

			data, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}
			f.Close()

			blob, err := encryptAndPack(data, k)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(filename+".sbox", blob, 0664)
			if err != nil {
				return err
			}
		}
	}

	if *decrypt {
		k, err := decodeKey(*key)
		if err != nil {
			return err
		}

		for _, filename := range flag.Args() {
			f, err := os.Open(filename)
			if err != nil {
				return err
			}

			blob, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}
			f.Close()

			data, err := unpackAndDecrypt(k, blob)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(filename+".decrypted", data, 0664)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func main() {
	err := _main()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
