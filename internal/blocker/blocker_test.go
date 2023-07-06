package blocker_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"reverseproxy/internal/blocker"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderBlocker_Block(t *testing.T) {
	tests := map[string]struct {
		headerMap map[string]string
		request   *http.Request
		expected  bool
	}{
		"MatchingHeader": {
			headerMap: map[string]string{
				"X-Header": "Value",
			},
			request: &http.Request{
				Header: http.Header{
					"X-Header": []string{"Value"},
				},
			},
			expected: true,
		},
		"NonMatchingHeader": {
			headerMap: map[string]string{
				"X-Header": "Value",
			},
			request: &http.Request{
				Header: http.Header{
					"X-Header": []string{"OtherValue"},
				},
			},
			expected: false,
		},
		"EmptyRequestHeader": {
			headerMap: map[string]string{
				"X-Header": "Value",
			},
			request:  &http.Request{},
			expected: false,
		},
		"EmptyHeaderMap": {
			headerMap: map[string]string{},
			request: &http.Request{
				Header: http.Header{
					"X-Header": []string{"Value"},
				},
			},
			expected: false,
		},
	}
	b := &blocker.HeaderBlocker{}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b.HeaderMap = tt.headerMap
			result, err := b.Block(context.TODO(), tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestHeaderBlocker_Name(t *testing.T) {
	b := &blocker.HeaderBlocker{}
	if b.Name() != "Header Blocker" {
		t.Errorf("Expected name to be \"Header Blocker\", but got %v", b.Name())
	}
}

func TestMethodBlocker_Block(t *testing.T) {
	tests := map[string]struct {
		methods  []string
		request  *http.Request
		expected bool
	}{
		"MatchingMethod": {
			methods:  []string{http.MethodPost, http.MethodPut},
			request:  &http.Request{Method: http.MethodPut},
			expected: true,
		},
		"NonMatchingMethod": {
			methods:  []string{http.MethodPost, http.MethodPut},
			request:  &http.Request{Method: http.MethodGet},
			expected: false,
		},
		"EmptyMethods": {
			methods:  []string{},
			request:  &http.Request{Method: http.MethodGet},
			expected: false,
		},
	}
	b := &blocker.MethodBlocker{}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b.Method = tt.methods
			result, err := b.Block(context.TODO(), tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMethodBlocker_Name(t *testing.T) {
	b := &blocker.MethodBlocker{}
	if b.Name() != "Method Blocker" {
		t.Errorf("Expected name to be \"Method Blocker\", but got %v", b.Name())
	}
}

func TestPathBlocker_Block(t *testing.T) {
	tests := map[string]struct {
		paths    []string
		request  *http.Request
		expected bool
	}{
		"MatchingPath": {
			paths:    []string{"/blocked"},
			request:  &http.Request{URL: &url.URL{Path: "/blocked"}},
			expected: true,
		},
		"NonMatchingPath": {
			paths:    []string{"/blocked"},
			request:  &http.Request{URL: &url.URL{Path: "/allowed"}},
			expected: false,
		},
		"EmptyPaths": {
			paths:    []string{},
			request:  &http.Request{URL: &url.URL{Path: "/blocked"}},
			expected: false,
		},
	}
	b := &blocker.PathBlocker{}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b.Path = tt.paths
			result, err := b.Block(context.TODO(), tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPathBlocker_Name(t *testing.T) {
	b := &blocker.PathBlocker{}
	if b.Name() != "Path Blocker" {
		t.Errorf("Expected name to be \"Path Blocker\", but got %v", b.Name())
	}
}

func TestQueryParamBlocker(t *testing.T) {
	tests := map[string]struct {
		paramsMap map[string]string
		request   *http.Request
		expected  bool
	}{
		"MatchingQueryParam": {
			paramsMap: map[string]string{
				"key": "value",
			},
			request:  &http.Request{URL: &url.URL{RawQuery: "key=value"}},
			expected: true,
		},
		"NonMatchingQueryParam": {
			paramsMap: map[string]string{
				"key": "value",
			},
			request:  &http.Request{URL: &url.URL{RawQuery: "key=other"}},
			expected: false,
		},
		"EmptyParamsMap": {
			paramsMap: map[string]string{},
			request:   &http.Request{},
			expected:  false,
		},
	}
	b := &blocker.QueryParamBlocker{}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b.ParamsMap = tt.paramsMap
			result, err := b.Block(context.TODO(), tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQueryParamBlocker_Name(t *testing.T) {
	b := &blocker.QueryParamBlocker{}
	if b.Name() != "Query Param Blocker" {
		t.Errorf("Expected name to be \"Query Param Blocker\", but got %v", b.Name())
	}
}
