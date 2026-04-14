package validator

import (
	"html"
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/usecase/dto"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func TestSanitizePosterInput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *dto.PosterInputFlatDTO
		expected *dto.PosterInputFlatDTO
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name: "all text fields with XSS",
			input: &dto.PosterInputFlatDTO{
				Description:   "<script>alert('xss')</script>Hello World",
				CategoryAlias: "Category<script>",
				Address:       "123 Main St <img src=x onerror=alert(1)>",
				City:          "Moscow<script>",
				District:      lo.ToPtr("Downtown<script>"),
				Alias:         lo.ToPtr("my-alias<script>"),
				Features:      []string{"parking<script>", "pool", "<b>gym</b>"},
			},
			expected: &dto.PosterInputFlatDTO{
				Description:   html.EscapeString("<script>alert('xss')</script>Hello World"),
				CategoryAlias: html.EscapeString("Category<script>"),
				Address:       html.EscapeString("123 Main St <img src=x onerror=alert(1)>"),
				City:          html.EscapeString("Moscow<script>"),
				District:      lo.ToPtr(html.EscapeString("Downtown<script>")),
				Alias:         lo.ToPtr(html.EscapeString("my-alias<script>")),
				Features:      []string{html.EscapeString("parking<script>"), html.EscapeString("pool"), html.EscapeString("<b>gym</b>")},
			},
		},
		{
			name: "empty strings",
			input: &dto.PosterInputFlatDTO{
				Description:   "",
				CategoryAlias: "",
				Address:       "",
				City:          "",
				District:      lo.ToPtr(""),
				Alias:         lo.ToPtr(""),
				Features:      []string{},
			},
			expected: &dto.PosterInputFlatDTO{
				Description:   "",
				CategoryAlias: "",
				Address:       "",
				City:          "",
				District:      lo.ToPtr(""),
				Alias:         lo.ToPtr(""),
				Features:      []string{},
			},
		},
		{
			name: "special characters",
			input: &dto.PosterInputFlatDTO{
				Description:   "Hello & World < > \" '",
				CategoryAlias: "Category & Co",
				Address:       "Street & Avenue",
				City:          "New York & London",
				District:      lo.ToPtr("District & Region"),
				Alias:         lo.ToPtr("alias&symbol"),
				Features:      []string{"feature&1", "feature<2>", "feature\"3\""},
			},
			expected: &dto.PosterInputFlatDTO{
				Description:   html.EscapeString("Hello & World < > \" '"),
				CategoryAlias: html.EscapeString("Category & Co"),
				Address:       html.EscapeString("Street & Avenue"),
				City:          html.EscapeString("New York & London"),
				District:      lo.ToPtr(html.EscapeString("District & Region")),
				Alias:         lo.ToPtr(html.EscapeString("alias&symbol")),
				Features:      []string{html.EscapeString("feature&1"), html.EscapeString("feature<2>"), html.EscapeString("feature\"3\"")},
			},
		},
		{
			name: "nil pointers",
			input: &dto.PosterInputFlatDTO{
				Description:   "Test",
				CategoryAlias: "Test",
				Address:       "Test",
				City:          "Test",
				District:      nil,
				Alias:         nil,
				Features:      []string{"test"},
			},
			expected: &dto.PosterInputFlatDTO{
				Description:   "Test",
				CategoryAlias: "Test",
				Address:       "Test",
				City:          "Test",
				District:      nil,
				Alias:         nil,
				Features:      []string{"test"},
			},
		},
		{
			name: "already escaped string",
			input: &dto.PosterInputFlatDTO{
				Description:   "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
				CategoryAlias: "Category",
				Address:       "Address",
				City:          "City",
				Features:      []string{"feature"},
			},
			expected: &dto.PosterInputFlatDTO{
				Description:   html.EscapeString("&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;"),
				CategoryAlias: "Category",
				Address:       "Address",
				City:          "City",
				Features:      []string{"feature"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			SanitizePosterInput(tt.input)

			if tt.input == nil {
				require.Nil(t, tt.input)
				return
			}

			require.Equal(t, tt.expected.Description, tt.input.Description)
			require.Equal(t, tt.expected.CategoryAlias, tt.input.CategoryAlias)
			require.Equal(t, tt.expected.Address, tt.input.Address)
			require.Equal(t, tt.expected.City, tt.input.City)

			if tt.expected.District != nil {
				require.NotNil(t, tt.input.District)
				require.Equal(t, *tt.expected.District, *tt.input.District)
			} else {
				require.Nil(t, tt.input.District)
			}

			if tt.expected.Alias != nil {
				require.NotNil(t, tt.input.Alias)
				require.Equal(t, *tt.expected.Alias, *tt.input.Alias)
			} else {
				require.Nil(t, tt.input.Alias)
			}

			require.Equal(t, tt.expected.Features, tt.input.Features)
		})
	}
}

