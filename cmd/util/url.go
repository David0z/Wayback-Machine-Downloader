package util

import (
	"regexp"
	"strings"
)

func RemoveSlashFromString(url string) string {
	return strings.ReplaceAll(url, "/", "")
}

func SanitizeFileName(fileName string) string {
	invalidChars := `[<>:"/\\|?*]`
	re := regexp.MustCompile(invalidChars)
	return re.ReplaceAllString(fileName, "_")
}
