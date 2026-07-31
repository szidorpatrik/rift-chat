package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSwaggerDocsRoutes(t *testing.T) {
	cfg := &ApiConfig{}
	mux := http.NewServeMux()
	cfg.RegisterRoutes(mux)

	t.Run("docs html", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/docs", nil)
		mux.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
		if got := recorder.Header().Get("Content-Type"); !strings.HasPrefix(got, "text/html") {
			t.Fatalf("content-type = %q, want html", got)
		}
		if !strings.Contains(recorder.Body.String(), "SwaggerUIBundle") {
			t.Fatalf("docs html missing Swagger UI bootstrap")
		}
	})

	t.Run("openapi json", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/docs/openapi.json", nil)
		mux.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
		if got := recorder.Header().Get("Content-Type"); !strings.HasPrefix(got, "application/json") {
			t.Fatalf("content-type = %q, want json", got)
		}
		body := recorder.Body.String()
		if !strings.Contains(body, `"openapi": "3.0.3"`) {
			t.Fatalf("openapi json missing version")
		}
		if !strings.Contains(body, `/api/v1/login`) {
			t.Fatalf("openapi json missing paths")
		}
	})
}
