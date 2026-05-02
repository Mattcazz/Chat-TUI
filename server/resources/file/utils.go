package file

import (
	"path/filepath"
	"regexp"
	"strings"
)

var invalidFileNameCharsRegex = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

func SanitizeFileName(fileName string) string {
	base := filepath.Base(fileName)
	base = strings.TrimSpace(base)
	if base == "" || base == "." || base == string(filepath.Separator) {
		return "file"
	}

	sanitized := invalidFileNameCharsRegex.ReplaceAllString(base, "_")
	sanitized = strings.Trim(sanitized, " .")
	if sanitized == "" {
		return "file"
	}

	return sanitized
}

func GetFileExtension(fileName string) string {
	base := filepath.Base(fileName)
	ext := filepath.Ext(base)

	if ext == "" {
		return ""
	}

	extWithoutDot := strings.TrimPrefix(ext, ".")

	validExt := ""
	for _, char := range extWithoutDot {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' {
			validExt += string(char)
		}
	}

	if len(validExt) > 10 {
		validExt = validExt[:10]
	}

	if validExt == "" {
		return ""
	}

	return "." + strings.ToLower(validExt)
}
