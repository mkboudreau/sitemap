package format

import (
	"fmt"
	"io"
	"strings"

	"github.com/mkboudreau/sitemap/domain"
)

type TabbedSiteFormatter struct{}

func (f *TabbedSiteFormatter) Format(w io.Writer, site *domain.Site) {
	f.traverse(w, site, 0)
}

func (f *TabbedSiteFormatter) traverse(w io.Writer, site *domain.Site, depth int) {
	tabs := getTabCount(depth)
	fmt.Fprintf(w, "%vURL:   \t%v\n", tabs, site.Url)

	if site.Assets == nil || len(site.Assets) == 0 {
		//fmt.Fprintf(w, "%vAssets: <none>\n", tabs)
	} else {
		fmt.Fprintf(w, "%vAssets:\t%v\n", tabs, strings.Join(site.Assets, ", "))
	}

	if site.Links == nil || len(site.Links) == 0 {
		//fmt.Fprintf(w, "%vLinks: <none>\n", tabs)
	} else {
		fmt.Fprintf(w, "%vLinks: \n", tabs)
		for _, link := range site.Links {
			f.traverse(w, link, (depth + 1))
		}
		fmt.Fprintln(w, "")
	}
	//fmt.Fprintln(w, "")
}

func getTabCount(count int) string {
	tabs := make([]string, count)
	for i := 0; i < count; i++ {
		tabs[i] = "\t"
	}
	return strings.Join(tabs, "")
}
