package gen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/juju/errors"
	"go/types"
	"strings"
)

// generateEndpointFactory generates Endpoint factory function code from function signature.
// The function signature must meet the following requirements:
//   - The first parameter must be context.Context.
//   - The second parameter must be an exported struct that is serializable.
//   - The first return value must be an exported struct or slice of structs that is serializable.
//     When returning a slice of structs, the response content will be encapsulated in the Items field.
//     When returning a single struct, the response content will be encapsulated in the Data field.
//   - The second return value must be of type error. The result of the Error method of this interface
//     will be used as the response msg.
func generateEndpointFactory(file *jen.File, svc *types.Named, f *types.Func) error {
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

func generateEndpointSet(file *jen.File, svc *types.Named) error {
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

// GenerateEndpoints generates endpoint factory for a given service
func GenerateEndpoints(f *jen.File, svc *types.Named) error {
	var (
		iface *types.Interface
		ok    bool
	)

	// Get the underlying interface of the named type
	underlying := svc.Underlying()

	// Check if the underlying type is an interface
	if iface, ok = underlying.(*types.Interface); !ok {
		// If the underlying type is not an interface, return an error
		return fmt.Errorf("%s is not an interface", svc.Obj().Name())
	}

	// Iterate over all the methods of the interface
	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)
		// Check if the method is exported
		if method.Exported() {
			err := generateEndpointFactory(f, svc, method)
			if err != nil {
				return errors.Annotatef(err, "generate endpoint factory func %s failed, error %+v", method, err)
			}
		}
	}
	err := generateEndpointSet(f, svc)
	if err != nil {
		return err
	}
	return nil
}
