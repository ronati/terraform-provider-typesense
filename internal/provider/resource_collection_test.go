package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCollectionResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccCollectionResourceConfig("test_collection"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection"),
					resource.TestCheckResourceAttr("typesense_collection.test", "default_sorting_field", "num_employees"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "2"),
					resource.TestCheckResourceAttrSet("typesense_collection.test", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "typesense_collection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing - add a new field
			{
				Config: testAccCollectionResourceConfigUpdated("test_collection"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "3"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccCollectionResource_WithNestedFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create collection with nested fields enabled
			{
				Config: testAccCollectionResourceConfigNested("test_collection_nested"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_nested"),
					resource.TestCheckResourceAttr("typesense_collection.test", "enable_nested_fields", "true"),
					resource.TestCheckResourceAttrSet("typesense_collection.test", "id"),
				),
			},
		},
	})
}

func TestAccCollectionResource_WithSymbolsAndTokens(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create collection with custom symbols and tokens
			{
				Config: testAccCollectionResourceConfigSymbolsTokens("test_collection_symbols"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_symbols"),
					resource.TestCheckResourceAttr("typesense_collection.test", "symbols_to_index.#", "2"),
					resource.TestCheckResourceAttr("typesense_collection.test", "token_separators.#", "1"),
					resource.TestCheckResourceAttrSet("typesense_collection.test", "id"),
				),
			},
		},
	})
}

func TestAccCollectionResource_WithOptionalFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create collection with optional fields
			{
				Config: testAccCollectionResourceConfigOptionalFields("test_collection_optional"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_optional"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "3"),
					resource.TestCheckResourceAttrSet("typesense_collection.test", "id"),
				),
			},
		},
	})
}

func TestAccCollectionResource_WithArrayFields(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create collection with array field types
			{
				Config: testAccCollectionResourceConfigArrayFields("test_collection_arrays"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_arrays"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "4"),
					resource.TestCheckResourceAttrSet("typesense_collection.test", "id"),
				),
			},
		},
	})
}

func TestAccCollectionResource_FieldDefaults(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionResourceConfigMinimalFields("test_collection_defaults"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_defaults"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "2"),
					resource.TestCheckResourceAttrSet("typesense_collection.test", "id"),
				),
			},
		},
	})
}

func TestAccCollectionResource_FieldDefaultsExplicit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionResourceConfigExplicitFieldDefaults("test_collection_explicit"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_explicit"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "2"),
					resource.TestCheckResourceAttrSet("typesense_collection.test", "id"),
				),
			},
		},
	})
}

func TestAccCollectionResource_FieldWithAllAttributes(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionResourceConfigAllFieldAttributes("test_collection_all_attrs"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_all_attrs"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "2"),
					resource.TestCheckResourceAttrSet("typesense_collection.test", "id"),
				),
			},
		},
	})
}

func TestAccCollectionResource_WithEmbedAndImage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionResourceConfigWithEmbed("test_collection_embed"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_embed"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("typesense_collection.test", "fields.*", map[string]string{
						"name": "thumbnailImage",
						"type": "image",
					}),
					resource.TestCheckTypeSetElemNestedAttrs("typesense_collection.test", "fields.*", map[string]string{
						"name":                          "embedding",
						"type":                          "float[]",
						"num_dim":                       "512",
						"embed.%":                       "2",
						"embed.from.#":                  "1",
						"embed.from.0":                  "thumbnailImage",
						"embed.model_config.%":          "1",
						"embed.model_config.model_name": "ts/clip-vit-b-p32",
					}),
					resource.TestCheckResourceAttrSet("typesense_collection.test", "id"),
				),
			},
		},
	})
}

func testAccCollectionResourceConfig(name string) string {
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
`, name)
}

func testAccCollectionResourceConfigUpdated(name string) string {
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

  fields {
    name = "country"
    type = "string"
    optional = true
  }

  default_sorting_field = "num_employees"
}
`, name)
}

func testAccCollectionResourceConfigNested(name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name                   = %[1]q
  enable_nested_fields   = true

  fields {
    name = "company_name"
    type = "string"
  }

  fields {
    name = "metadata"
    type = "object"
  }

  fields {
    name = "score"
    type = "int32"
    sort = true
  }

  default_sorting_field = "score"
}
`, name)
}

func testAccCollectionResourceConfigSymbolsTokens(name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name              = %[1]q
  symbols_to_index  = ["+", "-"]
  token_separators  = ["/"]

  fields {
    name = "title"
    type = "string"
  }

  fields {
    name = "rating"
    type = "int32"
    sort = true
  }

  default_sorting_field = "rating"
}
`, name)
}

func testAccCollectionResourceConfigOptionalFields(name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "title"
    type = "string"
  }

  fields {
    name = "description"
    type = "string"
    optional = true
  }

  fields {
    name = "rank"
    type = "int32"
    sort = true
  }

  default_sorting_field = "rank"
}
`, name)
}

func testAccCollectionResourceConfigArrayFields(name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "title"
    type = "string"
  }

  fields {
    name = "tags"
    type = "string[]"
  }

  fields {
    name = "scores"
    type = "int32[]"
  }

  fields {
    name = "rating"
    type = "int32"
    sort = true
  }

  default_sorting_field = "rating"
}
`, name)
}

