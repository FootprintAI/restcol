/* package swagdef provides methods to operate json data and swag definition */
package swagdef

import (
	"encoding/json"
	"fmt"

	"github.com/go-openapi/spec"

	modelcollections "github.com/footprintai/restcol/pkg/models/collections"
	jsruntime "github.com/footprintai/restcol/pkg/runtime/js"
)

// swagJsonMessageToSwagDef returns *spec.Definition by giving swag definition in json, ex:
//
//	{
//		"definitions": {
//		"foo": {}
//	}
func swagJsonMessageToSwagDef(j json.RawMessage) (*spec.Definitions, error) {
	swagSpec := &spec.Swagger{}

	if err := swagSpec.UnmarshalJSON([]byte(j)); err != nil {
		return nil, err
	}
	return &swagSpec.SwaggerProps.Definitions, nil
}

func RawJsonMessageToSwagDef(rawJson []byte) (*spec.Definitions, error) {
	rt := jsruntime.NewJSRuntime("swag")
	defer rt.Close()

	if err := rt.Load("swagdef.js"); err != nil {
		return nil, err
	}
	scriptInRaw := fmt.Sprintf(`doconvert('%s')`, string(rawJson))
	swagJson, err := rt.Run(scriptInRaw)
	if err != nil {
		return nil, err
	}
	return swagJsonMessageToSwagDef(json.RawMessage(fmt.Sprintf("{%s}", swagJson)))

}

// ModelFieldsSchemaToSwagDef converts from dot-notation json structure
// into swag.Definitions
func ModelFieldsSchemaToSwagDef(fields []modelcollections.ModelFieldSchema, withPrefixs ...string) (*spec.Definitions, error) {

	fieldsWithJson, err := modelcollections.ModelFieldsSchema(fields).ToJSON(withPrefixs...)
	if err != nil {
		return nil, err
	}
	return RawJsonMessageToSwagDef(fieldsWithJson)

}
