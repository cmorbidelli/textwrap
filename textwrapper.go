// Package textwrap is an implementation of Python's textwrap library.
// textwrap exports functions for wrapping, indenting, and dedenting
// text.  The wrapping functions' behaviors can be customized in a
// number of ways, such as dropping leading and trailing whitespace,
// expanding tabs, or only wrapping a certain number of lines.  For
// use cases demanding a large number of wrapping operations, the
// package exports the TextWrapper struct and its methods.
//
// A TextWrapper is initiated with default values, so functions that
// use a TextWrapper also accept optional arguments to override the
// defaults.  For example:
//
//     Wrap("The quick brown fox jumped over the lazy jog",
//          Width(10), ExpandTabs(false))
//
// will wrap the text to a line width of 10 without expanding tabs,
// but will otherwise use the default behavior for a TextWrapper.
// One "option" corresponds to each of TextWrapper's exported fields.
package textwrap

import (
	"regexp"
	"strings"
)

// The following variables define whitespace and regexps that are
// used throughout the package.  Although it is not recommended,
// other values can be substituted to change the functions'
// behaviors.  This may be useful, for instance, when dealing with
// different character sets.  Note that, while Regexp structs are
// safe for concurrent use by multiple goroutines, no effort has been
// taken to make the other global variables concurrency-safe.
// These variables have been exported in an effort to make textwrap
// more generally useful, but I have not verified that any of the
// functions will make sense in the context of non-Latin character
// sets. Anyone wishing to modify the global variables should do so
// carefully and verify that they produce the desired results.
var (
	Space           = " "
	Tab             = "\t"
	Newline         = "\n"
	OtherWhitespace = Tab + Newline + "\v\f\r"
	Whitespace      = Space + OtherWhitespace

	// WhitespaceRe matches any whitespace character except Space.
	// It is used to replace characters with spaces if
	// ReplaceWhitespace is true.
	WhitespaceRe = regexp.MustCompile("[" + OtherWhitespace + "]")
	// SentenceEndingRe matches any non-whitespace character, followed
	// by a sentence-ending punctuation mark and at least one space
	// It is only used if FixSentenceEndings is true.
	SentenceEndingRe = regexp.MustCompile("([^" + Whitespace + "]" +
		"[.!?]['\"]?) [ ]*")
	// ChunksHyphenRe is used to break text into chunks for wrapping if
	// BreakOnHyphens is true.
	ChunksHyphenRe = regexp.MustCompile("(\u2014|[^" + Whitespace + "]+-|" +
		"[^" + Whitespace +
		"]+|[" + Whitespace + "]+)")
	// ChunksNoHyphenRe is used if BreakOnHyphens is false.
	ChunksNoHyphenRe = regexp.MustCompile("(\u2014|[^" + Whitespace + "]+|" +
		"[" + Whitespace + "]+)")
	// ConsWhitespaceRe is used by Shorten to find consecutive
	// whitespace characters.
	ConsWhitespaceRe = regexp.MustCompile("[" + Whitespace + "]+")
	// LeadWhitespaceRe is used by Dedent to find leading whitespace.
	LeadWhitespaceRe = regexp.MustCompile("^[" + Whitespace + "]*")
)

// bonus functions to simplify stripping whitespace from chunks of text
func Strip(s string) string {
	return strings.Trim(s, Whitespace)
}

func Lstrip(s string) string {
	return strings.TrimLeft(s, Whitespace)
}

func Rstrip(s string) string {
	return strings.TrimRight(s, Whitespace)
}

// TextWrapper contains values that govern wrapping behavior.
type TextWrapper struct {
	// Line width.  Default value is 70.
	Width int

	// If true, tab characters are replaced with TabSize number of
	// spaces before wrapping.  Default value is true.
	ExpandTabs bool

	// Tab size in spaces.  Default value is 8 spaces.
	TabSize int

	// If true, each whitespace character is replaced with a space.
	// Whitespace characters are defined as: '\t', '\n', '\v' '\f',
	// '\r', and ' '.  Default value is true.
	ReplaceWhitespace bool

	// If true, leading and trailing whitespace is dropped from each
	// line, after wrapping but before indenting.  Whitespace-only
	// lines are dropped.  Ignores whitespace at the beginning of
	// text if followed by non-whitespace.  Default value is true.
	DropWhitespace bool

	// Prepended to the first line of text.  Default value is "".
	InitialIndent string

	// Prepended to lines following the first.  Default value is "".
	SubsequentIndent string

	// Attempts to place two spaces after the end of each sentence
	// using the sentenceEnding regexp.  Unfortunately, the regexp
	// can't currently distinguish punctuation within a sentence from
	// sentence endings, so (for instance) it will also match with
	// "Mr. Rogers".  Default is false.
	FixSentenceEndings bool

	// If true, words too long to fit on the line will be broken.
	// Otherwise, they will be placed on a separate line.  Default
	// value is true.
	BreakLongWords bool

	// Allows wrapping to occur on hyphens.  Default value is true.
	BreakOnHyphens bool

	// If greater than zero, will limit output to MaxLines lines.
	// Default value is 0.
	MaxLines int

	// If MaxLines is greater than zero and the text has to be
	// truncated, the last line will end with Placeholder.  Default
	// value is " [...]".
	Placeholder string
}

