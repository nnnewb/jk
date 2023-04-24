package gen

import "go/types"

// IsSerializable check is given type serializable
func IsSerializable(t types.Type) bool {
	switch t.(type) {
	case *types.Basic:
		return isBasicSerializableType(t)
	case *types.Pointer:
		return isSerializablePointerType(t)
	case *types.Slice:
		return isSerializableSliceType(t)
	case *types.Map:
		return isSerializableMapType(t)
	case *types.Struct:
		return isSerializableStructureType(t)
	default:
		return false
	}
}

func isSerializableStructureType(t types.Type) bool {
	if s, ok := t.(*types.Struct); ok {
		for i := 0; i < s.NumFields(); i++ {
			field := s.Field(i)
			switch ft := field.Type().(type) {
			case *types.Named:
				return IsSerializable(ft.Underlying())
			case *types.Basic:
				if !isBasicSerializableType(ft) {
					return false
				}
			case *types.Struct:
				if !isSerializableStructureType(ft) {
					return false
				}
			case *types.Pointer:
				if !isSerializablePointerType(ft) {
					return false
				}
			case *types.Map:
				if !isSerializableMapType(ft) {
					return false
				}
			case *types.Slice:
				if !isSerializableSliceType(ft) {
					return false
				}
			default:
				return false
			}
		}
		return true
	}
	return false
}

func isSerializablePointerType(t types.Type) bool {
	if p, ok := t.(*types.Pointer); ok {
		switch vt := p.Elem().(type) {
		case *types.Basic:
			return isBasicSerializableType(vt)
		case *types.Map:
			return isSerializableMapType(vt)
		case *types.Slice:
			return isSerializableSliceType(vt)
		case *types.Struct:
			return isSerializableStructureType(vt)
		default:
			return false
		}
	}
	return false
}

func isSerializableSliceType(t types.Type) bool {
	if s, ok := t.(*types.Slice); ok {
		switch et := s.Elem().(type) {
		case *types.Basic:
			return isBasicSerializableType(et)
		case *types.Map:
			return isSerializableMapType(et)
		case *types.Slice:
			return isSerializableSliceType(et)
		case *types.Struct:
			return isSerializableStructureType(et)
		case *types.Pointer:
			return isSerializablePointerType(et)
		default:
			return false
		}
	}
	return false
}

func isSerializableMapType(t types.Type) bool {
	if m, ok := t.(*types.Map); ok {
		if !isBasicSerializableType(m.Key()) {
			return false
		}

		switch vt := m.Elem().(type) {
		case *types.Basic:
			return isBasicSerializableType(vt)
		case *types.Map:
			return isSerializableMapType(vt)
		case *types.Slice:
			return isSerializableSliceType(vt)
		case *types.Struct:
			return isSerializableStructureType(vt)
		case *types.Pointer:
			return isSerializablePointerType(vt)
		default:
			return false
		}
	}
	return false
}

func isBasicSerializableType(t types.Type) bool {
	if b, ok := t.(*types.Basic); ok {
		switch b.Kind() {
		case types.Int,
			types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64,
			types.String, types.Bool:
			return true
		default:
			return false
		}
	}
	return false
}
