package dotnotation

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
)

type DotNotation struct {
	Parts []string `json:"parts"`
}

var (
	_ sql.Scanner   = &DotNotation{}
	_ driver.Valuer = DotNotation{}
)

func (d *DotNotation) Scan(in any) error {
	if _, isStr := in.(string); !isStr {
		return errors.New("dotnotation: expect str but not")
	}
	dCopied := Parse(in.(string))
	(*d) = (*dCopied)
	return nil
}

func (d DotNotation) Value() (driver.Value, error) {
	return d.String(), nil
}

func New(parts ...string) *DotNotation {
	return &DotNotation{
		Parts: parts,
	}
}

func (d *DotNotation) String() string {
	return strings.Join(d.Parts, ".")
}

func (d *DotNotation) AddSuffix(suffixs ...string) *DotNotation {
	return &DotNotation{
		Parts: append(d.Parts, suffixs...),
	}
}

func (d *DotNotation) AddPrefix(prefixs ...string) *DotNotation {
	return &DotNotation{
		Parts: append(prefixs, d.Parts...),
	}
}

func (d *DotNotation) Less(d2 *DotNotation) bool {
	if len(d.Parts) > len(d2.Parts) {
		return false
	}
	if len(d.Parts) < len(d2.Parts) {
		return true
	}
	for partIndex := 0; partIndex < len(d.Parts); partIndex++ {
		if d.Parts[partIndex] > d2.Parts[partIndex] {
			return false
		}
		if d.Parts[partIndex] < d2.Parts[partIndex] {
			return true
		}
	}
	return false
}

func Parse(s string) *DotNotation {
	tokens := strings.Split(s, ".")
	return New(tokens...)
}
