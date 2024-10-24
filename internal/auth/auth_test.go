package auth

import (
	"net/http"
	"testing"
)

// GetAPIKey -
func TestGetAPIKey(t *testing.T) {
	t.Run("testing GetAPIKey()", func(t *testing.T) {
		expected := "ApiKey mockkey"

		header := http.Header{}
		header.Set("Authorization", expected)

		result, err := GetAPIKey(header)
		if err != nil {
			t.Error(err)
		}

		if result != "mockkey" {
			t.Errorf("\nexpected: %v\nactual: %v", "mockkey", result)
		}
	})
}
