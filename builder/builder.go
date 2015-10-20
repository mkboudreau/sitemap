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

type linkResponse struct {
	site string
	link string
}

func NewSitemapBuilder(rate time.Duration, timeout time.Duration, workers int) *SitemapBuilder {
	return &SitemapBuilder{
		rate:             rate,
		timeout:          timeout,
		workers:          workers,
		interruptChannel: make(chan struct{}),
	}
}

func (builder *SitemapBuilder) Interrupt() {
	close(builder.interruptChannel)
}

func (builder *SitemapBuilder) Build(startingURL string) *domain.Sitemap {
	var workers sync.WaitGroup

	siteChannel := make(chan *domain.Site, builder.workers*5)
	responseChannel := make(chan *linkResponse, builder.workers)

	defer close(siteChannel)
	defer close(responseChannel)

	url := translateLink("", startingURL)
	top := domain.NewSite(url)
	sitemap := domain.NewSitemap(top)

	go builder.processResponses(sitemap, responseChannel, siteChannel)

	siteChannel <- top

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
					select {
					case <-builder.interruptChannel:
						break loop
					case <-time.After(remaining):
					}
				}
				last = time.Now() // reset
				select {
				case <-builder.interruptChannel:
					log.Println("Interrupt:", index)
					break loop
				case site := <-siteChannel:
					builder.procesSite(site, responseChannel)
				case <-time.After(builder.timeout):
					log.Println("Timeout:", index)
					break loop
				}
			}
		}(i)
	}

	workers.Wait()

	return sitemap
}

func (builder *SitemapBuilder) processResponses(sitemap *domain.Sitemap, approvedLink <-chan *linkResponse, siteChannel chan<- *domain.Site) {
	for {
		select {
		case <-builder.interruptChannel:
			log.Println("Interruptted")
			return
		case response, ok := <-approvedLink:
			if !ok {
				break
			}
			if linkedSite, isNew := sitemap.AddLink(response.site, response.link); isNew {
				go func() {
					select {
					case <-builder.interruptChannel:
					case siteChannel <- linkedSite:
					}
				}()
			}
		}
	}
}

func (builder *SitemapBuilder) procesSite(site *domain.Site, approvedLink chan<- *linkResponse) {
	log.Println("adding: ", site.Url)
	page, err := scrapePage(site.Url)
	if err != nil {
		log.Printf("could not retrieve links from url [%v]: %v", site.Url, err)
		return
	}

	for _, asset := range page.Assets {
		if translatedAsset := translateLink(site.Url, asset); translatedAsset != "" {
			site.Assets = append(site.Assets, translatedAsset)
		} else {
			fmt.Println("Ignoring asset", asset, "which translated to", translatedAsset)
		}
	}

	for _, link := range page.Links {
		if translatedLink := translateLink(site.Url, link); translatedLink != "" {
			if validateLinkInSameDomain(site.Url, link) {
				newLink := &linkResponse{site: site.Url, link: translatedLink}
				approvedLink <- newLink
			}
		}
	}

}
