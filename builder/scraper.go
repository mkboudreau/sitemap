package builder

import "github.com/PuerkitoBio/goquery"

type page struct {
	Links  []string
	Assets []string
}

func scrapePage(url string) (*page, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	page := new(page)
	page.Links = scrapeLinksFromDocument(doc)
	page.Assets = scrapeAssetsFromDocument(doc)
	return page, nil
}

func scrapeLinksFromDocument(doc *goquery.Document) []string {
	linkMap := make(map[string]bool)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr("href"); ok {
			linkMap[href] = true
		}
	})

	links := make([]string, len(linkMap))
	i := 0
	for k, _ := range linkMap {
		links[i] = k
		i++
	}

	return links
}

func scrapeAssetsFromDocument(doc *goquery.Document) []string {
	assetMap := make(map[string]bool)

	doc.Find("script").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok {
			assetMap[src] = true
		}
	})
	doc.Find("link").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("href"); ok {
			assetMap[src] = true
		}
	})
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		if src, ok := s.Attr("src"); ok {
			assetMap[src] = true
		}
	})

	assets := make([]string, len(assetMap))
	i := 0
	for k, _ := range assetMap {
		assets[i] = k
		i++
	}

	return assets
}
