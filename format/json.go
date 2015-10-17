package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/mkboudreau/sitemap/domain"
)

type JsonSiteFormatter struct{}

func (f *JsonSiteFormatter) Format(w io.Writer, site *domain.Site) {
	if err := json.NewEncoder(w).Encode(site); err != nil {
		fmt.Fprintf(w, "formatter failed to encode as json: %v", err)
	}
}
