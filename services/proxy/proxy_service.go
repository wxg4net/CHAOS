package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type proxyService struct {
}

func NewProxyService() Service {
	return &proxyService{}
}
func newReverseProxy(target string) (*httputil.ReverseProxy, error) {
	// 解析目标URL
	targetUrl, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	// 创建并返回反向代理
	return httputil.NewSingleHostReverseProxy(targetUrl), nil
}

func (u proxyService) ProxyUrl(proxyAddress string, proxyUrl string) error {
	if !strings.Contains(strings.ToLower(proxyUrl), "http") {
		proxyUrl = fmt.Sprintf("https://%s", proxyUrl)
	}

	proxy, err := newReverseProxy(proxyUrl)
	if err != nil {
		return err
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	if err := http.ListenAndServe(proxyAddress, nil); err != nil {
		return err
	}
	return nil
}
