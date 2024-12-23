package is_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/patrickward/datacop/is"
)

func TestRequired(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"non-empty string", "test", true},
		{"empty string", "", false},
		{"non-empty slice", []int{1, 2, 3}, true},
		{"empty slice", []int{}, false},
		{"non-zero time", time.Now(), true},
		{"zero time", time.Time{}, false},
		{"nil value", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, is.Required(tt.value))
		})
	}
}

func TestNotZero(t *testing.T) {
	tests := []struct {
		name  string
		value int
		want  bool
	}{
		{"non-zero int", 1, true},
		{"zero int", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, is.NotZero(tt.value))
		})
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		value   string
		want    bool
	}{
		{"match alphanumeric", `^[a-zA-Z0-9]+$`, "username123", true},
		{"no match alphanumeric", `^[a-zA-Z0-9]+$`, "username@123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, is.Match(tt.pattern)(tt.value))
		})
	}
}
