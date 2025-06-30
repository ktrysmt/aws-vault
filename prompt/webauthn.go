package prompt

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

// WebAuthn related types and interfaces

// AuthenticatorAssertionResponse represents a WebAuthn assertion response
type AuthenticatorAssertionResponse struct {
	ClientDataJSON    []byte
	AuthenticatorData []byte
	Signature         []byte
	UserHandle        []byte
}

// Assertion represents a complete WebAuthn assertion
type Assertion struct {
	Response AuthenticatorAssertionResponse
	Type     string
	ID       []byte
}

// WebAuthnClient interface for platform-specific WebAuthn implementations
type WebAuthnClient interface {
	GetAssertion(challenge []byte, rpID string, credentialID []byte) (*Assertion, error)
	IsAvailable() bool
}

// PasskeyAuthenticator handles passkey-based MFA authentication
type PasskeyAuthenticator struct {
	client WebAuthnClient
	rpID   string
}

// NewPasskeyAuthenticator creates a new passkey authenticator
func NewPasskeyAuthenticator(client WebAuthnClient, rpID string) *PasskeyAuthenticator {
	return &PasskeyAuthenticator{
		client: client,
		rpID:   rpID,
	}
}

// GetMFAToken generates an MFA token using passkey authentication
func (p *PasskeyAuthenticator) GetMFAToken(mfaSerial string) (string, error) {
	if !p.client.IsAvailable() {
		return "", fmt.Errorf("WebAuthn client not available")
	}

	challenge := generateChallenge(mfaSerial)
	credentialID := deriveCredentialID(mfaSerial)

	assertion, err := p.client.GetAssertion(challenge, p.rpID, credentialID)
	if err != nil {
		return "", fmt.Errorf("WebAuthn assertion failed: %w", err)
	}

	return convertAssertionToMFAToken(assertion)
}

// generateChallenge creates a deterministic challenge from the MFA serial
func generateChallenge(mfaSerial string) []byte {
	hash := sha256.Sum256([]byte(mfaSerial + ":challenge"))
	return hash[:]
}

// deriveCredentialID creates a credential ID from the MFA serial
func deriveCredentialID(mfaSerial string) []byte {
	hash := sha256.Sum256([]byte(mfaSerial + ":credential"))
	return hash[:16] // Use first 16 bytes as credential ID
}

// convertAssertionToMFAToken converts WebAuthn assertion to 6-digit MFA token
func convertAssertionToMFAToken(assertion *Assertion) (string, error) {
	if len(assertion.Response.Signature) < 4 {
		return "", fmt.Errorf("invalid assertion signature")
	}

	// Use first 4 bytes of signature to generate 6-digit token
	hash := sha256.Sum256(assertion.Response.Signature)
	token := binary.BigEndian.Uint32(hash[:4]) % 1000000
	return fmt.Sprintf("%06d", token), nil
}