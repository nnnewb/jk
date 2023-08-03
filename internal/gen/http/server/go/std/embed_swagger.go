package stdsvr

import (
	"fmt"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/nnnewb/jk/internal/domain"
)

func GenerateEmbedSwaggerJSON(f *jen.File, service *domain.Service) {
	// //go:embed swagger.json
	f.Commentf("//go:embed swagger.json")
	// var swagger embed.FS
	f.Var().Id("httpEmbedSwagger").Qual("embed", "FS")

	f.Func().
		Params(jen.Id("s").Op("*").Id("HTTPServerSet")).
		Id("RegisterEmbedSwaggerUI").
		Params(jen.Id("mux").Op("*").Qual("net/http", "ServeMux")).
		BlockFunc(func(g *jen.Group) {
			// m.Handle("/swagger/service/spec/", http.FileServer(http.FS(swagger)))
			g.Id("mux").Dot("Handle").Call(
				jen.Line().Lit(fmt.Sprintf("/swagger/%s/spec/*rest", strcase.ToKebab(service.Name()))),
				jen.Line().Qual("net/http", "StripPrefix").Call(
					jen.Lit(fmt.Sprintf("/swagger/%s/spec/", strcase.ToKebab(service.Name()))),
					jen.Qual("net/http", "FileServer").
						Call(jen.Qual("net/http", "FS").Call(jen.Id("httpEmbedSwagger")))),
			)

			// m.Handle("/swagger/service/swagger-ui/*", httpSwagger.Handler(httpSwagger.URL("/swagger/service/swagger.json")))
			g.Id("mux").Dot("Handle").Call(
				jen.Line().Lit(fmt.Sprintf("/swagger/%s/swagger-ui/*rest", strcase.ToKebab(service.Name()))),
				jen.Line().Qual("github.com/swaggo/http-swagger/v2", "Handler").
					Call(jen.Qual("github.com/swaggo/http-swagger/v2", "URL").
						Call(jen.Lit(fmt.Sprintf("/swagger/%s/spec/swagger.json", strcase.ToKebab(service.Name()))))),
			)
		})
}
