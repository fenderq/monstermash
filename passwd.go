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
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
)

func GetPasswordFromFile(fname string) ([]byte, error) {
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

func GetPasswordFromUser() ([]byte, error) {
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

	if Debug == true {
		log.Println("password:", string(s1))
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
