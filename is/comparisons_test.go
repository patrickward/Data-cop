package is_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/patrickward/datacop"
	"github.com/patrickward/datacop/is"
)

func TestBetween(t *testing.T) {
	tests := []struct {
		name  string
		min   int
		max   int
		value int
		want  bool
	}{
		{"value within range", 10, 20, 15, true},
		{"value below range", 10, 20, 5, false},
		{"value above range", 10, 20, 25, false},
		{"value at min boundary", 10, 20, 10, true},
		{"value at max boundary", 10, 20, 20, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, is.Between(tt.min, tt.max)(tt.value))
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name  string
		other any
		value any
		want  bool
	}{
		{"equal integers", 10, 10, true},
		{"unequal integers", 10, 5, false},
		{"equal strings", "test", "test", true},
		{"unequal strings", "test", "TEST", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, is.Equal(tt.other)(tt.value))
		})
	}
}

func TestEqualStrings(t *testing.T) {
	tests := []struct {
		name  string
		value string
		other string
		want  bool
	}{
		{"equal strings", "test", "test", true},
		{"unequal strings", "test", "TEST", false},
		{"equal strings with spaces", " test ", "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, is.EqualStrings(tt.value, tt.other))
		})
	}
}

func TestInValidation(t *testing.T) {
	tests := []struct {
		name     string
		in       interface{}
		value    interface{}
		expected bool
	}{
		// Integer values
		{"int match", is.In(1, 2, 3), 2, true},
		{"int no match", is.In(1, 2, 3), 4, false},
		{"int8", is.In[int8](1, 2, 3), int8(2), true},
		{"int16", is.In[int16](1, 2, 3), int16(4), false},
		{"int32", is.In[int32](1, 2, 3), int32(2), true},
		{"int64", is.In[int64](1, 2, 3), int64(4), false},

		// String values
		{"string match", is.In("a", "b", "c"), "b", true},
		{"string no match", is.In("a", "b", "c"), "d", false},
		{"empty string", is.In("a", "", "c"), "", true},

		// Mixed types should fail
		{"type mismatch", is.In(1, 2, 3), "2", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := tt.in.(datacop.ValidationFunc)
			result := validator(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAllInValidation(t *testing.T) {
	tests := []struct {
		name     string
		allIn    interface{}
		value    interface{}
		expected bool
	}{
		// Integer slices
		{"int all match", is.AllIn(1, 2, 3), []int{1, 2}, true},
		{"int some match", is.AllIn(1, 2, 3), []int{1, 4}, false},
		{"int empty allowed", is.AllIn(1, 2, 3), []int{}, true},
		{"int8 match", is.AllIn[int8](1, 2, 3), []int8{1, 2}, true},
		{"int16 no match", is.AllIn[int16](1, 2, 3), []int16{1, 4}, false},

		// String slices
		{"string all match", is.AllIn("a", "b", "c"), []string{"a", "b"}, true},
		{"string some match", is.AllIn("a", "b", "c"), []string{"a", "d"}, false},
		{"string empty allowed", is.AllIn("a", "b", "c"), []string{}, true},

		// Type mismatches
		{"wrong slice type", is.AllIn(1, 2, 3), []string{"1", "2"}, false},
		{"not a slice", is.AllIn(1, 2, 3), 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := tt.allIn.(datacop.ValidationFunc)
			result := validator(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNoDuplicatesValidation(t *testing.T) {
	tests := []struct {
		name     string
		noDups   interface{}
		value    interface{}
		expected bool
	}{
		// Integer slices
		{"int no duplicates", is.NoDuplicates[int](), []int{1, 2, 3}, true},
		{"int with duplicates", is.NoDuplicates[int](), []int{1, 2, 2}, false},
		{"int empty slice", is.NoDuplicates[int](), []int{}, true},
		{"int8 no duplicates", is.NoDuplicates[int8](), []int8{1, 2, 3}, true},
		{"int16 with duplicates", is.NoDuplicates[int16](), []int16{1, 2, 2}, false},

		// String slices
		{"string no duplicates", is.NoDuplicates[string](), []string{"a", "b", "c"}, true},
		{"string with duplicates", is.NoDuplicates[string](), []string{"a", "b", "b"}, false},
		{"string empty slice", is.NoDuplicates[string](), []string{}, true},

		// Type mismatches
		{"wrong slice type", is.NoDuplicates[int](), []string{"1", "2"}, false},
		{"not a slice", is.NoDuplicates[int](), 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := tt.noDups.(datacop.ValidationFunc)
			result := validator(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		name  string
		min   int
		value string
		want  bool
	}{
		{"length greater than min", 5, "hello", true},
		{"length equal to min", 5, "hello", true},
		{"length less than min", 5, "hi", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, is.MinLength(tt.min)(tt.value))
		})
	}
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		name  string
		max   int
		value string
		want  bool
	}{
		{"length less than max", 5, "hi", true},
		{"length equal to max", 5, "hello", true},
		{"length greater than max", 5, "hello world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, is.MaxLength(tt.max)(tt.value))
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		min      interface{}
		value    interface{}
		expected bool
	}{
		// Integer types
		{"int min", is.Min(5), 10, true},
		{"int equal", is.Min(5), 5, true},
		{"int below", is.Min(5), 3, false},
		{"int8", is.Min[int8](5), int8(10), true},
		{"int16", is.Min[int16](5), int16(3), false},
		{"int32", is.Min[int32](5), int32(5), true},
		{"int64", is.Min[int64](5), int64(3), false},

		// Unsigned integer types
		{"uint", is.Min[uint](5), uint(10), true},
		{"uint8", is.Min[uint8](5), uint8(3), false},
		{"uint16", is.Min[uint16](5), uint16(5), true},
		{"uint32", is.Min[uint32](5), uint32(10), true},
		{"uint64", is.Min[uint64](5), uint64(3), false},

		// Float types
		{"float32", is.Min[float32](5.5), float32(10.5), true},
		{"float32 equal", is.Min[float32](5.5), float32(5.5), true},
		{"float32 below", is.Min[float32](5.5), float32(3.5), false},
		{"float64", is.Min[float64](5.5), 10.5, true},
		{"float64 equal", is.Min[float64](5.5), 5.5, true},
		{"float64 below", is.Min[float64](5.5), 3.5, false},

		// String types
		{"string", is.Min("b"), "c", true},
		{"string equal", is.Min("b"), "b", true},
		{"string below", is.Min("b"), "a", false},

		// Type mismatch
		{"type mismatch int/float", is.Min(5), 5.5, false},
		{"type mismatch float/int", is.Min(5.5), 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := tt.min.(datacop.ValidationFunc)
			result := validator(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		max      interface{}
		value    interface{}
		expected bool
	}{
		// Integer types
		{"int max", is.Max(5), 3, true},
		{"int equal", is.Max(5), 5, true},
		{"int above", is.Max(5), 10, false},
		{"int8", is.Max[int8](5), int8(3), true},
		{"int16", is.Max[int16](5), int16(10), false},
		{"int32", is.Max[int32](5), int32(5), true},
		{"int64", is.Max[int64](5), int64(10), false},

		// Unsigned integer types
		{"uint", is.Max[uint](5), uint(3), true},
		{"uint8", is.Max[uint8](5), uint8(10), false},
		{"uint16", is.Max[uint16](5), uint16(5), true},
		{"uint32", is.Max[uint32](5), uint32(3), true},
		{"uint64", is.Max[uint64](5), uint64(10), false},

		// Float types
		{"float32", is.Max[float32](5.5), float32(3.5), true},
		{"float32 equal", is.Max[float32](5.5), float32(5.5), true},
		{"float32 above", is.Max[float32](5.5), float32(10.5), false},
		{"float64", is.Max[float64](5.5), 3.5, true},
		{"float64 equal", is.Max[float64](5.5), 5.5, true},
		{"float64 above", is.Max[float64](5.5), 10.5, false},

		// String types
		{"string", is.Max("b"), "a", true},
		{"string equal", is.Max("b"), "b", true},
		{"string above", is.Max("b"), "c", false},

		// Type mismatch
		{"type mismatch int/float", is.Max(5), 5.5, false},
		{"type mismatch float/int", is.Max(5.5), 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := tt.max.(datacop.ValidationFunc)
			result := validator(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGreaterThan(t *testing.T) {
	tests := []struct {
		name     string
		max      interface{}
		value    interface{}
		expected bool
	}{
		// Integer types
		{"int greater", is.GreaterThan(5), 10, true},
		{"int equal", is.GreaterThan(5), 5, false},
		{"int below", is.GreaterThan(5), 3, false},
		{"int8", is.GreaterThan[int8](5), int8(10), true},
		{"int16", is.GreaterThan[int16](5), int16(3), false},
		{"int32", is.GreaterThan[int32](5), int32(5), false},
		{"int64", is.GreaterThan[int64](5), int64(10), true},

		// Unsigned integer types
		{"uint", is.GreaterThan[uint](5), uint(10), true},
		{"uint8", is.GreaterThan[uint8](5), uint8(3), false},
		{"uint16", is.GreaterThan[uint16](5), uint16(5), false},
		{"uint32", is.GreaterThan[uint32](5), uint32(10), true},
		{"uint64", is.GreaterThan[uint64](5), uint64(3), false},

		// Float types
		{"float32", is.GreaterThan[float32](5.5), float32(10.5), true},
		{"float32 equal", is.GreaterThan[float32](5.5), float32(5.5), false},
		{"float32 below", is.GreaterThan[float32](5.5), float32(3.5), false},
		{"float64", is.GreaterThan[float64](5.5), 10.5, true},
		{"float64 equal", is.GreaterThan[float64](5.5), 5.5, false},
		{"float64 below", is.GreaterThan[float64](5.5), 3.5, false},

		// String types
		{"string", is.GreaterThan("b"), "c", true},
		{"string equal", is.GreaterThan("b"), "b", false},
		{"string below", is.GreaterThan("b"), "a", false},

		// Type mismatch
		{"type mismatch int/float", is.GreaterThan(5), 5.5, false},
		{"type mismatch float/int", is.GreaterThan(5.5), 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := tt.max.(datacop.ValidationFunc)
			result := validator(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLessThan(t *testing.T) {
	tests := []struct {
		name     string
		max      interface{}
		value    interface{}
		expected bool
	}{
		// Integer types
		{"int less", is.LessThan(5), 3, true},
		{"int equal", is.LessThan(5), 5, false},
		{"int above", is.LessThan(5), 10, false},
		{"int8", is.LessThan[int8](5), int8(3), true},
		{"int16", is.LessThan[int16](5), int16(10), false},
		{"int32", is.LessThan[int32](5), int32(5), false},
		{"int64", is.LessThan[int64](5), int64(3), true},

		// Unsigned integer types
		{"uint", is.LessThan[uint](5), uint(3), true},
		{"uint8", is.LessThan[uint8](5), uint8(10), false},
		{"uint16", is.LessThan[uint16](5), uint16(5), false},
		{"uint32", is.LessThan[uint32](5), uint32(3), true},
		{"uint64", is.LessThan[uint64](5), uint64(10), false},

		// Float types
		{"float32", is.LessThan[float32](5.5), float32(3.5), true},
		{"float32 equal", is.LessThan[float32](5.5), float32(5.5), false},
		{"float32 above", is.LessThan[float32](5.5), float32(10.5), false},
		{"float64", is.LessThan[float64](5.5), 3.5, true},
		{"float64 equal", is.LessThan[float64](5.5), 5.5, false},
		{"float64 above", is.LessThan[float64](5.5), 10.5, false},

		// String types
		{"string", is.LessThan("b"), "a", true},
		{"string equal", is.LessThan("b"), "b", false},
		{"string above", is.LessThan("b"), "c", false},

		// Type mismatch
		{"type mismatch int/float", is.LessThan(5), 5.5, false},
		{"type mismatch float/int", is.LessThan(5.5), 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := tt.max.(datacop.ValidationFunc)
			result := validator(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGreaterOrEqual(t *testing.T) {
	tests := []struct {
		name     string
		min      interface{}
		value    interface{}
		expected bool
	}{
		// Integer types
		{"int greater", is.GreaterOrEqual(5), 10, true},
		{"int equal", is.GreaterOrEqual(5), 5, true},
		{"int below", is.GreaterOrEqual(5), 3, false},
		{"int8", is.GreaterOrEqual[int8](5), int8(10), true},
		{"int16", is.GreaterOrEqual[int16](5), int16(3), false},
		{"int32", is.GreaterOrEqual[int32](5), int32(5), true},
		{"int64", is.GreaterOrEqual[int64](5), int64(10), true},

		// Unsigned integer types
		{"uint", is.GreaterOrEqual[uint](5), uint(10), true},
		{"uint8", is.GreaterOrEqual[uint8](5), uint8(3), false},
		{"uint16", is.GreaterOrEqual[uint16](5), uint16(5), true},
		{"uint32", is.GreaterOrEqual[uint32](5), uint32(10), true},
		{"uint64", is.GreaterOrEqual[uint64](5), uint64(3), false},

		// Float types
		{"float32", is.GreaterOrEqual[float32](5.5), float32(10.5), true},
		{"float32 equal", is.GreaterOrEqual[float32](5.5), float32(5.5), true},
		{"float32 below", is.GreaterOrEqual[float32](5.5), float32(3.5), false},
		{"float64", is.GreaterOrEqual[float64](5.5), 10.5, true},
		{"float64 equal", is.GreaterOrEqual[float64](5.5), 5.5, true},
		{"float64 below", is.GreaterOrEqual[float64](5.5), 3.5, false},

		// String types
		{"string greater", is.GreaterOrEqual("b"), "c", true},
		{"string equal", is.GreaterOrEqual("b"), "b", true},
		{"string below", is.GreaterOrEqual("b"), "a", false},

		// Type mismatch
		{"type mismatch int/float", is.GreaterOrEqual(5), 5.5, false},
		{"type mismatch float/int", is.GreaterOrEqual(5.5), 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := tt.min.(datacop.ValidationFunc)
			result := validator(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLessOrEqual(t *testing.T) {
	tests := []struct {
		name     string
		max      interface{}
		value    interface{}
		expected bool
	}{
		// Integer types
		{"int less", is.LessOrEqual(5), 3, true},
		{"int equal", is.LessOrEqual(5), 5, true},
		{"int above", is.LessOrEqual(5), 10, false},
		{"int8", is.LessOrEqual[int8](5), int8(3), true},
		{"int16", is.LessOrEqual[int16](5), int16(10), false},
		{"int32", is.LessOrEqual[int32](5), int32(5), true},
		{"int64", is.LessOrEqual[int64](5), int64(3), true},

		// Unsigned integer types
		{"uint", is.LessOrEqual[uint](5), uint(3), true},
		{"uint8", is.LessOrEqual[uint8](5), uint8(10), false},
		{"uint16", is.LessOrEqual[uint16](5), uint16(5), true},
		{"uint32", is.LessOrEqual[uint32](5), uint32(3), true},
		{"uint64", is.LessOrEqual[uint64](5), uint64(10), false},

		// Float types
		{"float32", is.LessOrEqual[float32](5.5), float32(3.5), true},
		{"float32 equal", is.LessOrEqual[float32](5.5), float32(5.5), true},
		{"float32 above", is.LessOrEqual[float32](5.5), float32(10.5), false},
		{"float64", is.LessOrEqual[float64](5.5), 3.5, true},
		{"float64 equal", is.LessOrEqual[float64](5.5), 5.5, true},
		{"float64 above", is.LessOrEqual[float64](5.5), 10.5, false},

		// String types
		{"string less", is.LessOrEqual("b"), "a", true},
		{"string equal", is.LessOrEqual("b"), "b", true},
		{"string above", is.LessOrEqual("b"), "c", false},

		// Type mismatch
		{"type mismatch int/float", is.LessOrEqual(5), 5.5, false},
		{"type mismatch float/int", is.LessOrEqual(5.5), 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := tt.max.(datacop.ValidationFunc)
			result := validator(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}
