package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSynonymResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSynonymResourceConfig("test_collection_for_synonym", "test_synonym", "coat", "blazer", "jacket"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_synonym.test", "name", "test_synonym"),
					resource.TestCheckResourceAttr("typesense_synonym.test", "collection_name", "test_collection_for_synonym"),
					resource.TestCheckResourceAttr("typesense_synonym.test", "synonyms.#", "3"),
					resource.TestCheckResourceAttrSet("typesense_synonym.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "typesense_synonym.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccSynonymResourceConfig("test_collection_for_synonym", "test_synonym", "coat", "blazer", "jacket", "sweater"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_synonym.test", "name", "test_synonym"),
					resource.TestCheckResourceAttr("typesense_synonym.test", "collection_name", "test_collection_for_synonym"),
					resource.TestCheckResourceAttr("typesense_synonym.test", "synonyms.#", "4"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccSynonymResource_OneWay(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create one-way synonym
			{
				Config: testAccSynonymResourceConfigOneWay("test_collection_oneway", "oneway_synonym", "sneaker", "shoe", "running shoe"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_synonym.test", "name", "oneway_synonym"),
					resource.TestCheckResourceAttr("typesense_synonym.test", "collection_name", "test_collection_oneway"),
					resource.TestCheckResourceAttr("typesense_synonym.test", "root", "sneaker"),
					resource.TestCheckResourceAttr("typesense_synonym.test", "synonyms.#", "2"),
					resource.TestCheckResourceAttrSet("typesense_synonym.test", "id"),
				),
			},
			// Update one-way synonym
			{
				Config: testAccSynonymResourceConfigOneWay("test_collection_oneway", "oneway_synonym", "sneaker", "shoe", "running shoe", "trainer"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_synonym.test", "name", "oneway_synonym"),
					resource.TestCheckResourceAttr("typesense_synonym.test", "root", "sneaker"),
					resource.TestCheckResourceAttr("typesense_synonym.test", "synonyms.#", "3"),
				),
			},
		},
	})
}

func TestAccSynonymResource_Multiple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create multiple synonyms for the same collection
			{
				Config: testAccSynonymResourceConfigMultiple("test_collection_multi_syn"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_synonym.test1", "name", "clothing_synonym"),
					resource.TestCheckResourceAttr("typesense_synonym.test1", "collection_name", "test_collection_multi_syn"),
					resource.TestCheckResourceAttr("typesense_synonym.test2", "name", "color_synonym"),
					resource.TestCheckResourceAttr("typesense_synonym.test2", "collection_name", "test_collection_multi_syn"),
					resource.TestCheckResourceAttrSet("typesense_synonym.test1", "id"),
					resource.TestCheckResourceAttrSet("typesense_synonym.test2", "id"),
				),
			},
		},
	})
}

func testAccSynonymResourceConfig(collectionName, synonymName string, synonyms ...string) string {
	synonymList := ""
	for i, syn := range synonyms {
		if i > 0 {
			synonymList += ", "
		}
		synonymList += fmt.Sprintf("%q", syn)
	}

	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "product_name"
    type = "string"
  }

  fields {
    name = "price"
    type = "int32"
  }

  default_sorting_field = "price"
}

resource "typesense_synonym" "test" {
  name            = %[2]q
  collection_name = typesense_collection.test.name
  synonyms        = [%[3]s]
}
`, collectionName, synonymName, synonymList)
}

func testAccSynonymResourceConfigOneWay(collectionName, synonymName, root string, synonyms ...string) string {
	synonymList := ""
	for i, syn := range synonyms {
		if i > 0 {
			synonymList += ", "
		}
		synonymList += fmt.Sprintf("%q", syn)
	}

	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "product_name"
    type = "string"
  }

  fields {
    name = "price"
    type = "int32"
  }

  default_sorting_field = "price"
}

resource "typesense_synonym" "test" {
  name            = %[2]q
  collection_name = typesense_collection.test.name
  root            = %[3]q
  synonyms        = [%[4]s]
}
`, collectionName, synonymName, root, synonymList)
}

func testAccSynonymResourceConfigMultiple(collectionName string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "product_name"
    type = "string"
  }

  fields {
    name = "category"
    type = "string"
  }

  fields {
    name = "price"
    type = "int32"
  }

  default_sorting_field = "price"
}

resource "typesense_synonym" "test1" {
  name            = "clothing_synonym"
  collection_name = typesense_collection.test.name
  synonyms        = ["shirt", "blouse", "top"]
}

resource "typesense_synonym" "test2" {
  name            = "color_synonym"
  collection_name = typesense_collection.test.name
  synonyms        = ["red", "crimson", "scarlet"]
}
`, collectionName)
}
