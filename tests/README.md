# E2E Tests for AWS Vault Passkey Integration

This directory contains end-to-end tests and integration tests for the passkey functionality.

## Files

### `quick_test.sh`
Comprehensive test suite that validates all passkey functionality:
- ✅ Passkey prompt method registration
- ✅ Unit test execution
- ✅ Functional testing with mock WebAuthn client
- ✅ Configuration parsing

**Usage:**
```bash
# From project root
./tests/quick_test.sh

# Or from tests directory
cd tests && ./quick_test.sh
```

### `test_passkey.go`
Standalone test program for interactive testing of passkey functionality.
Demonstrates:
- Prompt method availability check
- Mock WebAuthn client setup
- MFA token generation

**Usage:**
```bash
# From project root
go run tests/test_passkey.go

# Or from tests directory
cd tests && go run test_passkey.go
```

## Development Workflow

### Quick Development Test
```bash
# Fast feedback loop during development
./tests/quick_test.sh
```

### Manual Testing
```bash
# Interactive testing
go run tests/test_passkey.go

# CLI help verification
go run main.go --help | grep passkey
```

### Continuous Testing
```bash
# Watch for file changes and auto-test
find . -name "*.go" | entr -c ./tests/quick_test.sh
```

## Test Coverage

The tests cover:

1. **Registration Testing**: Verify passkey is registered as a prompt method
2. **Unit Testing**: Execute all passkey-related unit tests
3. **Integration Testing**: Test WebAuthn client interface with mocks
4. **Configuration Testing**: Validate passkey settings parsing
5. **CLI Integration**: Confirm passkey appears in help and CLI options

## Notes

- All tests use mock WebAuthn clients, so no actual hardware authentication is required
- Tests are designed to run quickly (< 5 seconds) for rapid development feedback
- The `+build ignore` tag in `test_passkey.go` prevents it from being included in regular builds