package gemini

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestSearchProductsFallsBackWhenConfiguredModelDoesNotExist(t *testing.T) {
	var requestedPaths []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestedPaths = append(requestedPaths, r.URL.Path)

		if strings.Contains(r.URL.Path, "gemma-4-31b") {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":{"code":404,"message":"model not found","status":"NOT_FOUND"}}`))
			return
		}

		if !strings.Contains(r.URL.Path, "gemini-2.5-flash") {
			t.Fatalf("unexpected model path: %s", r.URL.Path)
		}

		var req generateContentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.GenerationConfig.ResponseMimeType != "application/json" {
			t.Fatalf("expected JSON response mime type, got %q", req.GenerationConfig.ResponseMimeType)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"candidates":[{"content":{"parts":[{"text":"{\"ids\":[2,5]}"}]}}]}`))
	}))
	defer server.Close()

	client := NewClientWithOptions(ClientOptions{
		APIKey:         "test-key",
		Model:          "models/gemma-4-31b",
		FallbackModels: []string{"gemini-2.5-flash"},
		BaseURL:        server.URL,
		HTTPClient:     server.Client(),
	})

	ids, err := client.SearchProducts(context.Background(), "кухня", "[]")
	if err != nil {
		t.Fatalf("search products: %v", err)
	}
	if !reflect.DeepEqual(ids, []int64{2, 5}) {
		t.Fatalf("unexpected ids: %#v", ids)
	}
	if len(requestedPaths) != 2 {
		t.Fatalf("expected primary plus fallback request, got %d", len(requestedPaths))
	}
}

func TestSearchProductsDoesNotFallbackOnAuthError(t *testing.T) {
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requests++
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":{"code":403,"message":"bad key","status":"PERMISSION_DENIED"}}`))
	}))
	defer server.Close()

	client := NewClientWithOptions(ClientOptions{
		APIKey:         "bad-key",
		Model:          "gemini-2.5-flash",
		FallbackModels: []string{"gemini-2.5-flash-lite"},
		BaseURL:        server.URL,
		HTTPClient:     server.Client(),
	})

	_, err := client.SearchProducts(context.Background(), "кухня", "[]")
	if err == nil {
		t.Fatal("expected error")
	}
	if requests != 1 {
		t.Fatalf("auth errors must not try fallback models, got %d requests", requests)
	}
}

func TestSearchProductsDisabledWithoutAPIKey(t *testing.T) {
	client := NewClientWithOptions(ClientOptions{Model: "gemini-2.5-flash"})

	_, err := client.SearchProducts(context.Background(), "кухня", "[]")
	if !errors.Is(err, ErrDisabled) {
		t.Fatalf("expected ErrDisabled, got %v", err)
	}
}
