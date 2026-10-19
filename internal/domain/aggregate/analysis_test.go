package aggregate

import (
	"errors"
	"reflect"
	"testing"

	"github.com/KrisNathan/linguascreen-api-go/internal/domain/valueobject"
)

func span(start, end int) valueobject.Span { return valueobject.Span{Start: start, End: end} }

func word(surface, dictForm, translation string) Token {
	return Token{Surface: surface, Saveable: true, DictForm: dictForm, Reading: "r", Translation: translation}
}

func align(t *testing.T, text string, tokens []Token) ([]Token, error) {
	t.Helper()
	s, err := valueobject.NewSentence(text)
	if err != nil {
		t.Fatalf("NewSentence(%q): %v", text, err)
	}
	a := &Analysis{Tokens: tokens}
	err = a.Align(*s)
	return a.Tokens, err
}

func TestAlignSpans(t *testing.T) {
	got, err := align(t, "よく本を読みます。", []Token{
		word("よく", "よく", "often"),
		word("本", "本", "book"),
		{Surface: "を"},
		word("読みます", "読む", "to read"),
		{Surface: "。"},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []valueobject.Span{span(0, 2), span(2, 3), span(3, 4), span(4, 8), span(8, 9)}
	for i, tok := range got {
		if tok.Span != want[i] {
			t.Errorf("token %d %q: span %v, want %v", i, tok.Surface, tok.Span, want[i])
		}
	}
}

func TestAlignFillsGaps(t *testing.T) {
	got, err := align(t, "I read books.", []Token{
		word("I", "I", "I"),
		word("read", "read", "to read"),
		word("books", "book", "books"),
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []Token{
		{Surface: "I", Span: span(0, 1), Saveable: true, DictForm: "I", Reading: "r", Translation: "I"},
		{Surface: " ", Span: span(1, 2)},
		{Surface: "read", Span: span(2, 6), Saveable: true, DictForm: "read", Reading: "r", Translation: "to read"},
		{Surface: " ", Span: span(6, 7)},
		{Surface: "books", Span: span(7, 12), Saveable: true, DictForm: "book", Reading: "r", Translation: "books"},
		{Surface: ".", Span: span(12, 13)},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

func TestAlignAlteredSurface(t *testing.T) {
	_, err := align(t, "走った", []Token{word("走る", "走る", "to run")})
	if !errors.Is(err, ErrMisaligned) {
		t.Errorf("err = %v, want ErrMisaligned", err)
	}
}

func TestAlignCountsRunes(t *testing.T) {
	got, err := align(t, "😀本", []Token{{Surface: "😀"}, word("本", "本", "book")})
	if err != nil {
		t.Fatal(err)
	}
	if got[0].Span != span(0, 1) || got[1].Span != span(1, 2) {
		t.Errorf("spans %v %v, want [0,1) [1,2)", got[0].Span, got[1].Span)
	}
}

func TestAlignDropsEmptyAndDemotes(t *testing.T) {
	got, err := align(t, "本を", []Token{
		{Surface: ""},
		word("本", "本", ""), // missing translation
		{Surface: "を", Saveable: false, DictForm: "を", Parts: []Part{{Surface: "を"}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []Token{{Surface: "本", Span: span(0, 1)}, {Surface: "を", Span: span(1, 2)}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

func TestAlignParts(t *testing.T) {
	station := word("火车站", "火车站", "train station")
	station.Parts = []Part{{Surface: "火车", DictForm: "火车"}, {Surface: ""}, {Surface: "站", DictForm: "站"}}
	got, err := align(t, "去火车站", []Token{word("去", "去", "to go"), station})
	if err != nil {
		t.Fatal(err)
	}
	parts := got[1].Parts
	if len(parts) != 2 || parts[0].Span != span(1, 3) || parts[1].Span != span(3, 4) {
		t.Errorf("parts %+v, want 火车 [1,3) and 站 [3,4)", parts)
	}
}

func TestAlignDropsBadParts(t *testing.T) {
	station := word("火车站", "火车站", "train station")
	station.Parts = []Part{{Surface: "火车"}, {Surface: "车站"}}
	got, err := align(t, "火车站", []Token{station})
	if err != nil {
		t.Fatal(err)
	}
	if got[0].Parts != nil {
		t.Errorf("parts %+v, want nil", got[0].Parts)
	}
}
