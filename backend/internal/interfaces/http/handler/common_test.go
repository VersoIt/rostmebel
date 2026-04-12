package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rostmebel/backend/internal/domain/apperror"
)

func TestRespondWithErrorUsesStructuredContract(t *testing.T) {
	rec := httptest.NewRecorder()

	respondWithError(rec, apperror.New(apperror.CodeOrderRateLimited, "Order rate limit exceeded", map[string]any{
		"limit": 5,
	}))

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("expected status %d, got %d", http.StatusTooManyRequests, rec.Code)
	}

	var payload struct {
		Error struct {
			Code    string         `json:"code"`
			Message string         `json:"message"`
			Meta    map[string]any `json:"meta"`
		} `json:"error"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.Error.Code != apperror.CodeOrderRateLimited {
		t.Fatalf("expected code %q, got %q", apperror.CodeOrderRateLimited, payload.Error.Code)
	}
	if payload.Error.Message == "" {
		t.Fatal("expected non-empty error message")
	}
	if payload.Error.Meta["limit"].(float64) != 5 {
		t.Fatalf("expected limit metadata, got %#v", payload.Error.Meta)
	}
}

func TestDecodeAndValidateReturnsFieldMetadata(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")

	var body struct {
		ClientName  string `json:"client_name" validate:"required"`
		ClientPhone string `json:"client_phone" validate:"required"`
	}
	err := decodeAndValidate(req, &body)
	if err == nil {
		t.Fatal("expected validation error")
	}

	appErr, ok := apperror.From(err)
	if !ok {
		t.Fatalf("expected app error, got %T", err)
	}
	if appErr.Code != apperror.CodeValidationFailed {
		t.Fatalf("expected validation code, got %q", appErr.Code)
	}

	fields, ok := appErr.Meta["fields"].([]map[string]string)
	if !ok {
		t.Fatalf("expected fields metadata, got %#v", appErr.Meta["fields"])
	}
	if len(fields) != 2 {
		t.Fatalf("expected 2 field errors, got %d", len(fields))
	}
	if fields[0]["field"] != "client_name" || fields[1]["field"] != "client_phone" {
		t.Fatalf("expected json field names, got %#v", fields)
	}
}
