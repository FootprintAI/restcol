package schema

import (
	"reflect"
	"slices"
	"strings"

	apppb "github.com/footprintai/restcol/api/pb"

	encoding "github.com/footprintai/restcol/pkg/encoding"
	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
)

//
//import "context"
//
//type SchemaFinder struct {
//	collectionCURD *collectionsstorage.CollectionCURD
//}
//
//func NewSchemaFinder(collectionCURD *collectionsstorage.CollectionCURD) *SchemaFinder {
//	return &SchemaFinder{
//		collectionCURD: collectionCURD,
//	}
//}
//
//func (s *SchemaFinder) GetCollection(ctx context.Context)

type SchemaBuilder struct{}

func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{}
}

func (s *SchemaBuilder) Parse(rawByte []byte) (apppb.DataFormat, *appmodelcollections.ModelSchema, error) {
	dec, err := encoding.GetDecoder(apppb.DataFormat_DATA_FORMAT_AUTO)
	if err != nil {
		return apppb.DataFormat_DATA_FORMAT_UNKNOWN, nil, err
	}
	placeHolder := make(map[string]interface{})
	format, err := dec.Decode(rawByte, &placeHolder)
	_ = format
	if err != nil {
		return apppb.DataFormat_DATA_FORMAT_UNKNOWN, nil, err
	}

	dotNotationMap := make(map[string]interface{})
	if err := TraverseMap(
		placeHolder,
		nil,
		func(prefixes []string, current string, val any) error {
			path := strings.Join(append(prefixes, current), ".")
			dotNotationMap[path] = val
			return nil
		},
	); err != nil {
		return apppb.DataFormat_DATA_FORMAT_UNKNOWN, nil, err
	}

	var sortedKeys []string
	for k := range dotNotationMap {
		sortedKeys = append(sortedKeys, k)
	}
	slices.Sort(sortedKeys)

	var fields []appmodelcollections.ModelFieldSchema
	for _, k := range sortedKeys {
		val := dotNotationMap[k]
		swagVal, err := appmodelcollections.NewSwagValue(val)
		if err != nil {
			return apppb.DataFormat_DATA_FORMAT_UNKNOWN, nil, err
		}

		fields = append(fields, appmodelcollections.ModelFieldSchema{
			FieldName:      k,
			FieldValueType: swagVal.Type(),
			FieldExample:   swagVal,
		})
	}
	return format, &appmodelcollections.ModelSchema{
		Fields: fields,
	}, nil

}

type CallbackFunc func(prefixes []string, current string, value interface{}) error

func TraverseMap(m map[string]interface{}, prefixes []string, callback CallbackFunc) error {
	for key := range m {
		valueOf := reflect.ValueOf(m[key])
		switch valueOf.Kind() {
		case reflect.Map:
			if err := TraverseMap(m[key].(map[string]interface{}), append(prefixes, key), callback); err != nil {
				return err
			}
		default:
			if err := callback(prefixes, key, m[key]); err != nil {
				return err
			}
		}
	}
	return nil
}
