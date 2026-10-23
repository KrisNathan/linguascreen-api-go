package valueobject

type Token struct {
	Surface                        string
	Span                           Span // zero until aligned
	Saveable                       bool
	DictForm, Reading, Translation string
	Parts                          []Part
}
