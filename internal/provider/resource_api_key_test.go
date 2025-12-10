package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApiKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApiKeyResourceConfig("test-key", "Test API Key", "*", "*"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_api_key.test", "description", "Test API Key"),
					resource.TestCheckResourceAttr("typesense_api_key.test", "actions.#", "1"),
					resource.TestCheckResourceAttr("typesense_api_key.test", "actions.0", "*"),
					resource.TestCheckResourceAttr("typesense_api_key.test", "collections.#", "1"),
					resource.TestCheckResourceAttr("typesense_api_key.test", "collections.0", "*"),
					resource.TestCheckResourceAttrSet("typesense_api_key.test", "id"),
					resource.TestCheckResourceAttrSet("typesense_api_key.test", "value"),
					resource.TestCheckResourceAttrSet("typesense_api_key.test", "value_prefix"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "typesense_api_key.test",
				ImportState:       true,
				ImportStateVerify: true,
				// value is only available on creation, not on read
				ImportStateVerifyIgnore: []string{"value"},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccApiKeyResource_WithSpecificPermissions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with specific permissions
			{
				Config: testAccApiKeyResourceConfigSpecific("search-key", "Search Only Key", "documents:search", "products"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_api_key.test", "description", "Search Only Key"),
					resource.TestCheckResourceAttr("typesense_api_key.test", "actions.#", "1"),
					resource.TestCheckResourceAttr("typesense_api_key.test", "actions.0", "documents:search"),
					resource.TestCheckResourceAttr("typesense_api_key.test", "collections.#", "1"),
					resource.TestCheckResourceAttr("typesense_api_key.test", "collections.0", "products"),
					resource.TestCheckResourceAttrSet("typesense_api_key.test", "id"),
					resource.TestCheckResourceAttrSet("typesense_api_key.test", "value"),
				),
			},
		},
	})
}

func TestAccApiKeyResource_WithExpiration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with expiration
			{
				Config: testAccApiKeyResourceConfigWithExpiration("expiring-key", "Expiring Key", "1735689600"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_api_key.test", "description", "Expiring Key"),
					resource.TestCheckResourceAttr("typesense_api_key.test", "expires_at", "1735689600"),
					resource.TestCheckResourceAttrSet("typesense_api_key.test", "id"),
					resource.TestCheckResourceAttrSet("typesense_api_key.test", "value"),
				),
			},
		},
	})
}

func TestAccApiKeyResource_MultipleActions(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create with multiple actions
			{
				Config: testAccApiKeyResourceConfigMultipleActions("multi-action-key", "Multi Action Key"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_api_key.test", "description", "Multi Action Key"),
					resource.TestCheckResourceAttr("typesense_api_key.test", "actions.#", "2"),
					resource.TestCheckResourceAttrSet("typesense_api_key.test", "id"),
					resource.TestCheckResourceAttrSet("typesense_api_key.test", "value"),
				),
			},
		},
	})
}

func testAccApiKeyResourceConfig(name, description, action, collection string) string {
	return fmt.Sprintf(`
resource "typesense_api_key" "test" {
  description = %[2]q
  actions     = [%[3]q]
  collections = [%[4]q]
}
`, name, description, action, collection)
}

func testAccApiKeyResourceConfigSpecific(name, description, action, collection string) string {
	return fmt.Sprintf(`
resource "typesense_api_key" "test" {
  description = %[2]q
  actions     = [%[3]q]
  collections = [%[4]q]
}
`, name, description, action, collection)
}

func testAccApiKeyResourceConfigWithExpiration(name, description, expiresAt string) string {
	return fmt.Sprintf(`
resource "typesense_api_key" "test" {
  description = %[2]q
  actions     = ["*"]
  collections = ["*"]
  expires_at  = %[3]s
}
`, name, description, expiresAt)
}

func testAccApiKeyResourceConfigMultipleActions(name, description string) string {
	return fmt.Sprintf(`
resource "typesense_api_key" "test" {
  description = %[2]q
  actions     = ["documents:search", "documents:get"]
  collections = ["*"]
}
`, name, description)
}
