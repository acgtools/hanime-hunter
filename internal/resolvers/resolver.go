package resolvers

import (
	"fmt"
	"net/url"
	"sync"
)

type Resolver interface {
	Resolve(u string) ([]*HAnime, error)
}

var Resolvers = newResolverMap()

type ResolverMap struct {
	m         sync.Mutex
	resolvers map[string]Resolver
}

func newResolverMap() *ResolverMap {
	return &ResolverMap{
		m:         sync.Mutex{},
		resolvers: make(map[string]Resolver),
	}
}

func (r *ResolverMap) Register(domain string, resolver Resolver) {
	r.m.Lock()
	r.resolvers[domain] = resolver
	r.m.Unlock()
}

func Resolve(u string) ([]*HAnime, error) {
	urlRes, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("resovle url: %w", err)
	}

	domain := urlRes.Host

	resolver := Resolvers.resolvers[domain]

	return resolver.Resolve(u)
}
