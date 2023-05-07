package gen

import (
	"fmt"
	"github.com/go-openapi/spec"
	"github.com/iancoleman/strcase"
	"github.com/juju/errors"
	"go/types"
	"io"
	"reflect"
	"strings"
)

func GenerateSwagger(wr io.Writer, svc *types.Named, apiVer, Ver string) error {
	root := spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Swagger:  "2.0",
			Consumes: []string{"application/json"},
			Produces: []string{"application/json"},
			Schemes:  []string{"http", "https"},
			Host:     "localhost",
			BasePath: fmt.Sprintf("/api/%s/%s", apiVer, svc.Obj().Name()),
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Title:   svc.Obj().Name(), // TODO
					Version: Ver,              // TODO
				},
			},
			Paths: generatePaths(svc),
			SecurityDefinitions: spec.SecurityDefinitions{
				"api_key": {
					SecuritySchemeProps: spec.SecuritySchemeProps{
						Description: "simple api key",
						Type:        "apiKey",
						Name:        "X-Authentication",
						In:          "header",
					},
				},
			},
			Tags: []spec.Tag{
				spec.NewTag(svc.Obj().Name(), "", nil),
			},
		},
	}

	result, err := root.MarshalJSON()
	if err != nil {
		return errors.Trace(err)
	}

	_, err = wr.Write(result)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}

func generatePaths(svc *types.Named) *spec.Paths {
	ret := &spec.Paths{}
	ret.Paths = map[string]spec.PathItem{}
	iface, ok := svc.Underlying().(*types.Interface)
	if !ok {
		panic(errors.Errorf("svc %s underlying type is not interface", svc.Obj().Name()))
	}

	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)
		if !method.Exported() {
			continue
		}

		item := spec.PathItem{}
		item.Post = spec.
			NewOperation(strcase.ToKebab(method.Name())).
			WithConsumes("application/json").
			WithProduces("application/json").
			WithDefaultResponse(generateResponse(method)).
			WithTags(svc.Obj().Name())
		item.Post.Parameters = append(item.Post.Parameters, generateParameters(method)...)
		ret.Paths["/"+strcase.ToKebab(method.Name())] = item
	}
	return ret
}

func generateParameters(fun *types.Func) []spec.Parameter {
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

			jsonTag := reflect.StructTag(t.Tag(i)).Get("json")
			if jsonTag == "-" {
				continue
			}

			var jsonName string
			for _, v := range strings.Split(jsonTag, ",") {
				if v != "omitempty" && strings.TrimSpace(v) != "" {
					jsonName = v
					break
				}
			}

			if jsonName == "" {
				jsonName = field.Name()
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
