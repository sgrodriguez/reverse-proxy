package proxy_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"reverseproxy/proxy"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReverseProxy(t *testing.T) {
	// Create sample masker
	masker := &MockMasker{
		fn: func(text []byte) ([]byte, error) {
			return text, nil
		},
	}
	// Create sample blocker
	blocker := &MockBlocker{
		fn: func() (bool, error) {
			return false, nil
		},
	}
	// Create logger
	logger := zerolog.Nop()
	// Create a test server
	targetServerResponse := "Hello World"
	targetServerResponseCode := http.StatusOK
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(targetServerResponseCode)
		w.Write([]byte(targetServerResponse))
	}))
	defer targetServer.Close()
	// Create the reverse proxy
	reverseProxy, err := proxy.New(targetServer.URL,
		8080,
		[]proxy.Masker{masker},
		[]proxy.Blocker{blocker},
		logger)
	require.NoError(t, err)
	// Start the reverse proxy server
	cancel, err := reverseProxy.Start()
	defer cancel()
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080", nil)
	require.NoError(t, err)
	client := http.Client{}
	t.Run("Test No Masker and Blocker", func(t *testing.T) {
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, "invalid status code")
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		assert.Equal(t, targetServerResponse, buf.String(), "invalid response body")
	})
	t.Run("Test Masker", func(t *testing.T) {
		// Modify the masker to return a different value
		masker.fn = func(text []byte) ([]byte, error) {
			return []byte("Masked"), nil
		}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode, "invalid status code")
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		assert.Equal(t, "Masked", buf.String(), "invalid response body")
		t.Run("masker do not mask methods different from GET", func(t *testing.T) {
			newReq, err := http.NewRequest(http.MethodPost, "http://localhost:8080", nil)
			require.NoError(t, err)
			resp, err := client.Do(newReq)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode, "invalid status code")
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			assert.Equal(t, targetServerResponse, buf.String(), "invalid response body")
		})
	})
	t.Run("Test masker return error", func(t *testing.T) {
		masker.fn = func(text []byte) ([]byte, error) {
			return nil, errors.New("masker error")
		}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode, "invalid status code")
	})
	t.Run("Test blocker", func(t *testing.T) {
		blocker.fn = func() (bool, error) {
			return true, nil
		}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusForbidden, resp.StatusCode, "invalid status code")
	})
	t.Run("Test return error", func(t *testing.T) {
		blocker.fn = func() (bool, error) {
			return true, errors.New("blocker error")
		}
		resp, err := client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusInternalServerError, resp.StatusCode, "invalid status code")
	})
}

func TestReverseProxy_NoTargetSever(t *testing.T) {
	// Create the reverse proxy
	logger := zerolog.Nop()
	reverseProxy, err := proxy.New("http://localhost:123987",
		8081,
		[]proxy.Masker{},
		[]proxy.Blocker{},
		logger)
	require.NoError(t, err)
	// Start the reverse proxy server
	cancel, err := reverseProxy.Start()
	defer cancel()
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8081", nil)
	require.NoError(t, err)
	client := http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadGateway, resp.StatusCode, "invalid status code")
}

type MockMasker struct {
	fn func([]byte) ([]byte, error)
}

func (m *MockMasker) Mask(ctx context.Context, text []byte) ([]byte, error) {
	return m.fn(text)
}

func (m *MockMasker) Name() string {
	return "Test Masker"
}

type MockBlocker struct {
	fn func() (bool, error)
}

func (b *MockBlocker) Block(ctx context.Context, r *http.Request) (bool, error) {
	return b.fn()
}

func (b *MockBlocker) Name() string {
	return "Test Blocker"
}
