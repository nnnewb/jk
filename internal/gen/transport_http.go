package gen

import (
	"fmt"
	"go/types"
	"log"
	"net/http"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/nnnewb/jk/internal/domain"
)

func generateClientSet(f *jen.File, service *domain.Service, apiVer string) {
	httpPopulateDefaultAnnotations(service)
	interfaceType := service.Interface.Underlying().(*types.Interface)

	f.Type().Id("HTTPClientSet").StructFunc(func(g *jen.Group) {
		for i := 0; i < interfaceType.NumMethods(); i++ {
			method := interfaceType.Method(i)
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
				for _, methodData := range service.Methods {
					method := methodData.Func
					if !method.Exported() {
						continue
					}

					signature := method.Type().(*types.Signature)
					respPtrType := signature.Results().At(0).Type().(*types.Pointer)
					respType := respPtrType.Elem().(*types.Named)

					// check http-method annotation
					httpRequestEncoder := jen.Qual("github.com/go-kit/kit/transport/http", "EncodeJSONRequest")
					httpMethod := jen.Qual("net/http", "MethodPost")
					switch methodData.Annotations.HTTPMethod {
					case http.MethodGet:
						httpMethod = jen.Qual("net/http", "MethodGet")
						httpRequestEncoder = jen.Id("httpQueryStringEncoder")
					case http.MethodDelete:
						httpMethod = jen.Qual("net/http", "MethodDelete")
						httpRequestEncoder = jen.Id("httpQueryStringEncoder")
					case http.MethodPost:
						httpMethod = jen.Qual("net/http", "MethodPost")
					case http.MethodPatch:
						httpMethod = jen.Qual("net/http", "MethodPatch")
					case http.MethodPut:
						httpMethod = jen.Qual("net/http", "MethodPut")
					default:
						log.Printf("unexpected http-method annotation %s for method %s, fallback to POST", strings.ToUpper(methodData.Annotations.HTTPMethod), method.Name())
					}

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
							jen.Line().Add(httpMethod),
							jen.Line().Op("&").Qual("net/url", "URL").Values(jen.DictFunc(func(d jen.Dict) {
								d[jen.Id("Scheme")] = jen.Id("scheme")
								d[jen.Id("Host")] = jen.Qual("fmt", "Sprintf").Call(
									jen.Lit("%s:%d"),
									jen.Id("host"),
									jen.Id("port"))
								d[jen.Id("Path")] = jen.Lit(methodData.Annotations.HTTPPath)
							})),
							jen.Line().Add(httpRequestEncoder),
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
				for i := 0; i < interfaceType.NumMethods(); i++ {
					// XXXEndpoint: s.XXXClient.Endpoint(),
					method := interfaceType.Method(i)
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

func generateServerSet(f *jen.File, service *domain.Service, apiVer string) {
	httpPopulateDefaultAnnotations(service)
	interfaceType := service.Interface.Underlying().(*types.Interface)
	svcName := service.Interface.Obj().Name()

	// type HTTPServerSet {
	f.Type().Id("HTTPServerSet").StructFunc(func(g *jen.Group) {
		for i := 0; i < interfaceType.NumMethods(); i++ {
			method := interfaceType.Method(i)
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
			// options = append(options, khttp.ServerBefore(khttp.PopulateRequestContext))
			g.Id("options").Op("=").Append(
				jen.Id("options"),
				jen.Qual("github.com/go-kit/kit/transport/http", "ServerBefore").
					Call(jen.Qual("github.com/go-kit/kit/transport/http", "PopulateRequestContext")),
			)

			// options = append(options, khttp.ServerErrorEncoder(beautifyErrorEncoder))
			g.Id("options").Op("=").Append(
				jen.Id("options"),
				jen.Qual("github.com/go-kit/kit/transport/http", "ServerErrorEncoder").Call(jen.Id("beautifyErrorEncoder")),
			)

			// return HTTPServerSet{
			g.Return(jen.Id("HTTPServerSet").Values(jen.DictFunc(func(d jen.Dict) {
				for _, methodData := range service.Methods {
					method := methodData.Func
					if !method.Exported() {
						continue
					}

					signature := method.Type().(*types.Signature)
					reqPtrType := signature.Params().At(1).Type().(*types.Pointer)
					reqType := reqPtrType.Elem().(*types.Named)

					// check http-method annotation
					httpRequestDecoder := jen.Id("httpJSONRequestDecoder")
					switch strings.ToUpper(methodData.Annotations.HTTPMethod) {
					case http.MethodGet, http.MethodDelete:
						httpRequestDecoder = jen.Id("httpQueryStringRequestDecoder")
					case http.MethodPost, http.MethodPatch, http.MethodPut:
					default:
						log.Printf("unexpected http-method annotation %s for method %s, fallback to POST", strings.ToUpper(methodData.Annotations.HTTPMethod), method.Name())
					}

					// XXXServer: khttp.NewServer(
					//   endpointSet.XXXEndpoint,
					//   httpRequestDecoder[REQ],
					//   httpJSONResponseEncoder[RESP],
					//   append(options, khttp.ServerBefore(khttp.PopulateRequestContext))...,
					// )
					d[jen.Id(method.Name()+"Server")] = jen.Qual("github.com/go-kit/kit/transport/http", "NewServer").
						Call(
							jen.Line().Id("endpointSet").Dot(method.Name()+"Endpoint"),
							jen.Line().Add(httpRequestDecoder).Types(jen.Qual(reqType.Obj().Pkg().Path(), reqType.Obj().Name())),
							jen.Line().Qual("github.com/go-kit/kit/transport/http", "EncodeJSONResponse"),
							jen.Line().Id("options").Op("..."))
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

			for _, methodData := range service.Methods {
				method := methodData.Func
				if !method.Exported() {
					continue
				}

				// check http-method annotation
				httpMethod := jen.Qual("net/http", "MethodPost")
				switch methodData.Annotations.HTTPMethod {
				case http.MethodGet:
					httpMethod = jen.Qual("net/http", "MethodGet")
				case http.MethodDelete:
					httpMethod = jen.Qual("net/http", "MethodDelete")
				case http.MethodPost:
					httpMethod = jen.Qual("net/http", "MethodPost")
				case http.MethodPatch:
					httpMethod = jen.Qual("net/http", "MethodPatch")
				case http.MethodPut:
					httpMethod = jen.Qual("net/http", "MethodPut")
				default:
					log.Printf("unexpected http-method annotation %s for method %s, fallback to POST", strings.ToUpper(methodData.Annotations.HTTPMethod), method.Name())
				}

				// m.Handler(http.MethodPost, "/path/to/endpoint", s.XXXServer)
				g.Id("m").Dot("Handler").Call(
					jen.Line().Add(httpMethod),
					jen.Line().Lit(methodData.Annotations.HTTPPath),
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
			// return response, nil
			g.Return(jen.Op("&").Id("response"), jen.Nil())
		}).Line()
}

func generateHTTPQueryStringEncoder(f *jen.File) {
	// func httpQueryStringRequestEncoder[T any](c context.Context, r *http.Request, request T) error {
	f.Func().
		Id("httpQueryStringEncoder").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("r").Op("*").Qual("net/http", "Request"),
			jen.Id("request").Any(),
		).
		Params(jen.Error()).
		BlockFunc(func(g *jen.Group) {
			// var values url.Values
			g.Var().Id("values").Qual("net/url", "Values")
			// err := schema.NewEncoder().Encode(request, values)
			g.Err().Op(":=").
				Qual("github.com/gorilla/schema", "NewEncoder").Call().
				Dot("Encode").Call(jen.Id("request"), jen.Id("values"))
			// if err != nil {
			//     return err
			// }
			g.If(jen.Err().Op("!=").Nil()).
				Block(jen.Return(jen.Err()))

			// r.URL.RawQuery = values.String()
			g.Id("r").Dot("URL").Dot("RawQuery").Op("=").
				Id("values").Dot("Encode").Call()

			// return nil
			g.Return(jen.Nil())
		}).Line()
}

func generateHTTPQueryStringRequestDecoder(f *jen.File) {
	// func httpQueryStringRequestDecoder[T any](ctx context.Context, req *http.Request) (any, error) {
	f.Func().
		Id("httpQueryStringRequestDecoder").
		Types(jen.Id("T").Any()).
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Id("req").Op("*").Qual("net/http", "Request"),
		).
		Params(
			jen.Any(),
			jen.Error()).
		BlockFunc(func(g *jen.Group) {
			// var request T
			g.Var().Id("request").Id("T")
			// defer req.Body.Close()
			g.Defer().Id("req").Dot("Body").Dot("Close").Call()
			// err := schema.NewDecoder().Decode(&request, req.URL.RawQuery)
			g.Err().Op(":=").
				Qual("github.com/gorilla/schema", "NewDecoder").Call().
				Dot("Decode").Call(jen.Op("&").Id("request"), jen.Id("req").Dot("URL").Dot("Query").Call())

			// if err != nil { return nil, err }
			g.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err()))
			// return &request, nil
			g.Return(jen.Op("&").Id("request"), jen.Nil())
		}).Line()
}

func generateBeautifyErrorEncoder(f *jen.File) {
	// func(ctx context.Context, err error, w http.ResponseWriter)
	f.Func().
		Id("beautifyErrorEncoder").
		Params(
			jen.Id("ctx").Qual("context", "Context"),
			jen.Err().Error(),
			jen.Id("wr").Qual("net/http", "ResponseWriter"),
		).
		BlockFunc(func(g *jen.Group) {
			g.Var().Id("resp").Struct(
				jen.Id("Code").Int().Tag(map[string]string{"json": "code"}),
				jen.Id("Message").String().Tag(map[string]string{"json": "message"}),
			).Line()

			g.Id("resp").Dot("Code").Op("=").Lit(-1)
			g.Id("resp").Dot("Message").Op("=").
				Qual("fmt", "Sprintf").Call(jen.Lit("error occurred: %v"), jen.Err())

			g.Qual("encoding/json", "NewEncoder").Call(jen.Id("wr")).
				Dot("Encode").Call(jen.Id("resp"))
		}).Line()
}

func generateEmbedSwaggerJSON(f *jen.File) {
	// //go:embed swagger.json
	f.Commentf("//go:embed swagger.json")
	// var swagger embed.FS
	f.Var().Id("swagger").Qual("embed", "FS")
}

func GenerateHTTPTransport(f *jen.File, svc *domain.Service, ver string) {
	generateBeautifyErrorEncoder(f)
	generateHTTPJSONRequestDecoder(f)
	generateHTTPJSONResponseDecoder(f)
	generateHTTPQueryStringEncoder(f)
	generateHTTPQueryStringRequestDecoder(f)
	generateEmbedSwaggerJSON(f)
	generateClientSet(f, svc, ver)
	generateServerSet(f, svc, ver)
}
