package hanimetv

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"github.com/acgtools/hanime-hunter/pkg/util"
	"github.com/charmbracelet/log"
	"golang.org/x/net/html"
)

func init() {
	resolvers.Resolvers.Register("hanime.tv", New())
}

func New() resolvers.Resolver {
	return &resolver{}
}

var _ resolvers.Resolver = (*resolver)(nil)

type resolver struct{}

const videoAPIURL = "https://hanime.tv/api/v8/video?id="

func (re *resolver) Resolve(u string, opt *resolvers.Option) ([]*resolvers.HAnime, error) {
	if strings.Contains(u, "playlists") {
		return resolvePlaylist(u)
	}

	urlRes, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", u, err)
	}
	site := urlRes.Host

	slug, err := getVideoID(urlRes.Path)
	if err != nil {
		return nil, err
	}

	v, err := getVideoInfo(slug)
	if err != nil {
		return nil, err
	}

	log.Infof("Anime found: %s, Searching episodes, Please wait a moment...", v.HentaiFranchise.Title)

	res := make([]*resolvers.HAnime, 0)
	episodes := make([]string, 0)

	if !opt.Series {
		vidMap, eps := getVidMap(v)

		episodes = append(episodes, eps[0])
		log.Infof("Episodes found: %#q", episodes)

		res = append(res, &resolvers.HAnime{
			URL:    u,
			Site:   site,
			Title:  v.HentaiFranchise.Title,
			Videos: vidMap,
		})

		return res, nil
	}

	for _, fv := range v.HentaiFranchiseHentaiVideos {
		video, err := getVideoInfo(fv.Slug)
		if err != nil {
			return nil, err
		}

		vidMap, eps := getVidMap(video)
		episodes = append(episodes, eps[0])
		res = append(res, &resolvers.HAnime{
			Site:   site,
			Title:  video.HentaiFranchise.Slug,
			Videos: vidMap,
		})
	}

	log.Infof("Episodes found: %#q", episodes)

	return res, nil
}

func resolvePlaylist(u string) ([]*resolvers.HAnime, error) {
	slugs, err := getPlaylistSlugs(u)
	if err != nil {
		return nil, err
	}

	res := make([]*resolvers.HAnime, 0)

	for _, s := range slugs {
		v, err := getVideoInfo(s)
		if err != nil {
			return nil, err
		}

		log.Infof("Anime found: %s, Searching episodes, Please wait a moment...", v.HentaiFranchise.Title)

		vidMap, eps := getVidMap(v)

		log.Infof("Episodes found: %#q", eps[0])

		res = append(res, &resolvers.HAnime{
			Title:  v.HentaiFranchise.Title,
			Videos: vidMap,
		})

		time.Sleep(time.Duration(util.RandomInt63n(900, 3000)) * time.Millisecond) //nolint:gomnd
	}

	return res, nil
}

func getPlaylistSlugs(u string) ([]string, error) {
	doc, err := util.GetHTMLPage(NewClient(), u, map[string]string{"User-Agent": resolvers.UA})
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	listDivs := util.FindTagByNameAttrs(doc, "div", true, []html.Attribute{{Key: "class", Val: "playlists__panel panel__content"}})
	if len(listDivs) == 0 {
		return nil, fmt.Errorf("playlist not found in %q", u)
	}

	aTags := util.FindTagByNameAttrs(listDivs[0], "a", true, []html.Attribute{{Key: "class", Val: "flex row"}})
	if len(aTags) == 0 {
		return nil, fmt.Errorf("anime not found in %q", u)
	}

	res := make([]string, 0)
	for _, a := range aTags {
		href := util.GetAttrVal(a, "href")
		urlRes, _ := url.Parse(href)
		path := urlRes.Path
		if !strings.HasPrefix(path, "/videos/hentai/") {
			continue
		}
		res = append(res, strings.TrimPrefix(path, "/videos/hentai/"))
	}

	return res, nil
}

func getVidMap(v *Video) (map[string]*resolvers.Video, []string) {
	vidMap := make(map[string]*resolvers.Video)
	eps := make([]string, 0)

	for _, s := range v.VideosManifest.Servers[0].Streams {
		if s.Height == "1080" {
			continue
		}
		quality := s.Height + "p"

		eps = append(eps, v.HentaiVideo.Slug)
		log.Debugf("video %s %s found", v.HentaiVideo.Slug, quality)

		vidMap[quality] = &resolvers.Video{
			ID:      strconv.FormatInt(s.ID, 10),
			Quality: quality,
			URL:     s.URL,
			IsM3U8:  true,
			Title:   v.HentaiVideo.Slug,
			Size:    s.Size,
			Ext:     "mp4",
		}
	}

	return vidMap, eps
}

func getVideoID(path string) (string, error) {
	if !strings.HasPrefix(path, "/videos/hentai/") {
		return "", fmt.Errorf("video ID not found in %q", path)
	}

	params := strings.Split(path, "/")
	if len(params) != 4 { //nolint:gomnd
		return "", fmt.Errorf("video ID not found in %q", path)
	}

	return params[3], nil
}

func getVideoInfo(slug string) (*Video, error) {
	resp, err := util.Get(NewClient(), fmt.Sprintf("%s%s", videoAPIURL, slug), map[string]string{"User-Agent": resolvers.UA})
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response from anime %q: %w", slug, err)
	}

	v := &Video{}
	if err := json.Unmarshal(data, v); err != nil {
		return nil, fmt.Errorf("parse response json from anime %q : %w", slug, err)
	}

	return v, nil
}

func NewClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout: 30 * time.Second, //nolint:gomnd
			Proxy:               http.ProxyFromEnvironment,
		},
		Timeout: 15 * time.Minute, //nolint:gomnd
	}
}
