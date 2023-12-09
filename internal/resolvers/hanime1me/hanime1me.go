package hanime1me

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/acgtools/hanime-hunter/pkg/util"
	"github.com/charmbracelet/log"
	"golang.org/x/net/html"
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
		return nil, fmt.Errorf("parse url %q: %w", u, err)
	}

	title, series, err := getAniInfo(u)
	if err != nil {
		return nil, fmt.Errorf("get anime info from %q: %w", u, err)
	}

	log.Infof("Anime found: %s, Searching episodes, Please wait a moment...", title)

	res := make([]*resolvers.HAnime, 0)

	if !opt.Series {
		videos, eps, err := getDLInfo(vid)
		if err != nil {
			return nil, fmt.Errorf("get download info of %q: %w", vid, err)
		}

		log.Infof("Episodes found: %q", eps[0])

		res = append(res, &resolvers.HAnime{
			URL:    u,
			Site:   site,
			Title:  title,
			Videos: videos,
		})
	} else {
		titles := make([]string, 0)

		for _, s := range series {
			_, vID, _ := getSiteAndVID(s) // no need to check err
			videos, eps, err := getDLInfo(vID)
			if err != nil {
				return nil, fmt.Errorf("get download info of %q: %w", vID, err)
			}

			titles = append(titles, eps[0])

			res = append(res, &resolvers.HAnime{
				URL:    s,
				Site:   site,
				Title:  title,
				Videos: videos,
			})
		}

		log.Infof("Episodes found %#q", titles)
	}

	return res, nil
}

func getSiteAndVID(u string) (string, string, error) {
	urlRes, err := url.Parse(u)
	if err != nil {
		return "", "", fmt.Errorf("parse url %q: %w", u, err)
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
		return "", nil, fmt.Errorf("get ani page %q: %w", u, err)
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
		return nil, fmt.Errorf("parse anime info page %q to HTML doc: %w", u, err)
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

func getDLInfo(vid string) (map[string]*resolvers.Video, []string, error) {
	doc, err := getDLPage(vid)
	if err != nil {
		return nil, nil, fmt.Errorf("get download page: %w", err)
	}

	tables := util.FindTagByNameAttrs(doc, "table", true, []html.Attribute{{Key: "class", Val: "download-table"}})
	if len(tables) == 0 {
		return nil, nil, errors.New("download info not found")
	}

	vidMap := make(map[string]*resolvers.Video)
	episodes := make([]string, 0)

	// hanime1.me only have one table in dl page
	aTags := util.FindTagByNameAttrs(tables[0], "a", false, nil)
	for _, a := range aTags {
		link := util.GetAttrVal(a, "href")
		id := getID(link)
		quality := strings.Split(id, "-")[1]
		title := util.GetAttrVal(a, "download")
		size, ext, err := getVideoInfo(link)
		if err != nil {
			return nil, nil, err
		}

		episodes = append(episodes, title)
		log.Debugf("Episode found: %s - %s - %s", title, quality, ext)

		vidMap[quality] = &resolvers.Video{
			ID:      id,
			Quality: quality,
			URL:     link,
			Title:   title,
			Size:    size,
			Ext:     ext,
		}
	}

	return vidMap, episodes, nil
}

const ua = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36"

func getVideoInfo(u string) (int64, string, error) {
	resp, err := request(http.MethodGet, u)
	if err != nil {
		return 0, "", fmt.Errorf("get dl info from %q: %w", u, err)
	}
	defer resp.Body.Close()

	ext := strings.Split(resp.Header.Get("Content-Type"), "/")[1]
	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("get video size from %q: %w", u, err)
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

func request(method string, u string) (*http.Response, error) {
	client := newClient()

	req, err := http.NewRequest(method, u, nil) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("create http request for %q: %w", u, err)
	}

	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send http request to %q: %w", u, err)
	}

	return resp, nil
}

func newClient() *http.Client {
	tlsConfig := http.DefaultTransport.(*http.Transport).TLSClientConfig //nolint:forcetypeassert

	c := &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout: 30 * time.Second, //nolint:gomnd
			DisableKeepAlives:   false,

			Proxy: http.ProxyFromEnvironment,

			TLSClientConfig: &tls.Config{
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_AES_128_GCM_SHA256,
					tls.VersionTLS13, //nolint:gosec
					tls.VersionTLS10,
				},
			},
			DialTLSContext: func(_ context.Context, network, addr string) (net.Conn, error) {
				conn, err := tls.Dial(network, addr, tlsConfig)
				return conn, err //nolint:wrapcheck
			},
		},
	}

	return c
}
