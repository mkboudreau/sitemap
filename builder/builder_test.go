package builder

import (
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildInterrupt(t *testing.T) {
	testServer := serverWithData(t, &builderTestTemplateData{})
	defer testServer.Close()

	builder := NewSitemapBuilder(1*time.Second, 5*time.Second, 1)
	builder.Interrupt()

	url := testServer.URL
	sitemap := builder.Build(url)

	assert.NotNil(t, sitemap)
	assert.NotNil(t, sitemap.Top)
	assert.EqualValues(t, url, sitemap.Top.Url)
	assert.Empty(t, sitemap.Top.Assets)
	assert.Empty(t, sitemap.Top.Links)
}

func TestHttpNoLinks(t *testing.T) {
	testServer := serverWithData(t, &builderTestTemplateData{})
	defer testServer.Close()

	builder := NewSitemapBuilder(0, 0, 1)

	url := testServer.URL
	sitemap := builder.Build(url)

	assert.NotNil(t, sitemap)
	assert.NotNil(t, sitemap.Top)
	assert.EqualValues(t, url, sitemap.Top.Url)
	assert.Empty(t, sitemap.Top.Assets)
	assert.Empty(t, sitemap.Top.Links)
}
func TestHttpOneLink(t *testing.T) {
	testServer := serverWithData(t, &builderTestTemplateData{Links: []string{"/link.html"}})
	defer testServer.Close()

	builder := NewSitemapBuilder(0, 0, 1)

	url := testServer.URL
	sitemap := builder.Build(url)

	assert.NotNil(t, sitemap)
	assert.NotNil(t, sitemap.Top)
	assert.EqualValues(t, url, sitemap.Top.Url)
	assert.Empty(t, sitemap.Top.Assets)
	assert.Len(t, sitemap.Top.Links, 1)
}
func TestHttpTwoAssets(t *testing.T) {
	testServer := serverWithData(t, &builderTestTemplateData{Styles: []string{"/assetone.css", "/assettwo.css"}})
	defer testServer.Close()

	builder := NewSitemapBuilder(0, 0, 1)

	url := testServer.URL
	sitemap := builder.Build(url)

	assert.NotNil(t, sitemap)
	assert.NotNil(t, sitemap.Top)
	assert.EqualValues(t, url, sitemap.Top.Url)
	assert.Len(t, sitemap.Top.Assets, 2)
	assert.Empty(t, sitemap.Top.Links)
}
func TestHttpOneAssetOneLink(t *testing.T) {
	testServer := serverWithData(t, &builderTestTemplateData{Styles: []string{"/assetone.css"}, Links: []string{"/link.html"}})
	defer testServer.Close()

	builder := NewSitemapBuilder(0, 0, 1)

	url := testServer.URL
	sitemap := builder.Build(url)

	assert.NotNil(t, sitemap)
	assert.NotNil(t, sitemap.Top)
	assert.EqualValues(t, url, sitemap.Top.Url)
	assert.Len(t, sitemap.Top.Assets, 1)
	assert.Len(t, sitemap.Top.Links, 1)
}

func serverWithData(t *testing.T, testdata *builderTestTemplateData) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("builder").Parse(builderTestTemplate)
		if err != nil {
			fmt.Fprintf(w, "ERROR: %v", err)
		}

		if strings.Contains(r.URL.Path, "link.html") {
			if err := tmpl.Execute(w, &builderTestTemplateData{}); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			if err := tmpl.Execute(w, testdata); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}))
}

var builderTestTemplate = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>My Title</title>
		{{range .Styles}}<link href="{{.}}" />{{end}}
	</head>
	<body>
		<h1>HELLO</h1>
		<div>
			<p>Hello World</p>
			<p>{{range .Links}}<a href="{{.}}">some link</a>{{end}}</p>
			<p>{{range .Images}}<img src="{{.}}" />{{end}}</p>
		</div>
	</body>
</html>`

type builderTestTemplateData struct {
	Links  []string
	Styles []string
	Images []string
}
