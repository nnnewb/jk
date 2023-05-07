package gen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"go/types"
)

func generateContextKeyDeclaration(f *jen.File) {
	f.Type().Id("httpTransportRequestKeyType").Struct()
	f.Var().Id("httpTransportRequestKey").Id("httpTransportRequestKeyType")
}

func generateMakeHandlerFunc(f *jen.File) {
	// makeHandlerFunc[REQ,RESP any](f func(context.Context, REQ) (RESP, error)) http.HandlerFunc {
	f.Func().Id("makeHandlerFunc").
		Types(jen.Id("REQ").Any(), jen.Id("RESP").Any()).
		Params(jen.Id("f").
			Func().
			Params(jen.Qual("context", "Context"), jen.Id("REQ")).
			Params(jen.Id("RESP"), jen.Error())).
		Qual("net/http", "HandlerFunc").
		BlockFunc(func(g *jen.Group) {
			// return func(wr http.ResponseWriter, request *http.Request) {
			g.Return(
				jen.Func().
					Params(
						jen.Id("wr").Qual("net/http", "ResponseWriter"),
						jen.Id("request").Op("*").Qual("net/http", "Request"),
					)).
				BlockFunc(func(g *jen.Group) {
					// defer request.Body.Close()
					g.Defer().Id("request").Dot("Body").Dot("Close").Call()
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
					// ctx := context.WithValue(request.Context(), httpTransportRequestKey, req)
					g.Id("ctx").Op(":=").Qual("context", "WithValue").
						Call(
							jen.Id("request").Dot("Context").Call(),
							jen.Id("httpTransportRequestKey"),
							jen.Id("request"))
					// resp, err := svc.Method(request.Context(), payload)
					g.List(jen.Id("resp"), jen.Id("err")).Op(":=").
						Id("f").Call(
						jen.Id("ctx"),
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
}

func generateGetRequestFromContext(j *jen.File) {
	j.Commentf("// GetRequestFromContext get *http.Request from context.Context object. if no request associated, return nil.")
	j.Func().Id("GetRequestFromContext").
		Params(jen.Id("ctx").Qual("context", "Context")).
		Params(jen.Op("*").Qual("net/http", "Request")).
		BlockFunc(func(g *jen.Group) {
			g.List(jen.Id("req"), jen.Id("_")).Op(":=").
				Id("ctx").Dot("Value").Call(jen.Id("httpTransportRequestKey")).
				Assert(jen.Op("*").Qual("net/http", "Request"))
			g.Return(jen.Id("req"))
		}).Line()
}

func generateGetOperationNameFromContext(f *jen.File) {
	// TODO: useful with tracing
	f.Func().Id("GetOperationNameFromContext").
		Params(jen.Id("ctx").Qual("context", "Context"), jen.Id("defaultOperationName").String()).
		String().
		BlockFunc(func(g *jen.Group) {
			g.Id("req").Op(":=").Id("GetRequestFromContext").Call(jen.Id("ctx"))
			g.If(jen.Id("req").Op("!=").Nil()).
				BlockFunc(func(g *jen.Group) {
					g.Return(jen.Id("req").Dot("URL").Dot("Path"))
				})
			g.Return(jen.Lit("no http request associated, can not get operation name"))
		}).Line()
}

func generateMakeRemoteEndpoint(f *jen.File) {
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
}

func generateNewClient(f *jen.File, svc *types.Named) {
	var (
		iface   = svc.Underlying().(*types.Interface)
		svcName = svc.Obj().Name()
		pkgPath = svc.Obj().Pkg().Path()
	)
	f.Func().Id("NewClient").
		Params(jen.Id("host").String(), jen.Id("client").Op("*").Qual("net/http", "Client")).
		Params(jen.Qual(pkgPath, svcName)).
		BlockFunc(func(g *jen.Group) {
			for i := 0; i < iface.NumMethods(); i++ {
				method := iface.Method(i)
				if !method.Exported() {
					continue
				}
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
					if !method.Exported() {
						continue
					}
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
}

func generateEmbedSwaggerJSON(f *jen.File) {
	// //go:embed swagger.json
	f.Commentf("//go:embed swagger.json")
	// var swagger embed.FS
	f.Var().Id("swagger").Qual("embed", "FS")
}

func generateRegister(f *jen.File, svc *types.Named) {
	var (
		iface   = svc.Underlying().(*types.Interface)
		svcName = svc.Obj().Name()
		pkgPath = svc.Obj().Pkg().Path()
	)

	f.Func().Id("Register").
		Params(jen.Id("svc").Qual(pkgPath, svcName), jen.Id("m").Op("*").Qual("github.com/julienschmidt/httprouter", "Router")).
		BlockFunc(func(g *jen.Group) {
			// m.Handle("/swagger/service/spec/", http.FileServer(http.FS(swagger)))
			g.Id("m").Dot("Handler").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.Lit(fmt.Sprintf("/swagger/%s/spec/*rest", strcase.ToKebab(svcName))),

				jen.Qual("net/http", "StripPrefix").Call(
					jen.Lit(fmt.Sprintf("/swagger/%s/spec/", strcase.ToKebab(svcName))),
					jen.Qual("net/http", "FileServer").
						Call(jen.Qual("net/http", "FS").Call(jen.Id("swagger")))),
			)
			// m.Handle("/swagger/service/swagger-ui/*", httpSwagger.Handler(httpSwagger.URL("/swagger/service/swagger.json")))
			g.Id("m").Dot("Handler").Call(
				jen.Qual("net/http", "MethodGet"),
				jen.Lit(fmt.Sprintf("/swagger/%s/swagger-ui/*rest", strcase.ToKebab(svcName))),

				jen.Qual("github.com/swaggo/http-swagger/v2", "Handler").
					Call(jen.Qual("github.com/swaggo/http-swagger/v2", "URL").
						Call(jen.Lit(fmt.Sprintf("/swagger/%s/spec/swagger.json", strcase.ToKebab(svcName))))),
			)

			// TODO: 需要支持可选的 API 路径配置
			for i := 0; i < iface.NumMethods(); i++ {
				method := iface.Method(i)
				if !method.Exported() {
					continue
				}

				apiEndpointName := strcase.ToKebab(method.Name())
				apiServiceName := strcase.ToKebab(svcName)
				g.Id("m").Dot("Handler").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit(fmt.Sprintf("/api/v1/%s/%s", apiServiceName, apiEndpointName)),

					jen.Id("makeHandlerFunc").Call(jen.Id("svc").Dot(method.Name())),
				)
			}
		}).Line()
}

func generateTraceServer(f *jen.File, svc *types.Named) {
	panic("not implemented") // TODO: implement
}

func generateTraceClient(f *jen.File, svc *types.Named) {
	panic("not implemented") // TODO: implement
}

func GenerateHTTPTransport(f *jen.File, svc *types.Named) {
	generateContextKeyDeclaration(f)
	generateGetRequestFromContext(f)
	generateMakeHandlerFunc(f)
	generateMakeRemoteEndpoint(f)
	generateNewClient(f, svc)
	generateEmbedSwaggerJSON(f)
	generateRegister(f, svc)
}
