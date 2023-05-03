package main

const START = "/*"
const END = "*/"
const Tab = "    "

func MultilineComment(content string) string {
	s := ""
	s += START
	s += "\n"

	s += " *"
	s += Tab + Tab
	s += string(content[0])

	for i := 1; i < len(content); i++ {

		if content[i] == '+' && content[i-1] == '\n' {
			s += " *"
			s += Tab + Tab

		}

		if content[i] == '|' && content[i-1] == '\n' {
			s += " *"
			s += Tab + Tab
		}

		s += string(content[i])
	}

	s += " "
	s += END

	return s
}
