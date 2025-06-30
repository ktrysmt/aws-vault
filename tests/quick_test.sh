#!/bin/bash

# Move to project root directory
cd "$(dirname "$0")/.."

echo "🚀 Quick AWS Vault Passkey Test"
echo "================================"
echo "📁 Working directory: $(pwd)"
echo ""

# Test 1: Check if passkey is in help
echo "📖 Testing help output..."
if go run main.go --help 2>&1 | grep -q "passkey"; then
    echo "✅ Passkey found in help"
else
    echo "❌ Passkey not found in help"
    exit 1
fi

# Test 2: Run unit tests
echo "🧪 Running unit tests..."
if go test ./prompt -run TestPasskeyPrompt -v > /dev/null 2>&1; then
    echo "✅ Passkey unit tests pass"
else
    echo "❌ Passkey unit tests failed"
    exit 1
fi

# Test 3: Test passkey functionality
echo "🔐 Testing passkey functionality..."
if go run tests/test_passkey.go 2>/dev/null | grep -q "Generated MFA token"; then
    echo "✅ Passkey functionality works"
else
    echo "❌ Passkey functionality failed"
    exit 1
fi

# Test 4: Check config parsing
echo "⚙️  Testing configuration..."
if go test ./vault -run TestPasskeyConfigurationParsing -v > /dev/null 2>&1; then
    echo "✅ Passkey configuration parsing works"
else
    echo "❌ Passkey configuration parsing failed"
    exit 1
fi

echo ""
echo "🎉 All tests passed! Passkey integration is working correctly."
echo "💡 Use 'go run main.go --prompt passkey <command>' for testing"