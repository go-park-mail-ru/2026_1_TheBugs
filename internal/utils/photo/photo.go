package photo

import (
	"fmt"
	"strings"
)

func GetKeyFromPath(path string) string {
	return strings.TrimPrefix(path, "/")
}

func MakeUrlFromPath(path, publicHost, bucket string) string {
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}
	return fmt.Sprintf("%s/%s%s", publicHost, bucket, path)
}
