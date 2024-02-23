package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJsonDecoder(t *testing.T) {
	m := make(map[string]interface{})
	d := jsonDecoder{}
	assert.NoError(t, d.Decode([]byte(testJsonData), &m))
	assert.EqualValues(t, m["foo"], "bar")
	assert.EqualValues(t, m["foofloat"], 3.14)
}

func TestCsvDecoder(t *testing.T) {
	m := [][]string{}
	d := csvDecoder{}
	assert.NoError(t, d.Decode([]byte(testCsvData), &m))
	assert.EqualValues(t, m[0][0], "foo")
	assert.EqualValues(t, m[0][1], "foofloat")
	assert.EqualValues(t, m[1][0], "bar")
	assert.EqualValues(t, m[1][1], "3.14")
}

func TestXmlDecoder(t *testing.T) {
	m := make(map[string]interface{})
	d := xmlDecoder{}
	assert.NoError(t, d.Decode([]byte(testXmlData), &m))
	assert.NotNil(t, m["document"])
	md := m["document"].(map[string]interface{})
	assert.NotNil(t, md)
	assert.EqualValues(t, md["foo"], "bar")
	assert.EqualValues(t, md["foofloat"], 3.14)
}

var (
	testJsonData = `{"foo": "bar", "foofloat": 3.14}`
	testCsvData  = `foo,foofloat
bar,3.14`
	testXmlData = `<?xml version="1.0" encoding="UTF-8"?>
<document>
<foo>bar</foo>
<foofloat>3.14</foofloat>
</document>`
)
