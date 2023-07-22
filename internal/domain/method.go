package domain

import (
	"go/ast"
	"go/types"
)

type MethodAnnotations struct {
	HTTPMethod string `jk:"http-method"`
	HTTPPath   string `jk:"http-path"`
}

type Method struct {
	parent      *Service
	Func        *types.Func
	Field       *ast.Field
	Annotations *MethodAnnotations
}
