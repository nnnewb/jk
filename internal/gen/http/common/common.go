package common

import (
	"fmt"
	"net/http"
	"path"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/nnnewb/jk/internal/domain"
)

func HTTPPopulateDefaultAnnotations(service *domain.Service) {
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

func GetJsonName(tag string) (string, bool) {
	jsonTag := reflect.StructTag(tag).Get("json")
	if jsonTag == "-" {
		return "", false
	}

	var jsonName string
	for _, v := range strings.Split(jsonTag, ",") {
		if v != "omitempty" && strings.TrimSpace(v) != "" {
			jsonName = v
			break
		}
	}

	if jsonName == "" {
		return "", false
	}

	return jsonName, true
}
