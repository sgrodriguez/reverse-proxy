package proxy

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// Blocker ...
type Blocker interface {
	Block(ctx context.Context, r *http.Request) (bool, error)
	Name() string
}

// Masker ...
type Masker interface {
	Mask(ctx context.Context, text []byte) ([]byte, error)
	Name() string
}

// ReverseProxy ...
type ReverseProxy struct {
	TargetURL        string
	Port             int
	proxy            *httputil.ReverseProxy
	proxyHandlerFunc func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request)
	log              zerolog.Logger

	Blockers []Blocker
	Maskers  []Masker
}

// New creates a new reverse proxy
func New(
	targetURL string,
	reverseProxyPort int,
	m []Masker,
	b []Blocker,
	log zerolog.Logger) (*ReverseProxy, error) {

	rp := &ReverseProxy{TargetURL: targetURL,
		Port:     reverseProxyPort,
		log:      log,
		Maskers:  m,
		Blockers: b,
	}

	target, err := url.Parse(rp.TargetURL)
	if err != nil {
		return nil, err
	}
	rp.proxy = httputil.NewSingleHostReverseProxy(target)
	rp.proxyHandlerFunc = func(p *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			r.Host = target.Host
			ctx := r.Context()
			for _, b := range rp.Blockers {
				if ok, err := b.Block(ctx, r); err != nil {
					log.Info().Err(err).Str("blocker_name", b.Name()).Msg("blocker error")
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				} else if ok {
					log.Info().Str("blocker_name", b.Name()).Msg("request blocked")
					http.Error(w, "blocked", http.StatusForbidden)
					return
				}
			}
			p.ServeHTTP(w, r)
		}
	}

	rp.proxy.ErrorHandler = func(rw http.ResponseWriter, r *http.Request, err error) {
		log.Error().Err(err).Msg("proxy handler error")
		if _, ok := err.(*net.OpError); ok {
			rw.WriteHeader(http.StatusBadGateway)
			rw.Write([]byte{})
			return
		}
		if _, ok := err.(*url.Error); ok {
			rw.WriteHeader(http.StatusBadGateway)
			rw.Write([]byte{})
			return
		}
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte{})
	}

	rp.proxy.Transport = http.DefaultTransport
	rp.proxy.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// We can add later support for compression
	rp.proxy.Transport.(*http.Transport).DisableCompression = true

	rp.proxy.ModifyResponse = func(r *http.Response) error {
		// Only inspect request with GET method
		ctx := r.Request.Context()
		if r.Request.Method == http.MethodGet {
			// read response body
			resBody, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			var masked []byte
			masked = resBody
			for _, m := range rp.Maskers {
				masked, err = m.Mask(ctx, masked)
				if err != nil {
					log.Err(err).Str("masker_name", m.Name()).Msg("masker error")
					// We can leak some sensitive information if we dont return an error here
					return err
				}
			}
			r.Body = io.NopCloser(bytes.NewReader(masked))
		}
		return nil
	}

	return rp, nil
}

// Start start server exit if any error occurs
func (rp *ReverseProxy) Start() (cancel func(), err error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", withLoggingHandlerFunc(rp.log, rp.proxyHandlerFunc(rp.proxy)))
	// wait for sigint or sigterm to kill server
	q := make(chan struct{})
	cancel = func() {
		q <- struct{}{}
	}
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", rp.Port),
		Handler: mux,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			rp.log.Fatal().Err(err).Msg("server error")
		}
	}()
	go func() {
		<-q
		ctx, cc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cc()
		err := srv.Shutdown(ctx)
		if err != nil {
			rp.log.Fatal().Err(err).Msg("server shutdown error")
		}
	}()
	return cancel, nil
}

func withLoggingHandlerFunc(log zerolog.Logger, handler http.HandlerFunc) http.HandlerFunc {
	loggingFn := func(rw http.ResponseWriter, req *http.Request) {
		// Create custom response writer
		lrw := &loggingResponseWriter{rw, 0, nil}
		// Create a buffer to log the req body the content
		var reqBodyBuffer strings.Builder
		reqBodyTeeReader := io.TeeReader(req.Body, &reqBodyBuffer)
		req.Body = io.NopCloser(reqBodyTeeReader)
		start := time.Now()
		// Call handler
		handler.ServeHTTP(lrw, req)
		duration := time.Since(start)
		go func() {
			// Read response body
			defer lrw.body.Close()
			resBody, _ := io.ReadAll(lrw.body)

			dicReqHeader := zerolog.Dict()
			for name, values := range req.Header {
				for _, value := range values {
					dicReqHeader.Str(name, value)
				}
			}
			dicResHeader := zerolog.Dict()
			for name, values := range lrw.Header() {
				for _, value := range values {
					dicResHeader.Str(name, value)
				}
			}
			log.Debug().
				Dur("duration_ms", duration).
				Int("status", lrw.statusCode).
				Str("method", req.Method).
				Str("path", req.URL.Path).
				Str("host", req.Host).
				Str("remote_addr", req.RemoteAddr).
				Str("user_agent", req.UserAgent()).
				Dict("response_headers", dicResHeader).
				Dict("request_headers", dicReqHeader).
				Str("request_body", reqBodyBuffer.String()).
				Str("response_body", string(resBody)).
				Msg("request received")
		}()
	}
	return http.HandlerFunc(loggingFn)
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       io.ReadCloser
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	lrw.body = io.NopCloser(bytes.NewReader(b))
	return lrw.ResponseWriter.Write(b)
}

func (lrw *loggingResponseWriter) Header() http.Header {
	return lrw.ResponseWriter.Header()
}
