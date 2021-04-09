package textwrap

import (
	"regexp"
	"strings"
)

var whitespace = "\t\n\v\f\r \u0085\u00A0"

type option func(*TextWrapper)

func Width(i int) option {
	return func(t *TextWrapper) {
		t.Width = i
	}
}

func ExpandTabs(b bool) option {
	return func(t *TextWrapper) {
		t.ExpandTabs = b
	}
}

func TabSize(i int) option {
	return func(t *TextWrapper) {
		t.TabSize = i
	}
}

func ReplaceWhitespace(b bool) option {
	return func(t *TextWrapper) {
		t.ReplaceWhitespace = b
	}
}

func DropWhitespace(b bool) option {
	return func(t *TextWrapper) {
		t.DropWhitespace = b
	}
}

func InitialIndent(s string) option {
	return func(t *TextWrapper) {
		t.InitialIndent = s
	}
}

func SubsequentIndent(s string) option {
	return func(t *TextWrapper) {
		t.SubsequentIndent = s
	}
}

func FixSentenceEndings(b bool) option {
	return func(t *TextWrapper) {
		t.FixSentenceEndings = b
	}
}

func BreakLongWords(b bool) option {
	return func(t *TextWrapper) {
		t.BreakLongWords = b
	}
}

func BreakOnHyphens(b bool) option {
	return func(t *TextWrapper) {
		t.BreakOnHyphens = b
	}
}

func MaxLines(i int) option {
	return func(t *TextWrapper) {
		t.MaxLines = i
	}
}

func Placeholder(s string) option {
	return func(t *TextWrapper) {
		t.Placeholder = s
	}
}

type TextWrapper struct {
	Width              int
	ExpandTabs         bool   //done
	TabSize            int    //done
	ReplaceWhitespace  bool   //done
	DropWhitespace     bool   //done
	InitialIndent      string //done
	SubsequentIndent   string //done
	FixSentenceEndings bool   //done
	BreakLongWords     bool
	BreakOnHyphens     bool
	MaxLines           int
	Placeholder        string

  whitespace       string
  sentenceEndings  string
  chunks           string
}

func NewTextWrapper(opts ...option) TextWrapper {
	t := TextWrapper{
		Width:              70,
		ExpandTabs:         true,
		TabSize:            8,
		ReplaceWhitespace:  true,
		DropWhitespace:     true,
		InitialIndent:      "",
		SubsequentIndent:   "",
		FixSentenceEndings: false,
		BreakLongWords:     true,
		BreakOnHyphens:     true,
		MaxLines:           -1,
		Placeholder:        " [...]",

    whitespace:      regexp.MustCompile("[" + whitespace + "]+")
    sentenceEndings: regexp.MustCompile("([^" + whitespace + "]" +
                                        "[.!?]['\"]? )([^ ])"
    chunks:          regexp.MustCompile("(\u2014|" +
                                        "[^" + whitespace + "]+|" +
                                        "[" + whitespace + "]+)")
	}

	for _, opt := range opts {
		opt(&t)
	}

	return t
}

func (t TextWrapper) Wrap(text string) []string {
	//since all methods use Wrap, can probably create the regexps in
	//NewTextWrapper

	if t.ExpandTabs {
		tab := strings.Repeat(" ", t.TabSize)
		text = strings.Replace(text, "\t", tab, -1)
	}

	if t.ReplaceWhitespace {
		text = t.whitespace.ReplaceAllString(text, " ")
	}

	if t.FixSentenceEndings {
		text = t.sentenceEndings.ReplaceAllString(text, "${1} ${2}")
	}

	//chunks := re.FindAllString(text, -1)
	//var current string
	lines := []string{}
	//for i, chunk := range chunks {

	//magic

	/*indent := len(t.SubsequentIndent)
	  if len(lines) == 0 {
	      indent = len(t.InitialIndent)
	  }
	  max := t.Width - indent
	  if len(lines) == t.MaxLines {
	      max = t.Width - indent - len(t.Placeholder)
	  }


	  if len(lines) == t.MaxLines && len(current) + len(chunk) > max {
	      current += t.Placeholder
	      lines = append(lines, current)
	      break
	  }*/
	//end magic

	//}

	if t.DropWhitespace {
		var i int
		for _, line := range lines {
			if t.whitespace.FindString(line) == line {
				continue
			} else if i == 0 {
				lines[i] = strings.TrimRight(line, whitespace)
				i++
			} else {
				lines[i] = strings.Trim(line, whitespace)
				i++
			}
		}
		lines = lines[:i]
	}

	if t.InitialIndent == "" && t.SubsequentIndent == "" {
		return lines
	}

	for i, line := range lines {
		if t.whitespace.FindString(line) == line {
			continue
		} else if i == 0 {
			lines[i] = t.InitialIndent + line
		} else {
			lines[i] = t.SubsequentIndent + line
		}
	}

	return lines
}

func (t TextWrapper) Fill(text string) string {
	return strings.Join(t.Wrap(text), "\n")
}

func Wrap(text string, opts ...option) []string {
	t := NewTextWrapper(opts...)
	return t.Wrap(text)
}

func Fill(text string, opts ...option) string {
	t := NewTextWrapper(opts...)
	return t.Fill(text)
}

func Shorten(text string, opts ...option) string {
	t := NewTextWrapper(opts...)
	t.MaxLines = 1
	re := regexp.MustCompile("[" + whitespace + "]+")
	text = re.ReplaceAllString(text, " ")

	return t.Fill(text)
}

func Dedent(text string) string {
	lines := strings.Split(text, "\n")
	var indent string
	start := true
	re := regexp.MustCompile("^[" + whitespace + "]*")
	for i, line := range lines {
		white := re.FindString(line)

		if white == line {
			lines[i] = ""
		} else if start {
			indent, start = white, false
		} else if len(indent) == 0 {
			continue
		} else {
			s, t := []rune(indent), []rune(white)

			var j int
			for ; j < len(s) && j < len(t); j++ {
				if s[j] != t[j] {
					break
				}
			}

			indent = string(s[:i])
		}
	}

	for i, _ := range lines {
		lines[i] = strings.TrimPrefix(lines[i], indent)
	}

	return strings.Join(lines, "\n")
}

func Indent(text, pref string, pred func(string) bool) string {
	lines := strings.Split(text, "\n")
	re := regexp.MustCompile("^[" + whitespace + "]*$")
	for i, line := range lines {
		if re.FindString(line) != "" {
			continue
		}

		if pred == nil || pred(line) {
			lines[i] = pref + line
		}
	}

	return strings.Join(lines, "\n")
}

func Center(text string, pad rune, width int) string {
	if len([]rune(text)) >= width {
		return text
	}

	sides := width - len([]rune(text))
	left := strings.Repeat(string(pad), sides/2)
	right := left
	if sides%2 != 0 {
		right += string(pad)
	}

	return left + text + right
}
