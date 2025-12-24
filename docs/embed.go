package docs

import "embed"

//go:embed openapi.yaml postman_collection.json
var Files embed.FS
