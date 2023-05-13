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

			g.Id(method.Name()+"Client").Op("*").Qual("github.com/go-kit/kit/transport/http", "Client")
		}
	}).Line()

	// func NewHTTPClientSet(scheme, host string, port int, options ...http.ClientOptions) HTTPClientSet {
	f.Func().
		Id("NewHTTPClientSet").
		Params(
			jen.Id("scheme"),
			jen.Id("host").String(),
			jen.Id("port").Int(),
			jen.Id("options").Op("...").Qual("github.com/go-kit/kit/transport/http", "ClientOption")).
		Id("HTTPClientSet").
		BlockFunc(func(g *jen.Group) {
			// return HTTPClientSet{
			g.Return(jen.Id("HTTPClientSet")).Values(jen.DictFunc(func(d jen.Dict) {
				for i := 0; i < iface.NumMethods(); i++ {
					method := iface.Method(i)
					if !method.Exported() {
						continue
					}
					signature := method.Type().(*types.Signature)
					respType := signature.Results().At(0).Type().(*types.Named)

					// XXXClient: http.NewClient(
					//   http.MethodPost,
					//   url.URL{
					//     Scheme: scheme,
					//     Host: fmt.Sprintf("%s:%d", host, port),
					//     Path: "/api/v1/SERVICE/ENDPOINT",
					//   },
					//   httpJSONRequestEncoder[REQ],
					//   httpJSONResponseDecoder[RESP],
					//   options...,
					// )
					d[jen.Id(method.Name()+"Client")] = jen.
						Qual("github.com/go-kit/kit/transport/http", "NewClient").
						Call(
							jen.Line().Qual("net/http", "MethodPost"),
							jen.Line().Op("&").Qual("net/url", "URL").Values(jen.DictFunc(func(d jen.Dict) {
								d[jen.Id("Scheme")] = jen.Id("scheme")
								d[jen.Id("Host")] = jen.Qual("fmt", "Sprintf").Call(
									jen.Lit("%s:%d"),
									jen.Id("host"),
									jen.Id("port"))
								d[jen.Id("Path")] = jen.Lit(
									fmt.Sprintf(
										"/api/v1/%s/%s",
										strcase.ToKebab(svc.Obj().Name()),
										strcase.ToKebab(method.Name())))
							})),
							jen.Line().Qual("github.com/go-kit/kit/transport/http", "EncodeJSONRequest"),
							jen.Line().Id("httpJSONResponseDecoder").Types(jen.Qual(respType.Obj().Pkg().Path(), respType.Obj().Name())),
							jen.Line().Id("options").Op("..."),
						)
				}
			}))
		}).Line()

	// func (s HTTPClientSet) EndpointSet() EndpointSet {
	f.Func().
		Params(jen.Id("s").Id("HTTPClientSet")).
		Id("EndpointSet").
		Params().
		Id("EndpointSet").
		BlockFunc(func(g *jen.Group) {
			// return EndpointSet{
			g.Return(jen.Id("EndpointSet").Values(jen.DictFunc(func(d jen.Dict) {
				for i := 0; i < iface.NumMethods(); i++ {
					// XXXEndpoint: s.XXXClient.Endpoint(),
					method := iface.Method(i)
					if !method.Exported() {
						continue
					}

					d[jen.Id(method.Name()+"Endpoint")] = jen.
						Id("s").
						Dot(method.Name() + "Client").
						Dot("Endpoint").Call()
				}
			})))
		}).Line()
}

