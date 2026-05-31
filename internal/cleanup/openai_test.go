package cleanup

import "testing"

func TestExtractCleanupText(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{name: "output_text", raw: `{"output_text":"upraveny text"}`, want: "upraveny text"},
		{name: "text", raw: `{"text":"upraveny text"}`, want: "upraveny text"},
		{name: "output content", raw: `{"output":[{"content":[{"type":"output_text","text":"radek 1"},{"type":"output_text","text":"radek 2"}]}]}`, want: "radek 1\nradek 2"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := extractCleanupText([]byte(tc.raw))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("unexpected text: got %q want %q", got, tc.want)
			}
		})
	}
}
