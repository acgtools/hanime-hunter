package resolvers

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/charmbracelet/log"
)

type Resolver interface {
	Resolve(u string, opt *Option) ([]*HAnime, error)
}

var Resolvers = newResolverMap()

type ResolverMap struct {
	m         sync.Mutex
	resolvers map[string]Resolver
}

type Option struct {
	Series   bool
	PlayList bool
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

func Resolve(u string, opt *Option) ([]*HAnime, error) {
	urlRes, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("resolve url: %w", err)
	}

	domain := urlRes.Host

	log.Infof("Site: %s", domain)

	resolver := Resolvers.resolvers[domain]

	return resolver.Resolve(u, opt) //nolint:wrapcheck
}
