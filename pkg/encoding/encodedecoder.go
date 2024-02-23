package encoding

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/clbanning/mxj/v2"

	pb "github.com/footprintai/restcol/api/pb/proto"
	sderrors "github.com/sdinsure/agent/pkg/errors"
)

type Decoder interface {
	Decode(data []byte, v any) error
}

var (
	errNotImpl   = sderrors.NewNotImplError(errors.New("not impl"))
	errBadParams = sderrors.NewNotImplError(errors.New("bad params"))
)

func GetDecoder(dataformat pb.DataFormat) (Decoder, *sderrors.Error) {
	switch dataformat {
	case pb.DataFormat_DATA_FORMAT_AUTO:
		return newAutoDecoder(), nil
	case pb.DataFormat_DATA_FORMAT_JSON:
		return jsonDecoder{}, nil
	case pb.DataFormat_DATA_FORMAT_CSV:
		return csvDecoder{}, nil
	case pb.DataFormat_DATA_FORMAT_XML:
		return xmlDecoder{}, nil
	case pb.DataFormat_DATA_FORMAT_URL:
		return notImpl{}, errNotImpl
	case pb.DataFormat_DATA_FORMAT_MEDIA:
		return notImpl{}, errNotImpl
	default:
		return notImpl{}, errNotImpl
	}
}

type notImpl struct{}

func (n notImpl) Decode(data []byte, v any) error {
	return fmt.Errorf("not impl")
}

type jsonDecoder struct{}

var (
	_ Decoder = jsonDecoder{}
)

func (j jsonDecoder) Decode(data []byte, v any) error {
	return json.NewDecoder(bytes.NewReader(data)).Decode(v)
}

type csvDecoder struct{}

var (
	_ Decoder = csvDecoder{}
)

func (j csvDecoder) Decode(in []byte, v any) error {
	// check dst type, should be [][]string
	ss, isSliceType := v.(*[][]string)
	if !isSliceType {
		return errors.New("invalid type: require *[][]string for csv")
	}
	rs, err := csv.NewReader(bytes.NewReader(in)).ReadAll()
	if err != nil {
		return err
	}
	*ss = append(*ss, rs...)
	return nil
}

type xmlDecoder struct{}

var (
	_ Decoder = xmlDecoder{}
)

func (x xmlDecoder) Decode(data []byte, v any) error {
	mp, isMapType := v.(*map[string]interface{})
	if !isMapType {
		return errors.New("invalid type: require *map[string]interface{} for xml")
	}

	m, err := mxj.NewMapXmlReader(bytes.NewReader(data), true /*casting the value to its real type, e.g. float*/)
	if err != nil {
		return err
	}
	(*mp) = m
	return nil
}

type chainDecoder struct {
	chain []Decoder
}

var (
	_ Decoder = chainDecoder{}
)

func (c chainDecoder) Decode(data []byte, v any) error {
	for _, decoder := range c.chain {
		if err := decoder.Decode(data, v); err == nil {
			return nil
		}
	}
	return sderrors.NewBadParamsError(errors.New("not support format"))
}

func newAutoDecoder() Decoder {
	return chainDecoder{
		chain: []Decoder{
			jsonDecoder{},
			csvDecoder{},
			xmlDecoder{},
		},
	}
}
