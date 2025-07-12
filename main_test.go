package main

import (
	"testing"
)

func TestBump(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		component string
		want      string
		wantErr   bool
	}{
		// Standard semantic versioning
		{"bump patch", "1.2.3", "patch", "1.2.4", false},
		{"bump minor", "1.2.3", "minor", "1.3.0", false},
		{"bump major", "1.2.3", "major", "2.0.0", false},

		// Pre-release version (should treat as "1.2.3")
		{"prerelease", "1.2.3-alpha", "patch", "", true},

		// Short version schemas
		{"major only", "5", "major", "6", false},
		{"major.minor", "2.7", "minor", "2.8", false},
		{"major.minor bump major", "2.7", "major", "3.0", false},

		// Too many version elements
		// {"four elements", "1.2.3.4", "patch", "", true},

		// Invalid component
		{"invalid component", "1.2.3", "build", "", true},

		// Non-numeric version
		{"non-numeric", "a.b.c", "patch", "", true},

		// Component index out of bounds
		{"component out of bounds", "1", "minor", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bump(tt.version, tt.component)
			if (err != nil) != tt.wantErr {
				t.Errorf("bump() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("bump() = %v, want %v", got, tt.want)
			}
		})
	}
}