#!/bin/bash
# Import an existing API key by its ID
terraform import typesense_api_key.search_key 123

# Note: Replace '123' with the actual API key ID from your Typesense server
# You can get the API key ID by listing keys via the Typesense API or admin dashboard