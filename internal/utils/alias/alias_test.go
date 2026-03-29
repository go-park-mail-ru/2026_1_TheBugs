package alias

import (
	"testing"

	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
)

func strPtr(v string) *string { return &v }
func intRef(v int) *int       { return &v }

func TestGenerateAlias(t *testing.T) {
	tests := []struct {
		name   string
		poster *entity.PosterInput
		want   string
	}{
		{
			name: "address only",
			poster: &entity.PosterInput{
				UserID:     1,
				CategoryID: 2,
				CityID:     3,
				Address:    "ул. Ленина, 10",
			},
			want: "ul-lenina-10-" + shortHash(
				"user:1|category:2|city:3|address:ул. Ленина, 10|district:|company:",
				8,
			),
		},
		{
			name: "address and district",
			poster: &entity.PosterInput{
				UserID:     10,
				CategoryID: 20,
				CityID:     30,
				Address:    "пр-т Мира 5",
				District:   strPtr("Центральный"),
			},
			want: "pr-t-mira-5-centralnyj-" + shortHash(
				"user:10|category:20|city:30|address:пр-т Мира 5|district:Центральный|company:",
				8,
			),
		},
		{
			name: "address and district slugify to empty",
			poster: &entity.PosterInput{
				UserID:     7,
				CategoryID: 8,
				CityID:     9,
				Address:    "!!!",
				District:   strPtr("   "),
			},
			want: "poster-" + shortHash(
				"user:7|category:8|city:9|address:!!!|district:   |company:",
				8,
			),
		},
		{
			name: "company id affects hash",
			poster: &entity.PosterInput{
				UserID:     11,
				CategoryID: 12,
				CityID:     13,
				Address:    "Тверская 1",
				CompanyID:  intRef(999),
			},
			want: "tverskaja-1-" + shortHash(
				"user:11|category:12|city:13|address:Тверская 1|district:|company:999",
				8,
			),
		},
		{
			name: "trim and normalize dashes",
			poster: &entity.PosterInput{
				UserID:     21,
				CategoryID: 22,
				CityID:     23,
				Address:    "  ул___Пушкина   15  ",
			},
			want: "ul-pushkina-15-" + shortHash(
				"user:21|category:22|city:23|address:  ул___Пушкина   15  |district:|company:",
				8,
			),
		},
		{
			name: "same base but different hash source",
			poster: &entity.PosterInput{
				UserID:     1,
				CategoryID: 1,
				CityID:     1,
				Address:    "Ленина 1",
			},
			want: "lenina-1-" + shortHash(
				"user:1|category:1|city:1|address:Ленина 1|district:|company:",
				8,
			),
		},
		{
			name: "same slug different original address changes hash",
			poster: &entity.PosterInput{
				UserID:     1,
				CategoryID: 1,
				CityID:     1,
				Address:    "Ленина-1",
			},
			want: "lenina-1-" + shortHash(
				"user:1|category:1|city:1|address:Ленина-1|district:|company:",
				8,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			got := GenerateAlias(test.poster)
			if got != test.want {
				t.Fatalf("GenerateAlias() = %q, want %q", got, test.want)
			}
		})
	}
}
