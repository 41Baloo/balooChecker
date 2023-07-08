package proxy

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Returns true if response is correct
func ValidateResponse(resp []byte) bool {
	return (string(resp)[41:55] == "Example Domain")
}

func ConnectHTTP(proxy string, timeout time.Duration) ([]byte, error) {
	proxyURL, err := url.Parse("http://" + proxy)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: timeout,
	}

	response, err := client.Get("https://example.com/")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ConnectHTTPS(proxy string, timeout time.Duration) ([]byte, error) {
	proxyURL, err := url.Parse("https://" + proxy)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: timeout,
	}

	response, err := client.Get("https://example.com/")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func ConnectSOCKS5(proxy string, timeout time.Duration) ([]byte, error) {
	proxyURL, err := url.Parse("socks5://" + proxy)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		Timeout: timeout,
	}

	response, err := client.Get("https://example.com/")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}