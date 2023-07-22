package utils

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"

	"emperror.dev/errors"
	"github.com/dave/jennifer/jen"
	"github.com/nnnewb/battery/slices"
)

func InitializeFileCommon(f *jen.File) {
	f.ImportAlias("github.com/swaggo/http-swagger/v2", "httpSwagger")
	f.ImportName("emperror.dev/errors", "errors")
	f.ImportName("github.com/go-kit/kit/endpoint", "endpoint")
	f.ImportAlias("github.com/go-kit/kit/transport/http", "khttp")
	f.ImportAlias("github.com/gorilla/schema", "schema")
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
			jsonTags := strings.Split(jsonTag, ",")
			if len(jsonTags) == 0 {
				return errors.Errorf("the Code field has no \"json\" tag, please add the `json:\"code\"` tag to the Code field")
			}
			// The field name must be "code" before the first comma
			if jsonTags[0] != "code" {
				return errors.Errorf("the first word before the comma in the \"json\" tag of the Code field must be \"code\", to ensure consistency of the response structure")
			}
			// The code field cannot have omitempty and string tags, to ensure that the serialized result must exist and be of the JSON Number type
			if slices.Any(jsonTags, func(v string) bool { return v == "omitempty" || v == "string" }) {
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
			jsonTags := strings.Split(jsonTag, ",")
			if len(jsonTags) == 0 {
				return errors.Errorf("the Message field has no \"json\" tag, please add the `json:\"message\"` tag to the Message field")
			}
			// The field name must be "message" before the first comma
			if jsonTags[0] != "message" {
				return errors.Errorf("the first word before the comma in the \"json\" tag of the Message field must be \"message\", to ensure consistency of the response structure")
			}
			// The message field cannot have omitempty tag, to ensure that this field always appears
			if slices.Any(jsonTags, func(v string) bool { return v == "omitempty" }) {
				return errors.Errorf("the \"json\" tag of the Message field cannot contain \"omitempty\", to ensure that this field always appears")
			}
			return nil
		}
	}
	return fmt.Errorf("the Message field was not found in the %s.%s struct", named.Obj().Pkg().Path(), named.Obj().Name())
}
