package valueobject

type Span struct{ Start, End int } // [Start, End) in runes (Unicode code points)

// Slice returns the part of text covered by s, or false if s is empty or out of range.
func (s Span) Slice(text string) (string, bool) {
	runes := []rune(text)
	if s.Start < 0 || s.End <= s.Start || s.End > len(runes) {
		return "", false
	}
	return string(runes[s.Start:s.End]), true
}
