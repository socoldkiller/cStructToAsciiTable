package main

import (
	"fmt"
	"regexp"
	"strings"
)

type CVariable struct {
	parseString string
	Struct      string `json:"struct"`
	Typename    string `json:"typename"`
	Pointer     string `json:"pointer"`
	Variable    string `json:"variable"`

	Value        string `json:"value"`
	errorInfo    string
	isParseError bool
	process      func(string)
}

func skipSpace(ctx string) string {
	return strings.TrimSpace(ctx)
}

func (c *CVariable) toString() string {
	var varData string
	if c.Struct == "" {
		varData = fmt.Sprintf("%s%s %s", c.Typename, c.Pointer, c.Variable)
	} else {
		varData = fmt.Sprintf("%s %s%s %s", c.Struct, c.Typename, c.Pointer, c.Variable)
	}
	return varData
}

func (c *CVariable) parseError(ctx string) {
	// fmt.Printf("\"%s\" parse error\n", ctx)
	c.process = nil
	c.isParseError = true
}

func (c *CVariable) parseStruct(ctx string) {
	re := regexp.MustCompile("^struct")
	c.process = c.parseTypename
	if loc := re.FindStringIndex(ctx); loc != nil {
		c.Struct = ctx[:loc[1]]
		c.parseString = c.parseString[loc[1]:]
	}
}

func (c *CVariable) parseTypename(ctx string) {
	re := regexp.MustCompile("^[_a-zA-Z][_a-zA-Z0-9]*")
	if loc := re.FindStringIndex(ctx); loc != nil {
		c.Typename = ctx[:loc[1]]
		c.parseString = c.parseString[loc[1]:]
		c.process = c.parsePointer
		return
	}
	c.process = c.parseError
}

func (c *CVariable) parsePointer(ctx string) {

	if len(ctx) == 0 {
		c.process = c.parseError
	}

	// 无论是否是指针，都要跳转到值
	if ctx[0] == '*' {
		c.Pointer = "*"
		c.parseString = c.parseString[1:]
	}

	c.process = c.parseVariable
}

func (c *CVariable) parseVariable(ctx string) {
	re := regexp.MustCompile("^[_a-zA-Z][_a-zA-Z0-9]*")
	if loc := re.FindStringIndex(ctx); loc != nil {
		c.Variable = ctx[:loc[1]]
		c.parseString = c.parseString[loc[1]:]
		c.process = c.parseTheEnd
		return
	}
	c.process = c.parseError
}

func (c *CVariable) parseTheEnd(ctx string) {

	if len(ctx) == 0 {
		c.process = c.parseError
		c.parseString = ""
		c.errorInfo = "missing \";\""
		return
	}

	if ctx[0] == ';' {
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
	c.process = c.parseStruct
	for {
		process := c.process
		if process == nil {
			break
		}
		c.parseString = skipSpace(c.parseString)
		process(c.parseString)
	}

	return c
}

type CStruct struct {
	StructName        string `json:"structName"`
	parseString       string
	AnonymousFunction bool        `json:"anonymousFunction"`
	VarList           []CVariable `json:"varList"`
}

func (c *CStruct) parseVarList() {
	var variable CVariable
	for {
		if c.parseString == "" {
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

func (c *CStruct) StringTo(str string) {
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
