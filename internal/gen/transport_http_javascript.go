package gen

import (
	"fmt"
	"go/types"
	"io"

	"emperror.dev/errors"
	"github.com/iancoleman/strcase"
	"github.com/nnnewb/jk/internal/domain"
)

type Writer interface {
	io.Writer
	io.StringWriter
}

func generateProps(wr Writer) error {
	_, err := wr.WriteString(`base_url: "",`)
	if err != nil {
		return err
	}
	return nil
}

func generateAjax(wr Writer) error {
	_, err := wr.WriteString(`
    /**
     * send request
     * @param method
     * @param url
     * @param config {{headers?: HeadersInit, params?: URLSearchParams, data?: any}} request configuration
     * @returns { Promise<Response> }
     */
    ajax: function (method, url, config) {
        method = method.toLowerCase();
        let u = null;
        if (this.base_url && this.base_url.length !== 0) {
            u = new URL(url, this.base_url)
        } else {
            u = new URL(url);
        }

		let params = undefined;
		if (config.params) {
			params = config.params;
		}

        u.searchParams = params;

		let body = undefined;
		if (config.data) {
			body = JSON.stringify(config.data);
		}

		let headers = undefined;
		if (config.headers) {
			headers = config.headers;
		}

        return fetch(u, {method, headers, body});
    },
`)
	if err != nil {
		return err
	}
	return nil
}

func generateAPIPath(wr Writer, method *domain.Method) error {
	_, err := wr.WriteString("/**\n * \n")
	if err != nil {
		return err
	}

	err = generateJSDoc(wr, method)
	if err != nil {
		return err
	}

	_, err = wr.WriteString(" */\n")
	if err != nil {
		return err
	}

	_, err = wr.WriteString(strcase.ToSnake(method.Func.Name()))
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(wr, ": function(payload, config) {\nreturn this.ajax(\n'%s', \n'%s', \nconfig || {},\n);\n},\n", method.Annotations.HTTPMethod, method.Annotations.HTTPPath)
	if err != nil {
		return err
	}
	return nil
}

func generateJSDoc(wr Writer, method *domain.Method) error {
	signature := method.Func.Type().(*types.Signature)
	reqType := signature.Params().At(1)

	_, err := wr.WriteString(" * @param payload {")
	if err != nil {
		return err
	}

	err = generateJSDocSchema(wr, reqType.Type())
	if err != nil {
		return err
	}

	_, err = wr.WriteString("} \n")
	if err != nil {
		return err
	}

	_, err = wr.WriteString(" * @param config {{headers?: HeadersInit}} request configuration\n")
	if err != nil {
		return err
	}

	return nil
}

func generateJSDocSchema(wr Writer, typ types.Type) error {
	switch t := typ.(type) {
	case *types.Named:
		return generateJSDocSchema(wr, t.Underlying())
	case *types.Pointer:
		return generateJSDocSchema(wr, t.Elem())
	case *types.Array:
		_, err := wr.WriteString("Array<")
		if err != nil {
			return err
		}
		err = generateJSDocSchema(wr, t.Elem())
		if err != nil {
			return err
		}
		_, err = wr.WriteString(">")
		if err != nil {
			return err
		}
	case *types.Slice:
		_, err := wr.WriteString("Array<")
		if err != nil {
			return err
		}

		err = generateJSDocSchema(wr, t.Elem())
		if err != nil {
			return err
		}

		_, err = wr.WriteString(">")
		if err != nil {
			return err
		}
	case *types.Map:
		_, err := wr.WriteString("Map<")
		if err != nil {
			return err
		}

		err = generateJSDocSchema(wr, t.Key())
		if err != nil {
			return err
		}

		_, err = wr.WriteString(",")
		if err != nil {
			return err
		}

		err = generateJSDocSchema(wr, t.Elem())
		if err != nil {
			return err
		}

		_, err = wr.WriteString(">")
		if err != nil {
			return err
		}
	case *types.Basic:
		switch t.Kind() {
		case types.Bool:
			_, err := wr.WriteString("boolean")
			if err != nil {
				return err
			}
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64:
			_, err := wr.WriteString("number")
			if err != nil {
				return err
			}
		case types.String:
			_, err := wr.WriteString("string")
			if err != nil {
				return err
			}
		default:
			panic(errors.Errorf("unserializable basic type %v", t.Kind()))
		}
	case *types.Struct:
		_, err := wr.WriteString("{")
		if err != nil {
			return err
		}

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
				_, err := wr.WriteString(jsonName)
				if err != nil {
					return err
				}

				_, err = wr.WriteString(":")
				if err != nil {
					return err
				}

				err = generateJSDocSchema(wr, field.Type())
				if err != nil {
					return err
				}

				_, err = wr.WriteString(",")
				if err != nil {
					return err
				}
			}
		}

		_, err = wr.WriteString("}")
		if err != nil {
			return err
		}
	default:
		panic(errors.Errorf("unserializable type %v", typ))
	}
	return nil
}

func GenerateJavascriptClient(wr Writer, service *domain.Service) error {
	httpPopulateDefaultAnnotations(service)
	_, err := wr.WriteString("export default {\n")
	if err != nil {
		return err
	}

	err = generateProps(wr)
	if err != nil {
		return err
	}

	err = generateAjax(wr)
	if err != nil {
		return err
	}

	for _, method := range service.Methods {
		err := generateAPIPath(wr, method)
		if err != nil {
			return err
		}
	}

	_, err = wr.WriteString("};")
	if err != nil {
		return err
	}
	return nil
}
