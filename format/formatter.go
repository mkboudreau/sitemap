package format

import (
	"io"

	"github.com/mkboudreau/sitemap/domain"
)

type SiteFormatter interface {
	Format(w io.Writer, site *domain.Site)
}
