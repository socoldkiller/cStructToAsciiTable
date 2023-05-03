package main

import (
	"fmt"
	"regexp"
	"strings"
)

func skipSpace(ctx string) string {
	return strings.TrimSpace(ctx)
}

type CVariable struct {
	parseString    string
	StructKeywords string `json:"struct_keywords,omitempty"`
	Typename       string `json:"typename,omitempty"`
	Pointer        string `json:"pointer,omitempty"`
	Variable       string `json:"variable,omitempty"`
	Value          string `json:"value,omitempty"`
	errorInfo      string
	isParseError   bool
	process        func()
	Struct         *CStruct `json:"struct,omitempty"`
}

func (c *CVariable) toString() string {
	var varData string
	if c.StructKeywords == "" {
		varData = fmt.Sprintf("%s%s %s", c.Typename, c.Pointer, c.Variable)
	} else {
		varData = fmt.Sprintf("%s %s%s %s", c.StructKeywords, c.Typename, c.Pointer, c.Variable)
	}
	return varData
}

func (c *CVariable) parseStruct() {
	c.Struct = &CStruct{}
	c.Struct.Parse(c.parseString)
	c.parseString = c.Struct.parseString
	c.process = c.parseVariable
}

func (c *CVariable) parseError() {
	// fmt.Printf("\"%s\" parse error\n", ctx)
	c.process = nil
	c.isParseError = true
}

// bug.
func (c *CVariable) parseStructKeywords() {
	re := regexp.MustCompile("^struct")
	c.process = c.parseTypename
	tmp := c.parseString
	f := true
	if loc := re.FindStringIndex(c.parseString); loc != nil {
		f = false
		c.StructKeywords = c.parseString[:loc[1]]
		c.parseString = c.parseString[loc[1]:]
	}

	c.parseString = skipSpace(c.parseString)
	c.parseTypename()
	c.parseString = skipSpace(c.parseString)

	if c.parseString[0] == '{' {
		c.parseString = tmp
		c.process = c.parseStruct
		return
	}
	c.process = c.parseTypename
	if f == false {
		c.parseString = tmp[6:]
		return
	}
	c.parseString = tmp

}

func (c *CVariable) parseTypename() {
	re := regexp.MustCompile("^[_a-zA-Z][_a-zA-Z0-9]*")
	if loc := re.FindStringIndex(c.parseString); loc != nil {
		c.Typename = c.parseString[:loc[1]]
		c.parseString = c.parseString[loc[1]:]
		c.process = c.parsePointer
		return
	}
	c.process = c.parseError
}

func (c *CVariable) parsePointer() {

	if len(c.parseString) == 0 {
		c.process = c.parseError
	}

	// 无论是否是指针，都要跳转到值
	if c.parseString[0] == '*' {
		c.Pointer = "*"
		c.parseString = c.parseString[1:]
	}

	c.process = c.parseVariable
}

func (c *CVariable) parseVariable() {
	re := regexp.MustCompile("^[_a-zA-Z][_a-zA-Z0-9]*")
	if loc := re.FindStringIndex(c.parseString); loc != nil {
		c.Variable = c.parseString[:loc[1]]
		c.parseString = c.parseString[loc[1]:]
		c.process = c.parseTheEnd
		return
	}

	if c.parseString[0] == ';' {
		c.process = c.parseTheEnd
		c.Variable = ""
		return
	}

	c.process = c.parseError
}

func (c *CVariable) parseTheEnd() {

	if len(c.parseString) == 0 {
		c.process = c.parseError
		c.parseString = ""
		c.errorInfo = "missing \";\""
		return
	}

	if c.parseString[0] == ';' {
		c.process = nil
		c.parseString = c.parseString[1:]
		//c.parseString = ""
		return
	}

	c.process = c.parseError
}

func (c *CVariable) StringTo(str string) *CVariable {
	c.isParseError = false
	c.parseString = str
	c.process = c.parseStructKeywords
	for {
		process := c.process
		if process == nil {
			break
		}
		c.parseString = skipSpace(c.parseString)
		process()
	}
	return c
}

type CStruct struct {
	StructName        string `json:"structName,omitempty"`
	parseString       string
	AnonymousFunction bool        `json:"anonymousFunction"`
	VarList           []CVariable `json:"varList,omitempty"`
}

func (c *CStruct) parseVarList() {
	for {
		var variable CVariable

		c.parseString = skipSpace(c.parseString)
		if c.parseString == "" || c.parseString[0] == '}' {
			break
		}
		variable.StringTo(c.parseString)
		if variable.isParseError == true {
			break
		}
		c.VarList = append(c.VarList, variable)
		c.parseString = variable.parseString
	}
}

func (c *CStruct) parseStructName() {
	index := strings.Index(c.parseString, "struct")
	c.parseString = c.parseString[index+6:]
	c.parseString = skipSpace(c.parseString)
	re := regexp.MustCompile("^[_a-zA-Z][_a-zA-Z0-9]*")
	if loc := re.FindStringIndex(c.parseString); loc != nil {
		c.StructName = c.parseString[:loc[1]]
		c.parseString = c.parseString[loc[1]:]
		c.AnonymousFunction = false
		return
	}

	c.AnonymousFunction = true

}

func (c *CStruct) Parse(str string) {
	c.parseString = str

	c.parseStructName()
	c.parseString = skipSpace(c.parseString)

	if c.parseString[0] == '{' {
		c.parseString = c.parseString[1:]
	}
	c.parseVarList()
	c.parseString = skipSpace(c.parseString)

	if c.parseString[0] == '}' {
		c.parseString = c.parseString[1:]
	}
}
