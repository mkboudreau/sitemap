package format

import (
	"fmt"
	"io"

	"github.com/mkboudreau/sitemap/domain"
)

type DigraphSiteFormatter struct{}

func (f *DigraphSiteFormatter) Format(w io.Writer, site *domain.Site) {
	fmt.Fprintf(w, "digraph {\n")
	f.printSite(w, site)
	fmt.Fprintf(w, "}\n")
}
func (f *DigraphSiteFormatter) printSite(w io.Writer, site *domain.Site) {
	fmt.Fprintf(w, "\tsubgraph cluster_0 {\n")
	fmt.Fprintf(w, "\t\tlabel=\"Assets\";\n")
	f.printAllSiteAssets(w, site)
	fmt.Fprintf(w, "\t}\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tsubgraph cluster_1 {\n")
	fmt.Fprintf(w, "\t\tlabel=\"Links\";\n")
	f.printAllSiteLinks(w, site)
	fmt.Fprintf(w, "\t}\n")
}
func (f *DigraphSiteFormatter) printAllSiteAssets(w io.Writer, site *domain.Site) {
	var assetMap = make(map[string]bool)
	assets := f.traverseAllSiteAssets(site)
	for _, asset := range assets {
		if !assetMap[asset] {
			assetMap[asset] = true
			fmt.Fprintf(w, asset)
		}
	}
}
func (f *DigraphSiteFormatter) traverseAllSiteAssets(site *domain.Site) []string {
	var assets []string
	if site.Assets != nil && len(site.Assets) > 0 {
		for _, asset := range site.Assets {
			assets = append(assets, fmt.Sprintf("\t\t%q -> %q;\n", site.Url, asset))
		}
	}
	if site.Links != nil && len(site.Links) > 0 {
		for _, link := range site.Links {
			assets = append(assets, f.traverseAllSiteAssets(link)...)
		}
	}
	return assets
}
func (f *DigraphSiteFormatter) printAllSiteLinks(w io.Writer, site *domain.Site) {
	var linkMap = make(map[string]bool)
	links := f.traverseAllSiteLinks(site)
	for _, link := range links {
		if !linkMap[link] {
			linkMap[link] = true
			fmt.Fprintf(w, link)
		}
	}
}
func (f *DigraphSiteFormatter) traverseAllSiteLinks(site *domain.Site) []string {
	var links []string
	if site.Links != nil && len(site.Links) > 0 {
		for _, link := range site.Links {
			links = append(links, fmt.Sprintf("\t\t%q -> %q;\n", site.Url, link.Url))
		}
	}
	if site.Links != nil && len(site.Links) > 0 {
		for _, link := range site.Links {
			links = append(links, f.traverseAllSiteLinks(link)...)
		}
	}
	return links
}
