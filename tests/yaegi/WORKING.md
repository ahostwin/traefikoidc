./ref/traefikoidc/tests/yaegi/check.sh
==========================================
Yaegi Plugin Loading Test
==========================================
Plugin directory: ./ref/traefikoidc
Test directory: ./ref/traefikoidc/tests/yaegi

Yaegi found: ~/.gvm/pkgsets/go1.25.4/global/bin/yaegi

Found go.mod:
module github.com/lukaszraczylo/traefikoidc

go 1.24.0

Building test program...
Working directory: ./ref/traefikoidc/tests/yaegi
Getting dependencies...
Building yaegi_check...
✓ Test program built successfully

Running Yaegi plugin loading test...
----------------------------------------
==========================================
Yaegi Plugin Loading Test
==========================================
Plugin directory: ./ref/traefikoidc

✓ Manifest loaded:
  Import: github.com/lukaszraczylo/traefikoidc
  BasePkg: traefikoidc

  Module path (from go.mod): github.com/lukaszraczylo/traefikoidc

Copying plugin source files to GoPath structure...
  ✓ Copied: audience_test.go
  ✓ Copied: auth_flow.go
  ✓ Copied: auth_flow_behaviour_test.go
  ✓ Copied: auth_flow_pkce_test.go
  ✓ Copied: autocleanup.go
  ✓ Copied: autocleanup_additional_test.go
  ✓ Copied: azure_oidc_test.go
  ✓ Copied: background_tasks_ultra_test.go
  ✓ Copied: cache_bench_test.go
  ✓ Copied: cache_compat.go
  ✓ Copied: cache_manager.go
  ✓ Copied: cache_test.go
  ✓ Copied: config_marshalling.go
  ✓ Copied: coverage_boost_final_test.go
  ✓ Copied: csrf_session_test.go
  ✓ Copied: custom_claims_test.go
  ✓ Copied: dcr_storage_compat.go
  ✓ Copied: dcr_storage_test.go
  ✓ Copied: dynamic_client_registration.go
  ✓ Copied: dynamic_client_registration_test.go
  ✓ Copied: edge_cases_suite_test.go
  ✓ Copied: enhanced_mocks_suite_test.go
  ✓ Copied: enhanced_mocks_test.go
  ✓ Copied: error_recovery.go
  ✓ Copied: error_recovery_bench_test.go
  ✓ Copied: error_recovery_test.go
  ✓ Copied: goroutine_manager.go
  ✓ Copied: goroutine_manager_test.go
  ✓ Copied: helpers.go
  ✓ Copied: http_client_factory.go
  ✓ Copied: http_client_factory_unit_test.go
  ✓ Copied: http_client_pool.go
  ✓ Copied: http_client_pool_test.go
  ✓ Copied: input_validation.go
  ✓ Copied: input_validation_test.go
  ✓ Copied: issue67_regression_test.go
  ✓ Copied: jwk.go
  ✓ Copied: jwk_caching_test.go
  ✓ Copied: jwt.go
  ✓ Copied: logger_singleton.go
  ✓ Copied: logout.go
  ✓ Copied: logout_test.go
  ✓ Copied: main.go
  ✓ Copied: main_bench_test.go
  ✓ Copied: main_coverage_boost2_test.go
  ✓ Copied: main_coverage_boost_test.go
  ✓ Copied: main_exchange_test.go
  ✓ Copied: main_goroutine_leak_test.go
  ✓ Copied: main_initialization_test.go
  ✓ Copied: main_refresh_test.go
  ✓ Copied: main_servehttp_test.go
  ✓ Copied: main_simple_test.go
  ✓ Copied: main_test.go
  ✓ Copied: memory_leak_bench_test.go
  ✓ Copied: memory_leak_fixes.go
  ✓ Copied: memory_leak_test.go
  ✓ Copied: memory_monitor.go
  ✓ Copied: memory_optimizations.go
  ✓ Copied: metadata_cache.go
  ✓ Copied: middleware.go
  ✓ Copied: middleware_edge_cases_test.go
  ✓ Copied: mocks_test.go
  ✓ Copied: pkce_flow_test.go
  ✓ Copied: profiling.go
  ✓ Copied: profiling_test.go
  ✓ Copied: redis_integration_test.go
  ✓ Copied: refresh_coordinator.go
  ✓ Copied: refresh_coordinator_test.go
  ✓ Copied: refresh_race_test.go
  ✓ Copied: scope_filter.go
  ✓ Copied: scope_filter_test.go
  ✓ Copied: security_edge_cases_test.go
  ✓ Copied: security_monitoring.go
  ✓ Copied: security_monitoring_test.go
  ✓ Copied: session.go
  ✓ Copied: session_behaviour_test.go
  ✓ Copied: session_bench_test.go
  ✓ Copied: session_chunk_cleanup.go
  ✓ Copied: session_chunk_manager.go
  ✓ Copied: session_test.go
  ✓ Copied: settings.go
  ✓ Copied: sharded_cache.go
  ✓ Copied: singleton_resources.go
  ✓ Copied: singleton_resources_test.go
  ✓ Copied: test_config.go
  ✓ Copied: test_framework_test.go
  ✓ Copied: test_helpers_adapter_test.go
  ✓ Copied: test_infrastructure.go
  ✓ Copied: test_main_test.go
  ✓ Copied: test_utils_test.go
  ✓ Copied: testify_mocks_test.go
  ✓ Copied: testutil_example_test.go
  ✓ Copied: token_bench_test.go
  ✓ Copied: token_introspection.go
  ✓ Copied: token_manager.go
  ✓ Copied: token_resilience.go
  ✓ Copied: token_test.go
  ✓ Copied: token_validation_suite_test.go
  ✓ Copied: token_validator.go
  ✓ Copied: types.go
  ✓ Copied: universal_cache.go
  ✓ Copied: universal_cache_serialization_test.go
  ✓ Copied: universal_cache_singleton.go
  ✓ Copied: url_helpers.go
  ✓ Copied: url_helpers_ultra_test.go
  ✓ Copied: utilities.go

Copying plugin subdirectories...
  ✓ Copied directory: .git
  ✓ Copied directory: .github
  ✓ Copied directory: internal
  ✓ Copied directory: session

Copying vendor dependencies...
  ✓ Vendor dependencies copied

Creating Yaegi interpreter...
  GoPath: /tmp/yaegi-test-gopath-3062936617
✓ Interpreter created

Importing standard library...
✓ Standard library imported

==========================================
Importing plugin: github.com/ahostwin/traefikoidc
==========================================
✓ Plugin imported successfully

==========================================
Test 1: CreateConfig() function
==========================================
Executing: traefikoidc.CreateConfig()
✓ CreateConfig() called successfully

==========================================
Test 2: Config struct field access
==========================================
Testing access to: ProviderURL, ClientID, DynamicClientRegistration
✓ Config struct fields accessible

==========================================
Test 3: New() function resolution
==========================================
Checking if traefikoidc.New can be resolved...
✓ New() function resolved successfully

==========================================
✅ ALL TESTS PASSED!
==========================================
The plugin should load successfully in Traefik.
If Traefik still fails, check:
  1. Plugin volume mount in docker-compose.yaml
  2. Traefik logs for specific errors
  3. Plugin path matches module name in manifest
==========================================
----------------------------------------
✓ Plugin loads successfully with Yaegi!

If Traefik still can't load the plugin, check:
  1. Plugin volume mount in docker-compose.yaml
  2. Traefik logs for specific errors
  3. Plugin path matches module name in go.mod
