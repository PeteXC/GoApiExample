package respond

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithError(t *testing.T) {
	var tests = []struct {
		msg    string
		status int
	}{
		{
			msg:    "not found",
			status: http.StatusNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.msg, func(t *testing.T) {
			// Arrange.
			w := httptest.NewRecorder()

			WithError(w, tt.msg, tt.status)

			if tt.status != w.Result().StatusCode {
				t.Errorf("expected status: %d, got %d", tt.status, w.Result().StatusCode)
			}
		})

	}
}
