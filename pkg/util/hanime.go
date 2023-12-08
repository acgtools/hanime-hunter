package util

import (
	"github.com/acgtools/hanime-hunter/internal/resolvers"
	"sort"
)

func SortAniVideos(videos map[string]*resolvers.Video) []*resolvers.Video {
	res := make([]*resolvers.Video, 0, len(videos))

	for _, v := range videos {
		res = append(res, v)
	}

	sort.SliceStable(res, func(i, j int) bool {
		return res[i].Size > res[j].Size
	})

	return res
}
