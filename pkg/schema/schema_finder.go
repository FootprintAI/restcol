package schema

import (
	"reflect"
	"sort"

	sdinsurelogger "github.com/sdinsure/agent/pkg/logger"

	apppb "github.com/footprintai/restcol/api/pb"
	encoding "github.com/footprintai/restcol/pkg/encoding"
	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
	dotnotation "github.com/footprintai/restcol/pkg/notation/dot"
)

type SchemaBuilder struct {
	log sdinsurelogger.Logger
}

func NewSchemaBuilder(log sdinsurelogger.Logger) *SchemaBuilder {
	return &SchemaBuilder{log: log}
}

func (s *SchemaBuilder) Parse(rawByte []byte) (apppb.DataFormat, *appmodelcollections.ModelSchema, map[string]interface{}, error) {
	dec, err := encoding.GetDecoder(apppb.DataFormat_DATA_FORMAT_AUTO)
	if err != nil {
		return apppb.DataFormat_DATA_FORMAT_UNKNOWN, nil, nil, err
	}
	structHolder := make(map[string]interface{})
	format, err := dec.Decode(rawByte, &structHolder)
	if err != nil {
		return apppb.DataFormat_DATA_FORMAT_UNKNOWN, nil, nil, err
	}
	modelSchema, err := s.Flatten(structHolder)
	if err != nil {
		return apppb.DataFormat_DATA_FORMAT_UNKNOWN, nil, nil, err
	}
	return format, modelSchema, structHolder, nil
}

func (s *SchemaBuilder) Flatten(structHolder map[string]interface{}) (*appmodelcollections.ModelSchema, error) {
	dotNotationMap := make(map[*dotnotation.DotNotation]interface{})
	if err := TraverseMap(
		s.log,
		structHolder,
		nil,
		func(prefixes []string, current string, val any) error {
			nota := dotnotation.New(prefixes...).AddSuffix(current)
			dotNotationMap[nota] = val
			return nil
		},
	); err != nil {
		return nil, err
	}

	var sortedKeys []*dotnotation.DotNotation
	for k := range dotNotationMap {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i].Less(sortedKeys[j])
	})

	var fields []*appmodelcollections.ModelFieldSchema
	for _, k := range sortedKeys {
		val := dotNotationMap[k]
		swagVal, err := appmodelcollections.NewSwagValue(val)
		if err != nil {
			return nil, err
		}

		fields = append(fields, &appmodelcollections.ModelFieldSchema{
			FieldName:      k,
			FieldValueType: swagVal.Type(),
			FieldExample:   swagVal,
		})
	}
	return &appmodelcollections.ModelSchema{
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

func (s *SchemaBuilder) Equals(s1, s2 *appmodelcollections.ModelSchema) bool {

	if len(s1.Fields) != len(s2.Fields) {
		return false
	}

	// lookup is keyed by fieldname and each value is with FieldValueType
	lookup := make(map[string]appmodelcollections.SwagValueType)

	for _, field := range s1.Fields {
		lookup[field.FieldName.String()] = field.FieldValueType
	}

	for _, field := range s2.Fields {
		lookupValue, lookupExists := lookup[field.FieldName.String()]
		if !lookupExists {
			return false
		}
		if lookupValue != field.FieldValueType {
			return false
		}
	}
	return true
}
