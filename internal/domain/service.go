package domain

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"emperror.dev/errors"
	"github.com/nnnewb/jk/internal/utils"
)

type ServiceAnnotations struct {
	SwaggerInfoAPIVersion string `jk:"swagger-info-api-version"`
	SwaggerInfoAPITitle   string `jk:"swagger-info-api-title"`
	HTTPBasePath          string `jk:"http-base-path"`
}

type Service struct {
	Interface   *types.Named
	GenDecl     *ast.GenDecl        // 如果是 type ( /* document here */ xxx interface )
	TypeSpec    *ast.TypeSpec       // 如果是 /* document here */ type xxx interface
	Annotations *ServiceAnnotations // 以@开头写在注释里的注解

	Methods []*Method // 预先解析好的 method 列表
}

func (s *Service) Name() string {
	return s.Interface.Obj().Name()
}

func ParseInterfaceData(pkg *types.Package, astPkg *ast.Package, name string) (*Service, error) {
	ret := &Service{}

	ast.Inspect(astPkg, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.GenDecl:
			if n.Tok == token.TYPE {
				for _, spec := range n.Specs {
					typeSpec := spec.(*ast.TypeSpec)
					if typeSpec.Name != nil && typeSpec.Name.Name == name {
						ret.GenDecl = n
						ret.TypeSpec = typeSpec
						break
					}
				}
			}
			return false
		default:
			return true
		}
	})

	if ret.TypeSpec == nil {
		return nil, errors.Errorf("Type %s not found", name)
	}

	if _, ok := ret.TypeSpec.Type.(*ast.InterfaceType); !ok {
		return nil, fmt.Errorf("%s is not an interface type, got %T", name, ret.TypeSpec.Type)
	}

	obj := pkg.Scope().Lookup(name)
	if obj == nil {
		return nil, fmt.Errorf("no %s found in package %s", name, pkg.Path())
	}

	named := obj.Type().(*types.Named)
	if _, ok := named.Underlying().(*types.Interface); !ok {
		return nil, fmt.Errorf("%s is not an interface type, got %T", name, named.Underlying())
	}

	ret.Interface = obj.Type().(*types.Named)

	var cg *ast.CommentGroup
	if len(ret.GenDecl.Specs) > 1 {
		cg = ret.TypeSpec.Doc
	} else {
		cg = ret.GenDecl.Doc
	}

	ret.Annotations = &ServiceAnnotations{}
	err := utils.UnmarshalAnnotations(cg, ret.Annotations)
	if err != nil {
		return nil, err
	}

	// 预先解析所有方法
	interfaceType := ret.Interface.Underlying().(*types.Interface)
	ret.Methods = make([]*Method, 0, ret.Interface.NumMethods())
	for i := 0; i < interfaceType.NumMethods(); i++ {
		m := interfaceType.Method(i)
		astInterfaceType := ret.TypeSpec.Type.(*ast.InterfaceType)
		for _, field := range astInterfaceType.Methods.List {
			if field.Names[0].Name == m.Name() {
				method := &Method{
					Func:        m,
					Field:       field,
					Annotations: &MethodAnnotations{},
				}
				err := utils.UnmarshalAnnotations(field.Doc, method.Annotations)
				if err != nil {
					return nil, err
				}
				ret.Methods = append(ret.Methods, method)
			}
		}
	}

	return ret, nil
}
