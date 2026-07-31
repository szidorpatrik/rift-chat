package api

import (
	_ "embed"
	"net/http"
)

//go:embed openapi.json
var openAPISpec []byte

const swaggerDocsHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Rift Chat API Docs</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    body {
      margin: 0;
      background: #f5f7fb;
      color: #0f172a;
    }

    .docs-shell {
      min-height: 100vh;
      padding: 20px;
      box-sizing: border-box;
    }

    #swagger-ui {
      color: #0f172a;
      max-width: 1440px;
      margin: 0 auto;
      background: #ffffff;
      border: 1px solid rgba(15, 23, 42, 0.12);
      border-radius: 3px;
      overflow: hidden;
      box-shadow: 0 28px 64px rgba(15, 23, 42, 0.18);
    }

    .swagger-ui .topbar {
      display: none;
    }

    .swagger-ui .info,
    .swagger-ui .scheme-container,
    .swagger-ui .opblock,
    .swagger-ui .model-box,
    .swagger-ui .responses-wrapper,
    .swagger-ui .responses-inner,
    .swagger-ui .response-col_status,
    .swagger-ui .response-col_links,
    .swagger-ui .opblock-summary,
    .swagger-ui .opblock-summary-description,
    .swagger-ui .tab,
    .swagger-ui .parameter__name,
    .swagger-ui .parameter__type,
    .swagger-ui .parameter__in,
    .swagger-ui .response-col_description,
    .swagger-ui .responses-table,
    .swagger-ui .opblock-description-wrapper,
    .swagger-ui .markdown,
    .swagger-ui .renderedMarkdown {
      color: var(--text-color);
    }

    .swagger-ui .opblock,
    .swagger-ui .scheme-container,
    .swagger-ui .model-box,
    .swagger-ui .responses-wrapper,
    .swagger-ui .responses-inner {
      background: #ffffff;
      border-color: rgba(15, 23, 42, 0.12);
    }

    .swagger-ui .opblock .opblock-section-header {
      background: #f1f5f9;
    }

    .swagger-ui .opblock .tab-header {
      color: #0f172a !important;
    }

    .swagger-ui .btn,
    .swagger-ui .btn.authorize,
    .swagger-ui .btn.try-out__btn,
    .swagger-ui select,
    .swagger-ui input {
      border-color: rgba(15, 23, 42, 0.12);
    }

    .swagger-ui .opblock .opblock-summary-method {
      border-radius: 3px;
    }

    .swagger-ui .markdown code,
    .swagger-ui .renderedMarkdown code {
      background: rgba(99, 102, 241, 0.12);
      color: #0f172a;
    }

    .swagger-ui .opblock-body pre.microlight {
      color: #e5e7eb !important;
    }

    .swagger-ui .responses-inner h4,
    .swagger-ui .responses-inner h5,
    .swagger-ui .responses-inner .response-col_description__inner,
    .swagger-ui .responses-inner .response-col_status,
    .swagger-ui .responses-inner .response-col_links,
    .swagger-ui .responses-inner .response-col_description,
    .swagger-ui .opblock-section-header label,
    .swagger-ui .parameter__name,
    .swagger-ui .parameter__type,
    .swagger-ui .parameter__in {
      color: #0f172a;
    }
  </style>
</head>
<body>
  <div class="docs-shell">
    <div id="swagger-ui"></div>
  </div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: '/api/v1/docs/openapi.json',
      dom_id: '#swagger-ui',
      deepLinking: true,
      displayRequestDuration: true,
      persistAuthorization: true
    })
  </script>
</body>
</html>`

func (cfg *ApiConfig) SwaggerDocsHandler(w http.ResponseWriter, r *http.Request) {
	_ = cfg
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(swaggerDocsHTML))
}

func (cfg *ApiConfig) SwaggerSpecHandler(w http.ResponseWriter, r *http.Request) {
	_ = cfg
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(openAPISpec)
}
