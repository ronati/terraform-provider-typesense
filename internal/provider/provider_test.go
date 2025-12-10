// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){

	"typesense": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// Check if required environment variables are set
	// If not set, use default values for local testing
	if os.Getenv("TYPESENSE_API_KEY") == "" {
		os.Setenv("TYPESENSE_API_KEY", "test-api-key")
	}
	if os.Getenv("TYPESENSE_API_ADDRESS") == "" {
		os.Setenv("TYPESENSE_API_ADDRESS", "http://localhost:8108")
	}

	// Validate that environment variables are now set
	if v := os.Getenv("TYPESENSE_API_KEY"); v == "" {
		t.Fatal("TYPESENSE_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("TYPESENSE_API_ADDRESS"); v == "" {
		t.Fatal("TYPESENSE_API_ADDRESS must be set for acceptance tests")
	}
}
