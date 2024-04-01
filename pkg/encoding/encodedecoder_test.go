package encoding

import (
	"testing"

	pb "github.com/footprintai/restcol/api/pb/proto"
	"github.com/stretchr/testify/assert"
)

func TestJsonDecoder(t *testing.T) {
	m := make(map[string]interface{})
	d := jsonDecoder{}
	f, e := d.Decode([]byte(testJsonData), &m)
	assert.NoError(t, e)
	assert.EqualValues(t, f, pb.DataFormat_DATA_FORMAT_JSON)
	assert.EqualValues(t, m["foo"], "bar")
	assert.EqualValues(t, m["foofloat"], 3.14)
}

func TestCsvDecoder(t *testing.T) {
	m := [][]string{}
	d := csvDecoder{}
	f, e := d.Decode([]byte(testCsvData), &m)
	assert.NoError(t, e)
	assert.EqualValues(t, f, pb.DataFormat_DATA_FORMAT_CSV)
	assert.EqualValues(t, m[0][0], "foo")
	assert.EqualValues(t, m[0][1], "foofloat")
	assert.EqualValues(t, m[1][0], "bar")
	assert.EqualValues(t, m[1][1], "3.14")
}

func TestXmlDecoder(t *testing.T) {
	m := make(map[string]interface{})
	d := xmlDecoder{}
	f, e := d.Decode([]byte(testXmlData), &m)
	assert.NoError(t, e)
	assert.EqualValues(t, f, pb.DataFormat_DATA_FORMAT_XML)
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
