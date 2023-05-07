package main

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
)

func SetTableFormat(table *tablewriter.Table) {
	table.SetAlignment(tablewriter.ALIGN_LEFT) // Set Alignment
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
}

func getTable(c *CVar, tables map[string][][]string) {
	var writer bytes.Buffer
	// table, _ := tablewriter.NewCSVReader(&writer, csv.NewReader(strings.NewReader(str)), true)
	table := tablewriter.NewWriter(&writer)
	SetTableFormat(table)
	var t [][]string
	for _, cvar := range c.CVarList {
		t = append(t, []string{cvar.getTypeName(), cvar.VarName})
		if cvar.CVarList != nil {
			getTable(cvar, tables)
		}
	}

	tables[c.TypeName] = t
}

func getTableFormatString(data [][]string) string {
	var writer bytes.Buffer
	// table, _ := tablewriter.NewCSVReader(&writer, csv.NewReader(strings.NewReader(str)), true)
	table := tablewriter.NewWriter(&writer)
	SetTableFormat(table)
	table.AppendBulk(data)
	table.Render()
	return writer.String()
}

func main() {
	s := make(map[string][][]string)
	var c CVar
	file, err := os.ReadFile("test.abc")
	if err != nil {
		return
	}
	data := string(file)
	c.parse(data)
	getTable(&c, s)
	for _, v := range s {
		fmt.Printf("%s", MultilineComment(getTableFormatString(v)))
		fmt.Printf("\n")
	}
}
