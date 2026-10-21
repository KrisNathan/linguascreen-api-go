package valueobject

import (
	"crypto/sha256"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/KrisNathan/linguascreen-api-go/internal/errs"
	"golang.org/x/text/unicode/norm"
)

const MaxSentenceLength = 300 // runes

type Sentence struct {
	Text string
}

func NewSentence(raw string) (*Sentence, error) {
	s := normalize(raw)
	if s == "" {
		return nil, errs.WithHint(errs.ErrInvalidInput, "sentence is empty")
	}

	if sentenceLen := utf8.RuneCountInString(s); sentenceLen > MaxSentenceLength {
		return nil, errs.WithHint(errs.ErrInvalidInput, fmt.Sprintf("sentence is %d chars long, max is %d", sentenceLen, MaxSentenceLength))
	}

	return &Sentence{Text: s}, nil
}

func (s *Sentence) Hash() string {
	hashBytes := sha256.Sum256([]byte(s.Text))
	hashString := fmt.Sprintf("%x", hashBytes)
	return hashString
}

func normalize(raw string) string {
	var b strings.Builder
	var prev rune // last rune written, 0 at start
	pendingSpace := false
	for _, r := range norm.NFC.String(raw) {
		switch {
		case unicode.IsSpace(r):
			pendingSpace = true
		case r == '\u200d': // zero-width joiner, keep for emoji
			b.WriteRune(r)
			prev = r
		case unicode.IsControl(r), unicode.Is(unicode.Cf, r):
			// dropped, no fallthrough
		default:
			// A whitespace run becomes one space, or nothing between two CJK
			// characters or at the start.
			if pendingSpace && prev != 0 && !(isCJK(prev) && isCJK(r)) {
				b.WriteByte(' ')
			}
			pendingSpace = false
			b.WriteRune(r)
			prev = r
		}
	}
	return b.String() // trailing whitespace is never written
}

func isCJK(r rune) bool {
	return unicode.Is(unicode.Han, r) ||
		(r >= 0x3000 && r <= 0x30FF) || // CJK punctuation, hiragana, katakana, ー
		(r >= 0xFF00 && r <= 0xFFEF) // fullwidth and halfwidth forms
}
