package gen

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/iancoleman/strcase"
	"github.com/nnnewb/jk/internal/domain"
)

func httpPopulateDefaultAnnotations(service *domain.Service) {
	if service.Annotations.SwaggerInfoAPIVersion == "" {
		service.Annotations.SwaggerInfoAPIVersion = "v0.1.0"
	}

	if service.Annotations.SwaggerInfoAPITitle == "" {
		service.Annotations.SwaggerInfoAPITitle = service.Interface.Obj().Name()
	}

	if service.Annotations.HTTPBasePath == "" {
		service.Annotations.HTTPBasePath = fmt.Sprintf("/api/v1/%s/", strcase.ToKebab(service.Interface.Obj().Name()))
	}

	for _, method := range service.Methods {
		if !method.Func.Exported() {
			continue
		}

		method.Annotations.HTTPMethod = strings.ToUpper(method.Annotations.HTTPMethod)
		switch method.Annotations.HTTPMethod {
		case http.MethodGet, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodPost:
		default:
			method.Annotations.HTTPMethod = http.MethodPost
		}

		if method.Annotations.HTTPPath == "" {
			method.Annotations.HTTPPath = path.Join(service.Annotations.HTTPBasePath, strcase.ToKebab(method.Func.Name()))
		}
	}
}

func generateEmbedSwaggerJSON(f *jen.File) {
	// //go:embed swagger.json
	f.Commentf("//go:embed swagger.json")
	// var swagger embed.FS
	f.Var().Id("swagger").Qual("embed", "FS")
}

func GenerateHTTPTransportClient(f *jen.File, service *domain.Service) {
	generateHTTPJSONResponseDecoder(f)
	generateHTTPQueryStringEncoder(f)
	generateClientSet(f, service)
}

func GenerateHTTPTransportServer(f *jen.File, svc *domain.Service) {
	generateBeautifyErrorEncoder(f)
	generateHTTPJSONRequestDecoder(f)
	generateHTTPQueryStringRequestDecoder(f)
	generateEmbedSwaggerJSON(f)
	generateServerSet(f, svc)
}
