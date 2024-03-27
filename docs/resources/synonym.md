---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "typesense_synonym Resource - typesense"
subcategory: ""
description: |-
  The synonyms feature allows you to define search terms that should be considered equivalent. For eg: when you define a synonym for sneaker as shoe, searching for sneaker will now return all records with the word shoe in them, in addition to records with the word sneaker.
---

# typesense_synonym (Resource)

The synonyms feature allows you to define search terms that should be considered equivalent. For eg: when you define a synonym for sneaker as shoe, searching for sneaker will now return all records with the word shoe in them, in addition to records with the word sneaker.

## Example Usage

```terraform
resource "typesense_synonym" "my_synonym" {
  name            = "my-synonym"
  collection_name = typesense_collection.my_collection.name
  root            = "smart phone"

  synonyms = ["iphone", "android"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `collection_name` (String) Collection name
- `name` (String) Name identifier
- `synonyms` (List of String) Array of words that should be considered as synonyms.

### Optional

- `root` (String) For 1-way synonyms, indicates the root word that words in the synonyms parameter map to

### Read-Only

- `id` (String) Id identifier

## Import

Import is supported using the following syntax:

```shell
terraform import typesense_synonym.my_synonym my-synonym
```