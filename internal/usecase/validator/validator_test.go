package validator

import (
	"strings"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/samber/lo"
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

func TestNormolizatePhone(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		inp  string
		out  string
	}{
		{
			name: "With spaces",
			inp:  "+7 (893) 121-12-12",
			out:  "+7 893 121 12 12",
		},
		{
			name: "With 8 prefix",
			inp:  "8 (893) 121-12-12",
			out:  "8 893 121 12 12",
		},
		{
			name: "With equeal",
			inp:  "8 893 121 12 12",
			out:  "8 893 121 12 12",
		},
		{
			name: "With one trim",
			inp:  "+7 893 121 12-12",
			out:  "+7 893 121 12 12",
		},
		{
			name: "With all trim",
			inp:  "+7 893-121-12-12",
			out:  "+7 893 121 12 12",
		},
		{
			name: "Withount spaces",
			inp:  "88931211212",
			out:  "8 893 121 12 12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := NormolizePhoneNumber(tt.inp)
			require.Equal(t, tt.out, result)
		})
	}

}

func TestValidatePhone(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		inp  string
		out  bool
	}{
		{
			name: "With spaces",
			inp:  "+7 (893) 121-12-12",
			out:  true,
		},
		{
			name: "With 8 prefix",
			inp:  "8 (893) 121-12-12",
			out:  true,
		},
		{
			name: "With equeal",
			inp:  "8 893 121 12 12",
			out:  true,
		},
		{
			name: "Withount spaces",
			inp:  "88931211212",
			out:  true,
		},
		{
			name: "With one trim",
			inp:  "+7 893 121 12-12",
			out:  true,
		},
		{
			name: "With all trim",
			inp:  "+7 893-121-12-12",
			out:  true,
		},
		{
			name: "Wrong country",
			inp:  "+10 893-121-12-12",
			out:  false,
		},
		{
			name: "Wrong ()",
			inp:  "+7 (893)-121-12-12",
			out:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ValidatePhone(tt.inp)
			require.Equal(t, tt.out, result)
		})
	}
}

func TestValidateProfile(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		phone     string
		firstname string
		lastname  string
		out       bool
	}{
		{
			name:      "valid all",
			phone:     "+7 (893) 121-12-12",
			firstname: "John",
			lastname:  "Doe",
			out:       false,
		},
		{
			name:      "invalid phone",
			phone:     "invalid",
			firstname: "John",
			lastname:  "Doe",
			out:       true,
		},
		{
			name:      "invalid firstname",
			phone:     "+7 (893) 121-12-12",
			firstname: strings.Repeat("a", MaxNameLenght+1),
			lastname:  "Doe",
			out:       true,
		},
		{
			name:      "invalid lastname",
			phone:     "+7 (893) 121-12-12",
			firstname: "John",
			lastname:  strings.Repeat("a", MaxNameLenght+1),
			out:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ValidateProfile(tt.phone, tt.firstname, tt.lastname)
			if !tt.out {
				require.NoError(t, result)
			} else {
				require.Error(t, result)
			}
		})
	}
}

