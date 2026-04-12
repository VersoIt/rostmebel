package httpx

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultRequestTimeout      = 30 * time.Second
	defaultDialTimeout         = 10 * time.Second
	defaultKeepAlive           = 30 * time.Second
	defaultTLSHandshakeTimeout = 10 * time.Second
	defaultIdleConnTimeout     = 90 * time.Second
	defaultExpectContinue      = 1 * time.Second
	defaultResponseHeader      = 15 * time.Second

	defaultMaxIdleConns        = 100
	defaultMaxIdleConnsPerHost = 10
	defaultMaxConnsPerHost     = 50
)

type ProxyConfig struct {
	// Scheme: http, https, socks5
	Scheme string

	Host string
	Port string

	Username string
	Password string
}

func (p ProxyConfig) Enabled() bool {
	return strings.TrimSpace(p.Host) != "" && strings.TrimSpace(p.Port) != ""
}

func (p ProxyConfig) URL() (*url.URL, error) {
	if !p.Enabled() {
		return nil, nil
	}

	scheme := strings.TrimSpace(strings.ToLower(p.Scheme))
	if scheme == "" {
		scheme = "http"
	}

	switch scheme {
	case "http", "https", "socks5", "socks5h":
	default:
		return nil, fmt.Errorf("unsupported proxy scheme %q", p.Scheme)
	}

	u := &url.URL{
		Scheme: scheme,
		Host:   net.JoinHostPort(strings.TrimSpace(p.Host), strings.TrimSpace(p.Port)),
	}

	if p.Username != "" {
		if p.Password != "" {
			u.User = url.UserPassword(p.Username, p.Password)
		} else {
			u.User = url.User(p.Username)
		}
	}

	return u, nil
}

type ClientOptions struct {
	Timeout               time.Duration
	DialTimeout           time.Duration
	KeepAlive             time.Duration
	ResponseHeaderTimeout time.Duration
	MaxIdleConns          int
	MaxIdleConnsPerHost   int
	MaxConnsPerHost       int
	IdleConnTimeout       time.Duration
	TLSHandshakeTimeout   time.Duration
	ExpectContinueTimeout time.Duration

	Proxy *ProxyConfig

	// Обычно оставляем false.
	InsecureSkipVerify bool
}

func NewHTTPClient(opts ClientOptions) (*http.Client, error) {
	if opts.Timeout <= 0 {
		opts.Timeout = defaultRequestTimeout
	}
	if opts.DialTimeout <= 0 {
		opts.DialTimeout = defaultDialTimeout
	}
	if opts.KeepAlive <= 0 {
		opts.KeepAlive = defaultKeepAlive
	}
	if opts.ResponseHeaderTimeout <= 0 {
		opts.ResponseHeaderTimeout = defaultResponseHeader
	}
	if opts.MaxIdleConns <= 0 {
		opts.MaxIdleConns = defaultMaxIdleConns
	}
	if opts.MaxIdleConnsPerHost <= 0 {
		opts.MaxIdleConnsPerHost = defaultMaxIdleConnsPerHost
	}
	if opts.MaxConnsPerHost <= 0 {
		opts.MaxConnsPerHost = defaultMaxConnsPerHost
	}
	if opts.IdleConnTimeout <= 0 {
		opts.IdleConnTimeout = defaultIdleConnTimeout
	}
	if opts.TLSHandshakeTimeout <= 0 {
		opts.TLSHandshakeTimeout = defaultTLSHandshakeTimeout
	}
	if opts.ExpectContinueTimeout <= 0 {
		opts.ExpectContinueTimeout = defaultExpectContinue
	}

	dialer := &net.Dialer{
		Timeout:   opts.DialTimeout,
		KeepAlive: opts.KeepAlive,
	}

	transport := &http.Transport{
		Proxy:                 nil,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          opts.MaxIdleConns,
		MaxIdleConnsPerHost:   opts.MaxIdleConnsPerHost,
		MaxConnsPerHost:       opts.MaxConnsPerHost,
		IdleConnTimeout:       opts.IdleConnTimeout,
		TLSHandshakeTimeout:   opts.TLSHandshakeTimeout,
		ExpectContinueTimeout: opts.ExpectContinueTimeout,
		ResponseHeaderTimeout: opts.ResponseHeaderTimeout,
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: opts.InsecureSkipVerify,
		},
	}

	if opts.Proxy != nil && opts.Proxy.Enabled() {
		proxyURL, err := opts.Proxy.URL()
		if err != nil {
			return nil, fmt.Errorf("build proxy url: %w", err)
		}

		switch proxyURL.Scheme {
		case "http", "https":
			transport.Proxy = http.ProxyURL(proxyURL)
		case "socks5", "socks5h":
			socksDialer, err := newSOCKS5Dialer(proxyURL, dialer)
			if err != nil {
				return nil, fmt.Errorf("build socks5 dialer: %w", err)
			}
			transport.Proxy = nil
			transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
				type dialContext interface {
					DialContext(ctx context.Context, network, addr string) (net.Conn, error)
				}
				if dc, ok := socksDialer.(dialContext); ok {
					return dc.DialContext(ctx, network, addr)
				}
				return nil, errors.New("socks5 dialer does not support context")
			}
		}
	}

	return &http.Client{
		Timeout:   opts.Timeout,
		Transport: transport,
	}, nil
}

func CloseIdleConnections(client *http.Client) {
	if client == nil {
		return
	}
	if tr, ok := client.Transport.(interface{ CloseIdleConnections() }); ok {
		tr.CloseIdleConnections()
	}
}

func ReadBodySnippet(r io.Reader, maxBytes int64) string {
	if r == nil || maxBytes <= 0 {
		return ""
	}

	data, err := io.ReadAll(io.LimitReader(r, maxBytes))
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}
