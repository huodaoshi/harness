package knowledgeindexing

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultFetchTimeout = 60 * time.Second
	defaultMaxBodyBytes = 32 << 20
	maxRedirects        = 5
)

// FetchOptions configures outbound HTTP for source_type=2 URL ingest.
type FetchOptions struct {
	AllowHosts                 []string
	AllowedContentTypePrefixes []string
}

func FetchHostAllowed(hostname string, allowHosts []string) bool {
	if len(allowHosts) == 0 {
		return true
	}
	h := strings.TrimSpace(strings.ToLower(hostname))
	if h == "" {
		return false
	}
	for _, a := range allowHosts {
		if strings.EqualFold(strings.TrimSpace(a), h) {
			return true
		}
	}
	return false
}

func FetchContentTypeAllowed(contentType string, allowedPrefixes []string) bool {
	if len(allowedPrefixes) == 0 {
		return true
	}
	base := strings.TrimSpace(strings.ToLower(strings.Split(contentType, ";")[0]))
	if base == "" {
		return false
	}
	for _, p := range allowedPrefixes {
		p = strings.TrimSpace(strings.ToLower(p))
		if p != "" && strings.HasPrefix(base, p) {
			return true
		}
	}
	return false
}

func fetchSource(ctx context.Context, rawURL string, opts FetchOptions) ([]byte, string, string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, "", "", fmt.Errorf("knowledgeindexing: fetch: parse url: %w", err)
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return nil, "", "", fmt.Errorf("knowledgeindexing: fetch: unsupported scheme %q", u.Scheme)
	}
	host := strings.ToLower(u.Hostname())
	if host == "" || isBlockedHost(host) {
		return nil, "", "", fmt.Errorf("knowledgeindexing: fetch: blocked host %q", host)
	}
	if len(opts.AllowHosts) > 0 && !FetchHostAllowed(u.Hostname(), opts.AllowHosts) {
		return nil, "", "", fmt.Errorf("knowledgeindexing: fetch: host %q not in allowlist", host)
	}
	if err := forbidPrivateURLHost(ctx, u.Hostname()); err != nil {
		return nil, "", "", err
	}

	client := &http.Client{
		Timeout: defaultFetchTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("knowledgeindexing: fetch: too many redirects")
			}
			rh := req.URL.Hostname()
			if rh == "" || isBlockedHost(strings.ToLower(rh)) {
				return fmt.Errorf("knowledgeindexing: fetch: blocked redirect host %q", rh)
			}
			if len(opts.AllowHosts) > 0 && !FetchHostAllowed(rh, opts.AllowHosts) {
				return fmt.Errorf("knowledgeindexing: fetch: redirect host %q not in allowlist", strings.ToLower(rh))
			}
			return forbidPrivateURLHost(req.Context(), rh)
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, "", "", fmt.Errorf("knowledgeindexing: fetch: new request: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", "", fmt.Errorf("knowledgeindexing: fetch: get: %w", err)
	}
	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", "", fmt.Errorf("knowledgeindexing: fetch: status %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if len(opts.AllowedContentTypePrefixes) > 0 && !FetchContentTypeAllowed(ct, opts.AllowedContentTypePrefixes) {
		return nil, "", "", fmt.Errorf("knowledgeindexing: fetch: content-type %q not allowed", ct)
	}

	limited := io.LimitReader(resp.Body, defaultMaxBodyBytes+1)
	data, err := io.ReadAll(limited)
	if err != nil {
		return nil, "", "", fmt.Errorf("knowledgeindexing: fetch: read body: %w", err)
	}
	if int64(len(data)) > defaultMaxBodyBytes {
		return nil, "", "", fmt.Errorf("knowledgeindexing: fetch: body exceeds %d bytes", defaultMaxBodyBytes)
	}
	return data, ct, finalURL, nil
}

func isBlockedHost(host string) bool {
	switch host {
	case "localhost", "metadata.google.internal", "metadata", "169.254.169.254":
		return true
	default:
		return strings.HasSuffix(host, ".internal")
	}
}

func forbidPrivateURLHost(ctx context.Context, host string) error {
	if host == "" {
		return fmt.Errorf("knowledgeindexing: fetch: empty host")
	}
	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil {
		return fmt.Errorf("knowledgeindexing: fetch: resolve %q: %w", host, err)
	}
	for _, ipa := range ips {
		ip := ipa.IP
		if ip == nil {
			continue
		}
		if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return fmt.Errorf("knowledgeindexing: fetch: resolved IP %s for %q is not allowed", ip, host)
		}
	}
	return nil
}
