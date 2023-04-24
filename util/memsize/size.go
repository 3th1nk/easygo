package memsize

import (
	"reflect"
	"sync"
	"unsafe"
)

func Size(v interface{}) int {
	s1, s2 := getSize(reflect.ValueOf(v), make(map[uintptr]bool, 32))
	return s1 + s2
}

func getSize(v reflect.Value, sizedPtr map[uintptr]bool) (typeSize, size int) {
	kind := v.Kind()
	if size, ok := valSizeMap[kind]; ok {
		return 0, size
	}

	typ := v.Type()
	switch kind {
	case reflect.String:
		return strTypeSize, v.Len()

	case reflect.Map:
		typeSize = getTypeSize(typ)
		if keySize, valSize := valSizeMap[typ.Key().Kind()], valSizeMap[typ.Elem().Kind()]; keySize != 0 && valSize != 0 {
			size = (keySize + valSize) * v.Len()
		} else {
			for _, key := range v.MapKeys() {
				s1, s2 := getSize(key, sizedPtr)
				if s2 < 0 {
					return -1, -1
				}
				typeSize += s1
				size += s2

				s1, s2 = getSize(v.MapIndex(key), sizedPtr)
				if s2 < 0 {
					return -1, -1
				}
				typeSize += s1
				size += s2
			}
		}
		return

	case reflect.Slice, reflect.Array:
		typeSize = getTypeSize(typ)
		if elemSize, ok := valSizeMap[typ.Elem().Kind()]; ok {
			size = elemSize * v.Len()
		} else {
			for i, n := 0, v.Len(); i < n; i++ {
				s1, s2 := getSize(v.Index(i), sizedPtr)
				if s2 < 0 {
					return -1, -1
				}
				typeSize += s1
				size += s2
			}
		}
		return

	case reflect.Ptr:
		typeSize = getTypeSize(typ)
		if !v.IsNil() {
			if p := v.Pointer(); !sizedPtr[p] {
				sizedPtr[p] = true
			}
			s1, s2 := getSize(v.Elem(), sizedPtr)
			typeSize += s1
			size += s2
		}
		return

	case reflect.Interface:
		typeSize = getTypeSize(typ)
		if !v.IsNil() {
			s1, s2 := getSize(v.Elem(), sizedPtr)
			typeSize += s1
			size += s2
		}
		return

	case reflect.Struct:
		for i, n := 0, v.NumField(); i < n; i++ {
			s1, s2 := getSize(v.Field(i), sizedPtr)
			if s2 < 0 {
				return -1, -1
			}
			typeSize += s1
			size += s2
		}
		return

	default:
		typeSize = getTypeSize(typ)
		return
	}
}

var (
	strTypeSize = int(reflect.TypeOf("").Size())
	valSizeMap  = map[reflect.Kind]int{
		reflect.Bool:          int(reflect.TypeOf(true).Size()),
		reflect.Int:           int(reflect.TypeOf(int(0)).Size()),
		reflect.Int8:          int(reflect.TypeOf(int8(0)).Size()),
		reflect.Int16:         int(reflect.TypeOf(int16(0)).Size()),
		reflect.Int32:         int(reflect.TypeOf(int32(0)).Size()),
		reflect.Int64:         int(reflect.TypeOf(int64(0)).Size()),
		reflect.Uint:          int(reflect.TypeOf(uint(0)).Size()),
		reflect.Uint8:         int(reflect.TypeOf(uint8(0)).Size()),
		reflect.Uint16:        int(reflect.TypeOf(uint16(0)).Size()),
		reflect.Uint32:        int(reflect.TypeOf(uint32(0)).Size()),
		reflect.Uint64:        int(reflect.TypeOf(uint64(0)).Size()),
		reflect.Uintptr:       int(reflect.TypeOf(uintptr(0)).Size()),
		reflect.Float32:       int(reflect.TypeOf(float32(0)).Size()),
		reflect.Float64:       int(reflect.TypeOf(float64(0)).Size()),
		reflect.Complex64:     int(reflect.TypeOf(complex64(0)).Size()),
		reflect.Complex128:    int(reflect.TypeOf(complex128(0)).Size()),
		reflect.Func:          int(reflect.TypeOf(func() {}).Size()),
		reflect.UnsafePointer: int(reflect.TypeOf(unsafe.Pointer(nil)).Size()),
	}
	typeSizeLock = sync.RWMutex{}
	typeSizeMap  = make(map[reflect.Type]int, 128)
)

func getTypeSize(t reflect.Type) (size int) {
	typeSizeLock.RLock()
	size, ok := typeSizeMap[t]
	typeSizeLock.RUnlock()
	if !ok {
		size = int(t.Size())
		typeSizeLock.Lock()
		typeSizeMap[t] = size
		typeSizeLock.Unlock()
	}
	return
}
