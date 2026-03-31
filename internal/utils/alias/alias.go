package alias

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"

	"github.com/andoma-go/translit"
	"github.com/go-park-mail-ru/2026_1_TheBugs/internal/entity"
)

func GenerateAlias(poster *entity.PosterInput) string {
	parts := make([]string, 0, 2)

	if s := slugify(translit.Ru(poster.Address)); s != "" {
		parts = append(parts, s)
	}

	if poster.District != nil {
		if s := slugify(translit.Ru(*poster.District)); s != "" {
			parts = append(parts, s)
		}
	}

	base := strings.Join(parts, "-")
	if base == "" {
		base = "poster"
	}

	hashSource := fmt.Sprintf(
		"user:%d|category:%s|city:%d|address:%s|district:%s|company:%s|geo:%f,%f|area:%f",
		poster.UserID,
		poster.CategoryAlias,
		poster.CityID,
		poster.Address,
		stringPtr(poster.District),
		intPtr(poster.CompanyID),
		poster.Geo.Lat,
		poster.Geo.Lon,
		poster.Area,
	)
	// key := make([]byte, 8)
	// rand.Read(key)
	// hashSource += string(key)

	hash := shortHash(hashSource, 8)

	return base + "-" + hash
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))

	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")

	reDash := regexp.MustCompile(`-+`)
	s = reDash.ReplaceAllString(s, "-")

	return s
}

func shortHash(s string, n int) string {
	sum := sha1.Sum([]byte(s))
	hexStr := hex.EncodeToString(sum[:])

	if n > len(hexStr) {
		n = len(hexStr)
	}

	return hexStr[:n]
}

func stringPtr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func intPtr(v *int) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%d", *v)
}
