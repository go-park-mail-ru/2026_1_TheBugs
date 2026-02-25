package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "valid email",
			in:   "test@example.com",
			out:  true,
		},
		{
			name: "valid with numbers",
			in:   "user123@domain.co.uk",
			out:  true,
		},
		{
			name: "valid with special chars",
			in:   "user.name+tag@gmail.com",
			out:  true,
		},
		{
			name: "invalid no at",
			in:   "testexample.com",
			out:  false,
		},
		{
			name: "invalid short tld",
			in:   "test@example.c",
			out:  false,
		},
		{
			name: "invalid empty",
			in:   "",
			out:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ValidateEmail(tt.in)
			require.Equal(t, tt.out, result)
		})
	}
}

func TestValidatePwd(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "valid strong password",
			in:   "Passw0rd123",
			out:  true,
		},
		{
			name: "valid minimum length",
			in:   "Ab1defgH",
			out:  true,
		},
		{
			name: "valid with special chars",
			in:   "P@ssw0rd1",
			out:  true,
		},
		{
			name: "invalid too short",
			in:   "Ab1d",
			out:  false,
		},
		{
			name: "invalid no uppercase",
			in:   "password123",
			out:  false,
		},
		{
			name: "invalid no lowercase",
			in:   "PASSWORD123",
			out:  false,
		},
		{
			name: "invalid no digit",
			in:   "PasswordAB",
			out:  false,
		},
		{
			name: "invalid empty",
			in:   "",
			out:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ValidatePwd(tt.in)
			require.Equal(t, tt.out, result)
		})
	}
}
