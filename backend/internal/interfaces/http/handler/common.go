package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/rostmebel/backend/internal/domain/apperror"
)

var validate = validator.New()

func init() {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, err error) {
	appErr, ok := apperror.From(err)
	if !ok {
		appErr = apperror.Internal(err)
	}

	respondWithJSON(w, statusForError(appErr), map[string]any{
		"error": appErr,
	})
}

func decodeAndValidate(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return apperror.Wrap(err, apperror.CodeInvalidJSON, "Invalid JSON request body", jsonDecodeMeta(err))
	}
	if err := validate.Struct(v); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			return apperror.New(apperror.CodeValidationFailed, "Request validation failed", map[string]any{
				"fields": validationErrorMeta(validationErrors),
			})
		}
		return apperror.Wrap(err, apperror.CodeValidationFailed, "Request validation failed", nil)
	}
	return nil
}

func jsonDecodeMeta(err error) map[string]any {
	if err == nil {
		return nil
	}

	const unknownFieldPrefix = `json: unknown field "`
	if message := err.Error(); strings.HasPrefix(message, unknownFieldPrefix) {
		field := strings.TrimSuffix(strings.TrimPrefix(message, unknownFieldPrefix), `"`)
		return map[string]any{
			"field":  field,
			"reason": "unknown_field",
		}
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return map[string]any{
			"offset": syntaxErr.Offset,
			"reason": "syntax_error",
		}
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return map[string]any{
			"field":  typeErr.Field,
			"value":  typeErr.Value,
			"reason": "type_mismatch",
		}
	}

	return map[string]any{
		"reason": "malformed_json",
	}
}

func validationErrorMeta(validationErrors validator.ValidationErrors) []map[string]string {
	fields := make([]map[string]string, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		fields = append(fields, map[string]string{
			"field": fieldErr.Field(),
			"rule":  fieldErr.Tag(),
			"param": fieldErr.Param(),
		})
	}
	return fields
}

func statusForError(err *apperror.Error) int {
	switch err.Code {
	case apperror.CodeInvalidJSON,
		apperror.CodeValidationFailed,
		apperror.CodeInvalidID,
		apperror.CodeInvalidQuery,
		apperror.CodeUploadFileMissing,
		apperror.CodeUploadInvalidType,
		apperror.CodeReviewInvalidPhone:
		return http.StatusBadRequest
	case apperror.CodeUnauthorized,
		apperror.CodeAuthInvalidCredentials,
		apperror.CodeAuthInvalidRefreshToken:
		return http.StatusUnauthorized
	case apperror.CodeOrderIPBlocked,
		apperror.CodeReviewNotAllowed:
		return http.StatusForbidden
	case apperror.CodeProjectNotFound,
		apperror.CodeOrderNotFound:
		return http.StatusNotFound
	case apperror.CodeOrderRateLimited:
		return http.StatusTooManyRequests
	case apperror.CodeUploadFileTooLarge:
		return http.StatusRequestEntityTooLarge
	default:
		return http.StatusInternalServerError
	}
}
