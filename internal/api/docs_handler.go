package api

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"

	"store-service/docs"
)

const swaggerIndex = `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Store Service API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = () => {
      window.ui = SwaggerUIBundle({
        url: '/docs/openapi.yaml',
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis],
        layout: "BaseLayout"
      });
    };
  </script>
</body>
</html>`

func registerDocsRoutes(r chi.Router) {
	sub, _ := fs.Sub(docs.Files, ".")
	fileServer := http.FileServer(http.FS(sub))

	index := func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(swaggerIndex))
	}

	r.Get("/docs", index)
	r.Get("/docs/", index)

	r.Handle("/docs/*", http.StripPrefix("/docs", fileServer))
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs", http.StatusTemporaryRedirect)
	})
}
