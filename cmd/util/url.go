package util

import "strings"

func RemoveSlashFromString(url string) string {
	return strings.ReplaceAll(url, "/", "")
}
