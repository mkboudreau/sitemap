package domain

import (
	"fmt"
	"sync"
)

type Site struct {
	Url    string   `json:"url"`
	Assets []string `json:"assets,omitempty"`
	Links  []*Site  `json:"links,omitempty"`

	allLinksForSite map[string]bool `json:"-"`
	mutex           *sync.RWMutex   `json:"-"`
}

func (s *Site) String() string {
	return fmt.Sprintf("\nURL=%v; Assets=%v; Links=%v", s.Url, s.Assets, s.Links)
}

func (s *Site) copyAndFlattenSite() *Site {
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
	return &Site{Url: url, mutex: &sync.RWMutex{}, allLinksForSite: make(map[string]bool)}
}

func (s *Site) AddLink(newSite *Site) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.allLinksForSite[newSite.Url] {
		s.allLinksForSite[newSite.Url] = true
	}

	s.Links = append(s.Links, newSite)
}

type Sitemap struct {
	Top      *Site
	allsites map[string]*Site
	mutex    *sync.RWMutex
}

func NewSitemap(top *Site) *Sitemap {
	sitemap := &Sitemap{Top: top, allsites: make(map[string]*Site), mutex: &sync.RWMutex{}}
	sitemap.allsites[top.Url] = top
	return sitemap
}

func (s *Sitemap) AddLink(parentURL string, childURL string) (*Site, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var isNew bool

	parent := s.allsites[parentURL]
	if parent == nil {
		return nil, false
	}

	child := s.allsites[childURL]
	if child != nil {
		child = child.copyAndFlattenSite()
	} else {
		child = NewSite(childURL)
		s.allsites[childURL] = child
		isNew = true
	}

	parent.AddLink(child)

	return child, isNew
}
