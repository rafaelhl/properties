// Copyright 2013 Frank Schroeder. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goproperties

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Load reads a buffer into a Properties struct.
func Load(buf []byte, enc Encoding) (*Properties, error) {
	return loadBuf(buf, enc)
}

// LoadFile reads a file into a Properties struct.
func LoadFile(filename string, enc Encoding) (*Properties, error) {
	return loadFiles([]string{filename}, enc, false)
}

// LoadFiles reads multiple files in the given order into
// a Properties struct. If 'ignoreMissing' is true then
// non-existent files will not be reported as error.
func LoadFiles(filenames []string, enc Encoding, ignoreMissing bool) (*Properties, error) {
	return loadFiles(filenames, enc, ignoreMissing)
}

// MustLoadFile reads a file into a Properties struct and
// panics on error.
func MustLoadFile(filename string, enc Encoding) *Properties {
	return mustLoadFiles([]string{filename}, enc, false)
}

// MustLoadFiles reads multiple files in the given order into
// a Properties struct and panics on error. If 'ignoreMissing'
// is true then non-existent files will not be reported as error.
func MustLoadFiles(filenames []string, enc Encoding, ignoreMissing bool) *Properties {
	return mustLoadFiles(filenames, enc, ignoreMissing)
}

// ----------------------------------------------------------------------------

type Encoding uint

const (
	UTF8 Encoding = 1 << iota
	ISO_8859_1
)

func loadBuf(buf []byte, enc Encoding) (*Properties, error) {
	return parse(convert(buf, enc))
}

func loadFiles(filenames []string, enc Encoding, ignoreMissing bool) (*Properties, error) {
	buff := make([]byte, 0, 4096)

	for _, filename := range filenames {
		buf, err := ioutil.ReadFile(filename)
		if err != nil {
			if ignoreMissing && os.IsNotExist(err) {
				// TODO(frank): should we log that we are skipping the file?
				continue
			}
			return nil, err
		}

		// concatenate the buffers and add a new line in case
		// the previous file didn't end with a new line
		buff = append(append(buff, buf...), '\n')
	}

	return loadBuf(buff, enc)
}

func mustLoadFiles(filenames []string, enc Encoding, ignoreMissing bool) *Properties {
	p, err := loadFiles(filenames, enc, ignoreMissing)
	if err != nil {
		panic(err)
	}
	return p
}

// Interprets a byte buffer either as an ISO-8859-1 or UTF-8 encoded string.
// For ISO-8859-1 we can convert each byte straight into a rune since the
// first 256 unicode code points cover ISO-8859-1.
func convert(buf []byte, enc Encoding) string {
	switch enc {
	case UTF8:
		return string(buf)
	case ISO_8859_1:
		runes := make([]rune, len(buf))
		for i, b := range buf {
			runes[i] = rune(b)
		}
		return string(runes)
	default:
		panic(fmt.Sprintf("unsupported encoding %v", enc))
	}
}