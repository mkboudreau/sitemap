package format

import "testing"

func TestJsonFormatter(t *testing.T) {
	formatTestRunner(t, jsonFormatTestCases)
}

var jsonFormatTestCases = []formatTestCase{
	{&JsonSiteFormatter{}, TestSiteOneLink, jsonResultOneLink},
	{&JsonSiteFormatter{}, TestSiteTwoLinks, jsonResultTwoLinks},
	{&JsonSiteFormatter{}, TestSiteOneAsset, jsonResultOneAsset},
	{&JsonSiteFormatter{}, TestSiteTwoAssets, jsonResultTwoAssets},
	{&JsonSiteFormatter{}, TestSiteOneEach, jsonResultOneEach},
	{&JsonSiteFormatter{}, TestSiteTwoEach, jsonResultTwoEach},
	{&JsonSiteFormatter{}, TestSiteNested, jsonResultNested},
}

var jsonResultOneAsset = `{"url":"http://abc.com","assets":["http://abc.com/123"]}
`
var jsonResultTwoAssets = `{"url":"http://abc.com","assets":["http://abc.com/123","http://abc.com/456"]}
`
var jsonResultOneLink = `{"url":"http://abc.com","links":[{"url":"http://abc.com/123"}]}
`
var jsonResultTwoLinks = `{"url":"http://abc.com","links":[{"url":"http://abc.com/123"},{"url":"http://abc.com/456"}]}
`
var jsonResultOneEach = `{"url":"http://abc.com","assets":["http://abc.com/123"],"links":[{"url":"http://abc.com/456"}]}
`
var jsonResultTwoEach = `{"url":"http://abc.com","assets":["http://abc.com/123","http://abc.com/456"],"links":[{"url":"http://abc.com/123"},{"url":"http://abc.com/456"}]}
`
var jsonResultNested = `{"url":"http://abc.com","links":[{"url":"http://abc.com/123","links":[{"url":"http://abc.com/456"}]}]}
`
