package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliasResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccAliasResourceConfig("test_collection_for_alias", "test_alias", "test_collection_for_alias"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_alias.test", "name", "test_alias"),
					resource.TestCheckResourceAttr("typesense_alias.test", "collection_name", "test_collection_for_alias"),
					resource.TestCheckResourceAttrSet("typesense_alias.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "typesense_alias.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccAliasResourceConfig("test_collection_for_alias_updated", "test_alias", "test_collection_for_alias_updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_alias.test", "name", "test_alias"),
					resource.TestCheckResourceAttr("typesense_alias.test", "collection_name", "test_collection_for_alias_updated"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccAliasResource_MultipleAliases(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create multiple aliases pointing to the same collection
			{
				Config: testAccAliasResourceConfigMultiple("test_collection_multi", "alias1", "alias2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_alias.test1", "name", "alias1"),
					resource.TestCheckResourceAttr("typesense_alias.test1", "collection_name", "test_collection_multi"),
					resource.TestCheckResourceAttr("typesense_alias.test2", "name", "alias2"),
					resource.TestCheckResourceAttr("typesense_alias.test2", "collection_name", "test_collection_multi"),
					resource.TestCheckResourceAttrSet("typesense_alias.test1", "id"),
					resource.TestCheckResourceAttrSet("typesense_alias.test2", "id"),
				),
			},
		},
	})
}

func testAccAliasResourceConfig(collectionName, aliasName, targetCollection string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "company_name"
    type = "string"
  }

  fields {
    name = "num_employees"
    type = "int32"
    sort = true
  }

  default_sorting_field = "num_employees"
}

resource "typesense_alias" "test" {
  name            = %[2]q
  collection_name = %[3]q

  depends_on = [typesense_collection.test]
}
`, collectionName, aliasName, targetCollection)
}

func testAccAliasResourceConfigMultiple(collectionName, alias1Name, alias2Name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "title"
    type = "string"
  }

  fields {
    name = "points"
    type = "int32"
    sort = true
  }

  default_sorting_field = "points"
}

resource "typesense_alias" "test1" {
  name            = %[2]q
  collection_name = typesense_collection.test.name
}

resource "typesense_alias" "test2" {
  name            = %[3]q
  collection_name = typesense_collection.test.name
}
`, collectionName, alias1Name, alias2Name)
}
