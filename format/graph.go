package format

import (
	"fmt"
	"io"

	"github.com/mkboudreau/sitemap/domain"
)

type GraphSiteFormatter struct{}

func (f *GraphSiteFormatter) Format(w io.Writer, site *domain.Site) {
	fmt.Fprintf(w, "graph {\n")
	f.printSite(w, site)
	fmt.Fprintf(w, "}\n")
}
func (f *GraphSiteFormatter) printSite(w io.Writer, site *domain.Site) {
	var linkMap = make(map[string]bool)
	links := f.traverseSite(site)
	for _, link := range links {
		if !linkMap[link] {
			linkMap[link] = true
			fmt.Fprintf(w, link)
		}
	}
}
func (f *GraphSiteFormatter) traverseSite(site *domain.Site) []string {
	var links []string
	if site.Assets != nil && len(site.Assets) > 0 {
		for _, asset := range site.Assets {
			links = append(links, fmt.Sprintf("\t\t%q -- %q;\n", site.Url, asset))
		}
	}
	if site.Links != nil && len(site.Links) > 0 {
		for _, link := range site.Links {
			links = append(links, fmt.Sprintf("\t\t%q -- %q;\n", site.Url, link.Url))
		}
	}
	if site.Links != nil && len(site.Links) > 0 {
		for _, link := range site.Links {
			links = append(links, f.traverseSite(link)...)
		}
	}
	return links
}
