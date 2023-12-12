package util

import (
	"fmt"
	"net/http"

	"golang.org/x/net/html"
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

func GetHTMLPage(client *http.Client, u string, headers map[string]string) (*html.Node, error) {
	resp, err := Get(client, u, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html of %q: %w", u, err)
	}

	return doc, nil
}
