package request

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

var DefaultHeaders = map[string]string{
	"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
	"Accept-Charset":  "UTF-8,*;q=0.5",
	"Accept-Encoding": "gzip,deflate,sdch",
	"Accept-Language": "en-US,en;q=0.8",
	"User-Agent":      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.81 Safari/537.36",
}

// Request sent http request to download anime with fake UA
func Request(method, url string, headers map[string]string) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			DisableCompression:  true,
			TLSHandshakeTimeout: 10 * time.Second,                      //nolint:gomnd
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}, //nolint:gosec
		},
		//Timeout: 15 * time.Minute, //nolint:gomnd
		Timeout: 3 * time.Second, //nolint:gomnd
	}

	req, err := http.NewRequest(method, url, nil) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	for k, v := range DefaultHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get http response: %w", err)
	}

	return resp, nil
}
