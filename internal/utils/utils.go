package utils

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/juju/errors"
	"github.com/nnnewb/battery/slices"
)

func InitializeFileCommon(f *jen.File) {
	f.ImportAlias("github.com/swaggo/http-swagger/v2", "httpSwagger")
	f.ImportName("github.com/juju/errors", "errors")
	f.ImportName("github.com/go-kit/kit/endpoint", "endpoint")
	f.ImportAlias("github.com/go-kit/kit/transport/http", "khttp")
}

// CheckParams checks if the function signature meets the following requirements:
//   - The first parameter must be context.Context.
//   - The second parameter must be an exported struct.
func CheckParams(params *types.Tuple) error {
	if params.Len() != 2 {
		return errors.New("the function must have 2 parameter, and the second parameter must be a pointer to an exported and serializable struct type, and the first parameter must be an context.Context")
	}

	var (
		ptr     *types.Pointer
		named   *types.Named
		ok      bool
		ctxType = params.At(0).Type()
		reqType = params.At(1).Type()
	)

	// Check first parameter.
	if named, ok = ctxType.(*types.Named); !ok {
		return fmt.Errorf("first parameter must be context.Context")
	}

	if named.Obj().Pkg().Path() != "context" || named.Obj().Name() != "Context" {
		return fmt.Errorf("first parameter must be context.Context")
	}

	// Check second parameter.
	if ptr, ok = reqType.(*types.Pointer); !ok {
		return fmt.Errorf("second parameter must be pointer to exported and serializable struct, but got %s", reqType)
	}

	if named, ok = ptr.Elem().(*types.Named); !ok {
		return fmt.Errorf("second parameter must be pointer to exported and serializable struct, but got %s", reqType)
	}

	if named.Obj().Exported() == false {
		return fmt.Errorf("second parameter must be pointer to exported and serializable struct, but got %s", reqType)
	}

	if !IsSerializable(named.Underlying()) {
		return fmt.Errorf("second parameter must be pointer to exported and serializable struct, but got %s (unserializable)", reqType)
	}

	return nil
}

// CheckResults checks if the function signature meets the following requirements:
//   - The first return value must be exported serializable struct.
//     response struct must have Code (int) and Message (string) field.
//   - The second return value must be of type error.
func CheckResults(results *types.Tuple) error {
	// Check first return value.
	if results.Len() != 2 {
		return errors.New("the function must have 2 return values, and the first return value must be a pointer to an exported and serializable struct type, and the second return value must be an error")
	}

	var (
		ptr      *types.Pointer
		named    *types.Named
		ok       bool
		respType = results.At(0).Type()
		errType  = results.At(1).Type()
	)

	// Check if the first return value is an exported struct or slice of structs
	if ptr, ok = respType.(*types.Pointer); !ok {
		return fmt.Errorf("the first return value must be a pointer to an exported and serializable struct, but got %s", respType)
	}

	if named, ok = ptr.Elem().(*types.Named); !ok {
		return fmt.Errorf("the first return value must be a pointer to an exported and serializable struct, but got %s", respType)
	}

	if !named.Obj().Exported() {
		return fmt.Errorf("the first return value must be a pointer to an exported and serializable struct, but got %s", respType)
	}

	if !IsSerializable(named.Underlying()) {
		return fmt.Errorf("the first return value must be a pointer to an exported and serializable struct, but got %s", respType)
	}

	if err := checkCodeField(named); err != nil {
		return err
	}

	if err := checkMessageField(named); err != nil {
		return err
	}

	if named, ok = errType.(*types.Named); !ok {
		return fmt.Errorf("the type of the second return value must be error, but got %s", errType)
	}

	if !types.Identical(types.Universe.Lookup("error").Type(), named) {
		return fmt.Errorf("the type of the second return value must be error, but got %s", errType)
	}

	return nil
}

func checkCodeField(named *types.Named) error {
	var (
		p  *types.Struct
		ok bool
	)
	if p, ok = named.Underlying().(*types.Struct); !ok {
		return fmt.Errorf("%s.%s is not a struct type", named.Obj().Pkg().Path(), named.Obj().Name())
	}
	for i := 0; i < p.NumFields(); i++ {
		field := p.Field(i)
		if !field.Exported() {
			continue
		}
		if field.Name() == "Code" {
			b, ok := field.Type().(*types.Basic)
			if ok && b.Kind() != types.Int {
				return fmt.Errorf("the type of field %s.%s must be int", named.Obj().Name(), field.Name())
			}
			fieldTag := p.Tag(i)
			jsonTag := reflect.StructTag(fieldTag).Get("json")
			jsonTags := slices.Slice[string](strings.Split(jsonTag, ","))
			if len(jsonTags) == 0 {
				return errors.Errorf("the Code field has no \"json\" tag, please add the `json:\"code\"` tag to the Code field")
			}
			// The field name must be "code" before the first comma
			if jsonTags[0] != "code" {
				return errors.Errorf("the first word before the comma in the \"json\" tag of the Code field must be \"code\", to ensure consistency of the response structure")
			}
			// The code field cannot have omitempty and string tags, to ensure that the serialized result must exist and be of the JSON Number type
			if jsonTags.Any(func(v string) bool { return v == "omitempty" || v == "string" }) {
				return errors.Errorf("the \"json\" tag of the Code field cannot contain \"omitempty\" and \"string\", to ensure that this field always appears and has the correct type")
			}
			return nil
		}
	}
	return fmt.Errorf("the Code field was not found in the %s.%s struct", named.Obj().Pkg().Path(), named.Obj().Name())
}

