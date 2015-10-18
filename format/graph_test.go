package format

import "testing"

func TestGraphFormatter(t *testing.T) {
	formatTestRunner(t, graphFormatTestCases)
}

var graphFormatTestCases = []formatTestCase{
	{&GraphSiteFormatter{}, TestSiteOneLink, graphResultOneLink},
	{&GraphSiteFormatter{}, TestSiteTwoLinks, graphResultTwoLinks},
	{&GraphSiteFormatter{}, TestSiteOneAsset, graphResultOneAsset},
	{&GraphSiteFormatter{}, TestSiteTwoAssets, graphResultTwoAssets},
	{&GraphSiteFormatter{}, TestSiteOneEach, graphResultOneEach},
	{&GraphSiteFormatter{}, TestSiteTwoEach, graphResultTwoEach},
	{&GraphSiteFormatter{}, TestSiteNested, graphResultNested},
}

var graphResultOneAsset = `graph {
		"http://abc.com" -- "http://abc.com/123";
}
`
var graphResultTwoAssets = `graph {
		"http://abc.com" -- "http://abc.com/123";
		"http://abc.com" -- "http://abc.com/456";
}
`
var graphResultOneLink = `graph {
		"http://abc.com" -- "http://abc.com/123";
}
`
var graphResultTwoLinks = `graph {
		"http://abc.com" -- "http://abc.com/123";
		"http://abc.com" -- "http://abc.com/456";
}
`
var graphResultOneEach = `graph {
		"http://abc.com" -- "http://abc.com/123";
		"http://abc.com" -- "http://abc.com/456";
}
`
var graphResultTwoEach = `graph {
		"http://abc.com" -- "http://abc.com/123";
		"http://abc.com" -- "http://abc.com/456";
}
`
var graphResultNested = `graph {
		"http://abc.com" -- "http://abc.com/123";
		"http://abc.com/123" -- "http://abc.com/456";
}
`
