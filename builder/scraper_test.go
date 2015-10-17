package builder

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrapeOneLowerAnchor(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(oneLowercaseAnchorTagInHTML)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer testServer.Close()

	page, err := scrapePage(testServer.URL)

	assert.NoError(t, err)
	assert.NotNil(t, page.Links)
	assert.Len(t, page.Links, 1)
	assert.EqualValues(t, "abc.html", page.Links[0])
}
func TestScrapeOneUpperAnchor(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(oneUppercaseAnchorTagInHTML)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer testServer.Close()

	page, err := scrapePage(testServer.URL)

	assert.NoError(t, err)
	assert.NotNil(t, page.Links)
	assert.Len(t, page.Links, 1)
	assert.EqualValues(t, "abc.html", page.Links[0])
}
func TestScrapeTwoMixedCaseAnchors(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(mixedCaseAnchorTagInHTML)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer testServer.Close()

	page, err := scrapePage(testServer.URL)

	assert.NoError(t, err)
	assert.NotNil(t, page.Links)
	assert.Len(t, page.Links, 2)
}

func TestScrapeTwoDuplicateAnchors(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(duplicateAnchorTagInHTML)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer testServer.Close()

	page, err := scrapePage(testServer.URL)

	assert.NoError(t, err)
	assert.NotNil(t, page.Links)
	assert.Len(t, page.Links, 2)
}

var oneLowercaseAnchorTagInHTML = `
<html><head><title>test</title></head>
<body>
<h1>heading 1</h1>
<br />
<h2>heading 2</h2><br />
<p>hello world <img src="abc.gif"></p><br />
<p>hello anchor <a href="abc.html" /></p>
<br /><br /><br />
<p></p>
</body>
</html>`

var oneUppercaseAnchorTagInHTML = `
<html><head><title>test</title></head>
<body>
<h1>heading 1</h1>
<br />
<h2>heading 2</h2><br />
<p>hello world <img src="abc.gif"></p><br />
<p>hello anchor <A href="abc.html">abc</A></p>
<br /><br /><br />
<p></p>
</body>
</html>`

var mixedCaseAnchorTagInHTML = `
<html><head><title>test</title></head>
<body>
<h1>heading 1</h1>
<br />
<h2>heading 2</h2><br />
<p>hello world <img src="abc.gif"></p><br />
<p>hello anchor one<A href="abc.html">abc</A></p>
<p>hello anchor two<a href="xyz.html">xyz</a></p>
<br /><br /><br />
<p></p>
</body>
</html>`

var duplicateAnchorTagInHTML = `
<html><head><title>test</title></head>
<body>
<h1>heading 1</h1>
<br />
<h2>heading 2</h2><br />
<p>hello world <img src="abc.gif"></p><br />
<p>hello anchor one<A href="abc.html">abc</A></p>
<p>hello anchor one<A href="abc.html">abc</A></p>
<p>hello anchor one<a href="abc.html">abc</a></p>
<p>hello anchor two<A href="xyz.html">xyz</A></p>
<p>hello anchor two<a href="xyz.html">xyz</a></p>
<p>hello anchor two<a href="xyz.html">xyz</a></p>
<p>hello anchor one<A href="abc.html" /></p>
<p>hello anchor two<a href="xyz.html" /></p>
<p>hello anchor two<a href="xyz.html" /></p>
<p>hello anchor two<a href="xyz.html" /></p>
<br /><br /><br />
<p></p>
</body>
</html>`
