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
	"encoding/hex"
	"flag"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"strings"
)

type monsterMash struct {
	data []byte
	salt []byte
}

const (
	DefaultRounds = 200000
	FileBlockSize = 512
	FileMinSize = 1024 * 64
	FileMaxSize = 1024 * 1024 * 100
	PasswordCount = 10
	PasswordLength = 20
	SpaceAt = 5
)

var debug bool

func main() {
	flag.Usage = customUsage
	flag.BoolVar(&debug, "d", false, "enable debug mode")
	flag.Parse()

	fileName := flag.Arg(0)
	if fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Get data from file.
	data, err := getDataFromFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// Get password from user.
	passwd, err := getPassword()
	if err != nil {
		log.Fatal(err)
	}

	// Make passwords.
	s, err := makePasswords(data, passwd)
	if err != nil {
		log.Fatal(err)
	}

	// Print passwords.
	for i, v := range s {
		fmt.Printf("%02d: %s\n", i+1, v)
	}
}

func customUsage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage of %s: <filename>\n", os.Args[0])
	flag.PrintDefaults()
}

func getDataFromFile(filename string) (*monsterMash, error) {
	var mm monsterMash

	// Open our file.
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Validate file size.
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	size := fi.Size()
	if size < FileMinSize {
		err = fmt.Errorf("file too small: %d of %d bytes",
			size, FileMinSize)
		return nil, err
	} else if size > FileMaxSize {
		err = fmt.Errorf("file too large: %d bytes (%d max)",
			size, FileMaxSize)
		return nil, err
	}

	// Read a block of plaintext from file.
	mm.data = make([]byte, FileBlockSize)
	if _, err := io.ReadFull(file, mm.data); err != nil {
		return nil, err
	}

	// Hash the remaining file data.
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}

	// Create a 128bit salt from the hash.
	hash := hasher.Sum(nil)
	mm.salt = hash[:16]

	if debug == true {
		log.Println("salt:", hex.EncodeToString(mm.salt))
	}

	return &mm, nil
}

func getPassword() ([]byte, error) {
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

	if debug == true {
		log.Println("password:", string(s1))
	}

	return s1, nil
}

func makePasswords(mm *monsterMash, passwd []byte) ([]string, error) {
	keySize := 32
	ivSize := aes.BlockSize

	// Generate key and IV from password and salt.
	dk := pbkdf2.Key(passwd, mm.salt, DefaultRounds,
		keySize+ivSize, sha256.New)
	block, err := aes.NewCipher(dk[:keySize])
	if err != nil {
		return nil, err
	}

	if debug == true {
		log.Println("key:", hex.EncodeToString(dk[:keySize]))
		log.Println("iv:", hex.EncodeToString(dk[keySize:]))
	}

	// Encrypt data using AES256 in CBC mode.
	ciphertext := make([]byte, FileBlockSize)
	bs := cipher.NewCBCEncrypter(block, dk[keySize:])
	bs.CryptBlocks(ciphertext, mm.data)

	// Encode the ciphertext in base32.
	coder := base32.StdEncoding.WithPadding(base32.NoPadding)
	b32 := coder.EncodeToString(ciphertext)

	// Create a slice of password strings.
	var b strings.Builder
	var s []string
	for c := 0; c < PasswordCount; c++ {
		b.Reset()
		for i := 0; i < SpaceAt-1; i++ {
			from := i*SpaceAt + c*PasswordLength
			to := from + SpaceAt
			fmt.Fprintf(&b, "%s", b32[from:to])
			if i != SpaceAt-2 {
				fmt.Fprintf(&b, " ")
			}
		}
		s = append(s, b.String())
	}

	return s, nil
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