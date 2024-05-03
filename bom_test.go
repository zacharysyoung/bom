package bomsniffer

import (
	"strings"
	"testing"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/transform"
)

// 'Strings, bytes, runes and characters in Go', by Rob Pike,
// <https://go.dev/blog/strings>
const data = "It’s important to state right up front that a string holds *arbitrary* bytes."

const (
	_16be    = unicode.BigEndian
	_16le    = unicode.LittleEndian
	_16BOM   = unicode.UseBOM
	_16NoBOM = unicode.IgnoreBOM
	_32be    = utf32.BigEndian
	_32le    = utf32.LittleEndian
	_32BOM   = utf32.UseBOM
	_32NoBOM = utf32.IgnoreBOM
)

func TestRead(t *testing.T) {
	// testCases sets up the encoding of data with established
	// golang.org/x/text encoders to verify that sniffer correctly
	// detects a given BOM, and the lack of a BOM.
	var testCases = []struct {
		name    string
		encoder *encoding.Encoder
		want    BOM
	}{
		{"utf8-bom", unicode.UTF8BOM.NewEncoder(), UTF8},
		{"utf16-be-bom", unicode.UTF16(_16be, _16BOM).NewEncoder(), UTF16BE},
		{"utf16-le-bom", unicode.UTF16(_16le, _16BOM).NewEncoder(), UTF16LE},
		{"utf32-be-bom", utf32.UTF32(_32be, _32BOM).NewEncoder(), UTF32BE},
		{"utf32-le-bom", utf32.UTF32(_32le, _32BOM).NewEncoder(), UTF32LE},

		{"utf8", unicode.UTF8.NewEncoder(), Unknown},
		{"utf16-be", unicode.UTF16(_16be, _16NoBOM).NewEncoder(), Unknown},
		{"utf16-le", unicode.UTF16(_16le, _16NoBOM).NewEncoder(), Unknown},
		{"utf32-be", utf32.UTF32(_32be, _32NoBOM).NewEncoder(), Unknown},
		{"utf32-le", utf32.UTF32(_32le, _32NoBOM).NewEncoder(), Unknown},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := make([]byte, 1024)
			tr := transform.NewReader(strings.NewReader(data), tc.encoder)
			sniffer := NewReader(tr)
			n, err := sniffer.Read(p)
			if err != nil {
				t.Errorf("error encoding as %s, read %q before %v", tc.name, p[:n], err)
			}
			if got := sniffer.BOM(); got != tc.want {
				t.Errorf("got %s; want %s", got, tc.want)
			}
		})
	}
}

// TestIncremental reads UTF32-BE-BOM bytes two bytes at a time,
// twice, and asserts that after the first read the BOM is unknown,
// and that after the second read the BOM is known to be UTF32BE.
func TestIncremental(t *testing.T) {
	p := make([]byte, 2)
	tr := transform.NewReader(strings.NewReader(data), utf32.UTF32(_32be, _32BOM).NewEncoder())
	sniffer := NewReader(tr)

	for i, bom := range []BOM{Unknown, UTF32BE} {
		n, err := sniffer.Read(p)
		if err != nil {
			t.Fatalf("unepxected error after read #%d, n:%d buf:%v: %v", i+1, n, p, err)
		}
		if got := sniffer.BOM(); got != bom {
			t.Errorf("after read #%d, got BOM %s; want %s", i+1, got, bom)
		}
	}
}

// TestBufferLimit asserts that SnifferReader doesn't needlessly
// inspect bytes past the front of the file.
func TestBufferLimit(t *testing.T) {
	tr := transform.NewReader(strings.NewReader(data), unicode.UTF8.NewEncoder())
	sniffer := NewReader(tr)

	for i, wantN := range []int{32, 32, 15} {
		p := make([]byte, 32)
		n, err := sniffer.Read(p)
		if err != nil {
			t.Fatalf("unepxected error after read #%d, n:%d buf:%v: %v", i+1, n, p, err)
		}
		if n != wantN {
			t.Errorf("after read #%d, got %d bytes back; want %d", i+1, n, wantN)
		}
		if n = len(sniffer.buf); n > maxSize {
			t.Errorf("after read #%d, sniffer's internal buffer was %d bytes long; want %d", i+1, n, maxSize)
		}
	}
}
