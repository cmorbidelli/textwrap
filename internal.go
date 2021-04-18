package textwrap

import "strings"

var strip = strings.TrimSpace

func lStrip(s string) string {
    return strings.TrimLeft(s, "\t\n\v\f\r ")
}

func rStrip(s string) string {
    return strings.TrimRight(s, "\t\n\v\f\r ")
}

func isSpace(s string) bool {
    return strip(s) == ""
}

type line struct {
    chunks []string
    length int
}

func (l *line) push(chunk string) {
    l.length += len(chunk)
    l.chunks = append(l.chunks, chunk)
}

func (l *line) pop() {
    last := len(l.chunks) - 1
    l.length -= len(l.chunks[last])
    l.chunks = l.chunks[:last]
}

func (l *line) str() string {
    return strings.Join(l.chunks, "")
}
