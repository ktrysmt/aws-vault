package prompt_test

import (
	"errors"
	"testing"

	"github.com/99designs/aws-vault/v7/prompt"
)


func TestPasskeyAuthenticator_GetMFAToken_Success(t *testing.T) {
	// Given
	mockClient := &MockWebAuthnClient{
		assertion: &prompt.Assertion{
			Response: prompt.AuthenticatorAssertionResponse{
				Signature: []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0},
			},
		},
		available: true,
	}
	auth := prompt.NewPasskeyAuthenticator(mockClient, "aws-vault.local")

	// When
	token, err := auth.GetMFAToken("arn:aws:iam::123456789012:mfa/user")

	// Then
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(token) != 6 {
		t.Errorf("Expected 6-digit token, got %s", token)
	}
	// Token should be numeric
	for _, c := range token {
		if c < '0' || c > '9' {
			t.Errorf("Expected numeric token, got %s", token)
		}
	}
}

func TestPasskeyAuthenticator_GetMFAToken_ClientNotAvailable(t *testing.T) {
	// Given
	mockClient := &MockWebAuthnClient{
		available: false,
	}
	auth := prompt.NewPasskeyAuthenticator(mockClient, "aws-vault.local")

	// When
	_, err := auth.GetMFAToken("arn:aws:iam::123456789012:mfa/user")

	// Then
	if err == nil {
		t.Errorf("Expected error when client not available")
	}
}

func TestPasskeyAuthenticator_GetMFAToken_AssertionError(t *testing.T) {
	// Given
	mockClient := &MockWebAuthnClient{
		err:       errors.New("authenticator not responding"),
		available: true,
	}
	auth := prompt.NewPasskeyAuthenticator(mockClient, "aws-vault.local")

	// When
	_, err := auth.GetMFAToken("arn:aws:iam::123456789012:mfa/user")

	// Then
	if err == nil {
		t.Errorf("Expected error when assertion fails")
	}
}

func TestConvertAssertionToMFAToken_InvalidSignature(t *testing.T) {
	// Test with exported function would require refactoring
	// For now, test through the main interface
	mockClient := &MockWebAuthnClient{
		assertion: &prompt.Assertion{
			Response: prompt.AuthenticatorAssertionResponse{
				Signature: []byte{0x12}, // Too short
			},
		},
		available: true,
	}
	auth := prompt.NewPasskeyAuthenticator(mockClient, "aws-vault.local")

	_, err := auth.GetMFAToken("arn:aws:iam::123456789012:mfa/user")

	if err == nil {
		t.Errorf("Expected error for invalid signature")
	}
}

func TestPasskeyAuthenticator_DeterministicChallenge(t *testing.T) {
	// Test that same MFA serial produces same challenge
	mockClient := &MockWebAuthnClient{
		assertion: &prompt.Assertion{
			Response: prompt.AuthenticatorAssertionResponse{
				Signature: []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0},
			},
		},
		available: true,
	}
	auth := prompt.NewPasskeyAuthenticator(mockClient, "aws-vault.local")

	mfaSerial := "arn:aws:iam::123456789012:mfa/user"

	// Get token twice
	token1, err1 := auth.GetMFAToken(mfaSerial)
	token2, err2 := auth.GetMFAToken(mfaSerial)

	if err1 != nil || err2 != nil {
		t.Errorf("Expected no errors, got %v, %v", err1, err2)
	}

	if token1 != token2 {
		t.Errorf("Expected deterministic tokens for same MFA serial, got %s != %s", token1, token2)
	}
}