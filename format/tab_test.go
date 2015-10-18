package format

import "testing"

func TestTabFormatter(t *testing.T) {
	formatTestRunner(t, tabFormatTestCases)
}

var tabFormatTestCases = []formatTestCase{
	{&TabbedSiteFormatter{}, TestSiteOneLink, tabResultOneLink},
	{&TabbedSiteFormatter{}, TestSiteTwoLinks, tabResultTwoLinks},
	{&TabbedSiteFormatter{}, TestSiteOneAsset, tabResultOneAsset},
	{&TabbedSiteFormatter{}, TestSiteTwoAssets, tabResultTwoAssets},
	{&TabbedSiteFormatter{}, TestSiteOneEach, tabResultOneEach},
	{&TabbedSiteFormatter{}, TestSiteTwoEach, tabResultTwoEach},
	{&TabbedSiteFormatter{}, TestSiteNested, tabResultNested},
}

var tabResultOneAsset = `URL:   	http://abc.com
Assets:	http://abc.com/123
`
var tabResultTwoAssets = `URL:   	http://abc.com
Assets:	http://abc.com/123, http://abc.com/456
`
var tabResultOneLink = `URL:   	http://abc.com
Links: 
	URL:   	http://abc.com/123

`
var tabResultTwoLinks = `URL:   	http://abc.com
Links: 
	URL:   	http://abc.com/123
	URL:   	http://abc.com/456

`
var tabResultOneEach = `URL:   	http://abc.com
Assets:	http://abc.com/123
Links: 
	URL:   	http://abc.com/456

`
var tabResultTwoEach = `URL:   	http://abc.com
Assets:	http://abc.com/123, http://abc.com/456
Links: 
	URL:   	http://abc.com/123
	URL:   	http://abc.com/456

`
var tabResultNested = `URL:   	http://abc.com
Links: 
	URL:   	http://abc.com/123
	Links: 
		URL:   	http://abc.com/456


`
