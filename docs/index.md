---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "typesense Provider"
subcategory: ""
description: |-
  
---

# typesense Provider



## Example Usage

```terraform
provider "typesense" {
  api_key     = "xxxxxxxxxxxxxxxxxx"            // Or TYPESENSE_API_KEY enivoronment variable
  api_address = "https://your.typesense.server" // Or TYPESENSE_APP_ADDRESS enivoronment variable
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_address` (String) URL of the Typesense server. This can also be set via the `TYPESENSE_API_ADDRESS` environment variable.
- `api_key` (String, Sensitive) API Key to access the Typesense server. This can also be set via the `TYPESENSE_API_KEY` environment variable.
