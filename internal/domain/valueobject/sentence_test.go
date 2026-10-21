package valueobject

import "testing"

func normalizationTest(t *testing.T, raw, expected string) {
	s, err := NewSentence(raw)
	if err != nil {
		t.Errorf("failed to create new sentence '%v'", raw)
	}
	if s.Text != expected {
		t.Errorf("normalization failed expected '%v' got '%v'", expected, s.Text)
	}
}

func TestNormalizationJPNewline(t *testing.T) {
	normalizationTest(t, "よく本を\n読みます。", "よく本を読みます。")
}

func TestNormalizationENNewline(t *testing.T) {
	normalizationTest(t, "I read\nbooks", "I read books")
}

func TestNormalizationJPWhitespace(t *testing.T) {
	normalizationTest(t, "本 を 読む", "本を読む")
}

func TestHash(t *testing.T) {
	s, err := NewSentence("本を読む")
	if err != nil {
		t.Error("failed to create new sentence '本を読む'")
	}

	hashed := s.Hash()

	if hashed != "1ab62b222e45ffbce262890c4c4e0886cb77195321c22e3aafbf7564a9ec1b23" {
		t.Error("invalid sentence hash")
	}
}
