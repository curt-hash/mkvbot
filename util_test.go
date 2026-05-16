package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeFileName(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"Normal Title", "Normal Title"},
		{"Colon: To Dash", "Colon- To Dash"},
		{"Remove/Slashes\\Both", "RemoveSlashesBoth"},
		{`Strip<>"|*?`, "Strip"},
		{"Keep apostrophe's", "Keep apostrophe's"},
		{"Non-print \x01\x02 stripped", "Non-print  stripped"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, tt.want, sanitizeFileName(tt.in))
		})
	}
}
