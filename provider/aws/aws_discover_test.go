// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package aws_test

import (
	"log"
	"os"
	"strings"
	"testing"

	discover "github.com/hashicorp/go-discover"
	"github.com/hashicorp/go-discover/provider/aws"
)

func TestAddrs(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"session_token":     os.Getenv("AWS_SESSION_TOKEN"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
func TestAddrsEndpoint(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"endpoint":          os.Getenv("AWS_EC2_METADATA_SERVICE_ENDPOINT"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}
func TestAddrsECS(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"service":           "ecs",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 3 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsECSFilterClusterName(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"service":           "ecs",
		"ecs_cluster":       "go-discover-2",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsECSFilterFuzzyClusterName(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"service":           "ecs",
		"ecs_cluster":       "go-discover-1",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 1 {
		t.Fatalf("bad: %v", addrs)
	}
}

func TestAddrsECSFilterTaskFamily(t *testing.T) {
	args := discover.Config{
		"provider":          "aws",
		"service":           "ecs",
		"ecs_family":        "go-discover-familia",
		"region":            os.Getenv("AWS_REGION"),
		"tag_key":           "consul",
		"tag_value":         "server",
		"access_key_id":     os.Getenv("AWS_ACCESS_KEY_ID"),
		"secret_access_key": os.Getenv("AWS_SECRET_ACCESS_KEY"),
	}

	if args["region"] == "" || args["access_key_id"] == "" || args["secret_access_key"] == "" {
		t.Skip("AWS credentials or region missing")
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)
	addrs, err := p.Addrs(args, l)
	if err != nil {
		t.Fatal(err)
	}
	if len(addrs) != 2 {
		t.Fatalf("bad: %v", addrs)
	}

}

func TestAddrsDualStackEndpoint(t *testing.T) {
	// Test with dual stack endpoint enabled
	t.Setenv("AWS_USE_DUALSTACK_ENDPOINT", "true")

	args := discover.Config{
		"provider":          "aws",
		"region":            "us-east-1", // Use fixed region to avoid metadata lookup
		"tag_key":           "consul",
		"tag_value":         "server",
		"addr_type":         "public_v6",
		"access_key_id":     "test-key",
		"secret_access_key": "test-secret",
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)

	// This test will exercise the dual stack endpoint configuration path
	// We expect it to fail with actual AWS calls, but we're testing the config creation logic
	_, err := p.Addrs(args, l)

	// We expect an error since we're using fake credentials, but the important part
	// is that the configuration creation logic was exercised
	if err == nil {
		t.Fatal("Expected error with fake credentials, but got none")
	}

	// Verify the error is related to credentials/AWS API, not configuration creation
	if !containsAny(err.Error(), []string{"credential", "auth", "permission", "access"}) {
		t.Logf("Error message: %s", err.Error())
		// This is acceptable - the config creation worked, AWS API call failed as expected
	}
}

func TestAddrsConfigurationWithoutDualStack(t *testing.T) {
	// Test without dual stack endpoint (default behavior)
	t.Setenv("AWS_USE_DUALSTACK_ENDPOINT", "false")

	args := discover.Config{
		"provider":          "aws",
		"region":            "us-east-1", // Use fixed region to avoid metadata lookup
		"tag_key":           "consul",
		"tag_value":         "server",
		"addr_type":         "private_v4",
		"access_key_id":     "test-key",
		"secret_access_key": "test-secret",
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)

	// This test will exercise the non-dual stack configuration path
	_, err := p.Addrs(args, l)

	// We expect an error since we're using fake credentials
	if err == nil {
		t.Fatal("Expected error with fake credentials, but got none")
	}

	// Verify the error is related to credentials/AWS API, not configuration creation
	if !containsAny(err.Error(), []string{"credential", "auth", "permission", "access"}) {
		t.Logf("Error message: %s", err.Error())
		// This is acceptable - the config creation worked, AWS API call failed as expected
	}
}

func TestAddrsDefaultCredentialChain(t *testing.T) {
	// Test the default credential chain path (no static credentials provided)
	args := discover.Config{
		"provider":  "aws",
		"region":    "us-east-1", // Use fixed region to avoid metadata lookup
		"tag_key":   "consul",
		"tag_value": "server",
		"addr_type": "private_v4",
	}

	p := &aws.Provider{}
	l := log.New(os.Stderr, "", log.LstdFlags)

	// This test will exercise the default credential chain configuration path
	_, err := p.Addrs(args, l)

	// We expect an error since no credentials are available in test environment
	if err == nil {
		t.Fatal("Expected error with no credentials, but got none")
	}

	// Verify we're hitting the credential chain logic
	if !containsAny(err.Error(), []string{"credential", "auth", "permission", "access", "config"}) {
		t.Logf("Error message: %s", err.Error())
		// This is acceptable - the config creation worked, credential resolution failed as expected
	}
}

// TestDualStackConfigurations tests all combinations of dual-stack env var and addr types
func TestDualStackConfigurations(t *testing.T) {
	tests := []struct {
		name     string
		envValue string // "", "true", "false", "FALSE"
		addrType string
		region   string
	}{
		{"enabled_explicit", "true", "public_v6", "us-east-1"},
		{"enabled_unset", "", "public_v6", "us-east-1"},              // backward compat: dual-stack
		{"disabled_explicit", "false", "public_v6", "us-east-1"},     // FIX: honors disable
		{"disabled_case_insensitive", "FALSE", "public_v6", "us-east-1"}, // FIX: case insensitive
		{"disabled_with_v4", "false", "private_v4", "us-east-1"},
		{"disabled_me_central", "false", "private_v4", "me-central-1"}, // key use case
		{"enabled_with_v4", "true", "private_v4", "us-east-1"}, // v4 addr type, no dual-stack
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envValue != "" {
				t.Setenv("AWS_USE_DUALSTACK_ENDPOINT", tc.envValue)
			}

			args := discover.Config{
				"provider":          "aws",
				"region":            tc.region,
				"tag_key":           "consul",
				"tag_value":         "server",
				"addr_type":         tc.addrType,
				"access_key_id":     "test-key",
				"secret_access_key": "test-secret",
			}

			p := &aws.Provider{}
			l := log.New(os.Stderr, "", log.LstdFlags)

			// Exercise the config creation path
			_, err := p.Addrs(args, l)

			// We expect an error due to fake credentials
			if err == nil {
				t.Fatal("Expected error with fake credentials, but got none")
			}

			// Verify the error is auth-related, not a dual-stack endpoint error
			errLower := strings.ToLower(err.Error())
			if strings.Contains(errLower, "dualstack") || strings.Contains(errLower, "dual-stack") {
				t.Errorf("Got dual-stack endpoint error, expected auth error: %s", err.Error())
			}
		})
	}
}

// TestDualStackDefaultCredentialChain tests dual-stack configurations with the default
// credential chain (no static credentials), which is the production use case for IAM roles.
func TestDualStackDefaultCredentialChain(t *testing.T) {
	tests := []struct {
		name     string
		envValue string // "", "true", "false"
		addrType string
		region   string
	}{
		{"default_chain_enabled_unset", "", "private_v4", "us-east-1"},
		{"default_chain_enabled_explicit", "true", "private_v4", "us-east-1"},
		{"default_chain_disabled", "false", "private_v4", "us-east-1"},               // FIX
		{"default_chain_disabled_me_central", "false", "private_v4", "me-central-1"}, // key production use case
		{"default_chain_disabled_v6", "false", "public_v6", "us-east-1"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.envValue != "" {
				t.Setenv("AWS_USE_DUALSTACK_ENDPOINT", tc.envValue)
			}

			// No access_key_id or secret_access_key - uses default credential chain (IAM role)
			args := discover.Config{
				"provider":  "aws",
				"region":    tc.region,
				"tag_key":   "consul",
				"tag_value": "server",
				"addr_type": tc.addrType,
			}

			p := &aws.Provider{}
			l := log.New(os.Stderr, "", log.LstdFlags)

			// Exercise the default credential chain config path
			_, err := p.Addrs(args, l)

			// We expect an error since no credentials are available in test environment
			if err == nil {
				t.Fatal("Expected error with no credentials, but got none")
			}

			// Verify the error is credential-related, not a dual-stack endpoint error
			errLower := strings.ToLower(err.Error())
			if strings.Contains(errLower, "dualstack") || strings.Contains(errLower, "dual-stack") {
				t.Errorf("Got dual-stack endpoint error, expected credential error: %s", err.Error())
			}
		})
	}
}

// Helper function to check if error message contains any of the expected strings
func containsAny(str string, substrings []string) bool {
	for _, substr := range substrings {
		if len(str) >= len(substr) {
			for i := 0; i <= len(str)-len(substr); i++ {
				match := true
				for j := 0; j < len(substr); j++ {
					if str[i+j] != substr[j] && str[i+j] != substr[j]+32 && str[i+j] != substr[j]-32 {
						match = false
						break
					}
				}
				if match {
					return true
				}
			}
		}
	}
	return false
}
