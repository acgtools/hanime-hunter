package hanime1me

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/acgtools/hanime-hunter/pkg/util"
	"github.com/charmbracelet/log"
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

func (re *resolver) Resolve(u string, opt *resolvers.Option) ([]*resolvers.HAnime, error) {
	site, vid, err := getSiteAndVID(u)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	title, series, err := getAniInfo(u)
	if err != nil {
		return nil, fmt.Errorf("get ani info: %w", err)
	}

	log.Info("Anime Info", "title", title)

	res := make([]*resolvers.HAnime, 0)

	if !opt.Series {
		videos, err := getDLInfo(vid)
		if err != nil {
			return nil, fmt.Errorf("get download info: %w", err)
		}

		res = append(res, &resolvers.HAnime{
			URL:    u,
			Site:   site,
			Title:  title,
			Videos: videos,
		})
	} else {
		for _, s := range series {
			_, vID, _ := getSiteAndVID(s) // no need to check err
			videos, err := getDLInfo(vID)
			if err != nil {
				return nil, fmt.Errorf("get download info: %w", err)
			}

			res = append(res, &resolvers.HAnime{
				URL:    s,
				Site:   site,
				Title:  title,
				Videos: videos,
			})
		}
	}

	return res, nil
}

func getSiteAndVID(u string) (string, string, error) {
	urlRes, err := url.Parse(u)
	if err != nil {
		return "", "", fmt.Errorf("parse url: %w", err)
	}

	vid, ok := urlRes.Query()["v"]
	if !ok || len(vid) == 0 {
		return "", "", errors.New("vid not found")
	}
	return urlRes.Host, vid[0], nil
}

func getAniInfo(u string) (string, []string, error) {
	doc, err := getAniPage(u)
	if err != nil {
		return "", nil, fmt.Errorf("get ani page: %w", err)
	}

	series := util.FindTagByNameAttrs(doc, "div", true, []html.Attribute{{Key: "id", Val: "video-playlist-wrapper"}})
	seriesTag := series[0]

	titleTag := util.FindTagByNameAttrs(seriesTag, "h4", false, nil)
	title := titleTag[0].FirstChild.Data

	return title, getSeriesLinks(seriesTag), nil
}

func getAniPage(u string) (*html.Node, error) {
	resp, err := request(http.MethodGet, u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("get anime info page: %w", err)
	}

	return doc, nil
}

func getSeriesLinks(node *html.Node) []string {
	list := util.FindTagByNameAttrs(node, "div", true, []html.Attribute{{Key: "id", Val: "playlist-scroll"}})

	aTags := util.FindTagByNameAttrs(list[0], "a", false, nil)

	links := make([]string, 0)
	for _, a := range aTags {
		href := util.GetAttrVal(a, "href")
		if strings.Contains(href, "watch") {
			links = append(links, href)
		}
	}

	return links
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

func getVideoInfo(u string) (int64, string, error) {
	resp, err := request(http.MethodGet, u)
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

	resp, err := request(http.MethodGet, u)
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html: %w", err)
	}

	return doc, nil
}

func request(method string, u string) (*http.Response, error) {
	client := newClient()

	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send http request: %w", err)
	}

	return resp, nil
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
