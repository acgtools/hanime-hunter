package resolvers

import "sort"

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
	IsM3U8  bool
	Title   string
	Size    int64
	Ext     string
}

func SortAniVideos(videos map[string]*Video, asc bool) []*Video {
	res := make([]*Video, 0, len(videos))

	for _, v := range videos {
		res = append(res, v)
	}

	if asc {
		sort.SliceStable(res, func(i, j int) bool {
			return res[i].Size < res[j].Size
		})
	} else {
		sort.SliceStable(res, func(i, j int) bool {
			return res[i].Size > res[j].Size
		})
	}

	return res
}
