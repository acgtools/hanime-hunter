package util

import (
	"fmt"
	"net/http"
)

func Get(client *http.Client, u string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create req for %q: %w", u, err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send req to %q: %w", u, err)
	}

	return resp, nil
}
