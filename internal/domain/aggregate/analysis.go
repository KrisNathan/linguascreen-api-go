package aggregate

import (
	"errors"
	"strings"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

// ErrMisaligned means the LLM altered the sentence, so a token's surface can't be found in it.
var ErrMisaligned = errors.New("analysis does not match sentence")

type Analysis struct {
	Translation      string
	LanguageMismatch bool
	Tokens           []Token
}

func (a *Analysis) Align(sentence valueobject.Sentence) error {
	text := sentence.Text
	var out []Token
	byteAt, unitAt := 0, 0 // cursor in bytes and in UTF-16 code units

	addGap := func(gap string) {
		if gap == "" {
			return
		}
		end := unitAt + valueobject.UTF16Len(gap)
		out = append(out, Token{Surface: gap, Span: valueobject.Span{Start: unitAt, End: end}})
		byteAt += len(gap)
		unitAt = end
	}

	for _, t := range a.Tokens {
		if t.Surface == "" {
			continue
		}
		i := strings.Index(text[byteAt:], t.Surface)
		if i < 0 {
			return ErrMisaligned
		}
		addGap(text[byteAt : byteAt+i])

		t.Span = valueobject.Span{Start: unitAt, End: unitAt + valueobject.UTF16Len(t.Surface)}
		byteAt += len(t.Surface)
		unitAt = t.Span.End

		// reading may be empty on purpose, for languages without a separate reading system.
		if t.Saveable && (t.DictForm == "" || t.Translation == "") {
			t.Saveable = false
		}
		if t.Saveable {
			t.Parts = alignParts(t)
		} else {
			t = Token{Surface: t.Surface, Span: t.Span}
		}
		out = append(out, t)
	}
	addGap(text[byteAt:])

	a.Tokens = out
	return nil
}

// alignParts returns the token's parts with absolute spans, or nil if they
// don't join up to the token's surface exactly.
func alignParts(t Token) []Part {
	var parts []Part
	var joined strings.Builder
	for _, p := range t.Parts {
		if p.Surface != "" {
			parts = append(parts, p)
			joined.WriteString(p.Surface)
		}
	}
	if len(parts) == 0 || joined.String() != t.Surface {
		return nil
	}

	at := t.Span.Start
	for i := range parts {
		end := at + valueobject.UTF16Len(parts[i].Surface)
		parts[i].Span = valueobject.Span{Start: at, End: end}
		at = end
	}
	return parts
}
