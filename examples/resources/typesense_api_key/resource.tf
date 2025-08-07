provider "typesense" {
  api_address = "http://localhost:8108"
  api_key     = "xyz"
}

resource "typesense_api_key" "test_key" {
  description = "Test API key"
  actions     = ["documents:search"]
  collections = ["*"]
}

output "api_key_value" {
  value     = typesense_api_key.test_key.value
  sensitive = true
}

output "api_key_id" {
  value = typesense_api_key.test_key.id
}