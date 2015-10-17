package domain

import (
	"fmt"
	"sync"
)

type Site struct {
	Url    string        `json:"url"`
	Assets []string      `json:"assets,omitempty"`
	Links  []*Site       `json:"links,omitempty"`
	mutex  *sync.RWMutex `json:"-"`
}

func (s *Site) String() string {
	return fmt.Sprintf("\nURL=%v; Assets=%v; Links=%v", s.Url, s.Assets, s.Links)
}

func (s *Site) CopyAndFlattenSite() *Site {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	newSite := NewSite(s.Url)
	newSite.Assets = s.Assets
	var links []*Site
	for _, linkedSite := range s.Links {
		links = append(links, NewSite(linkedSite.Url))
	}
	newSite.Links = links
	return newSite
}

func NewSite(url string) *Site {
	return &Site{Url: url, mutex: &sync.RWMutex{}}
}
func (s *Site) AddLink(newSite *Site) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	link := s.Links[0]
	if link != nil {
		return false
	}
	s.Links = append(s.Links, newSite)
	return true
}

type Sitemap struct {
	Top      *Site
	allsites map[string]*Site
	mutex    *sync.RWMutex
}

func NewSitemap(top *Site) *Sitemap {
	return &Sitemap{Top: top, allsites: make(map[string]*Site), mutex: &sync.RWMutex{}}
}

func (sitemap *Sitemap) AddUrl(url string) (*Site, bool) {
	sitemap.mutex.Lock()
	defer sitemap.mutex.Unlock()

	site := sitemap.allsites[url]
	if site != nil {
		return site, false
	}
	site = NewSite(url)
	sitemap.allsites[url] = site

	return site, true
}
