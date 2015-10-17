package format

import (
	"fmt"
	"io"

	"github.com/mkboudreau/sitemap/domain"
)

type DigraphSiteFormatter struct{}

func (f *DigraphSiteFormatter) Format(w io.Writer, site *domain.Site) {
	fmt.Fprintf(w, "digraph {\n")
	f.traverseSite(w, site)
	fmt.Fprintf(w, "}\n")
}
func (f *DigraphSiteFormatter) traverseSite(w io.Writer, site *domain.Site) {
	fmt.Fprintf(w, "\tsubgraph cluster_0 {\n")
	fmt.Fprintf(w, "\t\tlabel=\"Assets\";\n")
	f.traverseAllSiteAssets(w, site)
	fmt.Fprintf(w, "\t}\n")
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "\tsubgraph cluster_1 {\n")
	fmt.Fprintf(w, "\t\tlabel=\"Links\";\n")
	f.traverseAllSiteLinks(w, site)
	fmt.Fprintf(w, "\t}\n")
}
func (f *DigraphSiteFormatter) traverseAllSiteAssets(w io.Writer, site *domain.Site) {
	if site.Assets != nil && len(site.Assets) > 0 {
		for _, asset := range site.Assets {
			fmt.Fprintf(w, "\t\t%q -> %q;\n", site.Url, asset)
		}
	}
	if site.Links != nil && len(site.Links) > 0 {
		for _, link := range site.Links {
			f.traverseAllSiteAssets(w, link)
		}
	}
}
func (f *DigraphSiteFormatter) traverseAllSiteLinks(w io.Writer, site *domain.Site) {
	if site.Links != nil && len(site.Links) > 0 {
		for _, link := range site.Links {
			fmt.Fprintf(w, "\t\t%q -> %q;\n", site.Url, link.Url)
		}
	}
	if site.Links != nil && len(site.Links) > 0 {
		for _, link := range site.Links {
			f.traverseAllSiteLinks(w, link)
		}
	}
}
