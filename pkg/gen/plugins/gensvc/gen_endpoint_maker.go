package gensvc

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/nnnewb/jk/pkg/gen/gencore"
	"github.com/nnnewb/jk/pkg/gen/genutils"
)

func genEndpointMaker(data *gencore.PluginData) error {
	gf := data.GetOrCreateFile("endpoint/endpoint_gen.go")
	req := data.Request
	f := jen.NewFile("endpoint")

	for i := 0; i < req.GetServiceInterface().NumMethods(); i++ {
		method := req.GetServiceInterface().Method(i)
		if genutils.IsOmitMethod(method, req.GetContextType(), req.GetErrorType()) {
			continue
		}

		methodName := method.Name()
		methodType := method.Type().(*types.Signature)

		// func MakeXxxEndpoint(svc service.ServiceInterface) endpoint.Endpoint {
		f.Func().Id(fmt.Sprintf("Make%sEndpoint", methodName)).Params(jen.Id("svc").Add(genutils.Qual(req.Svc.Type()))).Params(jen.Qual("github.com/go-kit/kit/endpoint", "Endpoint")).
			// return func(ctx context.Context, req interface{}) (interface{}, error) {
			Block(jen.Return(jen.Func().Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("req").Interface()).Params(jen.Interface(), jen.Error()).BlockFunc(func(g *jen.Group) {
				// request := req.(*XxxRequest)
				g.Id("request").Op(":=").Id("req").Assert(jen.Op("*").Id(fmt.Sprintf("%sRequest", methodName)))

				// result1,...,err := svc.Xxx(ctx, param1, ...)
				g.
					ListFunc(func(g *jen.Group) {
						for i := 0; i < methodType.Results().Len(); i++ {
							result := methodType.Results().At(i)
							g.Id(strcase.ToLowerCamel(result.Name()))
						}
					}).
					Op(":=").
					Id("svc").Dot(methodName).
					CallFunc(func(g *jen.Group) {
						g.Id("ctx")
						for i := 1; i < methodType.Params().Len(); i++ {
							param := methodType.Params().At(i)
							g.Id("request").Dot(strcase.ToCamel(param.Name()))
						}
					})

				// if err != nil { return nil, err }
				g.If(jen.Id("err").Op("!=").Nil()).
					Block(jen.Return(jen.Nil(), jen.Id("err")))

				// resp := &XxxResponse{Result1: result1, ...}
				g.Id("resp").Op(":=").Op("&").Id(fmt.Sprintf("%sResponse", methodName)).
					ValuesFunc(func(g *jen.Group) {
						g.Add(jen.DictFunc(func(dict jen.Dict) {
							for i := 0; i < methodType.Results().Len()-1; i++ {
								result := methodType.Results().At(i)
								dict[jen.Id(strcase.ToCamel(result.Name()))] = jen.Id(strcase.ToLowerCamel(result.Name()))
							}
						}))
					})

				// return resp, nil
				g.Return(jen.Id("resp"), jen.Nil())
			}))).
			Line()
	}

	return f.Render(gf)
}
