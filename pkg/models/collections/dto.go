package modelcollections

import (
	"time"

	apppb "github.com/footprintai/restcol/api/pb/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func NewPbCollectionMetadata(mc *ModelCollection) *apppb.CollectionMetadata {
	metadata := &apppb.CollectionMetadata{
		Pid:        mc.ModelProjectID.String(),
		Cid:        mc.ID.String(),
		XCreatedAt: timestamppb.New(mc.CreatedAt),
	}
	if deletedAt, _ := mc.DeletedAt.Value(); deletedAt != nil {
		metadata.XDeletedAt = timestamppb.New(deletedAt.(time.Time))
	}
	return metadata
}

func NewModelSchema(reqSchemas []*apppb.SchemaField) (ModelSchema, error) {
	var fields []ModelFieldSchema

	for _, reqSchema := range reqSchemas {
		valType := NewSwaggerValueType(reqSchema.Datatype)
		valueValue, err := NewSwagValue(reqSchema.Example)
		if err != nil {
			return ModelSchema{}, err
		}
		fields = append(fields, ModelFieldSchema{
			FieldName:      reqSchema.Name,
			FieldValueType: valType,
			FieldExample:   valueValue,
		})
	}

	return ModelSchema{
		Fields: fields,
	}, nil
}

func NewPbSchemaFields(m ModelSchema) []*apppb.SchemaField {
	var pbfields []*apppb.SchemaField

	for _, field := range m.Fields {
		pbfields = append(pbfields, &apppb.SchemaField{
			Name:     field.FieldName,
			Datatype: field.FieldValueType.Proto(),
			Example:  field.FieldExample.Proto(),
		})
	}

	return pbfields
}
