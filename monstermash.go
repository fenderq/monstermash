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
	"bufio"
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
	"math"
	"os"
	"strings"
)

const (
	DefaultRounds  = 200000
	FileMinSize    = 1024 * 64
	FileMaxSize    = 1024 * 1024 * 100
	PasswordCount  = 10
	PasswordLength = 20
	SpaceAt        = 5
	Version        = "1.4"
)

var Debug bool
var passwordFile string

func main() {
	flag.Usage = customUsage
	flag.BoolVar(&Debug, "d", false, "enable debug mode")
	flag.StringVar(&passwordFile, "f", "", "password file")
	flag.Parse()

	if Debug == true {
		log.Printf("%s v%s\n", os.Args[0], Version)
	}

	fileName := flag.Arg(0)
	if fileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	salt, err := GetSaltFromFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	var passwd []byte
	if passwordFile != "" {
		passwd, err = readPasswordFromFile(passwordFile)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		passwd, err = GetMasterPassword()
		if err != nil {
			log.Fatal(err)
		}
	}

	s, err := MakePasswords(salt, passwd)
	if err != nil {
		log.Fatal(err)
	}

	for i, v := range s {
		fmt.Printf("%02d: %s\n", i+1, v)
	}
}

func customUsage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage of %s: <filename>\n", os.Args[0])
	flag.PrintDefaults()
}

func GetSaltFromFile(filename string) ([]byte, error) {
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

	// Hash the file.
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}

	// Create a 128bit salt from the hash.
	hash := hasher.Sum(nil)
	salt := hash[:16]

	if Debug == true {
		log.Println("salt:", hex.EncodeToString(salt))
	}

	return salt, nil
}

func GetMasterPassword() ([]byte, error) {
	s1, err := readPasswordFromUser(os.Stdin, "enter password: ")
	if err != nil {
		return nil, err
	}
	s2, err := readPasswordFromUser(os.Stdin, "confirm password: ")
	if err != nil {
		return nil, err
	}
	if bytes.Equal(s1, s2) == false {
		return nil, fmt.Errorf("password mismatch")
	}

	if Debug == true {
		log.Println("password:", string(s1))
	}

	return s1, nil
}

func MakePasswords(salt, passwd []byte) ([]string, error) {
	keySize := 32
	ivSize := aes.BlockSize

	// Generate key and IV from password and salt.
	dk := pbkdf2.Key(passwd, salt, DefaultRounds,
		keySize+ivSize, sha256.New)

	key := dk[:keySize]
	iv := dk[keySize:]

	if Debug == true {
		log.Println("key:", hex.EncodeToString(key))
		log.Println("iv:", hex.EncodeToString(iv))
	}

	// Allocate enough bytes to produce desired base32 output.
	dataSize := int(math.Ceil(PasswordCount * PasswordLength * 5.0 / 8.0))
	data := make([]byte, dataSize)

	// Use AES-256 in CTR mode as a CSPRNG.
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(data, data)

	// Encode to base32.
	coder := base32.StdEncoding.WithPadding(base32.NoPadding)
	b32 := coder.EncodeToString(data)

	if Debug == true {
		log.Println("ciphertext:", hex.EncodeToString(data))
		log.Println("base32:", b32)
	}

	// Create a slice of password strings.
	var b strings.Builder
	var s []string
	for c := 0; c < PasswordCount; c++ {
		b.Reset()
		from := c * PasswordLength
		to := from + PasswordLength
		p := b32[from:to]
		for i, v := range p {
			if _, err := b.WriteRune(v); err != nil {
				return nil, err
			}
			if i%SpaceAt == SpaceAt-1 && i != PasswordLength-1 {
				if _, err := b.WriteRune(' '); err != nil {
					return nil, err
				}
			}
		}
		s = append(s, b.String())
	}

	return s, nil
}

func readPasswordFromUser(f *os.File, prompt string) ([]byte, error) {
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

func readPasswordFromFile(fname string) ([]byte, error) {
	var line string
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line = scanner.Text()
		break
	}
	if Debug == true {
		log.Println("password:", string(line))
	}
	return []byte(line), nil
}
