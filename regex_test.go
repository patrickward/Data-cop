package datacop_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/patrickward/datacop"
)

func TestEmail(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid email", "foo@example.com", true},
		{"invalid email", "invalid-email", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, datacop.Email(tt.value))
		})
	}
}

func TestPhone(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{"valid phone", "123-456-7890", true},
		{"invalid phone", "12345", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, datacop.Phone(tt.value))
		})
	}
}
