package resolvers

type HAnime struct {
	URL    string
	Site   string
	Title  string
	Videos map[string]*Video
}

type Video struct {
	ID      string
	Quality string
	URL     string
	Title   string
	Size    int64
	Ext     string
}
