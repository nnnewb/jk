package doc

import (
	"fmt"
	"go/types"
	"io"
	"strings"

	"emperror.dev/errors"
	"github.com/go-openapi/spec"
	"github.com/iancoleman/strcase"
	"github.com/nnnewb/jk/internal/domain"
	"github.com/nnnewb/jk/internal/gen/http/common"
	"github.com/nnnewb/jk/internal/utils"
)

func GenerateSwagger(wr io.Writer, service *domain.Service) error {
	common.HTTPPopulateDefaultAnnotations(service)

	paths, err := generatePaths(service)
	if err != nil {
		return err
	}

	root := spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger:  "2.0",
			Consumes: []string{"application/json"},
			Produces: []string{"application/json"},
			Schemes:  []string{"http", "https"},
			Host:     "localhost",
			BasePath: service.Annotations.HTTPBasePath,
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Title:   service.Annotations.SwaggerInfoAPITitle,
					Version: service.Annotations.SwaggerInfoAPIVersion,
				},
			},
			Paths: paths,
		},
	}

	result, err := root.MarshalJSON()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = wr.Write(result)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func generatePaths(service *domain.Service) (*spec.Paths, error) {
	ret := &spec.Paths{}
	ret.Paths = map[string]spec.PathItem{}

	for _, method := range service.Methods {
		if !method.Func.Exported() {
			continue
		}

		item := spec.PathItem{}
		operation := spec.
			NewOperation(strcase.ToKebab(method.Func.Name())).
			WithProduces("application/json").
			WithDefaultResponse(generateResponse(method.Func)).
			WithTags(service.Interface.Obj().Name())

		switch strings.ToLower(method.Annotations.HTTPMethod) {
		case "get":
			parameters := generateQueryParameters(method.Func)
			item.Get = operation
			item.Get.Parameters = append(item.Get.Parameters, parameters...)
		case "delete":
			parameters := generateQueryParameters(method.Func)
			item.Delete = operation
			item.Delete.Parameters = append(item.Delete.Parameters, parameters...)
		case "put":
			parameters := generatePostParameters(method.Func)
			operation.WithConsumes("application/json")
			item.Put = operation
			item.Put.Parameters = append(item.Put.Parameters, parameters...)
		case "patch":
			parameters := generatePostParameters(method.Func)
			operation.WithConsumes("application/json")
			item.Patch = operation
			item.Patch.Parameters = append(item.Patch.Parameters, parameters...)
		case "post":
			fallthrough
		default:
			parameters := generatePostParameters(method.Func)
			operation.WithConsumes("application/json")
			item.Post = operation
			item.Post.Parameters = append(item.Post.Parameters, parameters...)
		}

		ret.Paths[method.Annotations.HTTPPath] = item
	}

	return ret, nil
}

func generateQueryParameters(fun *types.Func) []spec.Parameter {
	signature := fun.Type().(*types.Signature)
	paramType := signature.Params().At(1).Type()

	var structType *types.Struct

	ptr := paramType.(*types.Pointer)
	if named, ok := ptr.Elem().(*types.Named); ok {
		structType = named.Underlying().(*types.Struct)
	} else {
		structType = ptr.Elem().(*types.Struct)
	}

	params := make([]spec.Parameter, 0, structType.NumFields())
	for i := 0; i < structType.NumFields(); i++ {
		f := structType.Field(i)
		if !f.Exported() {
			continue
		}

		if !utils.IsQueryStringSerializable(f.Type()) {
			panic(fmt.Errorf("unserializable query string parameter type %s", f.Type()))
		}

		var (
			jsonName string
			ok       bool
		)
		if jsonName, ok = common.GetJsonName(structType.Tag(i)); !ok {
			jsonName = f.Name()
		} else if jsonName == "-" {
			continue
		}

		param := spec.QueryParam(jsonName)
		switch ft := f.Type().(type) {
		case *types.Basic:
			switch ft.Kind() {
			case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
				types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
				param.Type = "integer"
			case types.Float32, types.Float64:
				param.Type = "number"
			case types.String:
				param.Type = "string"
			case types.Bool:
				param.Type = "boolean"
			default:
				panic(errors.Errorf("unserializable basic type %v", ft.Kind()))
			}
		}
		params = append(params, *param)
	}

	return params
}

func generatePostParameters(fun *types.Func) []spec.Parameter {
	signature := fun.Type().(*types.Signature)
	reqType := signature.Params().At(1)
	param := spec.Parameter{}
	param.Name = "payload"
	param.Schema = generateSchemaFromType(reqType.Type())
	param.In = "body"
	return []spec.Parameter{param}
}

func generateResponse(fun *types.Func) *spec.Response {
	signature := fun.Type().(*types.Signature)
	respType := signature.Results().At(0)
	return spec.
		NewResponse().
		WithSchema(generateSchemaFromType(respType.Type()))
}

func generateSchemaFromType(typ types.Type) *spec.Schema {
	switch t := typ.(type) {
	case *types.Named:
		return generateSchemaFromType(t.Underlying())
	case *types.Pointer:
		return generateSchemaFromType(t.Elem())
	case *types.Array:
		return spec.ArrayProperty(generateSchemaFromType(t.Elem()))
	case *types.Slice:
		return spec.ArrayProperty(generateSchemaFromType(t.Elem()))
	case *types.Map:
		return spec.MapProperty(generateSchemaFromType(t.Elem()))
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return spec.BoolProperty()
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64:
			return spec.Int64Property()
		case types.Float32, types.Float64:
			return spec.Float64Property()
		case types.String:
			return spec.StringProperty()
		default:
			panic(errors.Errorf("unserializable basic type %v", t.Kind()))
		}
	case *types.Struct:
		ret := &spec.Schema{}
		ret.Properties = make(spec.SchemaProperties)
		// Iterate over all the fields of the struct
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			if !field.Exported() {
				continue
			}

			var (
				jsonName string
				ok       bool
			)
			if jsonName, ok = common.GetJsonName(t.Tag(i)); !ok {
				jsonName = field.Name()
			} else if jsonName == "-" {
				continue
			}

			// Check if the field is exported
			if field.Exported() {
				ret.Properties[jsonName] = *generateSchemaFromType(field.Type())
			}
		}
		return ret
	default:
		panic(errors.Errorf("unserializable type %v", typ))
	}
}
