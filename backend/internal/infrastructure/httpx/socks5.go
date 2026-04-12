package httpx

import (
	"fmt"
	"net"
	"net/url"

	"golang.org/x/net/proxy"
)

func newSOCKS5Dialer(proxyURL *url.URL, forwardDialer *net.Dialer) (proxy.Dialer, error) {
	if proxyURL == nil {
		return nil, fmt.Errorf("proxy url is nil")
	}

	var auth *proxy.Auth
	if proxyURL.User != nil {
		password, _ := proxyURL.User.Password()
		auth = &proxy.Auth{
			User:     proxyURL.User.Username(),
			Password: password,
		}
	}

	return proxy.SOCKS5("tcp", proxyURL.Host, auth, forwardDialer)
}
