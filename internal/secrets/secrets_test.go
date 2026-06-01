package secrets

import "testing"

func TestProviderSpecSupportsMultipleProviders(t *testing.T) {
	cases := []struct {
		provider   string
		wantEnv    string
		wantSecret string
	}{
		{provider: "openai", wantEnv: OpenAIEnvKey, wantSecret: OpenAISecretKey},
		{provider: "gemini", wantEnv: GeminiEnvKey, wantSecret: GeminiSecretKey},
		{provider: "mistral", wantEnv: MistralEnvKey, wantSecret: MistralSecretKey},
		{provider: "deepseek", wantEnv: DeepSeekEnvKey, wantSecret: DeepSeekSecretKey},
	}

	for _, tc := range cases {
		envKey, secretKey, err := providerSpec(tc.provider)
		if err != nil {
			t.Fatalf("provider %s: unexpected error: %v", tc.provider, err)
		}
		if envKey != tc.wantEnv {
			t.Fatalf("provider %s: unexpected env key got %q want %q", tc.provider, envKey, tc.wantEnv)
		}
		if secretKey != tc.wantSecret {
			t.Fatalf("provider %s: unexpected secret key got %q want %q", tc.provider, secretKey, tc.wantSecret)
		}
	}
}

func TestRoleSpecSupportsTranscriptionAndCleanup(t *testing.T) {
	cases := []struct {
		role       string
		wantEnv    string
		wantSecret string
	}{
		{role: "transcription", wantEnv: TranscriptionEnvKey, wantSecret: TranscriptionSecretKey},
		{role: "cleanup", wantEnv: CleanupEnvKey, wantSecret: CleanupSecretKey},
	}

	for _, tc := range cases {
		envKey, secretKey, err := roleSpec(tc.role)
		if err != nil {
			t.Fatalf("role %s: unexpected error: %v", tc.role, err)
		}
		if envKey != tc.wantEnv {
			t.Fatalf("role %s: unexpected env key got %q want %q", tc.role, envKey, tc.wantEnv)
		}
		if secretKey != tc.wantSecret {
			t.Fatalf("role %s: unexpected secret key got %q want %q", tc.role, secretKey, tc.wantSecret)
		}
	}
}
