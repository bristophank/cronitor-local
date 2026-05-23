package notify

import (
	"net/http"
	"strings"
)

// rewriteTransport returns an http.RoundTripper that redirects all requests
// to the given base URL, preserving path and query. This allows tests to
// intercept outbound HTTP calls without modifying production code.
type hostRewriter struct {
	baseURL string
}

func rewriteTransport(baseURL string) http.RoundTripper {
	return &hostRewriter{baseURL: strings.TrimRight(baseURL, "/")}
}

func (h *hostRewriter) RoundTrip(req *http.Request) (*http.Response, error) {
	cloned := req.Clone(req.Context())
	original := req.URL
	cloned.URL = cloned.URL.JoinPath("") // copy
	cloned.URL.Scheme = "http"

	// parse base
	base := h.baseURL
	if idx := strings.Index(base, "://"); idx >= 0 {
		base = base[idx+3:]
	}
	parts := strings.SplitN(base, "/", 2)
	cloned.URL.Host = parts[0]
	cloned.URL.Scheme = "http"
	cloned.Host = parts[0]
	_ = original

	return http.DefaultTransport.RoundTrip(cloned)
}
