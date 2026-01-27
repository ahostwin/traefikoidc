# Yaegi Plugin Loading Test

This directory contains tests to verify that the traefikoidc plugin can be loaded by Yaegi (the Go interpreter used by Traefik).

## Purpose

Traefik uses Yaegi to interpret and load plugins at runtime. If Yaegi can't load the plugin, Traefik will fail with errors like:
- "invalid middleware type or middleware does not exist"
- Plugin loading errors in Traefik logs

This test simulates what Traefik does when loading the plugin, helping to identify issues before deploying to Traefik.

## Prerequisites

1. **Go 1.24+** installed
2. **Yaegi** installed:
   ```bash
   go install github.com/traefik/yaegi/cmd/yaegi@latest
   ```

## Running the Test

From the plugin root directory (`ref/traefikoidc/`):

```bash
cd tests/yaegi
./check.sh
```

Or manually:

```bash
cd tests/yaegi
go mod tidy
go build -o yaegi_check yaegi_check.go
./yaegi_check ../..
```

## What It Tests

1. **Package Import**: Verifies Yaegi can import the `traefikoidc` package
2. **CreateConfig()**: Tests if the `CreateConfig()` function can be called
3. **Struct Access**: Verifies config struct fields are accessible
4. **ClientRegistrationMetadata**: Tests access to nested DCR metadata struct
5. **Field Creation**: Tests creating and accessing `ClientRegistrationMetadata` instances

## Common Issues

### "Failed to import package"
- Check that `go.mod` module name matches the import path
- Verify all dependencies are available
- Check for unsupported Go features in the code

### "Failed to access config fields"
- May indicate struct tag issues
- Could be a Yaegi limitation with certain struct patterns

### "Failed to create/access ClientRegistrationMetadata"
- Often related to struct tag problems (json/yaml tags)
- May indicate issues with nested struct definitions

## Troubleshooting

If the test passes but Traefik still can't load the plugin:

1. **Check Traefik logs** for specific errors:
   ```bash
   docker logs traefik 2>&1 | grep -i "plugin\|traefikoidc\|error"
   ```

2. **Verify plugin volume mount** in docker-compose.yaml:
   ```yaml
   - ../ref/traefikoidc:/plugins-local/src/github.com/lukaszraczylo/traefikoidc:ro
   ```

3. **Check plugin path** matches module name in go.mod

4. **Verify Traefik version** supports the plugin features used

## Files

- `check.sh` - Main test script
- `yaegi_check.go` - Go program that uses Yaegi to test plugin loading
- `go.mod` - Dependencies for the test program
- `README.md` - This file
