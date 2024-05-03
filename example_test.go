package bomsniffer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// Create a sniffing Reader around bytes with a UTF-8BOM prefix
// and hand it off to bufio.NewScanner.
//
// Since Go source cannot contain a BOM, this example skips the
// first three bytes that are a UTF-8-encode BOM in the first
// line; the following lines don't need any special processing.
func ExampleReader() {
	in := "\uFEFF" + `line 1
line 2
line 3
`

	sniffer := NewReader(strings.NewReader(in))
	scanner := bufio.NewScanner(sniffer)

	scanner.Scan()
	fmt.Println(scanner.Text()[3:]) // discard BOM bytes

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	fmt.Println("\nText has BOM", sniffer.BOM())

	// Output:
	// line 1
	// line 2
	// line 3
	//
	// Text has BOM UTF-8
}

// Like the previous example, but use a [golang.org/x/text/transform.Reader]
// that decodes UTF-8-encoded bytes with a BOM prefix to UTF-8
// bytes without a BOM.  Then create a sniffing Reader and Scanner
// to print the lines and assert the BOM was not detected.
func ExampleReader_transformer() {
	in := "\uFEFF" + `line 1
line 2
line 3
`

	r := transform.NewReader(
		strings.NewReader(in),
		unicode.UTF8BOM.NewDecoder())

	buf := &bytes.Buffer{}

	io.Copy(buf, r)

	sniffer := NewReader(buf)
	scanner := bufio.NewScanner(sniffer)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	fmt.Println("\nText has BOM", sniffer.BOM())

	// Output:
	// line 1
	// line 2
	// line 3
	//
	// Text has BOM Unknown
}
