package modelcollections

import (
	"time"

	apppb "github.com/footprintai/restcol/api/pb"
	dotnotation "github.com/footprintai/restcol/pkg/notation/dot"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func NewPbCollectionMetadata(mc *ModelCollection) *apppb.CollectionMetadata {
	metadata := &apppb.CollectionMetadata{
		ProjectId:    mc.ModelProjectID.String(),
		CollectionId: mc.ID.String(),
		XCreatedAt:   timestamppb.New(mc.CreatedAt),
	}
	if deletedAt, _ := mc.DeletedAt.Value(); deletedAt != nil {
		metadata.XDeletedAt = timestamppb.New(deletedAt.(time.Time))
	}
	return metadata
}

func NewModelSchema(reqSchemas []*apppb.SchemaField) (*ModelSchema, error) {
	var fields []*ModelFieldSchema

	for _, reqSchema := range reqSchemas {
		valType := NewSwaggerValueType(reqSchema.Datatype)
		valueValue, err := NewSwagValue(reqSchema.Example)
		if err != nil {
			return nil, err
		}
		fields = append(fields, &ModelFieldSchema{
			FieldName:      dotnotation.Parse(reqSchema.Name),
			FieldValueType: valType,
			FieldExample:   valueValue,
		})
	}

	return &ModelSchema{
		Fields: fields,
	}, nil
}

func NewPbSchemaFields(m *ModelSchema) []*apppb.SchemaField {
	if m == nil {
		return nil
	}

	var pbfields []*apppb.SchemaField
	for _, field := range m.Fields {
		pbfields = append(pbfields, &apppb.SchemaField{
			Name:     field.FieldName.String(),
			Datatype: field.FieldValueType.Proto(),
			Example:  field.FieldExample.Proto(),
		})
	}

	return pbfields
}
