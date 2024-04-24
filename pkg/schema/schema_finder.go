package schema

import (
	"reflect"
	"slices"
	"strings"

	sdinsurelogger "github.com/sdinsure/agent/pkg/logger"

	apppb "github.com/footprintai/restcol/api/pb"
	encoding "github.com/footprintai/restcol/pkg/encoding"
	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
)

type SchemaBuilder struct {
	log sdinsurelogger.Logger
}

func NewSchemaBuilder(log sdinsurelogger.Logger) *SchemaBuilder {
	return &SchemaBuilder{log: log}
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
		s.log,
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

func TraverseMap(log sdinsurelogger.Logger, m map[string]interface{}, prefixes []string, callback CallbackFunc) error {
	for key := range m {
		valueOf := reflect.ValueOf(m[key])
		switch valueOf.Kind() {
		case reflect.Invalid:
			log.Warn("schemafinder: skipped field: %s as the value is invalid\n", key)
			continue

		case reflect.Map:
			if err := TraverseMap(log, m[key].(map[string]interface{}), append(prefixes, key), callback); err != nil {
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
