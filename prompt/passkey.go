package prompt

import (
	"fmt"
	"os"
)

// Default configuration for passkey authentication
const (
	DefaultPasskeyRPID   = "aws-vault.local"
	DefaultPasskeyOrigin = "https://aws-vault.local"
)

var (
	// Global WebAuthn client instance for testing
	testWebAuthnClient WebAuthnClient
)

// SetTestWebAuthnClient sets a test client for testing purposes
func SetTestWebAuthnClient(client WebAuthnClient) {
	testWebAuthnClient = client
}

// PasskeyMfaPrompt provides MFA authentication using WebAuthn/FIDO2 passkeys
func PasskeyMfaPrompt(mfaSerial string) (string, error) {
	var client WebAuthnClient

	// Use test client if set (for testing)
	if testWebAuthnClient != nil {
		client = testWebAuthnClient
	} else {
		// In production, create platform-specific client
		client = createPlatformWebAuthnClient()
	}

	if !client.IsAvailable() {
		return "", fmt.Errorf("WebAuthn/passkey authentication not available on this platform")
	}

	// Get configuration from environment or use defaults
	rpID := getEnvOrDefault("AWS_VAULT_PASSKEY_RP_ID", DefaultPasskeyRPID)

	auth := NewPasskeyAuthenticator(client, rpID)
	return auth.GetMFAToken(mfaSerial)
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// createPlatformWebAuthnClient creates a platform-specific WebAuthn client
func createPlatformWebAuthnClient() WebAuthnClient {
	// For now, return a stub client that indicates unavailability
	// Platform-specific implementations will be added later
	return &StubWebAuthnClient{}
}

// StubWebAuthnClient is a stub implementation for platforms without WebAuthn support
type StubWebAuthnClient struct{}

func (s *StubWebAuthnClient) GetAssertion(challenge []byte, rpID string, credentialID []byte) (*Assertion, error) {
	return nil, fmt.Errorf("WebAuthn not implemented for this platform")
}

func (s *StubWebAuthnClient) IsAvailable() bool {
	return false
}

func init() {
	// Register passkey method
	Methods["passkey"] = PasskeyMfaPrompt
}