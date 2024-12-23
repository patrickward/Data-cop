package is

import (
	"cmp"
	"strings"
	"unicode/utf8"

	"golang.org/x/exp/constraints"

	"github.com/patrickward/datacop"
)

// Between checks if a value is between a minimum and maximum value
// Example usage:
//
// Between(10, 20)(15) // returns true
// Between(10, 20)(25) // returns false
func Between[T constraints.Ordered](min, max T) datacop.ValidationFunc {
	return func(value any) bool {
		v, ok := value.(T)
		if !ok {
			return false
		}
		return v >= min && v <= max
	}
}

// Equal checks if values are equal
// Example usage:
// Equal(10)(10) // returns true
// Equal(10)(5) // returns false
func Equal[T comparable](other T) datacop.ValidationFunc {
	return func(value any) bool {
		v, ok := value.(T)
		if !ok {
			return false
		}
		return v == other
	}
}

// EqualStrings checks if two strings are equal
//
// Example usage:
// EqualStrings("test", "test") // returns true
// EqualStrings("test", "TEST") // returns false
func EqualStrings(value, other string) bool {
	return strings.TrimSpace(value) == strings.TrimSpace(other)
}

// In checks if a value is in a set of allowed values
//
// Example usage:
// In(1, 2, 3)(2) // returns true
// In(1, 2, 3)(4) // returns false
// In("a", "b", "c")("b") // returns true
// In("a", "b", "c")("d") // returns false
func In[T comparable](allowed ...T) datacop.ValidationFunc {
	return func(value any) bool {
		v, ok := value.(T)
		if !ok {
			return false
		}
		for _, a := range allowed {
			if v == a {
				return true
			}
		}
		return false
	}
}

// AllIn checks if all values in a slice are in a set of allowed values
//
// Example usage:
// AllIn(1, 2, 3)([]int{1, 2}) // returns true
// AllIn(1, 2, 3)([]int{1, 4}) // returns false
// AllIn("a", "b", "c")([]string{"a", "b"}) // returns true
// AllIn("a", "b", "c")([]string{"a", "d"}) // returns false
func AllIn[T comparable](allowed ...T) datacop.ValidationFunc {
	return func(value any) bool {
		values, ok := value.([]T)
		if !ok {
			return false
		}

		// Create a map for O(1) lookups (keeping this optimization)
		allowedMap := make(map[T]struct{}, len(allowed))
		for _, a := range allowed {
			allowedMap[a] = struct{}{}
		}

		for _, v := range values {
			if _, exists := allowedMap[v]; !exists {
				return false
			}
		}
		return true
	}
}

// NoDuplicates checks if a slice contains any duplicate values
//
// Example usage:
// NoDuplicates()([]int{1, 2, 3}) // returns true
// NoDuplicates()([]int{1, 2, 2}) // returns false
func NoDuplicates[T comparable]() datacop.ValidationFunc {
	return func(value any) bool {
		values, ok := value.([]T)
		if !ok {
			return false
		}

		seen := make(map[T]struct{}, len(values))
		for _, v := range values {
			if _, exists := seen[v]; exists {
				return false
			}
			seen[v] = struct{}{}
		}
		return true
	}
}

// MinLength returns a validation function that checks minimum string length
//
// Example usage:
// MinLength(5)("hello") // returns true
// MinLength(5)("hi") // returns false
func MinLength(min int) datacop.ValidationFunc {
	return func(value any) bool {
		str, ok := value.(string)
		if !ok {
			return false
		}
		return utf8.RuneCountInString(strings.TrimSpace(str)) >= min
	}
}

// MaxLength returns a validation function that checks maximum string length
//
// Example usage:
// MaxLength(5)("hello") // returns false
// MaxLength(5)("hi") // returns true
func MaxLength(max int) datacop.ValidationFunc {
	return func(value any) bool {
		str, ok := value.(string)
		if !ok {
			return false
		}
		return utf8.RuneCountInString(strings.TrimSpace(str)) <= max
	}
}

// Min returns a validation function that checks minimum value
//
// Example usage:
// Min(10)(15) // returns true
// Min(10)(5) // returns false
func Min[T cmp.Ordered](min T) datacop.ValidationFunc {
	return func(value any) bool {
		v, ok := value.(T)
		if !ok {
			return false
		}
		return v >= min
	}
}

// Max returns a validation function that checks maximum value
//
// Example usage:
// Max(10)(5) // returns true
// Max(10)(15) // returns false
func Max[T cmp.Ordered](max T) datacop.ValidationFunc {
	return func(value any) bool {
		v, ok := value.(T)
		if !ok {
			return false
		}
		return v <= max
	}
}

// GreaterThan returns a validation function that checks if a value is greater than a specified value
//
// Example usage:
//
// GreaterThan(10)(15) // returns true
// GreaterThan(10)(5) // returns false
func GreaterThan[T cmp.Ordered](n T) datacop.ValidationFunc {
	return func(value any) bool {
		v, ok := value.(T)
		if !ok {
			return false
		}
		return v > n
	}
}

// LessThan returns a validation function that checks if a value is less than a specified value
//
// Example usage:
// LessThan(10)(5) // returns true
// LessThan(10)(15) // returns false
func LessThan[T cmp.Ordered](n T) datacop.ValidationFunc {
	return func(value any) bool {
		v, ok := value.(T)
		if !ok {
			return false
		}
		return v < n
	}
}

// GreaterOrEqual returns a validation function that checks if a value is greater than or equal to a specified value
//
// Example usage:
// GreaterOrEqual(10)(15) // returns true
// GreaterOrEqual(10)(10) // returns true
// GreaterOrEqual(10)(5) // returns false
func GreaterOrEqual[T cmp.Ordered](n T) datacop.ValidationFunc {
	return func(value any) bool {
		v, ok := value.(T)
		if !ok {
			return false
		}
		return v >= n
	}
}

// LessOrEqual returns a validation function that checks if a value is less than or equal to a specified value
//
// Example usage:
// LessOrEqual(10)(5) // returns true
// LessOrEqual(10)(10) // returns true
// LessOrEqual(10)(15) // returns false
func LessOrEqual[T cmp.Ordered](n T) datacop.ValidationFunc {
	return func(value any) bool {
		v, ok := value.(T)
		if !ok {
			return false
		}
		return v <= n
	}
}
