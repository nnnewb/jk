package gen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/juju/errors"
	"go/types"
	"strings"
)

func generateEndpointFactory(file *jen.File) {
	file.Func().
		Id("makeEndpointFromFunc").
		Types(jen.List(jen.Id("REQ"), jen.Id("RESP").Any())).
		Params(jen.Id("f").Func().Params(jen.Qual("context", "Context"), jen.Id("REQ")).Params(jen.Id("RESP"), jen.Error())).
		Qual("github.com/go-kit/kit/endpoint", "Endpoint").
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Func().
				Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("r").Interface()).
				Params(jen.Interface(), jen.Error()).
				BlockFunc(func(g *jen.Group) {
					g.Id("req").Op(":=").Id("r").Assert(jen.Id("REQ"))
					g.Return(jen.Id("f").Call(jen.Id("ctx"), jen.Id("req")))
				}),
			)
		})
}

func generateEndpointSet(file *jen.File, svc *types.Named) {
	iface := svc.Underlying().(*types.Interface)
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
		if !method.Exported() {
			continue
		}

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

	file.Func().
		Id("NewEndpointSet").
		Params(jen.Id("svc").Qual(svc.Obj().Pkg().Path(), svc.Obj().Name())).
		Id("EndpointSet").
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Id("EndpointSet").Values(jen.DictFunc(func(d jen.Dict) {
				for i := 0; i < iface.NumMethods(); i++ {
					method := iface.Method(i)
					if !method.Exported() {
						continue
					}

					endpointName := method.Name() + "Endpoint"
					d[jen.Id(endpointName)] = jen.Id("makeEndpointFromFunc").Call(jen.Id("svc").Dot(method.Name()))
				}
			})))
		}).Line()

	file.Func().
		Params(jen.Id("s").Id("EndpointSet")).
		Id("With").
		Params(
			jen.Id("outer").Qual("github.com/go-kit/kit/endpoint", "Middleware"),
			jen.Id("others").Op("...").Qual("github.com/go-kit/kit/endpoint", "Middleware"),
		).
		Id("EndpointSet").
		BlockFunc(func(g *jen.Group) {
			g.ReturnFunc(func(g *jen.Group) {
				g.Id("EndpointSet").Values(jen.DictFunc(func(d jen.Dict) {
					for i := 0; i < iface.NumMethods(); i++ {
						method := iface.Method(i)
						if !method.Exported() {
							continue
						}

						endpointName := method.Name() + "Endpoint"
						d[jen.Id(endpointName)] = jen.Qual("github.com/go-kit/kit/endpoint", "Chain").Call(
							jen.Id("outer"),
							jen.Id("others").Op("..."),
						).Call(jen.Id("s").Dot(endpointName))
					}
				}))
			})
		})
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
		return errors.Errorf("%s is not an interface", svc.Obj().Name())
	}

	for i := 0; i < iface.NumMethods(); i++ {
		method := iface.Method(i)
		if !method.Exported() {
			continue
		}
		signature := method.Type().(*types.Signature)

		if err := checkParams(signature.Params()); err != nil {
			return errors.Annotatef(err, "check method signature: %s", method.FullName())
		}

		if err := checkResults(signature.Results()); err != nil {
			return errors.Annotatef(err, "check method signature: %s", method.FullName())
		}
	}

	generateEndpointFactory(f)
	generateEndpointSet(f, svc)
	return nil
}
