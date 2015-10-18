package format

import "testing"

func TestDigraphFormatter(t *testing.T) {
	formatTestRunner(t, digraphFormatTestCases)
}

var digraphFormatTestCases = []formatTestCase{
	{&DigraphSiteFormatter{}, TestSiteOneLink, digraphResultOneLink},
	{&DigraphSiteFormatter{}, TestSiteTwoLinks, digraphResultTwoLinks},
	{&DigraphSiteFormatter{}, TestSiteOneAsset, digraphResultOneAsset},
	{&DigraphSiteFormatter{}, TestSiteTwoAssets, digraphResultTwoAssets},
	{&DigraphSiteFormatter{}, TestSiteOneEach, digraphResultOneEach},
	{&DigraphSiteFormatter{}, TestSiteTwoEach, digraphResultTwoEach},
	{&DigraphSiteFormatter{}, TestSiteNested, digraphResultNested},
}

var digraphResultOneAsset = `digraph {
	subgraph cluster_0 {
		label="Assets";
		"http://abc.com" -> "http://abc.com/123";
	}

	subgraph cluster_1 {
		label="Links";
	}
}
`
var digraphResultTwoAssets = `digraph {
	subgraph cluster_0 {
		label="Assets";
		"http://abc.com" -> "http://abc.com/123";
		"http://abc.com" -> "http://abc.com/456";
	}

	subgraph cluster_1 {
		label="Links";
	}
}
`
var digraphResultOneLink = `digraph {
	subgraph cluster_0 {
		label="Assets";
	}

	subgraph cluster_1 {
		label="Links";
		"http://abc.com" -> "http://abc.com/123";
	}
}
`
var digraphResultTwoLinks = `digraph {
	subgraph cluster_0 {
		label="Assets";
	}

	subgraph cluster_1 {
		label="Links";
		"http://abc.com" -> "http://abc.com/123";
		"http://abc.com" -> "http://abc.com/456";
	}
}
`
var digraphResultOneEach = `digraph {
	subgraph cluster_0 {
		label="Assets";
		"http://abc.com" -> "http://abc.com/123";
	}

	subgraph cluster_1 {
		label="Links";
		"http://abc.com" -> "http://abc.com/456";
	}
}
`
var digraphResultTwoEach = `digraph {
	subgraph cluster_0 {
		label="Assets";
		"http://abc.com" -> "http://abc.com/123";
		"http://abc.com" -> "http://abc.com/456";
	}

	subgraph cluster_1 {
		label="Links";
		"http://abc.com" -> "http://abc.com/123";
		"http://abc.com" -> "http://abc.com/456";
	}
}
`
var digraphResultNested = `digraph {
	subgraph cluster_0 {
		label="Assets";
	}

	subgraph cluster_1 {
		label="Links";
		"http://abc.com" -> "http://abc.com/123";
		"http://abc.com/123" -> "http://abc.com/456";
	}
}
`
