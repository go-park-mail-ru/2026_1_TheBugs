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
		{
			name: "invalid too long",
			in:   "" + string(make([]byte, 255)) + "@example.com",
			out:  false,
		},
		{
			name: "invalid with spaces",
			in:   "test @example.com",
			out:  false,
		},
		{
			name: "invalid with special chars",
			in:   "😊test@example.com",
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
		{
			name: "invalid too long",
			in:   "A" + string(make([]byte, 64)) + "1a",
			out:  false,
		},
		{
			name: "invalid no special char",
			in:   "😊Password😊😊😊😊123df",
			out:  false,
		},
		{
			name: "invalid no special only char",
			in:   "😊😊😊😊😊😊😊😊😊😊😊😊😊😊😊",
			out:  false,
		},
		{
			name: "invalid with space",
			in:   "Password 123",
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

type cred struct {
	email string
	pwd   string
}

func TestValidateCred(t *testing.T) {
	tests := []struct {
		name string
		in   cred
		out  bool
	}{
		{
			name: "valid strong password",
			in: cred{
				email: "test@example.com",
				pwd:   "Passw0rd123",
			},
			out: true,
		},
		{
			name: "valid minimum length",
			in: cred{
				email: "test@example.com",
				pwd:   "Ab1defgH",
			},
			out: true,
		},
		{
			name: "valid with special chars",
			in: cred{
				email: "test@example.com",
				pwd:   "P@ssw0rd1",
			},
			out: true,
		},
		{
			name: "invalid too short",
			in: cred{
				email: "test@example.com",
				pwd:   "Ab1d",
			},
			out: false,
		},
		{
			name: "invalid no uppercase",
			in: cred{
				email: "test@example.com",
				pwd:   "password123",
			},
			out: false,
		},
		{
			name: "invalid no lowercase",
			in: cred{
				email: "test@example.com",
				pwd:   "PASSWORD123",
			},
			out: false,
		},
		{
			name: "invalid no digit",
			in: cred{
				email: "test@example.com",
				pwd:   "PasswordAB",
			},
			out: false,
		},
		{
			name: "invalid empty pwd",
			in: cred{
				email: "test@example.com",
				pwd:   "",
			},
			out: false,
		},
		{
			name: "invalid empty",
			in: cred{
				email: "",
				pwd:   "",
			},
			out: false,
		},
		{
			name: "invalid too long",
			in: cred{
				email: "test@example.com",
				pwd:   "A" + string(make([]byte, 64)) + "1a",
			},
			out: false,
		},
		{
			name: "invalid no special char",
			in: cred{
				email: "test@example.com",
				pwd:   "😊Password😊😊😊😊123df",
			},
			out: false,
		},
		{
			name: "invalid with space",
			in: cred{
				email: "test@example.com",
				pwd:   "Password  123df",
			},
			out: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ValidateCred(tt.in.email, tt.in.pwd)
			if tt.out {
				require.NoError(t, result)
				return
			}
			require.Error(t, result)
		})
	}
}
