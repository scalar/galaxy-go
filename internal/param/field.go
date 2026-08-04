package param

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type APIUnion struct{}

type FieldLike interface{ field() }

// Field is a wrapper used for all values sent to the API,
// to distinguish zero values from null or omitted fields.
//
// It also allows sending arbitrary deserializable values.
//
// To instantiate a Field, use the helpers exported from
// the package root: `F()`, `Null()`, `Raw()`, etc.
type Field[T any] struct {
	FieldLike
	Value   T
	Null    bool
	Present bool
	Raw     any
}

func (f Field[T]) String() string {
	if s, ok := any(f.Value).(fmt.Stringer); ok {
		return s.String()
	}
	return fmt.Sprintf("%v", f.Value)
}

func MarshalUnion(_ any, variants ...any) ([]byte, error) {
	var selected any
	count := 0
	for _, variant := range variants {
		if isZeroUnionVariant(variant) {
			continue
		}
		selected = variant
		count += 1
	}
	if count == 0 {
		return []byte("{}"), nil
	}
	if count > 1 {
		return nil, fmt.Errorf("param: cannot marshal union with multiple variants")
	}
	return json.Marshal(selected)
}

func isZeroUnionVariant(value any) bool {
	if value == nil {
		return true
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}
