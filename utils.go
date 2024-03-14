package conditions

import (
	"fmt"
	"reflect"
	"time"
)

// Helper function to check if a slice contains a specific key
func contains[T comparable](s []T, e string) bool {
	for _, a := range s {
		if reflect.DeepEqual(a, e) {
			return true
		}
	}
	return false
}

func isInCollection(collection any, element any) bool {
	val := reflect.ValueOf(collection)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(val.Index(i).Interface(), element) {
				return true
			}
		}
	case reflect.Map:
		value := val.MapIndex(reflect.ValueOf(element))
		if value.IsValid() {
			return true
		}
	}
	return false
}

// Returns 0 if equal, -1 if v1 < v2, 1 if v1 > v2, and an error if incomparable.
func compareNumbersOrDates(v1, v2 any) (int, error) {
	switch v1Typed := v1.(type) {
	case float64, float32, int, int64, int32, int16, int8, uint, uint64, uint32, uint16, uint8:
		f1, _ := toFloat64(v1)
		f2, _ := toFloat64(v2)

		if f1 < f2 {
			return -1, nil
		} else if f1 > f2 {
			return 1, nil
		}
		return 0, nil
	case time.Time:
		t2, ok := v2.(time.Time)
		if !ok {
			return 0, fmt.Errorf("cannot compare time.Time with non-time.Time")
		}
		if v1Typed.Before(t2) {
			return -1, nil
		} else if v1Typed.After(t2) {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported type for comparison")
	}
}

func toFloat64(value any) (float64, bool) {
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(rv.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return float64(rv.Uint()), true
	case reflect.Float32, reflect.Float64:
		return rv.Float(), true
	default:
		return 0, false // Not a numeric type
	}
}
