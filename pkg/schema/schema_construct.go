package schema

import (
	"errors"
	"fmt"

	appmodelcollections "github.com/footprintai/restcol/pkg/models/collections"
)

func Build(fields []*appmodelcollections.ModelFieldSchema) (map[string]interface{}, error) {
	constructedData := make(map[string]interface{})

	for _, field := range fields {
		parts := field.FieldName.Parts
		var branchParts []string
		var leavePart string
		if len(parts) == 0 {
			// something wrong
			return nil, errors.New("schema: field.parts is empty")
		} else if len(parts) == 1 {
			leavePart = parts[0]
		} else {
			branchParts = parts[0 : len(parts)-1]
			leavePart = parts[len(parts)-1]
		}
		constructedDataPtr := constructedData
		// trace branch
		for _, branchPart := range branchParts {
			_, exist := constructedDataPtr[branchPart]
			if !exist {
				constructedDataPtr[branchPart] = make(map[string]interface{})
			}
			constructedDataPtr = constructedDataPtr[branchPart].(map[string]interface{})
		}
		// handle leave
		if _, exist := constructedDataPtr[leavePart]; exist {
			return nil, fmt.Errorf("schema: paths(%+v) duplicated", parts)
			// something wrong
		}
		constructedDataPtr[leavePart] = field.FieldExample.Proto().AsInterface()
	}
	return constructedData, nil
}
