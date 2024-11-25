package proxy

type Service interface {
	ProxyUrl(proxyAddress string, proxyUrl string) error
}
