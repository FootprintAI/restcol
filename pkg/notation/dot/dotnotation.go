package dotnotation

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"strings"
)

type DotNotation struct {
	parts []string
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
		parts: parts,
	}
}

func (d *DotNotation) Parts() []string {
	return d.parts
}

func (d *DotNotation) String() string {
	return strings.ToLower(strings.Join(d.parts, "."))
}

func (d *DotNotation) AddSuffix(suffixs ...string) *DotNotation {
	return &DotNotation{
		parts: append(d.parts, suffixs...),
	}
}

func (d *DotNotation) AddPrefix(prefixs ...string) *DotNotation {
	return &DotNotation{
		parts: append(prefixs, d.parts...),
	}
}

func (d *DotNotation) Less(d2 *DotNotation) bool {
	if len(d.parts) > len(d2.parts) {
		return false
	}
	if len(d.parts) < len(d2.parts) {
		return true
	}
	for partIndex := 0; partIndex < len(d.parts); partIndex++ {
		if d.parts[partIndex] > d2.parts[partIndex] {
			return false
		}
		if d.parts[partIndex] < d2.parts[partIndex] {
			return true
		}
	}
	return false
}

func Parse(s string) *DotNotation {
	tokens := strings.Split(s, ".")
	return New(tokens...)
}
