package transcription

import "testing"

func TestParseParakeetWSResponse(t *testing.T) {
	cases := []struct {
		name    string
		raw     string
		want    string
		wantErr bool
		done    bool
	}{
		{name: "partial", raw: `{"type":"partial","text":""}`, done: false},
		{name: "result", raw: `{"type":"result","text":"ahoj"}`, want: "ahoj", done: true},
		{name: "error", raw: `{"error":"boom"}`, wantErr: true},
		{name: "empty result", raw: `{"type":"result","text":""}`, wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, done, err := parseParakeetWSResponse([]byte(tc.raw))
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("unexpected text got %q want %q", got, tc.want)
			}
			if done != tc.done {
				t.Fatalf("unexpected done got %v want %v", done, tc.done)
			}
		})
	}
}