func TestValidatePhoto(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		fileInput *dto.FileInput
		out       bool
	}{
		{
			name: "valid jpg small",
			fileInput: &dto.FileInput{
				Filename:    "test.jpg",
				Size:        1024,
				ContentType: "image/jpeg",
			},
			out: true,
		},
		{
			name: "valid png",
			fileInput: &dto.FileInput{
				Filename:    "test.png",
				Size:        50000,
				ContentType: "image/png",
			},
			out: true,
		},
		{
			name: "valid svg",
			fileInput: &dto.FileInput{
				Filename:    "test.svg",
				Size:        1000,
				ContentType: "image/svg+xml",
			},
			out: true,
		},
		{
			name:      "invalid nil",
			fileInput: nil,
			out:       false,
		},
		{
			name: "invalid size too small",
			fileInput: &dto.FileInput{
				Filename:    "test.jpg",
				Size:        0,
				ContentType: "image/jpeg",
			},
			out: false,
		},
		{
			name: "invalid size too big",
			fileInput: &dto.FileInput{
				Filename:    "test.jpg",
				Size:        maxPhotoSize + 1,
				ContentType: "image/jpeg",
			},
			out: false,
		},
		{
			name: "invalid wrong ext",
			fileInput: &dto.FileInput{
				Filename:    "test.gif",
				Size:        1024,
				ContentType: "image/jpeg",
			},
			out: false,
		},
		{
			name: "invalid wrong content-type",
			fileInput: &dto.FileInput{
				Filename:    "test.jpg",
				Size:        1024,
				ContentType: "text/plain",
			},
			out: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ValidatePhoto(tt.fileInput)
			require.Equal(t, tt.out, result)
		})
	}
}

func TestValidateAddress(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "valid short",
			in:   "ул. Ленина 1",
			out:  true,
		},
		{
			name: "valid long",
			in:   "пр-т Мира, дом 123, корпус 2, квартира 45",
			out:  true,
		},
		{
			name: "invalid too short",
			in:   strings.Repeat("a", minAddressLength-1),
			out:  false,
		},
		{
			name: "invalid too long",
			in:   strings.Repeat("a", maxAddressLength+1),
			out:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ValidateAddress(tt.in)
			require.Equal(t, tt.out, result)
		})
	}
}

func TestValidateDistrict(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		district *string
		out      bool
	}{
		{
			name:     "valid nil",
			district: nil,
			out:      true,
		},
		{
			name:     "valid short",
			district: lo.ToPtr("Центр"),
			out:      true,
		},
		{
			name:     "valid max length",
			district: lo.ToPtr(strings.Repeat("a", maxDistrictLength)),
			out:      true,
		},
		{
			name:     "invalid too long",
			district: lo.ToPtr(strings.Repeat("a", maxDistrictLength+1)),
			out:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ValidateDistrict(tt.district)
			require.Equal(t, tt.out, result)
		})
	}
}

func TestValidateFeatures(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		features []string
		out      bool
	}{
		{
			name:     "valid empty",
			features: []string{},
			out:      true,
		},
		{
			name:     "valid normal",
			features: []string{"wifi", "parking"},
			out:      true,
		},
		{
			name:     "invalid empty string",
			features: []string{"wifi", ""},
			out:      false,
		},
		{
			name:     "invalid too long",
			features: []string{strings.Repeat("a", 51)},
			out:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ValidateFeatures(tt.features)
			require.Equal(t, tt.out, result)
		})
	}
}

func TestValidatePosterBase(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		poster *dto.PosterInputFlatDTO
		out    bool
	}{
		{
			name: "valid",
			poster: &dto.PosterInputFlatDTO{
				Price:       100000,
				Description: "test",
				Area:        50,
				GeoLat:      55.75,
				GeoLon:      37.61,
				Address:     "ул. Ленина 1",
				FloorCount:  5,
				Features:    []string{"wifi"},
			},
			out: false,
		},
		{
			name: "invalid price zero",
			poster: &dto.PosterInputFlatDTO{
				Price:       0,
				Description: "test",
				Area:        50,
				GeoLat:      55.75,
				GeoLon:      37.61,
				Address:     "ул. Ленина 1",
				FloorCount:  5,
			},
			out: true,
		},
		{
			name: "invalid geo lat",
			poster: &dto.PosterInputFlatDTO{
				Price:       100000,
				Description: "test",
				Area:        50,
				GeoLat:      -190,
				GeoLon:      90.61,
				Address:     "ул. Ленина 1",
				FloorCount:  5,
			},
			out: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := ValidatePosterBase(tt.poster)
			if !tt.out {
				require.NoError(t, result)
			} else {
				require.Error(t, result)
			}
		})
	}
}
