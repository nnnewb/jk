package utils

import (
	"errors"
	"go/importer"
	"go/token"
	"go/types"
	"log"

	"github.com/dave/jennifer/jen"
)

// IsContext check given type is context.Context.
func IsContext(tp types.Type) bool {
	switch t := tp.(type) {
	case *types.Named:
		// FIXME: somehow types.Identical(ctxType, t) always return false here.
		return t.Obj().Pkg().Name() == "context" && t.Obj().Name() == "Context"
	case *types.Interface:
		pkg, err := importer.Default().Import("context")
		if err != nil {
			panic(err)
		}

		// ctxObj should be types.Named
		ctxObj := pkg.Scope().Lookup("Context")
		if ctxObj == nil {
			panic(errors.New("unable to find context.Context, this should not happen"))
		}

		ctxType := ctxObj.Type().(*types.Named)
		return types.Identical(ctxType.Underlying(), t)
	default:
		return false
	}
}

// IsError check given type is error.
func IsError(tp types.Type) bool {
	return types.Identical(types.Universe.Lookup("error").Type(), tp)
}

// IsPrimitive check given type is primitive type
//
// These type are primitive:
//   bool float32 float64 string
//   int  int8    int16   int32  int64
//   uint uint8   uint16  uint32 uint64
//
// Be aware, unsafe pointer and complex consider as not primitive for now.
func IsPrimitive(tp types.Type) bool {
	switch t := tp.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool, types.Float32, types.Float64, types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.String:
			return true
		default:
			// TODO: uintptr and complex are exclude from primitive
			return false
		}
	default:
		return false
	}
}

// IsSerializableType check given type is serializable.
func IsSerializableType(tp types.Type) bool {
	switch t := tp.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Bool, types.Float32, types.Float64, types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.String:
			return true
		default:
			return false
		}
	case *types.Struct:
		// empty struct not allowed
		if t.NumFields() == 0 {
			return false
		}

		for i := 0; i < t.NumFields(); i++ {
			f := t.Field(i)
			if f.Exported() {
				if !IsSerializableType(f.Type()) {
					return false
				}
			}
		}
		return true
	case *types.Pointer:
		// TODO: only allowed pointer to struct, need reconsider
		return IsStruct(t.Elem())
	case *types.Named:
		return IsSerializableType(t.Underlying())
	case *types.Map:
		// well, strict key type sounds reasonable for me.
		// TODO: available value type exclude map/slice due to gRPC limits. there is workaround but I'm too lazy to fix this shit.
		if (IsIntegral(t.Key()) || IsString(t.Key())) && !IsMap(t.Elem()) && !IsSlice(t.Elem()) {
			return IsSerializableType(t.Key()) && IsSerializableType(t.Elem())
		}
		return false
	case *types.Slice:
		if IsMap(t.Elem()) {
			// TODO: slice of map consider as non-serializable due to gRPC limits. there is workaround but I'm too lazy to fix this shit.
			return false
		}
		return IsSerializableType(t.Elem())
	case *types.Array:
		if IsMap(t.Elem()) {
			// TODO: slice of map consider as non-serializable due to gRPC limits. there is workaround but I'm too lazy to fix this shit.
			return false
		}
		return IsSerializableType(t.Elem())
	case *types.Interface, *types.Chan, *types.Signature, *types.TypeParam, *types.Union, *types.Tuple:
		return false
	default:
		return false
	}
}

// IsIntegral check given type is integral.
//
// these types are consider as integral:
//   int int8 int16 int32 int64
//   uint uint8 uint16 uint32 uint64
//   byte rune bool
func IsIntegral(tp types.Type) bool {
	switch t := tp.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Bool:
			// byte and rune are just alias name for uint8/int32
			return true
		default:
			return false
		}
	default:
		return false
	}
}

// IsBoolean check given type is boolean.
func IsBoolean(tp types.Type) bool {
	if t, ok := tp.(*types.Basic); ok {
		return t.Kind() == types.Bool
	}
	return false
}

// IsString check given type is string.
func IsString(tp types.Type) bool {
	if t, ok := tp.(*types.Basic); ok {
		return t.Kind() == types.String
	}
	return false
}

// IsMap check given type is map
func IsMap(tp types.Type) bool {
	_, ok := tp.(*types.Map)
	return ok
}

// IsSlice check given type is array or slice.
func IsSlice(tp types.Type) bool {
	switch tp.(type) {
	case *types.Array, *types.Slice:
		return true
	default:
		return false
	}
}

// IsStruct check given type is struct.
func IsStruct(tp types.Type) bool {
	_, ok := tp.(*types.Struct)
	return ok
}

// IsInterface check given type is interface.
func IsInterface(tp types.Type) bool {
	_, ok := tp.(*types.Interface)
	return ok
}

// IsFuncSignature check given type is types.Signature.
func IsFuncSignature(tp types.Type) bool {
	_, ok := tp.(*types.Signature)
	return ok
}

// IsNamed check given type is types.Named
func IsNamed(tp types.Type) bool {
	_, ok := tp.(*types.Named)
	return ok
}

// IsPointer check given type is types.Pointer
func IsPointer(tp types.Type) bool {
	_, ok := tp.(*types.Pointer)
	return ok
}

// IsNamed check given type is named struct
func IsNamedStruct(tp types.Type) bool {
	return IsNamed(tp) && IsStruct(tp.Underlying())
}

// IsIndirectNamedStruct check given type is pointer and point to named struct
func IsIndirectNamedStruct(tp types.Type) bool {
	return IsPointer(tp) && IsNamedStruct(tp.(*types.Pointer).Elem())
}

