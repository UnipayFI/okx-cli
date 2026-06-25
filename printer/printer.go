package printer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/UnipayFI/okx-cli/config"
	"github.com/olekukonko/tablewriter"
)

// TableWriter is implemented by every response wrapper that knows how to render
// itself as a table. The header is the column titles; each Row is one record.
type TableWriter interface {
	Header() []string
	Row() [][]any
}

// Print writes v to stdout. With --json enabled it always emits indented JSON;
// otherwise it renders v as a table when v implements TableWriter, falling back
// to JSON for everything else.
func Print(v any) {
	if config.Config.OutputJSON {
		printJSON(v)
		return
	}
	if tw, ok := v.(TableWriter); ok {
		PrintTable(tw)
		return
	}
	printJSON(v)
}

func PrintTable(writer TableWriter) {
	table := tablewriter.NewTable(os.Stdout, tablewriter.WithEastAsian(false))
	table.Header(writer.Header())
	table.Bulk(writer.Row())
	table.Render()
}

func printJSON(v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr, "json encode error:", err)
		return
	}
	fmt.Println(string(b))
}
