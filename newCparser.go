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

func skipAll(ctx string) string {
	ctx = skipSpace(ctx)
	ctx = skipComment(ctx)
	ctx = skipSpace(ctx)
	return ctx
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
	parseString     string
	KeyWords        string  `json:"key-words,omitempty"`
	TypeName        string  `json:"type-name,omitempty"`
	Pointer         string  `json:"pointer,omitempty"`
	VarName         string  `json:"var-name,omitempty"`
	CVarList        []*CVar `json:"c-var-list,omitempty"`
	ArrayLengthName string  `json:"array-length-name,omitempty"`
	process         Process
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
	c.parseString = skipAll(c.parseString)
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

	switch c.parseString[0] {
	case '[':
		c.process = c.parseArray
	case ',':
		c.process = c.parseComma
	default:
		c.process = c.parseEnd
	}

}
func (c *CVar) parsePointer() {
	assert(c.parseString[0] == '*', "parsePointer error")
	c.Pointer += "*"
	c.parseString = c.parseString[1:]
	c.parseString = skipAll(c.parseString)
	switch c.parseString[0] {
	case '*':
		c.process = c.parsePointer
	default:
		c.process = c.parseVarName
	}
}

func (c *CVar) parseArray() {
	assert(c.parseString[0] == '[', "parseArray error")
	i := 0
	for {
		i++
		if c.parseString[i] == ']' {
			c.parseString = c.parseString[i+1:]
			break

		}
		if i > len(c.parseString) {
			c.process = c.parseError
			return

		}
		c.ArrayLengthName += string(c.parseString[i])
	}

	switch c.parseString[0] {
	case ',':
		c.process = c.parseComma
	case ';':
		c.process = c.parseEnd
	}

}

func (c *CVar) parseCVarList() {
	for {
		c.parseString = skipAll(c.parseString)
		if len(c.parseString) == 0 {
			c.process = c.parseError
			return
		}
		cvar := &CVar{}
		switch c.parseString[0] {
		case ';':
			c.parseEnd()
			continue
		case '}':
			c.process = c.parseRightBracket
			return
		case ',':
			c.parseString = c.parseString[1:]
			c.parseString = skipAll(c.parseString)
			cvar.parseString = c.parseString
			switch cvar.parseString[0] {
			case '*':
				cvar._parse(cvar.parseString, cvar.parsePointer)
			default:
				cvar._parse(cvar.parseString, cvar.parseVarName)
			}
			cvar.TypeName = c.CVarList[0].TypeName
			c.parseString = skipAll(c.parseString)
		default:
			cvar.parse(c.parseString)
		}

		if cvar.isParserError {
			c.process = c.parseError
			return
		}

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
	c.parseString = skipAll(c.parseString)
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

func (c *CVar) parseComma() {
	assert(c.parseString[0] == ',', "parseComma error")
	c.process = nil
}

func (c *CVar) parse(parseStr string) {
	c._parse(parseStr, c.parseKeyWords)
}
func (c *CVar) _parse(parseStr string, startProcess Process) {
	c.parseString = parseStr
	c.process = startProcess
	for {
		if c.process == nil {
			break
		}
		c.parseString = skipAll(c.parseString)
		if len(c.parseString) == 0 {
			c.process = c.parseError
		}
		c.process()
	}
}

func (c *CVar) getTypeName() string {
	if c.ArrayLengthName != "" {
		return fmt.Sprintf("%s%s[%s]", c.TypeName, c.Pointer, c.ArrayLengthName)
	} else {
		return fmt.Sprintf("%s%s", c.TypeName, c.Pointer)

	}
}
