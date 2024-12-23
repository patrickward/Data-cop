package datacop

import (
	"reflect"
	"regexp"
	"strings"
	"time"
)

// Required checks if a value is non-empty
//
// Example usage:
// Required("some value") // returns true
// Required("") // returns false
// Required([]int{1, 2, 3}) // returns true
// Required([]int{}) // returns false
func Required(value any) bool {
	if value == nil {
		return false
	}

	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		return strings.TrimSpace(v.String()) != ""
	case reflect.Slice, reflect.Array:
		return v.Len() > 0
	case reflect.Map:
		return v.Len() > 0
	case reflect.Struct:
		if t, ok := value.(time.Time); ok {
			return !t.IsZero()
		}
		// For other structs, we could either:
		// 1. Consider them always required (return true)
		// 2. Check if they're zero value (return !reflect.DeepEqual(value, reflect.Zero(v.Type()).Interface()))
		return true
	case reflect.Ptr:
		if v.IsNil() {
			return false
		}
		return Required(v.Elem().Interface())
	default:
		// For basic types (int, float, etc), check if they're zero value
		return value != reflect.Zero(v.Type()).Interface()
	}
}

// NotZero checks if a numeric value is not a zero value
func NotZero[T comparable](value T) bool {
	return value != *new(T)
}

// Match returns a validation function that checks if a string matches a pattern
//
// Example usage:
// Match(`^[a-zA-Z0-9]+$`)(username) // returns true if username is alphanumeric
// Match(`^[a-zA-Z0-9]+$`)(email) // returns false if email is not alphanumeric
func Match(pattern string) ValidationFunc {
	regex := regexp.MustCompile(pattern)
	return func(value any) bool {
		str, ok := value.(string)
		if !ok {
			return false
		}
		return regex.MatchString(str)
	}
}
