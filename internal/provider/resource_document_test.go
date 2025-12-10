package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDocumentResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccDocumentResourceConfig("test_collection_for_doc", "doc1", "Product One", 100),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_document.test", "name", "doc1"),
					resource.TestCheckResourceAttr("typesense_document.test", "collection_name", "test_collection_for_doc"),
					resource.TestCheckResourceAttrSet("typesense_document.test", "id"),
					resource.TestCheckResourceAttrSet("typesense_document.test", "document"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "typesense_document.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccDocumentResourceConfig("test_collection_for_doc", "doc1", "Product One Updated", 150),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_document.test", "name", "doc1"),
					resource.TestCheckResourceAttr("typesense_document.test", "collection_name", "test_collection_for_doc"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccDocumentResource_ComplexJSON(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create document with complex JSON
			{
				Config: testAccDocumentResourceConfigComplex("test_collection_complex", "complex_doc1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_document.test", "name", "complex_doc1"),
					resource.TestCheckResourceAttr("typesense_document.test", "collection_name", "test_collection_complex"),
					resource.TestCheckResourceAttrSet("typesense_document.test", "id"),
					resource.TestCheckResourceAttrSet("typesense_document.test", "document"),
				),
			},
		},
	})
}

func TestAccDocumentResource_Multiple(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create multiple documents
			{
				Config: testAccDocumentResourceConfigMultiple("test_collection_multi_doc"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_document.test1", "name", "doc1"),
					resource.TestCheckResourceAttr("typesense_document.test1", "collection_name", "test_collection_multi_doc"),
					resource.TestCheckResourceAttr("typesense_document.test2", "name", "doc2"),
					resource.TestCheckResourceAttr("typesense_document.test2", "collection_name", "test_collection_multi_doc"),
					resource.TestCheckResourceAttrSet("typesense_document.test1", "id"),
					resource.TestCheckResourceAttrSet("typesense_document.test2", "id"),
				),
			},
		},
	})
}

func TestAccDocumentResource_WithArrays(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create document with array fields
			{
				Config: testAccDocumentResourceConfigWithArrays("test_collection_arrays", "array_doc1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_document.test", "name", "array_doc1"),
					resource.TestCheckResourceAttr("typesense_document.test", "collection_name", "test_collection_arrays"),
					resource.TestCheckResourceAttrSet("typesense_document.test", "id"),
				),
			},
		},
	})
}

func testAccDocumentResourceConfig(collectionName, docName, productName string, price int) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "product_name"
    type = "string"
    store = true
  }

  fields {
    name = "price"
    type = "int32"
    sort = true
    store = true
  }

  default_sorting_field = "price"
}

resource "typesense_document" "test" {
  name            = %[2]q
  collection_name = typesense_collection.test.name
  document = jsonencode({
    product_name = %[3]q
    price        = %[4]d
  })
}
`, collectionName, docName, productName, price)
}

func testAccDocumentResourceConfigComplex(collectionName, docName string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "title"
    type = "string"
    store = true
  }

  fields {
    name = "description"
    type = "string"
    store = true
  }

  fields {
    name = "rating"
    type = "float"
    store = true
  }

  fields {
    name = "in_stock"
    type = "bool"
    store = true
  }

  fields {
    name = "views"
    type = "int32"
    sort = true
    store = true
  }

  default_sorting_field = "views"
}

resource "typesense_document" "test" {
  name            = %[2]q
  collection_name = typesense_collection.test.name
  document = jsonencode({
    title       = "Complex Product"
    description = "This is a complex product with multiple fields"
    rating      = 4.5
    in_stock    = true
    views       = 1000
  })
}
`, collectionName, docName)
}

func testAccDocumentResourceConfigMultiple(collectionName string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "title"
    type = "string"
    store = true
  }

  fields {
    name = "count"
    type = "int32"
    sort = true
    store = true
  }

  default_sorting_field = "count"
}

resource "typesense_document" "test1" {
  name            = "doc1"
  collection_name = typesense_collection.test.name
  document = jsonencode({
    title = "Document One"
    count = 10
  })
}

resource "typesense_document" "test2" {
  name            = "doc2"
  collection_name = typesense_collection.test.name
  document = jsonencode({
    title = "Document Two"
    count = 20
  })
}
`, collectionName)
}

func testAccDocumentResourceConfigWithArrays(collectionName, docName string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "title"
    type = "string"
    store = true
  }

  fields {
    name = "tags"
    type = "string[]"
    store = true
  }

  fields {
    name = "scores"
    type = "int32[]"
    store = true
  }

  fields {
    name = "rating"
    type = "int32"
    sort = true
    store = true
  }

  default_sorting_field = "rating"
}

resource "typesense_document" "test" {
  name            = %[2]q
  collection_name = typesense_collection.test.name
  document = jsonencode({
    title  = "Product with Arrays"
    tags   = ["electronics", "gadget", "popular"]
    scores = [85, 90, 95]
    rating = 90
  })
}
`, collectionName, docName)
}
