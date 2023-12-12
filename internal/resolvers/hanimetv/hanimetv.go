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

func (re *resolver) Resolve(u string, _ *resolvers.Option) ([]*resolvers.HAnime, error) {
	urlRes, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", u, err)
	}
	site := urlRes.Host

	vid, err := getVideoID(urlRes.Path)
	if err != nil {
		return nil, err
	}

	aniTitle, video, err := getVideoInfo(vid)
	if err != nil {
		return nil, err
	}

	res := make([]*resolvers.HAnime, 0)

	res = append(res, &resolvers.HAnime{
		URL:    u,
		Site:   site,
		Title:  aniTitle,
		Videos: video,
	})

	return res, nil
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

func getVideoInfo(slug string) (string, map[string]*resolvers.Video, error) {
	resp, err := request(http.MethodGet, fmt.Sprintf("%s%s", videoAPIURL, slug))
	if err != nil {
		return "", nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("read response from anime %q: %w", slug, err)
	}

	v := &Video{}
	if err := json.Unmarshal(data, v); err != nil {
		return "", nil, fmt.Errorf("parse response json from anime %q : %w", slug, err)
	}

	vidMap := make(map[string]*resolvers.Video)

	for _, s := range v.VideosManifest.Servers[0].Streams {
		if s.Height == "1080" {
			continue
		}
		quality := s.Height + "p"

		vidMap[quality] = &resolvers.Video{
			ID:      strconv.FormatInt(s.ID, 10),
			Quality: quality,
			URL:     s.URL,
			Title:   v.HentaiVideo.Name,
			Size:    s.Size,
			Ext:     "mp4",
		}
	}

	return v.HentaiFranchise.Name, vidMap, nil
}

func request(method string, u string) (*http.Response, error) {
	client := newClient()

	req, err := http.NewRequest(method, u, nil) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("create http request for %q: %w", u, err)
	}

	const ua = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36"

	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send http request to %q: %w", u, err)
	}

	return resp, nil
}

func newClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSHandshakeTimeout: 30 * time.Second, //nolint:gomnd
			Proxy:               http.ProxyFromEnvironment,
		},
	}
}
