package gen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"go/types"
)

func InitializeFileCommon(f *jen.File) {
	f.ImportAlias("github.com/swaggo/http-swagger/v2", "httpSwagger")
	f.ImportName("github.com/juju/errors", "errors")
	f.ImportName("github.com/go-kit/kit/endpoint", "endpoint")
}

// checkParams checks if the function signature meets the following requirements:
//   - The first parameter must be context.Context.
//   - The second parameter must be an exported struct.
func checkParams(params *types.Tuple) error {
	// Check first parameter.
	if named, ok := params.At(0).Type().(*types.Named); ok {
		if named.Obj().Pkg().Path() != "context" || named.Obj().Name() != "Context" {
			return fmt.Errorf("first parameter must be context.Context")
		}
	} else {
		return fmt.Errorf("first parameter must be context.Context")
	}

	// Check second parameter.
	if named, ok := params.At(1).Type().(*types.Named); ok {
		if named.Obj().Exported() == false {
			return fmt.Errorf("second parameter must be an exported struct")
		}

		if !IsSerializable(named.Underlying()) {
			return fmt.Errorf("second parameter must be serializable")
		}
	} else {
		return fmt.Errorf("second parameter must be an exported struct")
	}

	return nil
}

// checkResults checks if the function signature meets the following requirements:
//   - The first return value must be an exported struct or slice of structs.
//   - The first return value must be serializable.
//   - The second return value must be of type error.
func checkResults(results *types.Tuple) error {
	// Check first return value.
	if results.Len() < 1 {
		return fmt.Errorf("function must have at least one return value")
	}

	// Check if the first return value is an exported struct or slice of structs
	if named, ok := results.At(0).Type().(*types.Named); ok {
		if named.Obj().Exported() == false {
			return fmt.Errorf("first return value must be an exported struct or slice of structs")
		}

		// Check if the first return value is serializable
		if !IsSerializable(named.Underlying()) {
			return fmt.Errorf("first return value must be serializable")
		}
	} else if slice, ok := results.At(0).Type().(*types.Slice); ok {
		if named, ok := slice.Elem().(*types.Named); ok {
			if named.Obj().Exported() == false {
				return fmt.Errorf("first return value must be an exported struct or slice of structs")
			}

			// Check if the first return value is serializable
			if !IsSerializable(named.Underlying()) {
				return fmt.Errorf("first return value must be serializable")
			}
		} else {
			return fmt.Errorf("first return value must be an exported struct or slice of structs")
		}
	} else {
		return fmt.Errorf("first return value must be an exported struct or slice of structs")
	}

	// Check second return value.
	if results.Len() < 2 {
		return fmt.Errorf("function must have at least two return values")
	}

	// Check if the second return value is of type error
	if named, ok := results.At(1).Type().(*types.Named); ok {
		if !types.Identical(types.Universe.Lookup("error").Type(), named) {
			return fmt.Errorf("second return value must be of type error")
		}
	} else {
		return fmt.Errorf("second return value must be of type error")
	}

	return nil
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
