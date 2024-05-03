package bomsniffer

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// Take the canonical example from package encoding/csv, prepend in
// with a UTF-8-encoded BOM, wrap it with a sniffing Reader, and pass
// that to csv.NewReader.
//
// Making examples with a BOM is difficult because Go source cannot
// contain a BOM, not even the Output: comment, so this example
// discards the first record which actually contains the BOM.
func ExampleReader() {
	in := "\uFEFF" + `first_name,last_name,username
"Rob","Pike",rob
Ken,Thompson,ken
"Robert","Griesemer","gri"
`
	sniffer := NewReader(strings.NewReader(in))
	r := csv.NewReader(sniffer)

	_, _ = r.Read() // discard first record (see description above)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		fmt.Println(record)
	}

	fmt.Println("\nCSV has BOM", sniffer.BOM())

	// Output:
	// [Rob Pike rob]
	// [Ken Thompson ken]
	// [Robert Griesemer gri]
	//
	// CSV has BOM UTF8
}
