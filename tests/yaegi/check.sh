#!/bin/bash
# Yaegi Plugin Loading Test
# This script uses Yaegi to test if the traefikoidc plugin can be loaded correctly
# It simulates what Traefik does when loading a plugin

# Don't exit on error - we want to handle errors gracefully

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PLUGIN_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"
TEST_DIR="$SCRIPT_DIR"

echo "=========================================="
echo "Yaegi Plugin Loading Test"
echo "=========================================="
echo "Plugin directory: $PLUGIN_DIR"
echo "Test directory: $TEST_DIR"
echo ""

# Check if yaegi is installed
if ! command -v yaegi &> /dev/null; then
    echo "ERROR: yaegi is not installed"
    echo "Install it with: go install github.com/traefik/yaegi/cmd/yaegi@latest"
    exit 1
fi

echo "Yaegi found: $(which yaegi)"
echo ""

# Check if we're in the right directory
if [ ! -f "$PLUGIN_DIR/go.mod" ]; then
    echo "ERROR: go.mod not found in $PLUGIN_DIR"
    exit 1
fi

echo "Found go.mod:"
cat "$PLUGIN_DIR/go.mod" | head -3
echo ""

# Build the test program
echo "Building test program..."
echo "Working directory: $TEST_DIR"
cd "$TEST_DIR"

# Check if go.mod exists, if not create it
if [ ! -f go.mod ]; then
    echo "WARNING: go.mod not found in test directory, creating it..."
    go mod init github.com/lukaszraczylo/traefikoidc/tests/yaegi 2>&1
fi

# Get dependencies
echo "Getting dependencies..."
go get github.com/traefik/yaegi/interp 2>&1 | grep -v "^go: " || true
go get github.com/traefik/yaegi/stdlib 2>&1 | grep -v "^go: " || true
go mod tidy 2>&1 | grep -v "^go: " || true

echo "Building yaegi_check..."
if ! go build -v -o yaegi_check yaegi_check.go 2>&1; then
    echo "ERROR: Failed to build test program"
    echo "Attempting to run directly with 'go run'..."
    echo "----------------------------------------"
    go run yaegi_check.go "$PLUGIN_DIR" 2>&1
    EXIT_CODE=$?
else
    echo "✓ Test program built successfully"
    echo ""
    
    # Run the test
    echo "Running Yaegi plugin loading test..."
    echo "----------------------------------------"
    ./yaegi_check "$PLUGIN_DIR" 2>&1
    EXIT_CODE=$?
fi

echo "----------------------------------------"
if [ $EXIT_CODE -eq 0 ]; then
    echo "✓ Plugin loads successfully with Yaegi!"
    echo ""
    echo "If Traefik still can't load the plugin, check:"
    echo "  1. Plugin volume mount in docker-compose.yaml"
    echo "  2. Traefik logs for specific errors"
    echo "  3. Plugin path matches module name in go.mod"
else
    echo "✗ Plugin failed to load with Yaegi"
    echo ""
    echo "This indicates a problem with the plugin code that prevents Yaegi from interpreting it."
    echo "Common issues:"
    echo "  - Unsupported Go features"
    echo "  - Syntax errors"
    echo "  - Missing dependencies"
    echo "  - Import path issues"
fi

# Cleanup
rm -f yaegi_check

exit $EXIT_CODE
