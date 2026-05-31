package secrets

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/99designs/keyring"
)

const ServiceName = "whisprgo"

const (
	OpenAIEnvKey        = "OPENAI_API_KEY"
	GeminiEnvKey        = "GEMINI_API_KEY"
	MistralEnvKey       = "MISTRAL_API_KEY"
	TranscriptionEnvKey = "WHISPRGO_TRANSCRIPTION_API_KEY"
	CleanupEnvKey       = "WHISPRGO_CLEANUP_API_KEY"

	OpenAISecretKey        = "provider.openai.api_key"
	GeminiSecretKey        = "provider.gemini.api_key"
	MistralSecretKey       = "provider.mistral.api_key"
	TranscriptionSecretKey = "role.transcription.api_key"
	CleanupSecretKey       = "role.cleanup.api_key"
)

type Source string

const (
	SourceEnv     Source = "environment"
	SourceKeyring Source = "keyring"
)

func providerSpec(provider string) (string, string, error) {
	switch strings.ToLower(strings.TrimSpace(provider)) {
	case "openai":
		return OpenAIEnvKey, OpenAISecretKey, nil
	case "gemini":
		return GeminiEnvKey, GeminiSecretKey, nil
	case "mistral":
		return MistralEnvKey, MistralSecretKey, nil
	default:
		return "", "", fmt.Errorf("unsupported provider: %s", provider)
	}
}

func openKeyring() (keyring.Keyring, error) {
	cfg := keyring.Config{
		ServiceName: ServiceName,
		AllowedBackends: []keyring.BackendType{
			keyring.SecretServiceBackend,
			keyring.KeychainBackend,
			keyring.WinCredBackend,
			keyring.KWalletBackend,
			keyring.PassBackend,
		},
		FileDir:                  os.TempDir(),
		FilePasswordFunc:         nil,
		KeychainTrustApplication: true,
	}
	return keyring.Open(cfg)
}

func roleSpec(role string) (string, string, error) {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "transcription":
		return TranscriptionEnvKey, TranscriptionSecretKey, nil
	case "cleanup":
		return CleanupEnvKey, CleanupSecretKey, nil
	default:
		return "", "", fmt.Errorf("unsupported role: %s", role)
	}
}

func getBySpec(envKey, secretKey string) (string, Source, error) {
	if v := strings.TrimSpace(os.Getenv(envKey)); v != "" {
		return v, SourceEnv, nil
	}
	kr, err := openKeyring()
	if err != nil {
		return "", "", fmt.Errorf("%s is not set in env and keyring is unavailable", envKey)
	}
	item, err := kr.Get(secretKey)
	if err != nil {
		return "", "", fmt.Errorf("%s is not set in env or keyring", envKey)
	}
	v := strings.TrimSpace(string(item.Data))
	if v == "" {
		return "", "", fmt.Errorf("%s is empty in keyring", envKey)
	}
	return v, SourceKeyring, nil
}

func sourceBySpec(envKey, secretKey string) (Source, error) {
	if strings.TrimSpace(os.Getenv(envKey)) != "" {
		return SourceEnv, nil
	}
	kr, err := openKeyring()
	if err != nil {
		return "", fmt.Errorf("keyring unavailable")
	}
	item, err := kr.Get(secretKey)
	if err != nil {
		return "", fmt.Errorf("secret missing")
	}
	if strings.TrimSpace(string(item.Data)) == "" {
		return "", fmt.Errorf("secret empty")
	}
	return SourceKeyring, nil
}

func Get(provider string) (string, Source, error) {
	envKey, secretKey, err := providerSpec(provider)
	if err != nil {
		return "", "", err
	}
	return getBySpec(envKey, secretKey)
}

func SourceFor(provider string) (Source, error) {
	envKey, secretKey, err := providerSpec(provider)
	if err != nil {
		return "", err
	}
	return sourceBySpec(envKey, secretKey)
}

func Set(provider, value string) error {
	_, secretKey, err := providerSpec(provider)
	if err != nil {
		return err
	}
	kr, err := openKeyring()
	if err != nil {
		return errors.New("keyring is unavailable; use environment variable instead")
	}
	return kr.Set(keyring.Item{Key: secretKey, Data: []byte(strings.TrimSpace(value))})
}

func SetRole(role, value string) error {
	_, secretKey, err := roleSpec(role)
	if err != nil {
		return err
	}
	kr, err := openKeyring()
	if err != nil {
		return errors.New("keyring is unavailable; use environment variable instead")
	}
	return kr.Set(keyring.Item{Key: secretKey, Data: []byte(strings.TrimSpace(value))})
}

func Delete(provider string) error {
	_, secretKey, err := providerSpec(provider)
	if err != nil {
		return err
	}
	kr, err := openKeyring()
	if err != nil {
		return errors.New("keyring is unavailable")
	}
	if err := kr.Remove(secretKey); err != nil && !errors.Is(err, keyring.ErrKeyNotFound) {
		return err
	}
	return nil
}

func DeleteRole(role string) error {
	_, secretKey, err := roleSpec(role)
	if err != nil {
		return err
	}
	kr, err := openKeyring()
	if err != nil {
		return errors.New("keyring is unavailable")
	}
	if err := kr.Remove(secretKey); err != nil && !errors.Is(err, keyring.ErrKeyNotFound) {
		return err
	}
	return nil
}

func Status(provider string) string {
	envKey, secretKey, err := providerSpec(provider)
	if err != nil {
		return "not configured"
	}
	if strings.TrimSpace(os.Getenv(envKey)) != "" {
		return fmt.Sprintf("configured via environment variable %s", envKey)
	}
	kr, err := openKeyring()
	if err != nil {
		return "not configured"
	}
	if _, err := kr.Get(secretKey); err == nil {
		return "configured via keyring"
	}
	return "not configured"
}

func StatusForRole(role, providerFallback string) string {
	roleEnv, roleSecret, err := roleSpec(role)
	if err != nil {
		return "not configured"
	}
	if strings.TrimSpace(os.Getenv(roleEnv)) != "" {
		return fmt.Sprintf("configured via environment variable %s", roleEnv)
	}
	kr, err := openKeyring()
	if err == nil {
		if _, err := kr.Get(roleSecret); err == nil {
			return "configured via keyring (role)"
		}
	}
	if providerFallback != "" {
		return fmt.Sprintf("fallback to provider %s: %s", providerFallback, Status(providerFallback))
	}
	return "not configured"
}

func GetForRole(role, providerFallback string) (string, Source, error) {
	roleEnv, roleSecret, err := roleSpec(role)
	if err != nil {
		return "", "", err
	}
	if v, src, err := getBySpec(roleEnv, roleSecret); err == nil {
		return v, src, nil
	}
	if strings.TrimSpace(providerFallback) != "" {
		return Get(providerFallback)
	}
	return "", "", fmt.Errorf("%s is not configured", role)
}
