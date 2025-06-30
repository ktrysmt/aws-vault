package cli

import (
	"github.com/alecthomas/kingpin/v2"

	"github.com/99designs/keyring"
	"github.com/99designs/aws-vault/v7/prompt"
)

// TestWebAuthnClient for integration testing
type TestWebAuthnClient struct {
	assertion *prompt.Assertion
	available bool
}

func (t *TestWebAuthnClient) GetAssertion(challenge []byte, rpID string, credentialID []byte) (*prompt.Assertion, error) {
	return t.assertion, nil
}

func (t *TestWebAuthnClient) IsAvailable() bool {
	return t.available
}

func ExampleExecCommand() {
	app := kingpin.New("aws-vault", "")
	awsVault := ConfigureGlobals(app)
	awsVault.keyringImpl = keyring.NewArrayKeyring([]keyring.Item{
		{Key: "llamas", Data: []byte(`{"AccessKeyID":"ABC","SecretAccessKey":"XYZ"}`)},
	})
	ConfigureExecCommand(app, awsVault)
	kingpin.MustParse(app.Parse([]string{
		"--debug", "exec", "--no-session", "llamas", "--", "sh", "-c", "echo $AWS_ACCESS_KEY_ID",
	}))

	// Output:
	// ABC
}

func ExampleExecCommand_withPasskey() {
	app := kingpin.New("aws-vault", "")
	awsVault := ConfigureGlobals(app)
	awsVault.keyringImpl = keyring.NewArrayKeyring([]keyring.Item{
		{Key: "passkey-profile", Data: []byte(`{"AccessKeyID":"XYZ","SecretAccessKey":"ABC"}`)},
	})
	
	// Set up mock WebAuthn client for testing
	prompt.SetTestWebAuthnClient(&TestWebAuthnClient{
		assertion: &prompt.Assertion{
			Response: prompt.AuthenticatorAssertionResponse{
				Signature: []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0},
			},
		},
		available: true,
	})
	defer prompt.SetTestWebAuthnClient(nil)
	
	ConfigureExecCommand(app, awsVault)
	kingpin.MustParse(app.Parse([]string{
		"--debug", "--prompt", "passkey", "exec", "--no-session", "passkey-profile", "--", "sh", "-c", "echo $AWS_ACCESS_KEY_ID",
	}))

	// Output:
	// XYZ
}
