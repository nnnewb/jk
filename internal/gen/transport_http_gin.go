package gen

import (
	"fmt"
	"net/http"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/nnnewb/jk/internal/domain"
)

func generateCommonCode(f *jen.File) {
	// type GenericEndpoint[Request, Response any] func(ctx context.Context, req *Request) (*Response, error)
	// f.Type().Id("GenericEndpoint").Types(jen.Id("Request").Any(), jen.Id("Response").Any()).
	// 	Func().
	// 	Params(
	// 		jen.Id("ctx").Qual("context", "Context"),
	// 		jen.Id("req").Op("*").Id("Request"),
	// 	).
	// 	Params(
	// 		jen.Op("*").Id("Response"),
	// 		jen.Error(),
	// 	).
	// 	Line()

	// type RequestDecoder func(c *gin.Context, req any) error
	f.Type().Id("RequestDecoder").
		Func().
		Params(
			jen.Id("c").Op("*").Qual("github.com/gin-gonic/gin", "Context"),
			jen.Id("req").Any()).
		Error().
		Line()

	// func QueryStringDecoder(c *gin.Context, req any) error {
	// 	return c.BindQuery(req)
	// }
	f.Func().Id("QueryStringDecoder").
		Params(
			jen.Id("c").Op("*").Qual("github.com/gin-gonic/gin", "Context"),
			jen.Id("req").Any(),
		).
		Error().
		Block(
			jen.Return(
				jen.Id("c").Dot("BindQuery").Call(jen.Id("req")),
			),
		).
		Line()

	// func JSONBodyDecoder(c *gin.Context, req any) error {
	// 	return c.BindJSON(req)
	// }
	f.Func().Id("JSONBodyDecoder").
		Params(
			jen.Id("c").Op("*").Qual("github.com/gin-gonic/gin", "Context"),
			jen.Id("req").Any(),
		).
		Error().
		Block(
			jen.Return(
				jen.Id("c").Dot("BindJSON").Call(jen.Id("req")),
			),
		).
		Line()

	// type ResponseEncoder func(c *gin.Context, resp any)
	f.Type().Id("ResponseEncoder").
		Func().
		Params(
			jen.Id("c").Op("*").Qual("github.com/gin-gonic/gin", "Context"),
			jen.Id("resp").Any(),
		).
		Line()

	// func JSONBodyEncoder(c *gin.Context, resp any) {
	// 	c.JSON(200, resp)
	// }
	f.Func().Id("JSONBodyEncoder").
		Params(
			jen.Id("c").Op("*").Qual("github.com/gin-gonic/gin", "Context"),
			jen.Id("resp").Any(),
		).
		Block(
			jen.Id("c").Dot("JSON").Call(jen.Lit(200), jen.Id("resp")),
		).
		Line()

	// func Handler[Request, Response any](ep GenericEndpoint[Request, Response], decoder RequestDecoder, encoder ResponseEncoder) gin.HandlerFunc {
	// 	return func(c *gin.Context) {
	// 		var req = new(Request)
	// 		err := decoder(c, req)
	// 		if err != nil {
	// 			c.AbortWithStatusJSON(400, gin.H{
	// 				"code":    -1,
	// 				"message": fmt.Sprintf("unable to parse request payload, error %v", err),
	// 			})
	// 			return
	// 		}
	//
	// 		resp, err := ep(c.Request.Context(), req)
	// 		encoder(c, resp)
	// 	}
	// }
	f.Func().Id("Handler").
		Types(jen.Id("Request").Any()).
		Params(
			jen.Id("ep").Qual("github.com/go-kit/kit/endpoint", "Endpoint"),
			jen.Id("decoder").Id("RequestDecoder"),
			jen.Id("encoder").Id("ResponseEncoder"),
		).
		Qual("github.com/gin-gonic/gin", "HandlerFunc").
		BlockFunc(func(g *jen.Group) {
			g.Return(
				jen.Func().Params(jen.Id("c").Op("*").Qual("github.com/gin-gonic/gin", "Context")).
					BlockFunc(func(g *jen.Group) {
						g.Var().Id("req").Op("=").New(jen.Id("Request"))
						g.Err().Op(":=").Id("decoder").Call(jen.Id("c"), jen.Id("req"))
						g.If(jen.Err().Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
							g.Id("c").Dot("AbortWithStatusJSON").
								Call(
									jen.Lit(400),
									jen.Qual("github.com/gin-gonic/gin", "H").
										Values(jen.Dict{
											jen.Lit("code"):    jen.Lit(-1),
											jen.Lit("message"): jen.Qual("fmt", "Sprintf").Call(jen.Lit("unable to parse request payload, error %v"), jen.Err()),
										}),
								)
							g.Return()
						}).Line()

						g.List(jen.Id("resp"), jen.Err()).Op(":=").Id("ep").
							Call(
								jen.Id("c").Dot("Request").Dot("Context").Call(),
								jen.Id("req"),
							)
						g.Id("encoder").Call(jen.Id("c"), jen.Id("resp"))
					}))
		})
}

