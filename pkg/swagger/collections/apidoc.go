package collectiondoc

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	openapiloads "github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	api "github.com/footprintai/restcol/api"
	modelcollections "github.com/footprintai/restcol/pkg/models/collections"
	swagdef "github.com/footprintai/restcol/pkg/runtime/swagdef"
)

type CollectionSwaggerDoc struct {
	Collections []*modelcollections.ModelCollection
}

func NewCollectionSwaggerDoc(collections ...*modelcollections.ModelCollection) *CollectionSwaggerDoc {
	return &CollectionSwaggerDoc{
		Collections: collections,
	}
}

func (c *CollectionSwaggerDoc) RenderDoc() (string, error) {

	embedFsLoader := func(path string) (json.RawMessage, error) {
		rawBytes, err := api.OpenApiV2Fs.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return json.RawMessage(rawBytes), nil
	}

	newSwagDoc := func() (*spec.Swagger, error) {
		swagDoc, err := openapiloads.Spec(
			"openapiv2/proto/restcol.swagger.json",
			openapiloads.WithDocLoader(
				embedFsLoader,
			),
		)
		if err != nil {
			return nil, err
		}
		return swagDoc.Spec(), nil
	}

	var pathSpec []*spec.Swagger
	for _, col := range c.Collections {
		swagSpec, _ := newSwagDoc()
		specClone, _ := copyPathsWithFilter(
			swagSpec,
			cidPathFilter,
			documentTagFilter,
		)
		if err := replacePathsWithCollection(col, specClone); err != nil {
			return "", err
		}
		pathSpec = append(pathSpec, specClone)
	}

	// merge multiple specs
	newSwagSpec, _ := newSwagDoc()
	if err := mergeSwagPaths(newSwagSpec, pathSpec...); err != nil {
		return "", err
	}
	swagJsonBlob, err := json.Marshal(newSwagSpec)
	if err != nil {
		return "", err
	}
	return string(swagJsonBlob), nil
}

type PathFilterFunc func(path string) bool

// cidPathFilter filters a path with {cid}
func cidPathFilter(path string) bool {
	if strings.Contains(path, "{cid}") {
		return true
	}
	return false
}

type TagFilterFunc func(tags []string) bool

func documentTagFilter(tags []string) bool {
	for _, tag := range tags {
		if tag == "document" {
			return true
		}
	}
	return false
}

func copyPathsWithFilter(origSpec *spec.Swagger, pathFilter PathFilterFunc, tagFilter TagFilterFunc) (*spec.Swagger, error) {

	// clone spec.Swagger
	retSpec := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Paths: &spec.Paths{
				Paths: make(map[string]spec.PathItem),
			},
		},
	}

	for path, pathItem := range origSpec.SwaggerProps.Paths.Paths {
		if !pathFilter(path) {
			continue
		}
		if pathItem.PathItemProps.Get != nil {
			if tagFilter(pathItem.PathItemProps.Get.OperationProps.Tags) {
				retSpec.SwaggerProps.Paths.Paths[path] = pathItem
			}
		}
		if pathItem.PathItemProps.Put != nil {
			if tagFilter(pathItem.PathItemProps.Put.OperationProps.Tags) {
				retSpec.SwaggerProps.Paths.Paths[path] = pathItem
			}
		}
		if pathItem.PathItemProps.Post != nil {
			if tagFilter(pathItem.PathItemProps.Post.OperationProps.Tags) {
				retSpec.SwaggerProps.Paths.Paths[path] = pathItem
			}
		}
		if pathItem.PathItemProps.Delete != nil {
			if tagFilter(pathItem.PathItemProps.Delete.OperationProps.Tags) {
				retSpec.SwaggerProps.Paths.Paths[path] = pathItem
			}
		}

	}

	return retSpec, nil
}