func generateServerSet(f *jen.File, svc *types.Named) {
	iface := svc.Underlying().(*types.Interface)
	svcName := svc.Obj().Name()

	// type HTTPServerSet {
	f.Type().Id("HTTPServerSet").StructFunc(func(g *jen.Group) {
		for i := 0; i < iface.NumMethods(); i++ {
			method := iface.Method(i)
			if !method.Exported() {
				continue
			}

			// XXXServer: *http.Server,
			g.Id(method.Name()+"Server").Op("*").Qual("github.com/go-kit/kit/transport/http", "Server")
		}
	}).Line()

	// func NewHTTPServerSet(endpointSet EndpointSet, options ...http.ServerOption) HTTPServerSet {
	f.Func().
		Id("NewHTTPServerSet").
		Params(
			jen.Id("endpointSet").Id("EndpointSet"),
			jen.Id("options").Op("...").Qual("github.com/go-kit/kit/transport/http", "ServerOption")).
		Id("HTTPServerSet").
		BlockFunc(func(g *jen.Group) {
			// return HTTPServerSet{
			g.Return(jen.Id("HTTPServerSet").Values(jen.DictFunc(func(d jen.Dict) {
				for i := 0; i < iface.NumMethods(); i++ {
					method := iface.Method(i)
					if !method.Exported() {
						continue
					}
					signature := method.Type().(*types.Signature)
					reqType := signature.Params().At(1).Type().(*types.Named)

					// XXXServer: khttp.NewServer(
					//   endpointSet.XXXEndpoint,
					//   httpJSONRequestDecoder[REQ],
					//   httpJSONResponseEncoder[RESP],
					//   append(options, khttp.ServerBefore(khttp.PopulateRequestContext))...,
					// )
					d[jen.Id(method.Name()+"Server")] = jen.Qual("github.com/go-kit/kit/transport/http", "NewServer").
						Call(
							jen.Line().Id("endpointSet").Dot(method.Name()+"Endpoint"),
							jen.Line().Id("httpJSONRequestDecoder").Types(jen.Qual(reqType.Obj().Pkg().Path(), reqType.Obj().Name())),
							jen.Line().Qual("github.com/go-kit/kit/transport/http", "EncodeJSONResponse"),
							jen.Line().Append(
								jen.Id("options"),
								jen.Qual("github.com/go-kit/kit/transport/http", "ServerBefore").
									Call(jen.Qual("github.com/go-kit/kit/transport/http", "PopulateRequestContext")),
							).Op("..."))
				}
			})))
		}).Line()

	// func (s HTTPServerSet) Handler() http.Handler {
	f.Func().
		Params(jen.Id("s").Id("HTTPServerSet")).
		Id("Handler").
		Params().
		Qual("net/http", "Handler").
		BlockFunc(func(g *jen.Group) {
			// var m *httprouter.Router
			g.Var().Id("m").Op("=").Id("new").Call(jen.Qual("github.com/julienschmidt/httprouter", "Router"))
			// m.Handle("/swagger/service/spec/", http.FileServer(http.FS(swagger)))
			g.Id("m").Dot("Handler").Call(
				jen.Line().Qual("net/http", "MethodGet"),
				jen.Line().Lit(fmt.Sprintf("/swagger/%s/spec/*rest", strcase.ToKebab(svcName))),
				jen.Line().Qual("net/http", "StripPrefix").Call(
					jen.Lit(fmt.Sprintf("/swagger/%s/spec/", strcase.ToKebab(svcName))),
					jen.Qual("net/http", "FileServer").
						Call(jen.Qual("net/http", "FS").Call(jen.Id("swagger")))),
			)
			// m.Handle("/swagger/service/swagger-ui/*", httpSwagger.Handler(httpSwagger.URL("/swagger/service/swagger.json")))
			g.Id("m").Dot("Handler").Call(
				jen.Line().Qual("net/http", "MethodGet"),
				jen.Line().Lit(fmt.Sprintf("/swagger/%s/swagger-ui/*rest", strcase.ToKebab(svcName))),
				jen.Line().Qual("github.com/swaggo/http-swagger/v2", "Handler").
					Call(jen.Qual("github.com/swaggo/http-swagger/v2", "URL").
						Call(jen.Lit(fmt.Sprintf("/swagger/%s/spec/swagger.json", strcase.ToKebab(svcName))))),
			)

			// TODO: 需要支持可选的 API 路径配置
			for i := 0; i < iface.NumMethods(); i++ {
				method := iface.Method(i)
				if !method.Exported() {
					continue
				}

				// m.Handler(http.MethodPost, "/api/v1/SERVICE/ENDPOINT", s.XXXServer)
				apiEndpointName := strcase.ToKebab(method.Name())
				apiServiceName := strcase.ToKebab(svcName)
				g.Id("m").Dot("Handler").Call(
					jen.Line().Qual("net/http", "MethodPost"),
					jen.Line().Lit(fmt.Sprintf("/api/v1/%s/%s", apiServiceName, apiEndpointName)),
					jen.Line().Id("s").Dot(method.Name()+"Server"),
				)
			}

			g.Return(jen.Id("m"))
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
			// return request,nil
			g.Return(jen.Id("request"), jen.Nil())
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
			// return response, nil
			g.Return(jen.Id("response"), jen.Nil())
		}).Line()
}

func generateEmbedSwaggerJSON(f *jen.File) {
	// //go:embed swagger.json
	f.Commentf("//go:embed swagger.json")
	// var swagger embed.FS
	f.Var().Id("swagger").Qual("embed", "FS")
}

func GenerateHTTPTransport(f *jen.File, svc *types.Named) {
	generateHTTPJSONRequestDecoder(f)
	generateHTTPJSONResponseDecoder(f)
	generateEmbedSwaggerJSON(f)
	generateClientSet(f, svc)
	generateServerSet(f, svc)
}
