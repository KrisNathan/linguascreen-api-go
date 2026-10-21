package valueobject

import "testing"

func TestSpanSlice(t *testing.T) {
	text := "よく本を読みます😀。"
	tests := []struct {
		span Span
		want string
		ok   bool
	}{
		{Span{4, 8}, "読みます", true},
		{Span{0, 1}, "よ", true},
		{Span{8, 9}, "😀", true},
		{Span{8, 10}, "😀。", true},
		{Span{3, 3}, "", false},
		{Span{5, 4}, "", false},
		{Span{-1, 2}, "", false},
		{Span{9, 11}, "", false},
	}
	for _, tt := range tests {
		got, ok := tt.span.Slice(text)
		if got != tt.want || ok != tt.ok {
			t.Errorf("%v.Slice() = %q, %v, want %q, %v", tt.span, got, ok, tt.want, tt.ok)
		}
	}
}
