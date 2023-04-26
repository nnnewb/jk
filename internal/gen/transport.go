package gen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"go/types"
)

func GenerateTransportLayerHTTP(f *jen.File, svc *types.Named) error {
	var (
		iface   = svc.Underlying().(*types.Interface)
		svcName = svc.Obj().Name()
		pkgPath = svc.Obj().Pkg().Path()
	)

	// makeHandlerFunc[REQ,RESP any](f func(context.Context, REQ) (RESP, error)) http.HandlerFunc {
	f.Func().Id("makeHandlerFunc").
		Types(jen.Id("REQ").Any(), jen.Id("RESP").Any()).
		Params(jen.Id("f").
			Func().
			Params(jen.Qual("context", "Context"), jen.Id("REQ")).
			Params(jen.Id("RESP"), jen.Error())).
		Params(jen.Qual("net/http", "HandlerFunc")).
		BlockFunc(func(g *jen.Group) {
			// return func(wr http.ResponseWriter, request *http.Request) {
			g.Return(
				jen.Func().
					Params(
						jen.Id("wr").Qual("net/http", "ResponseWriter"),
						jen.Id("request").Op("*").Qual("net/http", "Request"),
					)).
				BlockFunc(func(g *jen.Group) {
					// var payload REQ
					g.Var().Id("payload").Id("REQ")
					// err := json.NewDecoder(req.Body).Decode(&payload)
					g.Id("err").Op(":=").
						Qual("encoding/json", "NewDecoder").Call(jen.Id("request").Dot("Body")).
						Dot("Decode").Call(jen.Op("&").Id("payload"))
					// if err != nil {
					g.If(jen.Id("err").Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
						g.Panic(jen.Qual("github.com/juju/errors", "Errorf").Call(
							jen.Lit("unexpected unmarshal error %+v"), jen.Id("err")))
					}).Line()
					// resp, err := svc.Method(request.Context(), payload)
					g.List(jen.Id("resp"), jen.Id("err")).Op(":=").
						Id("f").Call(
						jen.Id("request").Dot("Context").Call(),
						jen.Id("payload"))
					// err = json.NewEncoder(wr).Encode(resp)
					g.Id("err").Op("=").
						Qual("encoding/json", "NewEncoder").Call(jen.Id("wr")).
						Dot("Encode").Call(jen.Id("resp"))
					// if err != nil {
					g.If(jen.Id("err").Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
						g.Panic(jen.Qual("github.com/juju/errors", "Errorf").Call(
							jen.Lit("unexpected unmarshal error %+v"), jen.Id("err")))
					})
				})
		}).Line()

	f.Func().Id("makeRemoteEndpoint").
		Types(jen.Id("REQ").Any(), jen.Id("RESP").Any()).
		Params(jen.Id("remoteUrl").String(), jen.Id("client").Op("*").Qual("net/http", "Client")).
		Qual("github.com/go-kit/kit/endpoint", "Endpoint").
		BlockFunc(func(g *jen.Group) {
			// return func(ctx context.Context, req interface{}) (interface{}, error) {
			g.Return(jen.Func().
				Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("req").Interface()).
				Params(jen.Interface(), jen.Error()).
				BlockFunc(func(g *jen.Group) {
					// buffer := bytes.NewBufferString("")
					g.Id("buffer").Op(":=").Qual("bytes", "NewBufferString").Call(jen.Lit(""))
					// err := json.NewEncoder(buffer).Encode(req)
					g.Id("err").Op(":=").
						Qual("encoding/json", "NewEncoder").Call(jen.Id("buffer")).
						Dot("Encode").Call(jen.Id("req"))
					// if err != nil {
					g.If(jen.Id("err").Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
						g.Panic(jen.Qual("github.com/juju/errors", "Errorf").Call(
							jen.Lit("unexpected marshal error %+v"), jen.Id("err")))
					}).Line()
					// request := http.NewRequestWithContext(ctx, http.MethodPost, remoteUrl, buffer)
					g.List(jen.Id("request"), jen.Id("err")).Op(":=").Qual("net/http", "NewRequestWithContext").
						Call(
							jen.Id("ctx"),
							jen.Qual("net/http", "MethodPost"),
							jen.Id("remoteUrl"),
							jen.Id("buffer"),
						)
					// if err != nil {
					g.If(jen.Id("err").Op("!=").Nil()).
						BlockFunc(func(g *jen.Group) {
							// return nil, errors.Trace(err)
							g.Return(jen.Nil(), jen.Qual("github.com/juju/errors", "Trace").Call(jen.Id("err")))
						}).Line()
					// response, err := client.Do(request)
					g.List(jen.Id("response"), jen.Id("err")).Op(":=").
						Id("client").Dot("Do").Call(jen.Id("request"))
					// if err != nil {
					g.If(jen.Id("err").Op("!=").Nil()).
						BlockFunc(func(g *jen.Group) {
							// return nil, errors.Trace(err)
							g.Return(jen.Nil(), jen.Qual("github.com/juju/errors", "Trace").Call(jen.Id("err")))
						}).Line()
					// defer response.Body.Close()
					g.Defer().Id("response").Dot("Body").Dot("Close").Call()
					// if response.StatusCode != http.StatusOK {
					g.If(jen.Id("response").Dot("StatusCode")).Op("!=").Qual("net/http", "StatusOK").
						BlockFunc(func(g *jen.Group) {
							// return nil, errors.Errorf("call remote endpoint failed, http status %d %s", response.StatusCode, response.Status)
							g.Return(
								jen.Nil(),
								jen.Qual("github.com/juju/errors", "Errorf").Call(
									jen.Lit("call remote endpoint failed, http status %d %s"),
									jen.Id("response").Dot("StatusCode"),
									jen.Id("response").Dot("Status"),
								),
							)
						}).Line()
					// var resp RESP
					g.Var().Id("resp").Id("RESP")
					// err = json.NewDecoder(response.Body).Decode(&resp)
					g.Id("err").Op("=").Qual("encoding/json", "NewDecoder").
						Call(jen.Id("response").Dot("Body")).
						Dot("Decode").Call(jen.Op("&").Id("resp"))
					// if err != nil {
					g.If(jen.Id("err").Op("!=").Nil()).
						BlockFunc(func(g *jen.Group) {
							// return nil, errors.Trace(err)
							g.Return(jen.Nil(), jen.Qual("github.com/juju/errors", "Trace").Call(jen.Id("err")))
						}).Line()
					// return resp, nil
					g.Return(jen.Id("resp"), jen.Nil())
				}))
		}).Line()

	f.Func().Id("NewClient").
		Params(jen.Id("host").String(), jen.Id("client").Op("*").Qual("net/http", "Client")).
		Params(jen.Qual(pkgPath, svcName)).
		BlockFunc(func(g *jen.Group) {
			for i := 0; i < iface.NumMethods(); i++ {
				method := iface.Method(i)
				g.Id("url"+method.Name()).Op(":=").Qual("net/url", "URL").Values(
					jen.Dict{
						jen.Id("Scheme"): jen.Lit("https"),
						jen.Id("Host"):   jen.Id("host"),
						jen.Id("Path"): jen.Lit(
							fmt.Sprintf(
								"/api/v1/%s/%s",
								strcase.ToKebab(svcName),
								strcase.ToKebab(method.Name()),
							)),
					},
				)
			}
			g.Return(jen.Id("EndpointSet").Values(jen.DictFunc(func(d jen.Dict) {
				for i := 0; i < iface.NumMethods(); i++ {
					method := iface.Method(i)
					signature := method.Type().(*types.Signature)
					params := signature.Params()
					results := signature.Results()
					d[jen.Id(method.Name()+"Endpoint")] = jen.
						Id("makeRemoteEndpoint").
						Types(generateTypeCode(params.At(1).Type()), generateTypeCode(results.At(0).Type())).
						Call(
							jen.Id("url"+method.Name()).Dot("String").Call(),
							jen.Id("client"),
						)
				}
			})))
		}).Line()

	f.Func().Id("Register").
		Params(jen.Id("svc").Qual(pkgPath, svcName), jen.Id("m").Op("*").Qual("net/http", "ServeMux")).
		Params(jen.Op("*").Qual("net/http", "ServeMux")).
		BlockFunc(func(g *jen.Group) {

			// TODO: 需要支持可选的 API 路径配置
			for i := 0; i < iface.NumMethods(); i++ {
				method := iface.Method(i)

				apiEndpointName := strcase.ToKebab(method.Name())
				apiServiceName := strcase.ToKebab(svcName)
				g.Id("m").Dot("HandleFunc").Call(
					jen.Lit(fmt.Sprintf("/api/v1/%s/%s", apiServiceName, apiEndpointName)),
					jen.Id("makeHandlerFunc").Call(jen.Id("svc").Dot(method.Name())),
				)
			}
			g.Return(jen.Id("m"))
		}).Line()
	return nil
}
