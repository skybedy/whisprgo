package cleanup

import "testing"

func TestExtractDeepSeekCleanupText(t *testing.T) {
	raw := `{"choices":[{"message":{"content":"Ahoj svete"}}]}`
	got, err := extractDeepSeekCleanupText([]byte(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "Ahoj svete" {
		t.Fatalf("unexpected text: got %q", got)
	}
}
