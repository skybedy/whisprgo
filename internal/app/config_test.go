package app

import "testing"

func TestIsSecretKeyPath(t *testing.T) {
	cases := []string{"openai_api_key", "api_key", "provider.openai.api_key", "secrets.openai"}
	for _, c := range cases {
		if !isSecretKeyPath(c) {
			t.Fatalf("expected secret key path to be blocked: %s", c)
		}
	}
	if isSecretKeyPath("transcription.model") {
		t.Fatalf("non-secret key should not be blocked")
	}
}
