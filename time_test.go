package datacop_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/patrickward/datacop"
)

func TestBefore(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		t     time.Time
		value time.Time
		want  bool
	}{
		{"before time", now, now.Add(-time.Hour), true},
		{"after time", now, now.Add(time.Hour), false},
		{"equal time", now, now, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, datacop.Before(tt.t)(tt.value))
		})
	}
}

func TestAfter(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		t     time.Time
		value time.Time
		want  bool
	}{
		{"after time", now, now.Add(time.Hour), true},
		{"before time", now, now.Add(-time.Hour), false},
		{"equal time", now, now, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, datacop.After(tt.t)(tt.value))
		})
	}
}

func TestBetweenTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		start time.Time
		end   time.Time
		value time.Time
		want  bool
	}{
		{"within range", now.Add(-time.Hour), now.Add(time.Hour), now, true},
		{"before range", now.Add(-time.Hour), now.Add(time.Hour), now.Add(-2 * time.Hour), false},
		{"after range", now.Add(-time.Hour), now.Add(time.Hour), now.Add(2 * time.Hour), false},
		{"at start boundary", now.Add(-time.Hour), now.Add(time.Hour), now.Add(-time.Hour), false},
		{"at end boundary", now.Add(-time.Hour), now.Add(time.Hour), now.Add(time.Hour), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, datacop.BetweenTime(tt.start, tt.end)(tt.value))
		})
	}
}
