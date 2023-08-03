package std

import (
	"go/types"
	"log"
	"net/http"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/nnnewb/jk/internal/domain"
	"github.com/nnnewb/jk/internal/gen/http/common"
)

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

func generateClientSet(f *jen.File, service *domain.Service) {
	common.HTTPPopulateDefaultAnnotations(service)
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

func GenerateHTTPTransportClient(f *jen.File, service *domain.Service) {
	generateHTTPJSONResponseDecoder(f)
	generateHTTPQueryStringEncoder(f)
	generateClientSet(f, service)
}