func testAccCollectionResourceConfigMinimalFields(name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "title"
    type = "string"
  }

  fields {
    name = "rating"
    type = "int32"
    sort = true
  }

  default_sorting_field = "rating"
}
`, name)
}

func testAccCollectionResourceConfigExplicitFieldDefaults(name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name     = "title"
    type     = "string"
    facet    = false
    index    = true
    optional = false
    infix    = false
    stem     = false
    store    = true
    locale   = ""
  }

  fields {
    name     = "rating"
    type     = "int32"
    sort     = true
    facet    = false
    index    = true
    optional = false
    infix    = false
    stem     = false
    store    = true
  }

  default_sorting_field = "rating"
}
`, name)
}

func testAccCollectionResourceConfigAllFieldAttributes(name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name     = "title"
    type     = "string"
    facet    = true
    index    = true
    optional = true
    sort     = false
    infix    = true
    stem     = true
    store    = true
    locale   = "en"
  }

  fields {
    name     = "rating"
    type     = "int32"
    sort     = true
    facet    = true
    index    = true
    optional = false
    infix    = false
    stem     = false
    store    = true
  }

  default_sorting_field = "rating"
}
`, name)
}

func testAccCollectionResourceConfigWithEmbed(name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "thumbnailImage"
    type = "image"
  }

  fields {
    name    = "embedding"
    type    = "float[]"
    num_dim = 512

    embed {
      from = ["thumbnailImage"]

      model_config {
        model_name = "ts/clip-vit-b-p32"
      }
    }
  }
}
`, name)
}

func TestAccCollectionResource_VectorFieldWithNumDim(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionResourceConfigVectorField("test_collection_vector"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_vector"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "3"),
					resource.TestCheckResourceAttrSet("typesense_collection.test", "id"),
				),
			},
			{
				ResourceName:      "typesense_collection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCollectionResource_VectorFieldUpdate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionResourceConfigBasic("test_collection_vector_update"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_vector_update"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "2"),
				),
			},
			{
				Config: testAccCollectionResourceConfigAddVectorField("test_collection_vector_update"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_vector_update"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "3"),
				),
			},
			{
				Config: testAccCollectionResourceConfigBasic("test_collection_vector_update"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_vector_update"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "2"),
				),
			},
		},
	})
}

func TestAccCollectionResource_VectorFieldNumDimChange(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionResourceConfigVectorFieldDim("test_collection_dim_change", 128),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_dim_change"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "3"),
				),
			},
			{
				Config: testAccCollectionResourceConfigVectorFieldDim("test_collection_dim_change", 256),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_dim_change"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "3"),
				),
			},
		},
	})
}

func TestAccCollectionResource_FloatArrayWithoutNumDim(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCollectionResourceConfigFloatArrayNoNumDim("test_collection_float_array"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("typesense_collection.test", "name", "test_collection_float_array"),
					resource.TestCheckResourceAttr("typesense_collection.test", "fields.#", "3"),
					resource.TestCheckResourceAttrSet("typesense_collection.test", "id"),
				),
			},
			{
				ResourceName:      "typesense_collection.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCollectionResourceConfigVectorField(name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "title"
    type = "string"
  }

  fields {
    name    = "embedding"
    type    = "float[]"
    num_dim = 768
  }

  fields {
    name = "rating"
    type = "int32"
    sort = true
  }

  default_sorting_field = "rating"
}
`, name)
}

func testAccCollectionResourceConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "title"
    type = "string"
  }

  fields {
    name = "rating"
    type = "int32"
    sort = true
  }

  default_sorting_field = "rating"
}
`, name)
}

func testAccCollectionResourceConfigAddVectorField(name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "title"
    type = "string"
  }

  fields {
    name    = "embedding"
    type    = "float[]"
    num_dim = 512
  }

  fields {
    name = "rating"
    type = "int32"
    sort = true
  }

  default_sorting_field = "rating"
}
`, name)
}

func testAccCollectionResourceConfigVectorFieldDim(name string, dim int) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "title"
    type = "string"
  }

  fields {
    name    = "embedding"
    type    = "float[]"
    num_dim = %[2]d
  }

  fields {
    name = "rating"
    type = "int32"
    sort = true
  }

  default_sorting_field = "rating"
}
`, name, dim)
}

func testAccCollectionResourceConfigFloatArrayNoNumDim(name string) string {
	return fmt.Sprintf(`
resource "typesense_collection" "test" {
  name = %[1]q

  fields {
    name = "title"
    type = "string"
  }

  fields {
    name = "scores"
    type = "float[]"
  }

  fields {
    name = "rating"
    type = "int32"
    sort = true
  }

  default_sorting_field = "rating"
}
`, name)
}