func TestSanitizeUserProfile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    *dto.CreateUserDTO
		expected *dto.CreateUserDTO
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name: "all fields with XSS",
			input: &dto.CreateUserDTO{
				Email:     "<script>alert('xss')</script>user@example.com",
				FirstName: "<b>John</b>",
				LastName:  "<i>Doe</i>",
				Phone:     "+7<script>(123)456-78-90",
			},
			expected: &dto.CreateUserDTO{
				Email:     html.EscapeString("<script>alert('xss')</script>user@example.com"),
				FirstName: html.EscapeString("<b>John</b>"),
				LastName:  html.EscapeString("<i>Doe</i>"),
				Phone:     html.EscapeString("+7<script>(123)456-78-90"),
			},
		},
		{
			name: "special characters",
			input: &dto.CreateUserDTO{
				Email:     "user&name@example.com",
				FirstName: "John & Doe",
				LastName:  "O'Reilly",
				Phone:     "+7 (123) 456-78-90",
			},
			expected: &dto.CreateUserDTO{
				Email:     html.EscapeString("user&name@example.com"),
				FirstName: html.EscapeString("John & Doe"),
				LastName:  html.EscapeString("O'Reilly"),
				Phone:     html.EscapeString("+7 (123) 456-78-90"),
			},
		},
		{
			name: "empty strings",
			input: &dto.CreateUserDTO{
				Email:     "",
				FirstName: "",
				LastName:  "",
				Phone:     "",
			},
			expected: &dto.CreateUserDTO{
				Email:     "",
				FirstName: "",
				LastName:  "",
				Phone:     "",
			},
		},
		{
			name: "already escaped",
			input: &dto.CreateUserDTO{
				Email:     "&lt;script&gt;alert(1)&lt;/script&gt;",
				FirstName: "&lt;b&gt;John&lt;/b&gt;",
				LastName:  "Doe",
				Phone:     "+71234567890",
			},
			expected: &dto.CreateUserDTO{
				Email:     html.EscapeString("&lt;script&gt;alert(1)&lt;/script&gt;"),
				FirstName: html.EscapeString("&lt;b&gt;John&lt;/b&gt;"),
				LastName:  "Doe",
				Phone:     "+71234567890",
			},
		},
		{
			name: "unicode and emoji",
			input: &dto.CreateUserDTO{
				Email:     "user😊@example.com",
				FirstName: "Иван",
				LastName:  "Петров",
				Phone:     "+7 (123) 456-78-90",
			},
			expected: &dto.CreateUserDTO{
				Email:     "user😊@example.com",
				FirstName: "Иван",
				LastName:  "Петров",
				Phone:     "+7 (123) 456-78-90",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			SanitizeUserProfile(tt.input)

			if tt.input == nil {
				require.Nil(t, tt.input)
				return
			}

			require.Equal(t, tt.expected.Email, tt.input.Email)
			require.Equal(t, tt.expected.FirstName, tt.input.FirstName)
			require.Equal(t, tt.expected.LastName, tt.input.LastName)
			require.Equal(t, tt.expected.Phone, tt.input.Phone)
		})
	}
}

func TestSanitizePosterInput_Security(t *testing.T) {
	t.Parallel()

	maliciousInputs := []struct {
		name   string
		field  string
		attack string
	}{
		{"XSS script", "Description", "<script>alert('xss')</script>"},
		{"XSS img onerror", "Address", "<img src=x onerror=alert(1)>"},
		{"XSS svg", "City", "<svg onload=alert(1)>"},
		{"XSS iframe", "CategoryAlias", "<iframe src='javascript:alert(1)'>"},
		{"XSS body", "Feature", "<body onload=alert(1)>"},
		{"XSS div", "District", "<div onmouseover='alert(1)'>"},
	}

	for _, malicious := range maliciousInputs {
		t.Run(malicious.name, func(t *testing.T) {
			t.Parallel()

			dto := &dto.PosterInputFlatDTO{
				Description:   "",
				CategoryAlias: "",
				Address:       "",
				City:          "",
				District:      lo.ToPtr(""),
				Alias:         lo.ToPtr(""),
				Features:      []string{""},
			}

			switch malicious.field {
			case "Description":
				dto.Description = malicious.attack
			case "CategoryAlias":
				dto.CategoryAlias = malicious.attack
			case "Address":
				dto.Address = malicious.attack
			case "City":
				dto.City = malicious.attack
			case "District":
				*dto.District = malicious.attack
			case "Feature":
				dto.Features[0] = malicious.attack
			}

			SanitizePosterInput(dto)

			switch malicious.field {
			case "Description":
				require.NotContains(t, dto.Description, "<script>")
				require.NotContains(t, dto.Description, "<img")
				require.Contains(t, dto.Description, "&lt;script&gt;")
			case "CategoryAlias":
				require.NotContains(t, dto.CategoryAlias, "<iframe")
				require.Contains(t, dto.CategoryAlias, "&lt;iframe")
			case "Address":
				require.NotContains(t, dto.Address, "<img")
				require.Contains(t, dto.Address, "&lt;img")
			case "City":
				require.NotContains(t, dto.City, "<svg")
				require.Contains(t, dto.City, "&lt;svg")
			case "District":
				require.NotContains(t, *dto.District, "<div")
				require.Contains(t, *dto.District, "&lt;div")
			}
		})
	}
}
