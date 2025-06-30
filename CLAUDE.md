# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

AWS Vault is a Go CLI tool that securely stores AWS credentials in your operating system's keystore and generates temporary credentials for shell environments and applications. It supports multiple authentication flows including assume role, MFA (including passkey/WebAuthn), AWS SSO, and web identity federation.

## Development Commands

### Build and Test
```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
go test ./...

# Run tests with race detection
go test -race ./...

# Run linting (requires golangci-lint v1.52.0)
golangci-lint run

# Clean build artifacts
make clean
```

### Local Development
```bash
# Install development version
go install

# Run specific test file
go test ./vault -run TestSpecificFunction

# Run CLI commands during development (faster than building)
go run main.go list
go run main.go exec --prompt=passkey profile-name -- env | grep AWS

# Run E2E tests
./tests/quick_test.sh

# Test passkey functionality interactively
go run tests/test_passkey.go

# Generate cross-platform binaries
make build-all
```

## Architecture Overview

### Core Components

**CLI Layer (`cli/`)**
- Each command has its own file with `Configure*Command()` functions
- Uses Kingpin v2 for CLI framework
- Main commands: `add`, `exec`, `list`, `login`, `export`, `clear`, `remove`, `rotate`, `proxy`
- `exec` command is the primary functionality for executing commands with temporary credentials

**Vault Core (`vault/`)**
- **Credential Providers**: Multiple provider implementations for different AWS credential sources
  - `sessiontokenprovider.go` - STS GetSessionToken (with MFA support)
  - `assumeroleprovider.go` - STS AssumeRole for cross-account access
  - `ssorolecredentialsprovider.go` - AWS SSO integration
  - `assumerolewithwebidentityprovider.go` - Web identity federation
  - `credentialprocessprovider.go` - External credential processes
- **Storage**: `credentialkeyring.go` for long-term credentials, `sessionkeyring.go` for temporary sessions
- **Caching**: `cachedsessionprovider.go` wraps providers with session caching
- **Configuration**: `config.go` handles AWS config file parsing and profile management (includes passkey config)

**Prompt System (`prompt/`)**
- **MFA Prompts**: Multiple authentication methods for MFA tokens
  - `passkey.go` - WebAuthn/FIDO2 passkey authentication
  - `webauthn.go` - WebAuthn protocol implementation and interfaces
  - Platform-specific prompt implementations for traditional MFA
- **Testing Support**: Mock clients and dependency injection for unit testing

**Server Components (`server/`)**
- Local metadata servers that mimic EC2/ECS metadata endpoints
- Allows applications to automatically refresh credentials without environment variables
- Platform-specific proxy implementations for network interception

### Authentication Flow Architecture

The credential resolution follows this hierarchy:
1. **Master Credentials** → Stored in OS keystore via `credentialkeyring.go`
2. **Session Tokens** → Generated via STS, cached in `sessionkeyring.go`
3. **Role Assumption** → Uses session tokens to assume roles
4. **Alternative Sources** → SSO, web identity, or credential processes

### Keystore Integration

Cross-platform credential storage using `github.com/99designs/keyring`:
- **macOS**: Keychain (with code signing for security)
- **Windows**: Credential Manager
- **Linux**: Secret Service (GNOME Keyring), KWallet, Pass
- **Fallback**: Encrypted file storage

### Configuration Management

AWS configuration is parsed from standard `~/.aws/config` files with support for:
- Profile inheritance and source profiles
- Role chaining (multi-level role assumptions)
- MFA serial numbers and token processes
- Custom credential processes
- Environment variable overrides
- **Passkey configuration**: `passkey_rp_id`, `passkey_origin`, `passkey_credential_id`, `passkey_user_id`

### Passkey Configuration Example
```ini
[profile example-with-passkey]
region = us-east-1
mfa_serial = arn:aws:iam::123456789012:mfa/user
passkey_rp_id = aws.amazon.com
passkey_origin = https://signin.aws.amazon.com
passkey_credential_id = base64-encoded-credential-id
passkey_user_id = user-handle
```

## Key Design Patterns

**Provider Pattern**: All credential sources implement `aws.CredentialsProvider` interface, allowing pluggable authentication methods.

**Decorator Pattern**: `cachedsessionprovider.go` wraps other providers to add caching functionality without modifying their core logic.

**Factory Pattern**: `vault.go` contains provider factories that create appropriate credential providers based on profile configuration.

**Command Pattern**: CLI commands are self-contained units with their own configuration and execution logic.

## Important Implementation Details

### MFA Handling
- MFA prompts are handled by `prompt/` package with platform-specific implementations
- **Traditional Methods**: Hardware token support (YubiKey) and software tokens
- **Passkey Support**: WebAuthn/FIDO2 passkey authentication following W3C WebAuthn specification
- **Prompt Method Selection**: Use `--prompt=passkey` CLI flag or configure in profile
- Session tokens are cached to minimize MFA prompts
- **Testing**: Mock WebAuthn clients available for unit testing without hardware

### Session Caching Strategy
- Temporary credentials cached in keystore with expiration tracking
- 5-minute expiration buffer to avoid credential expiry during execution
- Cache keys based on profile configuration hash for isolation

### Cross-Platform Builds
The Makefile supports building for multiple platforms with proper binary naming and code signing for macOS. Distribution includes DMG creation and checksums.

### Error Handling Patterns
- Extensive use of Go's error wrapping for context
- AWS SDK v2 error types for specific AWS error handling
- Graceful degradation when keystore unavailable

## Testing Strategy

### Unit Testing
- Unit tests alongside implementation files (`*_test.go`)
- **Classic TDD Approach**: Following Kent Beck's Red-Green-Refactor methodology
- **Mock Implementations**: `MockWebAuthnClient` for passkey testing without hardware
- **Dependency Injection**: Testable design allowing mock client substitution

### Integration Testing
- Integration tests in CLI commands testing end-to-end flows
- **E2E Tests**: Located in `tests/` directory for better organization
  - `tests/quick_test.sh` - Comprehensive automated testing
  - `tests/test_passkey.go` - Interactive passkey functionality testing
- Race condition testing with `-race` flag
- Cross-platform CI testing on Ubuntu and macOS

### Testing Commands
```bash
# Run all unit tests
go test ./...

# Test specific functionality
go test ./prompt -run TestPasskey

# Run with race detection
go test -race ./...

# E2E testing suite
./tests/quick_test.sh

# Interactive passkey testing
go run tests/test_passkey.go
```

## Security Considerations

- Credential storage only in OS-native secure keystores
- Temporary credentials with shortest possible lifespan
- No credential logging or debugging output
- Code signing for macOS distribution
- Principle of least privilege in role configurations

### Passkey Security
- **FIDO2/WebAuthn Compliance**: Implements W3C WebAuthn Level 3 specification
- **Hardware Security**: Passkeys stored in secure hardware (TPM, Secure Enclave)
- **Phishing Resistance**: Origin verification prevents phishing attacks
- **Replay Protection**: Challenge-response mechanism prevents replay attacks
- **No Shared Secrets**: Private keys never leave the authenticator device

### Implementation Status
- **Current Version**: Development implementation with mock WebAuthn client
- **Production Readiness**: Requires platform-specific WebAuthn client implementation
- **Security Testing**: Comprehensive unit tests with mock authentication flows