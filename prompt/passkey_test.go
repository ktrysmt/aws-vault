package prompt_test

import (
	"testing"

	"github.com/99designs/aws-vault/v7/prompt"
)

// MockWebAuthnClient for testing
type MockWebAuthnClient struct {
	assertion *prompt.Assertion
	err       error
	available bool
}

func (m *MockWebAuthnClient) GetAssertion(challenge []byte, rpID string, credentialID []byte) (*prompt.Assertion, error) {
	return m.assertion, m.err
}

func (m *MockWebAuthnClient) IsAvailable() bool {
	return m.available
}

func TestPasskeyPrompt_ShouldBeAvailable(t *testing.T) {
	methods := prompt.Available()

	// First failing test - passkey should be available
	found := false
	for _, method := range methods {
		if method == "passkey" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected passkey to be available in prompt methods")
	}
}

func TestPasskeyPrompt_ShouldReturnToken_WithMockClient(t *testing.T) {
	// Setup mock client
	mockClient := &MockWebAuthnClient{
		assertion: &prompt.Assertion{
			Response: prompt.AuthenticatorAssertionResponse{
				Signature: []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0},
			},
		},
		available: true,
	}
	
	prompt.SetTestWebAuthnClient(mockClient)
	defer prompt.SetTestWebAuthnClient(nil) // cleanup

	token, err := prompt.Method("passkey")("arn:aws:iam::123456789012:mfa/user")

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(token) != 6 {
		t.Errorf("Expected 6-digit token, got %s", token)
	}
}

func TestPasskeyPrompt_ShouldReturnError_WhenNotAvailable(t *testing.T) {
	// Reset test client to use stub
	prompt.SetTestWebAuthnClient(nil)

	_, err := prompt.Method("passkey")("arn:aws:iam::123456789012:mfa/user")

	if err == nil {
		t.Errorf("Expected error when WebAuthn not available")
	}
}