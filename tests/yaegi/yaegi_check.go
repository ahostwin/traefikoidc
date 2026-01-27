package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/unsafe"
	"gopkg.in/yaml.v3"
)

// Manifest represents the plugin manifest structure
type Manifest struct {
	Import  string `yaml:"import"`
	BasePkg string `yaml:"basePkg"`
	Type    string `yaml:"type"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <plugin-directory>\n", os.Args[0])
		os.Exit(1)
	}

	pluginDir := os.Args[1]
	absPluginDir, err := filepath.Abs(pluginDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to get absolute path: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("==========================================")
	fmt.Println("Yaegi Plugin Loading Test")
	fmt.Println("==========================================")
	fmt.Printf("Plugin directory: %s\n", absPluginDir)
	fmt.Println()

	// Verify plugin directory exists
	if _, err := os.Stat(absPluginDir); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "ERROR: Plugin directory does not exist: %s\n", absPluginDir)
		os.Exit(1)
	}

	// Read manifest
	manifestPath := filepath.Join(absPluginDir, ".traefik.yml")
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to open manifest: %v\n", err)
		os.Exit(1)
	}
	defer manifestFile.Close()

	var manifest Manifest
	decoder := yaml.NewDecoder(manifestFile)
	if err := decoder.Decode(&manifest); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to decode manifest: %v\n", err)
		os.Exit(1)
	}

	if manifest.Import == "" {
		fmt.Fprintf(os.Stderr, "ERROR: Manifest missing 'import' field\n")
		os.Exit(1)
	}

	// Determine basePkg (like Traefik does)
	basePkg := manifest.BasePkg
	if basePkg == "" {
		// Extract last part of import path, replacing hyphens with underscores
		parts := strings.Split(manifest.Import, "/")
		basePkg = strings.ReplaceAll(parts[len(parts)-1], "-", "_")
	}

	fmt.Printf("✓ Manifest loaded:\n")
	fmt.Printf("  Import: %s\n", manifest.Import)
	fmt.Printf("  BasePkg: %s\n", basePkg)
	fmt.Println()

	// Read go.mod to get actual module path
	goModPath := filepath.Join(absPluginDir, "go.mod")
	goModData, err := os.ReadFile(goModPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to read go.mod: %v\n", err)
		os.Exit(1)
	}

	// Extract module path from go.mod (first line after "module ")
	modulePath := ""
	lines := strings.Split(string(goModData), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			modulePath = strings.TrimSpace(strings.TrimPrefix(line, "module "))
			break
		}
	}
	if modulePath == "" {
		fmt.Fprintf(os.Stderr, "ERROR: Could not find module path in go.mod\n")
		os.Exit(1)
	}

	fmt.Printf("  Module path (from go.mod): %s\n", modulePath)
	fmt.Println()

	// Set up GoPath structure (like Traefik does for local plugins)
	// Create a temporary GoPath with src directory
	tempGoPath, err := os.MkdirTemp("", "yaegi-test-gopath-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to create temp GoPath: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tempGoPath)

	// Create src directory structure using actual module path
	srcDir := filepath.Join(tempGoPath, "src")
	pluginSrcDir := filepath.Join(srcDir, filepath.FromSlash(modulePath))
	if err := os.MkdirAll(pluginSrcDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to create plugin src directory: %v\n", err)
		os.Exit(1)
	}

	// Also create directory for manifest import path (in case it's different)
	// and symlink it to the actual module path
	if manifest.Import != modulePath {
		manifestSrcDir := filepath.Join(srcDir, filepath.FromSlash(manifest.Import))
		if err := os.MkdirAll(filepath.Dir(manifestSrcDir), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: Failed to create manifest import directory: %v\n", err)
		} else {
			// Remove if exists, then create symlink
			os.Remove(manifestSrcDir)
			if err := os.Symlink(pluginSrcDir, manifestSrcDir); err != nil {
				fmt.Fprintf(os.Stderr, "WARNING: Failed to create symlink from %s to %s: %v\n", manifest.Import, modulePath, err)
				// If symlink fails, just copy the directory
				if err := copyDir(pluginSrcDir, manifestSrcDir); err != nil {
					fmt.Fprintf(os.Stderr, "WARNING: Failed to copy directory: %v\n", err)
				}
			}
		}
	}

	// Copy manifest to GoPath structure
	manifestDest := filepath.Join(pluginSrcDir, ".traefik.yml")
	if err := copyFile(manifestPath, manifestDest); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to copy manifest: %v\n", err)
		os.Exit(1)
	}

	// Copy all .go files from root directory
	goFiles, err := filepath.Glob(filepath.Join(absPluginDir, "*.go"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to find Go files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Copying plugin source files to GoPath structure...")
	for _, goFile := range goFiles {
		fileName := filepath.Base(goFile)
		destPath := filepath.Join(pluginSrcDir, fileName)
		if err := copyFile(goFile, destPath); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to copy %s: %v\n", fileName, err)
			os.Exit(1)
		}
		fmt.Printf("  ✓ Copied: %s\n", fileName)
	}

	// Copy subdirectories (internal, session, etc.) - needed for internal packages
	entries, err := os.ReadDir(absPluginDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to read plugin directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Copying plugin subdirectories...")
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		dirName := entry.Name()
		// Skip vendor, tests, examples, docs, etc.
		if dirName == "vendor" || dirName == "tests" || dirName == "examples" || dirName == "docs" || dirName == "integration" || dirName == "regression" {
			continue
		}
		srcDir := filepath.Join(absPluginDir, dirName)
		destDir := filepath.Join(pluginSrcDir, dirName)
		if err := copyDir(srcDir, destDir); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to copy directory %s: %v\n", dirName, err)
			os.Exit(1)
		}
		fmt.Printf("  ✓ Copied directory: %s\n", dirName)
	}

	// Copy vendor directory if it exists (required for dependencies)
	vendorDir := filepath.Join(absPluginDir, "vendor")
	if _, err := os.Stat(vendorDir); err == nil {
		fmt.Println()
		fmt.Println("Copying vendor dependencies...")
		vendorDest := filepath.Join(tempGoPath, "src")
		if err := copyDir(vendorDir, vendorDest); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to copy vendor directory: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("  ✓ Vendor dependencies copied")
	}
	fmt.Println()

	// Create Yaegi interpreter with GoPath (like Traefik does)
	fmt.Println("Creating Yaegi interpreter...")
	fmt.Printf("  GoPath: %s\n", tempGoPath)
	i := interp.New(interp.Options{
		GoPath: tempGoPath,
		Env:    os.Environ(),
	})
	fmt.Println("✓ Interpreter created")

	// Import standard library
	fmt.Println()
	fmt.Println("Importing standard library...")
	if err := i.Use(stdlib.Symbols); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to import stdlib: %v\n", err)
		os.Exit(1)
	}
	if err := i.Use(unsafe.Symbols); err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Failed to import unsafe: %v\n", err)
	}
	fmt.Println("✓ Standard library imported")

	// Import plugin package (like Traefik does)
	fmt.Println()
	fmt.Println("==========================================")
	fmt.Printf("Importing plugin: %s\n", manifest.Import)
	fmt.Println("==========================================")
	importStmt := fmt.Sprintf(`import "%s"`, manifest.Import)
	_, err = i.Eval(importStmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ ERROR: Failed to import plugin: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nThis might indicate:\n")
		fmt.Fprintf(os.Stderr, "  - Syntax errors in plugin source\n")
		fmt.Fprintf(os.Stderr, "  - Unsupported Go features\n")
		fmt.Fprintf(os.Stderr, "  - Missing dependencies\n")
		fmt.Fprintf(os.Stderr, "  - Incorrect GoPath structure\n")
		os.Exit(1)
	}
	fmt.Println("✓ Plugin imported successfully")
	fmt.Println()

	// Test CreateConfig function (like Traefik does)
	fmt.Println("==========================================")
	fmt.Println("Test 1: CreateConfig() function")
	fmt.Println("==========================================")
	createConfigStmt := fmt.Sprintf(`%s.CreateConfig()`, basePkg)
	fmt.Printf("Executing: %s\n", createConfigStmt)
	_, err = i.Eval(fmt.Sprintf(`config := %s`, createConfigStmt))
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ ERROR: Failed to call CreateConfig(): %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ CreateConfig() called successfully")
	fmt.Println()

	// Test accessing config fields
	fmt.Println("==========================================")
	fmt.Println("Test 2: Config struct field access")
	fmt.Println("==========================================")
	fmt.Println("Testing access to: ProviderURL, ClientID, DynamicClientRegistration")
	testCode := fmt.Sprintf(`
		config := %s.CreateConfig()
		_ = config.ProviderURL
		_ = config.ClientID
		_ = config.DynamicClientRegistration
	`, basePkg)
	_, err = i.Eval(testCode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ ERROR: Failed to access config fields: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("✓ Config struct fields accessible")
	fmt.Println()

	// Test that New function exists (like Traefik does)
	fmt.Println("==========================================")
	fmt.Println("Test 3: New() function resolution")
	fmt.Println("==========================================")
	fmt.Printf("Checking if %s.New can be resolved...\n", basePkg)
	newStmt := fmt.Sprintf(`%s.New`, basePkg)
	_, err = i.Eval(newStmt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ ERROR: Failed to resolve New function: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nTraefik requires a New() function with signature:\n")
		fmt.Fprintf(os.Stderr, "  func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error)\n")
		os.Exit(1)
	}
	fmt.Println("✓ New() function resolved successfully")
	fmt.Println()

	fmt.Println("==========================================")
	fmt.Println("✅ ALL TESTS PASSED!")
	fmt.Println("==========================================")
	fmt.Println("The plugin should load successfully in Traefik.")
	fmt.Println("If Traefik still fails, check:")
	fmt.Println("  1. Plugin volume mount in docker-compose.yaml")
	fmt.Println("  2. Traefik logs for specific errors")
	fmt.Println("  3. Plugin path matches module name in manifest")
	fmt.Println("==========================================")
}

func copyFile(src, dst string) error {
	srcData, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, srcData, 0644)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path from src
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return copyFile(path, destPath)
	})
}
