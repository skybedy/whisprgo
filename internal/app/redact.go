package app

import "regexp"

var (
	queryKeyPattern      = regexp.MustCompile(`(?i)([?&]key=)([^&\s"]+)`)
	bearerTokenPattern   = regexp.MustCompile(`(?i)(bearer\s+)([^\s",]+)`)
	authorizationPattern = regexp.MustCompile(`(?i)(authorization\s*[:=]\s*bearer\s+)([^\s",]+)`)
)

// SanitizeText masks secrets before they reach logs or stderr.
func SanitizeText(input string) string {
	out := queryKeyPattern.ReplaceAllString(input, `${1}***`)
	out = authorizationPattern.ReplaceAllString(out, `${1}***`)
	out = bearerTokenPattern.ReplaceAllString(out, `${1}***`)
	return out
}
