package typecheck

import (
	"go/types"
	"log"
)

func isStructOrBasic(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Basic:
		return canUseBasic(t)
	case *types.Struct:
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			if field.Exported() && !CanUse(field.Type()) {
				log.Printf("%s is exported and can not use", field.Name())
				return false
			}
		}
		return true
	case *types.Named:
		return isStructOrBasicOrPointerToThem(t.Underlying())
	default:
		return false
	}
}

func isStructOrBasicOrPointerToThem(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Basic:
		return canUseBasic(t)
	case *types.Struct:
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			if field.Exported() && !CanUse(field.Type()) {
				log.Printf("%s is exported and can not use", field.Name())
				return false
			}
		}
		return true
	case *types.Pointer:
		return isStructOrBasic(t.Elem())
	case *types.Named:
		return isStructOrBasicOrPointerToThem(t.Underlying())
	default:
		return false
	}
}

func canUseBasic(typ *types.Basic) bool {
	switch typ.Kind() {
	case types.Bool, types.Float32, types.Float64, types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
		types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
		types.String:
		return true
	default:
		return false
	}
}

func isIntegralOrString(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64, types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.String:
			return true
		default:
			return false
		}
	default:
		return false
	}
}

func CanUse(typ types.Type) bool {
	switch t := typ.(type) {
	case *types.Basic:
		return canUseBasic(t)
	case *types.Array:
		return isStructOrBasicOrPointerToThem(t.Elem())
	case *types.Slice:
		return isStructOrBasicOrPointerToThem(t.Elem())
	case *types.Map:
		return isIntegralOrString(t.Key()) && isStructOrBasicOrPointerToThem(t.Elem())
	case *types.Named:
		return CanUse(t.Underlying())
	case *types.Pointer:
		return isStructOrBasic(t.Elem())
	case *types.Struct:
		return isStructOrBasicOrPointerToThem(t)
	default:
		return false
	}
}
