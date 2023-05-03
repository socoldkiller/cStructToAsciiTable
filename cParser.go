package main

import (
	"fmt"
	"regexp"
	"strings"
)

type CVariable struct {
	parseString string
	typename    string
	variable    string
	pointer     string
	_struct     string
	value       string
	process     func(string) int
}

func skip_space(ctx string, index *int) string {
	i := 0
	for {

		if ctx[i] == ' ' || ctx[i] == '\n' {

			*index = i
		}

	}
	return ctx[i:]
}

func (c *CVariable) toString() string {
	var varData string
	if c._struct == "" {
		varData = fmt.Sprintf("%s%s %s", c.typename, c.pointer, c.variable)
	} else {
		varData = fmt.Sprintf("%s %s%s %s", c._struct, c.typename, c.pointer, c.variable)
	}
	return varData
}

func (c *CVariable) parseError(ctx string) int {
	fmt.Printf("%s parse error\n", ctx)
	c.process = nil
	return -1
}

func (c *CVariable) parseStruct(ctx string) int {
	re := regexp.MustCompile("\\s*struct\\s+")
	c.process = c.parseTypename
	if loc := re.FindStringIndex(ctx); loc != nil {
		c._struct = "struct"
		return loc[1]
	}
	return 0
}

func (c *CVariable) parseTypename(ctx string) int {
	re := regexp.MustCompile("^\\s*[_a-zA-Z][_a-zA-Z0-9]*\\s*")
	if loc := re.FindStringIndex(ctx); loc != nil {
		c.typename = ctx[loc[0]:loc[1]]
		c.typename = strings.TrimSpace(c.typename)
		c.process = c.parsePointer
		return loc[1]
	}
	c.process = c.parseError
	return 0
}

func (c *CVariable) parsePointer(ctx string) int {
	re := regexp.MustCompile("^\\s*\\*\\s*")
	c.process = c.parseVariable
	if loc := re.FindStringIndex(ctx); loc != nil {
		c.pointer = "*"
		return loc[1]
	}
	return 0
}

func (c *CVariable) parseVariable(ctx string) int {
	re := regexp.MustCompile("^\\s*[_a-zA-Z][_a-zA-Z0-9]*\\s*")
	if loc := re.FindStringIndex(ctx); loc != nil {
		c.variable = ctx[loc[0]:loc[1]]
		c.variable = strings.TrimSpace(c.variable)
		c.process = c.parseTheEnd
		return loc[1]
	}
	c.process = c.parseError
	return 0
}

func (c *CVariable) parseTheEnd(ctx string) int {

	if len(ctx) == 0 {
		fmt.Printf("parse error,missing semicolon\n")
		c.process = nil
		return 0
	}

	if ctx[0] == ';' {
		c.process = nil
		return 1
	}

	c.process = c.parseError
	return 0

}

func (c *CVariable) StringTo(str string) *CVariable {
	next := 0
	nowIndex := 0
	c.process = c.parseStruct
	for {
		process := c.process
		if process == nil {
			break
		}

		next = process(str[nowIndex:])
		nowIndex += next
	}

	return c
}
