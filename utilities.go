package textwrap

import (
    "regexp"
    "strings"
)

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
    text = t.Whitespace.ReplaceAllString(text, " ")

    return t.Fill(text)
}

func Dedent(text string) string {
    lines := strings.Split(text, "\n")
    re := regexp.MustCompile("^[[:space:]]*")

    var indent string
    start := true
    for i, line := range lines {
        if isSpace(line) {
            lines[i] = ""
        } else if start {
            indent, start = re.FindString(line), false
        } else if len(indent) != 0 {
            s, t := []rune(indent), []rune(line)
            var j int
            for ; j < len(s) && j < len(t); j++ {
                if s[j] != t[j] {
                    break
                }
            }

            indent = string(s[:j])
        }
    }

    for i, _ := range lines {
        lines[i] = strings.TrimPrefix(lines[i], indent)
    }

    return strings.Join(lines, "\n")
}

func Indent(text, pref string, pred func(string) bool) string {
    lines := strings.Split(text, "\n")
    for i, line := range lines {
        if isSpace(line) {
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
    left := strings.Repeat(string(pad), sides / 2)
    right := left
    if sides % 2 != 0 {
        right += string(pad)
    }

    return left + text + right
}
