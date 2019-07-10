//
// Copyright (c) 2019 Steven Roberts <sroberts@fenderq.com>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
//

package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base32"
	"flag"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
)

const DefaultRounds int = 262144
const FileBufferSize int = 512
const PasswordCount int = 10
const PasswordLength int = 20
const SplitAt int = 5

func main() {
	var verbose bool
	flag.BoolVar(&verbose, "v", false, "enable verbose mode")
	flag.Parse()

	filename := flag.Arg(0)
	if filename == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Get data from file.
	data, err := dataFromFile(filename, FileBufferSize)
	if err != nil {
		log.Fatal(err)
	}

	// Get password from user.
	passwd, err := doPassword()
	if err != nil {
		log.Fatal(err)
	}

	makePaperPasswords(data, passwd)
}

func makePaperPasswords(data []byte, passwd []byte) ([]string, error) {
	hash := sha256.Sum256(data)
	salt := hash[:8]
	ks := 32
	ivs := aes.BlockSize
	dk := pbkdf2.Key(passwd, salt, DefaultRounds, ks+ivs, sha256.New)
	block, err := aes.NewCipher(dk[:ks])
	if err != nil {
		return nil, err
	}
	bs := cipher.NewCBCEncrypter(block, dk[ks:])
	bs.CryptBlocks(data, data)
	coder := base32.StdEncoding.WithPadding(base32.NoPadding)
	s := coder.EncodeToString(data)
	for c := 0; c < PasswordCount; c++ {
		for i := 0; i < SplitAt-1; i++ {
			from := i*SplitAt + c*PasswordLength
			to := from + SplitAt
			fmt.Printf("%s", s[from:to])
			if i == SplitAt-2 {
				fmt.Printf("\n")
			} else {
				fmt.Printf(" ")
			}
		}
	}
	return nil, nil
}

func doPassword() ([]byte, error) {
	s1, err := readPassword(os.Stdin, "enter password: ")
	if err != nil {
		return nil, err
	}
	s2, err := readPassword(os.Stdin, "confirm password: ")
	if err != nil {
		return nil, err
	}
	if bytes.Equal(s1, s2) == false {
		return nil, fmt.Errorf("password mismatch")
	}
	return s1, nil
}

func readPassword(f *os.File, prompt string) ([]byte, error) {
	fd := int(f.Fd())
	if terminal.IsTerminal(fd) == false {
		return nil, fmt.Errorf("invalid terminal")
	}
	fmt.Printf("%s", prompt)
	data, err := terminal.ReadPassword(fd)
	fmt.Printf("\n")
	if err != nil {
		return nil, err
	}
	return data, nil
}

func dataFromFile(filename string, size int) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data := make([]byte, size)
	n, err := io.ReadFull(file, data)
	if err != nil {
		if n != size {
			err = fmt.Errorf("%s: %d of %d bytes", err, n, size)
		}
		return nil, err
	}
	return data, nil
}
