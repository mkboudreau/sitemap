package format

import (
	"bytes"
	"testing"

	"github.com/mkboudreau/sitemap/domain"
)

var (
	TestSiteOneLink   *domain.Site = newTestSiteLinksOnly("http://abc.com", "http://abc.com/123")
	TestSiteTwoLinks               = newTestSiteLinksOnly("http://abc.com", "http://abc.com/123", "http://abc.com/456")
	TestSiteOneAsset               = newTestSiteAssetsOnly("http://abc.com", "http://abc.com/123")
	TestSiteTwoAssets              = newTestSiteAssetsOnly("http://abc.com", "http://abc.com/123", "http://abc.com/456")
	TestSiteOneEach                = newTestSite("http://abc.com", []string{"http://abc.com/123"}, "http://abc.com/456")
	TestSiteTwoEach                = newTestSite("http://abc.com", []string{"http://abc.com/123", "http://abc.com/456"}, "http://abc.com/123", "http://abc.com/456")
	TestSiteNested                 = newTestSiteWithSite("http://abc.com", newTestSiteLinksOnly("http://abc.com/123", "http://abc.com/456"))
)

type formatTestCase struct {
	formatter   SiteFormatter
	site        *domain.Site
	expectation string
}

func formatTestRunner(t *testing.T, testcases []formatTestCase) {
	for _, testcase := range testcases {
		b := bytes.NewBuffer(nil)
		testcase.formatter.Format(b, testcase.site)
		actual := b.String()
		if actual != testcase.expectation {
			t.Logf("Format result for formatter does not match expectation:\n[%v]\nactual:\n[%v]\n", testcase.expectation, actual)
			t.Fail()
		}
	}
}

func newTestSiteLinksOnly(url string, links ...string) *domain.Site {
	sites := make([]*domain.Site, len(links))
	for i, l := range links {
		sites[i] = domain.NewSite(l)
	}
	return &domain.Site{
		Url:   url,
		Links: sites,
	}
}

func newTestSiteWithSite(url string, sites ...*domain.Site) *domain.Site {
	return &domain.Site{
		Url:   url,
		Links: sites,
	}
}

func newTestSiteAssetsOnly(url string, assets ...string) *domain.Site {
	return &domain.Site{
		Url:    url,
		Assets: assets,
	}
}

func newTestSite(url string, assets []string, links ...string) *domain.Site {
	site := newTestSiteLinksOnly(url, links...)
	site.Assets = assets
	return site
}
