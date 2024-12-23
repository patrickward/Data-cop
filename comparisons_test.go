package datacop_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/patrickward/datacop"
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
			assert.Equal(t, tt.want, datacop.Between(tt.min, tt.max)(tt.value))
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
			assert.Equal(t, tt.want, datacop.Equal(tt.other)(tt.value))
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
			assert.Equal(t, tt.want, datacop.EqualStrings(tt.value, tt.other))
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
		{"int match", datacop.In(1, 2, 3), 2, true},
		{"int no match", datacop.In(1, 2, 3), 4, false},
		{"int8", datacop.In[int8](1, 2, 3), int8(2), true},
		{"int16", datacop.In[int16](1, 2, 3), int16(4), false},
		{"int32", datacop.In[int32](1, 2, 3), int32(2), true},
		{"int64", datacop.In[int64](1, 2, 3), int64(4), false},

		// String values
		{"string match", datacop.In("a", "b", "c"), "b", true},
		{"string no match", datacop.In("a", "b", "c"), "d", false},
		{"empty string", datacop.In("a", "", "c"), "", true},

		// Mixed types should fail
		{"type mismatch", datacop.In(1, 2, 3), "2", false},
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
		{"int all match", datacop.AllIn(1, 2, 3), []int{1, 2}, true},
		{"int some match", datacop.AllIn(1, 2, 3), []int{1, 4}, false},
		{"int empty allowed", datacop.AllIn(1, 2, 3), []int{}, true},
		{"int8 match", datacop.AllIn[int8](1, 2, 3), []int8{1, 2}, true},
		{"int16 no match", datacop.AllIn[int16](1, 2, 3), []int16{1, 4}, false},

		// String slices
		{"string all match", datacop.AllIn("a", "b", "c"), []string{"a", "b"}, true},
		{"string some match", datacop.AllIn("a", "b", "c"), []string{"a", "d"}, false},
		{"string empty allowed", datacop.AllIn("a", "b", "c"), []string{}, true},

		// Type mismatches
		{"wrong slice type", datacop.AllIn(1, 2, 3), []string{"1", "2"}, false},
		{"not a slice", datacop.AllIn(1, 2, 3), 1, false},
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
		{"int no duplicates", datacop.NoDuplicates[int](), []int{1, 2, 3}, true},
		{"int with duplicates", datacop.NoDuplicates[int](), []int{1, 2, 2}, false},
		{"int empty slice", datacop.NoDuplicates[int](), []int{}, true},
		{"int8 no duplicates", datacop.NoDuplicates[int8](), []int8{1, 2, 3}, true},
		{"int16 with duplicates", datacop.NoDuplicates[int16](), []int16{1, 2, 2}, false},

		// String slices
		{"string no duplicates", datacop.NoDuplicates[string](), []string{"a", "b", "c"}, true},
		{"string with duplicates", datacop.NoDuplicates[string](), []string{"a", "b", "b"}, false},
		{"string empty slice", datacop.NoDuplicates[string](), []string{}, true},

		// Type mismatches
		{"wrong slice type", datacop.NoDuplicates[int](), []string{"1", "2"}, false},
		{"not a slice", datacop.NoDuplicates[int](), 1, false},
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
			assert.Equal(t, tt.want, datacop.MinLength(tt.min)(tt.value))
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
			assert.Equal(t, tt.want, datacop.MaxLength(tt.max)(tt.value))
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
		{"int min", datacop.Min(5), 10, true},
		{"int equal", datacop.Min(5), 5, true},
		{"int below", datacop.Min(5), 3, false},
		{"int8", datacop.Min[int8](5), int8(10), true},
		{"int16", datacop.Min[int16](5), int16(3), false},
		{"int32", datacop.Min[int32](5), int32(5), true},
		{"int64", datacop.Min[int64](5), int64(3), false},

		// Unsigned integer types
		{"uint", datacop.Min[uint](5), uint(10), true},
		{"uint8", datacop.Min[uint8](5), uint8(3), false},
		{"uint16", datacop.Min[uint16](5), uint16(5), true},
		{"uint32", datacop.Min[uint32](5), uint32(10), true},
		{"uint64", datacop.Min[uint64](5), uint64(3), false},

		// Float types
		{"float32", datacop.Min[float32](5.5), float32(10.5), true},
		{"float32 equal", datacop.Min[float32](5.5), float32(5.5), true},
		{"float32 below", datacop.Min[float32](5.5), float32(3.5), false},
		{"float64", datacop.Min[float64](5.5), 10.5, true},
		{"float64 equal", datacop.Min[float64](5.5), 5.5, true},
		{"float64 below", datacop.Min[float64](5.5), 3.5, false},

		// String types
		{"string", datacop.Min("b"), "c", true},
		{"string equal", datacop.Min("b"), "b", true},
		{"string below", datacop.Min("b"), "a", false},

		// Type mismatch
		{"type mismatch int/float", datacop.Min(5), 5.5, false},
		{"type mismatch float/int", datacop.Min(5.5), 6, false},
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
		{"int max", datacop.Max(5), 3, true},
		{"int equal", datacop.Max(5), 5, true},
		{"int above", datacop.Max(5), 10, false},
		{"int8", datacop.Max[int8](5), int8(3), true},
		{"int16", datacop.Max[int16](5), int16(10), false},
		{"int32", datacop.Max[int32](5), int32(5), true},
		{"int64", datacop.Max[int64](5), int64(10), false},

		// Unsigned integer types
		{"uint", datacop.Max[uint](5), uint(3), true},
		{"uint8", datacop.Max[uint8](5), uint8(10), false},
		{"uint16", datacop.Max[uint16](5), uint16(5), true},
		{"uint32", datacop.Max[uint32](5), uint32(3), true},
		{"uint64", datacop.Max[uint64](5), uint64(10), false},

		// Float types
		{"float32", datacop.Max[float32](5.5), float32(3.5), true},
		{"float32 equal", datacop.Max[float32](5.5), float32(5.5), true},
		{"float32 above", datacop.Max[float32](5.5), float32(10.5), false},
		{"float64", datacop.Max[float64](5.5), 3.5, true},
		{"float64 equal", datacop.Max[float64](5.5), 5.5, true},
		{"float64 above", datacop.Max[float64](5.5), 10.5, false},

		// String types
		{"string", datacop.Max("b"), "a", true},
		{"string equal", datacop.Max("b"), "b", true},
		{"string above", datacop.Max("b"), "c", false},

		// Type mismatch
		{"type mismatch int/float", datacop.Max(5), 5.5, false},
		{"type mismatch float/int", datacop.Max(5.5), 6, false},
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
		{"int greater", datacop.GreaterThan(5), 10, true},
		{"int equal", datacop.GreaterThan(5), 5, false},
		{"int below", datacop.GreaterThan(5), 3, false},
		{"int8", datacop.GreaterThan[int8](5), int8(10), true},
		{"int16", datacop.GreaterThan[int16](5), int16(3), false},
		{"int32", datacop.GreaterThan[int32](5), int32(5), false},
		{"int64", datacop.GreaterThan[int64](5), int64(10), true},

		// Unsigned integer types
		{"uint", datacop.GreaterThan[uint](5), uint(10), true},
		{"uint8", datacop.GreaterThan[uint8](5), uint8(3), false},
		{"uint16", datacop.GreaterThan[uint16](5), uint16(5), false},
		{"uint32", datacop.GreaterThan[uint32](5), uint32(10), true},
		{"uint64", datacop.GreaterThan[uint64](5), uint64(3), false},

		// Float types
		{"float32", datacop.GreaterThan[float32](5.5), float32(10.5), true},
		{"float32 equal", datacop.GreaterThan[float32](5.5), float32(5.5), false},
		{"float32 below", datacop.GreaterThan[float32](5.5), float32(3.5), false},
		{"float64", datacop.GreaterThan[float64](5.5), 10.5, true},
		{"float64 equal", datacop.GreaterThan[float64](5.5), 5.5, false},
		{"float64 below", datacop.GreaterThan[float64](5.5), 3.5, false},

		// String types
		{"string", datacop.GreaterThan("b"), "c", true},
		{"string equal", datacop.GreaterThan("b"), "b", false},
		{"string below", datacop.GreaterThan("b"), "a", false},

		// Type mismatch
		{"type mismatch int/float", datacop.GreaterThan(5), 5.5, false},
		{"type mismatch float/int", datacop.GreaterThan(5.5), 6, false},
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
		{"int less", datacop.LessThan(5), 3, true},
		{"int equal", datacop.LessThan(5), 5, false},
		{"int above", datacop.LessThan(5), 10, false},
		{"int8", datacop.LessThan[int8](5), int8(3), true},
		{"int16", datacop.LessThan[int16](5), int16(10), false},
		{"int32", datacop.LessThan[int32](5), int32(5), false},
		{"int64", datacop.LessThan[int64](5), int64(3), true},

		// Unsigned integer types
		{"uint", datacop.LessThan[uint](5), uint(3), true},
		{"uint8", datacop.LessThan[uint8](5), uint8(10), false},
		{"uint16", datacop.LessThan[uint16](5), uint16(5), false},
		{"uint32", datacop.LessThan[uint32](5), uint32(3), true},
		{"uint64", datacop.LessThan[uint64](5), uint64(10), false},

		// Float types
		{"float32", datacop.LessThan[float32](5.5), float32(3.5), true},
		{"float32 equal", datacop.LessThan[float32](5.5), float32(5.5), false},
		{"float32 above", datacop.LessThan[float32](5.5), float32(10.5), false},
		{"float64", datacop.LessThan[float64](5.5), 3.5, true},
		{"float64 equal", datacop.LessThan[float64](5.5), 5.5, false},
		{"float64 above", datacop.LessThan[float64](5.5), 10.5, false},

		// String types
		{"string", datacop.LessThan("b"), "a", true},
		{"string equal", datacop.LessThan("b"), "b", false},
		{"string above", datacop.LessThan("b"), "c", false},

		// Type mismatch
		{"type mismatch int/float", datacop.LessThan(5), 5.5, false},
		{"type mismatch float/int", datacop.LessThan(5.5), 6, false},
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
		{"int greater", datacop.GreaterOrEqual(5), 10, true},
		{"int equal", datacop.GreaterOrEqual(5), 5, true},
		{"int below", datacop.GreaterOrEqual(5), 3, false},
		{"int8", datacop.GreaterOrEqual[int8](5), int8(10), true},
		{"int16", datacop.GreaterOrEqual[int16](5), int16(3), false},
		{"int32", datacop.GreaterOrEqual[int32](5), int32(5), true},
		{"int64", datacop.GreaterOrEqual[int64](5), int64(10), true},

		// Unsigned integer types
		{"uint", datacop.GreaterOrEqual[uint](5), uint(10), true},
		{"uint8", datacop.GreaterOrEqual[uint8](5), uint8(3), false},
		{"uint16", datacop.GreaterOrEqual[uint16](5), uint16(5), true},
		{"uint32", datacop.GreaterOrEqual[uint32](5), uint32(10), true},
		{"uint64", datacop.GreaterOrEqual[uint64](5), uint64(3), false},

		// Float types
		{"float32", datacop.GreaterOrEqual[float32](5.5), float32(10.5), true},
		{"float32 equal", datacop.GreaterOrEqual[float32](5.5), float32(5.5), true},
		{"float32 below", datacop.GreaterOrEqual[float32](5.5), float32(3.5), false},
		{"float64", datacop.GreaterOrEqual[float64](5.5), 10.5, true},
		{"float64 equal", datacop.GreaterOrEqual[float64](5.5), 5.5, true},
		{"float64 below", datacop.GreaterOrEqual[float64](5.5), 3.5, false},

		// String types
		{"string greater", datacop.GreaterOrEqual("b"), "c", true},
		{"string equal", datacop.GreaterOrEqual("b"), "b", true},
		{"string below", datacop.GreaterOrEqual("b"), "a", false},

		// Type mismatch
		{"type mismatch int/float", datacop.GreaterOrEqual(5), 5.5, false},
		{"type mismatch float/int", datacop.GreaterOrEqual(5.5), 6, false},
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
		{"int less", datacop.LessOrEqual(5), 3, true},
		{"int equal", datacop.LessOrEqual(5), 5, true},
		{"int above", datacop.LessOrEqual(5), 10, false},
		{"int8", datacop.LessOrEqual[int8](5), int8(3), true},
		{"int16", datacop.LessOrEqual[int16](5), int16(10), false},
		{"int32", datacop.LessOrEqual[int32](5), int32(5), true},
		{"int64", datacop.LessOrEqual[int64](5), int64(3), true},

		// Unsigned integer types
		{"uint", datacop.LessOrEqual[uint](5), uint(3), true},
		{"uint8", datacop.LessOrEqual[uint8](5), uint8(10), false},
		{"uint16", datacop.LessOrEqual[uint16](5), uint16(5), true},
		{"uint32", datacop.LessOrEqual[uint32](5), uint32(3), true},
		{"uint64", datacop.LessOrEqual[uint64](5), uint64(10), false},

		// Float types
		{"float32", datacop.LessOrEqual[float32](5.5), float32(3.5), true},
		{"float32 equal", datacop.LessOrEqual[float32](5.5), float32(5.5), true},
		{"float32 above", datacop.LessOrEqual[float32](5.5), float32(10.5), false},
		{"float64", datacop.LessOrEqual[float64](5.5), 3.5, true},
		{"float64 equal", datacop.LessOrEqual[float64](5.5), 5.5, true},
		{"float64 above", datacop.LessOrEqual[float64](5.5), 10.5, false},

		// String types
		{"string less", datacop.LessOrEqual("b"), "a", true},
		{"string equal", datacop.LessOrEqual("b"), "b", true},
		{"string above", datacop.LessOrEqual("b"), "c", false},

		// Type mismatch
		{"type mismatch int/float", datacop.LessOrEqual(5), 5.5, false},
		{"type mismatch float/int", datacop.LessOrEqual(5.5), 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := tt.max.(datacop.ValidationFunc)
			result := validator(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}
