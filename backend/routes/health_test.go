package routes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {

	handler := Health.Handler

	t.Run("Returns 200 on GET", func(t *testing.T) {
		req := genRequest(t, "GET")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})

	t.Run("Returns 200 on POST", func(t *testing.T) {
		req := genRequest(t, "POST")

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}
	})
}

func genRequest(t *testing.T, m string) *http.Request {
	req, err := http.NewRequest(m, "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.WithValue(
		context.WithValue(
			context.WithValue(
				context.Background(),
				"api_name", "test name",
			),
			"api_version", "0.0.1",
		),
		"api_description",
		"test description",
	)

	return req.WithContext(ctx)
}
