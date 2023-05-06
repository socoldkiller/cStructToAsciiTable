package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"
)

func assert(flag bool, str string) {
	if !flag {
		panic(str)
	}
}

func skipSpace(ctx string) string {
	return strings.TrimSpace(ctx)
}

func skipComment(ctx string) string {

	if len(ctx) < 2 {
		return ctx
	}

	i := 0
	if ctx[0:2] == "/*" {
		for {
			i++
			if ctx[i:i+2] == "*/" {
				return ctx[i+2:]

			}
		}
	}
	return ctx
}

type Process func()

type CVar struct {
	parseString string
	KeyWords    string  `json:"key_words,omitempty"`
	TypeName    string  `json:"type_name,omitempty"`
	Pointer     string  `json:"pointer,omitempty"`
	VarName     string  `json:"var_name,omitempty"`
	CVarList    []*CVar `json:"c_var_list,omitempty"`
	process     Process

	isParserError   bool
	parserErrorInfo string
}

func (c *CVar) getParserErrorInfo() {
	if c.isParserError {
		log.Fatalf("%s\n", c.parserErrorInfo)
	}
	fmt.Printf("parse ok\n")
}

func (c *CVar) parseKeyWords() {
	re := regexp.MustCompile("^struct|^enum|^union")
	if loc := re.FindStringIndex(c.parseString); loc != nil {
		c.KeyWords = c.parseString[:loc[1]]
		c.parseString = c.parseString[loc[1]:]
		c.process = c.parseTypeName
		return
	}

	c.process = c.parseTypeName
}

func (c *CVar) parseTypeName() {
	// ctx := c.parseString
	re := regexp.MustCompile("^[_a-zA-Z][_a-zA-Z0-9]*")
	if loc := re.FindStringIndex(c.parseString); loc != nil {
		c.TypeName = c.parseString[:loc[1]]
		c.parseString = c.parseString[loc[1]:]
	}
	c.parseString = skipSpace(c.parseString)

	if len(c.parseString) == 0 {
		c.process = c.parseError
		return
	}

	switch c.parseString[0] {
	case '{':
		c.process = c.parseLeftBracket
	case '*':
		c.process = c.parsePointer
	default:
		c.process = c.parseVarName
	}

}

func (c *CVar) parseVarName() {
	// ctx := c.parseString
	re := regexp.MustCompile("^[_a-zA-Z][_a-zA-Z0-9]*")
	if loc := re.FindStringIndex(c.parseString); loc != nil {
		c.VarName = c.parseString[:loc[1]]
		c.parseString = c.parseString[loc[1]:]
	}
	c.process = c.parseEnd
}
func (c *CVar) parsePointer() {

	if c.parseString[0] == '*' {
		c.Pointer = "*"
		c.parseString = c.parseString[1:]
	}
	c.process = c.parseVarName
}

func (c *CVar) parseCVarList() {
	for {
		c.parseString = skipSpace(c.parseString)

		if len(c.parseString) == 0 {
			c.process = c.parseError
			return
		}

		if c.parseString[0] == '}' {
			c.process = c.parseRightBracket
			return
		}
		cvar := &CVar{}
		cvar.parse(c.parseString)
		c.CVarList = append(c.CVarList, cvar)
		c.parseString = cvar.parseString
	}
}

func (c *CVar) parseLeftBracket() {
	assert(c.parseString[0] == '{', "parser error,parseLeftBracket")
	if c.parseString[0] == '{' {
		c.parseString = c.parseString[1:]
	}
	c.process = c.parseCVarList
}

func (c *CVar) parseRightBracket() {
	assert(c.parseString[0] == '}', "parser error,parseRightBracket")
	c.parseString = c.parseString[1:]
	c.parseString = skipSpace(c.parseString)
	if c.parseString[0] == ';' {
		c.process = c.parseEnd
	} else if unicode.IsLetter(rune(c.parseString[0])) {
		c.process = c.parseVarName
	} else {
		c.process = c.parseError
	}

}
func (c *CVar) parseEnd() {
	if c.parseString[0] == ';' {
		c.process = nil
		c.parseString = c.parseString[1:]
		return
	}
	c.process = c.parseError
}

func (c *CVar) parseError() {
	c.parserErrorInfo = c.parseString
	c.isParserError = true
	c.process = nil
}

func (c *CVar) parse(parseStr string) {
	c.parseString = parseStr
	c.process = c.parseKeyWords
	for {
		process := c.process
		if process == nil {
			break
		}
		c.parseString = skipComment(c.parseString)
		c.parseString = skipSpace(c.parseString)
		if len(c.parseString) == 0 {
			c.process = c.parseError
		}
		process()
	}
}
