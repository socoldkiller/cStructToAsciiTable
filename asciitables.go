package main

import (
	"bytes"
	"github.com/olekukonko/tablewriter"
)

func SetTableFormat(table *tablewriter.Table) {
	table.SetAlignment(tablewriter.ALIGN_LEFT) // Set Alignment
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	table.SetHeader([]string{"TypeName", "VarName", "Extra"})
}

func getTable(c *CVar, tables map[string][][]string) {
	var t [][]string
	t = append(t, []string{c.getTypeName(), c.getVarName(), c.Comment})
	for _, cvar := range c.CVarList {
		t = append(t, []string{cvar.getTypeName(), cvar.getVarName(), cvar.Comment})
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
