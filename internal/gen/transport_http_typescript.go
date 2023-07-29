package gen

import (
	"fmt"
	"go/types"
	"io"
	"net/http"
	"strings"

	"emperror.dev/errors"
	"github.com/iancoleman/strcase"
	"github.com/nnnewb/jk/internal/domain"
)

func generateNamedInterfaceDeclaration(wr io.Writer, memo map[string]bool, named *types.Named) error {
	if !named.Obj().Exported() {
		return nil
	}

	if memo[named.Obj().Name()] {
		return nil
	}

	switch t := named.Underlying().(type) {
	case *types.Pointer:
		if underlying, ok := t.Elem().(*types.Named); ok {
			return generateNamedInterfaceDeclaration(wr, memo, underlying)
		}
	case *types.Slice:
		if underlying, ok := t.Elem().(*types.Named); ok {
			return generateNamedInterfaceDeclaration(wr, memo, underlying)
		}
	case *types.Map:
		if underlying, ok := t.Elem().(*types.Named); ok {
			return generateNamedInterfaceDeclaration(wr, memo, underlying)
		}
	case *types.Named:
		return generateNamedInterfaceDeclaration(wr, memo, t)
	case *types.Struct:
		for i := 0; i < t.NumFields(); i++ {
			f := t.Field(i)
			if !f.Exported() {
				continue
			}

			switch ft := f.Type().(type) {
			case *types.Pointer:
				if underlying, ok := ft.Elem().(*types.Named); ok {
					err := generateNamedInterfaceDeclaration(wr, memo, underlying)
					if err != nil {
						return err
					}
				}
			case *types.Slice:
				if underlying, ok := ft.Elem().(*types.Named); ok {
					err := generateNamedInterfaceDeclaration(wr, memo, underlying)
					if err != nil {
						return err
					}
				}
			case *types.Map:
				if underlying, ok := ft.Elem().(*types.Named); ok {
					err := generateNamedInterfaceDeclaration(wr, memo, underlying)
					if err != nil {
						return err
					}
				}
			case *types.Named:
				err := generateNamedInterfaceDeclaration(wr, memo, ft)
				if err != nil {
					return err
				}
			}
		}
	default:
		return nil
	}

	memo[named.Obj().Name()] = true
	return generateTypescriptSchema(wr, named, 0)
}

func generateInterfaceDeclaration(wr io.Writer, service *domain.Service) error {
	allNamedType := make(map[string]bool)
	for _, method := range service.Methods {
		if allNamedType[method.RequestTypeName()] {
			continue
		}
		err := generateNamedInterfaceDeclaration(wr, allNamedType, method.RequestType().(*types.Pointer).Elem().(*types.Named))
		if err != nil {
			return err
		}
		if allNamedType[method.ResponseTypeName()] {
			continue
		}
		err = generateNamedInterfaceDeclaration(wr, allNamedType, method.ResponseType().(*types.Pointer).Elem().(*types.Named))
		if err != nil {
			return err
		}
	}

	return nil
}

func generateAjaxTypescript(wr io.Writer) error {
	_, err := fmt.Fprint(wr, `
`)
	if err != nil {
		return err
	}
	return nil
}

func generateAPIPathTypescript(wr io.Writer, method *domain.Method) error {
	var initPayload string
	switch method.Annotations.HTTPMethod {
	case http.MethodGet, http.MethodDelete:
		initPayload = `Object.getOwnPropertyNames(payload).map(prop => u.searchParams.append(prop, payload[prop]));`
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		initPayload = `init.body = JSON.stringify(payload);`
	}
	_, err := fmt.Fprintf(wr, `
	%s: async function(payload: %s, init?: RequestInit): Promise<%s> {
		const u = new URL("%s", this.baseURL);
		%s
		init.method = "%s";
		const req = new Request(u, init);
		const resp = await fetch(req, init);
		return await resp.json();
	},`,
		strcase.ToSnake(method.Func.Name()),
		method.RequestTypeName(),
		method.ResponseTypeName(),
		method.Annotations.HTTPPath,
		initPayload,
		method.Annotations.HTTPMethod,
	)
	if err != nil {
		return err
	}
	return nil
}

func generateTypescriptSchema(wr io.Writer, typ types.Type, depth int) error {
	switch t := typ.(type) {
	case *types.Named:
		if depth == 0 {
			_, err := fmt.Fprintf(wr, "interface %s {\n", t.Obj().Name())
			if err != nil {
				return err
			}

			err = generateTypescriptSchema(wr, t.Underlying(), depth+1)
			if err != nil {
				return err
			}

			_, err = io.WriteString(wr, "}\n")
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprintf(wr, "%s", t.Obj().Name())
			if err != nil {
				return err
			}
		}
	case *types.Pointer:
		return generateTypescriptSchema(wr, t.Elem(), depth+1)
	case *types.Array:
		_, err := io.WriteString(wr, "Array<")
		if err != nil {
			return err
		}
		err = generateTypescriptSchema(wr, t.Elem(), depth+1)
		if err != nil {
			return err
		}
		_, err = io.WriteString(wr, ">")
		if err != nil {
			return err
		}
	case *types.Slice:
		_, err := io.WriteString(wr, "Array<")
		if err != nil {
			return err
		}

		err = generateTypescriptSchema(wr, t.Elem(), depth+1)
		if err != nil {
			return err
		}

		_, err = io.WriteString(wr, ">")
		if err != nil {
			return err
		}
	case *types.Map:
		_, err := io.WriteString(wr, "Map<")
		if err != nil {
			return err
		}

		err = generateTypescriptSchema(wr, t.Key(), depth+1)
		if err != nil {
			return err
		}

		_, err = io.WriteString(wr, ",")
		if err != nil {
			return err
		}

		err = generateTypescriptSchema(wr, t.Elem(), depth+1)
		if err != nil {
			return err
		}

		_, err = io.WriteString(wr, ">")
		if err != nil {
			return err
		}
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			_, err := io.WriteString(wr, "boolean")
			if err != nil {
				return err
			}
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64:
			_, err := io.WriteString(wr, "number")
			if err != nil {
				return err
			}
		case types.String:
			_, err := io.WriteString(wr, "string")
			if err != nil {
				return err
			}
		default:
			panic(errors.Errorf("unserializable basic type %v", t.Kind()))
		}
	case *types.Struct:
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
			if jsonName, ok = getJsonName(t.Tag(i)); !ok {
				jsonName = field.Name()
			} else if jsonName == "-" {
				continue
			}

			// Check if the field is exported
			if field.Exported() {
				_, err := fmt.Fprintf(wr, "%s%s: ", strings.Repeat(" ", 2), jsonName)
				if err != nil {
					return err
				}

				err = generateTypescriptSchema(wr, field.Type(), depth+1)
				if err != nil {
					return err
				}

				_, err = io.WriteString(wr, ";\n")
				if err != nil {
					return err
				}
			}
		}
	default:
		panic(errors.Errorf("unserializable type %v", typ))
	}
	return nil
}

func GenerateTypeScriptClient(wr io.Writer, service *domain.Service) error {
	httpPopulateDefaultAnnotations(service)
	err := generateInterfaceDeclaration(wr, service)
	if err != nil {
		return err
	}

	_, err = io.WriteString(wr, `
export default {
	baseURL: "",
`)

	if err != nil {
		return err
	}

	err = generateAjaxTypescript(wr)
	if err != nil {
		return err
	}

	for _, method := range service.Methods {
		err := generateAPIPathTypescript(wr, method)
		if err != nil {
			return err
		}
	}

	_, err = io.WriteString(wr, "};")
	if err != nil {
		return err
	}
	return nil
}
