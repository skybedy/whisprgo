package app

import (
	"strings"
	"testing"
)

func TestSanitizeTextMasksAPIKeysAndBearerTokens(t *testing.T) {
	input := `cleanup request failed: Post "https://example.com/v1?key=abc123&x=1": Authorization: Bearer supersecret`
	got := SanitizeText(input)

	if got == input {
		t.Fatalf("expected sanitized output to differ from input")
	}
	if want := `https://example.com/v1?key=***&x=1`; !contains(got, want) {
		t.Fatalf("expected masked query key, got: %s", got)
	}
	if want := `Authorization: Bearer ***`; !contains(got, want) {
		t.Fatalf("expected masked bearer token, got: %s", got)
	}
	if contains(got, "abc123") || contains(got, "supersecret") {
		t.Fatalf("secret leaked in sanitized output: %s", got)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
