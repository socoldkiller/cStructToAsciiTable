package main

import (
	"bytes"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
)

func SetTableFormat(table *tablewriter.Table) {
	table.SetAlignment(tablewriter.ALIGN_LEFT) // Set Alignment
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader([]string{"TypeName", "VarName"})
}

func getTable(c *CVar, tables map[string][][]string) {
	var t [][]string
	t = append(t, []string{c.getTypeName(), c.VarName})
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
	table := tablewriter.NewWriter(&writer)
	SetTableFormat(table)
	table.AppendBulk(data)
	table.Render()
	return writer.String()
}

func main() {
	s := make(map[string][][]string)
	var c CVar
	if len(os.Args) < 2 {
		log.Fatal("args must Greater than 1")
	}
	cStr := os.Args[1]
	c.parse(cStr)
	getTable(&c, s)
	for _, v := range s {
		fmt.Printf("%s", MultilineComment(getTableFormatString(v)))
		fmt.Printf("\n\n\n")
	}
}
