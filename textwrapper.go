package textwrap

import (
    "regexp"
    "strings"
)

type TextWrapper struct {
    Width              int
    ExpandTabs         bool
    TabSize            int
    ReplaceWhitespace  bool
    DropWhitespace     bool
    InitialIndent      string
    SubsequentIndent   string
    FixSentenceEndings bool
    BreakLongWords     bool
    BreakOnHyphens     bool
    MaxLines           int
    Placeholder        string

    Whitespace         *regexp.Regexp
    SentenceEnding     *regexp.Regexp
    ChunksHyphen       *regexp.Regexp
    ChunksNoHyphen     *regexp.Regexp
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
        MaxLines:           0,
        Placeholder:        " [...]",

        Whitespace:       regexp.MustCompile("[[:space:]]+"),
        SentenceEnding:   regexp.MustCompile("([^[:space:]]" +
                                             "[.!?]['\"]?) [ ]?"),
        ChunksHyphen:     regexp.MustCompile("(\u2014|" +
                                             "[^[:space:]]+-|" +
                                             "[^[:space:]]+|" +
                                             "[[:space:]]+)"),
		    ChunksNoHyphen:   regexp.MustCompile("(\u2014|" +
                                             "[^[:space:]]+|" +
                                             "[:space:]+)"),
    }

    for _, opt := range opts {
        opt(&t)
    }

    return t
}

func (t *TextWrapper) validate() {
    if t.Width < 1 {
        panic("Width must be at least 1")
    } else if t.ExpandTabs && t.TabSize < 0 {
        panic("Tab size must be at least 0 to expand tabs")
    }

    if t.MaxLines > 0 {
        indent := t.SubsequentIndent
        if t.MaxLines == 1 {
            indent = t.InitialIndent
        }

        if len(indent) + len(lStrip(t.Placeholder)) > t.Width {
            panic("Placeholder is too wide to fit on indented line")
        }
    }
}

func (t *TextWrapper) Wrap(text string) []string {
    t.validate()

    if t.ExpandTabs {
        tab := strings.Repeat(" ", t.TabSize)
        text = strings.Replace(text, "\t", tab, -1)
    }

    if t.ReplaceWhitespace {
        text = t.Whitespace.ReplaceAllString(text, " ")
    }

    if t.FixSentenceEndings {
        text = t.SentenceEnding.ReplaceAllString(text, "${1}  ")
    }

    var chunks, lines []string
    if t.BreakOnHyphens {
        chunks = t.ChunksHyphen.FindAllString(text, -1)
    } else {
        chunks = t.ChunksNoHyphen.FindAllString(text, -1)
    }

    for i := 0; i < len(chunks); i++ {
        if len(lines) > 0 && t.DropWhitespace && isSpace(chunks[i]) {
            i++
        }

        var indent string
        if len(lines) > 0 {
            indent = t.SubsequentIndent
        } else {
            indent = t.InitialIndent
        }

        width := t.Width - len([]rune(indent))

        curLine := line{}
        for ; i < len(chunks); i++ {
            if curLine.length + len([]rune(chunks[i])) < width {
                curLine.push(chunks[i])
						} else {
                //i--
                break
            }
        }

        if i + 1 < len(chunks) && len([]rune(chunks[i+1])) > width {
            i++
            if t.BreakLongWords {
                c := []rune(chunks[i])
                spaceLeft := 1
                if width >= 1 {
                    spaceLeft = width - curLine.length
                }
                curLine.push(string(c[:spaceLeft]))
                chunks[i] = string(c[spaceLeft:])
                i--
            } else if curLine.length == 0 {
                curLine.push(chunks[i])
            } else {
							  i--
						}
        }

        if curLine.length == 0 {
            continue
        }

        if t.DropWhitespace &&
           isSpace(curLine.chunks[len(curLine.chunks) - 1]) {
            curLine.pop()
        }

        if t.MaxLines < 1 || len(lines) + 1 < t.MaxLines ||
           (t.DropWhitespace && i == len(chunks) - 1 &&
           isSpace(chunks[i])) {
            lines = append(lines, indent + curLine.str())
        } else {
            for j := len(curLine.chunks) - 1; j > 0; j-- {
                if isSpace(curLine.chunks[j]) &&
                   curLine.length + len([]rune(t.Placeholder)) < width {
                    curLine.push(t.Placeholder)
                    lines = append(lines, indent + curLine.str())
                    break
                }
                curLine.pop()
            }

            if l := len(lines); l > 0 {
                prevLine := rStrip(lines[l - 1]) + t.Placeholder
                if len([]rune(prevLine)) <= t.Width {
                    lines[l - 1] = prevLine
                    break
                }
            }

            lines = append(lines, indent + lStrip(t.Placeholder))
            break
        }
    }

    return lines
}

func (t *TextWrapper) Fill(text string) string {
    return strings.Join(t.Wrap(text), "\n")
}
