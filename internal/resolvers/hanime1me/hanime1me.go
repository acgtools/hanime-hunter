package hanime1me

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/acgtools/hanime-hunter/pkg/util"
	"golang.org/x/net/html"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func init() {
	resolvers.Resolvers.Register("hanime1.me", New())
}

func New() resolvers.Resolver {
	return &resolver{}
}

var _ resolvers.Resolver = (*resolver)(nil)

type resolver struct{}

func (re *resolver) Resolve(u string) ([]*resolvers.HAnime, error) {
	urlRes, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	vid, ok := urlRes.Query()["v"]
	if !ok {
		return nil, errors.New("video id not found")
	}

	videos, err := getDLInfo(vid[0])
	if err != nil {
		return nil, fmt.Errorf("get dl info: %w", err)
	}

	res := make([]*resolvers.HAnime, 0)

	res = append(res, &resolvers.HAnime{
		URL:    u,
		Site:   urlRes.Host,
		Videos: videos,
	})

	return res, nil
}

func getDLInfo(vid string) (map[string]*resolvers.Video, error) {
	doc, err := getDLPage(vid)
	if err != nil {
		return nil, fmt.Errorf("get download page: %w", err)
	}

	tables := util.FindTagByNameAttrs(doc, "table", true, []html.Attribute{{Key: "class", Val: "download-table"}})
	if len(tables) == 0 {
		return nil, errors.New("download info not found")
	}

	vidMap := make(map[string]*resolvers.Video)

	// hanime1.me only have one table in dl page
	aTags := util.FindTagByNameAttrs(tables[0], "a", false, nil)
	for _, a := range aTags {
		link := util.GetAttrVal(a, "href")
		id := getID(link)
		quality := strings.Split(id, "-")[1]
		title := util.GetAttrVal(a, "download")
		size, ext, err := getVideoInfo(link)
		if err != nil {
			return nil, err
		}

		vidMap[id] = &resolvers.Video{
			ID:      id,
			Quality: quality,
			URL:     link,
			Title:   title,
			Size:    size,
			Ext:     ext,
		}
	}

	return vidMap, nil
}

const ua = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36"

func getVideoInfo(link string) (int64, string, error) {
	client := newClient()

	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return 0, "", fmt.Errorf("create dl link request: %w", err)
	}
	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", fmt.Errorf("get dl info: %w", err)
	}
	defer resp.Body.Close()

	ext := strings.Split(resp.Header.Get("Content-Type"), "/")[1]
	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("get video size: %w", err)
	}

	return size, ext, nil
}

func getID(link string) string {
	r := regexp.MustCompile(`[^/]+-\d+p`)
	return r.FindString(link)
}

func getDLPage(vid string) (*html.Node, error) {
	u := "https://hanime1.me/download?v=" + vid
	client := newClient()

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send http request: %w", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	return doc, nil
}

func newClient() *http.Client {
	tlsConfig := http.DefaultTransport.(*http.Transport).TLSClientConfig

	c := &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout: 30 * time.Second,
			DisableKeepAlives:   false,

			Proxy: http.ProxyFromEnvironment,

			TLSClientConfig: &tls.Config{
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_AES_128_GCM_SHA256,
					tls.VersionTLS13,
					tls.VersionTLS10,
				},
			},
			DialTLSContext: func(_ context.Context, network, addr string) (net.Conn, error) {
				conn, err := tls.Dial(network, addr, tlsConfig)
				return conn, err
			},
		},
	}

	return c
}
