package genrpc

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/nnnewb/jk/pkg/gen/gencore"
	"github.com/nnnewb/jk/pkg/gen/genutils"
)

func genBindingStruct(req *gencore.GenRequest) jen.Code {
	// type XxxBinding struct { svc Xxx }
	return jen.Type().Id(req.Svc.Name() + "Binding").StructFunc(func(g *jen.Group) {
		// endpoint set
		for i := 0; i < req.GetServiceInterface().NumMethods(); i++ {
			method := req.GetServiceInterface().Method(i)
			if genutils.IsOmitMethod(method, req.GetContextType(), req.GetErrorType()) {
				continue
			}

			g.Id(method.Name()+"Endpoint").Qual("github.com/go-kit/kit/endpoint", "Endpoint")
		}
	}).Line()
}

func genBindingStructConstructor(req *gencore.GenRequest) jen.Code {
	return jen.
		Func().Id("New" + req.Svc.Name() + "Binding").
		Params(jen.Id("svc").Add(genutils.Qual(req.Svc.Type()))).
		Params(jen.Op("*").Id(req.Svc.Name() + "Binding")).
		Block(
			jen.Return(
				jen.Op("&").Id(req.Svc.Name() + "Binding").Values(jen.DictFunc(func(d jen.Dict) {
					for i := 0; i < req.GetServiceInterface().NumMethods(); i++ {
						method := req.GetServiceInterface().Method(i)
						if genutils.IsOmitMethod(method, req.GetContextType(), req.GetErrorType()) {
							continue
						}

						d[jen.Id(method.Name()+"Endpoint")] = jen.Qual(req.Svc.Pkg().Path()+"/endpoint", "Make"+method.Name()+"Endpoint").Call(jen.Id("svc"))
					}
				})),
			),
		)
}

func GenerateBindings(data *gencore.PluginData) error {
	req := data.Request
	gf := data.GetOrCreateFile("transport/netrpc/bindings_gen.go")
	f := jen.NewFile("netrpc")
	f.ImportAlias("github.com/go-kit/kit/endpoint", "kitendpoint")

	// generate type XxxBinding strcut {}
	f.Add(genBindingStruct(req)).Line()

	// generate func NewXxxBinding() *XxxBinding {}
	f.Add(genBindingStructConstructor(req)).Line()

	// generate func (*XxxBinding) Method(req Req, reply *Reply) error {}
	for i := 0; i < req.GetServiceInterface().NumMethods(); i++ {
		method := req.GetServiceInterface().Method(i)
		if genutils.IsOmitMethod(method, req.GetContextType(), req.GetErrorType()) {
			continue
		}

		f.Func().
			Params(jen.Id("b").Id(fmt.Sprintf("%sBinding", req.Svc.Name()))).
			Id(method.Name()).
			Params(
				jen.Id("req").Qual(req.Svc.Pkg().Path()+"/endpoint", fmt.Sprintf("%sRequest", method.Name())),
				jen.Id("response").Op("*").Qual(req.Svc.Pkg().Path()+"/endpoint", fmt.Sprintf("%sResponse", method.Name())),
			).
			Params(jen.Error()).
			BlockFunc(func(g *jen.Group) {
				// resp, err := b.svc.XxxEndpoint(context.Background(), req)
				g.List(jen.Id("resp"), jen.Id("err")).Op(":=").Add(jen.Id("b").Dot(method.Name()+"Endpoint")).Call(jen.Qual("context", "Background").Call(), jen.Id("req"))
				// if err != nil { return err }
				g.If(jen.Id("err").Op("!=").Nil()).Block(jen.Return(jen.Id("err")))
				// (*resp) = *resp.(*XxxResponse)
				g.Op("*").Id("response").Op("=").Op("*").Id("resp").Assert(jen.Op("*").Qual(req.Svc.Pkg().Path()+"/endpoint", fmt.Sprintf("%sResponse", method.Name())))
				// return nil
				g.Return(jen.Nil())
			}).Line()
	}

	return f.Render(gf)
}
