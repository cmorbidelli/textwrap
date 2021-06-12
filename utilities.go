package textwrap

import "strings"

// Wrap is a convenience function corresponding to TextWrapper.Wrap.
// It accepts all of the same options as a TextWrapper.  As each 
// call creates a new TextWrapper, programs that need to perform the
// operation repeatedly should use the Wrap method instead.
func Wrap(text string, opts ...option) []string {
    t := NewTextWrapper(opts...)
    return t.Wrap(text)
}

// Fill is a convenience function corresponding to TextWrapper.Fill.
// It accepts all of the same options as a TextWrapper.  As each 
// call creates a new TextWrapper, programs that need to perform the
// operation many times should use the Fill method instead.
func Fill(text string, opts ...option) string {
    t := NewTextWrapper(opts...)
    return t.Fill(text)
}

// Shorten attempts to fit text onto a single line by replacing any
// sequences of whitespace with a single Space, then returning the
// first line of wrapped text.  While it accepts all of the same
// options as NewTextWrapper, keep in mind that ExpandTabs, TabSize,
// ReplaceWhitespace, and Maxlines have no effect.
func Shorten(text string, opts ...option) string {
    t := NewTextWrapper(opts...)
    t.MaxLines = 1

    text = ConsWhitespaceRe.ReplaceAllString(text, Space)

    return t.Fill(text)
}

// Dedent removes the indent--that is, any leading whitespace shared
// by all lines--from each line of text.  Lines consisting entirely
// of whitespace are ignored.
func Dedent(text string) string {
    lines := strings.Split(text, Newline)
    var indent string

    start := true
    for i, line := range lines {
        if strip(line) == "" {
            lines[i] = ""
        } else if start {
            indent, start = LeadWhitespaceRe.FindString(line), false
        } else if len(indent) != 0 {
            s, t := []rune(indent), []rune(line)
            var j int
            for ; j < len(s) && j < len(t) && s[j] == t[j]; j++ {

            }

            indent = string(s[:j])
        }
    }

    for i, _ := range lines {
        lines[i] = strings.TrimPrefix(lines[i], indent)
    }

    return strings.Join(lines, Newline)
}

// Indent prepends pref to lines within text.  Lines consisting only
// of whitespace are ignored.  If pred is nil, each line is indented;
// otherwise, only lines for which pred(line) == true are indented.
func Indent(text, pref string, pred func(string) bool) string {
    lines := strings.Split(text, Newline)
    for i, line := range lines {
        if strip(line) == "" {
            continue
        }

        if pred == nil || pred(line) {
            lines[i] = pref + line
        }
    }

    return strings.Join(lines, Newline)
}
