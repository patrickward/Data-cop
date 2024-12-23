package datacop_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/patrickward/datacop"
)

func TestPassword(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid password", "Password1", true},
		{"missing uppercase", "password1", false},
		{"missing lowercase", "PASSWORD1", false},
		{"missing digit", "Password", false},
		{"too short", "Pass1", false},
		{"non-string value", 12345678, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, datacop.Password(tt.value))
		})
	}
}

func TestUsername(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  bool
	}{
		{"valid username", "user_name123", true},
		{"too short", "us", false},
		{"too long", "a" + string(make([]byte, 255)), false},
		{"invalid characters", "user@name", false},
		{"non-string value", 12345, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, datacop.Username(tt.value))
		})
	}
}
