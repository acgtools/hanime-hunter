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

const defaultAniTitle = "unknown"

func init() {
	resolvers.Resolvers.Register("hanime1.me", New())
}

func New() resolvers.Resolver {
	return &resolver{}
}

var _ resolvers.Resolver = (*resolver)(nil)

type resolver struct{}

func (re *resolver) Resolve(u string, opt *resolvers.Option) ([]*resolvers.HAnime, error) {
	if strings.Contains(u, "playlist") {
		return resolvePlaylist(u)
	}

	site, vid, err := getSiteAndVID(u)
	if err != nil {
		return nil, fmt.Errorf("parse url %q: %w", u, err)
	}

	title, series, err := getAniInfo(u)
	if err != nil {
		return nil, fmt.Errorf("get anime info from %q: %w", u, err)
	}

	if title == defaultAniTitle {
		log.Warn("failed to get anime title")
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

func resolvePlaylist(u string) ([]*resolvers.HAnime, error) {
	doc, err := getHTMLPage(u)
	if err != nil {
		return nil, err
	}

	playlist := util.FindTagByNameAttrs(doc, "div", true, []html.Attribute{{Key: "id", Val: "home-rows-wrapper"}})
	if len(playlist) == 0 {
		return nil, fmt.Errorf("palylist in %q not found", u)
	}

	// hanime1me only contains one list div
	aTags := util.FindTagByNameAttrs(playlist[0], "a", true, []html.Attribute{{Key: "class", Val: "playlist-show-links"}})

	res := make([]*resolvers.HAnime, 0)
	for _, a := range aTags {
		href := util.GetAttrVal(a, "href")
		if strings.Contains(href, "watch") {
			site, vid, err := getSiteAndVID(href)
			if err != nil {
				return nil, err
			}

			title, _, err := getAniInfo(href)
			if err != nil {
				return nil, err
			}

			log.Infof("Anime found: %s, Searching episodes, Please wait a moment...", title)

			videos, eps, err := getDLInfo(vid)
			if err != nil {
				return nil, err
			}

			log.Infof("Episodes found: %#q", eps[0])
			// reduce request frequency to avoid rate limit
			time.Sleep(time.Duration(util.RandomInt63n(900, 3000)) * time.Millisecond) //nolint:gomnd

			res = append(res, &resolvers.HAnime{
				URL:    href,
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
		return "", "", fmt.Errorf("parse url %q: %w", u, err)
	}

	vid, ok := urlRes.Query()["v"]
	if !ok || len(vid) == 0 {
		return "", "", errors.New("vid not found")
	}
	return urlRes.Host, vid[0], nil
}

func getAniInfo(u string) (string, []string, error) {
	doc, err := getHTMLPage(u)
	if err != nil {
		return "", nil, fmt.Errorf("get ani page %q: %w", u, err)
	}

	series := util.FindTagByNameAttrs(doc, "div", true, []html.Attribute{{Key: "id", Val: "video-playlist-wrapper"}})
	if len(series) == 0 {
		return "", nil, fmt.Errorf("get series info from %q error", u)
	}
	seriesTag := series[0]

	title := defaultAniTitle
	titleTag := util.FindTagByNameAttrs(seriesTag, "h4", false, nil)
	if len(titleTag) > 0 {
		title = titleTag[0].FirstChild.Data
	}

	return title, getSeriesLinks(seriesTag), nil
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
	u := "https://hanime1.me/download?v=" + vid
	doc, err := getHTMLPage(u)
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
		title := util.GetAttrVal(a, "download")

		id := getID(link)
		quality := ""
		if tmp := strings.Split(id, "-"); len(tmp) > 1 { // the video id may not contain quality
			quality = tmp[1]
		}

		size, ext, err := getVideoInfo(link)
		if err != nil {
			return nil, nil, err
		}

		episodes = append(episodes, title)
		log.Debugf("Video found: %s - %s - %s", title, quality, ext)

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

func getVideoInfo(u string) (int64, string, error) {
	client := newClient()
	headers := map[string]string{
		"User-Agent": resolvers.UA,
	}

	resp, err := util.Get(client, u, headers)
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

func getHTMLPage(u string) (*html.Node, error) {
	client := newClient()
	headers := map[string]string{
		"User-Agent": resolvers.UA,
	}

	resp, err := util.Get(client, u, headers)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("parse html of %q: %w", u, err)
	}

	return doc, nil
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
