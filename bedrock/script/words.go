package script

import (
	"fmt"
	"strings"
)

// words splits a command line into words, as sh does for the subset read:
// on blanks outside quotes, with ; a word of its own, the quotes removed,
// and variables expanded in double quotes and outside quotes, but not in
// single quotes.
func (p *parser) words(line string) ([]string, error) {
	var words []string
	var word strings.Builder
	inWord := false
	end := func() {
		if inWord {
			words = append(words, word.String())
			word.Reset()
			inWord = false
		}
	}
	for i := 0; i < len(line); i++ {
		c := line[i]
		switch c {
		case ' ', '\t':
			end()
		case ';':
			end()
			words = append(words, ";")
		case '\'':
			j := strings.IndexByte(line[i+1:], '\'')
			if j < 0 {
				return nil, fmt.Errorf("a quote not closed")
			}
			word.WriteString(line[i+1 : i+1+j])
			i, inWord = i+1+j, true
		case '"':
			j := strings.IndexByte(line[i+1:], '"')
			if j < 0 {
				return nil, fmt.Errorf("a quote not closed")
			}
			word.WriteString(p.expand(line[i+1 : i+1+j]))
			i, inWord = i+1+j, true
		case '$':
			n := variableEnd(line, i)
			word.WriteString(p.expand(line[i:n]))
			i, inWord = n-1, true
		default:
			word.WriteByte(c)
			inWord = true
		}
	}
	end()
	return words, nil
}

// expand expands the variables in text: those the script set, to their
// values, and others, its parameters and the environment's, to
// themselves, $NAME.
func (p *parser) expand(text string) string {
	var b strings.Builder
	for i := 0; i < len(text); i++ {
		if text[i] != '$' {
			b.WriteByte(text[i])
			continue
		}
		n := variableEnd(text, i)
		name := strings.Trim(text[i+1:n], "{}")
		if value, ok := p.vars[name]; ok {
			b.WriteString(value)
		} else {
			b.WriteString("$" + name)
		}
		i = n - 1
	}
	return b.String()
}

// variableEnd is the index after the variable reference that starts at i:
// $NAME, ${NAME}, $1 or $?.
func variableEnd(text string, i int) int {
	j := i + 1
	switch {
	case j < len(text) && text[j] == '{':
		if k := strings.IndexByte(text[j:], '}'); k >= 0 {
			return j + k + 1
		}
	case j < len(text) && (text[j] == '?' || text[j] >= '0' && text[j] <= '9'):
		return j + 1
	}
	for j < len(text) && (text[j] == '_' || text[j] >= 'A' && text[j] <= 'Z' || text[j] >= 'a' && text[j] <= 'z' || text[j] >= '0' && text[j] <= '9') {
		j++
	}
	return j
}
