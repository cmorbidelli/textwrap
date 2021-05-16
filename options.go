package textwrap

// options allow the user to override a TextWrapper's default values.
// The function calls can be passed as arguments to NewTextWrapper, 
// Wrap, Fill, and Shorten.  For example:
//     Fill(myText, Width(50), FixSentenceEndings(true))
// fills myText using a TextWrapper with the default Width and 
// FixSentenceEndings values replaced by the user-provided values.
// One option corresponds to each exported TextWrapper field.
type option func(*TextWrapper)

// A call to Width may be passed to NewTextWrapper or any wrapping
// function to override the default line width (70).
func Width(i int) option {
    return func(t *TextWrapper) {
        t.Width = i
    }
}

// A call to ExpandTabs may be passed to NewTextWrapper or any
// wrapping function to override the default value (true).  Passing
// ExpandTabs to Shorten has no effect, as each sequence of
// whitespace is replaced by a single space before wrapping.
func ExpandTabs(b bool) option {
    return func(t *TextWrapper) {
        t.ExpandTabs = b
    }
}

// A call to TabSize may be passed to NewTextWrapper or any wrapping
// function to override the default tab size (8).  Like ExpandTabs,
// passing TabSize to Shorten has no effect because each sequence of
// whitespace is replaced by a single space before wrapping.
func TabSize(i int) option {
    return func(t *TextWrapper) {
        t.TabSize = i
    }
}

// A call to ReplaceWhitespace may be passed to NewTextWrapper or any
// wrapping function to override the default value (true).  Passing
// ReplaceWhitespace to Shorten has no effect, as each sequence of
// whitespace is replaced by a single space before wrapping.
func ReplaceWhitespace(b bool) option {
    return func(t *TextWrapper) {
        t.ReplaceWhitespace = b
    }
}

// A call to DropWhitespace may be passed to NewTextWrapper or any
// wrapping function to override the default value (true).
func DropWhitespace(b bool) option {
    return func(t *TextWrapper) {
        t.DropWhitespace = b
    }
}

// A call to InitialIndent may be passed to NewTextWrapper or any
// wrapping function to override the default initial indent ("").
func InitialIndent(s string) option {
    return func(t *TextWrapper) {
        t.InitialIndent = s
    }
}

// A call to SubsequentIndent may be passed to NewTextWrapper or any
// wrapping function to override the default subsequent indent ("").
func SubsequentIndent(s string) option {
    return func(t *TextWrapper) {
        t.SubsequentIndent = s
    }
}

// A call to FixSentenceEndings may be passed to NewTextWrapper or
// any wrapping function to override the default value (false).
func FixSentenceEndings(b bool) option {
    return func(t *TextWrapper) {
        t.FixSentenceEndings = b
    }
}

// A call to BreakLongWords may be passed to NewTextWrapper or any
// wrapping function to override the default value (true).
func BreakLongWords(b bool) option {
    return func(t *TextWrapper) {
        t.BreakLongWords = b
    }
}

// A call to BreakOnHyphens may be passed to NewTextWrapper or any
// wrapping function to override the default value (true).
func BreakOnHyphens(b bool) option {
    return func(t *TextWrapper) {
        t.BreakOnHyphens = b
    }
}

// A call to MaxLines may be passed to NewTextWrapper or any wrapping
// function to specify a maximum number of lines to wrap.  The
// default value is zero, meaning that there is no maximum.  Passing
// MaxLines to Shorten has no effect, as Shorten always returns a
// single line.
func MaxLines(i int) option {
    return func(t *TextWrapper) {
        t.MaxLines = i
    }
}

// A call to Placeholder may be passed to NewTextWrapper or any
// wrapping function to override the default placeholder (" [...]").
func Placeholder(s string) option {
    return func(t *TextWrapper) {
        t.Placeholder = s
    }
}
