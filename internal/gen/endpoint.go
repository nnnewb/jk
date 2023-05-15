package gen

import (
	"go/types"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/juju/errors"
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
		reqType := params.At(1).Type().(*types.Named)
		respType := results.At(0).Type().(*types.Named)

		file.Func().
			Params(jen.Id(receiver).Id(receiverTyp)).
			Id(method.Name()).
			Params(
				jen.Id("ctx").Qual("context", "Context"),
				jen.Id("req").Qual(reqType.Obj().Pkg().Path(), reqType.Obj().Name())).
			Params(
				jen.Qual(respType.Obj().Pkg().Path(), respType.Obj().Name()),
				jen.Error()).
			BlockFunc(func(g *jen.Group) {
				g.List(jen.Id("resp"), jen.Err()).
					Op(":=").
					Id(receiver).Dot(endpointFunc).Call(jen.Id("ctx"), jen.Id("req")).
					Line()
				// if err != nil { return RESP{}, err }
				g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Qual(respType.Obj().Pkg().Path(), respType.Obj().Name()).Values(), jen.Err()))
				// return resp.(RESP), nil
				g.Return(
					jen.Id("resp").Assert(jen.Qual(respType.Obj().Pkg().Path(), respType.Obj().Name())),
					jen.Nil())
			}).Line()
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
