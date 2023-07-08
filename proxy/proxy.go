package proxy

const (
	PROXY_HTTP   = "HTTP"
	PROXY_HTTPS  = "HTTPS"
	PROXY_SOCKS4 = "SOCKS4"
	PROXY_SOCKS5 = "SOCKS5"
)

type PROXY struct {
	IP   string
	Type string
}
