package textwrap

// options allow the user to override a TextWrapper's default values.
// The function calls can be passed as arguments to NewTextWrapper, 
// Wrap, Fill, and Shorten.  For example:
//     Fill(myText, Width(50), FixSentenceEndings(true))
// fills myText using a TextWrapper with the default Width and 
// FixSentenceEndings values replaced by the user-provided values.
// One option corresponds to each exported TextWrapper field.
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
