package dotnotation

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDotNotation(t *testing.T) {
	for _, testcase := range []struct {
		parts []string
		repr  string
	}{
		{
			parts: []string{"a", "b", "c"},
			repr:  "a.b.c",
		},
		{
			parts: []string{"a"},
			repr:  "a",
		},
	} {
		dot := New(testcase.parts...)
		assert.EqualValues(t, testcase.repr, dot.String())

		gotDot := Parse(testcase.repr)
		assert.EqualValues(t, testcase.parts, gotDot.Parts)
	}
}

func TestDotNotationSort(t *testing.T) {
	notations := []*DotNotation{
		New("a", "b", "c"),
		New("a"),
		New("a", "b", "c1"),
		New("a", "b1", "c"),
		New("a", "b2", "c"),
	}

	sort.Slice(notations, func(i, j int) bool {
		return notations[i].Less(notations[j])
	})

	assert.EqualValues(t, "a", notations[0].String())
	assert.EqualValues(t, "a.b.c", notations[1].String())
	assert.EqualValues(t, "a.b.c1", notations[2].String())
	assert.EqualValues(t, "a.b1.c", notations[3].String())
	assert.EqualValues(t, "a.b2.c", notations[4].String())
}
