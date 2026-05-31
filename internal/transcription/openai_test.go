package transcription

import "testing"

func TestExtractTranscriptionText(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{name: "text", raw: `{"text":"ahoj"}`, want: "ahoj"},
		{name: "output_text", raw: `{"output_text":"ahoj svete"}`, want: "ahoj svete"},
		{name: "output content", raw: `{"output":[{"content":[{"type":"output_text","text":"line1"},{"type":"output_text","text":"line2"}]}]}`, want: "line1\nline2"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := extractTranscriptionText([]byte(tc.raw))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("unexpected text: got %q want %q", got, tc.want)
			}
		})
	}
}
