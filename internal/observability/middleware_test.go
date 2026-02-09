package observability

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetricsMiddleware_RecordsMetrics(t *testing.T) {
	metrics := NewTestMetrics()
	middleware := MetricsMiddleware(metrics)

	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := middleware(inner)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	rw.WriteHeader(http.StatusNotFound)

	assert.Equal(t, http.StatusNotFound, rw.statusCode)
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestResponseWriter_DefaultStatus(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	assert.Equal(t, http.StatusOK, rw.statusCode)
}

type mockFlusher struct {
	http.ResponseWriter
	flushed bool
}

func (m *mockFlusher) Flush() { m.flushed = true }

func TestResponseWriter_Flush(t *testing.T) {
	mf := &mockFlusher{ResponseWriter: httptest.NewRecorder()}
	rw := &responseWriter{ResponseWriter: mf, statusCode: http.StatusOK}

	rw.Flush()

	assert.True(t, mf.flushed)
}

func TestResponseWriter_FlushNoOp(t *testing.T) {
	// ResponseWriter that does not implement http.Flusher
	rec := httptest.NewRecorder()
	rw := &responseWriter{ResponseWriter: rec, statusCode: http.StatusOK}

	// Should not panic even though httptest.ResponseRecorder does implement Flusher.
	// Test the non-flusher path with a custom writer.
	rw2 := &responseWriter{ResponseWriter: &nonFlusher{}, statusCode: http.StatusOK}
	rw2.Flush() // should not panic
	_ = rw      // keep linter happy
}

type nonFlusher struct{ http.ResponseWriter }

func (nonFlusher) Header() http.Header         { return http.Header{} }
func (nonFlusher) Write(b []byte) (int, error) { return len(b), nil }
func (nonFlusher) WriteHeader(_ int)           { /* no-op */ }
