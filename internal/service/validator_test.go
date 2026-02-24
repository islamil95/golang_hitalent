package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateDepartmentName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		wantErr   bool
	}{
		{"ok", "Backend", "Backend", false},
		{"trim", "  Frontend  ", "Frontend", false},
		{"empty", "", "", true},
		{"only spaces", "   ", "", true},
		{"too long", string(make([]byte, 201)), "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateDepartmentName(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, ErrValidation, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestValidateEmployeeFullName(t *testing.T) {
	_, err := ValidateEmployeeFullName("")
	assert.Error(t, err)
	assert.Equal(t, ErrValidation, err)

	s, err := ValidateEmployeeFullName("  Ivan Ivanov  ")
	require.NoError(t, err)
	assert.Equal(t, "Ivan Ivanov", s)
}

func TestValidateEmployeePosition(t *testing.T) {
	_, err := ValidateEmployeePosition("")
	assert.Error(t, err)
	s, err := ValidateEmployeePosition("Developer")
	require.NoError(t, err)
	assert.Equal(t, "Developer", s)
}
