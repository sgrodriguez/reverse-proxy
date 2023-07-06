package blocker

import (
	"context"
	"net/http"
)

// HeaderBlocker ...
type HeaderBlocker struct {
	HeaderMap map[string]string
}

// Block every request that has a header with the given name and value.
func (hb *HeaderBlocker) Block(ctx context.Context, r *http.Request) (bool, error) {
	for k, v := range hb.HeaderMap {
		if r.Header.Get(k) == v {
			return true, nil
		}
	}
	return false, nil
}

// Name returns the name of the blocker.
func (hb *HeaderBlocker) Name() string {
	return "Header Blocker"
}

// MethodBlocker ...
type MethodBlocker struct {
	Method []string
}

// Block every request that has the given method.
func (mb *MethodBlocker) Block(ctx context.Context, r *http.Request) (bool, error) {
	for _, m := range mb.Method {
		if r.Method == m {
			return true, nil
		}
	}
	return false, nil
}

// Name returns the name of the blocker.
func (mb *MethodBlocker) Name() string {
	return "Method Blocker"
}

// PathBlocker ...
type PathBlocker struct {
	Path []string
}

// Block every request that has the given path.
func (pb *PathBlocker) Block(ctx context.Context, r *http.Request) (bool, error) {
	for _, p := range pb.Path {
		if r.URL.Path == p {
			return true, nil
		}
	}
	return false, nil
}

// Name returns the name of the blocker.
func (pb *PathBlocker) Name() string {
	return "Path Blocker"
}

// QueryParamBlocker ...
type QueryParamBlocker struct {
	ParamsMap map[string]string
}

// Block every request that has the given query parameter.
func (qpb *QueryParamBlocker) Block(ctx context.Context, r *http.Request) (bool, error) {
	for k, v := range qpb.ParamsMap {
		if r.URL.Query().Get(k) == v {
			return true, nil
		}
	}
	return false, nil
}

// Name returns the name of the blocker.
func (qpb *QueryParamBlocker) Name() string {
	return "Query Param Blocker"
}
