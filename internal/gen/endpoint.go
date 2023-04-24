package gen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"go/types"
	"strings"
)

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

// GenerateEndpointFactory generates Endpoint factory function code from function signature.
// The function signature must meet the following requirements:
//   - The first parameter must be context.Context.
//   - The second parameter must be an exported struct that is serializable.
//   - The first return value must be an exported struct or slice of structs that is serializable.
//     When returning a slice of structs, the response content will be encapsulated in the Items field.
//     When returning a single struct, the response content will be encapsulated in the Data field.
//   - The second return value must be of type error. The result of the Error method of this interface
//     will be used as the response msg.
func GenerateEndpointFactory(file *jen.File, svc *types.Named, f *types.Func) error {
	signature := f.Type().(*types.Signature)
	params := signature.Params()
	if params.Len() != 2 {
		return fmt.Errorf("function must have exactly two parameters")
	}

	if err := checkParams(params); err != nil {
		return err
	}

	if err := checkResults(signature.Results()); err != nil {
		return err
	}

	reqParam := signature.Params().At(1).Type().(*types.Named)
	file.
		Commentf("// Make%sEndpoint create endpoint.Endpoint for function %s.%s\n", f.Name(), f.Pkg().Path(), f.Name()).
		Func().Id(fmt.Sprintf("Make%sEndpoint", f.Name())).
		Params(jen.Id("svc").Qual(svc.Obj().Pkg().Path(), svc.Obj().Name())).
		Params(jen.Qual("github.com/go-kit/kit/endpoint", "Endpoint")).
		BlockFunc(func(g *jen.Group) {
			g.ReturnFunc(func(g *jen.Group) {
				g.Func().
					Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("req").Interface()).
					Params(jen.Interface(), jen.Error()).
					BlockFunc(func(g *jen.Group) {
						g.Id("request").Op(":=").
							Id("req").Assert(jen.Qual(reqParam.Obj().Pkg().Path(), reqParam.Obj().Name()))
						g.Return(
							jen.Id("svc").Dot(f.Name()).
								Call(jen.Id("ctx"), jen.Id("request")))
					})
			})
		})
	return nil
}

func GenerateEndpointSet(file *jen.File, svc *types.Named) error {
	var (
		iface *types.Interface
		ok    bool
	)
	if iface, ok = svc.Underlying().(*types.Interface); !ok {
		return fmt.Errorf("%s underlying type is not interface", svc.Obj().Name())
	}
	file.Type().Id("EndpointSet").StructFunc(func(g *jen.Group) {
		for i := 0; i < iface.NumMethods(); i++ {
			method := iface.Method(i)
			if !method.Exported() {
				continue
			}
			g.Id(method.Name()+"Endpoint").Qual("github.com/go-kit/kit/endpoint", "Endpoint")
		}
	}).Line()
	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)
		signature := method.Type().(*types.Signature)
		params := signature.Params()
		results := signature.Results()
		receiverTyp := "EndpointSet"
		receiver := strings.ToLower(svc.Obj().Name()[:1])
		endpointFunc := method.Name() + "Endpoint"

		file.Func().
			Params(jen.Id(receiver).Id(receiverTyp)).
			Id(method.Name()).
			ParamsFunc(func(g *jen.Group) {
				for i := 0; i < params.Len(); i++ {
					param := params.At(i)
					name := fmt.Sprintf("arg%d", i)
					if param.Name() != "" {
						name = param.Name()
					}

					g.Id(name).Add(generateTypeCode(param.Type()))
				}
			}).
			ParamsFunc(func(g *jen.Group) {
				for i := 0; i < results.Len(); i++ {
					result := results.At(i)
					name := fmt.Sprintf("arg%d", i)
					if result.Name() != "" {
						name = result.Name()
					}

					g.Id(name).Add(generateTypeCode(result.Type()))
				}
			}).
			BlockFunc(func(g *jen.Group) {
				g.List(jen.Id("resp"), jen.Id("err")).
					Op(":=").
					Id(receiver).Dot(endpointFunc).CallFunc(func(g *jen.Group) {
					for i := 0; i < params.Len(); i++ {
						param := params.At(i)
						if param.Name() != "" {
							g.Id(param.Name())
						} else {
							g.Id(fmt.Sprintf("arg%d", i))
						}
					}
				})
				g.Return(jen.Id("resp").Assert(generateTypeCode(results.At(0).Type())), jen.Id("err"))
			}).Line().Line()
	}
	return nil
}
