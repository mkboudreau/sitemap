package builder

import (
	"log"
	"net/url"
	"strings"
)

func translateLink(sitelink string, newlink string) string {
	//log.Println("... site link:", sitelink, "; new link:", newlink)
	newURL, newErr := normalizeToURL(newlink)
	if newErr != nil {
		log.Printf("Could not parse link %v: %v", newlink, newErr)
		return ""
	}
	siteURL, siteErr := normalizeToURL(sitelink)
	//log.Println("... normalized site link:", siteURL.String(), "; new link:", newURL.String())

	normalizeSchemeFromURL(newURL, siteURL)
	normalizeHostFromURL(newURL, siteURL)

	//log.Println("... more normalized site link:", siteURL.String(), "; new link:", newURL.String())

	if !schemeIsHttpCompatible(newURL.Scheme) {
		return ""
	}

	if siteErr != nil || siteURL == nil || siteURL.String() == "" {
		log.Printf("(ignore if this is top url) could not translate site url=%v; err=%v; returning new url=%v", siteURL, siteErr, newURL)
		return newURL.String()
	}

	if siteURL.String() == newURL.String() {
		return ""
	}

	resolvedURL := siteURL.ResolveReference(newURL)

	if resolvedURL.String() == siteURL.String() {
		return ""
	}

	return resolvedURL.String()
}

func normalizeSchemeFromURL(target *url.URL, source *url.URL) {
	if target.Scheme == "" {
		if source != nil && source.Scheme != "" && strings.HasPrefix(target.Path, "/") {
			target.Scheme = source.Scheme
		} else if target.Host != "" {
			target.Scheme = "http"
		}
	}
}
func normalizeHostFromURL(target *url.URL, source *url.URL) {
	if target.Host == "" {
		if source != nil && source.Host != "" && strings.HasPrefix(target.Path, "/") {
			target.Host = source.Host
		} else {
			//target.Host = "localhost"
		}
	}
}

func validateLinkInSameDomain(sitelink string, newlink string) bool {
	newURL, newErr := normalizeToURL(newlink)
	if newErr != nil {
		return false
	}
	siteURL, siteErr := normalizeToURL(sitelink)
	if siteErr != nil {
		return false
	}

	if siteURL.Host == newURL.Host {
		return true
	} else if strings.Contains(siteURL.Host, newURL.Host) {
		return true
	} else if strings.Contains(newURL.Host, siteURL.Host) {
		return false
	}
	return false
}

func normalizeToURL(link string) (*url.URL, error) {
	newURL, err := url.Parse(link)
	if err != nil {
		return nil, err
	}
	if newURL.Scheme == "" && newURL.Host == "" {
		fixHost(newURL)
	}
	return newURL, nil
}

func fixHost(url *url.URL) {
	if parts := strings.SplitN(url.Path, "/", 2); len(parts) > 1 {
		if len(strings.Split(parts[0], ".")) >= 2 && parts[0] != "" && parts[1] != "" {
			url.Host = parts[0]
			url.Path = parts[1]
		}
	} else {
		dotParts := strings.Split(url.Path, ".")
		if len(dotParts) == 0 {
			return
		} else if len(dotParts) > 2 || suffixlooksLikeHost(url.Path) {
			url.Host = url.Path
			url.Path = ""
		}
	}
}

func suffixlooksLikeHost(path string) bool {
	if strings.HasSuffix(path, ".com") ||
		strings.HasSuffix(path, ".org") ||
		strings.HasSuffix(path, ".net") ||
		strings.HasSuffix(path, ".gov") ||
		strings.HasSuffix(path, ".edu") ||
		strings.HasSuffix(path, ".us") ||
		strings.HasSuffix(path, ".io") {
		return true
	}
	return false
}

func schemeIsHttpCompatible(scheme string) bool {
	return scheme == "" || scheme == "http" || scheme == "https"
}