// NewTextWrapper returns a TextWrapper struct. Each field receives a
// default value unless the user provides a value in the form of an
// "option."  For example:
//
//     NewTextWrapper(Width(75), TabSize(4))
//
// sets t.Width and t.TabSize, but otherwise keeps the default values.
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
	}

	for _, opt := range opts {
		opt(&t)
	}

	return t
}

// Wrap splits text into lines of specified length.  The TextWrapper
// object contains fields that can be modified to control Wrap's
// behavior.  See TextWrapper for descriptions of the fields.
func (t *TextWrapper) Wrap(text string) []string {
	// First, Wrap checks if the values of TextWrapper's fields
	// make it impossible to wrap the text.  This can occur if:
	// (1) the line width is less than 1;
	if t.Width < 1 {
		panic("Width must be at least 1.")

		// (2) ExpandTabs is true and TabSize is less than zero; or
	} else if t.ExpandTabs && t.TabSize < 0 {
		panic("Tab size must be at least 0 to expand tabs.")
	}

	// (3) MaxLines is positive, but the last line is not wide enough
	//     to hold both the indent and the placeholder.
	if t.MaxLines > 0 {
		indent := t.SubsequentIndent
		if t.MaxLines == 1 {
			indent = t.InitialIndent
		}

		if len(indent)+len(Lstrip(t.Placeholder)) > t.Width {
			panic("Placeholder is too wide to fit on indented line.")
		}
	}
	// If one of these conditions is met, Wrap panics instead  of
	// restoring the default values because it is difficult to infer
	// the user's intent and simpler to assume that a mistake occured.

	// expands tabs if ExpandTabs is true
	if t.ExpandTabs {
		tabString := strings.Repeat(Space, t.TabSize)
		text = strings.Replace(text, Tab, tabString, -1)
	}

	// replaces whitespace if ReplaceWhitespace is true
	if t.ReplaceWhitespace {
		text = WhitespaceRe.ReplaceAllString(text, Space)
	}

	// attempts to fix sentence endings if FixSentenceEndings is true
	if t.FixSentenceEndings {
		text = SentenceEndingRe.ReplaceAllString(text, "${1}"+Space+Space)
	}

	// breaks text into chunks depending on BreakOnHyphens
	var chunks []string
	if t.BreakOnHyphens {
		chunks = ChunksHyphenRe.FindAllString(text, -1)
	} else {
		chunks = ChunksNoHyphenRe.FindAllString(text, -1)
	}

	// iterates through lines
	var lines []string
	for i := 0; i < len(chunks); i++ {
		// drops leading whitespace if DropWhitespace is true
		if len(lines) > 0 && t.DropWhitespace && Strip(chunks[i]) == "" {
			i++
		}

		// selects appropriate indent
		var indent string
		if len(lines) > 0 {
			indent = t.SubsequentIndent
		} else {
			indent = t.InitialIndent
		}

		// sets line width to allow room for indent and placeholder
		width := t.Width - len([]rune(indent))
		if t.MaxLines > 0 && len(lines) == t.MaxLines-1 {
			width -= len(t.Placeholder)
		}

		// appends chunks to current line until the next chunk would
		// exceed width, or text ends
		var curLen int
		var curLine []string
		for ; i < len(chunks); i++ {
			if curLen+len([]rune(chunks[i])) >= width {
				i--
				break
			}
			curLine = append(curLine, chunks[i])
			curLen += len([]rune(chunks[i]))
		}

		// peeks ahead to check if next chunk will need to be split
		// or placed on its own line
		if i+1 < len(chunks) && len([]rune(chunks[i+1])) > width {
			// if BreakLongWords is true, appends as much of the
			// chunk as possible to the current line, and leaves any
			// remainder for the next line
			if t.BreakLongWords {
				c := []rune(chunks[i+1])
				spaceLeft := 1
				if width >= 1 {
					spaceLeft = width - curLen
				}
				curLine = append(curLine, string(c[:spaceLeft]))
				curLen += spaceLeft
				chunks[i+1] = string(c[spaceLeft:])
				// or, if current line is empty, the chunk is appended
			} else if curLen == 0 && len(lines) != t.MaxLines-1 {
				i++
				curLine = append(curLine, chunks[i])
				curLen += len([]rune(chunks[i]))
			}
		}

		// if DropWhitespace is true, drops any trailing whitespace
		if last := len(curLine) - 1; t.DropWhitespace &&
			curLen > 0 && Strip(curLine[last]) == "" {
			curLen -= len([]rune(curLine[last]))
			curLine = curLine[:last]
		}

		// if the current line is MaxLine, applies any placeholder
		// and indent, appends the current line to lines, and exits
		// the main loop
		if t.MaxLines > 0 && len(lines) == t.MaxLines-1 {
			// if the line is empty, removes any leading whitespace
			// from the placeholder
			if curLen == 0 {
				curLine = append(curLine, Lstrip(t.Placeholder))
			} else {
				curLine = append(curLine, t.Placeholder)
			}
			lines = append(lines, indent+strings.Join(curLine, ""))
			break
			// or, if the current line is not empty, applies any indent
			// and appends the current line to lines
		} else if curLen > 0 {
			lines = append(lines, indent+strings.Join(curLine, ""))
		}
	}

	return lines
}

// Fill wraps the text and returns a single string consisting of
// the newline-separated lines.  The TextWrapper object contains
// fields that can be modified to control Wrap's behavior.  See
// TextWrapper for descriptions of the fields.
func (t *TextWrapper) Fill(text string) string {
	return strings.Join(t.Wrap(text), Newline)
}