// GetSliceElem check given type is slice/array and return element type.
// If given type is not slice/array, return nil.
func GetSliceElem(tp types.Type) types.Type {
	if sl, ok := tp.(*types.Slice); ok {
		return sl.Elem()
	} else if arr, ok := tp.(*types.Array); ok {
		return arr.Elem()
	} else {
		return nil
	}
}

// ZeroLit make zero value literal for input type
func ZeroLit(tp types.Type) jen.Code {
	switch t := tp.(type) {
	case *types.Array, *types.Slice, *types.Map, *types.Pointer, *types.Chan, *types.Tuple, *types.Union, *types.TypeParam, *types.Interface, *types.Signature:
		return jen.Nil()
	case *types.Named:
		return jen.Qual(t.Obj().Pkg().Path(), t.Obj().Name()).Add(ZeroLit(t.Underlying()))
	case *types.Struct:
		return jen.Block()
	case *types.Basic:
		switch t.Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64:
			return jen.Lit(0)
		case types.String:
			return jen.Lit("")
		case types.Bool:
			return jen.Lit(false)
		default:
			// anything else?
			return jen.Nil()
		}
	default:
		return jen.Nil()
	}
}

// TypeQual if given type is types.Named and not builtin type, return a jen.Qual for it.
func TypeQual(tp types.Type) jen.Code {
	switch t := tp.(type) {
	case *types.Named:
		if t.Obj().Pkg() != nil {
			return jen.Qual(t.Obj().Pkg().Path(), t.Obj().Name())
		}
		return jen.Id(t.Obj().Name())
	default:
		return jen.Id(tp.String())
	}
}

// CheckMethodSignature check method signature match our requirements.
//
// good method signature example:
//
//   ToUpper(ctx context.Context, text string) (text string, err error)
//
// requirements listed below:
//
// 0. must have 1 parameter at least and 1 result at least
//
// 1. context.Context must be first parameter
//
// 2. all parameter should be primitive type (include slice/map but element type should also be primitive)
//
// 3. last result must be error
//
// 4. all results should be primitive type (include slice/map but element type should also be primitive)
//
// 5. all results should be named (for readable response field name)
func CheckMethodSignature(method *types.Func) bool {
	methodSignature := method.Type().(*types.Signature)

	if methodSignature.TypeParams().Len() > 0 {
		log.Printf("ignore generic method %s, not supported", method.Name())
		return false
	}

	if methodSignature.Params().Len() == 0 || methodSignature.Results().Len() == 0 {
		log.Printf("ignore bad method %s, function must have 1 parameter(ctx context.Context) and 1 return value(error) at least", method.Name())
		return false
	}

	if !IsContext(methodSignature.Params().At(0).Type()) {
		log.Printf("ignore bad method %s, first parameter must be (ctx context.Context)", method.Name())
		return false
	}

	for i := 1; i < methodSignature.Params().Len(); i++ {
		param := methodSignature.Params().At(i)
		if !IsSerializableType(param.Type()) {
			log.Printf("ignore bad method %s, parameter type %s not support", method.Name(), param.Type())
			return false
		}
	}

	if !IsError(methodSignature.Results().At(methodSignature.Results().Len() - 1).Type()) {
		log.Printf("ignore bad method %s, last return type should be error", method.Name())
		return false
	}

	for i := 0; i < methodSignature.Results().Len()-1; i++ {
		result := methodSignature.Results().At(i)
		if !IsSerializableType(result.Type()) {
			log.Printf("ignore bad method %s, non-primitive result (%s(%T)) not support yet", method.Name(), result.Type(), result.Type())
			return false
		}
		if result.Name() == "" {
			log.Printf("ignore bad method %s, all results must be named", method.Name())
			return false
		}
	}

	return true
}

func FilterCorrespondPublicMethod(typ *types.Named) []*types.Func {
	svc := typ.Underlying().(*types.Interface)
	ret := make([]*types.Func, 0, svc.NumMethods())
	for i := 0; i < svc.NumMethods(); i++ {
		method := svc.Method(i)

		if !method.Exported() {
			continue
		}

		// preflight check, see comments of utils.CheckMethodSignature for more detail
		if !CheckMethodSignature(method) {
			continue
		}

		ret = append(ret, method)
	}

	return ret
}

func NewParamsStruct(f *types.Func, ctxType, errType types.Type) *types.Struct {
	s := f.Type().(*types.Signature)
	fields := make([]*types.Var, 0, s.Params().Len())
	for i := 0; i < s.Params().Len(); i++ {
		if types.Identical(s.Params().At(i).Type(), ctxType) || types.Identical(s.Params().At(i).Type(), errType) {
			continue
		}
		fields = append(fields, types.NewVar(token.NoPos, s.Params().At(i).Pkg(), s.Params().At(i).Name(), s.Params().At(i).Type()))
	}
	return types.NewStruct(fields, nil)
}

func NewResultsStruct(f *types.Func, ctxType, errType types.Type) *types.Struct {
	s := f.Type().(*types.Signature)
	fields := make([]*types.Var, 0, s.Results().Len())
	for i := 0; i < s.Results().Len(); i++ {
		if types.Identical(s.Results().At(i).Type(), ctxType) || types.Identical(s.Results().At(i).Type(), errType) {
			continue
		}
		fields = append(fields, types.NewVar(token.NoPos, s.Results().At(i).Pkg(), s.Results().At(i).Name(), s.Results().At(i).Type()))
	}
	return types.NewStruct(fields, nil)
}

func GenStruct(name string, pkg *types.Package, typ *types.Struct) jen.Code {
	return jen.Type().Id(name).StructFunc(func(g *jen.Group) {
		for i := 0; i < typ.NumFields(); i++ {
			field := typ.Field(i)
			if !field.Exported() {
				continue
			}

			g.Id(field.Name()).Add(TypeQual(pkg, field.Type()))
		}
	})
}
