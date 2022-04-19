package service

import (
	"fmt"
	"go/types"
	"log"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/nnnewb/jk/pkg/generator/driver"
	"github.com/nnnewb/jk/pkg/generator/utils"
)

func init() {
	driver.RegisterServiceGenDriver("", defaultServiceGenerator{})
}

type defaultServiceGenerator struct {
	req *driver.GenerateRequest
}

func (d defaultServiceGenerator) GenerateService(req *driver.GenerateRequest) error {
	d.req = req
	gf := req.GenFile("endpoint/endpoint_gen.go")
	f := jen.NewFile("endpoint")
	f.HeaderComment("This file is generated by jk, DO NOT EDIT.")

	var ctxType, errType types.Type
	req.TypesImporter.Import("context")

	dstPkg := types.NewPackage(req.Pkg.Path()+"/endpoint", "endpoint")
	for _, method := range utils.FilterCorrespondPublicMethod(req.Svc) {
		methodName := method.Name()
		methodType := method.Type().(*types.Signature)

		log.Printf("generate %s", methodName)
		// generate code
		// 1. create request/response struct
		// 2. create makeXxxEndpoint function to create endpoint for method
		reqStruct := utils.NewParamsStruct(method)
		respStruct := utils.NewResultsStruct(method)
		f.Add(utils.GenStruct(fmt.Sprintf("%sRequest", methodName), dstPkg, reqStruct)).Line()
		f.Add(utils.GenStruct(fmt.Sprintf("%sResponse", methodName), dstPkg, respStruct)).Line()

		// create endpoint constructor for method
		f.Func().Id(fmt.Sprintf("make%sEndpoint", methodName)).
			Params(jen.Id("svc").Qual(req.Pkg.Path(), req.SvcName)).
			Params(jen.Qual("github.com/go-kit/kit/endpoint", "Endpoint")).
			Block(
				jen.Return(jen.Func().
					Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("req").Interface()).
					Params(jen.Interface(), jen.Error()).
					BlockFunc(func(g *jen.Group) {
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
					}),
				),
			).
			Line()

	}

	return f.Render(gf.Writer)
}