func checkMessageField(named *types.Named) error {
	var (
		p  *types.Struct
		ok bool
	)
	if p, ok = named.Underlying().(*types.Struct); !ok {
		return fmt.Errorf("%s.%s is not a struct type", named.Obj().Pkg().Path(), named.Obj().Name())
	}
	for i := 0; i < p.NumFields(); i++ {
		field := p.Field(i)
		if !field.Exported() {
			continue
		}
		if field.Name() == "Message" {
			b, ok := field.Type().(*types.Basic)
			if ok && b.Kind() != types.String {
				return fmt.Errorf("the type of field %s.%s must be string", named.Obj().Name(), field.Name())
			}
			fieldTag := p.Tag(i)
			jsonTag := reflect.StructTag(fieldTag).Get("json")
			jsonTags := slices.Slice[string](strings.Split(jsonTag, ","))
			if len(jsonTags) == 0 {
				return errors.Errorf("the Message field has no \"json\" tag, please add the `json:\"message\"` tag to the Message field")
			}
			// The field name must be "message" before the first comma
			if jsonTags[0] != "message" {
				return errors.Errorf("the first word before the comma in the \"json\" tag of the Message field must be \"message\", to ensure consistency of the response structure")
			}
			// The message field cannot have omitempty tag, to ensure that this field always appears
			if jsonTags.Any(func(v string) bool { return v == "omitempty" }) {
				return errors.Errorf("the \"json\" tag of the Message field cannot contain \"omitempty\", to ensure that this field always appears")
			}
			return nil
		}
	}
	return fmt.Errorf("the Message field was not found in the %s.%s struct", named.Obj().Pkg().Path(), named.Obj().Name())
}

// Define a function named generateTypeCode that takes a types.Type as input and returns a *jen.Statement
func generateTypeCode(typ types.Type) *jen.Statement {
	// Check the underlying type of the input type
	switch t := typ.(type) {
	// If the type is a named type, return the name of the type
	case *types.Named:
		if t.Obj().Pkg() == nil {
			return jen.Id(t.Obj().Name())
		} else {
			return jen.Qual(t.Obj().Pkg().Path(), t.Obj().Name())
		}
	// If the type is a pointer type, return a pointer to the generated code for the element type
	case *types.Pointer:
		return jen.Op("*").Add(generateTypeCode(t.Elem()))
	// If the type is an array type, return an array of the generated code for the element type
	case *types.Array:
		return jen.Index(jen.Lit(int(t.Len()))).Add(generateTypeCode(t.Elem()))
	// If the type is a slice type, return a slice of the generated code for the element type
	case *types.Slice:
		return jen.Index().Add(generateTypeCode(t.Elem()))
	// If the type is a map type, return a map with the generated code for the key and value types
	case *types.Map:
		return jen.Map(generateTypeCode(t.Key())).Add(generateTypeCode(t.Elem()))
	// If the type is a struct type, generate code for the struct fields
	case *types.Struct:
		// Define a new *jen.Statement for the struct
		structCode := jen.StructFunc(func(g *jen.Group) {
			// Iterate over all the fields of the struct
			for i := 0; i < t.NumFields(); i++ {
				field := t.Field(i)
				// Check if the field is exported
				if field.Exported() {
					// Add the generated code for the field to the struct
					g.Id(field.Name()).Add(generateTypeCode(field.Type()))
				}
			}
		})
		// Return the generated code for the struct
		return structCode
	// If the type is a function type, generate code for the function signature
	case *types.Signature:
		// Define a new *jen.Statement for the function signature
		sigCode := jen.Func().Params()
		// Iterate over all the parameters of the function
		for i := 0; i < t.Params().Len(); i++ {
			param := t.Params().At(i)
			paramName := param.Name()
			if paramName == "" {
				paramName = fmt.Sprintf("arg%d", i+1)
			}
			// Add the generated code for the parameter to the function signature
			sigCode.Add(jen.Id(paramName)).Add(generateTypeCode(param.Type()))
		}
		// Generate code for the return type of the function
		returnType := t.Results()
		if returnType.Len() == 1 {
			sigCode.Add(generateTypeCode(returnType.At(0).Type()))
		} else if returnType.Len() > 1 {
			sigCode.Add(jen.Index().Add(generateTypeCode(returnType)))
		}
		// Return the generated code for the function signature
		return sigCode
	// If the type is a basic type, return the name of the type
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			return jen.Bool()
		case types.Int:
			return jen.Int()
		case types.Int8:
			return jen.Int8()
		case types.Int16:
			return jen.Int16()
		case types.Int32:
			return jen.Int32()
		case types.Int64:
			return jen.Int64()
		case types.Uint:
			return jen.Uint()
		case types.Uint8:
			return jen.Uint8()
		case types.Uint16:
			return jen.Uint16()
		case types.Uint32:
			return jen.Uint32()
		case types.Uint64:
			return jen.Uint64()
		case types.Uintptr:
			return jen.Uintptr()
		case types.Float32:
			return jen.Float32()
		case types.Float64:
			return jen.Float64()
		case types.Complex64:
			return jen.Complex64()
		case types.Complex128:
			return jen.Complex128()
		case types.String:
			return jen.String()
		default:
			return jen.Empty()
		}
	// If the type is not recognized, return an empty *jen.Statement
	default:
		return jen.Empty()
	}
}
