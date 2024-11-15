package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"snippetbox.abiral.net/internal/assert"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// mock of the httphander that we can pass our secureHeader middlware
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Pass the mock HTTP handler to our secureHeaders middleware. Because
	// secureHeaders returns a http.Handler we can call its ServeHTTP()
	// method, passing in the http.ResponseRecorder and dummy http.Request to
	// execute it.
	secureHeaders(next).ServeHTTP(rr, r)
	rs := rr.Result()

	// Check middlware has correct Content-Security-Policy in header
	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	print(rs.Header.Get("Content-Security-Policy"))
	print("\n")
	// assert.Equal(t, rs.Header.Get("Content-Security-Policy"), expectedValue)

	// Check middleware has correct X-Content-Type-Options in header
	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, rs.Header.Get("Referrer-Policy"), expectedValue)

	// Check middleware has correct X-XSS-Protection header
	expectedValue = "deny"
	assert.Equal(t, rs.Header.Get("X-Frame-Options"), expectedValue)

	expectedValue = "0"
	assert.Equal(t, rs.Header.Get("X-XSS-Protection"), expectedValue)

	assert.Equal(t, rs.StatusCode, http.StatusOK)

	defer rs.Body.Close()
	body, err := io.ReadAll(rs.Body)

	if err != nil {
		t.Fatal(err)
	}

	bytes.TrimSpace(body)

	assert.Equal(t, string(body), "OK")
}
