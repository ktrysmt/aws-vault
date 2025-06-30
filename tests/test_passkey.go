// +build ignore

package main

import (
	"fmt"
	"log"

	"github.com/99designs/aws-vault/v7/prompt"
)

func main() {
	fmt.Println("🔐 AWS Vault Passkey Test")
	fmt.Println("==========================")

	// Check available prompt methods
	methods := prompt.Available()
	fmt.Printf("📋 Available prompt methods: %v\n", methods)

	// Check if passkey is available
	found := false
	for _, method := range methods {
		if method == "passkey" {
			found = true
			break
		}
	}

	if found {
		fmt.Println("✅ Passkey prompt method is available!")

		// Set up test client
		mockClient := &MockWebAuthnClient{
			assertion: &prompt.Assertion{
				Response: prompt.AuthenticatorAssertionResponse{
					Signature: []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0},
				},
			},
			available: true,
		}
		prompt.SetTestWebAuthnClient(mockClient)

		// Test passkey prompt
		token, err := prompt.Method("passkey")("arn:aws:iam::123456789012:mfa/test")
		if err != nil {
			log.Printf("❌ Error: %v", err)
		} else {
			fmt.Printf("🎉 Generated MFA token: %s\n", token)
		}
	} else {
		fmt.Println("❌ Passkey prompt method not found!")
	}
}

// MockWebAuthnClient for testing
type MockWebAuthnClient struct {
	assertion *prompt.Assertion
	available bool
}

func (m *MockWebAuthnClient) GetAssertion(challenge []byte, rpID string, credentialID []byte) (*prompt.Assertion, error) {
	return m.assertion, nil
}

func (m *MockWebAuthnClient) IsAvailable() bool {
	return m.available
}