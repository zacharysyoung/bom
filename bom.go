// Package bomsniffer provides a simple, well-assuming io.Reader,
// that sniffs for one of the five UTF BOMs.
//
// Well-assuming in that it has only been tested with in-memory
// Readers, and a few times on disk—so, maybe insert "toy" into
// the above description as well.
package bomsniffer

import (
	"bytes"
	"fmt"
	"io"
)

// BOM represents one of the five UTF variants of a Byte Order Mark,
// or Unknown.
type BOM int

const (
	Unknown BOM = iota
	UTF8
	UTF16BE
	UTF16LE
	UTF32BE
	UTF32LE
)

// String returns the name of bom; implements Stringer.
func (bom BOM) String() string {
	switch bom {
	case Unknown:
		return "Unkown"
	case UTF8:
		return "UTF8"
	case UTF16BE:
		return "UTF16BE"
	case UTF16LE:
		return "UTF16LE"
	case UTF32BE:
		return "UTF32BE"
	case UTF32LE:
		return "UTF32LE"
	default:
		panic(fmt.Errorf("unexpected BOM value %d", bom))
	}
}

// Reader wraps an io.Reader and inspects the first few bytes
// read for the presence of a BOM prefix. Past the first few
// bytes, it stops looking.  It does not infer any encoding of
// the bytes read.
type Reader struct {
	bom BOM
	r   io.Reader
	buf []byte // scratch buffer
}

const maxSize = 4 // the max size of the scratch buffer

// NewReader returns a new Reader that sniffs as it reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{Unknown, r, make([]byte, 0, maxSize)}
}

// BOM returns a detected BOM prefix, or Unknown.
func (sr *Reader) BOM() BOM { return sr.bom }

// Read reads data into p, and inspects the first few bytes of
// p for a BOM prefix.
func (sr *Reader) Read(p []byte) (n int, err error) {
	n, err = sr.r.Read(p)

	if sr.bom == Unknown && len(sr.buf) < maxSize {
		n := len(p)
		if n >= maxSize {
			n = maxSize
		}
		sr.buf = append(sr.buf, p[:n]...)
		switch {
		case bytes.HasPrefix(sr.buf, []byte{0x00, 0x00, 0xFE, 0xFF}):
			sr.bom = UTF32BE
		case bytes.HasPrefix(sr.buf, []byte{0xFF, 0xFE, 0x00, 0x00}):
			sr.bom = UTF32LE
		case bytes.HasPrefix(sr.buf, []byte{0xEF, 0xBB, 0xBF}):
			sr.bom = UTF8
		case bytes.HasPrefix(sr.buf, []byte{0xFE, 0xFF}):
			sr.bom = UTF16BE
		case bytes.HasPrefix(sr.buf, []byte{0xFF, 0xFE}):
			sr.bom = UTF16LE
		}
	}

	return
}
