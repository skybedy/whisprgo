package cleanup

import "testing"

func TestExtractGeminiCleanupText(t *testing.T) {
	raw := `{"candidates":[{"content":{"parts":[{"text":"Ahoj"},{"text":"svete"}]}}]}`
	got, err := extractGeminiCleanupText([]byte(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Ahoj\nsvete" {
		t.Fatalf("unexpected text: got %q", got)
	}
}