func generateGinServerSet(f *jen.File, service *domain.Service) {
	f.Type().Id("GinServerSet").StructFunc(func(g *jen.Group) {
		for _, method := range service.Methods {
			g.Id(method.Func.Name()+"Handler").Qual("github.com/gin-gonic/gin", "HandlerFunc")
		}
	}).Line()

	f.Func().Id("NewGinServerSet").
		Params(jen.Id("eps").Id("EndpointSet")).
		Op("*").Id("GinServerSet").
		BlockFunc(func(g *jen.Group) {
			g.Return(jen.Op("&").Id("GinServerSet").
				Values(jen.DictFunc(func(d jen.Dict) {
					for _, method := range service.Methods {
						switch method.Annotations.HTTPMethod {
						case http.MethodGet, http.MethodDelete:
							d[jen.Id(method.Func.Name()+"Handler")] = jen.
								Id("Handler").
								Types(method.RequestTypeCodeJen()).
								Call(
									jen.Id("eps").Dot(method.Func.Name()+"Endpoint"),
									jen.Id("QueryStringDecoder"),
									jen.Id("JSONBodyEncoder"),
								)
						case http.MethodPost, http.MethodPut, http.MethodPatch:
							d[jen.Id(method.Func.Name()+"Handler")] = jen.
								Id("Handler").
								Types(method.RequestTypeCodeJen()).
								Call(
									jen.Id("eps").Dot(method.Func.Name()+"Endpoint"),
									jen.Id("JSONBodyDecoder"),
									jen.Id("JSONBodyEncoder"),
								)
						}
					}
				})))
		}).Line()

	f.Func().
		Params(jen.Id("s").Op("*").Id("GinServerSet")).
		Id("Register").
		Params(jen.Id("router").Qual("github.com/gin-gonic/gin", "IRouter")).
		BlockFunc(func(g *jen.Group) {
			for _, method := range service.Methods {
				g.Id("router").Dot(method.Annotations.HTTPMethod).
					Call(
						jen.Lit(method.Annotations.HTTPPath),
						jen.Id("s").Dot(method.Func.Name()+"Handler"),
					)
			}
		}).Line()
}

func GenerateGinEmbedSwaggerUI(f *jen.File, service *domain.Service) {
	// //go:embed swagger.json
	f.Commentf("//go:embed swagger.json")
	// var swagger embed.FS
	f.Var().Id("swagger").Qual("embed", "FS")

	f.Commentf("// RegisterEmbedSwaggerUI register embed swagger-ui urls")
	// func RegisterEmbedSwaggerUI(r gin.IRouter) {
	f.Func().Id("RegisterEmbedSwaggerUI").
		Params(jen.Id("r").Qual("github.com/gin-gonic/gin", "IRouter")).
		BlockFunc(func(g *jen.Group) {
			// fs = http.FS(swagger)
			g.Id("fs").Op(":=").Qual("net/http", "FS").Call(jen.Id("swagger"))
			// handler = http.FileServer(fs)
			g.Id("handler").Op(":=").Qual("net/http", "FileServer").Call(jen.Id("fs"))
			// handler = http.StripPrefix("", handler)
			prefix := fmt.Sprintf("/swagger/%s/spec/", strcase.ToKebab(service.Name()))
			g.Id("handler").Op("=").Qual("net/http", "StripPrefix").Call(jen.Lit(prefix), jen.Id("handler"))
			// r.GET("/swagger/SVC/spec/*rest", http.WrapH(handler))
			url := fmt.Sprintf("/swagger/%s/spec/*rest", strcase.ToKebab(service.Name()))
			g.Id("r").Dot("GET").Call(
				jen.Line().Lit(url),
				jen.Line().Qual("github.com/gin-gonic/gin", "WrapH").Call(jen.Id("handler")),
			)

			// u := httpSwagger.URL("/swagger/SVC/swagger.json")
			url = fmt.Sprintf("/swagger/%s/spec/swagger.json", strcase.ToKebab(service.Name()))
			g.Id("u").Op(":=").Qual("github.com/swaggo/http-swagger/v2", "URL").Call(jen.Lit(url))
			// handler = httpSwagger.Handler(u)
			g.Id("handler").Op("=").Qual("github.com/swaggo/http-swagger/v2", "Handler").Call(jen.Id("u"))
			// r.GET("/swagger/SVC/swagger-ui/*", gin.WrapH(handler))
			url = fmt.Sprintf("/swagger/%s/swagger-ui/*rest", strcase.ToKebab(service.Name()))
			g.Id("r").Dot("GET").Call(
				jen.Lit(url),
				jen.Qual("github.com/gin-gonic/gin", "WrapH").Call(jen.Id("handler")),
			)
		}).Line()
}

func GenerateGin(f *jen.File, service *domain.Service) error {
	httpPopulateDefaultAnnotations(service)
	generateCommonCode(f)
	generateGinServerSet(f, service)
	return nil
}
