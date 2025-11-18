package models

// Bbox represents a bounding box with coordinates
type Bbox struct {
	X0 float64 `json:"x0"`
	Y0 float64 `json:"y0"`
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
}

// Symbol represents a single character/symbol in OCR
type Symbol struct {
	Bbox          Bbox    `json:"bbox"`
	Text          string  `json:"text"`
	Confidence    float64 `json:"confidence"`
	IsSuperscript int     `json:"is_superscript"`
	IsSubscript   int     `json:"is_subscript"`
	IsDropcap     int     `json:"is_dropcap"`
}

// Choice represents an alternative recognition for a word
type Choice struct {
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
}

// Word represents a recognized word in OCR
type Word struct {
	Bbox       Bbox     `json:"bbox"`
	Text       string   `json:"text"`
	Confidence float64  `json:"confidence"`
	Choices    []Choice `json:"choices"`
	FontName   string   `json:"font_name"`
	Symbols    []Symbol `json:"symbols"`
}

// RowAttributes represents attributes of a text row
type RowAttributes struct {
	RowHeight  float64 `json:"rowHeight"`
	Descenders float64 `json:"descenders"`
	Ascenders  float64 `json:"ascenders"`
}

// Baseline represents the baseline of a text line
type Baseline struct {
	X0          float64 `json:"x0"`
	Y0          float64 `json:"y0"`
	X1          float64 `json:"x1"`
	Y1          float64 `json:"y1"`
	HasBaseline bool    `json:"has_baseline"`
}

// Line represents a line of text in OCR
type Line struct {
	Bbox          Bbox          `json:"bbox"`
	Text          string        `json:"text"`
	Confidence    float64       `json:"confidence"`
	RowAttributes RowAttributes `json:"rowAttributes"`
	Baseline      Baseline      `json:"baseline"`
	Words         []Word        `json:"words"`
}

// Paragraph represents a paragraph of text in OCR
type Paragraph struct {
	Bbox       Bbox    `json:"bbox"`
	Text       string  `json:"text"`
	Confidence float64 `json:"confidence"`
	IsLtr      int     `json:"is_ltr"`
	Lines      []Line  `json:"lines"`
}

// Block represents a block of text in OCR results
type Block struct {
	Bbox       Bbox        `json:"bbox"`
	Text       string      `json:"text"`
	Confidence float64     `json:"confidence"`
	Blocktype  int         `json:"blocktype"`
	Paragraphs []Paragraph `json:"paragraphs"`
}

// Page represents the full OCR result page
type Page struct {
	Blocks        []Block `json:"blocks"`
	Confidence    float64 `json:"confidence"`
	Oem           string  `json:"oem"`
	Osd           string  `json:"osd"`
	Psm           string  `json:"psm"`
	Text          string  `json:"text"`
	Version       string  `json:"version"`
	Hocr          string  `json:"hocr"`
	Tsv           string  `json:"tsv"`
	Box           string  `json:"box"`
	Unlv          string  `json:"unlv"`
	Sd            string  `json:"sd"`
	ImageColor    string  `json:"imageColor"`
	ImageGrey     string  `json:"imageGrey"`
	ImageBinary   string  `json:"imageBinary"`
	RotateRadians float64 `json:"rotateRadians"`
	Pdf           []int   `json:"pdf"`
	Debug         string  `json:"debug"`
}
