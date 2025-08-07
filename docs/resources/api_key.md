# typesense_api_key

API Key resource for accessing Typesense collections with specific permissions.

API keys provide fine-grained access control to Typesense collections and operations. They can be scoped to specific collections and limited to certain actions like search, insert, delete, etc.

## Example Usage

### Basic search-only API key

```terraform
resource "typesense_api_key" "search_key" {
  description = "Search-only key for frontend"
  actions     = ["documents:search"]
  collections = ["products", "categories"]
}
```

### Admin key with all permissions

```terraform
resource "typesense_api_key" "admin_key" {
  description = "Admin key with full access"
  actions     = ["*"]
  collections = ["*"]
}
```

### Key with expiration

```terraform
resource "typesense_api_key" "temporary_key" {
  description = "Temporary key for batch import"
  actions     = ["documents:create", "documents:upsert"]
  collections = ["products"]
  expires_at  = 1735689600  # Unix timestamp
}
```

### Multi-collection key with specific permissions

```terraform
resource "typesense_api_key" "app_key" {
  description = "Application key with read/write access"
  actions     = [
    "documents:search",
    "documents:create",
    "documents:update",
    "documents:delete"
  ]
  collections = ["users", "orders", "products"]
}
```

## Argument Reference

* `description` - (Required) A human-readable description of the API key's purpose.
* `actions` - (Required) List of actions this API key can perform. See [Available Actions](#available-actions) below.
* `collections` - (Required) List of collections this API key can access. Use `["*"]` for all collections.
* `expires_at` - (Optional) Unix timestamp when the API key expires. If not specified, the key never expires.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the API key.
* `value` - The actual API key value. This is only available after creation and is sensitive.
* `value_prefix` - First few characters of the API key for identification purposes.

## Available Actions

API keys support the following actions:

### Document Operations
* `documents:search` - Search for documents
* `documents:create` - Create new documents
* `documents:upsert` - Insert or update documents
* `documents:update` - Update existing documents
* `documents:delete` - Delete documents
* `documents:import` - Bulk import documents

### Collection Operations
* `collections:create` - Create new collections
* `collections:delete` - Delete collections
* `collections:list` - List collections
* `collections:get` - Retrieve collection schema

### Alias Operations
* `aliases:list` - List collection aliases
* `aliases:create` - Create collection aliases
* `aliases:delete` - Delete collection aliases

### Synonym Operations
* `synonyms:list` - List synonyms
* `synonyms:create` - Create synonyms
* `synonyms:delete` - Delete synonyms

### Special Actions
* `*` - Wildcard for all actions (admin access)

## Collection Patterns

The `collections` attribute supports various patterns:

* **Specific collections**: `["collection1", "collection2"]`
* **All collections**: `["*"]`
* **Regex patterns**: `["org_.*"]` (matches collections starting with "org_")
* **Prefix matching**: `["user_data_.*"]`

## Import

API keys can be imported using their ID:

```bash
terraform import typesense_api_key.example 123
```

## Notes

* **No Updates**: API keys cannot be updated. Any changes require destroying and recreating the resource.
* **Value Security**: The full API key value is only returned during creation for security reasons.
* **Expiration**: Keys with expiration times will automatically become invalid after the specified time.
* **Least Privilege**: Always follow the principle of least privilege - grant only the minimum permissions necessary.
* **Value Sensitivity**: The `value` attribute is marked as sensitive and will not appear in logs or console output.