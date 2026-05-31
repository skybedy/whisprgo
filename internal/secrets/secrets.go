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
	OpenAIEnvKey = "OPENAI_API_KEY"
	GeminiEnvKey = "GEMINI_API_KEY"

	OpenAISecretKey = "provider.openai.api_key"
	GeminiSecretKey = "provider.gemini.api_key"
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

func Get(provider string) (string, Source, error) {
	envKey, secretKey, err := providerSpec(provider)
	if err != nil {
		return "", "", err
	}
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
