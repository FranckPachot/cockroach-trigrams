package trgm

import (
	"strings"

	"github.com/euskadi31/go-tokenizer"
	"github.com/gofrs/uuid"
	"github.com/kljensen/snowball"
)

func Stemmed(str string) []string {
	t := tokenizer.New()
	tokens := t.Tokenize(str)
	stemmed := make([]string, len(tokens))
	for i, token := range tokens {
		stemmed[i], _ = snowball.Stem(token, "english", true)
	}
	return stemmed
}

func analyze(str string) string {
	t := tokenizer.New()
	tokens := t.Tokenize(str)
	stemmed := make([]string, len(tokens))
	for i, token := range tokens {
		stemmed[i], _ = snowball.Stem(token, "english", true)
	}
	return strings.Join(stemmed, " ")
}

type Food struct {
	id uuid.UUID `json:"-"`

	FDCID int64  `json:"fdcId"`
	Name  string `json:"description"`
	Brand string `json:"brandName"`
	Type  string `json:"dataType"`
}

func (f *Food) ID() uuid.UUID {
	if f.id == uuid.Nil {
		f.id = uuid.Must(uuid.NewV4())
	}
	return f.id
}

func (f *Food) FullName() string {
	if f.Brand == "" {
		return f.Name
	}
	return f.Brand + " " + f.Name
}

func (f *Food) StemmedFullName() []string {
	return Stemmed(f.FullName())
}

func (f *Food) AnalyzedFullName() string {
	return strings.Join(Stemmed(f.FullName()), " ")
}

func (f *Food) Weight() float64 {
	if f.Type == "Foundation" {
		return 1.5
	}
	return 1.0
}
