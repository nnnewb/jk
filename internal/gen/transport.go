package gen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"go/types"
)

func generateClientSet(f *jen.File, svc *types.Named) {
	iface := svc.Underlying().(*types.Interface)
	f.Type().Id("HTTPClientSet").StructFunc(func(g *jen.Group) {
		for i := 0; i < iface.NumMethods(); i++ {
			method := iface.Method(i)
			if !method.Exported() {
				continue
			}

			g.Id(method.Name()+"Client").Qual("github.com/go-kit/kit/transport/http", "Client")
		}
	}).Line()
}

func generateServerSet(f *jen.File, svc *types.Named) {
	iface := svc.Underlying().(*types.Interface)
	f.Type().Id("HTTPServerSet").StructFunc(func(g *jen.Group) {
		for i := 0; i < iface.NumMethods(); i++ {
			method := iface.Method(i)
			if !method.Exported() {
				continue
			}

			g.Id(method.Name()+"Server").Qual("github.com/go-kit/kit/transport/http", "Server")
		}
	}).Line()
}

func generateHTTPJSONRequestDecoder(f *jen.File) {
	// func httpJSONRequestDecoder[T any](ctx context.Context, req *http.Request) (any, error) {
	f.Func().
		Id("httpJSONRequestDecoder").
		Types(jen.Id("T").Any()).
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("req").Op("*").Qual("net/http", "Request")).
		Params(
			jen.Any(),
			jen.Error()).
		BlockFunc(func(g *jen.Group) {
			// var request T
			g.Var().Id("request").Id("T")
			// err := json.NewDecoder(req.Body).Decode(&request)
			g.Err().Op(":=").Qual("encoding/json", "NewDecoder").Call(jen.Id("req").Dot("Body")).Dot("Decode").
				Call(jen.Op("&").Id("request"))
			// if err != nil { return nil, err }
			g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err()))
			// return &request,nil
			g.Return(jen.Op("&").Id("request"), jen.Nil())
		}).Line()
}

func generateHTTPJSONResponseDecoder(f *jen.File) {
	// func httpJSONResponseDecoder[T any](ctx context.Context, req *http.Response) (any, error) {
	f.Func().
		Id("httpJSONResponseDecoder").
		Types(jen.Id("T").Any()).
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("resp").Op("*").Qual("net/http", "Response")).
		Params(
			jen.Any(),
			jen.Error()).
		BlockFunc(func(g *jen.Group) {
			// var response T
			g.Var().Id("response").Id("T")
			// defer resp.Body.Close()
			g.Defer().Id("resp").Dot("Body").Dot("Close").Call()
			// err := json.NewDecoder(resp.Body).Decode(&response)
			g.Err().Op(":=").Qual("encoding/json", "NewDecoder").Call(jen.Id("resp").Dot("Body")).Dot("Decode").
				Call(jen.Op("&").Id("response"))
			// if err != nil { return nil, err }
			g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err()))
			// return &response, nil
			g.Return(jen.Op("&").Id("response"), jen.Nil())
		}).Line()
}

func generateHTTPJSONRequestEncoder(f *jen.File) {
	// func httpJSONRequestEncoder[T any](ctx context.Context, req *http.Request, request any) (any, error) {
	f.Func().
		Id("httpJSONRequestEncoder").
		Types(jen.Id("T").Any()).
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("req").Op("*").Qual("net/http", "Request"),
			jen.Id("request").Any()).
		Error().
		BlockFunc(func(g *jen.Group) {
			// var buffer bytes.Buffer
			g.Var().Id("buffer").Qual("bytes", "Buffer")
			// err := json.NewEncoder(&buffer).Encode(request)
			g.Err().Op(":=").
				Qual("encoding/json", "NewEncoder").Call(jen.Op("&").Id("buffer")).
				Dot("Encode").Call(jen.Id("request"))
			// if err != nil { return err }
			g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Err()))
			// req.Body = &buffer
			g.Id("req").Dot("Body").Op("=").Qual("io", "NopCloser").Call(jen.Op("&").Id("buffer"))
			// return nil
			g.Return(jen.Nil())
		}).Line()
}

func generateHTTPJSONResponseEncoder(f *jen.File) {
	// func httpJSONResponseEncoder[T any](ctx context.Context, req *http.Request) (any, error) {
	f.Func().
		Id("httpJSONResponseEncoder").
		Types(jen.Id("T").Any()).
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("wr").Qual("net/http", "ResponseWriter"),
			jen.Id("resp").Any()).
		Error().
		BlockFunc(func(g *jen.Group) {
			// return json.NewEncoder(wr).Encode(resp)
			g.Return(
				jen.Qual("encoding/json", "NewEncoder").Call(jen.Id("wr")).
					Dot("Encode").Call(jen.Id("resp")))
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

func generateNewHTTPClient(f *jen.File, svc *types.Named) {
	var (
		iface   = svc.Underlying().(*types.Interface)
		svcName = svc.Obj().Name()
		pkgPath = svc.Obj().Pkg().Path()
	)
	f.Func().Id("NewHTTPClient").
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
	)

	f.Func().Id("Register").
		Params(jen.Id("svc").Id("EndpointSet"), jen.Id("m").Op("*").Qual("github.com/julienschmidt/httprouter", "Router")).
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
				signature := method.Type().(*types.Signature)
				reqType := signature.Params().At(1).Type().(*types.Named)
				respType := signature.Results().At(0).Type().(*types.Named)

				apiEndpointName := strcase.ToKebab(method.Name())
				apiServiceName := strcase.ToKebab(svcName)
				g.Id("m").Dot("Handler").Call(
					jen.Qual("net/http", "MethodPost"),
					jen.Lit(fmt.Sprintf("/api/v1/%s/%s", apiServiceName, apiEndpointName)),
					jen.Qual("github.com/go-kit/kit/transport/http", "NewServer").Call(
						jen.Id("svc").Dot(method.Name()+"Endpoint"),
						jen.Id("httpJSONRequestDecoder").Types(jen.Qual(reqType.Obj().Pkg().Path(), reqType.Obj().Name())),
						jen.Id("httpJSONResponseEncoder").Types(jen.Qual(respType.Obj().Pkg().Path(), respType.Obj().Name())),
					))
			}
		}).Line()
}

func GenerateHTTPTransport(f *jen.File, svc *types.Named) {
	generateClientSet(f, svc)
	generateServerSet(f, svc)
	generateHTTPJSONRequestDecoder(f)
	generateHTTPJSONRequestEncoder(f)
	generateHTTPJSONResponseEncoder(f)
	generateHTTPJSONResponseDecoder(f)
	generateMakeRemoteEndpoint(f)
	generateNewHTTPClient(f, svc)
	generateEmbedSwaggerJSON(f)
	generateRegister(f, svc)
}
