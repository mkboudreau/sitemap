package builder

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mkboudreau/sitemap/domain"
)

type SitemapBuilder struct {
	rate             time.Duration
	timeout          time.Duration
	workers          int
	interruptChannel chan struct{}
}

func NewSitemapBuilder(rate time.Duration, timeout time.Duration, workers int) *SitemapBuilder {
	return &SitemapBuilder{
		rate:             rate,
		timeout:          timeout,
		workers:          workers,
		interruptChannel: make(chan struct{}),
	}
}

func (builder *SitemapBuilder) interrupt() {
	close(builder.interruptChannel)
}

func (builder *SitemapBuilder) Build(startingURL string) *domain.Sitemap {
	var workers sync.WaitGroup

	siteChannel := make(chan *domain.Site, builder.workers*5)
	doneChannel := make(chan struct{})

	url := translateLink("", startingURL)
	top := domain.NewSite(url)
	sitemap := domain.NewSitemap(top)

	go addSiteToSitemap(top, sitemap, siteChannel)

	for i := 0; i < builder.workers; i++ {
		workers.Add(1)
		go func(index int) {
			defer workers.Done()

			last := time.Date(1980, time.January, 1, 1, 1, 1, 1, time.UTC)
		loop:
			for {
				elapsed := time.Since(last)
				remaining := builder.rate - elapsed
				if remaining > 0 {
					<-time.After(remaining)
				}
				last = time.Now() // reset
				select {
				case site := <-siteChannel:
					addSiteToSitemap(site, sitemap, siteChannel)
				case <-time.After(builder.timeout):
					log.Println("Timeout:", index)
					break loop
				case <-builder.interruptChannel:
					log.Println("Interrupt:", index)
					break loop
				}
			}
		}(i)
	}

	workers.Wait()
	close(doneChannel)

	return sitemap
}

func addSiteToSitemap(site *domain.Site, sitemap *domain.Sitemap, siteChannel chan<- *domain.Site) {
	log.Println("adding: ", site.Url)
	page, err := scrapePage(site.Url)
	if err != nil {
		log.Printf("could not retrieve links from url [%v]: %v", site.Url, err)
		return
	}

	allLinksForSite := make(map[string]bool)
	for _, link := range page.Links {
		if translatedLink := translateLink(site.Url, link); translatedLink != "" {
			if validateLinkInSameDomain(site.Url, link) {
				linkedSite, isNew := sitemap.AddUrl(translatedLink)
				/*
					if isNew {
						site.AddLink(linkedSite)
					} else {
						site.AddLink(linkedSite.CopyAndFlattenSite())
					}
				*/
				if isNew {
					// first time seen for entire site
					allLinksForSite[linkedSite.Url] = true
					site.Links = append(site.Links, linkedSite)
					go func() {
						siteChannel <- linkedSite
					}()
				} else if !allLinksForSite[linkedSite.Url] {
					// already processed somewhere else in the site
					allLinksForSite[linkedSite.Url] = true
					site.Links = append(site.Links, linkedSite.CopyAndFlattenSite())
				}
			}
		}
	}

	for _, asset := range page.Assets {
		if translatedAsset := translateLink(site.Url, asset); translatedAsset != "" {
			site.Assets = append(site.Assets, translatedAsset)
		} else {
			fmt.Println("Ignoring asset", asset, "which translated to", translatedAsset)
		}
	}
}
