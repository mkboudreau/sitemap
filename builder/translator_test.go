package builder

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var translationSameSiteTestCases = []struct {
	site, link, expected string
	valid                bool
}{
	{"", "/link.html", "http://localhost/link.html", true},
	{"http://www.site.com", "/link.html", "http://www.site.com/link.html", true},
	{"http://www.site.com/index.html", "link.html", "http://www.site.com/link.html", true},
	{"http://www.site.com/index.html", "sub/link.html", "http://www.site.com/sub/link.html", true},
	{"http://www.site.com/sub/index.html", "again/link.html", "http://www.site.com/sub/again/link.html", true},
	{"http://www.site.com/sub/index.html", "/fake/link.html", "http://www.site.com/fake/link.html", true},
	{"/index.html", "link.html", "http://localhost/link.html", true},
	{"site.com/index.html", "link.html", "http://site.com/link.html", true},
	{"www.site.com/index.html", "link.html", "http://www.site.com/link.html", true},
	{"http://www.site.com/sub/index.html", "link.html", "http://www.site.com/sub/link.html", true},
	{"http://www.site.com/sub/index.html", "/link.html", "http://www.site.com/link.html", true},
}
var translationDifferentSiteTestCases = []struct {
	site, link, expected string
	valid                bool
}{
	{"", "site.com/link.html", "http://site.com/link.html", false},
	{"", "https://site.com/link.html", "https://site.com/link.html", false},
	{"", "www.yahoo.com", "http://www.yahoo.com", false},
	{"http://www.site.com", "another.com/link.html", "http://another.com/link.html", false},
	{"http://www.site.com/index.html", "http://another.com/link.html", "http://another.com/link.html", false},
	{"site.com/index.html", "http://another.com/link.html", "http://another.com/link.html", false},
	{"site.com/index.html", "sub.site.com/link.html", "http://sub.site.com/link.html", false},
	{"sub.site.com/index.html", "site.com/link.html", "http://site.com/link.html", true},
}

func TestSameSiteTranslator(t *testing.T) {
	for i, testcase := range translationSameSiteTestCases {
		actual := translateLink(testcase.site, testcase.link)
		assert.EqualValues(t, testcase.expected, actual, fmt.Sprintf("Test Case %d", i))
	}
}

func TestDifferentSiteTranslator(t *testing.T) {
	for i, testcase := range translationDifferentSiteTestCases {
		actual := translateLink(testcase.site, testcase.link)
		assert.EqualValues(t, testcase.expected, actual, fmt.Sprintf("Test Case %d", i))
	}
}

func TestSameSiteValidator(t *testing.T) {
	for i, testcase := range translationSameSiteTestCases {
		actual := validateLinkInSameDomain(testcase.site, testcase.link)
		assert.Equal(t, testcase.valid, actual, fmt.Sprintf("Test Case %d", i))
	}
}

func TestDifferentSiteValidator(t *testing.T) {
	for i, testcase := range translationDifferentSiteTestCases {
		actual := validateLinkInSameDomain(testcase.site, testcase.link)
		assert.EqualValues(t, testcase.valid, actual, fmt.Sprintf("Test Case %d", i))
	}
}
