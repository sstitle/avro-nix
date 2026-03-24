package main

import (
	"fmt"
	"os"
	"sort"

	goavro "github.com/linkedin/goavro/v2"
	"github.com/pterm/pterm"
)

func main() {
	if len(os.Args) < 2 {
		pterm.Error.Println("Usage: read_item <file.avro>")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		pterm.Error.Printf("Could not open file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	ocfr, err := goavro.NewOCFReader(f)
	if err != nil {
		pterm.Error.Printf("Could not create Avro reader: %v\n", err)
		os.Exit(1)
	}

	var records []map[string]interface{}
	for ocfr.Scan() {
		datum, err := ocfr.Read()
		if err != nil {
			pterm.Error.Printf("Could not read record: %v\n", err)
			os.Exit(1)
		}
		records = append(records, datum.(map[string]interface{}))
	}

	if len(records) == 0 {
		pterm.Warning.Println("No records found.")
		return
	}

	fields := make([]string, 0, len(records[0]))
	for k := range records[0] {
		fields = append(fields, k)
	}
	sort.Strings(fields)

	header := make([]string, len(fields))
	for i, f := range fields {
		header[i] = pterm.Bold.Sprint(f)
	}

	tableData := pterm.TableData{header}
	for _, record := range records {
		row := make([]string, len(fields))
		for i, field := range fields {
			row[i] = fmt.Sprintf("%v", record[field])
		}
		tableData = append(tableData, row)
	}

	pterm.Println()
	pterm.DefaultHeader.WithFullWidth().Printf("Item Records")
	pterm.Println()
	pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(tableData).Render()
	pterm.Println()
	pterm.Info.Printf("%d record(s)\n", len(records))
}