// replacePathsWithCollection expands $cid, $pid with values defined in col
func replacePathsWithCollection(col *modelcollections.ModelCollection, specClone *spec.Swagger) error {
	// all response of a single collection now is under `apiRequestResponse`
	responseDefs, err := swagdef.ModelFieldsSchemaToSwagDef(col.Schemas[0].Fields, "apiRequestResponse")
	if err != nil {
		return err
	}
	apiRequestBodySchema := (*responseDefs)["apiRequestResponse"]
	apiResponseSchema := (*responseDefs)["apiRequestResponse"]

	var cidReplacer = regexp.MustCompile(`\{cid\}`)
	var pidReplacer = regexp.MustCompile(`\{pid\}`)
	pathClone := &spec.Paths{
		Paths: map[string]spec.PathItem{},
	}
	// things to do:
	// 1. replace {cid} and {pid} defined path param with collection values
	// 2. add example value on params with cid as example
	for path, pathItem := range specClone.SwaggerProps.Paths.Paths {
		replacedPath :=
			pidReplacer.ReplaceAllString(
				cidReplacer.ReplaceAllString(
					path, col.ID.String(),
				), col.ModelProjectID.String(),
			)
		cidParam := spec.Parameter{
			ParamProps: spec.ParamProps{
				Name: "cid",
			},
		}
		pidParam := spec.Parameter{
			ParamProps: spec.ParamProps{
				Name: "pid",
			},
		}
		delParams := []spec.Parameter{cidParam, pidParam}
		// replace pathItem properties
		if pathItem.PathItemProps.Get != nil {
			updateToSwagOperation(pathItem.PathItemProps.Get, col.Summary, delParams, nil, &apiResponseSchema)
		}
		if pathItem.PathItemProps.Put != nil {
			updateToSwagOperation(pathItem.PathItemProps.Put, col.Summary, delParams, &apiRequestBodySchema, &apiResponseSchema)
		}
		if pathItem.PathItemProps.Post != nil {
			updateToSwagOperation(pathItem.PathItemProps.Post, col.Summary, delParams, &apiRequestBodySchema, &apiResponseSchema)
		}
		if pathItem.PathItemProps.Delete != nil {
			updateToSwagOperation(pathItem.PathItemProps.Delete, col.Summary, delParams, nil, nil)
		}
		pathClone.Paths[replacedPath] = pathItem
	}
	specClone.SwaggerProps.Paths = pathClone
	return nil
}

func updateToSwagOperation(op *spec.Operation, newSummary string, delParams []spec.Parameter, apiRequestBodySchema *spec.Schema, apiResponseSchema *spec.Schema) error {
	op.OperationProps.Summary = newSummary

	curParamMap := make(map[string]spec.Parameter) // keyed by name
	for _, curParam := range op.OperationProps.Parameters {
		curParamMap[curParam.Name] = curParam
	}

	// delete ${cid} params as it is not required now
	for _, delParam := range delParams {
		if _, found := curParamMap[delParam.Name]; found {
			delete(curParamMap, delParam.Name)
		}
	}

	// add body params regarding to the data foramat
	if apiRequestBodySchema != nil {
		bodyParam, bodyExists := curParamMap["body"]
		if bodyExists && bodyParam.In == "body" {
			bodyParam.ParamProps.Schema = apiRequestBodySchema
			curParamMap["body"] = bodyParam
		}
	}
	var paramSlice []spec.Parameter
	for _, p := range curParamMap {
		paramSlice = append(paramSlice, p)
	}
	op.OperationProps.Parameters = paramSlice

	// update the response format
	if apiResponseSchema != nil {
		normalResponse := op.OperationProps.Responses.ResponsesProps.StatusCodeResponses[200]
		normalResponse.ResponseProps.Schema = apiResponseSchema
		op.OperationProps.Responses.ResponsesProps.StatusCodeResponses[200] = normalResponse
	}

	return nil
}

func mergeSwagPaths(dst *spec.Swagger, froms ...*spec.Swagger) error {
	for _, from := range froms {
		for fromPath, fromPathItem := range from.SwaggerProps.Paths.Paths {
			if _, exist := dst.SwaggerProps.Paths.Paths[fromPath]; exist {
				dstItem := dst.SwaggerProps.Paths.Paths[fromPath]
				if err := copyPathItemProps(&dstItem, fromPathItem); err != nil {
					return err
				}
				dst.SwaggerProps.Paths.Paths[fromPath] = dstItem
			} else {
				dst.SwaggerProps.Paths.Paths[fromPath] = fromPathItem
			}
		}
	}
	return nil
}

func copyPathItemProps(to *spec.PathItem, from spec.PathItem) error {
	if from.PathItemProps.Get != nil && to.PathItemProps.Get == nil {
		to.PathItemProps.Get = mergeOperationParameters(to.PathItemProps.Get, from.PathItemProps.Get)
	} else if from.PathItemProps.Post != nil && to.PathItemProps.Post == nil {
		to.PathItemProps.Post = mergeOperationParameters(to.PathItemProps.Post, from.PathItemProps.Post)
	} else if from.PathItemProps.Put != nil && to.PathItemProps.Put == nil {
		to.PathItemProps.Put = mergeOperationParameters(to.PathItemProps.Put, from.PathItemProps.Put)
	} else if from.PathItemProps.Delete != nil && to.PathItemProps.Delete == nil {
		to.PathItemProps.Delete = mergeOperationParameters(to.PathItemProps.Delete, from.PathItemProps.Delete)
	} else {
		return errors.New("swag: unable to merge pathitemprops.")
	}
	return nil

}

// mergeOperationParameters use $from operation as reference
// more like `cp <from> <to>` in shell
func mergeOperationParameters(to *spec.Operation, from *spec.Operation) *spec.Operation {
	return from
}
