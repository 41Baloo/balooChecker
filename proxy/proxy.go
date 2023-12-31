package proxy

const (
	PROXY_HTTP   = "HTTP"
	PROXY_HTTPS  = "HTTPS"
	PROXY_SOCKS4 = "SOCKS4"
	PROXY_SOCKS5 = "SOCKS5"
)

var (
	HTTP_FOUND   int32
	HTTPS_FOUND  int32
	SOCKS4_FOUND int32
	SOCKS5_FOUND int32

	FAKE_FOUND int32
)
