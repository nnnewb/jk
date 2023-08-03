package stdsvr

import (
	"go/types"
	"log"
	"net/http"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/nnnewb/battery/slices"
	"github.com/nnnewb/jk/internal/domain"
	"github.com/nnnewb/jk/internal/gen/http/common"
)

func generateServerSet(f *jen.File, service *domain.Service) {
	common.HTTPPopulateDefaultAnnotations(service)
	interfaceType := service.Interface.Underlying().(*types.Interface)

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
}

func generateRegister(f *jen.File, service *domain.Service) {
	common.HTTPPopulateDefaultAnnotations(service)
	// func (s HTTPServerSet) Register(mux *http.ServeMux) {
	f.Func().
		Params(jen.Id("s").Op("*").Id("HTTPServerSet")).
		Id("Register").
		Params(jen.Id("mux").Op("*").Qual("net/http", "ServeMux")).
		BlockFunc(func(g *jen.Group) {
			groupByPath := slices.GroupBy(service.Methods, func(m *domain.Method) string { return m.Annotations.HTTPPath })
			for path, methods := range groupByPath {
				// mux.HandleFunc("/path/to/endpoint", func(wr ResponseWriter, req *Request) {
				//   switch req.Method {
				//     case http.MethodPost:
				//       s.XXServer.ServeHTTP(wr, req)
				//       return
				//     default:
				//       wr.WriteHeader(http.StatusMethodNotAllowed)
				//   }
				// })
				g.Id("mux").Dot("HandleFunc").Call(
					jen.Line().Lit(path),
					jen.Func().Params(
						jen.Id("wr").Qual("net/http", "ResponseWriter"),
						jen.Id("req").Op("*").Qual("net/http", "Request"),
					).BlockFunc(func(g *jen.Group) {
						g.Switch(jen.Id("req").Dot("Method")).BlockFunc(func(g *jen.Group) {
							for _, method := range methods {
								g.Case(method.HTTPMethodJen())
								g.Id("s").
									Dot(method.Func.Name()+"Server").
									Dot("ServeHTTP").Call(jen.Id("wr"), jen.Id("req"))
								g.Return()
							}
							g.Default()
							g.Id("wr").Dot("WriteHeader").Call(
								jen.Qual("net/http", "StatusMethodNotAllowed"),
							)
							g.List(jen.Id("_"), jen.Err()).Op(":=").Qual("io", "WriteString").Call(
								jen.Id("wr"),
								jen.Lit(`"{\"code\": -1, \"message\": \"method not allowed\"}"`),
							)
							g.If(jen.Err().Op("!=").Nil()).Block(
								g.Qual("log", "Printf").
									Call(jen.Lit("write 405 error response failed, error %+v"), jen.Err()),
							)
							g.Return()
						})
					}),
				)
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
			// return request,nil
			g.Return(jen.Op("&").Id("request"), jen.Nil())
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
	// func beautifyErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
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

func GenerateHTTPTransportServer(f *jen.File, svc *domain.Service) {
	generateBeautifyErrorEncoder(f)
	generateHTTPJSONRequestDecoder(f)
	generateHTTPQueryStringRequestDecoder(f)
	generateServerSet(f, svc)
	generateRegister(f, svc)
}
