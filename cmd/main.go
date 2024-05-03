package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/zacharysyoung/bomsniffer"
)

func main() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)

	for _, fname := range []string{
		// (no BOM) the original the BOM-variants were created from
		"utf8.txt",

		"utf16lebom.txt",
		"utf32bebom.txt",

		// comment-out if your OS doesn't have it
		"/usr/share/dict/web2",
	} {
		f, err := os.Open(fname)
		if err != nil {
			panic(err)
		}
		sniffer := bomsniffer.NewReader(f)
		io.Copy(io.Discard, sniffer)
		fmt.Fprintf(w, "%s\t%s\n", fname, sniffer.BOM())
	}

	w.Flush()

	// Output:
	// utf8.txt             Unkown
	// utf16lebom.txt       UTF16LE
	// utf32bebom.txt       UTF32BE
	// /usr/share/dict/web2 Unkown
}
