package alias

import (
	"regexp"
	"strings"
	"testing"

	"github.com/andoma-go/translit"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/utils/geo"
	"github.com/stretchr/testify/require"
)

func strPtr(v string) *string { return &v }
func intRef(v int) *int       { return &v }

func expectedBase(p *entity.PosterInput) string {
	parts := make([]string, 0, 2)
	if s := slugify(translit.Ru(p.Address)); s != "" {
		parts = append(parts, s)
	}
	if p.District != nil {
		if s := slugify(translit.Ru(*p.District)); s != "" {
			parts = append(parts, s)
		}
	}
	base := strings.Join(parts, "-")
	if base == "" {
		base = "poster"
	}
	return base
}

func TestGenerateAliasSignature(t *testing.T) {
	tests := []struct {
		name   string
		poster *entity.PosterInput
	}{
		{
			name: "address only",
			poster: &entity.PosterInput{
				UserID:        1,
				CategoryAlias: "flat",
				CityID:        3,
				Address:       "ул. Ленина, 10",
				Geo:           geo.GeographyPoint{Lat: 55.75158, Lon: 12.6173},
				Area:          122,
			},
		},
		{
			name: "address and district",
			poster: &entity.PosterInput{
				UserID:        10,
				CategoryAlias: "flat",
				CityID:        30,
				Address:       "пр-т Мира 5",
				District:      strPtr("Центральный"),
				Geo:           geo.GeographyPoint{Lat: 25.7558, Lon: 37.6173},
				Area:          23,
			},
		},
		{
			name: "address and district slugify to empty",
			poster: &entity.PosterInput{
				UserID:        7,
				CategoryAlias: "flat",
				CityID:        9,
				Address:       "!!!",
				District:      strPtr("   "),
				Geo:           geo.GeographyPoint{Lat: 55.7558, Lon: 37.6173},
				Area:          112,
			},
		},
		{
			name: "company id affects hash",
			poster: &entity.PosterInput{
				UserID:        11,
				CategoryAlias: "flat",
				CityID:        13,
				Address:       "Тверская 1",
				CompanyID:     intRef(999),
				Geo:           geo.GeographyPoint{Lat: 55.7558, Lon: 37.6173},
				Area:          212,
			},
		},
		{
			name: "trim and normalize dashes",
			poster: &entity.PosterInput{
				UserID:        21,
				CategoryAlias: "flat",
				CityID:        23,
				Address:       "  ул___Пушкина   15  ",
				Geo:           geo.GeographyPoint{Lat: 55.7558, Lon: 37.6173},
				Area:          122,
			},
		},
		{
			name: "same base different hash source 1",
			poster: &entity.PosterInput{
				UserID:        1,
				CategoryAlias: "flat",
				CityID:        1,
				Geo:           geo.GeographyPoint{Lat: 55.7558, Lon: 37.6173},
				Address:       "Ленина 1",
				Area:          122,
			},
		},
		{
			name: "same base different hash source 2",
			poster: &entity.PosterInput{
				UserID:        1,
				CategoryAlias: "flat",
				CityID:        1,
				Geo:           geo.GeographyPoint{Lat: 55.7558, Lon: 37.6173},
				Address:       "Ленина-1",
				Area:          12,
			},
		},
	}

	re := regexp.MustCompile(`^(.+)-([0-9a-f]{8})$`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateAlias(tt.poster)

			m := re.FindStringSubmatch(got)
			if m == nil {
				t.Fatalf("GenerateAlias() = %q does not match expected pattern <base>-<8hex>", got)
			}
			basePart := m[1]
			hashPart := m[2]

			expBase := expectedBase(tt.poster)
			if basePart != expBase {
				t.Fatalf("base part = %q, want %q", basePart, expBase)
			}
			if len(hashPart) != 8 {
				t.Fatalf("hash part length = %d, want 8", len(hashPart))
			}
		})
	}
}

func TestGenerateAliasSameInput(t *testing.T) {
	poster := &entity.PosterInput{
		UserID:        123,
		CategoryAlias: "flat",
		CityID:        789,
		Address:       "ул. Ленина, 10",
		Geo:           geo.GeographyPoint{Lat: 55.75158, Lon: 12.6173},
		Area:          122,
	}
	res1 := GenerateAlias(poster)
	res2 := GenerateAlias(poster)
	require.Equal(t, res1, res2, "Expected same alias for same input")
}

func TestGenerateAliasDifferentInput(t *testing.T) {
	poster := &entity.PosterInput{
		UserID:        123,
		CategoryAlias: "flat",
		CityID:        789,
		Address:       "ул. Ленина, 10",
		Geo:           geo.GeographyPoint{Lat: 55.75158, Lon: 12.6173},
		Area:          122,
	}
	res1 := GenerateAlias(poster)
	poster.Geo.Lat = 12
	res2 := GenerateAlias(poster)
	require.NotEqual(t, res1, res2, "Expected different aliases for different inputs")
}
