package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"os"
	"strings"
)

func SetTableFormat(table *tablewriter.Table) {
	table.SetAlignment(tablewriter.ALIGN_LEFT) // Set Alignment
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_CENTER)
}

func getTableFormatString(str string) string {
	var b bytes.Buffer
	table, _ := tablewriter.NewCSVReader(&b, csv.NewReader(strings.NewReader(str)), true)
	SetTableFormat(table)
	table.Render()
	return b.String()
}

func main() {
	//c := os.Args[1]
	//b := getTableFormatString(c)
	//data := MultilineComment(b)
	//fmt.Println(data)
	//defer func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			fmt.Println(r)
	//		}
	//	}()
	//}()

	//c := CVariable{}
	//c.StringTo("struct int *p;int a;int b")
	//

	var c CStruct
	file, err := os.ReadFile("test.abc")
	if err != nil {
		return
	}

	c.StringTo(string(file))
	var b []byte
	b, _ = json.Marshal(&c)

	fmt.Printf("%s", b)
}
