package hanimetv

type Video struct {
	HentaiVideo     HentaiVideo     `json:"hentai_video"`
	HentaiFranchise HentaiFranchise `json:"hentai_franchise"`
	VideosManifest  VideosManifest  `json:"videos_manifest"`
}

type HentaiVideo struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type HentaiFranchise struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type VideosManifest struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	ID     int64  `json:"id"`
	Height string `json:"height"`
	Size   int64  `json:"filesize_mbs"`
	URL    string `json:"url"`
}
